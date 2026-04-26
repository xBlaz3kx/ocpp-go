package securefirmware

import (
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type secureFirmwareTestSuite struct {
	suite.Suite
}

func (suite *secureFirmwareTestSuite) TestSignedUpdateFirmwareRequestValidation() {
	t := suite.T()
	fw := Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	var requestTable = []tests.GenericTestEntry{
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: fw}, true},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RequestID: 42, Firmware: fw}, true},
		{SignedUpdateFirmwareRequest{RequestID: 42, Firmware: fw}, true},
		{SignedUpdateFirmwareRequest{Firmware: fw}, true},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{Location: "https://someurl", RetrieveDateTime: types.NewDateTime(time.Now())}}, true},
		{SignedUpdateFirmwareRequest{}, false},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(-1), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: fw}, false},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(-1), RequestID: 42, Firmware: fw}, false},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: -1, Firmware: fw}, false},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{Location: "https://someurl", InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{SignedUpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{Location: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *secureFirmwareTestSuite) TestSignedUpdateFirmwareResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{SignedUpdateFirmwareResponse{Status: UpdateFirmwareStatusAccepted}, true},
		{SignedUpdateFirmwareResponse{Status: UpdateFirmwareStatusRejected}, true},
		{SignedUpdateFirmwareResponse{Status: UpdateFirmwareStatusAccepted}, true},
		{SignedUpdateFirmwareResponse{}, false},
		{SignedUpdateFirmwareResponse{Status: "invalidFirmwareUpdateStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *secureFirmwareTestSuite) TestSignedUpdateFirmwareFeature() {
	feature := SignedUpdateFirmwareFeature{}
	suite.Equal(SignedUpdateFirmwareFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SignedUpdateFirmwareRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SignedUpdateFirmwareResponse{}), feature.GetResponseType())
}

func (suite *secureFirmwareTestSuite) TestNewSignedUpdateFirmwareRequest() {
	requestId := 42
	fw := Firmware{
		Location:         "https://someurl",
		RetrieveDateTime: types.NewDateTime(time.Now()),
	}
	req := NewSignedUpdateFirmwareRequest(requestId, fw)
	suite.NotNil(req)
	suite.Equal(SignedUpdateFirmwareFeatureName, req.GetFeatureName())
	suite.Equal(requestId, req.RequestID)
	suite.Equal(fw, req.Firmware)
}

func (suite *secureFirmwareTestSuite) TestNewSignedUpdateFirmwareResponse() {
	status := UpdateFirmwareStatusAccepted
	resp := NewSignedUpdateFirmwareResponse(status)
	suite.NotNil(resp)
	suite.Equal(SignedUpdateFirmwareFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}

func TestSecureFirmwareSuite(t *testing.T) {
	suite.Run(t, new(secureFirmwareTestSuite))
}