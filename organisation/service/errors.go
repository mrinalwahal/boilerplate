package service

import "fmt"

var (
	ErrInvalidOptions        = fmt.Errorf("invalid options")
	ErrInvalidorganisationID = fmt.Errorf("invalid organisation_id")
	ErrInvalidUserID         = fmt.Errorf("invalid user_id")
	ErrInvalidTitle          = fmt.Errorf("invalid title")
	ErrInvalidFilters        = fmt.Errorf("invalid filters")
	ErrInvalidDB             = fmt.Errorf("invalid db")
)
