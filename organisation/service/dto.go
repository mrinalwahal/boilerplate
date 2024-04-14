package service

import (
	"github.com/google/uuid"
)

// CreateOptions holds the options for creating a new organisation.
type CreateOptions struct {

	//	Title of the organisation.
	Title string

	// ID of the user who is creating the organisation.
	OwnerID uuid.UUID
}

func (o *CreateOptions) validate() error {
	if o.Title == "" {
		return ErrInvalidTitle
	}
	if o.OwnerID == uuid.Nil {
		return ErrInvalidUserID
	}
	return nil
}

type ListOptions struct {

	//	Title of the organisation.
	Title string
	//	Skip for pagination.
	Skip int
	//	Limit for pagination.
	Limit int
	//	Order by field.
	OrderBy string
	//	Order by direction.
	OrderDirection string
}

func (o *ListOptions) validate() error {
	if o.Skip < 0 {
		return ErrInvalidFilters
	}
	if o.Limit < 0 || o.Limit > 100 {
		return ErrInvalidFilters
	}
	return nil
}

type UpdateOptions struct {

	//	Title of the organisation.
	Title string
}

func (o *UpdateOptions) validate() error {
	if o.Title == "" {
		return ErrInvalidTitle
	}
	return nil
}
