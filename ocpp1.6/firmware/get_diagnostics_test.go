package firmware

import (
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type firmwareTestSuite struct {
	suite.Suite
}

func (suite *firmwareTestSuite) TestGetDiagnosticsRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{GetDiagnosticsRequest{Location: "ftp:some/path", Retries: tests.NewInt(10), RetryInterval: tests.NewInt(10), StartTime: types.NewDateTime(time.Now()), StopTime: types.NewDateTime(time.Now())}, true},
		{GetDiagnosticsRequest{Location: "ftp:some/path", Retries: tests.NewInt(10), RetryInterval: tests.NewInt(10), StartTime: types.NewDateTime(time.Now())}, true},
		{GetDiagnosticsRequest{Location: "ftp:some/path", Retries: tests.NewInt(10), RetryInterval: tests.NewInt(10)}, true},
		{GetDiagnosticsRequest{Location: "ftp:some/path", Retries: tests.NewInt(10)}, true},
		{GetDiagnosticsRequest{Location: "ftp:some/path"}, true},
		{GetDiagnosticsRequest{}, false},
		{GetDiagnosticsRequest{Location: "invalidUri"}, false},
		{GetDiagnosticsRequest{Location: "ftp:some/path", Retries: tests.NewInt(-1)}, false},
		{GetDiagnosticsRequest{Location: "ftp:some/path", RetryInterval: tests.NewInt(-1)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestGetDiagnosticsConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{GetDiagnosticsConfirmation{FileName: "someFileName"}, true},
		{GetDiagnosticsConfirmation{FileName: ""}, true},
		{GetDiagnosticsConfirmation{}, true},
		{GetDiagnosticsConfirmation{FileName: ">255............................................................................................................................................................................................................................................................"}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestGetDiagnosticsFeature() {
	feature := GetDiagnosticsFeature{}
	suite.Equal(GetDiagnosticsFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetDiagnosticsRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetDiagnosticsConfirmation{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewGetDiagnosticsRequest() {
	location := "ftp:some/path"
	req := NewGetDiagnosticsRequest(location)
	suite.NotNil(req)
	suite.Equal(GetDiagnosticsFeatureName, req.GetFeatureName())
	suite.Equal(location, req.Location)
}

func (suite *firmwareTestSuite) TestNewGetDiagnosticsConfirmation() {
	conf := NewGetDiagnosticsConfirmation()
	suite.NotNil(conf)
	suite.Equal(GetDiagnosticsFeatureName, conf.GetFeatureName())
}

func TestFirmwareSuite(t *testing.T) {
	suite.Run(t, new(firmwareTestSuite))
}