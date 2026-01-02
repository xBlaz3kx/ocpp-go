package ocppj_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

// BenchmarkClientState_AddPendingRequest benchmarks concurrent AddPendingRequest operations
func BenchmarkClientState_AddPendingRequest(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(requestID, req)
			counter++
			state.DeletePendingRequest(requestID)
		}
	})
}

// BenchmarkClientState_GetPendingRequest benchmarks concurrent GetPendingRequest operations
func BenchmarkClientState_GetPendingRequest(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	requestID := "test-request"
	state.AddPendingRequest(requestID, req)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = state.GetPendingRequest(requestID)
		}
	})
}

// BenchmarkClientState_HasPendingRequest benchmarks concurrent HasPendingRequest operations
func BenchmarkClientState_HasPendingRequest(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			state.AddPendingRequest("test-request", req)
			_ = state.HasPendingRequest()
		}
	})
}

// BenchmarkClientState_DeletePendingRequest benchmarks concurrent DeletePendingRequest operations
func BenchmarkClientState_DeletePendingRequest(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(requestID, req)
			state.DeletePendingRequest(requestID)
			counter++
		}
	})
}

// BenchmarkClientState_ClearPendingRequests benchmarks concurrent ClearPendingRequests operations
func BenchmarkClientState_ClearPendingRequests(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			state.AddPendingRequest("test-request", req)
			state.ClearPendingRequests()
		}
	})
}

// BenchmarkClientState_AddAndDelete benchmarks concurrent AddAndDelete operations
func BenchmarkClientState_AddAndDelete(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(requestID, req)
			state.DeletePendingRequest(requestID)
			counter++
		}
	})
}

// BenchmarkClientState_MixedOperations benchmarks concurrent mixed operations
func BenchmarkClientState_MixedOperations(b *testing.B) {
	state := ocppj.NewClientState()
	req := newMockRequest("benchmark")
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(requestID, req)
			_ = state.HasPendingRequest()
			_, _ = state.GetPendingRequest(requestID)
			state.DeletePendingRequest(requestID)
			counter++
		}
	})
}

// BenchmarkServerState_AddPendingRequest benchmarks concurrent AddPendingRequest operations
func BenchmarkServerState_AddPendingRequest(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
	}

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(clientID, requestID, req)
			state.DeletePendingRequest(clientID, requestID)
			counter++
		}
	})
}

// BenchmarkServerState_DeletePendingRequest benchmarks concurrent DeletePendingRequest operations
func BenchmarkServerState_DeletePendingRequest(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
	}
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(clientID, requestID, req)
			state.DeletePendingRequest(clientID, requestID)
			counter++
		}
	})
}

// BenchmarkServerState_GetClientState benchmarks concurrent GetClientState operations
func BenchmarkServerState_GetClientState(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")

	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
		state.AddPendingRequest(clientIDs[i], "pre-req", req)
	}

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			_ = state.GetClientState(clientID)
			counter++
		}
	})
}

// BenchmarkServerState_HasPendingRequest benchmarks concurrent HasPendingRequest operations
func BenchmarkServerState_HasPendingRequest(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)

	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
		state.AddPendingRequest(clientIDs[i], "pre-req", req)
	}

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			_ = state.HasPendingRequest(clientID)
			counter++
		}
	})
}

// BenchmarkServerState_HasPendingRequests benchmarks concurrent HasPendingRequests operations
func BenchmarkServerState_HasPendingRequests(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")

	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
		state.AddPendingRequest(clientIDs[i], "pre-req", req)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = state.HasPendingRequests()
		}
	})
}

// BenchmarkServerState_ClearClientPendingRequest benchmarks concurrent ClearClientPendingRequest operations
func BenchmarkServerState_ClearClientPendingRequest(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
	}
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			state.AddPendingRequest(clientID, "pre-req", req)
			state.ClearClientPendingRequest(clientID)
			counter++
		}
	})
}

// BenchmarkServerState_ClearAllPendingRequests benchmarks concurrent ClearAllPendingRequests operations
func BenchmarkServerState_ClearAllPendingRequests(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, clientID := range clientIDs {
				state.AddPendingRequest(clientID, "pre-req", req)
			}
			state.ClearAllPendingRequests()
		}
	})
}

// BenchmarkServerState_AddAndDelete benchmarks concurrent AddAndDelete operations
func BenchmarkServerState_AddAndDelete(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
	}
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(clientID, requestID, req)
			state.DeletePendingRequest(clientID, requestID)
			counter++
		}
	})
}

// BenchmarkServerState_MixedOperations benchmarks concurrent mixed operations
func BenchmarkServerState_MixedOperations(b *testing.B) {
	var mutex sync.RWMutex
	state := ocppj.NewServerState(&mutex)
	req := newMockRequest("benchmark")
	numClients := 10
	clientIDs := make([]string, numClients)
	for i := 0; i < numClients; i++ {
		clientIDs[i] = fmt.Sprintf("client%d", i)
	}
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			clientID := clientIDs[counter%len(clientIDs)]
			requestID := fmt.Sprintf("req-%d", counter)
			state.AddPendingRequest(clientID, requestID, req)
			_ = state.HasPendingRequest(clientID)
			_ = state.GetClientState(clientID)
			_ = state.HasPendingRequests()
			state.DeletePendingRequest(clientID, requestID)
			counter++
		}
	})
}
