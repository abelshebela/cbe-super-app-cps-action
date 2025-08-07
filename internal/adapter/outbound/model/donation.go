package model

import (
	"mime/multipart"

	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Donation struct {
	ID                  bson.ObjectID           `bson:"_id,omitempty"`
	CompanyID           string                  `json:"company_id"`
	CategoryID          string                  `json:"category_id"`
	Title               string                  `json:"title"`
	IsFeatured          bool                    `json:"is_featured"`
	Target              float64                 `json:"donation_amount"`
	DonationDescription string                  `json:"donation_description"`
	DonationImages      []*multipart.FileHeader `json:"donation_images"`
	EndDate             time.Time               `json:"end_date"`
	StartDate           time.Time               `json:"start_date"`
	CreatedAt           time.Time               `bson:"created_at"`
	DeletedAt           time.Time               `bson:"deleted_at"`
	LastModifiedAt      time.Time               `bson:"last_modified_at"`
}

type DonationCategory struct {
	CategoryName   string    `json:"category_name"`
	Icon           string    `json:"donation_icon"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `bson:"created_at"`
	LastModifiedAt time.Time `bson:"last_modified_at"`
}

type DonationCompany struct {
	CompanyName    string                `json:"company_name"`
	CompanyLogo    *multipart.FileHeader `json:"company_logo"`
	AccountNumber  string                `json:"account_number"`
	CreatedAt      time.Time             `bson:"created_at"`
	DeletedAt      time.Time             `bson:"deleted_at"`
	LastModifiedAt time.Time             `bson:"last_modified_at"`
}
