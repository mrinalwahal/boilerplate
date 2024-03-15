package authorization

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"
)

type Service interface {
}

type service struct {

	//	Database connection.
	db *gorm.DB
	//	Authorization Model
	model *model.Model
	//	Enforcer
	enforcer *casbin.Enforcer
}
