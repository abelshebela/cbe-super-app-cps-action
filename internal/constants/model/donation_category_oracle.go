package model

import "time"

type DonationCategoryOracle struct {
	ID             string    `json:"id,omitempty"`
	CategoryName   string    `json:"category_name"`
	Icon           string    `json:"donation_icon"`
	IsDeleted      bool      `json:"is_deleted"`
	Enabled        bool      `json:"enabled"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	LastModifiedAt time.Time `json:"last_modified_at,omitempty"`
}
