package todo

import (
	"github.com/google/uuid"
	"github.com/mrinalwahal/boilerplate/storage"
)

type Todo struct {
	storage.Base

	//	Title of the todo.
	Title string `json:"title" gorm:"not null"`

	//	ID of the user who created the todo.
	UserID uuid.UUID `json:"user_id" gorm:"not null"`
}
