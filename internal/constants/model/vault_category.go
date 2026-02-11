package model

import (
	"time"
)

type VaultCategory struct {
	ID           string     `json:"id" bson:"id"`
	Name         string     `json:"name" bson:"name"`
	CategoryType string     `json:"category_type" bson:"category_type"`
	CoverImage   string     `json:"cover_image" bson:"cover_image"`
	InterestType string     `json:"interest_type" bson:"interest_type"`
	Interest     string     `json:"interest" bson:"interest"`
	Deadlock     bool       `json:"deadlock" bson:"deadlock"`
	IsActive     bool       `json:"is_active" bson:"is_active"`
	IsDeleted    bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
