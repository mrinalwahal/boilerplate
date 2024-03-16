package db

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DB interface {
	Create(context.Context, *CreateOptions) (*Todo, error)
	Get(context.Context, uuid.UUID) (*Todo, error)
	List(context.Context, *ListOptions) ([]*Todo, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Todo, error)
	Delete(context.Context, uuid.UUID) error
}

func NewDB(db *gorm.DB) DB {
	return &DefaultDB{conn: db}
}

type DefaultDB struct {

	//	Database Connection
	conn *gorm.DB
}

func (db *DefaultDB) Create(ctx context.Context, options *CreateOptions) (*Todo, error) {
	txn := db.conn.WithContext(ctx)

	var payload Todo
	payload.Title = options.Title

	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *DefaultDB) Get(ctx context.Context, ID uuid.UUID) (*Todo, error) {
	txn := db.conn.WithContext(ctx)

	var payload Todo
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *DefaultDB) List(ctx context.Context, options *ListOptions) ([]*Todo, error) {
	txn := db.conn.WithContext(ctx)

	var payload []*Todo

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
	where := Todo{
		Title: options.Title,
	}

	if result := query.Where(&where).Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

func (db *DefaultDB) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*Todo, error) {
	txn := db.conn.WithContext(ctx)

	var payload Todo
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return db.Get(ctx, id)
}

func (db *DefaultDB) Delete(ctx context.Context, ID uuid.UUID) error {
	txn := db.conn.WithContext(ctx)

	var payload Todo
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}

func create(txn *gorm.DB, options *CreateOptions) (*Todo, error) {
	var payload Todo
	payload.Title = options.Title

	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func get(txn *gorm.DB, ID uuid.UUID) (*Todo, error) {
	var payload Todo
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func list(txn *gorm.DB, options *ListOptions) ([]*Todo, error) {
	var payload []*Todo

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
	where := Todo{
		Title: options.Title,
	}

	if result := query.Where(&where).Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

func update(txn *gorm.DB, id uuid.UUID, options *UpdateOptions) error {
	var payload Todo
	payload.ID = id
	result := txn.Model(&payload).Updates(options)
	return result.Error
}

func delete(txn *gorm.DB, ID uuid.UUID) error {
	var payload Todo
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}
