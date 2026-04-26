package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestSetChargingProfileRequestValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	requestTable := []tests.GenericTestEntry{
		{SetChargingProfileRequest{ConnectorId: 1, ChargingProfile: chargingProfile}, true},
		{SetChargingProfileRequest{ChargingProfile: chargingProfile}, true},
		{SetChargingProfileRequest{}, false},
		{SetChargingProfileRequest{ConnectorId: 1}, false},
		{SetChargingProfileRequest{ConnectorId: -1, ChargingProfile: chargingProfile}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestSetChargingProfileConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{SetChargingProfileConfirmation{Status: ChargingProfileStatusAccepted}, true},
		{SetChargingProfileConfirmation{Status: "invalidChargingProfileStatus"}, false},
		{SetChargingProfileConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestSetChargingProfileFeature() {
	feature := SetChargingProfileFeature{}
	suite.Equal(SetChargingProfileFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetChargingProfileRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetChargingProfileConfirmation{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewSetChargingProfileRequest() {
	connectorId := 1
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	req := NewSetChargingProfileRequest(connectorId, chargingProfile)
	suite.NotNil(req)
	suite.Equal(SetChargingProfileFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(chargingProfile, req.ChargingProfile)
}

func (suite *smartChargingTestSuite) TestNewSetChargingProfileConfirmation() {
	status := ChargingProfileStatusAccepted
	conf := NewSetChargingProfileConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(SetChargingProfileFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}