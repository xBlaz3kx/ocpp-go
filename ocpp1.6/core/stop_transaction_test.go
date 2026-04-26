package core

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestStopTransactionRequestValidation() {
	t := suite.T()
	transactionData := []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}
	var requestTable = []tests.GenericTestEntry{
		{StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: ReasonEVDisconnected, TransactionData: transactionData}, true},
		{StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: ReasonEVDisconnected, TransactionData: []types.MeterValue{}}, true},
		{StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: ReasonEVDisconnected}, true},
		{StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, true},
		{StopTransactionRequest{MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, true},
		{StopTransactionRequest{MeterStop: 100, Timestamp: types.NewDateTime(time.Now())}, true},
		{StopTransactionRequest{Timestamp: types.NewDateTime(time.Now())}, true},
		{StopTransactionRequest{MeterStop: 100}, false},
		{StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: "invalidReason"}, false},
		{StopTransactionRequest{IdTag: ">20..................", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, false},
		{StopTransactionRequest{MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, TransactionData: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{StopTransactionConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}}, true},
		{StopTransactionConfirmation{}, true},
		{StopTransactionConfirmation{IdTagInfo: &types.IdTagInfo{Status: "invalidAuthorizationStatus"}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestStopTransactionFeature() {
	feature := StopTransactionFeature{}
	suite.Equal(StopTransactionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(StopTransactionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(StopTransactionConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewStopTransactionRequest() {
	meterStop := 100
	timestamp := types.NewDateTime(time.Now())
	transactionId := 42
	req := NewStopTransactionRequest(meterStop, timestamp, transactionId)
	suite.NotNil(req)
	suite.Equal(StopTransactionFeatureName, req.GetFeatureName())
	suite.Equal(meterStop, req.MeterStop)
	suite.Equal(timestamp, req.Timestamp)
	suite.Equal(transactionId, req.TransactionId)
}

func (suite *coreTestSuite) TestNewStopTransactionConfirmation() {
	conf := NewStopTransactionConfirmation()
	suite.NotNil(conf)
	suite.Equal(StopTransactionFeatureName, conf.GetFeatureName())
}