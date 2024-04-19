package model

import "github.com/google/uuid"

type Role struct {
	Base

	// ID of the organisation to which the role belongs.
	//
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	//
	// It is a required field.
	OrgID uuid.UUID `json:"org_id" gorm:"not null;type:uuid"`

	// Name of the role.
	//
	// Example: "Admin"
	//
	// It is a required field.
	Name string `json:"name" gorm:"not null;check:(length(name)>0)"`

	// Permissions this role has.
	//
	// Example: [{"operation": "create", "entity": "organisation"}, {"operation": "read", "entity": "member"}]
	//
	// It is a required field.
	Permissions []Permission `json:"permissions" gorm:"not null;type:json"`
}

type Relation string

const (
	Owner  Relation = "owner"
	Viewer Relation = "viewer"
	Editor Relation = "editor"
	Admin  Relation = "admin"
	Member Relation = "member"
)

type Permission struct {
	Operation Operation `json:"operation" gorm:"not null;check:(length(operation)>0)"`
	Entity    Entity    `json:"entity" gorm:"not null;check:(length(entity)>0)"`
}

type Operation string

const (
	Create Operation = "create"
	Read   Operation = "read"
	Update Operation = "update"
	Delete Operation = "delete"
)

type Entity string

const (
	organisation Entity = "organisation"
)
