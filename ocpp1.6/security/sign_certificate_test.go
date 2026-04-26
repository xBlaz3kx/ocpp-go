package security

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *securityTestSuite) TestSignCertificateRequestValidation() {
	var requestTable = []tests.GenericTestEntry{
		{SignCertificateRequest{CSR: "deadc0de", CertificateType: types.ChargingStationCert}, true},
		{SignCertificateRequest{CSR: "deadc0de"}, true},
		{SignCertificateRequest{}, false},
		{SignCertificateRequest{CSR: "deadc0de", CertificateType: "invalidType"}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *securityTestSuite) TestSignCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SignCertificateResponse{Status: types.GenericStatusAccepted}, true},
		{SignCertificateResponse{}, false},
		{SignCertificateResponse{Status: "invalidStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *securityTestSuite) TestSignCertificateFeature() {
	feature := SignCertificateFeature{}
	suite.Equal(SignCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SignCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SignCertificateResponse{}), feature.GetResponseType())
}

func (suite *securityTestSuite) TestNewSignCertificateRequest() {
	csr := "deadc0de"
	req := NewSignCertificateRequest(csr)
	suite.NotNil(req)
	suite.Equal(SignCertificateFeatureName, req.GetFeatureName())
	suite.Equal(csr, req.CSR)
}

func (suite *securityTestSuite) TestNewSignCertificateResponse() {
	status := types.GenericStatusAccepted
	resp := NewSignCertificateResponse(status)
	suite.NotNil(resp)
	suite.Equal(SignCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}