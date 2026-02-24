package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DonationCategory struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryName   string        `json:"category_name" bson:"category_name"`
	Icon           string        `json:"donation_icon" bson:"donation_icon"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	Enabled        bool          `json:"enabled" bson:"enabled"`
}
