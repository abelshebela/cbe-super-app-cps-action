package donation_company

import (
	"mime/multipart"
)

type DonationCompanyRequest struct {
	CompanyName   string                `json:"company_name" bson:"company_name"`
	CompanyLogo   *multipart.FileHeader `json:"company_logo,omitempty" bson:"company_logo,omitempty"`
	AccountNumber string                `json:"account_number" bson:"account_number"`
	PhoneNumber   string                `json:"phone_number" bson:"phone_number"`
	Email         string                `json:"email" bson:"email"`
	Address       string                `json:"address" bson:"address"`
}

// DonationCompanyCPSRequest is used for CPS actions where logo is a URL string
type DonationCompanyCPSRequest struct {
	ID            string `json:"id,omitempty" bson:"id,omitempty"`
	CompanyName   string `json:"company_name" bson:"company_name"`
	CompanyLogo   string `json:"company_logo,omitempty" bson:"company_logo,omitempty"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	AccountHolderName string `json:"account_holder_name" bson:"account_holder_name"`
	PhoneNumber   string `json:"phone_number" bson:"phone_number"`
	Email         string `json:"email" bson:"email"`
	Address       string `json:"address" bson:"address"`
}
