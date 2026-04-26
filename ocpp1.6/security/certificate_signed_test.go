package security

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type securityTestSuite struct {
	suite.Suite
}

func (suite *securityTestSuite) TestCertificateSignedRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{CertificateSignedRequest{CertificateChain: "sampleCert"}, true},
		{CertificateSignedRequest{CertificateChain: ""}, false},
		{CertificateSignedRequest{}, false},
		{CertificateSignedRequest{CertificateChain: tests.NewLongString(100001)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *securityTestSuite) TestCertificateSignedConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{CertificateSignedResponse{Status: CertificateSignedStatusAccepted}, true},
		{CertificateSignedResponse{Status: CertificateSignedStatusRejected}, true},
		{CertificateSignedResponse{Status: "invalidCertificateSignedStatus"}, false},
		{CertificateSignedResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *securityTestSuite) TestCertificateSignedFeature() {
	feature := CertificateSignedFeature{}
	suite.Equal(CertificateSignedFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(CertificateSignedRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(CertificateSignedResponse{}), feature.GetResponseType())
}

func (suite *securityTestSuite) TestNewCertificateSignedRequest() {
	certificateChain := "sampleCert"
	req := NewCertificateSignedRequest(certificateChain)
	suite.NotNil(req)
	suite.Equal(CertificateSignedFeatureName, req.GetFeatureName())
	suite.Equal(certificateChain, req.CertificateChain)
}

func (suite *securityTestSuite) TestNewCertificateSignedResponse() {
	status := CertificateSignedStatusAccepted
	resp := NewCertificateSignedResponse(status)
	suite.NotNil(resp)
	suite.Equal(CertificateSignedFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}

func TestSecuritySuite(t *testing.T) {
	suite.Run(t, new(securityTestSuite))
}