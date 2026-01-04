package callback

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/ocpp-go/ocpp"
)

// MockResponse is a simple mock implementation of ocpp.Response for testing
type MockResponse struct {
	Value string
}

func (m *MockResponse) GetFeatureName() string {
	return "Mock"
}

// CallbackRegistryTestSuite contains tests for the CallbackRegistry
type CallbackRegistryTestSuite struct {
	suite.Suite
	registry *Registry
}

func (suite *CallbackRegistryTestSuite) SetupTest() {
	suite.registry = New()
}

// TestRegisterCallbackSuccess verifies that RegisterCallback succeeds when try() returns nil
func (suite *CallbackRegistryTestSuite) TestRegisterCallbackSuccess() {
	clientID := "client-1"
	requestID := "req-1"
	callbackCalled := atomic.Bool{}
	var receivedResponse ocpp.Response
	var receivedError error

	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Store(true)
		receivedResponse = confirmation
		receivedError = err
	}

	try := func() (string, error) {
		return requestID, nil // Success
	}

	err := suite.registry.RegisterCallback(clientID, try, callback)
	suite.Require().NoError(err)
	suite.Assert().False(callbackCalled.Load(), "Callback should not be called yet")

	// Get callback and invoke it
	cb, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Require().True(ok)
	suite.Require().NotNil(cb)

	mockResp := &MockResponse{Value: "test"}
	cb(mockResp, nil)

	suite.Assert().True(callbackCalled.Load())
	suite.Assert().Equal(mockResp, receivedResponse)
	suite.Assert().Nil(receivedError)
}

// TestRegisterCallbackFailure verifies that RegisterCallback removes callback when try() returns error
func (suite *CallbackRegistryTestSuite) TestRegisterCallbackFailure() {
	clientID := "client-1"
	requestID := "req-1"
	callbackCalled := atomic.Bool{}
	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Store(true)
	}

	try := func() (string, error) {
		return requestID, errors.New("try failed")
	}

	err := suite.registry.RegisterCallback(clientID, try, callback)
	suite.Require().Error(err)
	suite.Assert().Equal("try failed", err.Error())
	suite.Assert().False(callbackCalled.Load(), "Callback should not be called")

	// Verify callback was removed
	_, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Assert().False(ok, "Callback should have been removed")
}

// TestGetCallbackEmpty verifies that GetCallback returns false for non-existent callback
func (suite *CallbackRegistryTestSuite) TestGetCallbackEmpty() {
	clientID := "client-1"
	requestID := "non-existent"

	cb, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Assert().False(ok)
	suite.Assert().Nil(cb)
}

// TestGetCallbackByRequestID verifies that callbacks are retrieved by their specific request ID
func (suite *CallbackRegistryTestSuite) TestGetCallbackByRequestID() {
	clientID := "client-1"
	requestID1 := "req-1"
	requestID2 := "req-2"
	requestID3 := "req-3"

	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}
	callback3Called := atomic.Bool{}

	callback1 := func(confirmation ocpp.Response, err error) {
		callback1Called.Store(true)
	}
	callback2 := func(confirmation ocpp.Response, err error) {
		callback2Called.Store(true)
	}
	callback3 := func(confirmation ocpp.Response, err error) {
		callback3Called.Store(true)
	}

	// Register multiple callbacks with different request IDs
	try1 := func() (string, error) {
		return requestID1, nil
	}
	try2 := func() (string, error) {
		return requestID2, nil
	}
	try3 := func() (string, error) {
		return requestID3, nil
	}

	err := suite.registry.RegisterCallback(clientID, try1, callback1)
	suite.Require().NoError(err)
	err = suite.registry.RegisterCallback(clientID, try2, callback2)
	suite.Require().NoError(err)
	err = suite.registry.RegisterCallback(clientID, try3, callback3)
	suite.Require().NoError(err)

	// Retrieve callbacks by their specific request IDs (not in order)
	cb2, ok := suite.registry.GetCallback(clientID, requestID2)
	suite.Require().True(ok)
	cb2(nil, nil)
	suite.Assert().True(callback2Called.Load())
	suite.Assert().False(callback1Called.Load())
	suite.Assert().False(callback3Called.Load())

	cb1, ok := suite.registry.GetCallback(clientID, requestID1)
	suite.Require().True(ok)
	cb1(nil, nil)
	suite.Assert().True(callback1Called.Load())
	suite.Assert().False(callback3Called.Load())

	cb3, ok := suite.registry.GetCallback(clientID, requestID3)
	suite.Require().True(ok)
	cb3(nil, nil)
	suite.Assert().True(callback3Called.Load())

	// All callbacks should be removed now
	_, ok = suite.registry.GetCallback(clientID, requestID1)
	suite.Assert().False(ok)
	_, ok = suite.registry.GetCallback(clientID, requestID2)
	suite.Assert().False(ok)
	_, ok = suite.registry.GetCallback(clientID, requestID3)
	suite.Assert().False(ok)
}

