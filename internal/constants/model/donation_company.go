package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DonationCompany struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CompanyName    string        `json:"company_name" bson:"company_name"`
	CompanyLogo    string        `json:"company_logo" bson:"company_logo"`
	AccountNumber  string        `json:"account_number" bson:"account_number"`
	AccountHolderName string `json:"account_holder_name" bson:"account_holder_name"`
	Address        string        `json:"address" bson:"address"`
	PhoneNumber    string        `json:"phone_number" bson:"phone_number"`
	Email          string        `json:"email" bson:"email"`
	Enabled        bool          `json:"enabled" bson:"enabled"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}
