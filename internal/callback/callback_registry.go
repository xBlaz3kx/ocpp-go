package callback

import (
	"sync"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
)

// ClientID represents a client identifier
type ClientID string

// RequestID represents a request/message identifier
type RequestID string

// CallbackRegistry registers and retrieves callbacks based on client ID and request ID.
// This allows callbacks to be registered for specific requests and retrieved when responses arrive.
// The registry uses per-client locks to allow concurrent operations on different clients
// while preventing race conditions for individual client operations.
type Registry struct {
	globalMutex sync.RWMutex             // Protects the clientLocks and callbacks maps
	clientLocks map[ClientID]*sync.Mutex // Per-client locks for concurrent send operations
	callbacks   map[ClientID]map[RequestID]func(confirmation ocpp.Response, err error)
}

// New creates a new CallbackRegistry instance.
func New() *Registry {
	return &Registry{
		clientLocks: make(map[ClientID]*sync.Mutex),
		callbacks:   make(map[ClientID]map[RequestID]func(confirmation ocpp.Response, err error)),
	}
}

// getClientLock returns the lock for a specific client, creating one if it doesn't exist.
// This allows concurrent sends to different clients while serializing operations per-client.
func (cr *Registry) getClientLock(clientID ClientID) *sync.Mutex {
	cr.globalMutex.RLock()
	lock, exists := cr.clientLocks[clientID]
	cr.globalMutex.RUnlock()

	if exists {
		return lock
	}

	cr.globalMutex.Lock()
	defer cr.globalMutex.Unlock()
	// Double-check after acquiring write lock
	if lock, exists = cr.clientLocks[clientID]; exists {
		return lock
	}
	lock = &sync.Mutex{}
	cr.clientLocks[clientID] = lock
	return lock
}

// RegisterCallback registers a callback for a specific client ID and request ID.
// If the try function returns an error, the callback is not registered and the error is returned.
// If the try function succeeds, the callback is registered and can be retrieved later using GetCallback.
// The per-client lock ensures that the callback is registered before any response can be looked up,
// while still allowing concurrent sends to different clients.
func (cr *Registry) RegisterCallback(clientID string, try func() (string, error), callback func(confirmation ocpp.Response, err error)) error {
	cID := ClientID(clientID)

	// Get per-client lock - this serializes operations for the same client
	// but allows concurrent operations on different clients
	clientLock := cr.getClientLock(cID)
	clientLock.Lock()
	defer clientLock.Unlock()

	// Perform the try operation (typically SendRequest) while holding the client lock
	// This ensures the callback is registered before any response can be looked up
	requestId, err := try()
	if err != nil {
		return err
	}

	// Now acquire global lock only for map access
	cr.globalMutex.Lock()
	defer cr.globalMutex.Unlock()

	// Initialize the inner map if it doesn't exist
	if cr.callbacks[cID] == nil {
		cr.callbacks[cID] = make(map[RequestID]func(confirmation ocpp.Response, err error))
	}

	// Register the callback
	cr.callbacks[cID][RequestID(requestId)] = callback
	return nil
}

// GetCallback retrieves and removes the callback for a specific client ID and request ID.
// Returns the callback and true if found, nil and false otherwise.
// The per-client lock ensures this doesn't race with RegisterCallback for the same client.
func (cr *Registry) GetCallback(clientID string, requestID string) (func(confirmation ocpp.Response, err error), bool) {
	cID := ClientID(clientID)

	// Get per-client lock to ensure we don't race with RegisterCallback
	clientLock := cr.getClientLock(cID)
	clientLock.Lock()
	defer clientLock.Unlock()

	cr.globalMutex.Lock()
	defer cr.globalMutex.Unlock()

	clientCallbacks, ok := cr.callbacks[cID]
	if !ok {
		return nil, false
	}

	callback, ok := clientCallbacks[RequestID(requestID)]
	if !ok {
		return nil, false
	}

	// Remove the callback after retrieving it
	delete(clientCallbacks, RequestID(requestID))
	// Clean up empty client maps
	if len(clientCallbacks) == 0 {
		delete(cr.callbacks, cID)
	}

	return callback, true
}

// RemoveCallback removes a callback for a specific client ID and request ID without invoking it.
// Returns true if the callback was found and removed, false otherwise.
func (cr *Registry) RemoveCallback(clientID string, requestID string) bool {
	cID := ClientID(clientID)

	// Get per-client lock
	clientLock := cr.getClientLock(cID)
	clientLock.Lock()
	defer clientLock.Unlock()

	cr.globalMutex.Lock()
	defer cr.globalMutex.Unlock()

	clientCallbacks, ok := cr.callbacks[cID]
	if !ok {
		return false
	}

	_, ok = clientCallbacks[RequestID(requestID)]
	if !ok {
		return false
	}

	delete(clientCallbacks, RequestID(requestID))
	// Clean up empty client maps
	if len(clientCallbacks) == 0 {
		delete(cr.callbacks, cID)
	}

	return true
}

// GetAllCallbacks retrieves and removes all callbacks for a specific client ID.
// Returns a map of request ID to callback function, and true if any callbacks were found.
// The returned map will be empty if no callbacks exist for the client.
func (cr *Registry) GetAllCallbacks(clientID string) (map[RequestID]func(confirmation ocpp.Response, err error), bool) {
	cID := ClientID(clientID)

	// Get per-client lock
	clientLock := cr.getClientLock(cID)
	clientLock.Lock()
	defer clientLock.Unlock()

	cr.globalMutex.Lock()
	defer cr.globalMutex.Unlock()

	clientCallbacks, ok := cr.callbacks[cID]
	if !ok || len(clientCallbacks) == 0 {
		return make(map[RequestID]func(confirmation ocpp.Response, err error)), false
	}

	// Copy all callbacks to the result map
	result := make(map[RequestID]func(confirmation ocpp.Response, err error), len(clientCallbacks))
	for requestID, callback := range clientCallbacks {
		result[requestID] = callback
	}

	// Remove all callbacks for this client
	delete(cr.callbacks, cID)
	return result, true
}

// ClearCallbacks removes all callbacks for a specific client ID.
// Returns the number of callbacks that were removed.
func (cr *Registry) ClearCallbacks(clientID string) int {
	cID := ClientID(clientID)

	// Get per-client lock
	clientLock := cr.getClientLock(cID)
	clientLock.Lock()
	defer clientLock.Unlock()

	cr.globalMutex.Lock()
	defer cr.globalMutex.Unlock()

	clientCallbacks, ok := cr.callbacks[cID]
	if !ok {
		return 0
	}

	count := len(clientCallbacks)
	delete(cr.callbacks, cID)
	return count
}
