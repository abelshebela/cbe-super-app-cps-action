package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Icon struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Icon         string        `bson:"icon" json:"icon"`
	Enabled      bool          `bson:"enabled" json:"enabled"`
	IsDeleted    bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	LastModified time.Time     `bson:"last_modified" json:"last_modified"`
}
