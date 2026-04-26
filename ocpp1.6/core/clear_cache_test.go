package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestClearCacheRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ClearCacheRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestClearCacheConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ClearCacheConfirmation{Status: ClearCacheStatusAccepted}, true},
		{ClearCacheConfirmation{Status: ClearCacheStatusRejected}, true},
		{ClearCacheConfirmation{Status: "invalidClearCacheStatus"}, false},
		{ClearCacheConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestClearCacheFeature() {
	feature := ClearCacheFeature{}
	suite.Equal(ClearCacheFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearCacheRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearCacheConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewClearCacheRequest() {
	req := NewClearCacheRequest()
	suite.NotNil(req)
	suite.Equal(ClearCacheFeatureName, req.GetFeatureName())
}

func (suite *coreTestSuite) TestNewClearCacheConfirmation() {
	status := ClearCacheStatusAccepted
	conf := NewClearCacheConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(ClearCacheFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}