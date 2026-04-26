package firmware

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestDiagnosticsStatusNotificationRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{DiagnosticsStatusNotificationRequest{Status: DiagnosticsStatusUploaded}, true},
		{DiagnosticsStatusNotificationRequest{}, false},
		{DiagnosticsStatusNotificationRequest{Status: "invalidDiagnosticsStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestDiagnosticsStatusNotificationConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{DiagnosticsStatusNotificationConfirmation{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestDiagnosticsStatusNotificationFeature() {
	feature := DiagnosticsStatusNotificationFeature{}
	suite.Equal(DiagnosticsStatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(DiagnosticsStatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(DiagnosticsStatusNotificationConfirmation{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewDiagnosticsStatusNotificationRequest() {
	status := DiagnosticsStatusUploaded
	req := NewDiagnosticsStatusNotificationRequest(status)
	suite.NotNil(req)
	suite.Equal(DiagnosticsStatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(status, req.Status)
}

func (suite *firmwareTestSuite) TestNewDiagnosticsStatusNotificationConfirmation() {
	conf := NewDiagnosticsStatusNotificationConfirmation()
	suite.NotNil(conf)
	suite.Equal(DiagnosticsStatusNotificationFeatureName, conf.GetFeatureName())
}