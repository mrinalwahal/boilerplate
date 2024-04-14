//go:generate mockgen -destination=db_mock.go -source=db.go -package=db
package db

import (
	"context"

	"github.com/google/uuid"
	model "github.com/mrinalwahal/boilerplate/organisation/model"
)

// DB interface declares the signature of the database layer.
type DB interface {
	Create(context.Context, *CreateOptions) (*model.Organisation, error)
	List(context.Context, *ListOptions) ([]*model.Organisation, error)
	Get(context.Context, uuid.UUID) (*model.Organisation, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*model.Organisation, error)
	Delete(context.Context, uuid.UUID) error
}
