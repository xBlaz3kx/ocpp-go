package certificates

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *certificatesTestSuite) TestInstallCertificateRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate}, false},
		{InstallCertificateRequest{CertificateType: types.CentralSystemRootCertificate, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.CentralSystemRootCertificate, Certificate: ""}, false},
		{InstallCertificateRequest{Certificate: "0xdeadbeef"}, false},
		{InstallCertificateRequest{}, false},
		{InstallCertificateRequest{CertificateType: "invalidCertificateUse", Certificate: "0xdeadbeef"}, false},
		{InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate, Certificate: tests.NewLongString(5501)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *certificatesTestSuite) TestInstallCertificateConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{InstallCertificateResponse{Status: CertificateStatusAccepted}, true},
		{InstallCertificateResponse{Status: CertificateStatusRejected}, true},
		{InstallCertificateResponse{Status: CertificateStatusFailed}, true},
		{InstallCertificateResponse{}, false},
		{InstallCertificateResponse{Status: "invalidInstallCertificateStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *certificatesTestSuite) TestInstallCertificateFeature() {
	feature := InstallCertificateFeature{}
	suite.Equal(InstallCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(InstallCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(InstallCertificateResponse{}), feature.GetResponseType())
}

func (suite *certificatesTestSuite) TestNewInstallCertificateRequest() {
	certificateType := types.CentralSystemRootCertificate
	certificate := "0xdeadbeef"
	req := NewInstallCertificateRequest(certificateType, certificate)
	suite.NotNil(req)
	suite.Equal(InstallCertificateFeatureName, req.GetFeatureName())
	suite.Equal(certificateType, req.CertificateType)
	suite.Equal(certificate, req.Certificate)
}

func (suite *certificatesTestSuite) TestNewInstallCertificateResponse() {
	status := CertificateStatusAccepted
	resp := NewInstallCertificateResponse(status)
	suite.NotNil(resp)
	suite.Equal(InstallCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}