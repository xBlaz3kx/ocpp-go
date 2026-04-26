package logging

import (
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type loggingTestSuite struct {
	suite.Suite
}

func (suite *loggingTestSuite) TestGetLogRequestValidation() {
	t := suite.T()
	logParams := LogParameters{
		RemoteLocation:  "ftp://some/path",
		OldestTimestamp: types.NewDateTime(time.Now()),
		LatestTimestamp: types.NewDateTime(time.Now()),
	}
	var requestTable = []tests.GenericTestEntry{
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(3), RetryInterval: tests.NewInt(60), Log: logParams}, true},
		{GetLogRequest{LogType: LogTypeSecurity, RequestID: 1, Retries: tests.NewInt(3), RetryInterval: tests.NewInt(60), Log: logParams}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 0, Log: logParams}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: logParams}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: LogParameters{RemoteLocation: "ftp://some/path"}}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: LogParameters{RemoteLocation: "ftp://some/path", OldestTimestamp: types.NewDateTime(time.Now())}}, true},
		{GetLogRequest{}, false},
		{GetLogRequest{LogType: "invalidLogType", RequestID: 1, Log: logParams}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: -1, Log: logParams}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(-1), Log: logParams}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, RetryInterval: tests.NewInt(-1), Log: logParams}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: LogParameters{}}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: LogParameters{RemoteLocation: "notAValidURL"}}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: LogParameters{RemoteLocation: tests.NewLongString(513)}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *loggingTestSuite) TestGetLogResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{GetLogResponse{Status: LogStatusAccepted}, true},
		{GetLogResponse{Status: LogStatusRejected}, true},
		{GetLogResponse{Status: LogStatusAcceptedCanceled}, true},
		{GetLogResponse{Status: LogStatusAccepted, Filename: "logfile.txt"}, true},
		{GetLogResponse{Status: LogStatusAccepted, Filename: tests.NewLongString(256)}, true},
		{GetLogResponse{Status: LogStatusAccepted, Filename: tests.NewLongString(257)}, false},
		{GetLogResponse{}, false},
		{GetLogResponse{Status: "invalidLogStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *loggingTestSuite) TestGetLogFeature() {
	feature := GetLogFeature{}
	suite.Equal(GetLogFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetLogRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetLogResponse{}), feature.GetResponseType())
}

func (suite *loggingTestSuite) TestNewGetLogRequest() {
	logType := LogTypeDiagnostics
	requestID := 42
	logParams := LogParameters{RemoteLocation: "ftp://some/path"}
	req := NewGetLogRequest(logType, requestID, logParams)
	suite.NotNil(req)
	suite.Equal(GetLogFeatureName, req.GetFeatureName())
	suite.Equal(logType, req.LogType)
	suite.Equal(requestID, req.RequestID)
	suite.Equal(logParams, req.Log)
}

func (suite *loggingTestSuite) TestNewGetLogResponse() {
	status := LogStatusAccepted
	resp := NewGetLogResponse(status)
	suite.NotNil(resp)
	suite.Equal(GetLogFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}

func TestLoggingSuite(t *testing.T) {
	suite.Run(t, new(loggingTestSuite))
}
