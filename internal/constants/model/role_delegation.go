package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleDelegation struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Enable    bool          `json:"enabled" bson:"enabled"`
	StartAt   time.Time     `json:"start_at" bson:"start_at"`
	EndAt     time.Time     `json:"end_at" bson:"end_at"`
	JobTitle  string        `json:"job_title" bson:"job_title"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
