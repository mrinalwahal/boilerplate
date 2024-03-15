package authorization

import (
	"github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"
)

type Config struct {
	//	Database connection.
	DB *gorm.DB
	//	Authorization Model
	Model *model.Model
}
