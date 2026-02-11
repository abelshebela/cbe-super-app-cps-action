package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ActionRole struct {
	ID              bson.ObjectID     `json:"_id" bson:"_id"`
	ActionCode      string            `json:"action_code" bson:"action_code"`
	ActionName      string            `json:"action_name" bson:"action_name"`
	AssignedMakers  []bson.ObjectID   `json:"assigned_makers" bson:"assigned_makers"`
	AssignedChecker [][]bson.ObjectID `json:"assigned_checkers" bson:"assigned_checkers"`
	Enabled         bool              `json:"enabled" bson:"enabled"`
	UpdatedAt       time.Time         `json:"updated_at" bson:"updated_at"`
	CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
}
