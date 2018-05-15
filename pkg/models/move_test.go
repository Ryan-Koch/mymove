package models_test

import (
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &Move{}

	expErrors := map[string][]string{
		"orders_id": {"OrdersID can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestFetchMove() {
	t := suite.T()

	order1, _ := testdatagen.MakeOrder(suite.db)
	order2, _ := testdatagen.MakeOrder(suite.db)
	reqApp := app.MyApp

	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := Move{
		OrdersID:         order1.ID,
		SelectedMoveType: &selectedType,
	}
	verrs, err := suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	// All correct
	fetchedMove, err := FetchMove(suite.db, order1.ServiceMember.User, reqApp, move.ID)
	if err != nil {
		t.Error("Expected to get moveResult back.", err)
	}
	if fetchedMove.ID != move.ID {
		t.Error("Expected new move to match move.")
	}

	// Bad Move
	fetchedMove, err = FetchMove(suite.db, order1.ServiceMember.User, reqApp, uuid.Must(uuid.NewV4()))
	if err != ErrFetchNotFound {
		t.Error("Expected to get fetchnotfound.", err)
	}

	// Bad User
	fetchedMove, err = FetchMove(suite.db, order2.ServiceMember.User, reqApp, move.ID)
	if err != ErrFetchForbidden {
		t.Error("Expected to get a Forbidden back.", err)
	}

}
