package organisation

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

type UpdateOptions struct {

	//	Title of the organisation.
	Title string
}
