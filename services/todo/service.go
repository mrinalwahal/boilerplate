package todo

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(context.Context, *CreateOptions) (*Todo, error)
	Get(context.Context, uuid.UUID) (*Todo, error)
	List(context.Context, *ListOptions) ([]*Todo, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Todo, error)
	Delete(context.Context, uuid.UUID) error
}

// Initializes and gets the service with the supplied database connection.
func NewService(db DB) Service {
	return &service{db: db}
}

type service struct {

	//	Database layer connection.
	db DB
}

func (s *service) Create(ctx context.Context, options *CreateOptions) (*Todo, error) {
	return s.db.Create(ctx, options)
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*Todo, error) {
	return s.db.Get(ctx, ID)
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*Todo, error) {
	return s.db.List(ctx, options)
}

func (s *service) Update(ctx context.Context, ID uuid.UUID, options *UpdateOptions) (*Todo, error) {
	return s.db.Update(ctx, ID, options)
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	return s.db.Delete(ctx, ID)
}
