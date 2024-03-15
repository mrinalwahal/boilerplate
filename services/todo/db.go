package todo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