// TestMultipleClients verifies that different clients maintain separate callback maps
func (suite *CallbackRegistryTestSuite) TestMultipleClients() {
	clientID1 := "client-1"
	clientID2 := "client-2"
	requestID := "req-1"

	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}

	callback1 := func(confirmation ocpp.Response, err error) {
		callback1Called.Store(true)
	}

	callback2 := func(confirmation ocpp.Response, err error) {
		callback2Called.Store(true)
	}

	try1 := func() (string, error) {
		return requestID, nil
	}
	try2 := func() (string, error) {
		return requestID, nil
	}

	// Register callbacks for different clients with same request ID
	err := suite.registry.RegisterCallback(clientID1, try1, callback1)
	suite.Require().NoError(err)

	err = suite.registry.RegisterCallback(clientID2, try2, callback2)
	suite.Require().NoError(err)

	// Get callback from clientID1
	cb1, ok := suite.registry.GetCallback(clientID1, requestID)
	suite.Require().True(ok)
	cb1(nil, nil)
	suite.Assert().True(callback1Called.Load())
	suite.Assert().False(callback2Called.Load())

	// Get callback from clientID2
	cb2, ok := suite.registry.GetCallback(clientID2, requestID)
	suite.Require().True(ok)
	cb2(nil, nil)
	suite.Assert().True(callback2Called.Load())

	// Both should be empty now
	_, ok = suite.registry.GetCallback(clientID1, requestID)
	suite.Assert().False(ok)
	_, ok = suite.registry.GetCallback(clientID2, requestID)
	suite.Assert().False(ok)
}

// TestGetCallbackRemovesFromRegistry verifies that GetCallback removes the callback from the registry
func (suite *CallbackRegistryTestSuite) TestGetCallbackRemovesFromRegistry() {
	clientID := "client-1"
	requestID := "req-1"
	callbackCalled := atomic.Int64{}

	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Add(1)
	}

	try := func() (string, error) {
		return requestID, nil
	}

	// Register callback
	err := suite.registry.RegisterCallback(clientID, try, callback)
	suite.Require().NoError(err)

	// Get callback first time
	cb1, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Require().True(ok)
	cb1(nil, nil)
	suite.Assert().Equal(1, int(callbackCalled.Load()))

	// Try to get callback again - should fail
	_, ok = suite.registry.GetCallback(clientID, requestID)
	suite.Assert().False(ok)
	suite.Assert().Equal(1, int(callbackCalled.Load()), "Callback should only be called once")
}

