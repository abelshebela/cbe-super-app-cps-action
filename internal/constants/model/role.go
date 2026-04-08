package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	JobTitle  string        `json:"job_title" bson:"job_title"`
	Role      string        `json:"role" bson:"role"` // job role code (matches job_roles.code)
	Enabled   bool          `json:"enabled" bson:"enabled"`
	UpdateAt  time.Time     `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	Type      string        `json:"type" bson:"type"` // from job_roles.type (aggregation)
	RoleCode  string        `json:"role_code,omitempty" bson:"role_code,omitempty"`
	RoleName  string        `json:"role_name,omitempty" bson:"role_name,omitempty"`
}
