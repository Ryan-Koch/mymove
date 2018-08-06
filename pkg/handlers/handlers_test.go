package handlers

import (
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/uploader"
)

type HandlerSuite struct {
	suite.Suite
	db                 *pop.Connection
	logger             *zap.Logger
	filesToClose       []*runtime.File
	notificationSender notifications.NotificationSender
}

func (suite *HandlerSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *HandlerSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

func (suite *HandlerSuite) isNotErrResponse(response middleware.Responder) {
	r, ok := response.(*errResponse)
	if ok {
		suite.logger.Error("Received an unexpected error response from handler: ", zap.Error(r.err))
		// Formally lodge a complaint
		suite.IsType(&errResponse{}, response)
	}
}

func (suite *HandlerSuite) checkErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*errResponse)
	if !ok || errResponse.code != code {
		suite.T().Fatalf("Expected %s Response: %v", name, resp)
		debug.PrintStack()
	}
}

func (suite *HandlerSuite) checkResponseBadRequest(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

func (suite *HandlerSuite) checkResponseUnauthorized(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

func (suite *HandlerSuite) checkResponseForbidden(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

func (suite *HandlerSuite) checkResponseNotFound(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusNotFound, "NotFound")
}

func (suite *HandlerSuite) checkResponseInternalServerError(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

func (suite *HandlerSuite) checkResponseTeapot(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusTeapot, "Teapot")
}

// Request authenticated with a service member
func (suite *HandlerSuite) authenticateRequest(req *http.Request, serviceMember models.ServiceMember) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          serviceMember.UserID,
		IDToken:         "fake token",
		ServiceMemberID: serviceMember.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// Request only authenticated with a bare user - have no idea if they are a service member yet
func (suite *HandlerSuite) authenticateUserRequest(req *http.Request, user models.User) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

func (suite *HandlerSuite) authenticateOfficeRequest(req *http.Request, user models.OfficeUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

func (suite *HandlerSuite) authenticateTspRequest(req *http.Request, user models.TspUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		TspUserID:       user.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

func (suite *HandlerSuite) fixture(name string) *runtime.File {
	fixtureDir := "fixtures"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Error(err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)
	file, err := uploader.NewLocalFile(fixturePath)

	if err != nil {
		suite.T().Error(err)
	}
	suite.closeFile(file)
	return file
}

func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Data.Close()
	}
}

func (suite *HandlerSuite) closeFile(file *runtime.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

func TestHandlerSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{
		db:                 db,
		logger:             logger,
		notificationSender: notifications.NewStubNotificationSender(logger),
	}

	suite.Run(t, hs)
}
