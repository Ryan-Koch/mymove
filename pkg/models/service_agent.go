package models

import (
	"time"

	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Role represents the type of agent being recorded
type Role string

const (
	// RoleORIGIN capture enum value "ORIGIN"
	RoleORIGIN Role = "ORIGIN"
	// RoleDESTINATION capture enum value "DESTINATION"
	RoleDESTINATION Role = "DESTINATION"
)

// ServiceAgent represents an assigned agent for a shipment
type ServiceAgent struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ShipmentID       uuid.UUID `json:"shipment_id" db:"shipment_id"`
	Shipment         *Shipment `belongs_to:"shipment"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Role             Role      `json:"role" db:"role"`
	PointOfContact   string    `json:"point_of_contact" db:"point_of_contact"`
	Email            *string   `json:"email" db:"email"`
	PhoneNumber      *string   `json:"phone_number" db:"phone_number"`
	FaxNumber        *string   `json:"fax_number" db:"fax_number"`
	EmailIsPreferred *bool     `json:"email_is_preferred" db:"email_is_preferred"`
	PhoneIsPreferred *bool     `json:"phone_is_preferred" db:"phone_is_preferred"`
	Notes            *string   `json:"notes" db:"notes"`
}

// String is not required by pop and may be deleted
func (s ServiceAgent) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ServiceAgents is not required by pop and may be deleted
type ServiceAgents []ServiceAgent

// String is not required by pop and may be deleted
func (s ServiceAgents) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ServiceAgent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.ShipmentID, Name: "ShipmentID"},
		&validators.StringIsPresent{Field: string(s.Role), Name: "Role"},
		&validators.StringIsPresent{Field: s.PointOfContact, Name: "PointOfContact"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ServiceAgent) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ServiceAgent) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// CreateServiceAgent creates a ServiceAgent model from payload and queried fields.
func CreateServiceAgent(tx *pop.Connection,
	shipmentID uuid.UUID,
	role Role,
	pointOfContact *string,
	email *string,
	phoneNumber *string,
	emailIsPreferred *bool,
	phoneIsPreferred *bool,
	notes *string) (ServiceAgent, *validate.Errors, error) {

	var stringPointOfContact string
	if pointOfContact != nil {
		stringPointOfContact = string(*pointOfContact)
	}
	newServiceAgent := ServiceAgent{
		ShipmentID:       shipmentID,
		Role:             role,
		PointOfContact:   stringPointOfContact,
		Email:            email,
		PhoneNumber:      phoneNumber,
		EmailIsPreferred: emailIsPreferred,
		PhoneIsPreferred: phoneIsPreferred,
		Notes:            notes,
	}
	verrs, err := tx.ValidateAndCreate(&newServiceAgent)
	if err != nil {
		zap.L().Error("DB insertion error", zap.Error(err))
		return ServiceAgent{}, verrs, err
	} else if verrs.HasAny() {
		zap.L().Error("Validation errors", zap.Error(verrs))
		return ServiceAgent{}, verrs, errors.New("Validation error on Service Agent")
	}
	return newServiceAgent, verrs, err
}
