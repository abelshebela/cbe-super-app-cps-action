package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PermissionGroup struct {
	ID                 bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GroupName          string        `bson:"group_name,omitempty" json:"group_name,omitempty"`
	Permissions        interface{}   `bson:"permissions,omitempty" json:"permissions,omitempty"`
	PermissionCategory interface{}   `bson:"permission_category,omitempty" json:"permission_category,omitempty"`
	Role               string        `bson:"role,omitempty" json:"role,omitempty"`
	Realm              string        `bson:"realm,omitempty" json:"realm,omitempty"`
	Enabled            bool          `bson:"enabled,omitempty" json:"enabled,omitempty"`
	IsDeleted          bool          `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
	CreatedAt          time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LastModified       time.Time     `bson:"last_modified,omitempty" json:"last_modified,omitempty"`
}
