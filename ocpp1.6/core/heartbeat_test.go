package core

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestHeartbeatRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{HeartbeatRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestHeartbeatConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{HeartbeatConfirmation{CurrentTime: types.NewDateTime(time.Now())}, true},
		{HeartbeatConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestHeartbeatFeature() {
	feature := HeartbeatFeature{}
	suite.Equal(HeartbeatFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(HeartbeatRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(HeartbeatConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewHeartbeatRequest() {
	req := NewHeartbeatRequest()
	suite.NotNil(req)
	suite.Equal(HeartbeatFeatureName, req.GetFeatureName())
}

func (suite *coreTestSuite) TestNewHeartbeatConfirmation() {
	currentTime := types.NewDateTime(time.Now())
	conf := NewHeartbeatConfirmation(currentTime)
	suite.NotNil(conf)
	suite.Equal(HeartbeatFeatureName, conf.GetFeatureName())
	suite.Equal(currentTime, conf.CurrentTime)
}