package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"time"
)

type MerchantLookUpResponse struct {
	Status   string              `json:"status"`
	Company  types.Company       `json:"company"`
	Branches []types.Branch      `json:"branches"`
	Users    []types.UserAccount `json:"users"`
}

type MiniAppMerchantResponseDTO struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Type          string    `json:"type"`
	MerchantName  string    `json:"merchant_name"`
	AccountNumber string    `json:"account_number"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Enabled       bool      `json:"enabled"`
	IsDeleted     bool      `json:"is_deleted"`
	CreatedAt     time.Time `json:"created_at"`
	LastModified  time.Time `json:"last_modified"`
}

type PaginatedMiniAppResponseResponse struct {
	Data []MiniAppMerchantResponseDTO `json:"docs"`
	Meta types.PaginationMeta         `json:"meta"`
}
