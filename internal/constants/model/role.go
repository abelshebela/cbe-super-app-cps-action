package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Level     string        `json:"level" bson:"level"`
	Enabled   bool          `json:"enabled" bson:"enabled"`
	UpdateAt  time.Time     `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}
