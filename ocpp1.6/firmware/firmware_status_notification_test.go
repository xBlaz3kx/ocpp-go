package firmware

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{FirmwareStatusNotificationRequest{Status: FirmwareStatusDownloaded}, true},
		{FirmwareStatusNotificationRequest{}, false},
		{FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{FirmwareStatusNotificationConfirmation{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestFirmwareStatusNotificationFeature() {
	feature := FirmwareStatusNotificationFeature{}
	suite.Equal(FirmwareStatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(FirmwareStatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(FirmwareStatusNotificationConfirmation{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewFirmwareStatusNotificationRequest() {
	status := FirmwareStatusDownloaded
	req := NewFirmwareStatusNotificationRequest(status)
	suite.NotNil(req)
	suite.Equal(FirmwareStatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(status, req.Status)
}

func (suite *firmwareTestSuite) TestNewFirmwareStatusNotificationConfirmation() {
	conf := NewFirmwareStatusNotificationConfirmation()
	suite.NotNil(conf)
	suite.Equal(FirmwareStatusNotificationFeatureName, conf.GetFeatureName())
}