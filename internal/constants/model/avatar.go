package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Avatar struct {
	ID             bson.ObjectID `json:"id" bson:"_id"`
	Avatar         string        `json:"avatar" bson:"avatar"`
	Label          string        `json:"label" bson:"label"`
	Enable         bool          `json:"enable" bson:"enable"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time    `json:"deleted_at" bson:"deleted_at"`
}
