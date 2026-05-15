package model

import "time"

type DonationCategoryOracle struct {
	ID             string    `bson:"id,omitempty" json:"id,omitempty"`
	CategoryName   string    `bson:"category_name" json:"category_name"`
	Icon           string    `bson:"donation_icon" json:"donation_icon"`
	IsDeleted      bool      `bson:"is_deleted" json:"is_deleted"`
	Enabled        bool      `bson:"enabled" json:"enabled"`
	CreatedAt      time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LastModifiedAt time.Time `bson:"last_modified_at,omitempty" json:"last_modified_at,omitempty"`
}
