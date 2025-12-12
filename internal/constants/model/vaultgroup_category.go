package model

import "time"

type VaultGroupCategory struct {
	ID         string     `json:"id" bson:"id"`
	Name       string     `json:"name" bson:"name"`
	CoverImage string     `json:"cover_image" bson:"cover_image"`
	IsActive   bool       `json:"is_active" bson:"is_active"`
	IsDeleted  bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	CreatedBy  string     `json:"created_by" bson:"created_by"`
	UpdatedBy  string     `json:"updated_by" bson:"updated_by"`
}

type UpdateVaultGroupCategory struct {
	Name       string    `json:"name" bson:"name"`
	CoverImage *string   `json:"cover_image" bson:"cover_image"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy  *string   `json:"updated_by" bson:"updated_by"`
}
