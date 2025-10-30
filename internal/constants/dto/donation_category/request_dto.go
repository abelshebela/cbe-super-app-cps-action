package donation_category

import (
	"mime/multipart"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DonationCategoryRequest struct {
	CategoryName string                `json:"category_name" bson:"category_name"`
	Icon         *multipart.FileHeader `json:"donation_icon,omitempty" bson:"donation_icon,omitempty"`
}

// DonationCategoryCPSRequest is used for CPS actions where icon is a URL string
type DonationCategoryCPSRequest struct {
	ID           string `json:"id,omitempty" bson:"id,omitempty"`
	CategoryName string `json:"category_name" bson:"category_name"`
	Icon         string `json:"donation_icon,omitempty" bson:"donation_icon,omitempty"`
}

type EnableDonationCategoryRequest struct {
	Enabled bool `json:"enabled" bson:"enabled"`
}


