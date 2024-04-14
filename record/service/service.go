//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/pkg/middleware"
	"github.com/mrinalwahal/boilerplate/record/model"
	"gorm.io/gorm"
)

// Service interface declares the signature of the service layer.
type Service interface {
	Create(context.Context, *CreateOptions) (*model.Record, error)
	List(context.Context, *ListOptions) ([]*model.Record, error)
	Get(context.Context, uuid.UUID) (*model.Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*model.Record, error)
	Delete(context.Context, uuid.UUID) error
}

type Config struct {

	// Service connection.
	// The connection should already be open.
	//
	// This field is mandatory.
	DB *gorm.DB

	// Logger instance.
	//
	// This field is mandatory.
	Logger *slog.Logger
}

func NewService(config *Config) Service {
	if config == nil {
		panic("service: nil config")
	}

	service := service{
		db:  config.DB,
		log: config.Logger,
	}

	if service.log == nil {
		service.log = slog.Default()
	}

	service.log = service.log.With("service", "record")

	return &service
}

// service is the standard and default service layer implementation.
//
// It implements the Service interface.
type service struct {

	//	Service Connection
	db *gorm.DB

	// Logger instance
	log *slog.Logger
}

// Create operation creates a new record in the service.
func (service *service) Create(ctx context.Context, options *CreateOptions) (*model.Record, error) {
	service.log.LogAttrs(ctx, slog.LevelDebug, "creating a new record",
		slog.String("function", "create"),
	)
	if options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	txn := service.db.WithContext(ctx)

	//
	// This method has no Row Level Security (RLS) checks.
	//

	// Prepare the payload we have to send to the service transaction.
	var payload model.Record
	payload.Title = options.Title
	payload.UserID = options.UserID

	// Execute the transaction.
	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

// List operation fetches a list of records from the service.
func (service *service) List(ctx context.Context, options *ListOptions) ([]*model.Record, error) {
	service.log.LogAttrs(ctx, slog.LevelDebug, "fetching records",
		slog.String("function", "list"),
	)
	if options == nil {
		options = &ListOptions{}
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	txn := service.db.WithContext(ctx)

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can list it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload []*model.Record

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
	if options.Title != "" {
		query = query.Where(&model.Record{
			Title: options.Title,
		})
	}

	if result := query.Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

// Get operation fetches a record from the service.
func (service *service) Get(ctx context.Context, ID uuid.UUID) (*model.Record, error) {
	service.log.LogAttrs(ctx, slog.LevelDebug, "fetching a record",
		slog.String("function", "get"),
	)
	if ID == uuid.Nil {
		return nil, ErrInvalidRecordID
	}

	txn := service.db.WithContext(ctx)

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can get it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload model.Record
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

// Update operation updates a record in the service.
func (service *service) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*model.Record, error) {
	service.log.LogAttrs(ctx, slog.LevelDebug, "updating a record",
		slog.String("function", "update"),
	)
	if id == uuid.Nil {
		return nil, ErrInvalidRecordID
	}
	if options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	txn := service.db.WithContext(ctx)

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can update it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload model.Record
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return service.Get(ctx, id)
}

// Delete operation deletes a record from the service.
func (service *service) Delete(ctx context.Context, ID uuid.UUID) error {
	service.log.LogAttrs(ctx, slog.LevelDebug, "deleting a record",
		slog.String("function", "delete"),
	)
	if ID == uuid.Nil {
		return ErrInvalidRecordID
	}

	txn := service.db.WithContext(ctx)

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can delete it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload model.Record
	payload.ID = ID
	result := txn.Delete(&payload)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
