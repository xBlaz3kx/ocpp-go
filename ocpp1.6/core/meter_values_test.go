package core

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestMeterValuesRequestValidation() {
	var requestTable = []tests.GenericTestEntry{
		{MeterValuesRequest{ConnectorId: 1, TransactionId: tests.NewInt(1), MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, true},
		{MeterValuesRequest{ConnectorId: 1, MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, true},
		{MeterValuesRequest{MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, true},
		{MeterValuesRequest{ConnectorId: -1, MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, false},
		{MeterValuesRequest{ConnectorId: 1, MeterValue: []types.MeterValue{}}, false},
		{MeterValuesRequest{ConnectorId: 1}, false},
		{MeterValuesRequest{ConnectorId: 1, MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{}}}}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *coreTestSuite) TestMeterValuesConfirmationValidation() {
	var confirmationTable = []tests.GenericTestEntry{
		{MeterValuesConfirmation{}, true},
	}
	tests.ExecuteGenericTestTable(suite.T(), confirmationTable)
}

func (suite *coreTestSuite) TestMeterValuesFeature() {
	feature := MeterValuesFeature{}
	suite.Equal(MeterValuesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(MeterValuesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(MeterValuesConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewMeterValuesRequest() {
	connectorId := 1
	meterValues := []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}
	req := NewMeterValuesRequest(connectorId, meterValues)
	suite.NotNil(req)
	suite.Equal(MeterValuesFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(meterValues, req.MeterValue)
}

func (suite *coreTestSuite) TestNewMeterValuesConfirmation() {
	conf := NewMeterValuesConfirmation()
	suite.NotNil(conf)
	suite.Equal(MeterValuesFeatureName, conf.GetFeatureName())
}