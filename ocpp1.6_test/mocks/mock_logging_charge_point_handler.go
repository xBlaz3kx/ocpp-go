// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	logging "github.com/lorenzodonini/ocpp-go/ocpp1.6/logging"
	mock "github.com/stretchr/testify/mock"
)

// MockLogChargePointHandler is an autogenerated mock type for the ChargePointHandler type
type MockLogChargePointHandler struct {
	mock.Mock
}

type MockLogChargePointHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogChargePointHandler) EXPECT() *MockLogChargePointHandler_Expecter {
	return &MockLogChargePointHandler_Expecter{mock: &_m.Mock}
}

// OnGetLog provides a mock function with given fields: request
func (_m *MockLogChargePointHandler) OnGetLog(request *logging.GetLogRequest) (*logging.GetLogResponse, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for OnGetLog")
	}

	var r0 *logging.GetLogResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*logging.GetLogRequest) (*logging.GetLogResponse, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*logging.GetLogRequest) *logging.GetLogResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*logging.GetLogResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*logging.GetLogRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLogChargePointHandler_OnGetLog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnGetLog'
type MockLogChargePointHandler_OnGetLog_Call struct {
	*mock.Call
}

// OnGetLog is a helper method to define mock.On call
//   - request *logging.GetLogRequest
func (_e *MockLogChargePointHandler_Expecter) OnGetLog(request interface{}) *MockLogChargePointHandler_OnGetLog_Call {
	return &MockLogChargePointHandler_OnGetLog_Call{Call: _e.mock.On("OnGetLog", request)}
}

func (_c *MockLogChargePointHandler_OnGetLog_Call) Run(run func(request *logging.GetLogRequest)) *MockLogChargePointHandler_OnGetLog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*logging.GetLogRequest))
	})
	return _c
}

func (_c *MockLogChargePointHandler_OnGetLog_Call) Return(response *logging.GetLogResponse, err error) *MockLogChargePointHandler_OnGetLog_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockLogChargePointHandler_OnGetLog_Call) RunAndReturn(run func(*logging.GetLogRequest) (*logging.GetLogResponse, error)) *MockLogChargePointHandler_OnGetLog_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLogChargePointHandler creates a new instance of MockLogChargePointHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLogChargePointHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLogChargePointHandler {
	mock := &MockLogChargePointHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}