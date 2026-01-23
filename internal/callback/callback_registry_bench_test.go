package callback

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
)

const (
	benchNumClients         = 100
	benchCallbacksPerClient = 10
)

// setupBenchmarkRegistry creates a registry with pre-populated callbacks for multiple clients
func setupBenchmarkRegistry(b *testing.B, numClients, callbacksPerClient int) (*Registry, []string, []string) {
	registry := New()
	clientIDs := make([]string, numClients)
	requestIDs := make([]string, numClients*callbacksPerClient)

	callback := func(confirmation ocpp.Response, err error) {
		// No-op callback for benchmarking
	}

	idx := 0
	for i := 0; i < numClients; i++ {
		clientID := fmt.Sprintf("client-%d", i)
		clientIDs[i] = clientID

		for j := 0; j < callbacksPerClient; j++ {
			requestID := fmt.Sprintf("req-%d-%d", i, j)
			requestIDs[idx] = requestID
			idx++

			requestIDCopy := requestID // Capture for closure
			try := func() (string, error) {
				return requestIDCopy, nil
			}

			err := registry.RegisterCallback(clientID, try, callback)
			if err != nil {
				b.Fatalf("failed to register callback: %v", err)
			}
		}
	}

	return registry, clientIDs, requestIDs
}

// BenchmarkRegisterCallback_MultipleClients benchmarks RegisterCallback with multiple clients
func BenchmarkRegisterCallback_MultipleClients(b *testing.B) {
	registry := New()
	clientIDs := make([]string, benchNumClients)
	for i := 0; i < benchNumClients; i++ {
		clientIDs[i] = fmt.Sprintf("client-%d", i)
	}

	callback := func(confirmation ocpp.Response, err error) {
		// No-op callback for benchmarking
	}

	var requestIdx int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentIdx := int(atomic.AddInt64(&requestIdx, 1) - 1)
			currentClientIdx := currentIdx % benchNumClients
			currentRequestIdx := currentIdx / benchNumClients

			clientID := clientIDs[currentClientIdx]
			requestID := fmt.Sprintf("req-%d-%d", currentClientIdx, currentRequestIdx)

			requestIDCopy := requestID // Capture for closure
			try := func() (string, error) {
				return requestIDCopy, nil
			}

			_ = registry.RegisterCallback(clientID, try, callback)
		}
	})
}

// BenchmarkGetCallback_MultipleClients benchmarks GetCallback with multiple clients
func BenchmarkGetCallback_MultipleClients(b *testing.B) {
	registry, clientIDs, _ := setupBenchmarkRegistry(b, benchNumClients, benchCallbacksPerClient)

	var idx int64
	totalRequests := int64(benchNumClients * benchCallbacksPerClient)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentIdx := int(atomic.AddInt64(&idx, 1)-1) % int(totalRequests)
			clientIdx := currentIdx / benchCallbacksPerClient
			requestIdx := currentIdx % benchCallbacksPerClient
			clientID := clientIDs[clientIdx]
			requestID := fmt.Sprintf("req-%d-%d", clientIdx, requestIdx)

			_, _ = registry.GetCallback(clientID, requestID)

			// Re-register if we've exhausted all callbacks
			if currentIdx == int(totalRequests)-1 {
				callback := func(confirmation ocpp.Response, err error) {}
				for i := 0; i < benchNumClients; i++ {
					for j := 0; j < benchCallbacksPerClient; j++ {
						reqID := fmt.Sprintf("req-%d-%d", i, j)
						reqIDCopy := reqID
						try := func() (string, error) {
							return reqIDCopy, nil
						}
						_ = registry.RegisterCallback(clientIDs[i], try, callback)
					}
				}
			}
		}
	})
}

// BenchmarkRemoveCallback_MultipleClients benchmarks RemoveCallback with multiple clients
func BenchmarkRemoveCallback_MultipleClients(b *testing.B) {
	registry, clientIDs, _ := setupBenchmarkRegistry(b, benchNumClients, benchCallbacksPerClient)

	var idx int64
	totalRequests := int64(benchNumClients * benchCallbacksPerClient)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentIdx := int(atomic.AddInt64(&idx, 1)-1) % int(totalRequests)
			clientIdx := currentIdx / benchCallbacksPerClient
			requestIdx := currentIdx % benchCallbacksPerClient
			clientID := clientIDs[clientIdx]
			requestID := fmt.Sprintf("req-%d-%d", clientIdx, requestIdx)

			_ = registry.RemoveCallback(clientID, requestID)

			// Re-register if we've exhausted all callbacks
			if currentIdx == int(totalRequests)-1 {
				callback := func(confirmation ocpp.Response, err error) {}
				for i := 0; i < benchNumClients; i++ {
					for j := 0; j < benchCallbacksPerClient; j++ {
						reqID := fmt.Sprintf("req-%d-%d", i, j)
						reqIDCopy := reqID
						try := func() (string, error) {
							return reqIDCopy, nil
						}
						_ = registry.RegisterCallback(clientIDs[i], try, callback)
					}
				}
			}
		}
	})
}

