package organisation

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, string) (*Organisation, error)
	Get(context.Context, uuid.UUID) (*Organisation, error)
	List(context.Context, *ListOptions) ([]*Organisation, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Organisation, error)
	Delete(context.Context, uuid.UUID) error
}

type service struct {

	//	Database connection.
	db *gorm.DB
}

func (s *service) Create(ctx context.Context, title string) (*Organisation, error) {
	txn := s.db.WithContext(ctx)

	var payload Organisation
	payload.Title = title

	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*Organisation, error) {
	txn := s.db.WithContext(ctx)

	var payload Organisation
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*Organisation, error) {
	txn := s.db.WithContext(ctx)

	var payload []*Organisation

	query := txn
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
	where := Organisation{
		Title: options.Title,
	}

	if result := query.Where(&where).Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*Organisation, error) {
	txn := s.db.WithContext(ctx)

	var payload Organisation
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return s.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	txn := s.db.WithContext(ctx)

	var payload Organisation
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}
