package model

import "github.com/google/uuid"

type Membership struct {
	Base

	//	ID of the member who is part of the organisation.
	//
	//	Example: "550e8400-e29b-41d4-a716-446655440000"
	//
	//	It is a required field.
	UserID uuid.UUID `json:"user_id" gorm:"not null;type:uuid"`

	// ID of the organisation to which the member belongs.
	//
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	//
	// It is a required field.
	OrgID uuid.UUID `json:"org_id" gorm:"not null;type:uuid"`
}
