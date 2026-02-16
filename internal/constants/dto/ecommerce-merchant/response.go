package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"time"
)

// type MerchantLookUpResponse struct {
// 	Status   string              `json:"status"`
// 	Company  types.Company       `json:"company"`
// 	Branches []types.Branch      `json:"branches"`
// 	Users    []types.UserAccount `json:"users"`
// }

type MerchantLookupAPIResponse struct {
	Status      int                      `json:"status"`
	Message     string                   `json:"message"`
	Data        []MerchantLookUpResponse `json:"data"`
	IsOrganizer bool                     `json:"is_organizer"`
	IsDelivery  bool                     `json:"is_delivery"`
	IsEcommerce bool                     `json:"is_ecommerce"`
}

type Branch struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	BranchID          string  `json:"branch_id"`
	BusinessType      *string `json:"business_type"`
	AccountNumber     string  `json:"cps_account_number"`
	AccountHolderName string  `json:"account_holder_name"`
	APIKey            string  `json:"api_key"`
	Email             *string `json:"email"`
	Phone             *string `json:"phone"`
}

type MerchantLookUpResponse struct {
	ID            int      `json:"id"`
	Banner        string   `json:"banner"`
	BusinessType  *string  `json:"business_type"`
	City          string   `json:"city"`
	AccountNumber string   `json:"cps_account_number"`
	Description   *string  `json:"description"`
	IsFeatured    bool     `json:"is_featured"`
	Logo          string   `json:"logo"`
	MerchantID    string   `json:"merchant_id"`
	Name          string   `json:"name"`
	Street        string   `json:"street"`
	Branches      []Branch `json:"branches"`
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
