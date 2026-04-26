package logging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *loggingTestSuite) TestLogStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{LogStatusNotificationRequest{Status: UploadLogStatusUploaded, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusUploading, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusUploadFailure, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusIdle, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusBadMessage, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusNotSupportedOp, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusPermissionDenied, RequestID: 1}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusUploaded, RequestID: 0}, true},
		{LogStatusNotificationRequest{}, false},
		{LogStatusNotificationRequest{Status: "invalidUploadLogStatus", RequestID: 1}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *loggingTestSuite) TestLogStatusNotificationResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{LogStatusNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *loggingTestSuite) TestLogStatusNotificationFeature() {
	feature := LogStatusNotificationFeature{}
	suite.Equal(LogStatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(LogStatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(LogStatusNotificationResponse{}), feature.GetResponseType())
}

func (suite *loggingTestSuite) TestNewLogStatusNotificationRequest() {
	status := UploadLogStatusUploading
	requestID := 5
	req := NewLogStatusNotificationRequest(status, requestID)
	suite.NotNil(req)
	suite.Equal(LogStatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(status, req.Status)
	suite.Equal(requestID, req.RequestID)
}

func (suite *loggingTestSuite) TestNewLogStatusNotificationResponse() {
	resp := NewLogStatusNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(LogStatusNotificationFeatureName, resp.GetFeatureName())
}