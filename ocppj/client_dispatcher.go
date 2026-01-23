package ocppj

import (
	"fmt"
	"sync"
	"time"

	"github.com/xBlaz3kx/ocpp-go/logging"
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ws"
)

// DefaultClientDispatcher is a default implementation of the ClientDispatcher interface.
//
// The dispatcher implements the ClientState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultClientDispatcher struct {
	logger              logging.Logger
	requestQueue        RequestQueue
	requestChannel      chan bool
	readyForDispatch    chan bool
	pendingRequestState ClientState
	network             ws.Client
	mutex               sync.RWMutex
	onRequestCancel     func(requestID string, request ocpp.Request, err *ocpp.Error)
	timer               *time.Timer
	state               string
	timeout             time.Duration
}

const (
	defaultTimeoutTick    = 24 * time.Hour
	defaultMessageTimeout = 30 * time.Second
)

// NewDefaultClientDispatcher creates a new DefaultClientDispatcher struct.
// If logger is nil, a VoidLogger will be used.
func NewDefaultClientDispatcher(queue RequestQueue, logger logging.Logger) *DefaultClientDispatcher {
	if logger == nil {
		logger = &logging.VoidLogger{}
	}

	return &DefaultClientDispatcher{
		logger:              logger,
		requestQueue:        queue,
		requestChannel:      make(chan bool, 1),
		readyForDispatch:    make(chan bool, 1),
		pendingRequestState: NewClientState(),
		state:               "stopped",
		timeout:             defaultMessageTimeout,
	}
}

func (d *DefaultClientDispatcher) SetOnRequestCanceled(cb func(requestID string, request ocpp.Request, err *ocpp.Error)) {
	d.onRequestCancel = cb
}

func (d *DefaultClientDispatcher) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *DefaultClientDispatcher) Start() {
	d.mutex.Lock()
	d.state = "running"
	d.requestChannel = make(chan bool, 1)
	d.timer = time.NewTimer(defaultTimeoutTick) // Default to 24 hours tick
	d.mutex.Unlock()

	go d.messagePump()
}

func (d *DefaultClientDispatcher) IsRunning() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.requestChannel != nil
}

func (d *DefaultClientDispatcher) IsPaused() bool {
	d.mutex.Lock()
	defer d.mutex.Lock()

	return d.state == "paused"
}

func (d *DefaultClientDispatcher) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	close(d.requestChannel)
	// TODO: clear pending requests?
}

func (d *DefaultClientDispatcher) SetNetworkClient(client ws.Client) {
	d.network = client
}

func (d *DefaultClientDispatcher) SetPendingRequestState(state ClientState) {
	d.pendingRequestState = state
}

func (d *DefaultClientDispatcher) SendRequest(req RequestBundle) error {
	if d.network == nil {
		return fmt.Errorf("cannot SendRequest, no network client was set")
	}
	if err := d.requestQueue.Push(req); err != nil {
		return err
	}

	d.requestChannel <- true
	return nil
}

func (d *DefaultClientDispatcher) messagePump() {
	rdy := true // Ready to transmit at the beginning

	for {
		select {
		case _, ok := <-d.requestChannel:
			// New request was posted
			if !ok {
				// Channel closed, stopping dispatcher
				d.requestQueue.Init()
				return
			}
		case _, ok := <-d.timer.C:
			// Timeout elapsed
			if !ok {
				continue
			}
			if d.pendingRequestState.HasPendingRequest() {
				// Current request timed out. Removing request and triggering cancel callback
				el := d.requestQueue.Peek()
				bundle, _ := el.(RequestBundle)
				d.CompleteRequest(bundle.Call.UniqueId)
				if d.onRequestCancel != nil {
					d.onRequestCancel(bundle.Call.UniqueId, bundle.Call.Payload,
						ocpp.NewError(GenericError, "Request timed out", bundle.Call.UniqueId))
				}
			}
			// No request is currently pending -> set timer to high number
			d.timer.Reset(defaultTimeoutTick)
		case rdy = <-d.readyForDispatch:
			// Ready flag set, keep going
		}

		// Check if dispatcher is paused
		if d.IsPaused() {
			// Ignore dispatch events as long as dispatcher is paused
			continue
		}

		// Only dispatch request if able to send and request queue isn't empty
		if rdy && !d.requestQueue.IsEmpty() {
			d.dispatchNextRequest()
			rdy = false
			// Set timer
			if !d.timer.Stop() {
				<-d.timer.C
			}
			d.timer.Reset(d.timeout)
		}
	}
}

func (d *DefaultClientDispatcher) dispatchNextRequest() {
	// Get first element in queue
	el := d.requestQueue.Peek()

	bundle, canCast := el.(RequestBundle)
	if !canCast {
		d.logger.Errorf("failed to cast request queue element to RequestBundle")
		return
	}

	if bundle.Call == nil {
		d.logger.Errorf("request bundle has no Call associated")
		return
	}

	if bundle.Data == nil {
		d.logger.Errorf("request bundle has no Data associated")
		return
	}

	jsonMessage := bundle.Data
	d.pendingRequestState.AddPendingRequest(bundle.Call.UniqueId, bundle.Call.Payload)
	// Attempt to send over network
	err := d.network.Write(jsonMessage)
	if err != nil {
		// TODO: handle retransmission instead of skipping request altogether
		d.CompleteRequest(bundle.Call.GetUniqueId())
		if d.onRequestCancel != nil {
			d.onRequestCancel(bundle.Call.UniqueId, bundle.Call.Payload,
				ocpp.NewError(InternalError, err.Error(), bundle.Call.UniqueId))
		}
	}

	d.logger.Infof("dispatched request %s to server", bundle.Call.UniqueId)
	d.logger.Debugf("sent JSON message to server: %s", string(jsonMessage))
}

func (d *DefaultClientDispatcher) Pause() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.state = "paused"

	if !d.timer.Stop() {
		<-d.timer.C
	}

	d.timer.Reset(defaultTimeoutTick)

}

func (d *DefaultClientDispatcher) Resume() {
	d.mutex.Lock()
	d.state = "running"
	d.mutex.Unlock()

	if d.pendingRequestState.HasPendingRequest() {
		// There is a pending request already. Awaiting response, before dispatching new requests.
		d.timer.Reset(d.timeout)
	} else {
		// Can dispatch a new request. Notifying message pump.
		d.readyForDispatch <- true
	}
}

func (d *DefaultClientDispatcher) CompleteRequest(requestId string) {
	el := d.requestQueue.Peek()
	if el == nil {
		d.logger.Errorf("attempting to pop front of queue, but queue is empty")
		return
	}

	bundle, _ := el.(RequestBundle)
	if bundle.Call.UniqueId != requestId {
		d.logger.Errorf("internal state mismatch: received response for %v but expected response for %v", requestId, bundle.Call.UniqueId)
		return
	}

	d.requestQueue.Pop()
	d.pendingRequestState.DeletePendingRequest(requestId)
	d.logger.Debugf("removed request %v from front of queue", bundle.Call.UniqueId)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- true
}
