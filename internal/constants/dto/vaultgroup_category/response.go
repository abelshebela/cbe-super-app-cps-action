package vaultgroupcategory

import "time"

type VaultGroupCategoryResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	CoverImage string     `json:"cover_image"`
	IsActive   bool       `json:"is_active"`
	IsDeleted  bool       `json:"is_deleted"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
