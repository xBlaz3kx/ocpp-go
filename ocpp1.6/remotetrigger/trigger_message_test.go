package remotetrigger

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type remoteTriggerTestSuite struct {
	suite.Suite
}

func (suite *remoteTriggerTestSuite) TestTriggerMessageRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{TriggerMessageRequest{RequestedMessage: MessageTrigger(core.StatusNotificationFeatureName), ConnectorId: tests.NewInt(1)}, true},
		{TriggerMessageRequest{RequestedMessage: MessageTrigger(core.StatusNotificationFeatureName)}, true},
		{TriggerMessageRequest{}, false},
		{TriggerMessageRequest{RequestedMessage: MessageTrigger(core.StatusNotificationFeatureName), ConnectorId: tests.NewInt(0)}, true},
		{TriggerMessageRequest{RequestedMessage: MessageTrigger(core.StatusNotificationFeatureName), ConnectorId: tests.NewInt(-1)}, false},
		{TriggerMessageRequest{RequestedMessage: MessageTrigger(core.StartTransactionFeatureName)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *remoteTriggerTestSuite) TestTriggerMessageConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{TriggerMessageConfirmation{Status: TriggerMessageStatusAccepted}, true},
		{TriggerMessageConfirmation{Status: "invalidTriggerMessageStatus"}, false},
		{TriggerMessageConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *remoteTriggerTestSuite) TestTriggerMessageFeature() {
	feature := TriggerMessageFeature{}
	suite.Equal(TriggerMessageFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(TriggerMessageRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(TriggerMessageConfirmation{}), feature.GetResponseType())
}

func (suite *remoteTriggerTestSuite) TestNewTriggerMessageRequest() {
	requestedMessage := MessageTrigger(core.StatusNotificationFeatureName)
	req := NewTriggerMessageRequest(requestedMessage)
	suite.NotNil(req)
	suite.Equal(TriggerMessageFeatureName, req.GetFeatureName())
	suite.Equal(requestedMessage, req.RequestedMessage)
}

func (suite *remoteTriggerTestSuite) TestNewTriggerMessageConfirmation() {
	status := TriggerMessageStatusAccepted
	conf := NewTriggerMessageConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(TriggerMessageFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}

func TestRemoteTriggerSuite(t *testing.T) {
	suite.Run(t, new(remoteTriggerTestSuite))
}