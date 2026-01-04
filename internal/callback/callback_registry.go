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
type Registry struct {
	mutex     sync.RWMutex
	callbacks map[ClientID]map[RequestID]func(confirmation ocpp.Response, err error)
}

// New creates a new CallbackRegistry instance.
func New() *Registry {
	return &Registry{
		callbacks: make(map[ClientID]map[RequestID]func(confirmation ocpp.Response, err error)),
	}
}

// RegisterCallback registers a callback for a specific client ID and request ID.
// If the try function returns an error, the callback is not registered and the error is returned.
// If the try function succeeds, the callback is registered and can be retrieved later using GetCallback.
func (cr *Registry) RegisterCallback(clientID string, try func() (string, error), callback func(confirmation ocpp.Response, err error)) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	// Initialize the inner map if it doesn't exist
	if cr.callbacks[ClientID(clientID)] == nil {
		cr.callbacks[ClientID(clientID)] = make(map[RequestID]func(confirmation ocpp.Response, err error))
	}

	requestId, err := try()
	if err != nil {
		// Clean up empty client maps
		if len(cr.callbacks[ClientID(clientID)]) == 0 {
			delete(cr.callbacks, ClientID(clientID))
		}
		return err
	}

	// Register the callback
	cr.callbacks[ClientID(clientID)][RequestID(requestId)] = callback

	return nil
}

// GetCallback retrieves and removes the callback for a specific client ID and request ID.
// Returns the callback and true if found, nil and false otherwise.
func (cr *Registry) GetCallback(clientID string, requestID string) (func(confirmation ocpp.Response, err error), bool) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	clientCallbacks, ok := cr.callbacks[ClientID(clientID)]
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
		delete(cr.callbacks, ClientID(clientID))
	}

	return callback, true
}

// RemoveCallback removes a callback for a specific client ID and request ID without invoking it.
// Returns true if the callback was found and removed, false otherwise.
func (cr *Registry) RemoveCallback(clientID string, requestID string) bool {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	clientCallbacks, ok := cr.callbacks[ClientID(clientID)]
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
		delete(cr.callbacks, ClientID(clientID))
	}

	return true
}

// GetAllCallbacks retrieves and removes all callbacks for a specific client ID.
// Returns a map of request ID to callback function, and true if any callbacks were found.
// The returned map will be empty if no callbacks exist for the client.
func (cr *Registry) GetAllCallbacks(clientID string) (map[RequestID]func(confirmation ocpp.Response, err error), bool) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	clientId := ClientID(clientID)
	clientCallbacks, ok := cr.callbacks[clientId]
	if !ok || len(clientCallbacks) == 0 {
		return make(map[RequestID]func(confirmation ocpp.Response, err error)), false
	}

	// Copy all callbacks to the result map
	result := make(map[RequestID]func(confirmation ocpp.Response, err error), len(clientCallbacks))
	for requestID, callback := range clientCallbacks {
		result[requestID] = callback
	}

	// Remove all callbacks for this client
	delete(cr.callbacks, clientId)
	return result, true
}

// ClearCallbacks removes all callbacks for a specific client ID.
// Returns the number of callbacks that were removed.
func (cr *Registry) ClearCallbacks(clientID string) int {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	clientCallbacks, ok := cr.callbacks[ClientID(clientID)]
	if !ok {
		return 0
	}

	count := len(clientCallbacks)
	delete(cr.callbacks, ClientID(clientID))
	return count
}
