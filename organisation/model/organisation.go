package model

import "github.com/google/uuid"

type Organisation struct {
	Base

	// Title of the organisation.
	//
	// Example: "Test Organisation"
	//
	// It is a required field.
	Title string `json:"title" gorm:"not null;check:(length(title)>0)"`

	//	ID of the user who created the organisation.
	//
	//	Example: "550e8400-e29b-41d4-a716-446655440000"
	//
	//	It is a required field.
	OwnerID uuid.UUID `json:"owner_id" gorm:"not null;type:uuid"`
}
