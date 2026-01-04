package ocppj

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xBlaz3kx/ocpp-go/logging"
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ws"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type DefaultServerDispatcherOption func(*DefaultServerDispatcher)

func WithServerDispatcherTimeout(timeout time.Duration) DefaultServerDispatcherOption {
	return func(d *DefaultServerDispatcher) {
		d.timeout = timeout
	}
}

func WithLogger(logger logging.Logger) DefaultServerDispatcherOption {
	return func(d *DefaultServerDispatcher) {
		d.logger = logger
	}
}

func WithMeterProvider(provider metric.MeterProvider) DefaultServerDispatcherOption {
	return func(d *DefaultServerDispatcher) {
		if provider == nil {
			return
		}

		metrics, err := newDispatcherMetrics(provider, d.logger)
		if err != nil {
			d.logger.Errorf("failed to create dispatcher metrics: %v", err)
			return
		}
		d.metrics = metrics
	}
}

// DefaultServerDispatcher is a default implementation of the ServerDispatcher interface.
//
// The dispatcher implements the ClientState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultServerDispatcher struct {
	logger              logging.Logger
	queueMap            ServerQueueMap
	requestChannel      chan string
	readyForDispatch    chan string
	pendingRequestState ServerState
	timeout             time.Duration
	timerC              chan string
	running             atomic.Bool
	stoppedC            chan struct{}
	onRequestCancel     CanceledRequestHandler
	network             ws.Server
	mutex               sync.RWMutex
	metrics             *dispatcherMetrics
}

// Handler function to be invoked when a request gets canceled (either due to timeout or to other external factors).
type CanceledRequestHandler func(clientID string, requestID string, request ocpp.Request, err *ocpp.Error)

// Utility struct for passing a client context around and cancel pending requests.
type clientTimeoutContext struct {
	ctx    context.Context
	cancel func()
}

func (c clientTimeoutContext) isActive() bool {
	return c.cancel != nil
}

// NewDefaultServerDispatcher creates a new DefaultServerDispatcher struct.
// If logger is nil, a VoidLogger will be used.
func NewDefaultServerDispatcher(
	queueMap ServerQueueMap,
	opts ...DefaultServerDispatcherOption,
) *DefaultServerDispatcher {
	logger := &logging.VoidLogger{}

	dispatcherMetrics, err := newDispatcherMetrics(otel.GetMeterProvider(), logger)
	if err != nil {
		logger.Errorf("failed to create dispatcher metrics: %v", err)
		return nil
	}

	d := &DefaultServerDispatcher{
		logger:           logger,
		queueMap:         queueMap,
		requestChannel:   make(chan string, 20),
		readyForDispatch: make(chan string, 1),
		timerC:           make(chan string, 10),
		stoppedC:         make(chan struct{}, 1),
		timeout:          defaultMessageTimeout,
		metrics:          dispatcherMetrics,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	d.pendingRequestState = NewServerState(&d.mutex)

	dispatcherMetrics.ObserveQueues(queueMap)
	dispatcherMetrics.ObserveInFlightRequests(d.pendingRequestState.(*serverState))

	return d
}

func (d *DefaultServerDispatcher) Start() {
	d.running.Store(true)

	d.queueMap.Init()
	d.requestChannel = make(chan string, 30)
	d.readyForDispatch = make(chan string, 1)
	d.timerC = make(chan string, 10)
	d.stoppedC = make(chan struct{}, 1)

	go d.messagePump()
}

func (d *DefaultServerDispatcher) IsRunning() bool {
	return d.running.Load()
}

func (d *DefaultServerDispatcher) Stop() {
	d.running.Store(false)

	// Close all channels
	close(d.stoppedC)
	close(d.readyForDispatch)
	close(d.timerC)
	close(d.requestChannel)
}

func (d *DefaultServerDispatcher) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *DefaultServerDispatcher) CreateClient(clientID string) {
	if d.IsRunning() {
		_ = d.queueMap.GetOrCreate(clientID)
	}
}

