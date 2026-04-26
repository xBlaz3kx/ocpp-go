package extendedtriggermessage

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type extendedTriggerMessageTestSuite struct {
	suite.Suite
}

func (suite *extendedTriggerMessageTestSuite) TestExtendedTriggerMessageRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeBootNotification}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeLogStatusNotification}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeHeartbeat}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeMeterValues}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeSignChargingStationCertificate}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeFirmwareStatusNotification}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeStatusNotification}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeBootNotification, ConnectorId: tests.NewInt(1)}, true},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeBootNotification, ConnectorId: tests.NewInt(0)}, true},
		{ExtendedTriggerMessageRequest{}, false},
		{ExtendedTriggerMessageRequest{RequestedMessage: "invalidMessageType"}, false},
		{ExtendedTriggerMessageRequest{RequestedMessage: ExtendedTriggerMessageTypeBootNotification, ConnectorId: tests.NewInt(-1)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *extendedTriggerMessageTestSuite) TestExtendedTriggerMessageResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{ExtendedTriggerMessageResponse{Status: ExtendedTriggerMessageStatusAccepted}, true},
		{ExtendedTriggerMessageResponse{Status: ExtendedTriggerMessageStatusRejected}, true},
		{ExtendedTriggerMessageResponse{Status: ExtendedTriggerMessageStatusNotImplemented}, true},
		{ExtendedTriggerMessageResponse{}, false},
		{ExtendedTriggerMessageResponse{Status: "invalidStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *extendedTriggerMessageTestSuite) TestExtendedTriggerMessageFeature() {
	feature := ExtendedTriggerMessageFeature{}
	suite.Equal(ExtendedTriggerMessageFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ExtendedTriggerMessageRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ExtendedTriggerMessageResponse{}), feature.GetResponseType())
}

func (suite *extendedTriggerMessageTestSuite) TestNewExtendedTriggerMessageRequest() {
	msgType := ExtendedTriggerMessageTypeHeartbeat
	req := NewExtendedTriggerMessageRequest(msgType)
	suite.NotNil(req)
	suite.Equal(ExtendedTriggerMessageFeatureName, req.GetFeatureName())
	suite.Equal(msgType, req.RequestedMessage)
}

func (suite *extendedTriggerMessageTestSuite) TestNewExtendedTriggerMessageResponse() {
	status := ExtendedTriggerMessageStatusAccepted
	resp := NewExtendedTriggerMessageResponse(status)
	suite.NotNil(resp)
	suite.Equal(ExtendedTriggerMessageFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}

func TestExtendedTriggerMessageSuite(t *testing.T) {
	suite.Run(t, new(extendedTriggerMessageTestSuite))
}