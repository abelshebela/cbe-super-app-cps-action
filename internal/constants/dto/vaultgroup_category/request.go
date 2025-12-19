package vaultgroupcategory

import "mime/multipart"

type CreateVaultGroupCategoryRequest struct {
	Name         string                `form:"name"`
	CategoryType string                `form:"category_type"`
	CoverImage   *multipart.FileHeader `form:"cover_image"`
}

type UpdateVaultGroupCategoryRequest struct {
	Name         string                `form:"name,omitempty"`
	CategoryType string                `form:"category_type"`
	CoverImage   *multipart.FileHeader `form:"cover_image,omitempty"`
}
