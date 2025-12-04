package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"time"
)

type MiniAppMerchantResponseDTO struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Type          string    `json:"type"`
	MerchantName  string    `json:"merchant_name"`
	KYC           KYCDTO    `json:"kyc"`
	AccountNumber string    `json:"account_number"`
	Enabled       bool      `json:"enabled"`
	IsDeleted     bool      `json:"is_deleted"`
	CreatedAt     time.Time `json:"created_at"`
	LastModified  time.Time `json:"last_modified"`
}

type PaginatedMiniAppResponseResponse struct {
	Data []MiniAppMerchantResponseDTO `json:"docs"`
	Meta types.PaginationMeta         `json:"meta"`
}
