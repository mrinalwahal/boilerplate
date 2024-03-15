package authorization

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type RBAC struct {
	//	Database connection.
	db *gorm.DB
	//	Authorization Model
	model *model.Model
	//	Enforcer
	enforcer *casbin.Enforcer
}

func NewRBAC(config *Config) Service {

	// Initialize a GORM adapter and use it in a Casbin enforcer.
	a, _ := gormadapter.NewAdapterByDB(config.DB)
	e, _ := casbin.NewEnforcer(config.Model, a)

	// Load the policy from DB.
	// e.LoadPolicy()

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	return &RBAC{
		db:       config.DB,
		model:    &m,
		enforcer: e,
	}
}

func (r *RBAC) Enforce(sub, obj, act string) (bool, error) {
	return r.enforcer.Enforce(sub, obj, act)
}
