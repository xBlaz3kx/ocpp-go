package reservation

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type reservationTestSuite struct {
	suite.Suite
}

func (suite *reservationTestSuite) TestCancelReservationRequestValidation() {
	t := suite.T()
	requestTable := []tests.GenericTestEntry{
		{CancelReservationRequest{ReservationId: 42}, true},
		{CancelReservationRequest{}, true},
		{CancelReservationRequest{ReservationId: -1}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *reservationTestSuite) TestCancelReservationConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{CancelReservationConfirmation{Status: CancelReservationStatusAccepted}, true},
		{CancelReservationConfirmation{Status: "invalidCancelReservationStatus"}, false},
		{CancelReservationConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *reservationTestSuite) TestCancelReservationFeature() {
	feature := CancelReservationFeature{}
	suite.Equal(CancelReservationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(CancelReservationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(CancelReservationConfirmation{}), feature.GetResponseType())
}

func (suite *reservationTestSuite) TestNewCancelReservationRequest() {
	reservationId := 42
	req := NewCancelReservationRequest(reservationId)
	suite.NotNil(req)
	suite.Equal(CancelReservationFeatureName, req.GetFeatureName())
	suite.Equal(reservationId, req.ReservationId)
}

func (suite *reservationTestSuite) TestNewCancelReservationConfirmation() {
	status := CancelReservationStatusAccepted
	conf := NewCancelReservationConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(CancelReservationFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}

func TestReservationSuite(t *testing.T) {
	suite.Run(t, new(reservationTestSuite))
}