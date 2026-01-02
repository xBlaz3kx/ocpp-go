package ocppj_test

import (
	"fmt"
	"testing"

	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

// BenchmarkFIFOClientQueue_PushParallel benchmarks concurrent pushes
func BenchmarkFIFOClientQueue_PushParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0) // Unlimited capacity
	req := newMockRequest("benchmark")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Push(req)
		}
	})
}

// BenchmarkFIFOClientQueue_PopParallel benchmarks concurrent pops
func BenchmarkFIFOClientQueue_PopParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")

	// Pre-populate queue
	for i := 0; i < b.N; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Pop()
		}
	})
}

// BenchmarkFIFOClientQueue_PeekParallel benchmarks concurrent peeks
func BenchmarkFIFOClientQueue_PeekParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")
	_ = queue.Push(req)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Peek()
		}
	})
}

// BenchmarkFIFOClientQueue_SizeParallel benchmarks concurrent size checks
func BenchmarkFIFOClientQueue_SizeParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")
	for i := 0; i < 1000; i++ {
		_ = queue.Push(req)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Size()
		}
	})
}

// BenchmarkFIFOClientQueue_PushPopParallel benchmarks concurrent push-pop operations
func BenchmarkFIFOClientQueue_PushPopParallel(b *testing.B) {
	queue := ocppj.NewFIFOClientQueue(0)
	req := newMockRequest("benchmark")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queue.Push(req)
			_ = queue.Pop()
		}
	})
}

// BenchmarkFIFOQueueMap_GetParallel benchmarks concurrent gets
func BenchmarkFIFOQueueMap_GetParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	clientID := "client1"
	queueMap.GetOrCreate(clientID)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = queueMap.Get(clientID)
		}
	})
}

// BenchmarkFIFOQueueMap_GetOrCreateParallel benchmarks concurrent get-or-create operations
func BenchmarkFIFOQueueMap_GetOrCreateParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			clientID := fmt.Sprintf("client%d", i%10) // Reuse 10 clients
			_ = queueMap.GetOrCreate(clientID)
			i++
		}
	})
}

// BenchmarkFIFOQueueMap_SizeParallel benchmarks concurrent size checks
func BenchmarkFIFOQueueMap_SizeParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	// Pre-populate with queues and elements
	for i := 0; i < 100; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := queueMap.GetOrCreate(clientID)
		for j := 0; j < 10; j++ {
			_ = queue.Push(req)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queueMap.Size()
		}
	})
}

// BenchmarkFIFOQueueMap_SizePerClientParallel benchmarks concurrent size-per-client checks
func BenchmarkFIFOQueueMap_SizePerClientParallel(b *testing.B) {
	queueMap := ocppj.NewFIFOQueueMap(100)
	req := newMockRequest("benchmark")

	// Pre-populate with queues and elements
	for i := 0; i < 100; i++ {
		clientID := fmt.Sprintf("client%d", i)
		queue := queueMap.GetOrCreate(clientID)
		for j := 0; j < 10; j++ {
			_ = queue.Push(req)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = queueMap.SizePerClient()
		}
	})
}