// TestRegisterCallbackFailureDoesNotAffectOthers verifies that failure doesn't affect other callbacks
func (suite *CallbackRegistryTestSuite) TestRegisterCallbackFailureDoesNotAffectOthers() {
	clientID := "client-1"
	requestID1 := "req-1"
	requestID2 := "req-2"
	requestID3 := "req-3"

	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}
	callback3Called := atomic.Bool{}

	callback1 := func(confirmation ocpp.Response, err error) {
		callback1Called.Store(true)
	}
	callback2 := func(confirmation ocpp.Response, err error) {
		callback2Called.Store(true)
	}
	callback3 := func(confirmation ocpp.Response, err error) {
		callback3Called.Store(true)
	}

	trySuccess1 := func() (string, error) {
		return requestID1, nil
	}

	trySuccess2 := func() (string, error) {
		return requestID2, nil
	}
	tryFail := func() (string, error) {
		return requestID3, errors.New("failed")
	}

	// Register two callbacks successfully
	err := suite.registry.RegisterCallback(clientID, trySuccess1, callback1)
	suite.Require().NoError(err)
	err = suite.registry.RegisterCallback(clientID, trySuccess2, callback2)
	suite.Require().NoError(err)

	// Register a third callback that fails
	err = suite.registry.RegisterCallback(clientID, tryFail, callback3)
	suite.Require().Error(err)

	// First two callbacks should still be registered
	cb1, ok := suite.registry.GetCallback(clientID, requestID1)
	suite.Require().True(ok)
	cb1(nil, nil)
	suite.Assert().True(callback1Called.Load())

	cb2, ok := suite.registry.GetCallback(clientID, requestID2)
	suite.Require().True(ok)
	cb2(nil, nil)
	suite.Assert().True(callback2Called.Load())

	// Third callback should not be registered
	_, ok = suite.registry.GetCallback(clientID, requestID3)
	suite.Assert().False(ok)
	suite.Assert().False(callback3Called.Load())
}

// TestConcurrentRegisterCallback verifies thread safety of RegisterCallback
func (suite *CallbackRegistryTestSuite) TestConcurrentRegisterCallback() {
	clientID := "client-1"
	numCallbacks := 100
	var wg sync.WaitGroup
	callbacksRegistered := make(chan int, numCallbacks)

	// Concurrently register callbacks with unique request IDs
	for i := 0; i < numCallbacks; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			requestID := fmt.Sprintf("req-%d", index)
			callback := func(confirmation ocpp.Response, err error) {
				callbacksRegistered <- index
			}

			try := func() (string, error) {
				return requestID, nil
			}

			err := suite.registry.RegisterCallback(clientID, try, callback)
			suite.Require().NoError(err)
		}(i)
	}

	wg.Wait()

	// Get all callbacks by their request IDs
	received := make(map[int]bool)
	for i := 0; i < numCallbacks; i++ {
		requestID := fmt.Sprintf("req-%d", i)
		cb, ok := suite.registry.GetCallback(clientID, requestID)
		suite.Require().True(ok, "Callback for requestID %s should exist", requestID)
		cb(nil, nil)
	}

	// Collect all received indices
	for i := 0; i < numCallbacks; i++ {
		index := <-callbacksRegistered
		received[index] = true
	}

	// Verify all callbacks were registered and retrieved
	suite.Assert().Len(received, numCallbacks)

	// Registry should be empty now
	_, ok := suite.registry.GetCallback(clientID, "req-0")
	suite.Assert().False(ok)
}

// TestConcurrentGetCallback verifies thread safety of GetCallback
func (suite *CallbackRegistryTestSuite) TestConcurrentGetCallback() {
	clientID := "client-1"
	numCallbacks := 50
	var wg sync.WaitGroup
	dequeued := make(chan bool, numCallbacks)

	// Register callbacks with unique request IDs
	for i := 0; i < numCallbacks; i++ {
		requestID := fmt.Sprintf("req-%d", i)
		callback := func(confirmation ocpp.Response, err error) {
			dequeued <- true
		}

		requestIDCopy := requestID // Capture for closure
		try := func() (string, error) {
			return requestIDCopy, nil
		}

		err := suite.registry.RegisterCallback(clientID, try, callback)
		suite.Require().NoError(err)
	}

	// Concurrently get callbacks
	for i := 0; i < numCallbacks; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			requestID := fmt.Sprintf("req-%d", index)
			cb, ok := suite.registry.GetCallback(clientID, requestID)
			if ok && cb != nil {
				cb(nil, nil)
			}
		}(i)
	}

	wg.Wait()

	// Verify all callbacks were retrieved
	suite.Assert().Len(dequeued, numCallbacks)

	// Registry should be empty now
	_, ok := suite.registry.GetCallback(clientID, "req-0")
	suite.Assert().False(ok)
}

