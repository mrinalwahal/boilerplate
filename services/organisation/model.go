package organisation

import (
	"github.com/mrinalwahal/boilerplate/storage"
)

type Organisation struct {
	storage.Base

	//	Title of the organisation.
	Title string `json:"title" gorm:"not null"`
}
