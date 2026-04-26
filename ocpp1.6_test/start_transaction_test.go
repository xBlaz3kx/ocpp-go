package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV16TestSuite) TestStartTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	meterStart := 100
	reservationId := tests.NewInt(42)
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	parentIdTag := "parentTag1"
	status := types.AuthorizationStatusAccepted
	expiryDate := types.NewDateTime(time.Now().Add(time.Hour * 8))
	transactionId := 16
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","meterStart":%v,"reservationId":%v,"timestamp":"%v"}]`,
		messageId, core.StartTransactionFeatureName, connectorId, idTag, meterStart, *reservationId, timestamp.FormatTimestamp())
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"},"transactionId":%v}]`, messageId, expiryDate.FormatTimestamp(), parentIdTag, status, transactionId)
	startTransactionConfirmation := core.NewStartTransactionConfirmation(&types.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status}, transactionId)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnStartTransaction", mock.AnythingOfType("string"), mock.Anything).Return(startTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.StartTransactionRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, idTag, request.IdTag)
		assert.Equal(t, meterStart, request.MeterStart)
		assert.Equal(t, *reservationId, *request.ReservationId)
		assertDateTimeEquality(t, *timestamp, *request.Timestamp)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.StartTransaction(connectorId, idTag, meterStart, timestamp, func(request *core.StartTransactionRequest) {
		request.ReservationId = reservationId
	})
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, *expiryDate, *confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestStartTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	meterStart := 100
	reservationId := 42
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	authorizeRequest := core.NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","meterStart":%v,"reservationId":%v,"timestamp":"%v"}]`, messageId, core.StartTransactionFeatureName, connectorId, idTag, meterStart, reservationId, timestamp.FormatTimestamp())
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
