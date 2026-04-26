package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestChangeConfigurationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ChangeConfigurationRequest{Key: "someKey", Value: "someValue"}, true},
		{ChangeConfigurationRequest{Key: "someKey"}, false},
		{ChangeConfigurationRequest{Value: "someValue"}, false},
		{ChangeConfigurationRequest{}, false},
		{ChangeConfigurationRequest{Key: ">50................................................", Value: "someValue"}, false},
		{ChangeConfigurationRequest{Key: "someKey", Value: ">500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestChangeConfigurationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ChangeConfigurationConfirmation{Status: ConfigurationStatusAccepted}, true},
		{ChangeConfigurationConfirmation{Status: ConfigurationStatusRejected}, true},
		{ChangeConfigurationConfirmation{Status: ConfigurationStatusRebootRequired}, true},
		{ChangeConfigurationConfirmation{Status: ConfigurationStatusNotSupported}, true},
		{ChangeConfigurationConfirmation{Status: "invalidConfigurationStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestChangeConfigurationFeature() {
	feature := ChangeConfigurationFeature{}
	suite.Equal(ChangeConfigurationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ChangeConfigurationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ChangeConfigurationConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewChangeConfigurationRequest() {
	key := "someKey"
	value := "someValue"
	req := NewChangeConfigurationRequest(key, value)
	suite.NotNil(req)
	suite.Equal(ChangeConfigurationFeatureName, req.GetFeatureName())
	suite.Equal(key, req.Key)
	suite.Equal(value, req.Value)
}

func (suite *coreTestSuite) TestNewChangeConfigurationConfirmation() {
	status := ConfigurationStatusAccepted
	conf := NewChangeConfigurationConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(ChangeConfigurationFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}