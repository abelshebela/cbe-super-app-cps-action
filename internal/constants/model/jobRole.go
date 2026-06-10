package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRole struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	Code      string        `json:"code" bson:"code"`
	JobTitle  string        `json:"job_title" bson:"job_title"`
	Role      string        `json:"role" bson:"role"` // job role code (matches job_roles.code)
	Enabled   bool          `json:"enabled" bson:"enabled"`
	UpdateAt  time.Time     `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	// Type, RoleCode, RoleName are populated on read when using RoleRepository aggregation (GET /job_roles*), not stored on the roles document.
	Type      string    `json:"type" bson:"type"`
	RoleCode  string    `json:"role_code,omitempty" bson:"role_code,omitempty"`
	RoleName  string    `json:"role_name,omitempty" bson:"role_name,omitempty"`
	IsDeleted bool      `json:"is_deleted" bson:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
