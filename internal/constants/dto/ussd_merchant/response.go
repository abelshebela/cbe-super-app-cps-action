package ussd_merchant_dto

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UssdMerchantResponse struct {
	ID               bson.ObjectID `bson:"_id" json:"id"`
	Name             string        `bson:"name" json:"name"`
	PhoneNumber      string        `bson:"phone_number" json:"phone_number"`
	Email            string        `bson:"email" json:"email"`
	Service          string        `bson:"service" json:"service"`
	ServiceName      string        `bson:"service_name" json:"service_name"`
	AccountNumber    string        `bson:"account_number" json:"account_number"`
	Logo             string        `bson:"logo" json:"logo"`
	SettlementMethod string        `bson:"settlement_method" json:"settlement_method"`
	Enabled          bool          `bson:"enabled" json:"enabled"`
	CreatedAt        time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at" json:"updated_at"`
}
