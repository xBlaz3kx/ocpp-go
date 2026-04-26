package reservation

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *reservationTestSuite) TestReserveNowRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: "9999"}, true},
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{ReserveNowRequest{ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345"}, true},
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now())}, false},
		{ReserveNowRequest{ConnectorId: 1, IdTag: "12345"}, false},
		{ReserveNowRequest{ConnectorId: -1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, false},
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: -1}, true},
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: ">20.................."}, false},
		{ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: ">20.................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *reservationTestSuite) TestReserveNowConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{ReserveNowConfirmation{Status: ReservationStatusAccepted}, true},
		{ReserveNowConfirmation{Status: "invalidReserveNowStatus"}, false},
		{ReserveNowConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *reservationTestSuite) TestReserveNowFeature() {
	feature := ReserveNowFeature{}
	suite.Equal(ReserveNowFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ReserveNowRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ReserveNowConfirmation{}), feature.GetResponseType())
}

func (suite *reservationTestSuite) TestNewReserveNowRequest() {
	connectorId := 1
	expiryDate := types.NewDateTime(time.Now())
	idTag := "12345"
	reservationId := 42
	req := NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
	suite.NotNil(req)
	suite.Equal(ReserveNowFeatureName, req.GetFeatureName())
	suite.Equal(connectorId, req.ConnectorId)
	suite.Equal(expiryDate, req.ExpiryDate)
	suite.Equal(idTag, req.IdTag)
	suite.Equal(reservationId, req.ReservationId)
}

func (suite *reservationTestSuite) TestNewReserveNowConfirmation() {
	status := ReservationStatusAccepted
	conf := NewReserveNowConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(ReserveNowFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}