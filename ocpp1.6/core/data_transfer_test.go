package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *coreTestSuite) TestDataTransferRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{DataTransferRequest{VendorId: "12345"}, true},
		{DataTransferRequest{VendorId: "12345", MessageId: "6789"}, true},
		{DataTransferRequest{VendorId: "12345", MessageId: "6789", Data: "mockData"}, true},
		{DataTransferRequest{}, false},
		{DataTransferRequest{VendorId: ">255............................................................................................................................................................................................................................................................"}, false},
		{DataTransferRequest{VendorId: "12345", MessageId: ">50................................................"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestDataTransferConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{DataTransferConfirmation{Status: DataTransferStatusAccepted}, true},
		{DataTransferConfirmation{Status: DataTransferStatusRejected}, true},
		{DataTransferConfirmation{Status: DataTransferStatusUnknownMessageId}, true},
		{DataTransferConfirmation{Status: DataTransferStatusUnknownVendorId}, true},
		{DataTransferConfirmation{Status: "invalidDataTransferStatus"}, false},
		{DataTransferConfirmation{Status: DataTransferStatusAccepted, Data: "mockData"}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestDataTransferFeature() {
	feature := DataTransferFeature{}
	suite.Equal(DataTransferFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(DataTransferRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(DataTransferConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewDataTransferRequest() {
	vendorId := "12345"
	req := NewDataTransferRequest(vendorId)
	suite.NotNil(req)
	suite.Equal(DataTransferFeatureName, req.GetFeatureName())
	suite.Equal(vendorId, req.VendorId)
}

func (suite *coreTestSuite) TestNewDataTransferConfirmation() {
	status := DataTransferStatusAccepted
	conf := NewDataTransferConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(DataTransferFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}