package localauth

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type localAuthTestSuite struct {
	suite.Suite
}

func (suite *localAuthTestSuite) TestGetLocalListVersionRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{GetLocalListVersionRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *localAuthTestSuite) TestGetLocalListVersionConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{GetLocalListVersionConfirmation{ListVersion: 1}, true},
		{GetLocalListVersionConfirmation{ListVersion: 0}, true},
		{GetLocalListVersionConfirmation{}, true},
		{GetLocalListVersionConfirmation{ListVersion: -1}, true},
		{GetLocalListVersionConfirmation{ListVersion: -2}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *localAuthTestSuite) TestGetLocalListVersionFeature() {
	feature := GetLocalListVersionFeature{}
	suite.Equal(GetLocalListVersionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetLocalListVersionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetLocalListVersionConfirmation{}), feature.GetResponseType())
}

func (suite *localAuthTestSuite) TestNewGetLocalListVersionRequest() {
	req := NewGetLocalListVersionRequest()
	suite.NotNil(req)
	suite.Equal(GetLocalListVersionFeatureName, req.GetFeatureName())
}

func (suite *localAuthTestSuite) TestNewGetLocalListVersionConfirmation() {
	version := 1
	conf := NewGetLocalListVersionConfirmation(version)
	suite.NotNil(conf)
	suite.Equal(GetLocalListVersionFeatureName, conf.GetFeatureName())
	suite.Equal(version, conf.ListVersion)
}

func TestLocalAuthSuite(t *testing.T) {
	suite.Run(t, new(localAuthTestSuite))
}