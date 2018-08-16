package internal

import (
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/pjdufour-truss/pdfcpu/pkg/api"
	"github.com/pjdufour-truss/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) assertPDFPageCount(count int, file afero.File, storer storage.FileStorer) {
	pdfConfig := pdfcpu.NewInMemoryConfiguration()
	pdfConfig.FileSystem = storer.FileSystem()

	ctx, err := api.Read(file.Name(), pdfConfig)
	suite.parent.NoError(err)

	err = pdfcpu.ValidateXRefTable(ctx.XRefTable)
	suite.parent.NoError(err)

	suite.parent.Equal(2, ctx.PageCount)
}

func (suite *HandlerSuite) createHandlerContext() utils.HandlerContext {
	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	fakeS3 := storageTest.NewFakeS3Storage(true)
	context.SetFileStorer(fakeS3)

	return context
}

func (suite *HandlerSuite) TestCreatePPMAttachmentsHandler() {
	uploadKeyRe := regexp.MustCompile(`(user/.+/uploads/.+)\?`)

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)
	ppm := testdatagen.MakeDefaultPPM(suite.parent.Db)
	expDoc := testdatagen.MakeMovingExpenseDocument(suite.parent.Db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			PersonallyProcuredMoveID: &ppm.ID,
		},
	})

	// Context gives us our file storer and filesystem
	context := suite.createHandlerContext()

	// Open our test file
	f, err := os.Open("fixtures/test.pdf")
	suite.parent.NoError(err)

	// Backfill the uploaded orders file in filesystem
	uploadedOrdersUpload := ppm.Move.Orders.UploadedOrders.Uploads[0]
	_, err = context.Storage.Store(uploadedOrdersUpload.StorageKey, f, uploadedOrdersUpload.Checksum)
	suite.parent.NoError(err)

	// Create upload for expense document model
	loader := uploader.NewUploader(suite.parent.Db, suite.parent.Logger, context.Storage)
	loader.CreateUpload(&expDoc.MoveDocument.DocumentID, *officeUser.UserID, f)

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.parent.AuthenticateOfficeRequest(request, officeUser)

	params := ppmop.CreatePPMAttachmentsParams{
		PersonallyProcuredMoveID: *utils.FmtUUID(ppm.ID),
		HTTPRequest:              request,
	}

	handler := CreatePersonallyProcuredMoveAttachmentsHandler(context)
	response := handler.Handle(params)
	// assert we got back the 201 response
	suite.parent.IsNotErrResponse(response)
	createdResponse := response.(*ppmop.CreatePPMAttachmentsOK)
	createdPDFPayload := createdResponse.Payload
	suite.parent.NotNil(createdPDFPayload.URL)

	// Extract upload key from returned URL
	attachmentsURL := string(*createdPDFPayload.URL)
	uploadKey := uploadKeyRe.FindStringSubmatch(attachmentsURL)[1]

	merged, err := context.Storage.Fetch(uploadKey)
	suite.parent.NoError(err)
	mergedFile := merged.(afero.File)

	suite.assertPDFPageCount(2, mergedFile, context.Storage)
}
