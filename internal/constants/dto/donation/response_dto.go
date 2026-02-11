package donation

import "cbe-super-app-cps-action/internal/constants/types"

// types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

type DonationResponse struct {
	DonationCode        string                `json:"donation_code" bson:"donation_code"`
	CompanyID           string                `json:"company_id" bson:"company_id"`
	CategoryID          string                `json:"category_id" bson:"category_id"`
	Title               string                `json:"title" bson:"title"`
	IsFeatured          bool                  `json:"is_featured" bson:"is_featured"`
	Target              string                `json:"target" bson:"target"`
	DonationDescription string                `json:"donation_description" bson:"donation_description"`
	DonationImages      []types.DonationImage `json:"donation_images" bson:"donation_images"`
	CoverImage          string                `json:"cover_image,omitempty" bson:"cover_image,omitempty"` // URL for cover image
	EndDate             string                `json:"end_date" bson:"end_date"`
	StartDate           string                `json:"start_date" bson:"start_date"`
	Enabled             bool                  `json:"enabled" bson:"enabled"`
}
type Company struct {
	ID            string `json:"id" bson:"id"`
	CompanyName   string `json:"company_name" bson:"company_name"`
	CompanyLogo   string `json:"company_logo" bson:"company_logo"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	Enabled       bool   `json:"enabled" bson:"enabled"`
}

type Category struct {
	ID           string `json:"id" bson:"id"`
	CategoryName string `json:"category_name" bson:"category_name"`
	Icon         string `json:"icon" bson:"icon"`
}

type DonationListResponse struct {
	ID                  string                `json:"id" bson:"id"`
	DonationCode        string                `json:"donation_code" bson:"donation_code"`
	Company             Company               `json:"company" bson:"company"`
	Category            Category              `json:"category" bson:"category"`
	Title               string                `json:"title" bson:"title"`
	IsFeatured          bool                  `json:"is_featured" bson:"is_featured"`
	Target              string                `json:"target" bson:"target"`
	CurrentAmount       string                `json:"current_amount" bson:"current_amount"`
	DonationDescription string                `json:"donation_description" bson:"donation_description"`
	DonationImages      []types.DonationImage `json:"donation_images" bson:"donation_images"`
	CoverImage          string                `json:"cover_image,omitempty" bson:"cover_image,omitempty"` // URL for cover image
	EndDate             string                `json:"end_date" bson:"end_date"`
	StartDate           string                `json:"start_date" bson:"start_date"`
	IsDeleted           bool                  `json:"is_deleted" bson:"is_deleted"`
	CreatedAt           string                `json:"created_at" bson:"created_at"`
	LastModifiedAt      string                `json:"last_modified_at" bson:"last_modified_at"`
	Enabled             bool                  `json:"enabled" bson:"enabled"`
}
