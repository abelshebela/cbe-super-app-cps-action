package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	PermissionName string `bson:"permissionName" json:"permissionName"`
}

type PermissionGroup struct {
	ID                 bson.ObjectID        `bson:"_id,omitempty" json:"_id,omitempty"`
	GroupName          string               `bson:"groupName" json:"groupName"`
	Permissions        []Permission         `bson:"permissions,omitempty" json:"permissions,omitempty"`
	PermissionCategory []PermissionCategory `bson:"permissionCategory" json:"permissionCategory"`
	Role               string               `bson:"role,omitempty" json:"role,omitempty"`
	Realm              string               `bson:"realm,omitempty" json:"realm,omitempty"`
	Enabled            bool                 `bson:"enabled,omitempty" json:"enabled,omitempty"`
	IsDeleted          bool                 `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	CreatedAt          time.Time            `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	LastModified       time.Time            `bson:"lastModified,omitempty" json:"lastModified,omitempty"`
}
