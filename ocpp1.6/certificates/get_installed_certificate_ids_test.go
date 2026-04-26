package certificates

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *certificatesTestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{GetInstalledCertificateIdsRequest{CertificateType: types.CentralSystemRootCertificate}, true},
		{GetInstalledCertificateIdsRequest{}, false},
		{GetInstalledCertificateIdsRequest{CertificateType: "invalidCertificateUse"}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *certificatesTestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusAccepted}, true},
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusNotFound}, true},
		{GetInstalledCertificateIdsResponse{}, false},
		{GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *certificatesTestSuite) TestGetInstalledCertificateIdsFeature() {
	feature := GetInstalledCertificateIdsFeature{}
	suite.Equal(GetInstalledCertificateIdsFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetInstalledCertificateIdsRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetInstalledCertificateIdsResponse{}), feature.GetResponseType())
}

func (suite *certificatesTestSuite) TestNewGetInstalledCertificateIdsRequest() {
	certificateType := types.CentralSystemRootCertificate
	req := NewGetInstalledCertificateIdsRequest(certificateType)
	suite.NotNil(req)
	suite.Equal(GetInstalledCertificateIdsFeatureName, req.GetFeatureName())
	suite.Equal(certificateType, req.CertificateType)
}

func (suite *certificatesTestSuite) TestNewGetInstalledCertificateIdsResponse() {
	status := GetInstalledCertificateStatusAccepted
	resp := NewGetInstalledCertificateIdsResponse(status)
	suite.NotNil(resp)
	suite.Equal(GetInstalledCertificateIdsFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}