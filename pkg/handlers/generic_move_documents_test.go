package handlers

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateGenericMoveDocumentHandler() {
	move := testdatagen.MakeDefaultMove(suite.db)
	sm := move.Orders.ServiceMember

	upload := testdatagen.MakeUpload(suite.db, testdatagen.Assertions{
		Upload: models.Upload{
			UploaderID: sm.UserID,
		},
	})
	upload.DocumentID = nil
	suite.mustSave(&upload)
	uploadIds := []strfmt.UUID{*fmtUUID(upload.ID)}

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, sm)

	newMoveDocPayload := internalmessages.CreateGenericMoveDocumentPayload{
		UploadIds:        uploadIds,
		MoveDocumentType: internalmessages.MoveDocumentTypeOTHER,
		Title:            fmtString("awesome_document.pdf"),
		Notes:            fmtString("Some notes here"),
	}

	newMoveDocParams := movedocop.CreateGenericMoveDocumentParams{
		HTTPRequest:                      request,
		CreateGenericMoveDocumentPayload: &newMoveDocPayload,
		MoveID: strfmt.UUID(move.ID.String()),
	}

	context := NewHandlerContext(suite.db, suite.logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := CreateGenericMoveDocumentHandler(context)
	response := handler.Handle(newMoveDocParams)
	// assert we got back the 201 response
	suite.isNotErrResponse(response)
	createdResponse := response.(*movedocop.CreateGenericMoveDocumentOK)
	createdPayload := createdResponse.Payload
	suite.NotNil(createdPayload.ID)

	// Make sure the Upload was associated to the new document
	createdDocumentID := createdPayload.Document.ID
	var fetchedUpload models.Upload
	suite.db.Find(&fetchedUpload, upload.ID)
	suite.Equal(createdDocumentID.String(), fetchedUpload.DocumentID.String())

	// Next try the wrong user
	wrongUser := testdatagen.MakeDefaultServiceMember(suite.db)
	request = suite.authenticateRequest(request, wrongUser)
	newMoveDocParams.HTTPRequest = request

	badUserResponse := handler.Handle(newMoveDocParams)
	suite.checkResponseForbidden(badUserResponse)

	// Now try a bad move
	newMoveDocParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newMoveDocParams)
	suite.checkResponseNotFound(badMoveResponse)
}
