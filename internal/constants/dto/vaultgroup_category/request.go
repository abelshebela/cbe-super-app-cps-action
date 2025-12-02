package vaultgroupcategory

import "mime/multipart"

type CreateVaultGroupCategoryRequest struct {
	Name       string                `form:"name" validate:"required"`
	CoverImage *multipart.FileHeader `form:"cover_image" validate:"required"`
}

type UpdateVaultGroupCategoryRequest struct {
	Name       *string               `form:"name,omitempty"`
	CoverImage *multipart.FileHeader `form:"cover_image,omitempty"`
}
