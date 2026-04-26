package core

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestStartTransactionRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100, ReservationId: tests.NewInt(42), Timestamp: types.NewDateTime(time.Now())}, true},
		{StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, true},
		{StartTransactionRequest{ConnectorId: 1, IdTag: "12345", Timestamp: types.NewDateTime(time.Now())}, true},
		{StartTransactionRequest{ConnectorId: 0, IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{StartTransactionRequest{ConnectorId: -1, IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{StartTransactionRequest{ConnectorId: 1, IdTag: ">20..................", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{StartTransactionRequest{IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{StartTransactionRequest{ConnectorId: 1, MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestStartTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{StartTransactionConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}, TransactionId: 10}, true},
		{StartTransactionConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}}, true},
		{StartTransactionConfirmation{IdTagInfo: &types.IdTagInfo{Status: "invalidAuthorizationStatus"}, TransactionId: 10}, false},
		{StartTransactionConfirmation{TransactionId: 10}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestStartTransactionFeature() {
	feature := StartTransactionFeature{}
	suite.Equal(StartTransactionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(StartTransactionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(StartTransactionConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewStartTransactionRequest() {
	connectorId := 1
	idTag := "12345"
	meterStart := 100
	timestamp := types.NewDateTime(time.Now())
	req := NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	suite.NotNil(req)
	suite.Equal(StartTransactionFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(idTag, req.IdTag)
	suite.Equal(meterStart, req.MeterStart)
	suite.Equal(timestamp, req.Timestamp)
}

func (suite *coreTestSuite) TestNewStartTransactionConfirmation() {
	idTagInfo := &types.IdTagInfo{Status: types.AuthorizationStatusAccepted}
	transactionId := 10
	conf := NewStartTransactionConfirmation(idTagInfo, transactionId)
	suite.NotNil(conf)
	suite.Equal(StartTransactionFeatureName, conf.GetFeatureName())
	suite.Equal(idTagInfo, conf.IdTagInfo)
	suite.Equal(transactionId, conf.TransactionId)
}