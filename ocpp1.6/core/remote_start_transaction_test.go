package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestRemoteStartTransactionRequestValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	var requestTable = []tests.GenericTestEntry{
		{RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: tests.NewInt(1), ChargingProfile: chargingProfile}, true},
		{RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: tests.NewInt(1)}, true},
		{RemoteStartTransactionRequest{IdTag: "12345"}, true},
		{RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: tests.NewInt(-1)}, false},
		{RemoteStartTransactionRequest{}, false},
		{RemoteStartTransactionRequest{IdTag: ">20..................", ConnectorId: tests.NewInt(1)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestRemoteStartTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}, true},
		{RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}, true},
		{RemoteStartTransactionConfirmation{Status: "invalidRemoteStartTransactionStatus"}, false},
		{RemoteStartTransactionConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestRemoteStartTransactionFeature() {
	feature := RemoteStartTransactionFeature{}
	suite.Equal(RemoteStartTransactionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(RemoteStartTransactionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(RemoteStartTransactionConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewRemoteStartTransactionRequest() {
	idTag := "12345"
	req := NewRemoteStartTransactionRequest(idTag)
	suite.NotNil(req)
	suite.Equal(RemoteStartTransactionFeatureName, req.GetFeatureName())
	suite.Equal(idTag, req.IdTag)
}

func (suite *coreTestSuite) TestNewRemoteStartTransactionConfirmation() {
	status := types.RemoteStartStopStatusAccepted
	conf := NewRemoteStartTransactionConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(RemoteStartTransactionFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}