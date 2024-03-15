package membership

type CreateOptions struct {

	//	OrgID is the ID of the organisation.
	OrgID string `json:"org_id"`

	//	UserID is the ID of the user.
	UserID string `json:"user_id"`
}

type ListOptions struct {

	//	OrgID is the ID of the organisation.
	OrgID string `json:"org_id"`

	//	UserID is the ID of the user.
	UserID string `json:"user_id"`
}
