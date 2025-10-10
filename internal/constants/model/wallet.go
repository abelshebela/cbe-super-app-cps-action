package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Wallet struct {
	ID             bson.ObjectID `json:"id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Code           string        `json:"code" bson:"code"`
	Avatar         string        `json:"avatar" bson:"avatar"`
	Enabled        bool          `json:"enabled" bson:"enabled"`
	Self           bool          `json:"self" bson:"self"`
	Other          bool          `json:"other" bson:"other"`
	Agent          bool          `json:"agent" bson:"agent"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      time.Time     `json:"deleted_at" bson:"deleted_at"`
}
