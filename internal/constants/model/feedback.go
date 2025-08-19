package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Feedback struct {
	ID        bson.ObjectID             `json:"id" bson:"_id,omitempty"`
	UserID    string                    `json:"user_id" bson:"user_id"`
	Responses map[string]types.Response `json:"responses" bson:"responses"`
	CreatedAt time.Time                 `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at" bson:"updated_at"`
}
