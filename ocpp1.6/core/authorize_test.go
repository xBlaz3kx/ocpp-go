package core

import (
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type coreTestSuite struct {
	suite.Suite
}

func (suite *coreTestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{AuthorizeRequest{IdTag: "12345"}, true},
		{AuthorizeRequest{}, false},
		{AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *coreTestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{AuthorizeConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}}, true},
		{AuthorizeConfirmation{IdTagInfo: &types.IdTagInfo{Status: "invalidAuthorizationStatus"}}, false},
		{AuthorizeConfirmation{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *coreTestSuite) TestAuthorizeFeature() {
	feature := AuthorizeFeature{}
	suite.Equal(AuthorizeFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(AuthorizeRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(AuthorizeConfirmation{}), feature.GetResponseType())
}

func (suite *coreTestSuite) TestNewAuthorizationRequest() {
	idTag := "12345"
	req := NewAuthorizationRequest(idTag)
	suite.NotNil(req)
	suite.Equal(AuthorizeFeatureName, req.GetFeatureName())
	suite.Equal(idTag, req.IdTag)
}

func (suite *coreTestSuite) TestNewAuthorizationConfirmation() {
	idTagInfo := &types.IdTagInfo{Status: types.AuthorizationStatusAccepted}
	conf := NewAuthorizationConfirmation(idTagInfo)
	suite.NotNil(conf)
	suite.Equal(AuthorizeFeatureName, conf.GetFeatureName())
	suite.Equal(idTagInfo, conf.IdTagInfo)
}

func TestCoreSuite(t *testing.T) {
	suite.Run(t, new(coreTestSuite))
}