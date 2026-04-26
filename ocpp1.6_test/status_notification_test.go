package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

func (suite *OcppV16TestSuite) TestStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	status := core.ChargePointStatusAvailable
	cpErrorCode := core.NoError
	info := "mockInfo"
	vendorId := "mockVendorId"
	vendorErrorCode := "mockErrorCode"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"errorCode":"%v","info":"%v","status":"%v","timestamp":"%v","vendorId":"%v","vendorErrorCode":"%v"}]`, messageId, core.StatusNotificationFeatureName, connectorId, cpErrorCode, info, status, timestamp.FormatTimestamp(), vendorId, vendorErrorCode)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	statusNotificationConfirmation := core.NewStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(statusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.StatusNotificationRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, cpErrorCode, request.ErrorCode)
		assert.Equal(t, status, request.Status)
		assert.Equal(t, info, request.Info)
		assert.Equal(t, vendorId, request.VendorId)
		assert.Equal(t, vendorErrorCode, request.VendorErrorCode)
		assertDateTimeEquality(t, *timestamp, *request.Timestamp)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.StatusNotification(connectorId, cpErrorCode, status, func(request *core.StatusNotificationRequest) {
		request.Timestamp = timestamp
		request.Info = info
		request.VendorId = vendorId
		request.VendorErrorCode = vendorErrorCode
	})
	require.Nil(t, err)
	require.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	status := core.ChargePointStatusAvailable
	cpErrorCode := core.NoError
	info := "mockInfo"
	vendorId := "mockVendorId"
	vendorErrorCode := "mockErrorCode"
	statusNotificationRequest := core.NewStatusNotificationRequest(connectorId, cpErrorCode, status)
	statusNotificationRequest.Info = info
	statusNotificationRequest.Timestamp = timestamp
	statusNotificationRequest.VendorId = vendorId
	statusNotificationRequest.VendorErrorCode = vendorErrorCode
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"errorCode":"%v","info":"%v","status":"%v","timestamp":"%v","vendorId":"%v","vendorErrorCode":"%v"}]`, messageId, core.StatusNotificationFeatureName, connectorId, cpErrorCode, info, status, timestamp.FormatTimestamp(), vendorId, vendorErrorCode)
	testUnsupportedRequestFromCentralSystem(suite, statusNotificationRequest, requestJson, messageId)
}
