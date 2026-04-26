package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestResetRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ResetRequest{Type: ResetTypeHard}, true},
		{ResetRequest{Type: ResetTypeSoft}, true},
		{ResetRequest{Type: "invalidResetType"}, false},
		{ResetRequest{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *coreTestSuite) TestResetConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ResetConfirmation{Status: ResetStatusAccepted}, true},
		{ResetConfirmation{Status: ResetStatusRejected}, true},
		{ResetConfirmation{Status: "invalidResetStatus"}, false},
		{ResetConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *coreTestSuite) TestResetFeature() {
	feature := ResetFeature{}
	suite.Equal(ResetFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ResetRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ResetConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewResetRequest() {
	resetType := ResetTypeSoft
	req := NewResetRequest(resetType)
	suite.NotNil(req)
	suite.Equal(ResetFeatureName, req.GetFeatureName())
	suite.Equal(resetType, req.Type)
}

func (suite *coreTestSuite) TestNewResetConfirmation() {
	status := ResetStatusAccepted
	conf := NewResetConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(ResetFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}