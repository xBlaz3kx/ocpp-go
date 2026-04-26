package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV16TestSuite) TestRemoteStopTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	transactionId := 1
	status := types.RemoteStartStopStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":%v}]`, messageId, core.RemoteStopTransactionFeatureName, transactionId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	RemoteStopTransactionConfirmation := core.NewRemoteStopTransactionConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockChargePointCoreListener{}
	coreListener.On("OnRemoteStopTransaction", mock.Anything).Return(RemoteStopTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*core.RemoteStopTransactionRequest)
		require.NotNil(t, request)
		require.True(t, ok)
		assert.Equal(t, transactionId, request.TransactionId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.RemoteStopTransaction(wsId, func(confirmation *core.RemoteStopTransactionConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, transactionId)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestRemoteStopTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	transactionId := 1
	RemoteStopTransactionRequest := core.NewRemoteStopTransactionRequest(transactionId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":%v}]`, messageId, core.RemoteStopTransactionFeatureName, transactionId)
	testUnsupportedRequestFromChargePoint(suite, RemoteStopTransactionRequest, requestJson, messageId)
}