// TestGetCallbackWithError verifies that callbacks can be invoked with errors
func (suite *CallbackRegistryTestSuite) TestGetCallbackWithError() {
	clientID := "client-1"
	requestID := "req-1"
	var receivedError error
	var receivedResponse ocpp.Response

	callback := func(confirmation ocpp.Response, err error) {
		receivedResponse = confirmation
		receivedError = err
	}

	try := func() (string, error) {
		return requestID, nil
	}

	err := suite.registry.RegisterCallback(clientID, try, callback)
	suite.Require().NoError(err)

	// Get callback and invoke with error
	ocppErr := ocpp.NewError(ocpp.ErrorCode("GenericError"), "test error", "123")
	cb, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Require().True(ok)
	cb(nil, ocppErr)

	suite.Assert().Nil(receivedResponse)
	suite.Assert().NotNil(receivedError)
	suite.Assert().Equal(ocppErr, receivedError)
}

// TestGetCallbackWithResponse verifies that callbacks can be invoked with responses
func (suite *CallbackRegistryTestSuite) TestGetCallbackWithResponse() {
	clientID := "client-1"
	requestID := "req-1"
	var receivedError error
	var receivedResponse ocpp.Response

	callback := func(confirmation ocpp.Response, err error) {
		receivedResponse = confirmation
		receivedError = err
	}

	try := func() (string, error) {
		return requestID, nil
	}

	err := suite.registry.RegisterCallback(clientID, try, callback)
	suite.Require().NoError(err)

	// Get callback and invoke with response
	mockResp := &MockResponse{Value: "success"}
	cb, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Require().True(ok)
	cb(mockResp, nil)

	suite.Assert().Equal(mockResp, receivedResponse)
	suite.Assert().Nil(receivedError)
}

// TestMultipleClientsConcurrent verifies concurrent operations on different clients
func (suite *CallbackRegistryTestSuite) TestMultipleClientsConcurrent() {
	numClients := 10
	callbacksPerClient := 5
	var wg sync.WaitGroup

	// Register callbacks for multiple clients concurrently
	for clientIndex := 0; clientIndex < numClients; clientIndex++ {
		clientID := fmt.Sprintf("client-%d", clientIndex)
		for reqIndex := 0; reqIndex < callbacksPerClient; reqIndex++ {
			wg.Add(1)
			go func(cID string, requestIndex int) {
				defer wg.Done()
				requestID := fmt.Sprintf("req-%d", requestIndex)
				callback := func(confirmation ocpp.Response, err error) {
					// Just verify it's called
				}

				requestIDCopy := requestID // Capture for closure
				try := func() (string, error) {
					return requestIDCopy, nil
				}

				err := suite.registry.RegisterCallback(cID, try, callback)
				suite.Require().NoError(err)
			}(clientID, reqIndex)
		}
	}

	wg.Wait()

	// Verify all callbacks were registered
	for clientIndex := 0; clientIndex < numClients; clientIndex++ {
		clientID := fmt.Sprintf("client-%d", clientIndex)
		dequeuedCount := 0
		for reqIndex := 0; reqIndex < callbacksPerClient; reqIndex++ {
			requestID := fmt.Sprintf("req-%d", reqIndex)
			cb, ok := suite.registry.GetCallback(clientID, requestID)
			if ok {
				cb(nil, nil)
				dequeuedCount++
			}
		}
		suite.Assert().Equal(callbacksPerClient, dequeuedCount, "All callbacks should be retrieved for client %s", clientID)
	}
}

