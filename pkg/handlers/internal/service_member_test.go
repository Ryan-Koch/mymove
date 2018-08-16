package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowServiceMemberHandler() {
	// Given: A servicemember and a user
	user := testdatagen.MakeDefaultUser(suite.parent.Db)

	newServiceMember := models.ServiceMember{
		UserID: user.ID,
	}
	suite.parent.MustSave(&newServiceMember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: show ServiceMember is queried
	showHandler := ShowServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.parent.Assertions.IsType(&servicememberop.ShowServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOK)

	// And: Returned query to include our added servicemember
	suite.parent.Assertions.Equal(user.ID.String(), okResponse.Payload.UserID.String())
}

func (suite *HandlerSuite) TestShowServiceMemberWrongUser() {
	// Given: Servicemember trying to load another
	notLoggedInUser := testdatagen.MakeDefaultServiceMember(suite.parent.Db)
	loggedInUser := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	req := httptest.NewRequest("GET", fmt.Sprintf("/service_members/%s", notLoggedInUser.ID.String()), nil)
	req = suite.parent.AuthenticateRequest(req, loggedInUser)

	showServiceMemberParams := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(notLoggedInUser.ID.String()),
	}
	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := showHandler.Handle(showServiceMemberParams)

	suite.parent.Assertions.IsType(&utils.ErrResponse{}, response)
	errResponse := response.(*utils.ErrResponse)

	suite.parent.Assertions.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerAllValues() {
	// Given: A logged-in user
	user := testdatagen.MakeDefaultUser(suite.parent.Db)

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		UserID:                 strfmt.UUID(user.ID.String()),
		Edipi:                  swag.String("random string bla"),
		FirstName:              swag.String("random string bla"),
		MiddleName:             swag.String("random string bla"),
		LastName:               swag.String("random string bla"),
		Suffix:                 swag.String("random string bla"),
		Telephone:              swag.String("random string bla"),
		SecondaryTelephone:     swag.String("random string bla"),
		PersonalEmail:          swag.String("wml@example.com"),
		PhoneIsPreferred:       swag.Bool(false),
		TextMessageIsPreferred: swag.Bool(false),
		EmailIsPreferred:       swag.Bool(true),
		ResidentialAddress:     fakeAddressPayload(),
		BackupMailingAddress:   fakeAddressPayload(),
		SocialSecurityNumber:   (*strfmt.SSN)(swag.String("123-45-6789")),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&utils.CookieUpdateResponder{}, response)
	unwrappedResponse := response.(*utils.CookieUpdateResponder).Responder
	suite.parent.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)

	// Then: we expect a servicemember to have been created for the user
	query := suite.parent.Db.Where(fmt.Sprintf("user_id='%v'", user.ID))
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.parent.Assertions.Len(serviceMembers, 1)
}

func (suite *HandlerSuite) TestSubmitServiceMemberSSN() {
	// Given: A logged-in user
	user := testdatagen.MakeDefaultUser(suite.parent.Db)
	session := &auth.Session{
		UserID:          user.ID,
		ApplicationName: auth.MyApp,
	}

	// When: a new ServiceMember is posted
	ssn := "123-45-6789"
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		SocialSecurityNumber: (*strfmt.SSN)(swag.String(ssn)),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&utils.CookieUpdateResponder{}, response)
	unwrappedResponse := response.(*utils.CookieUpdateResponder).Responder
	suite.parent.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)
	okResponse := unwrappedResponse.(*servicememberop.CreateServiceMemberCreated)

	suite.parent.Assertions.True(*okResponse.Payload.HasSocialSecurityNumber)

	// Then: we expect a ServiceMember to have been created for the user
	query := suite.parent.Db.Where(fmt.Sprintf("user_id='%v'", user.ID))
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.parent.Assertions.Len(serviceMembers, 1)

	serviceMemberID, _ := uuid.FromString(okResponse.Payload.ID.String())

	session.ServiceMemberID = serviceMemberID
	serviceMember, err := models.FetchServiceMemberForUser(suite.parent.Db, session, serviceMemberID)
	suite.parent.Assertions.NoError(err)

	suite.parent.Assertions.True(serviceMember.SocialSecurityNumber.Matches(ssn))
}

