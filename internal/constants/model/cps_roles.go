package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSRoles struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name,omitempty" bson:"name,omitempty"`
	Enabled   *bool         `json:"enabled" bson:"enabled"`
	CreatedAt time.Time     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
