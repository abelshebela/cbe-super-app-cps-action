package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniAppProductCode struct {
	ID            bson.ObjectID `bson:"_id" json:"_id"`
	Name          string        `bson:"name" json:"name"`
	ChargeCode    string        `bson:"charge_code" json:"charge_code"`
	CommisionCode string        `bson:"commission_code" json:"commission_code"`
	ServiceCode   string        `bson:"service_code" json:"service_code"`
	IsEnabled     bool          `bson:"is_enabled" json:"is_enabled"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}
