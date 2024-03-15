package membership

import (
	"github.com/mrinalwahal/boilerplate/storage"
)

type Membership struct {
	storage.Base

	//	OrgID is the ID of the organisation.
	OrgID string `json:"org_id"`

	//	UserID is the ID of the user.
	UserID string `json:"user_id"`
}
