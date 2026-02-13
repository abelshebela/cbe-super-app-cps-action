package sqlc

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"time"
)

type VaultTiers struct {
	ID           string `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name         string `json:"name" bson:"name"`
	CategoryID   string `json:"category_id" bson:"category_id"`
	TierInterest string `json:"tier_interest" bson:"tier_interest"`
	MinAmount    string `json:"min_amount" bson:"min_amount"`
	MaxAmount    string `json:"max_amount" bson:"max_amount"`
}

type VaultCategory struct {
	ID               string              `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name             string              `json:"name" bson:"name"`
	CoverImageURL    string              `json:"cover_image_url" bson:"cover_image_url"`
	InterestType     string              `json:"interest_type" bson:"interest_type"`
	CategoryInterest string              `json:"category_interest" bson:"category_interest"`
	Deadlock         bool                `json:"deadlock" bson:"deadlock"`
	IsActive         bool                `json:"is_active" bson:"is_active"`
	CreatedAt        time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at" bson:"updated_at"`
	TotalCount       int64               `json:"total_count"`
	Tiers            []imodel.VaultTiers `json:"tiers" bson:"tiers"`
}

