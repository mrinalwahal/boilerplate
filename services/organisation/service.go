package organisation

import (
	membership "github.com/mrinalwahal/boilerplate/services/organisation/membership"
	"github.com/mrinalwahal/boilerplate/services/organisation/organisation"
)

type Service interface {
	Organisations() organisation.Service
	Memberships() membership.Service
}

type service struct{}

func (s *service) Organisations() organisation.Service {
	return organisation.GetService()
}

func (s *service) Memberships() membership.Service {
	return membership.GetService()
}
