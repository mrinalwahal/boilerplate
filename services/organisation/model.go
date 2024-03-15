package organisation

import (
	"github.com/mrinalwahal/boilerplate/storage"
)

type Organisation struct {
	storage.Base

	//	Name of the organisation.
	Name string `json:"name"`
}