func (d *DefaultServerDispatcher) DeleteClient(clientID string) {
	d.queueMap.Remove(clientID)
	if d.IsRunning() {
		d.requestChannel <- clientID
	}
}

func (d *DefaultServerDispatcher) SetNetworkServer(server ws.Server) {
	d.network = server
}

func (d *DefaultServerDispatcher) SetOnRequestCanceled(cb CanceledRequestHandler) {
	d.onRequestCancel = cb
}

func (d *DefaultServerDispatcher) SetPendingRequestState(state ServerState) {
	d.pendingRequestState = state
}

func (d *DefaultServerDispatcher) SendRequest(clientID string, req RequestBundle) error {
	if d.network == nil {
		return fmt.Errorf("cannot send request %v, no network server was set", req.Call.UniqueId)
	}
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		return fmt.Errorf("cannot send request %s, no client %s exists", req.Call.UniqueId, clientID)
	}
	if err := q.Push(req); err != nil {
		return err
	}

	d.requestChannel <- clientID

	return nil
}

// requestPump processes new outgoing requests for each client and makes sure they are processed sequentially.
// This method is executed by a dedicated coroutine as soon as the server is started and runs indefinitely.
func (d *DefaultServerDispatcher) messagePump() {
	var clientID string
	var ok bool
	var rdy bool
	var clientCtx clientTimeoutContext
	var clientQueue RequestQueue
	clientContextMap := map[string]clientTimeoutContext{} // Empty at the beginning

	// Dispatcher Loop
	for {
		select {
		case <-d.stoppedC:
			// server was stopped
			d.queueMap.Init()
			d.logger.Info("stopped processing requests")
			return
		case clientID, ok = <-d.requestChannel:
			// Request channel closed, stopping dispatcher
			if !ok {
				continue
			}
			// Check whether there is a request queue for the specified client
			clientQueue, ok = d.queueMap.Get(clientID)
			if !ok {
				// No client queue found (client was removed)
				// Deleting and canceling the context
				clientCtx = clientContextMap[clientID]
				delete(clientContextMap, clientID)
				if clientCtx.ctx != nil {
					clientCtx.cancel()
				}
				continue
			}

			// Check whether we can transmit to client
			clientCtx, ok = clientContextMap[clientID]
			// Ready to transmit if its the first request or previous request timed out
			rdy = !ok || !clientCtx.isActive()

		case clientID, ok = <-d.timerC:
			// Timeout elapsed
			if !ok {
				continue
			}
			// Canceling timeout context
			d.logger.Debugf("timeout for client %v, canceling message", clientID)
			clientCtx = clientContextMap[clientID]
			if clientCtx.isActive() {
				clientCtx.cancel()
				clientContextMap[clientID] = clientTimeoutContext{}
			}

			if d.pendingRequestState.HasPendingRequest(clientID) {
				// Current request for client timed out. Removing request and triggering cancel callback
				q, found := d.queueMap.Get(clientID)
				if !found {
					// Possible race condition: queue was already removed
					d.logger.Errorf("dispatcher timeout for client %s triggered, but no request queue found", clientID)
					continue
				}
				el := q.Peek()
				if el == nil {
					// Should never happen
					d.logger.Error("dispatcher timeout for client %s triggered, but no pending request found", clientID)
					continue
				}
				bundle, _ := el.(RequestBundle)
				d.CompleteRequest(clientID, bundle.Call.UniqueId)
				d.logger.Infof("request %v for %v timed out", bundle.Call.UniqueId, clientID)
				if d.onRequestCancel != nil {
					d.onRequestCancel(clientID, bundle.Call.UniqueId, bundle.Call.Payload,
						ocpp.NewError(GenericError, "Request timed out", bundle.Call.UniqueId))
				}
			}
		case clientID = <-d.readyForDispatch:
			// Cancel previous timeout (if any)
			clientCtx, ok = clientContextMap[clientID]
			if ok && clientCtx.isActive() {
				clientCtx.cancel()
				clientContextMap[clientID] = clientTimeoutContext{}
			}
			// client can now transmit again
			clientQueue, ok = d.queueMap.Get(clientID)
			if ok {
				// Ready to transmit
				rdy = true
			}
			d.logger.Debugf("%v ready to transmit again", clientID)
		}

		// Only dispatch request if able to send and request queue isn't empty
		if rdy && clientQueue != nil && !clientQueue.IsEmpty() {
			// Send request & set new context
			clientCtx = d.dispatchNextRequest(clientID)
			clientContextMap[clientID] = clientCtx
			if clientCtx.isActive() {
				go d.waitForTimeout(clientID, clientCtx)
			}
			// Update ready state
			rdy = false
		}
	}
}