func (suite *HandlerSuite) TestPatchServiceMemberHandler() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.parent.Db)

	// TODO: add more fields to change
	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.parent.MustSave(&newServiceMember)

	affiliation := internalmessages.AffiliationARMY
	rank := internalmessages.ServiceMemberRankE1
	ssn := utils.FmtSSN("555-55-5555")
	resAddress := fakeAddressPayload()
	backupAddress := fakeAddressPayload()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &newEdipi,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &affiliation,
		EmailIsPreferred:     swag.Bool(true),
		FirstName:            swag.String("Firstname"),
		LastName:             swag.String("Lastname"),
		MiddleName:           swag.String("Middlename"),
		PersonalEmail:        swag.String("name@domain.com"),
		PhoneIsPreferred:     swag.Bool(true),
		Rank:                 &rank,
		TextMessageIsPreferred: swag.Bool(true),
		SecondaryTelephone:     swag.String("555555555"),
		SocialSecurityNumber:   ssn,
		Suffix:                 swag.String("Sr."),
		Telephone:              swag.String("555555555"),
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.PatchServiceMemberOK)

	serviceMemberPayload := okResponse.Payload

	suite.parent.Assertions.Equal(*serviceMemberPayload.Edipi, newEdipi)
	suite.parent.Assertions.Equal(*serviceMemberPayload.Affiliation, affiliation)
	suite.parent.Assertions.Equal(*serviceMemberPayload.HasSocialSecurityNumber, true)
	suite.parent.Assertions.Equal(*serviceMemberPayload.TextMessageIsPreferred, true)
	suite.parent.Assertions.Equal(*serviceMemberPayload.ResidentialAddress.StreetAddress1, *resAddress.StreetAddress1)
	suite.parent.Assertions.Equal(*serviceMemberPayload.BackupMailingAddress.StreetAddress1, *backupAddress.StreetAddress1)

	// Then: we expect addresses to have been created
	addresses := []models.Address{}
	suite.parent.Db.All(&addresses)
	suite.parent.Assertions.Len(addresses, 2)

	// Then: we expect a SSN to have been created
	ssns := models.SocialSecurityNumbers{}
	suite.parent.Db.All(&ssns)
	suite.parent.Assertions.Len(ssns, 1)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerWrongUser() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.parent.Db)
	user2 := testdatagen.MakeDefaultUser(suite.parent.Db)

	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.parent.MustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateUserRequest(req, user2)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&utils.ErrResponse{}, response)
	errResponse := response.(*utils.ErrResponse)

	suite.parent.Assertions.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoServiceMember() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.parent.Db)

	servicememberUUID := uuid.Must(uuid.NewV4())

	var newEdipi = "9999999999"

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateUserRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(servicememberUUID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&utils.ErrResponse{}, response)
	errResponse := response.(*utils.ErrResponse)

	suite.parent.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoChange() {
	// Given: a logged in user with a servicemember
	user := testdatagen.MakeDefaultUser(suite.parent.Db)

	var origEdipi = "4444444444"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.parent.MustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.parent.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
}

func (suite *HandlerSuite) TestShowServiceMemberOrders() {
	order1 := testdatagen.MakeDefaultOrder(suite.parent.Db)
	order2Assertions := testdatagen.Assertions{
		Order: models.Order{
			ServiceMember:   order1.ServiceMember,
			ServiceMemberID: order1.ServiceMemberID,
		},
	}
	order2 := testdatagen.MakeOrder(suite.parent.Db, order2Assertions)

	req := httptest.NewRequest("GET", "/service_members/some_id/current_orders", nil)
	req = suite.parent.AuthenticateRequest(req, order1.ServiceMember)

	params := servicememberop.ShowServiceMemberOrdersParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(order1.ServiceMemberID.String()),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetFileStorer(fakeS3)
	handler := ShowServiceMemberOrdersHandler(context)

	response := handler.Handle(params)

	suite.parent.IsType(&servicememberop.ShowServiceMemberOrdersOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOrdersOK)
	responsePayload := okResponse.Payload

	// Should return the most recently created order
	suite.parent.Equal(order2.ID.String(), responsePayload.ID.String())
}
