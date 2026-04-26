package smartcharging

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type smartChargingTestSuite struct {
	suite.Suite
}

func (suite *smartChargingTestSuite) TestClearChargingProfileRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{ClearChargingProfileRequest{Id: tests.NewInt(1), ConnectorId: tests.NewInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, StackLevel: tests.NewInt(1)}, true},
		{ClearChargingProfileRequest{Id: tests.NewInt(1), ConnectorId: tests.NewInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile}, true},
		{ClearChargingProfileRequest{Id: tests.NewInt(1), ConnectorId: tests.NewInt(1)}, true},
		{ClearChargingProfileRequest{Id: tests.NewInt(1)}, true},
		{ClearChargingProfileRequest{}, true},
		{ClearChargingProfileRequest{ConnectorId: tests.NewInt(-1)}, false},
		{ClearChargingProfileRequest{Id: tests.NewInt(-1)}, true},
		{ClearChargingProfileRequest{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}, false},
		{ClearChargingProfileRequest{StackLevel: tests.NewInt(-1)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestClearChargingProfileConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{ClearChargingProfileConfirmation{Status: ClearChargingProfileStatusAccepted}, true},
		{ClearChargingProfileConfirmation{Status: "invalidClearChargingProfileStatus"}, false},
		{ClearChargingProfileConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestClearChargingProfileFeature() {
	feature := ClearChargingProfileFeature{}
	suite.Equal(ClearChargingProfileFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearChargingProfileRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearChargingProfileConfirmation{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewClearChargingProfileRequest() {
	req := NewClearChargingProfileRequest()
	suite.NotNil(req)
	suite.Equal(ClearChargingProfileFeatureName, req.GetFeatureName())
}

func (suite *smartChargingTestSuite) TestNewClearChargingProfileConfirmation() {
	status := ClearChargingProfileStatusAccepted
	conf := NewClearChargingProfileConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(ClearChargingProfileFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}

func TestSmartChargingSuite(t *testing.T) {
	suite.Run(t, new(smartChargingTestSuite))
}