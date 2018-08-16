package internal

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateDocumentsHandler() {
	t := suite.parent.T()

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	params := documentop.NewCreateDocumentParams()
	params.DocumentPayload = &internalmessages.PostDocumentPayload{
		ServiceMemberID: *utils.FmtUUID(serviceMember.ID),
	}

	req := &http.Request{}
	req = suite.parent.AuthenticateRequest(req, serviceMember)
	params.HTTPRequest = req

	handler := CreateDocumentHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	createdResponse, ok := response.(*documentop.CreateDocumentCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
	documentPayload := createdResponse.Payload

	if uuid.Must(uuid.FromString(documentPayload.ID.String())) == uuid.Nil {
		t.Errorf("got empty document uuid")
	}

	if uuid.Must(uuid.FromString(documentPayload.ServiceMemberID.String())) == uuid.Nil {
		t.Errorf("got empty serviceMember uuid")
	}

	if len(documentPayload.Uploads) != 0 {
		t.Errorf("wrong number of uploads, expected 0, got %d", len(documentPayload.Uploads))
	}

	document := models.Document{}
	err := suite.parent.Db.Find(&document, documentPayload.ID)
	if err != nil {
		t.Errorf("Couldn't find expected document.")
	}
}

func (suite *HandlerSuite) TestShowDocumentHandler() {
	t := suite.parent.T()

	upload := testdatagen.MakeDefaultUpload(suite.parent.Db)

	documentID := upload.DocumentID
	var document models.Document

	err := suite.parent.Db.Eager("ServiceMember.User").Find(&document, documentID)
	if err != nil {
		t.Fatalf("could not load document: %s", err)
	}

	params := documentop.NewShowDocumentParams()
	params.DocumentID = strfmt.UUID(documentID.String())

	req := &http.Request{}
	req = suite.parent.AuthenticateRequest(req, document.ServiceMember)
	params.HTTPRequest = req

	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)
	handler := ShowDocumentHandler(context)
	response := handler.Handle(params)

	showResponse, ok := response.(*documentop.ShowDocumentOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}
	documentPayload := showResponse.Payload

	responseDocumentUUID := documentPayload.ID.String()
	if responseDocumentUUID != documentID.String() {
		t.Errorf("wrong document uuid, expected %v, got %v", documentID, responseDocumentUUID)
	}

	if len(documentPayload.Uploads) != 1 {
		t.Errorf("wrong number of uploads, expected 1, got %d", len(documentPayload.Uploads))
	}

	uploadPayload := documentPayload.Uploads[0]
	expectedURL := fmt.Sprintf("https://example.com/dir/%s?contentType=application/pdf&signed=test", upload.StorageKey)
	if (*uploadPayload.URL).String() != expectedURL {
		t.Errorf("wrong URL for upload, expected %s, got %s", expectedURL, uploadPayload.URL)
	}
}
