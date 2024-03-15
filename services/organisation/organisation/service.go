package organisation

import (
	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/storage"
)

type Service interface {
	Create(string) (*Organisation, error)
	Get(uuid.UUID) (*Organisation, error)
}

type service struct{}

func (s *service) Create(name string) (*Organisation, error) {
	return &Organisation{Name: name}, nil
}

func (s *service) Get(ID uuid.UUID) (*Organisation, error) {
	return &Organisation{Base: storage.Base{ID: ID}}, nil
}
