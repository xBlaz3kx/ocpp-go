package certificates

import (
	"os"
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
)

func TestMain(m *testing.M) {
	// Register hashAlgorithm validator (required for CertificateHashData)
	_ = types.Validate.RegisterValidation("hashAlgorithm", func(fl validator.FieldLevel) bool {
		algorithm := types.HashAlgorithmType(fl.Field().String())
		switch algorithm {
		case types.SHA256, types.SHA384, types.SHA512:
			return true
		default:
			return false
		}
	})
	os.Exit(m.Run())
}

type certificatesTestSuite struct {
	suite.Suite
}

func (suite *certificatesTestSuite) TestDeleteCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{DeleteCertificateRequest{}, false},
		{DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *certificatesTestSuite) TestDeleteCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{DeleteCertificateResponse{Status: DeleteCertificateStatusAccepted}, true},
		{DeleteCertificateResponse{Status: DeleteCertificateStatusFailed}, true},
		{DeleteCertificateResponse{Status: DeleteCertificateStatusNotFound}, true},
		{DeleteCertificateResponse{Status: "invalidDeleteCertificateStatus"}, false},
		{DeleteCertificateResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *certificatesTestSuite) TestDeleteCertificateFeature() {
	feature := DeleteCertificateFeature{}
	suite.Equal(DeleteCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(DeleteCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(DeleteCertificateResponse{}), feature.GetResponseType())
}

func (suite *certificatesTestSuite) TestNewDeleteCertificateRequest() {
	certificateHashData := types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	req := NewDeleteCertificateRequest(certificateHashData)
	suite.NotNil(req)
	suite.Equal(DeleteCertificateFeatureName, req.GetFeatureName())
	suite.Equal(certificateHashData, req.CertificateHashData)
}

func (suite *certificatesTestSuite) TestNewDeleteCertificateResponse() {
	status := DeleteCertificateStatusAccepted
	resp := NewDeleteCertificateResponse(status)
	suite.NotNil(resp)
	suite.Equal(DeleteCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(status, resp.Status)
}

func TestCertificatesSuite(t *testing.T) {
	suite.Run(t, new(certificatesTestSuite))
}