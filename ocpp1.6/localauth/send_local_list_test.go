package localauth

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *localAuthTestSuite) TestSendLocalListRequestValidation() {
	t := suite.T()
	localAuthEntry := AuthorizationData{IdTag: "12345", IdTagInfo: &types.IdTagInfo{
		ExpiryDate:  types.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "000000",
		Status:      types.AuthorizationStatusAccepted,
	}}
	invalidAuthEntry := AuthorizationData{IdTag: "12345", IdTagInfo: &types.IdTagInfo{
		ExpiryDate:  types.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "000000",
		Status:      "invalidStatus",
	}}
	requestTable := []tests.GenericTestEntry{
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []AuthorizationData{localAuthEntry}}, true},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []AuthorizationData{}}, true},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential, ListVersion: 1}, true},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential, ListVersion: 0}, true},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential}, true},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential, ListVersion: -1}, false},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []AuthorizationData{invalidAuthEntry}}, false},
		{SendLocalListRequest{UpdateType: "invalidUpdateType", ListVersion: 1}, false},
		{SendLocalListRequest{ListVersion: 1}, false},
		{SendLocalListRequest{}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *localAuthTestSuite) TestSendLocalListConfirmationValidation() {
	t := suite.T()
	confirmationTable := []tests.GenericTestEntry{
		{SendLocalListConfirmation{Status: UpdateStatusAccepted}, true},
		{SendLocalListConfirmation{Status: "invalidStatus"}, false},
		{SendLocalListConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *localAuthTestSuite) TestSendLocalListFeature() {
	feature := SendLocalListFeature{}
	suite.Equal(SendLocalListFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SendLocalListRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SendLocalListConfirmation{}), feature.GetResponseType())
}

func (suite *localAuthTestSuite) TestNewSendLocalListRequest() {
	version := 1
	updateType := UpdateTypeDifferential
	req := NewSendLocalListRequest(version, updateType)
	suite.NotNil(req)
	suite.Equal(SendLocalListFeatureName, req.GetFeatureName())
	suite.Equal(version, req.ListVersion)
	suite.Equal(updateType, req.UpdateType)
}

func (suite *localAuthTestSuite) TestNewSendLocalListConfirmation() {
	status := UpdateStatusAccepted
	conf := NewSendLocalListConfirmation(status)
	suite.NotNil(conf)
	suite.Equal(SendLocalListFeatureName, conf.GetFeatureName())
	suite.Equal(status, conf.Status)
}