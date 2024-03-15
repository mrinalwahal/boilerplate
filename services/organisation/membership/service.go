package membership

import (
	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/storage"
)

type Service interface {
	Create(*CreateOptions) (*Membership, error)
	Get(uuid.UUID) (*Membership, error)
	List(*ListOptions) ([]*Membership, error)
}

type service struct{}

func (s *service) Create(options *CreateOptions) (*Membership, error) {
	return &Membership{
		OrgID:  options.OrgID,
		UserID: options.UserID,
	}, nil
}

func (s *service) Get(ID uuid.UUID) (*Membership, error) {
	return &Membership{Base: storage.Base{ID: ID}}, nil
}

func (s *service) List(options *ListOptions) ([]*Membership, error) {
	return []*Membership{
		{Base: storage.Base{ID: uuid.New()}},
	}, nil
}
