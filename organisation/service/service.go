//go:generate mockgen -destination=service_mock.go -source=service.go -package=service
package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/organisation/db/organisation"
	"github.com/mrinalwahal/boilerplate/organisation/model"
)

type Service interface {
	Create(context.Context, *CreateOptions) (*model.Organisation, error)
	List(context.Context, *ListOptions) ([]*model.Organisation, error)
	Get(context.Context, uuid.UUID) (*model.Organisation, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*model.Organisation, error)
	Delete(context.Context, uuid.UUID) error

	Members()
	Roles()
}

type Config struct {

	//	Database layer.
	DB organisation.DB

	//	Logger.
	Logger *slog.Logger
}

// Initializes and gets the service with the supplied database connection.
func NewService(config *Config) Service {

	if config == nil {
		panic("service: nil config")
	}

	svc := service{
		db:     config.DB,
		logger: config.Logger,
	}

	if svc.logger == nil {
		svc.logger = slog.Default()
	}

	svc.logger = svc.logger.With("layer", "service")

	return &svc
}

type service struct {

	//	Database layer service.
	db organisation.DB

	//	Logger.
	logger *slog.Logger
}

func (s *service) Create(ctx context.Context, options *CreateOptions) (*model.Organisation, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "creating a new organisation",
		slog.String("function", "create"),
	)
	if options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	return s.db.Create(ctx, &organisation.CreateOptions{
		Title:   options.Title,
		OwnerID: options.OwnerID,
	})
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*model.Organisation, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "listing all organisations",
		slog.String("function", "list"),
	)
	if options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	return s.db.List(ctx, &organisation.ListOptions{
		Title:          options.Title,
		Skip:           options.Skip,
		Limit:          options.Limit,
		OrderBy:        options.OrderBy,
		OrderDirection: options.OrderDirection,
	})
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*model.Organisation, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "retrieving a organisation",
		slog.String("function", "get"),
	)
	if ID == uuid.Nil {
		return nil, ErrInvalidOptions
	}
	return s.db.Get(ctx, ID)
}

func (s *service) Update(ctx context.Context, ID uuid.UUID, options *UpdateOptions) (*model.Organisation, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "updating a organisation",
		slog.String("function", "update"),
	)
	if ID == uuid.Nil {
		return nil, ErrInvalidorganisationID
	}
	if options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}
	return s.db.Update(ctx, ID, &organisation.UpdateOptions{
		Title: options.Title,
	})
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "deleting a organisation",
		slog.String("function", "delete"),
	)
	if ID == uuid.Nil {
		return ErrInvalidorganisationID
	}
	return s.db.Delete(ctx, ID)
}

func (s *service) Members() {}

func (s *service) Roles() {}
