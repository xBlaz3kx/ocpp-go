package security

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *securityTestSuite) TestSecurityEventNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{SecurityEventNotificationRequest{Type: "someType", Timestamp: types.NewDateTime(time.Now())}, true},
		{SecurityEventNotificationRequest{Type: "someType", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someInfo"}, true},
		{SecurityEventNotificationRequest{Type: "someType"}, false},
		{SecurityEventNotificationRequest{Timestamp: types.NewDateTime(time.Now())}, false},
		{SecurityEventNotificationRequest{}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *securityTestSuite) TestSecurityEventNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SecurityEventNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *securityTestSuite) TestSecurityEventNotificationFeature() {
	feature := SecurityEventNotificationFeature{}
	suite.Equal(SecurityEventNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SecurityEventNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SecurityEventNotificationResponse{}), feature.GetResponseType())
}

func (suite *securityTestSuite) TestNewSecurityEventNotificationRequest() {
	typ := "someType"
	timestamp := types.NewDateTime(time.Now())
	req := NewSecurityEventNotificationRequest(typ, timestamp)
	suite.NotNil(req)
	suite.Equal(SecurityEventNotificationFeatureName, req.GetFeatureName())
	suite.Equal(typ, req.Type)
	suite.Equal(timestamp, req.Timestamp)
}

func (suite *securityTestSuite) TestNewSecurityEventNotificationResponse() {
	resp := NewSecurityEventNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(SecurityEventNotificationFeatureName, resp.GetFeatureName())
}