// TestRemoveCallback verifies that RemoveCallback removes a callback without invoking it
func (suite *CallbackRegistryTestSuite) TestRemoveCallback() {
	clientID := "client-1"
	requestID := "req-1"
	callbackCalled := atomic.Bool{}

	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Store(true)
	}

	try := func() (string, error) {
		return requestID, nil
	}

	// Register callback
	err := suite.registry.RegisterCallback(clientID, try, callback)
	suite.Require().NoError(err)

	// Remove callback without invoking it
	removed := suite.registry.RemoveCallback(clientID, requestID)
	suite.Assert().True(removed)
	suite.Assert().False(callbackCalled.Load(), "Callback should not be called when removed")

	// Verify callback is gone
	_, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Assert().False(ok)

	// Try to remove non-existent callback
	removed = suite.registry.RemoveCallback(clientID, "non-existent")
	suite.Assert().False(removed)
}

// TestClearCallbacks verifies that ClearCallbacks removes all callbacks for a client
func (suite *CallbackRegistryTestSuite) TestClearCallbacks() {
	clientID := "client-1"
	numCallbacks := 10
	callbackCalled := atomic.Int64{}

	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Add(1)
	}

	// Register multiple callbacks
	for i := 0; i < numCallbacks; i++ {
		requestID := fmt.Sprintf("req-%d", i)
		requestIDCopy := requestID // Capture for closure
		try := func() (string, error) {
			return requestIDCopy, nil
		}
		err := suite.registry.RegisterCallback(clientID, try, callback)
		suite.Require().NoError(err)
	}

	// Clear all callbacks
	count := suite.registry.ClearCallbacks(clientID)
	suite.Assert().Equal(numCallbacks, count)
	suite.Assert().Equal(0, int(callbackCalled.Load()), "Callbacks should not be invoked when cleared")

	// Verify all callbacks are gone
	for i := 0; i < numCallbacks; i++ {
		requestID := fmt.Sprintf("req-%d", i)
		_, ok := suite.registry.GetCallback(clientID, requestID)
		suite.Assert().False(ok)
	}

	// Clear non-existent client
	count = suite.registry.ClearCallbacks("non-existent")
	suite.Assert().Equal(0, count)
}

// TestClearCallbacksDoesNotAffectOthers verifies that clearing one client doesn't affect others
func (suite *CallbackRegistryTestSuite) TestClearCallbacksDoesNotAffectOthers() {
	clientID1 := "client-1"
	clientID2 := "client-2"
	requestID := "req-1"

	callback1 := func(confirmation ocpp.Response, err error) {}
	callback2 := func(confirmation ocpp.Response, err error) {}

	try1 := func() (string, error) {
		return requestID, nil
	}
	try2 := func() (string, error) {
		return requestID, nil
	}

	// Register callbacks for both clients
	err := suite.registry.RegisterCallback(clientID1, try1, callback1)
	suite.Require().NoError(err)
	err = suite.registry.RegisterCallback(clientID2, try2, callback2)
	suite.Require().NoError(err)

	// Clear callbacks for clientID1
	count := suite.registry.ClearCallbacks(clientID1)
	suite.Assert().Equal(1, count)

	// clientID1 should be empty
	_, ok := suite.registry.GetCallback(clientID1, requestID)
	suite.Assert().False(ok)

	// clientID2 should still have its callback
	cb, ok := suite.registry.GetCallback(clientID2, requestID)
	suite.Require().True(ok)
	suite.Require().NotNil(cb)
}

