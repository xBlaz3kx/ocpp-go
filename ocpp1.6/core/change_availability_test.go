package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ChangeAvailabilityRequest{ConnectorId: 0, Type: AvailabilityTypeOperative}, true},
		{ChangeAvailabilityRequest{ConnectorId: 0, Type: AvailabilityTypeInoperative}, true},
		{ChangeAvailabilityRequest{ConnectorId: 0}, false},
		{ChangeAvailabilityRequest{Type: AvailabilityTypeOperative}, true},
		{ChangeAvailabilityRequest{Type: "invalidAvailabilityType"}, false},
		{ChangeAvailabilityRequest{ConnectorId: -1, Type: AvailabilityTypeOperative}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *coreTestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ChangeAvailabilityConfirmation{Status: AvailabilityStatusAccepted}, true},
		{ChangeAvailabilityConfirmation{Status: AvailabilityStatusRejected}, true},
		{ChangeAvailabilityConfirmation{Status: AvailabilityStatusScheduled}, true},
		{ChangeAvailabilityConfirmation{Status: "invalidAvailabilityStatus"}, false},
		{ChangeAvailabilityConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *coreTestSuite) TestChangeAvailabilityFeature() {
	feature := ChangeAvailabilityFeature{}
	suite.Equal(ChangeAvailabilityFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ChangeAvailabilityRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ChangeAvailabilityConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewChangeAvailabilityRequest() {
	connectorId := 1
	availabilityType := AvailabilityTypeOperative
	req := NewChangeAvailabilityRequest(connectorId, availabilityType)
	suite.NotNil(req)
	suite.Equal(ChangeAvailabilityFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(availabilityType, req.Type)
}

func (suite *coreTestSuite) TestNewChangeAvailabilityConfirmation() {
	status := AvailabilityStatusAccepted
	conf := NewChangeAvailabilityConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(ChangeAvailabilityFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}