package smartcharging

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestGetCompositeScheduleRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts}, true},
		{GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600}, true},
		{GetCompositeScheduleRequest{ConnectorId: 1}, true},
		{GetCompositeScheduleRequest{ConnectorId: 0}, true},
		{GetCompositeScheduleRequest{}, true},
		{GetCompositeScheduleRequest{ConnectorId: -1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{GetCompositeScheduleRequest{ConnectorId: 1, Duration: -1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestGetCompositeScheduleConfirmationValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	confirmationTable := []tests.GenericTestEntry{
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted, ConnectorId: tests.NewInt(1), ScheduleStart: types.NewDateTime(time.Now()), ChargingSchedule: chargingSchedule}, true},
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted, ConnectorId: tests.NewInt(1), ScheduleStart: types.NewDateTime(time.Now())}, true},
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted, ConnectorId: tests.NewInt(1)}, true},
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted, ConnectorId: tests.NewInt(0)}, true},
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted}, true},
		{GetCompositeScheduleConfirmation{}, false},
		{GetCompositeScheduleConfirmation{Status: "invalidGetCompositeScheduleStatus"}, false},
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted, ConnectorId: tests.NewInt(-1)}, false},
		{GetCompositeScheduleConfirmation{Status: GetCompositeScheduleStatusAccepted, ConnectorId: tests.NewInt(1), ChargingSchedule: types.NewChargingSchedule(types.ChargingRateUnitWatts)}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestGetCompositeScheduleFeature() {
	feature := GetCompositeScheduleFeature{}
	suite.Equal(GetCompositeScheduleFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetCompositeScheduleRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetCompositeScheduleConfirmation{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewGetCompositeScheduleRequest() {
	connectorId := 1
	duration := 600
	req := NewGetCompositeScheduleRequest(connectorId, duration)
	suite.NotNil(req)
	suite.Equal(GetCompositeScheduleFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(duration, req.Duration)
}

func (suite *smartChargingTestSuite) TestNewGetCompositeScheduleConfirmation() {
	status := GetCompositeScheduleStatusAccepted
	conf := NewGetCompositeScheduleConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(GetCompositeScheduleFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}