package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Donation struct {
	ID                  bson.ObjectID         `bson:"_id,omitempty" json:"id"`
	DonationCode        string                `json:"donation_code" bson:"donation_code"`
	CompanyID           bson.ObjectID         `json:"company_id" bson:"company_id"`
	CategoryID          bson.ObjectID         `json:"category_id" bson:"category_id"`
	Title               string                `json:"title" bson:"title"`
	IsFeatured          bool                  `json:"is_featured" bson:"is_featured"`
	Target              string                `json:"target" bson:"target"`
	CurrentAmount       string                `json:"current_amount" bson:"current_amount"`
	DonationDescription string                `json:"donation_description" bson:"donation_description"`
	DonationImages      []types.DonationImage `json:"donation_images" bson:"donation_images"`
	CoverImage          string                `json:"cover_image,omitempty" bson:"cover_image,omitempty"`
	EndDate             time.Time             `json:"end_date" bson:"end_date"`
	StartDate           time.Time             `json:"start_date" bson:"start_date"`
	IsDeleted           bool                  `json:"is_deleted" bson:"is_deleted"`
	CreatedAt           time.Time             `json:"created_at" bson:"created_at"`
	LastModifiedAt      time.Time             `json:"last_modified_at" bson:"last_modified_at"`
	Enabled             bool                  `json:"enabled" bson:"enabled"`
}
