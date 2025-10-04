package model

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"

	"cbe-super-app-cps-action/internal/constants/types"
)

type PaginatedDonationCategoryResponse struct {
	Docs []donation_category.DonationCategoryListResponse `json:"docs"`
	Meta types.PaginationMeta                             `json:"meta"`
}

type PaginatedDonationResponse struct {
	Docs []donation.DonationListResponse `json:"docs"`
	Meta types.PaginationMeta            `json:"meta"`
}

type PaginatedDonationCompanyResponse struct {
	Docs []donation_company.DonationCompanyListResponse `json:"docs"`
	Meta types.PaginationMeta                           `json:"meta"`
}

type PaginatedWalletResponse struct {
	Docs []Wallet             `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}

type PaginatedNotificationResponse struct {
	Docs []Notification       `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}

type PaginatedArchieveUserResponse struct {
	Docs []ArchivedUser       `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}

type PaginatedDepartmentResponse struct {
	Docs []Department             `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}


type PaginatedEventResponse struct {
	Docs []Event             `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}
