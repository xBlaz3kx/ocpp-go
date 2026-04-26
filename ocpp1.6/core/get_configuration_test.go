package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestGetConfigurationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetConfigurationRequest{Key: []string{"key1", "key2"}}, true},
		{GetConfigurationRequest{Key: []string{"key1", "key2", "key3", "key4", "key5", "key6"}}, true},
		{GetConfigurationRequest{Key: []string{"key1", "key2", "key2"}}, false},
		{GetConfigurationRequest{}, true},
		{GetConfigurationRequest{Key: []string{}}, true},
		{GetConfigurationRequest{Key: []string{">50................................................"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestGetConfigurationConfirmationValidation() {
	t := suite.T()
	value1 := "value1"
	value2 := "value2"
	longValue := ">500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."
	var confirmationTable = []tests.GenericTestEntry{
		{GetConfigurationConfirmation{ConfigurationKey: []ConfigurationKey{{Key: "key1", Readonly: true, Value: &value1}}}, true},
		{GetConfigurationConfirmation{ConfigurationKey: []ConfigurationKey{{Key: "key1", Readonly: true, Value: &value1}, {Key: "key2", Readonly: false, Value: &value2}}}, true},
		{GetConfigurationConfirmation{ConfigurationKey: []ConfigurationKey{{Key: "key1", Readonly: true, Value: &value1}}, UnknownKey: []string{"keyX"}}, true},
		{GetConfigurationConfirmation{ConfigurationKey: []ConfigurationKey{{Key: "key1", Readonly: false, Value: &value1}}, UnknownKey: []string{"keyX", "keyY"}}, true},
		{GetConfigurationConfirmation{UnknownKey: []string{"keyX"}}, true},
		{GetConfigurationConfirmation{UnknownKey: []string{">50................................................"}}, false},
		{GetConfigurationConfirmation{ConfigurationKey: []ConfigurationKey{{Key: ">50................................................", Readonly: true, Value: &value1}}}, false},
		{GetConfigurationConfirmation{ConfigurationKey: []ConfigurationKey{{Key: "key1", Readonly: true, Value: &longValue}}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestGetConfigurationFeature() {
	feature := GetConfigurationFeature{}
	suite.Equal(GetConfigurationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetConfigurationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetConfigurationConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewGetConfigurationRequest() {
	keys := []string{"key1", "key2"}
	req := NewGetConfigurationRequest(keys)
	suite.NotNil(req)
	suite.Equal(GetConfigurationFeatureName, req.GetFeatureName())
	suite.Equal(keys, req.Key)
}

func (suite *coreTestSuite) TestNewGetConfigurationConfirmation() {
	value := "value1"
	configKeys := []ConfigurationKey{{Key: "key1", Readonly: true, Value: &value}}
	conf := NewGetConfigurationConfirmation(configKeys)
	suite.NotNil(conf)
	suite.Equal(GetConfigurationFeatureName, conf.GetFeatureName())
	suite.Equal(configKeys, conf.ConfigurationKey)
}