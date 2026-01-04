package ocppj

import (
	"time"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ws"
)

// ClientDispatcher contains the state and logic for handling outgoing messages on a client endpoint.
// This allows the ocpp-j layer to delegate queueing and processing logic to an external entity.
//
// The dispatcher writes outgoing messages directly to the networking layer, using a previously set websocket client.
//
// A ClientState needs to be passed to the dispatcher, before starting it.
// The dispatcher is in charge of managing pending requests while handling the request flow.
type ClientDispatcher interface {
	// Starts the dispatcher. Depending on the implementation, this may
	// start a dedicated goroutine or simply allocate the necessary state.
	Start()
	// Sets the maximum timeout to be considered after sending a request.
	// If a response to the request is not received within the specified period, the request
	// is discarded and an error is returned to the caller.
	//
	// The timeout is reset upon a disconnection/reconnection.
	//
	// This function must be called before starting the dispatcher, otherwise it may lead to unexpected behavior.
	SetTimeout(timeout time.Duration)
	// Returns true, if the dispatcher is currently running, false otherwise.
	// If the dispatcher is paused, the function still returns true.
	IsRunning() bool
	// Returns true, if the dispatcher is currently paused, false otherwise.
	// If the dispatcher is not running at all, the function will still return false.
	IsPaused() bool
	// Dispatches a request. Depending on the implementation, this may first queue a request
	// and process it later, asynchronously, or write it directly to the networking layer.
	//
	// If no network client was set, or the request couldn't be processed, an error is returned.
	SendRequest(req RequestBundle) error
	// Notifies the dispatcher that a request has been completed (i.e. a response was received).
	// The dispatcher takes care of removing the request marked by the requestID from
	// the pending requests. It will then attempt to process the next queued request.
	CompleteRequest(requestID string)
	// Sets a callback to be invoked when a request gets canceled, due to network timeouts or internal errors.
	// The callback passes the original message ID and request struct of the failed request, along with an error.
	//
	// Calling Stop on the dispatcher will not trigger this callback.
	//
	// If no callback is set, a request will still be removed from the dispatcher when a timeout occurs.
	SetOnRequestCanceled(cb func(requestID string, request ocpp.Request, err *ocpp.Error))
	// Sets the network client, so the dispatcher may send requests using the networking layer directly.
	//
	// This needs to be set before calling the Start method. If not, sending requests will fail.
	SetNetworkClient(client ws.Client)
	// Sets the state manager for pending requests in the dispatcher.
	//
	// The state should only be accessed by the dispatcher while running.
	SetPendingRequestState(stateHandler ClientState)
	// Stops a running dispatcher. This will clear all state and empty the internal queues.
	//
	// If an onRequestCanceled callback is set, it won't be triggered by stopping the dispatcher.
	Stop()
	// Notifies that an external event (typically network-related) should pause
	// the dispatcher. Internal timers will be stopped an no further requests
	// will be set to pending. You may keep enqueuing requests.
	// Use the Resume method for re-starting the dispatcher.
	Pause()
	// Undoes a previous pause operation, restarting internal timers and the
	// regular request flow.
	//
	// If there was a pending request before pausing the dispatcher, a response/timeout
	// for this request shall be awaited anew.
	Resume()
}

// pendingRequest is used internally for associating metadata to a pending Request.
type pendingRequest struct {
	request ocpp.Request
}

// ServerDispatcher contains the state and logic for handling outgoing messages on a server endpoint.
// This allows the ocpp-j layer to delegate queueing and processing logic to an external entity.
//
// The dispatcher writes outgoing messages directly to the networking layer, using a previously set websocket server.
//
// A ClientState needs to be passed to the dispatcher, before starting it.
// The dispatcher is in charge of managing all pending requests to clients, while handling the request flow.
type ServerDispatcher interface {
	// Starts the dispatcher. Depending on the implementation, this may
	// start a dedicated goroutine or simply allocate the necessary state.
	Start()
	// Returns true, if the dispatcher is currently running, false otherwise.
	// If the dispatcher is paused, the function still returns true.
	IsRunning() bool
	// Sets the maximum timeout to be considered after sending a request.
	// If a response to the request is not received within the specified period, the request
	// is discarded and an error is returned to the caller.
	//
	// One timeout per client runs in the background.
	// The timeout is reset whenever a response comes in, the connection is closed, or the server is stopped.
	//
	// This function must be called before starting the dispatcher, otherwise it may lead to unexpected behavior.
	SetTimeout(timeout time.Duration)
	// Dispatches a request for a specific client. Depending on the implementation, this may first queue
	// a request and process it later (asynchronously), or write it directly to the networking layer.
	//
	// If no network server was set, or the request couldn't be processed, an error is returned.
	SendRequest(clientID string, req RequestBundle) error
	// Notifies the dispatcher that a request has been completed (i.e. a response was received),
	// for a specific client.
	// The dispatcher takes care of removing the request marked by the requestID from
	// that client's pending requests. It will then attempt to process the next queued request.
	CompleteRequest(clientID string, requestID string)
	// Sets a callback to be invoked when a request gets canceled, due to network timeouts.
	// The callback passes the original client ID, message ID, and request struct of the failed request,
	// along with an error.
	//
	// Calling Stop on the dispatcher will not trigger this callback.
	//
	// If no callback is set, a request will still be removed from the dispatcher when a timeout occurs.
	SetOnRequestCanceled(cb CanceledRequestHandler)
	// Sets the network server, so the dispatcher may send requests using the networking layer directly.
	//
	// This needs to be set before calling the Start method. If not, sending requests will fail.
	SetNetworkServer(server ws.Server)
	// Sets the state manager for pending requests in the dispatcher.
	//
	// The state should only be accessed by the dispatcher while running.
	SetPendingRequestState(stateHandler ServerState)
	// Stops a running dispatcher. This will clear all state and empty the internal queues.
	//
	// If an onRequestCanceled callback is set, it won't be triggered by stopping the dispatcher.
	Stop()
	// Notifies that it is now possible to dispatch requests for a new client.
	//
	// Internal queues are created and requests for the client are now accepted.
	CreateClient(clientID string)
	// Notifies that a client was invalidated (typically caused by a network event).
	//
	// The dispatcher will stop dispatching requests for that specific client.
	// Internal queues for that client are cleared and no further requests will be accepted.
	// Undelivered pending requests are also cleared.
	// The OnRequestCanceled callback will be invoked for each discarded request.
	DeleteClient(clientID string)
}
