package service

import (
	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/record/model"
	"gorm.io/gorm"
)

// create executes the transaction on the database to create a new record.
func (s *service) create(tx *gorm.DB, options *CreateOptions) (*model.Record, error) {
	var payload model.Record
	payload.Title = options.Title
	payload.UserID = options.UserID

	result := tx.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

// list executes the transaction on the database to fetch a list of records.
func (s *service) list(txn *gorm.DB, options *ListOptions) ([]*model.Record, error) {
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

// get executes the transaction on the database to fetch a record.
func (s *service) get(txn *gorm.DB, ID uuid.UUID) (*model.Record, error) {
	var payload model.Record
	payload.ID = ID

	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

// update executes the transaction on the database to update a record.
func (s *service) update(txn *gorm.DB, id uuid.UUID, options *UpdateOptions) error {
	var payload model.Record
	payload.ID = id

	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return result.Error
	}
	return nil
}

// delete executes the transaction on the database to delete a record.
func (s *service) delete(txn *gorm.DB, ID uuid.UUID) error {
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
