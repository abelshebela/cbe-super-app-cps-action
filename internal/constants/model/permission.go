package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PermissionName string        `bson:"permission_name" json:"permission_name"`
	CreatedAt      time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
