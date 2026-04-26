package core

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError, Info: "mockInfo", Status: ChargePointStatusAvailable, Timestamp: types.NewDateTime(time.Now()), VendorId: "mockId", VendorErrorCode: "mockErrorCode"}, true},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError, Status: ChargePointStatusAvailable}, true},
		{StatusNotificationRequest{ErrorCode: NoError, Status: ChargePointStatusAvailable}, true},
		{StatusNotificationRequest{ConnectorId: -1, ErrorCode: NoError, Status: ChargePointStatusAvailable}, false},
		{StatusNotificationRequest{ConnectorId: 0, Status: ChargePointStatusAvailable}, false},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError}, false},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: "invalidErrorCode", Status: ChargePointStatusAvailable}, false},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError, Status: "invalidChargePointStatus"}, false},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError, Info: ">50................................................", Status: ChargePointStatusAvailable}, false},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError, VendorErrorCode: ">50................................................", Status: ChargePointStatusAvailable}, false},
		{StatusNotificationRequest{ConnectorId: 0, ErrorCode: NoError, VendorId: ">255............................................................................................................................................................................................................................................................", Status: ChargePointStatusAvailable}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{StatusNotificationConfirmation{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestStatusNotificationFeature() {
	feature := StatusNotificationFeature{}
	suite.Equal(StatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(StatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(StatusNotificationConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewStatusNotificationRequest() {
	connectorId := 1
	errorCode := NoError
	status := ChargePointStatusAvailable
	req := NewStatusNotificationRequest(connectorId, errorCode, status)
	suite.NotNil(req)
	suite.Equal(StatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(errorCode, req.ErrorCode)
	suite.Equal(status, req.Status)
}

func (suite *coreTestSuite) TestNewStatusNotificationConfirmation() {
	conf := NewStatusNotificationConfirmation()
	suite.NotNil(conf)
	suite.Equal(StatusNotificationFeatureName, conf.GetFeatureName())
}