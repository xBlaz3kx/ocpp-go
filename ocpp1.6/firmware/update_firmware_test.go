package firmware

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestUpdateFirmwareRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{UpdateFirmwareRequest{Location: "ftp:some/path", Retries: tests.NewInt(10), RetryInterval: tests.NewInt(10), RetrieveDate: types.NewDateTime(time.Now())}, true},
		{UpdateFirmwareRequest{Location: "ftp:some/path", Retries: tests.NewInt(10), RetrieveDate: types.NewDateTime(time.Now())}, true},
		{UpdateFirmwareRequest{Location: "ftp:some/path", RetrieveDate: types.NewDateTime(time.Now())}, true},
		{UpdateFirmwareRequest{}, false},
		{UpdateFirmwareRequest{Location: "ftp:some/path"}, false},
		{UpdateFirmwareRequest{Location: "invalidUri", RetrieveDate: types.NewDateTime(time.Now())}, false},
		{UpdateFirmwareRequest{Location: "ftp:some/path", Retries: tests.NewInt(-1), RetrieveDate: types.NewDateTime(time.Now())}, false},
		{UpdateFirmwareRequest{Location: "ftp:some/path", RetryInterval: tests.NewInt(-1), RetrieveDate: types.NewDateTime(time.Now())}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestUpdateFirmwareConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{UpdateFirmwareConfirmation{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestUpdateFirmwareFeature() {
	feature := UpdateFirmwareFeature{}
	suite.Equal(UpdateFirmwareFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(UpdateFirmwareRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(UpdateFirmwareConfirmation{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewUpdateFirmwareRequest() {
	location := "ftp:some/path"
	retrieveDate := types.NewDateTime(time.Now())
	req := NewUpdateFirmwareRequest(location, retrieveDate)
	suite.NotNil(req)
	suite.Equal(UpdateFirmwareFeatureName, req.GetFeatureName())
	suite.Equal(location, req.Location)
	suite.Equal(retrieveDate, req.RetrieveDate)
}

func (suite *firmwareTestSuite) TestNewUpdateFirmwareConfirmation() {
	conf := NewUpdateFirmwareConfirmation()
	suite.NotNil(conf)
	suite.Equal(UpdateFirmwareFeatureName, conf.GetFeatureName())
}