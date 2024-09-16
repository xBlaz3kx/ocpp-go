package ocpp16_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{certificates.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CentralSystemRootCertificate}}, true},
		{certificates.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.ManufacturerRootCertificate}}, true},
		{certificates.GetInstalledCertificateIdsRequest{}, true},
		{certificates.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{"invalidCertificateUse"}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusAccepted}, true},
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusNotFound}, true},
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusAccepted}, true},
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusAccepted}, true},
		{certificates.GetInstalledCertificateIdsResponse{}, false},
		{certificates.GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusAccepted}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateTypes := []types.CertificateUse{types.CentralSystemRootCertificate}
	status := certificates.GetInstalledCertificateStatusAccepted

	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":["%v"]}]`, messageId, certificates.GetInstalledCertificateIdsFeatureName, certificateTypes[0])
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`,
		messageId, status)
	getInstalledCertificateIdsConfirmation := certificates.NewGetInstalledCertificateIdsResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := &MockChargePointCertificateHandler{}
	handler.On("OnGetInstalledCertificateIds", mock.Anything).Return(getInstalledCertificateIdsConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*certificates.GetInstalledCertificateIdsRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateTypes, request.CertificateTypes)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	suite.chargePoint.SetCertificateHandler(handler)
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetInstalledCertificateIds(wsId, func(confirmation *certificates.GetInstalledCertificateIdsResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *certificates.GetInstalledCertificateIdsRequest) {
		request.CertificateTypes = certificateTypes
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsInvalidEndpoint() {
	messageId := defaultMessageId
	certificateTypes := []types.CertificateUse{types.CentralSystemRootCertificate}
	GetInstalledCertificateIdsRequest := certificates.NewGetInstalledCertificateIdsRequest()
	GetInstalledCertificateIdsRequest.CertificateTypes = certificateTypes
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":["%v"]}]`, messageId, certificates.GetInstalledCertificateIdsFeatureName, certificateTypes[0])
	testUnsupportedRequestFromCentralSystem(suite, GetInstalledCertificateIdsRequest, requestJson, messageId)
}
