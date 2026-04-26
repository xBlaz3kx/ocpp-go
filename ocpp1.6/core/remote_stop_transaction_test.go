package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestRemoteStopTransactionRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{RemoteStopTransactionRequest{TransactionId: 1}, true},
		{RemoteStopTransactionRequest{}, true},
		{RemoteStopTransactionRequest{TransactionId: -1}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestRemoteStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}, true},
		{RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}, true},
		{RemoteStopTransactionConfirmation{Status: "invalidRemoteStopTransactionStatus"}, false},
		{RemoteStopTransactionConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestRemoteStopTransactionFeature() {
	feature := RemoteStopTransactionFeature{}
	suite.Equal(RemoteStopTransactionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(RemoteStopTransactionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(RemoteStopTransactionConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewRemoteStopTransactionRequest() {
	transactionId := 1
	req := NewRemoteStopTransactionRequest(transactionId)
	suite.NotNil(req)
	suite.Equal(RemoteStopTransactionFeatureName, req.GetFeatureName())
	suite.Equal(transactionId, req.TransactionId)
}

func (suite *coreTestSuite) TestNewRemoteStopTransactionConfirmation() {
	status := types.RemoteStartStopStatusAccepted
	conf := NewRemoteStopTransactionConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(RemoteStopTransactionFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}