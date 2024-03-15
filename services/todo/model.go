package todo

import (
	"github.com/mrinalwahal/boilerplate/storage"
)

type Todo struct {
	storage.Base

	//	Title of the todo.
	Title string `json:"title" gorm:"not null"`
}
