package todo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, string) (*Todo, error)
	Get(context.Context, uuid.UUID) (*Todo, error)
	List(context.Context, *ListOptions) ([]*Todo, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Todo, error)
	Delete(context.Context, uuid.UUID) error
}

type service struct {

	//	Database connection.
	db *gorm.DB
}

func (s *service) Create(ctx context.Context, title string) (*Todo, error) {
	var payload Todo
	payload.Title = title

	result := s.db.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*Todo, error) {
	var payload Todo
	payload.ID = ID
	result := s.db.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*Todo, error) {
	var payload []*Todo

	query := s.db
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.Skip > 0 {
		query = query.Offset(options.Skip)
	}
	if options.OrderBy != "" {
		query = query.Order(options.OrderBy + " " + options.OrderDirection)
	}

	//	Add conditions to the query.
	where := Todo{
		Title: options.Title,
	}

	if result := query.Where(&where).Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*Todo, error) {
	var payload Todo
	payload.ID = id
	if result := s.db.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return s.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	var payload Todo
	payload.ID = ID
	result := s.db.Delete(&payload)
	return result.Error
}