// BenchmarkGetAllCallbacks_MultipleClients benchmarks GetAllCallbacks with multiple clients
func BenchmarkGetAllCallbacks_MultipleClients(b *testing.B) {
	registry, clientIDs, _ := setupBenchmarkRegistry(b, benchNumClients, benchCallbacksPerClient)

	var clientIdx int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentClientIdx := int(atomic.AddInt64(&clientIdx, 1)-1) % benchNumClients
			clientID := clientIDs[currentClientIdx]
			_, _ = registry.GetAllCallbacks(clientID)

			// Re-populate for next iteration
			callback := func(confirmation ocpp.Response, err error) {}
			for j := 0; j < benchCallbacksPerClient; j++ {
				requestID := fmt.Sprintf("req-%d-%d", currentClientIdx, j)
				requestIDCopy := requestID
				try := func() (string, error) {
					return requestIDCopy, nil
				}
				_ = registry.RegisterCallback(clientID, try, callback)
			}
		}
	})
}

// BenchmarkClearCallbacks_MultipleClients benchmarks ClearCallbacks with multiple clients
func BenchmarkClearCallbacks_MultipleClients(b *testing.B) {
	registry, clientIDs, _ := setupBenchmarkRegistry(b, benchNumClients, benchCallbacksPerClient)

	var clientIdx int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentClientIdx := int(atomic.AddInt64(&clientIdx, 1)-1) % benchNumClients
			clientID := clientIDs[currentClientIdx]
			_ = registry.ClearCallbacks(clientID)

			// Re-populate for next iteration
			callback := func(confirmation ocpp.Response, err error) {}
			for j := 0; j < benchCallbacksPerClient; j++ {
				requestID := fmt.Sprintf("req-%d-%d", currentClientIdx, j)
				requestIDCopy := requestID
				try := func() (string, error) {
					return requestIDCopy, nil
				}
				_ = registry.RegisterCallback(clientID, try, callback)
			}
		}
	})
}

// BenchmarkRegisterAndGet_MultipleClients benchmarks the combined RegisterCallback and GetCallback operations
func BenchmarkRegisterAndGet_MultipleClients(b *testing.B) {
	registry := New()
	clientIDs := make([]string, benchNumClients)
	for i := 0; i < benchNumClients; i++ {
		clientIDs[i] = fmt.Sprintf("client-%d", i)
	}

	callback := func(confirmation ocpp.Response, err error) {
		// No-op callback for benchmarking
	}

	var requestIdx int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentIdx := int(atomic.AddInt64(&requestIdx, 1) - 1)
			currentClientIdx := currentIdx % benchNumClients
			currentRequestIdx := currentIdx / benchNumClients

			clientID := clientIDs[currentClientIdx]
			requestID := fmt.Sprintf("req-%d-%d", currentClientIdx, currentRequestIdx)

			requestIDCopy := requestID // Capture for closure
			try := func() (string, error) {
				return requestIDCopy, nil
			}

			_ = registry.RegisterCallback(clientID, try, callback)
			_, _ = registry.GetCallback(clientID, requestID)
		}
	})
}

// BenchmarkConcurrentOperations_MultipleClients benchmarks concurrent operations across multiple clients
func BenchmarkConcurrentOperations_MultipleClients(b *testing.B) {
	registry, clientIDs, _ := setupBenchmarkRegistry(b, benchNumClients, benchCallbacksPerClient)

	var idx int64
	totalRequests := int64(benchNumClients * benchCallbacksPerClient)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			currentIdx := int(atomic.AddInt64(&idx, 1)-1) % int(totalRequests)
			clientIdx := currentIdx / benchCallbacksPerClient
			requestIdx := currentIdx % benchCallbacksPerClient
			clientID := clientIDs[clientIdx]
			requestID := fmt.Sprintf("req-%d-%d", clientIdx, requestIdx)

			// Mix of operations
			switch currentIdx % 4 {
			case 0:
				// Register
				callback := func(confirmation ocpp.Response, err error) {}
				requestIDCopy := requestID
				try := func() (string, error) {
					return requestIDCopy, nil
				}
				_ = registry.RegisterCallback(clientID, try, callback)
			case 1:
				// Get
				_, _ = registry.GetCallback(clientID, requestID)
			case 2:
				// Remove
				_ = registry.RemoveCallback(clientID, requestID)
			case 3:
				// GetAll
				_, _ = registry.GetAllCallbacks(clientID)
				// Re-populate after GetAll
				callback := func(confirmation ocpp.Response, err error) {}
				for j := 0; j < benchCallbacksPerClient; j++ {
					reqID := fmt.Sprintf("req-%d-%d", clientIdx, j)
					reqIDCopy := reqID
					try := func() (string, error) {
						return reqIDCopy, nil
					}
					_ = registry.RegisterCallback(clientID, try, callback)
				}
			}
		}
	})
}
