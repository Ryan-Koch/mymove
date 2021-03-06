package scenario

import (
	"log"
	"time"

	"github.com/go-openapi/swag"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/storage"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

// E2eBasicScenario builds a basic set of data for e2e testing
type e2eBasicScenario NamedScenario

// E2eBasicScenario Is the thing
var E2eBasicScenario = e2eBasicScenario{"e2e_basic"}

// Run does that data load thing
func (e e2eBasicScenario) Run(db *pop.Connection, loader *uploader.Uploader, logger *zap.Logger, storer *storage.Filesystem) {

	/*
	 * Basic user with tsp access
	 */
	email := "tspuser1@example.com"
	tspUser := testdatagen.MakeTspUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("6cd03e5b-bee8-4e97-a340-fecb8f3d5465")),
			LoginGovEmail: email,
		},
		TspUser: models.TspUser{
			ID:    uuid.FromStringOrNil("1fb58b82-ab60-4f55-a654-0267200473a4"),
			Email: email,
		},
	})

	/*
	 * Basic user with office access
	 */
	email = "officeuser1@example.com"
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
			LoginGovEmail: email,
		},
		OfficeUser: models.OfficeUser{
			ID:    uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
			Email: email,
		},
	})

	/*
	 * Service member with uploaded orders and a new ppm
	 */
	email = "ppm@incomple.te"
	uuidStr := "e10d5964-c070-49cb-9bd1-eaf9f7348eb6"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	nowTime := time.Now()
	ppm0 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("94ced723-fabc-42af-b9ee-87f8986bb5c9"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0db80bd6-de75-439e-bf89-deaafa1d0dc8"),
			Locator: "VGHEIS",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &nowTime,
		},
		Uploader: loader,
	})
	ppm0.Move.Submit()
	models.SaveMoveDependencies(db, &ppm0.Move)

	/*
	 * A move, that will be canceled by the E2E test
	 */
	email = "ppm-to-cancel@example.com"
	uuidStr = "e10d5964-c070-49cb-9bd1-eaf9f7348eb7"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	nowTime = time.Now()
	ppmToCancel := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("94ced723-fabc-42af-b9ee-87f8986bb5ca"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0db80bd6-de75-439e-bf89-deaafa1d0dc9"),
			Locator: "CANCEL",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &nowTime,
		},
		Uploader: loader,
	})
	ppmToCancel.Move.Submit()
	models.SaveMoveDependencies(db, &ppmToCancel.Move)

	/*
	 * Service member with a ppm in progress
	 */
	email = "ppm.in@progre.ss"
	uuidStr = "20199d12-5165-4980-9ca7-19b5dc9f1032"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	pastTime := time.Now().AddDate(0, 0, -10)
	ppm1 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("466c41b9-50bf-462c-b3cd-1ae33a2dad9b"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("In Progress"),
			Edipi:         models.StringPointer("1617033988"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("c9df71f2-334f-4f0e-b2e7-050ddb22efa1"),
			Locator: "GBXYUI",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &pastTime,
		},
		Uploader: loader,
	})
	ppm1.Move.Submit()
	ppm1.Move.Approve()
	models.SaveMoveDependencies(db, &ppm1.Move)

	/*
	 * Service member with a ppm move approved, but not in progress
	 */
	email = "ppm@approv.ed"
	uuidStr = "1842091b-b9a0-4d4a-ba22-1e2f38f26317"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	futureTime := time.Now().AddDate(0, 0, 10)
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm2 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("9ce5a930-2446-48ec-a9c0-17bc65e8522d"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Approved"),
			Edipi:         models.StringPointer("7617033988"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("12345"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0a2580ef-180a-44b2-a40b-291fa9cc13cc"),
			Locator: "FDXTIU",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &futureTime,
		},
		Uploader: loader,
	})
	ppm2.Move.Submit()
	ppm2.Move.Approve()
	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppm2.Move.PersonallyProcuredMoves[0].Submit()
	ppm2.Move.PersonallyProcuredMoves[0].Approve()
	models.SaveMoveDependencies(db, &ppm2.Move)

	/*
	 * A PPM move that has been canceled.
	 */
	email = "ppm-canceled@example.com"
	uuidStr = "20102768-4d45-449c-a585-81bc386204b1"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	nowTime = time.Now()
	ppmCanceled := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2da0d5e6-4efb-4ea1-9443-bf9ef64ace65"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Canceled"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("6b88c856-5f41-427e-a480-a7fb6c87533b"),
			Locator: "PPMCAN",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &nowTime,
		},
		Uploader: loader,
	})
	ppmCanceled.Move.Submit()
	models.SaveMoveDependencies(db, &ppmCanceled.Move)
	ppmCanceled.Move.Cancel("reasons")
	models.SaveMoveDependencies(db, &ppmCanceled.Move)

	/*
	 * Service member with orders and a move
	 */
	email = "profile@comple.te"
	uuidStr = "13F3949D-0D53-4BE4-B1B1-AE4314793F34"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})

	testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0a1e72b0-1b9f-442b-a6d3-7b7cfa6bbb95"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Profile"),
			LastName:      models.StringPointer("Complete"),
			Edipi:         models.StringPointer("8893308161"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("173da49c-fcec-4d01-a622-3651e81c654e"),
			Locator: "BLABLA",
		},
		Uploader: loader,
	})

	/*
	 * Service member with orders and a move, but no move type selected to select HHG
	 */
	email = "sm_hhg@example.com"
	uuidStr = "4b389406-9258-4695-a091-0bf97b5a132f"

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})

	dutyStationAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "Fort Gordon",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30813",
			Country:        swag.String("United States"),
		},
	})

	dutyStation := testdatagen.MakeDutyStation(db, testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Sam Houston",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   dutyStationAddress.ID,
			Address:     dutyStationAddress,
		},
	})

	testdatagen.MakeMoveWithoutMoveType(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("b5d1f44b-5ceb-4a0e-9119-5687808996ff"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("HHGDude"),
			LastName:      models.StringPointer("UserPerson"),
			Edipi:         models.StringPointer("6833908163"),
			PersonalEmail: models.StringPointer(email),
			DutyStationID: &dutyStation.ID,
			DutyStation:   dutyStation,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("8718c8ac-e0c6-423b-bdc6-af971ee05b9a"),
			Locator: "REWGIE",
		},
	})

	/*
	 * Another service member with orders and a move, but no move type selected
	 */
	email = "sm_no_move_type@example.com"
	uuidStr = "9ceb8321-6a82-4f6d-8bb3-a1d85922a202"

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})

	testdatagen.MakeMoveWithoutMoveType(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("7554e347-2215-484f-9240-c61bae050220"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("HHGDude2"),
			LastName:      models.StringPointer("UserPerson2"),
			Edipi:         models.StringPointer("6833908164"),
			PersonalEmail: models.StringPointer(email),
			DutyStationID: &dutyStation.ID,
			DutyStation:   dutyStation,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("b2ecbbe5-36ad-49fc-86c8-66e55e0697a7"),
			Locator: "ZPGVED",
		},
	})

	/*
	 * Service member with uploaded orders and a new shipment move
	 */
	email = "hhg@incomple.te"

	hhg0 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("ebc176e0-bb34-47d4-ba37-ff13e2dd40b9")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0d719b18-81d6-474a-86aa-b87246fff65c"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("2ed0b5a2-26d9-49a3-a775-5220055e8ffe"),
			Locator:          "RLKBEM",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("0dfdbdda-c57e-4b29-994a-09fb8641fc75"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
	})

	hhg0.Move.Submit()
	models.SaveMoveDependencies(db, &hhg0.Move)

	/*
	 * Service member with uploaded orders and an approved shipment
	 */
	email = "hhg@award.ed"

	offer1 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("7980f0cf-63e3-4722-b5aa-ba46f8f7ac64")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("8a66beef-1cdf-4117-9db2-aad548f54430"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("56b8ef45-8145-487b-9b59-0e30d0d465fa"),
			Locator:          "BACON1",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("776b5a23-2830-4de0-bb6a-7698a25865cb"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAWARDED,
			HasDeliveryAddress: true,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg1 := offer1.Shipment
	hhg1.Move.Submit()
	models.SaveMoveDependencies(db, &hhg1.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to be accepted
	 */
	email = "hhg@fromawardedtoaccept.ed"

	packDate := time.Now().AddDate(0, 0, 1)
	pickupDate := time.Now().AddDate(0, 0, 5)
	deliveryDate := time.Now().AddDate(0, 0, 10)
	sourceOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ABCD",
		},
	})
	destOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "QRED",
		},
	})
	offer2 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("179598c5-a5ee-4da5-8259-29749f03a398")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("179598c5-a5ee-4da5-8259-29749f03a398"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForAccept"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			DepartmentIndicator: models.StringPointer("17"),
			TAC:                 models.StringPointer("NTA4"),
			SAC:                 models.StringPointer("1234567890 9876543210"),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("849a7880-4a82-4f76-acb4-63cf481e786b"),
			Locator:          "BACON2",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("5f86c201-1abf-4f9d-8dcb-d039cb1c6bfc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			ID:                          uuid.FromStringOrNil("53ebebef-be58-41ce-9635-a4930149190d"),
			Status:                      models.ShipmentStatusAWARDED,
			PmSurveyPlannedPackDate:     &packDate,
			PmSurveyConductedDate:       &packDate,
			PmSurveyPlannedPickupDate:   &pickupDate,
			PmSurveyPlannedDeliveryDate: &deliveryDate,
			SourceGBLOC:                 &sourceOffice.Gbloc,
			DestinationGBLOC:            &destOffice.Gbloc,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			TransportationServiceProvider:   tspUser.TransportationServiceProvider,
		},
	})

	_, err := testdatagen.MakeTSPPerformanceDeprecated(db,
		tspUser.TransportationServiceProvider,
		*offer2.Shipment.TrafficDistributionList,
		models.IntPointer(3),
		0.40,
		5,
		unit.DiscountRate(0.50),
		unit.DiscountRate(0.55))
	if err != nil {
		log.Panic(err)
	}

	hhg2 := offer2.Shipment
	hhg2.Move.Submit()
	models.SaveMoveDependencies(db, &hhg2.Move)

	/*
	 * Service member with accepted shipment
	 */
	email = "hhg@accept.ed"

	offer3 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("6a39dd2a-a23f-4967-a035-3bc9987c6848")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("6a39dd2a-a23f-4967-a035-3bc9987c6848"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("4752270d-4a6f-44ea-82f6-ae3cf3277c5d"),
			Locator:          "BACON3",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("e09f8b8b-67a6-4ce3-b5c3-bd48c82512fc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusACCEPTED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg3 := offer3.Shipment
	hhg3.Move.Submit()
	models.SaveMoveDependencies(db, &hhg3.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to have weight added
	 */
	email = "hhg@addweigh.ts"

	offer4 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("bf022aeb-3f14-4429-94d7-fe759f493aed")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("01fa956f-d17b-477e-8607-1db1dd891720"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("94739ee0-664c-47c5-afe9-0f5067a2e151"),
			Locator:          "BACON4",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("9ebc891b-f629-4ea1-9ebf-eef1971d69a3"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAWARDED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg4 := offer4.Shipment
	hhg4.Move.Submit()
	models.SaveMoveDependencies(db, &hhg4.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to have weight added
	 * This shipment is rejected by the e2e test.
	 */
	email = "hhg@reject.ing"

	offer5 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("76bdcff3-ade4-41ff-bf09-0b2474cec751")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("f4e362e9-9fdd-490b-a2fa-1fa4035b8f0d"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("7fca3fd0-08a6-480a-8a9c-16a65a100db9"),
			Locator:          "REJECT",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("1731c3e6-b510-43d0-be46-13e5a2032bad"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAWARDED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg5 := offer5.Shipment
	hhg5.Move.Submit()
	models.SaveMoveDependencies(db, &hhg5.Move)

	/*
	 * Service member with in transit shipment
	 */
	email = "hhg@in.transit"

	offer6 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("1239dd2a-a23f-4967-a035-3bc9987c6848")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2339dd2a-a23f-4967-a035-3bc9987c6824"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("3452270d-4a6f-44ea-82f6-ae3cf3277c5d"),
			Locator:          "NINOPK",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("459f8b8b-67a6-4ce3-b5c3-bd48c82512fc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusINTRANSIT,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg6 := offer6.Shipment
	hhg6.Move.Submit()
	models.SaveMoveDependencies(db, &hhg6.Move)

	/*
	 * Service member with approved shipment
	 */
	email = "hhg@approv.ed"

	offer7 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("68461d67-5385-4780-9cb6-417075343b0e")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2825cadf-410f-4f82-aa0f-4caaf000e63e"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("616560f2-7e35-4504-b7e6-69038fb0c015"),
			Locator:          "APPRVD",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("5fe59be4-45d0-47c7-b426-cf4db9882af7"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAPPROVED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg7 := offer7.Shipment
	hhg7.Move.Submit()
	models.SaveMoveDependencies(db, &hhg7.Move)

	/*
	 * Service member with approved basics and accepted shipment
	 */
	email = "hhg@accept.ed"

	offer8 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("f79fd68e-4461-4ba8-b630-9618b913e229")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("f79fd68e-4461-4ba8-b630-9618b913e229"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("29cd6b2f-9ef2-48be-b4ee-1c1e0a1456ef"),
			Locator:          "BACON5",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		Order: models.Order{
			OrdersNumber:        models.StringPointer("54321"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("d17e2e3e-9bff-4bb0-b301-f97ad03350c1"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusACCEPTED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg8 := offer8.Shipment
	hhg8.Move.Submit()
	models.SaveMoveDependencies(db, &hhg8.Move)

	/*
	 * Service member with uploaded orders, a new shipment move, and a service agent
	 */
	email = "hhg@incomplete.serviceagent"
	hhg9 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("412e76e0-bb34-47d4-ba37-ff13e2dd40b9")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("245a9b18-81d6-474a-86aa-b87246fff65c"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("1a3eb5a2-26d9-49a3-a775-5220055e8ffe"),
			Locator:          "LRKREK",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("873dbdda-c57e-4b29-994a-09fb8641fc75"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
	})
	testdatagen.MakeServiceAgent(db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			ShipmentID: hhg9.ID,
		},
	})
	testdatagen.MakeServiceAgent(db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			ShipmentID: hhg9.ID,
			Role:       models.RoleDESTINATION,
		},
	})
	hhg9.Move.Submit()
	models.SaveMoveDependencies(db, &hhg9.Move)

	/*
	 * Service member with delivered shipment
	 */
	email = "hhg@de.livered"

	offer10 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("3339dd2a-a23f-4967-a035-3bc9987c6848")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2559dd2a-a23f-4967-a035-3bc9987c6824"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("3442270d-4a6f-44ea-82f6-ae3cf3277c5d"),
			Locator:          "SCHNOO",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("466f8b8b-67a6-4ce3-b5c3-bd48c82512fc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusDELIVERED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg10 := offer10.Shipment
	hhg10.Move.Submit()
	models.SaveMoveDependencies(db, &hhg10.Move)

	/*
	 * Service member with completed shipment
	 */
	email = "hhg@com.pleted"

	offer11 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("4449dd2a-a23f-4967-a035-3bc9987c6848")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("5559dd2a-a23f-4967-a035-3bc9987c6824"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("9992270d-4a6f-44ea-82f6-ae3cf3277c5d"),
			Locator:          "NOCHKA",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("777f8b8b-67a6-4ce3-b5c3-bd48c82512fc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusCOMPLETED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg11 := offer11.Shipment
	hhg11.Move.Submit()
	models.SaveMoveDependencies(db, &hhg11.Move)

	/*
	 * Service member with approved basics and accepted shipment to be approved
	 */
	email = "hhg@delivered.tocomplete"

	offer12 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("ab9fd68e-4461-4ba8-b630-9618b913e229")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("ab9fd68e-4461-4ba8-b630-9618b913e229"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("abcd6b2f-9ef2-48be-b4ee-1c1e0a1456ef"),
			Locator:          "SSETZN",
			SelectedMoveType: models.StringPointer("HHG"),
			Status:           models.MoveStatusAPPROVED,
		},
		Order: models.Order{
			OrdersNumber:        models.StringPointer("54321"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("ab7e2e3e-9bff-4bb0-b301-f97ad03350c1"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusDELIVERED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg12 := offer12.Shipment
	hhg12.Move.Submit()
	models.SaveMoveDependencies(db, &hhg12.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to be accepted & able to generate GBL
	 */
	MakeHhgFromAwardedToAcceptedGBLReady(db, tspUser)

	/*
	 * Service member with uploaded orders and an approved shipment to be accepted & GBL generated
	 */
	MakeHhgWithGBL(db, tspUser, logger, storer)

	/*
	 * Service member with uploaded orders and an approved shipment
	 */
	email = "hhg@premo.ve"

	offer13 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("8f6b87f1-20ad-4c50-a855-ab66e222c7c3")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("1a98be36-5c4c-4056-b16f-d5a6c65b8569"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("01d85649-18c2-44ad-854d-da8884579f42"),
			Locator:          "PREMVE",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("fd76c4fc-a2fb-45b6-a3a6-7c35357ab79a"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAWARDED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg13 := offer13.Shipment
	hhg13.Move.Submit()
	models.SaveMoveDependencies(db, &hhg13.Move)

	/*
	 * Service member with uploaded orders and an approved shipment
	 */
	email = "hhg@dates.panel"

	offer14 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("444b87f1-20ad-4c50-a855-ab66e222c7c3")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("1222be36-5c4c-4056-b16f-d5a6c65b8569"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("11185649-18c2-44ad-854d-da8884579f42"),
			Locator:          "DATESP",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("feeec4fc-a2fb-45b6-a3a6-7c35357ab79a"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAWARDED,
			ActualDeliveryDate: nil,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg14 := offer14.Shipment
	hhg14.Move.Submit()
	models.SaveMoveDependencies(db, &hhg14.Move)

	/* Service member with an in progress for doc testing on TSP side
	 */
	email = "doc.viewer@tsp.org"

	offer15 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("027f183d-a45e-44fc-b890-cd5092a99ecb")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("b5df04d6-2a35-4294-9a67-8a2427eba0bc"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("9999888777"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("ccd45bd4-660c-4ddd-b6c6-062da0a647f9"),
			Locator:          "DOCVWR",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("7ad595da-9b34-4914-aeaa-9a540d13872f"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAWARDED,
			ActualDeliveryDate: nil,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg15 := offer15.Shipment
	hhg15.Move.Submit()
	models.SaveMoveDependencies(db, &hhg15.Move)

	/* Service member with an in progress for testing delivery address
	 */
	email = "duty.station@tsp.org"

	offer16 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("2b6036ce-acf1-40fc-86da-d2b32329054f")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("815a314e-3c30-430b-bae7-54cf9ded79d4"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("9999888777"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("761743d9-2259-4bee-b144-3bda29311446"),
			Locator:          "DTYSTN",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("a36a58b4-51ab-4d39-bcdc-b3ca3a59a4a1"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAWARDED,
			HasDeliveryAddress: false,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg16 := offer16.Shipment
	hhg16.Move.Submit()
	models.SaveMoveDependencies(db, &hhg16.Move)

	/*
	 * Service member with a in progress for doc upload testing on TSP side
	 */
	email = "doc.upload@tsp.org"

	offer17 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("0d033463-09fd-498b-869d-30cda1c95599")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("98b21b35-9709-4da8-9d42-a2e887cf1e6c"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("9999888777"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("787c0921-7696-400e-86f5-c1bcb8bb88a3"),
			Locator:          "DOCUPL",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("ad3ee670-6978-46e1-bcfc-686cdd4ffa87"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			ID:                 uuid.FromStringOrNil("65e00326-420e-436a-89fc-6aeb3f90b870"),
			Status:             models.ShipmentStatusAWARDED,
			ActualDeliveryDate: nil,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg17 := offer17.Shipment
	hhg17.Move.Submit()
	models.SaveMoveDependencies(db, &hhg17.Move)

	/*
	 * Service member with approved basics and awarded shipment (can't approve shipment yet)
	 */
	email = "hhg@cant.approve"

	offer18 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("0187d7e5-2ee7-410a-b42f-d889a78b0bff")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2f3ad6c4-2e6c-4c45-a7e5-8b220ebaabb6"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("BasicsApproveOnly"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("57e9275b-b433-474c-99f2-ac64966b3c9b"),
			Locator:          "BACON6",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		Order: models.Order{
			OrdersNumber:        models.StringPointer("54321"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("cb49e75e-7897-4a01-8cff-c13ae85ca5ba"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAWARDED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg18 := offer18.Shipment
	hhg18.Move.Submit()
	models.SaveMoveDependencies(db, &hhg18.Move)

	/*
	 * Service member with uploaded orders and an approved shipment
	 */
	email = "hhg@doc.uploads"

	offer19 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("5245b1ff-ae5a-4875-8a21-6b05c735b684")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("60cdcd83-6d6f-442f-a5b5-c256b312d000"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("533d176f-0bab-4c51-88cd-c899f6855b9d"),
			Locator:          "BACON7",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("60e65f0c-aa21-4d95-a825-9d323a3dc4f1"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAWARDED,
			HasDeliveryAddress: true,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg19 := offer19.Shipment
	hhg19.Move.Submit()
	models.SaveMoveDependencies(db, &hhg19.Move)

	/*
	 * Service member with uploaded orders and an approved shipment. Use this to test zeroing dates.
	 */
	email = "hhg@dates.panel"

	offer20 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("cf1b1f09-8ea2-4f68-872e-a056c3a5f22f")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("6c4bc296-927c-4c6b-a01e-1f064c5d5f9b"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("da9af941-253a-45e0-b012-8ee0385e28f8"),
			Locator:          "DATESZ",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("9728e6a1-0469-4718-9ba1-5d7baace1191"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAWARDED,
			ActualDeliveryDate: nil,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg20 := offer20.Shipment
	hhg20.Move.Submit()
	models.SaveMoveDependencies(db, &hhg20.Move)

	/* Service member with a doc for testing on TSP side
	 */
	email = "doc.owner@tsp.org"

	offer21 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("99bdaeed-a8a8-492e-9e28-7d0da6b1c907")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("61473913-36b8-425d-b46a-cee488a4ae71"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("2232332334"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("60098ff1-8dc9-4318-a2e8-47bc8aac11a4"),
			Locator:          "GOTDOC",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("7ad595da-9b34-4914-aeaa-9a540d13872f"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status:             models.ShipmentStatusAPPROVED,
			ActualDeliveryDate: nil,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
		Document: models.Document{
			ID:              uuid.FromStringOrNil("06886210-b151-4b15-951a-783d3d58f042"),
			ServiceMemberID: uuid.FromStringOrNil("61473913-36b8-425d-b46a-cee488a4ae71"),
		},
		MoveDocument: models.MoveDocument{
			ID:               uuid.FromStringOrNil("b660080d-0158-4214-99ca-216f82b26b3c"),
			DocumentID:       uuid.FromStringOrNil("06886210-b151-4b15-951a-783d3d58f042"),
			MoveDocumentType: "WEIGHT_TICKET",
			Status:           "OK",
			MoveID:           uuid.FromStringOrNil("60098ff1-8dc9-4318-a2e8-47bc8aac11a4"),
			Title:            "document_title",
		},
	})

	hhg21 := offer21.Shipment
	hhg21.Move.Submit()
	models.SaveMoveDependencies(db, &hhg21.Move)

	/*
	 * Service member with uploaded orders and an approved shipment with service agent
	 */
	email = "hhg@enter.premove"

	offer22 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("426b87f1-20ad-4c50-a855-ab66e222c7c3")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("4298be36-5c4c-4056-b16f-d5a6c65b8569"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Approved"),
			Edipi:         models.StringPointer("4424567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("42d85649-18c2-44ad-854d-da8884579f42"),
			Locator:          "ENTPMS",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("f426c4fc-a2fb-45b6-a3a6-7c35357ab79a"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAPPROVED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	testdatagen.MakeServiceAgent(db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			Shipment:   &offer22.Shipment,
			ShipmentID: offer22.ShipmentID,
		},
	})

	hhg22 := offer22.Shipment
	hhg22.Move.Submit()
	models.SaveMoveDependencies(db, &hhg22.Move)

	/*
	 * Service member with accepted move but needs to be assigned a service agent
	 */
	email = "hhg@assign.serviceagent"
	offer23 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("8ff1c3ca-4c51-40ad-9926-8add5463eb25")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("e52c90df-502f-4fa2-8838-ee0894725b4d"),
			FirstName:     models.StringPointer("Assign"),
			LastName:      models.StringPointer("ServiceAgent"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("33686dbe-cd64-4786-8aaa-a93dda278683"),
			Locator:          "ASSIGN",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("d40edb7e-24c9-4a21-8e4b-2e473471263e"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusACCEPTED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})
	hhg23 := offer23.Shipment
	hhg23.Move.Submit()
	models.SaveMoveDependencies(db, &hhg23.Move)

	/*
	 * Service member with in transit shipment
	 */
	email = "enter@delivery.date"

	offer24 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("1af7ca19-8511-4c6e-a93b-144811c0fa7c")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("ae29e24b-b048-4c17-88d6-a008b91d0f85"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("135af727-f570-4c7e-bf5b-878d717ef83c"),
			Locator:          "ENTDEL",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("ebfac7dc-acfa-4a88-bbbf-a2dd1a7f2657"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusINTRANSIT,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg24 := offer24.Shipment
	hhg24.Move.Submit()
	models.SaveMoveDependencies(db, &hhg24.Move)

	/*
	 * Service member with a cancelled HHG move.
	 */
	email = "hhg@cancel.ed"

	hhg25 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("05ea5bc3-fd77-4f42-bdc5-a984a81b3829")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("d27bcb66-fc51-42b6-a13b-c896d34c79fb"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Cancelled"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("da6bf1f4-a810-486d-befe-ddf8e9a4e2ef"),
			Locator:          "HHGCAN",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("d89dba9c-5ee9-40ee-8430-2e3eb13eedeb"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
	})

	hhg25.Move.Submit()
	models.SaveMoveDependencies(db, &hhg25.Move)
	hhg25.Move.Cancel("reasons")
	models.SaveMoveDependencies(db, &hhg25.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to have weight added in the office app
	 */
	email = "hhg@addweights.office"

	offer26 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("611aea22-1689-4e16-90e7-e55d49010069")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("033297aa-4f4d-4df1-a05d-d22d717f6d5b"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("2be4f6a3-82f5-4919-a257-39a24859058f"),
			Locator:          "WTSPNL",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("d2c24faf-3439-451f-b020-fc1492f6b4bf"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAWARDED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})

	hhg26 := offer26.Shipment
	hhg26.Move.Submit()
	models.SaveMoveDependencies(db, &hhg26.Move)

	/*
	* Service member to update dates from office app
	 */
	email = "hhg1@officeda.te"

	hhg27 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("5e2f7338-0f54-4ba9-99cc-796153da94f3")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("8cfe7777-d8a5-43ed-bb0e-5ba2ceda2251"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444555888"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("47e9c534-a93c-4986-ae8f-41fddefaa618"),
			Locator:          "ODATES",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("51004395-ecbf-4ab2-9edc-ec5041bbe390"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
	})

	hhg27.Move.Submit()
	models.SaveMoveDependencies(db, &hhg27.Move)

	/*
	 * Service member to update dates from office app
	 */
	email = "hhg2@officeda.te"

	hhg28 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("961108be-ace1-407c-b110-7e996e95d286")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("823a2177-3d68-43a5-a3ed-6b10454a6481"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444999888"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("762f2ec2-f362-4c14-b601-d7178c4862fe"),
			Locator:          "ODATE0",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("12e17ea7-9c94-4b61-a28d-5a81744a355c"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
	})
	hhg28.Move.Submit()
	models.SaveMoveDependencies(db, &hhg28.Move)
}

// MakeHhgFromAwardedToAcceptedGBLReady creates a scenario for an approved shipment ready for GBL generation
func MakeHhgFromAwardedToAcceptedGBLReady(db *pop.Connection, tspUser models.TspUser) models.Shipment {
	/*
	 * Service member with uploaded orders and an approved shipment to be accepted, able to generate GBL
	 */
	email := "hhg@govbilloflading.ready"

	packDate := time.Now().AddDate(0, 0, 1)
	pickupDate := time.Now().AddDate(0, 0, 5)
	deliveryDate := time.Now().AddDate(0, 0, 10)
	weightEstimate := unit.Pound(5000)
	sourceOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ABCD",
		},
	})
	destOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "QRED",
		},
	})
	offer9 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("658f3a78-b3a9-47f4-a820-af673103d62d")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("658f3a78-b3a9-47f4-a820-af673103d62d"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForGBL"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("NTA4"),
			SAC:                 models.StringPointer("1234567890 9876543210"),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("05a58b2e-07da-4b41-b4f8-d18ab68dddd5"),
			Locator:          "GBLGBL",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("b15fdc2b-52cd-4b3e-91ba-a36d6ab94a16"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			ID:                          uuid.FromStringOrNil("a4013cee-aa0a-41a3-b5f5-b9eed0758e1d 0xc42022c070"),
			Status:                      models.ShipmentStatusAPPROVED,
			PmSurveyConductedDate:       &packDate,
			PmSurveyMethod:              "PHONE",
			PmSurveyPlannedPackDate:     &packDate,
			PmSurveyPlannedPickupDate:   &pickupDate,
			PmSurveyPlannedDeliveryDate: &deliveryDate,
			PmSurveyWeightEstimate:      &weightEstimate,
			SourceGBLOC:                 &sourceOffice.Gbloc,
			DestinationGBLOC:            &destOffice.Gbloc,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			TransportationServiceProvider:   tspUser.TransportationServiceProvider,
			Accepted:                        models.BoolPointer(true),
		},
	})

	testdatagen.MakeTSPPerformanceDeprecated(db,
		tspUser.TransportationServiceProvider,
		*offer9.Shipment.TrafficDistributionList,
		models.IntPointer(3),
		0.40,
		5,
		unit.DiscountRate(0.50),
		unit.DiscountRate(0.55))

	testdatagen.MakeServiceAgent(db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			Shipment:   &offer9.Shipment,
			ShipmentID: offer9.ShipmentID,
		},
	})

	hhg2 := offer9.Shipment
	hhg2.Move.Submit()
	models.SaveMoveDependencies(db, &hhg2.Move)
	return offer9.Shipment
}

// MakeHhgWithGBL creates a scenario for an approved shipment with a GBL generated
func MakeHhgWithGBL(db *pop.Connection, tspUser models.TspUser, logger *zap.Logger, storer *storage.Filesystem) models.Shipment {
	/*
	 * Service member with uploaded orders and an approved shipment to be accepted, able to generate GBL
	 */
	email := "hhg@gov_bill_of_lading.created"

	packDate := time.Now().AddDate(0, 0, 1)
	pickupDate := time.Now().AddDate(0, 0, 5)
	deliveryDate := time.Now().AddDate(0, 0, 10)
	weightEstimate := unit.Pound(5000)
	sourceOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ABCD",
		},
	})
	destOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "QRED",
		},
	})
	offer := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("b7dccea1-d052-4a66-aed9-2fdacf461023")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("8a1a86c7-78d6-4897-806e-0e4c5546fdec"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("HasGBL"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			DepartmentIndicator: models.StringPointer("17"),
			TAC:                 models.StringPointer("NTA4"),
			SAC:                 models.StringPointer("1234567890 9876543210"),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("6eee3663-1973-40c5-b49e-e70e9325b895"),
			Locator:          "CONGBL",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("87fcebf6-63b8-40cb-bc40-b553f5b91b9c"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			ID:                          uuid.FromStringOrNil("0851706a-997f-46fb-84e4-2525a444ade0"),
			Status:                      models.ShipmentStatusAPPROVED,
			PmSurveyConductedDate:       &packDate,
			PmSurveyMethod:              "PHONE",
			PmSurveyPlannedPackDate:     &packDate,
			PmSurveyPlannedPickupDate:   &pickupDate,
			PmSurveyPlannedDeliveryDate: &deliveryDate,
			PmSurveyWeightEstimate:      &weightEstimate,
			SourceGBLOC:                 &sourceOffice.Gbloc,
			DestinationGBLOC:            &destOffice.Gbloc,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			TransportationServiceProvider:   tspUser.TransportationServiceProvider,
			Accepted:                        models.BoolPointer(true),
		},
	})

	testdatagen.MakeTSPPerformanceDeprecated(db,
		tspUser.TransportationServiceProvider,
		*offer.Shipment.TrafficDistributionList,
		models.IntPointer(3),
		0.40,
		5,
		unit.DiscountRate(0.50),
		unit.DiscountRate(0.55))

	testdatagen.MakeServiceAgent(db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			Shipment:   &offer.Shipment,
			ShipmentID: offer.ShipmentID,
		},
	})

	hhg := offer.Shipment
	hhgID := offer.ShipmentID
	hhg.Move.Submit()
	models.SaveMoveDependencies(db, &hhg.Move)

	// Create PDF for GBL
	gbl, _ := models.FetchGovBillOfLadingExtractor(db, hhgID)
	formLayout := paperwork.Form1203Layout

	// Read in bytes from Asset pkg
	data, _ := assets.Asset(formLayout.TemplateImagePath)
	f, _ := storer.FileSystem().Create("something.png")
	f.Write(data)
	f.Seek(0, 0)

	form, _ := paperwork.NewTemplateForm(f, formLayout.FieldsLayout)

	// Populate form fields with GBL data
	form.DrawData(gbl)
	aFile, _ := storer.FileSystem().Create(gbl.GBLNumber1)
	form.Output(aFile)

	uploader := uploaderpkg.NewUploader(db, logger, storer)
	upload, _, _ := uploader.CreateUpload(nil, *tspUser.UserID, aFile)
	uploads := []models.Upload{*upload}

	// Create GBL move document associated to the shipment
	hhg.Move.CreateMoveDocument(db,
		uploads,
		&hhgID,
		models.MoveDocumentTypeGOVBILLOFLADING,
		string("Government Bill Of Lading"),
		swag.String(""),
		string(apimessages.SelectedMoveTypeHHG),
	)

	return offer.Shipment
}
