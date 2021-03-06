package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakePPM creates a single Personally Procured Move and its associated Move and Orders
func MakePPM(db *pop.Connection, assertions Assertions) models.PersonallyProcuredMove {
	shirt := internalmessages.TShirtSizeM
	defaultAdvance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	mustCreate(db, &defaultAdvance)

	// Create new Move if not provided
	move := assertions.PersonallyProcuredMove.Move
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.PersonallyProcuredMove.MoveID) {
		move = MakeMove(db, assertions)
	}

	ppm := models.PersonallyProcuredMove{
		Move:                          move,
		MoveID:                        move.ID,
		Size:                          &shirt,
		WeightEstimate:                models.Int64Pointer(8000),
		PlannedMoveDate:               models.TimePointer(DateInsidePeakRateCycle),
		PickupPostalCode:              models.StringPointer("72017"),
		HasAdditionalPostalCode:       models.BoolPointer(false),
		AdditionalPickupPostalCode:    nil,
		DestinationPostalCode:         models.StringPointer("60605"),
		HasSit:                        models.BoolPointer(false),
		DaysInStorage:                 nil,
		Status:                        models.PPMStatusDRAFT,
		HasRequestedAdvance:           true,
		Advance:                       &defaultAdvance,
		AdvanceID:                     &defaultAdvance.ID,
		EstimatedStorageReimbursement: models.StringPointer("estimate sit"),
	}

	// Overwrite values with those from assertions
	mergeModels(&ppm, assertions.PersonallyProcuredMove)

	mustCreate(db, &ppm)

	// Add the ppm we just created to the move.ppm array
	ppm.Move.PersonallyProcuredMoves = append(ppm.Move.PersonallyProcuredMoves, ppm)

	return ppm
}

// MakeDefaultPPM makes a PPM with default values
func MakeDefaultPPM(db *pop.Connection) models.PersonallyProcuredMove {
	return MakePPM(db, Assertions{})
}

// MakePPMData creates 3 PPMs (and in turn a move and set of orders for each)
func MakePPMData(db *pop.Connection) {
	for i := 0; i < 3; i++ {
		MakeDefaultPPM(db)
	}
}
