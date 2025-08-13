package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Donation struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id"`
	DonationCode        string        `json:"donation_code" bson:"donation_code"`
	CompanyID           bson.ObjectID `json:"company_id" bson:"company_id"`
	CategoryID          bson.ObjectID `json:"category_id" bson:"category_id"`
	Title               string        `json:"title" bson:"title"`
	IsFeatured          bool          `json:"is_featured" bson:"is_featured"`
	Target              int           `json:"target" bson:"target"`
	DonationDescription string        `json:"donation_description" bson:"donation_description"`
	DonationImages      []string      `json:"donation_images" bson:"donation_images"`
	EndDate             time.Time     `json:"end_date" bson:"end_date"`
	StartDate           time.Time     `json:"start_date" bson:"start_date"`
	IsDeleted           bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt      time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}

type DonationCategory struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryName   string        `json:"category_name" bson:"category_name"`
	Icon           string        `json:"donation_icon" bson:"donation_icon"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}

type DonationCompany struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CompanyName    string        `json:"company_name" bson:"company_name"`
	CompanyLogo    string        `json:"company_logo" bson:"company_logo"`
	AccountNumber  string        `json:"account_number" bson:"account_number"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}
