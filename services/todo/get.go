package todo

import "gorm.io/gorm"

// Initializes and gets the service with the supplied database connection.
func GetService(db *gorm.DB) Service {
	return &service{db: db}
}
