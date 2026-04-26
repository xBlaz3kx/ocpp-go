package core

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test"}, true},
		{BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, true},
		{BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{BootNotificationRequest{ChargeBoxSerialNumber: ">25.......................", ChargePointModel: "test", ChargePointVendor: "test"}, false},
		{BootNotificationRequest{ChargePointModel: ">20..................", ChargePointVendor: "test"}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointSerialNumber: ">25.......................", ChargePointVendor: "test"}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: ">20.................."}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", FirmwareVersion: ">50................................................"}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Iccid: ">20.................."}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Imsi: ">20.................."}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterSerialNumber: ">25......................."}, false},
		{BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterType: ">25......................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: RegistrationStatusAccepted}, true},
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: RegistrationStatusPending}, true},
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: RegistrationStatusRejected}, true},
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Status: RegistrationStatusAccepted}, true},
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: -1, Status: RegistrationStatusRejected}, false},
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: "invalidRegistrationStatus"}, false},
		{BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60}, false},
		{BootNotificationConfirmation{Interval: 60, Status: RegistrationStatusAccepted}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestBootNotificationFeature() {
	feature := BootNotificationFeature{}
	suite.Equal(BootNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(BootNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(BootNotificationConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewBootNotificationRequest() {
	model := "model1"
	vendor := "ABL"
	req := NewBootNotificationRequest(model, vendor)
	suite.NotNil(req)
	suite.Equal(BootNotificationFeatureName, req.GetFeatureName())
	suite.Equal(model, req.ChargePointModel)
	suite.Equal(vendor, req.ChargePointVendor)
}

func (suite *coreTestSuite) TestNewBootNotificationConfirmation() {
	currentTime := types.NewDateTime(time.Now())
	interval := 60
	status := RegistrationStatusAccepted
	conf := NewBootNotificationConfirmation(currentTime, interval, status)
	suite.NotNil(conf)
	suite.Equal(BootNotificationFeatureName, conf.GetFeatureName())
	suite.Equal(currentTime, conf.CurrentTime)
	suite.Equal(interval, conf.Interval)
	suite.Equal(status, conf.Status)
}