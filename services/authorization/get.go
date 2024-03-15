package authorization

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func GetService(config *Config) Service {

	// Initialize a GORM adapter and use it in a Casbin enforcer.
	a, _ := gormadapter.NewAdapterByDB(config.DB)
	e, _ := casbin.NewEnforcer(config.Model, a)

	// Load the policy from DB.
	// e.LoadPolicy()

	return &service{
		db:       config.DB,
		model:    config.Model,
		enforcer: e,
	}
}

type Config struct {
	//	Database connection.
	DB *gorm.DB
	//	Authorization Model
	Model *model.Model
}