func (d *DefaultServerDispatcher) dispatchNextRequest(clientID string) (clientCtx clientTimeoutContext) {
	// Get first element in queue
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		d.logger.Errorf("failed to dispatch next request for %s, no request queue available", clientID)
		return
	}
	el := q.Peek()
	bundle, canCast := el.(RequestBundle)
	if !canCast {
		d.logger.Errorf("failed to cast request queue element to RequestBundle for client %s", clientID)
		return
	}

	jsonMessage := bundle.Data
	if bundle.Call == nil {
		d.logger.Errorf("request bundle has no Call associated")
		return
	}
	if bundle.Data == nil {
		d.logger.Errorf("request bundle has no Data associated")
		return
	}

	callID := bundle.Call.GetUniqueId()
	d.pendingRequestState.AddPendingRequest(clientID, callID, bundle.Call.Payload)
	err := d.network.Write(clientID, jsonMessage)
	if err != nil {
		d.logger.Errorf("error while sending message: %v", err)
		// TODO: handle retransmission instead of removing pending request
		d.CompleteRequest(clientID, callID)
		if d.onRequestCancel != nil {
			d.onRequestCancel(clientID, bundle.Call.UniqueId, bundle.Call.Payload,
				ocpp.NewError(InternalError, err.Error(), bundle.Call.UniqueId))
		}
		return
	}
	// Create and return context (only if timeout is set)
	if d.timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
		clientCtx = clientTimeoutContext{ctx: ctx, cancel: cancel}
	}
	d.logger.Infof("dispatched request %s for %s", callID, clientID)
	d.logger.Debugf("sent JSON message to %s: %s", clientID, string(jsonMessage))
	return
}

func (d *DefaultServerDispatcher) waitForTimeout(clientID string, clientCtx clientTimeoutContext) {
	defer clientCtx.cancel()
	d.logger.Debugf("started timeout timer for %s", clientID)
	select {
	case <-clientCtx.ctx.Done():
		err := clientCtx.ctx.Err()
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			// Timeout triggered, notifying messagePump
			d.mutex.RLock()
			running := d.running.Load()
			timerC := d.timerC
			d.mutex.RUnlock()
			if running && timerC != nil {
				timerC <- clientID
			}
		default:
			d.logger.Debugf("timeout canceled for %s", clientID)
		}
	case <-d.stoppedC:
		// server was stopped, every pending timeout gets canceled
	}
}

func (d *DefaultServerDispatcher) CompleteRequest(clientID string, requestID string) {
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		d.logger.Errorf("attempting to complete request for client %v, but no matching queue found", clientID)
		return
	}
	el := q.Peek()
	if el == nil {
		d.logger.Errorf("attempting to pop front of queue, but queue is empty")
		return
	}
	bundle, _ := el.(RequestBundle)
	callID := bundle.Call.GetUniqueId()
	if callID != requestID {
		d.logger.Errorf("internal state mismatch: processing response for %v but expected response for %v", requestID, callID)
		return
	}
	q.Pop()
	d.pendingRequestState.DeletePendingRequest(clientID, requestID)
	d.logger.Debugf("completed request %s for %s", callID, clientID)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- clientID
}