// TestGetAllCallbacks verifies that GetAllCallbacks retrieves and removes all callbacks for a client
func (suite *CallbackRegistryTestSuite) TestGetAllCallbacks() {
	clientID := "client-1"
	numCallbacks := 5
	callbackCalled := make(map[string]*atomic.Bool)

	// Register multiple callbacks
	for i := 0; i < numCallbacks; i++ {
		requestID := fmt.Sprintf("req-%d", i)
		callbackCalled[requestID] = &atomic.Bool{}
		requestIDCopy := requestID // Capture for closure
		callback := func(confirmation ocpp.Response, err error) {
			callbackCalled[requestIDCopy].Store(true)
		}

		try := func() (string, error) {
			return requestIDCopy, nil
		}

		err := suite.registry.RegisterCallback(clientID, try, callback)
		suite.Require().NoError(err)
	}

	// Get all callbacks
	allCallbacks, found := suite.registry.GetAllCallbacks(clientID)
	suite.Require().True(found)
	suite.Assert().Len(allCallbacks, numCallbacks)

	// Verify all callbacks are present
	for i := 0; i < numCallbacks; i++ {
		requestID := fmt.Sprintf("req-%d", i)
		cb, ok := allCallbacks[RequestID(requestID)]
		suite.Require().True(ok, "Callback for requestID %s should be present", requestID)
		suite.Require().NotNil(cb)
		cb(nil, nil)
		suite.Assert().True(callbackCalled[requestID].Load())
	}

	// Registry should be empty for this client
	_, ok := suite.registry.GetCallback(clientID, "req-0")
	suite.Assert().False(ok)

	// GetAllCallbacks on empty client should return empty map and false
	allCallbacks2, found2 := suite.registry.GetAllCallbacks(clientID)
	suite.Assert().False(found2)
	suite.Assert().NotNil(allCallbacks2)
	suite.Assert().Len(allCallbacks2, 0)
}

// TestGetAllCallbacksDoesNotAffectOthers verifies that getting all callbacks for one client doesn't affect others
func (suite *CallbackRegistryTestSuite) TestGetAllCallbacksDoesNotAffectOthers() {
	clientID1 := "client-1"
	clientID2 := "client-2"
	requestID1 := "req-1"
	requestID2 := "req-2"

	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}

	callback1 := func(confirmation ocpp.Response, err error) {
		callback1Called.Store(true)
	}
	callback2 := func(confirmation ocpp.Response, err error) {
		callback2Called.Store(true)
	}

	// Register callbacks for both clients
	try1 := func() (string, error) {
		return requestID1, nil
	}
	try2 := func() (string, error) {
		return requestID2, nil
	}
	try3 := func() (string, error) {
		return requestID1, nil
	}

	err := suite.registry.RegisterCallback(clientID1, try1, callback1)
	suite.Require().NoError(err)
	err = suite.registry.RegisterCallback(clientID1, try2, callback1)
	suite.Require().NoError(err)
	err = suite.registry.RegisterCallback(clientID2, try3, callback2)
	suite.Require().NoError(err)

	// Get all callbacks for clientID1
	allCallbacks, found := suite.registry.GetAllCallbacks(clientID1)
	suite.Require().True(found)
	suite.Assert().Len(allCallbacks, 2)

	// clientID1 should be empty
	_, ok := suite.registry.GetCallback(clientID1, requestID1)
	suite.Assert().False(ok)

	// clientID2 should still have its callback
	cb, ok := suite.registry.GetCallback(clientID2, requestID1)
	suite.Require().True(ok)
	suite.Require().NotNil(cb)
	cb(nil, nil)
	suite.Assert().True(callback2Called.Load())
}

// TestRegisterCallbackEmptyIDAfterFailure verifies that client map is cleaned up when last callback fails
func (suite *CallbackRegistryTestSuite) TestRegisterCallbackEmptyIDAfterFailure() {
	clientID := "client-1"
	requestID := "req-1"

	callback := func(confirmation ocpp.Response, err error) {}

	tryFail := func() (string, error) {
		return requestID, errors.New("failed")
	}

	// Register callback that fails
	err := suite.registry.RegisterCallback(clientID, tryFail, callback)
	suite.Require().Error(err)

	// Client map should be removed
	_, ok := suite.registry.GetCallback(clientID, requestID)
	suite.Assert().False(ok)
}

func TestCallbackRegistry(t *testing.T) {
	suite.Run(t, new(CallbackRegistryTestSuite))
}
