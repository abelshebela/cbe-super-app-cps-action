package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRole struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string        `json:"codde" bson:"code"`
	Name      string        `json:"name" bson:"name"`
	UpdatedAt time.Time     `json:"updted_at" bson:"updated_at"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}
