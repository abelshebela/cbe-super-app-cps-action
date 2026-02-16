package vault

import "mime/multipart"

// type CreateVaultGroupCategoryRequest struct {
// 	Name         string                `form:"name"`
// 	CategoryType string                `form:"category_type"`
// 	CoverImage   *multipart.FileHeader `form:"cover_image"`
// }

// type UpdateVaultGroupCategoryRequest struct {
// 	Name         string                `form:"name,omitempty"`
// 	CategoryType string                `form:"category_type"`
// 	CoverImage   *multipart.FileHeader `form:"cover_image,omitempty"`
// }

// Create request
type CreateTierDTO struct {
	Name         string `json:"name" validate:"required"`
	TierInterest string `json:"tier_interest" validate:"required"`
	Min          string `json:"min" validate:"required"`
	Max          string `json:"max" validate:"required"`
}

type CreateCategoryRequest struct {
	Name             string                `json:"name" validate:"required"`
	CoverImage       *multipart.FileHeader `form:"cover_image"`
	InterestType     string                `json:"interest_type" bson:"interest_type"`
	CategoryInterest string                `json:"category_interest" validate:"required"`
	Deadlock         *bool                 `json:"deadlock"`
	Tiers            []CreateTierDTO       `json:"tiers" validate:"required"`
}

type CreateWithdrawalRequest struct {
	LockedVaultID         string  `json:"locked_vault_id" validate:"required"`
	WithdrawalAmount      float64 `json:"withdrawal_amount" validate:"required"`
	WithdrawerName        string  `json:"withdrawer_name" validate:"required"`
	WithdrawerPhoneNumber string  `json:"withdrawer_phone_number" validate:"required"`
}

// Update request
type UpdateWithdrawalStatusRequest struct {
	WithdrawalStatus string `json:"withdrawal_status" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name             *string               `json:"name,omitempty"`
	CoverImage       *multipart.FileHeader `form:"cover_image,omitempty"`
	InterestType     *string               `json:"interest_type" bson:"interest_type"`
	CategoryInterest *string               `json:"category_interest,omitempty"`
	Deadlock         *bool                 `json:"deadlock,omitempty"`
	Tiers            []UpdateTierDTO       `json:"tiers,omitempty"`
}

type UpdateTierDTO struct {
	Name         *string `json:"name,omitempty"`
	Min          *string `json:"min,omitempty"`
	Max          *string `json:"max,omitempty"`
	TierInterest *string `json:"tier_interest,omitempty"`
}
