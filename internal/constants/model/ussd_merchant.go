package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UssdMerchant struct {
	ID               bson.ObjectID              `bson:"_id" json:"id"`
	MerchantCode     string                     `bson:"merchant_code" json:"merchant_code"`
	Credential       string                     `bson:"credential" json:"credential"`
	Name             string                     `bson:"name" json:"name"`
	SettlementMethod constants.SettlementMethod `bson:"settlement_method" json:"settlement_method"`
	PhoneNumber      string                     `bson:"phone_number" json:"phone_number"`
	Email            string                     `bson:"email" json:"email"`
	Service          string                     `bson:"service" json:"service"`
	AccountNumber    string                     `bson:"account_number" json:"account_number"`
	Logo             string                     `bson:"logo" json:"logo"`
	Enabled          bool                       `bson:"enabled" json:"enabled"`
	CreatedAt        time.Time                  `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time                  `bson:"updated_at" json:"updated_at"`
}
