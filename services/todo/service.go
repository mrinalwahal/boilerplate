package todo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, *CreateOptions) (*Todo, error)
	Get(context.Context, uuid.UUID) (*Todo, error)
	List(context.Context, *ListOptions) ([]*Todo, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Todo, error)
	Delete(context.Context, uuid.UUID) error
}

type service struct {

	//	Database connection.
	db *gorm.DB
}

func (s *service) Create(ctx context.Context, options *CreateOptions) (*Todo, error) {
	return create(s.db.WithContext(ctx), options)
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*Todo, error) {
	return get(s.db.WithContext(ctx), ID)
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*Todo, error) {
	return list(s.db.WithContext(ctx), options)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*Todo, error) {
	if err := update(s.db.WithContext(ctx), id, options); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	return delete(s.db.WithContext(ctx), ID)
}
