package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestUnlockConnectorRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{UnlockConnectorRequest{ConnectorId: 1}, true},
		{UnlockConnectorRequest{ConnectorId: -1}, false},
		{UnlockConnectorRequest{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *coreTestSuite) TestUnlockConnectorConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{UnlockConnectorConfirmation{Status: UnlockStatusUnlocked}, true},
		{UnlockConnectorConfirmation{Status: "invalidUnlockStatus"}, false},
		{UnlockConnectorConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *coreTestSuite) TestUnlockConnectorFeature() {
	feature := UnlockConnectorFeature{}
	suite.Equal(UnlockConnectorFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(UnlockConnectorRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(UnlockConnectorConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewUnlockConnectorRequest() {
	connectorId := 1
	req := NewUnlockConnectorRequest(connectorId)
	suite.NotNil(req)
	suite.Equal(UnlockConnectorFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
}

func (suite *coreTestSuite) TestNewUnlockConnectorConfirmation() {
	status := UnlockStatusUnlocked
	conf := NewUnlockConnectorConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(UnlockConnectorFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}