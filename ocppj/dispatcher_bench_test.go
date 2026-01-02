package ocppj_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

// setupClientDispatcher creates and configures a client dispatcher for benchmarking
// Returns the dispatcher and an example request bundle
func setupClientDispatcher() (*ocppj.DefaultClientDispatcher, ocppj.RequestBundle, error) {
	endpoint := &ocppj.Client{Id: "client1"}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queue := ocppj.NewFIFOClientQueue(1000)
	dispatcher := ocppj.NewDefaultClientDispatcher(queue, nil)
	state := ocppj.NewClientState()
	dispatcher.SetPendingRequestState(state)

	websocketClient := &MockWebsocketClient{}
	websocketClient.On("Write", mock.Anything).Return(nil)
	dispatcher.SetNetworkClient(websocketClient)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()

	bundle, err := createRequestBundle(endpoint, "benchmark")
	if err != nil {
		dispatcher.Stop()
		return nil, ocppj.RequestBundle{}, err
	}

	return dispatcher, bundle, nil
}

// setupServerDispatcher creates and configures a server dispatcher for benchmarking
// Returns the dispatcher and an example request bundle
func setupServerDispatcher() (*ocppj.DefaultServerDispatcher, ocppj.RequestBundle, error) {
	endpoint := &ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)
	queueMap := ocppj.NewFIFOQueueMap(1000)
	dispatcher := ocppj.NewDefaultServerDispatcher(queueMap, noop.NewMeterProvider(), nil)

	websocketServer := &MockWebsocketServer{}
	websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	dispatcher.SetNetworkServer(websocketServer)
	dispatcher.SetTimeout(30 * time.Second)
	dispatcher.Start()

	bundle, err := createRequestBundle(endpoint, "benchmark")
	if err != nil {
		dispatcher.Stop()
		return nil, ocppj.RequestBundle{}, err
	}

	return dispatcher, bundle, nil
}

// createRequestBundle creates a request bundle from an endpoint and request value
func createRequestBundle(endpoint interface {
	CreateCall(request ocpp.Request) (*ocppj.Call, error)
}, requestValue string) (ocppj.RequestBundle, error) {
	req := newMockRequest(requestValue)
	call, err := endpoint.CreateCall(req)
	if err != nil {
		return ocppj.RequestBundle{}, err
	}

	data, err := call.MarshalJSON()
	if err != nil {
		return ocppj.RequestBundle{}, err
	}

	return ocppj.RequestBundle{Call: call, Data: data}, nil
}

// BenchmarkClientDispatcher_SendRequest benchmarks sending requests through the client dispatcher
func BenchmarkClientDispatcher_SendRequest(b *testing.B) {
	dispatcher, bundle, err := setupClientDispatcher()
	if err != nil {
		b.Fatalf("failed to setup client dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dispatcher.SendRequest(bundle)
		}
	})
}

// BenchmarkClientDispatcher_CompleteRequest benchmarks completing requests in the client dispatcher
func BenchmarkClientDispatcher_CompleteRequest(b *testing.B) {
	dispatcher, _, err := setupClientDispatcher()
	if err != nil {
		b.Fatalf("failed to setup client dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	endpoint := &ocppj.Client{Id: "client1"}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)

	// Pre-populate queue with requests
	requestIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		newBundle, err := createRequestBundle(endpoint, "benchmark")
		if err != nil {
			b.Fatalf("failed to create request bundle: %v", err)
		}
		requestIDs[i] = newBundle.Call.UniqueId
		_ = dispatcher.SendRequest(newBundle)
	}

	// Wait for requests to be dispatched
	time.Sleep(100 * time.Millisecond)

	var mu sync.Mutex
	idx := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			currentIdx := idx
			idx++
			mu.Unlock()
			if currentIdx < len(requestIDs) {
				dispatcher.CompleteRequest(requestIDs[currentIdx])
			}
		}
	})
}

// BenchmarkClientDispatcher_SendAndComplete benchmarks the full cycle of sending and completing requests
func BenchmarkClientDispatcher_SendAndComplete(b *testing.B) {
	dispatcher, bundle, err := setupClientDispatcher()
	if err != nil {
		b.Fatalf("failed to setup client dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dispatcher.SendRequest(bundle)
			dispatcher.CompleteRequest(bundle.Call.UniqueId)
		}
	})
}

// BenchmarkServerDispatcher_SendRequest benchmarks sending requests through the server dispatcher
func BenchmarkServerDispatcher_SendRequest(b *testing.B) {
	dispatcher, bundle, err := setupServerDispatcher()
	if err != nil {
		b.Fatalf("failed to setup server dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	clientID := "client1"
	dispatcher.CreateClient(clientID)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dispatcher.SendRequest(clientID, bundle)
		}
	})
}

// BenchmarkServerDispatcher_CompleteRequest benchmarks completing requests in the server dispatcher
func BenchmarkServerDispatcher_CompleteRequest(b *testing.B) {
	dispatcher, _, err := setupServerDispatcher()
	if err != nil {
		b.Fatalf("failed to setup server dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	endpoint := &ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	endpoint.AddProfile(mockProfile)

	clientID := "client1"
	dispatcher.CreateClient(clientID)

	// Pre-populate queue with requests
	requestIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		newBundle, err := createRequestBundle(endpoint, "benchmark")
		if err != nil {
			b.Fatalf("failed to create request bundle: %v", err)
		}
		requestIDs[i] = newBundle.Call.UniqueId
		_ = dispatcher.SendRequest(clientID, newBundle)
	}

	// Wait for requests to be dispatched
	time.Sleep(100 * time.Millisecond)

	var mu sync.Mutex
	idx := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			currentIdx := idx
			idx++
			mu.Unlock()
			if currentIdx < len(requestIDs) {
				dispatcher.CompleteRequest(clientID, requestIDs[currentIdx])
			}
		}
	})
}

// BenchmarkServerDispatcher_SendAndComplete benchmarks the full cycle of sending and completing requests
func BenchmarkServerDispatcher_SendAndComplete(b *testing.B) {
	dispatcher, bundle, err := setupServerDispatcher()
	if err != nil {
		b.Fatalf("failed to setup server dispatcher: %v", err)
	}
	defer dispatcher.Stop()

	clientID := "client1"
	dispatcher.CreateClient(clientID)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = dispatcher.SendRequest(clientID, bundle)
			dispatcher.CompleteRequest(clientID, bundle.Call.UniqueId)
		}
	})
}
