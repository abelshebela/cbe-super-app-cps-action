package donation_company

import (
	"mime/multipart"
)

type DonationCompanyRequest struct {
	CompanyName        string                `json:"company_name" bson:"company_name"`
	CompanyCode        string                `json:"company_code" bson:"company_code"`
	CompanyDescription string                `json:"company_description" bson:"company_description"`
	CompanyLogo        *multipart.FileHeader `json:"company_logo,omitempty" bson:"company_logo,omitempty"`
	PhoneNumber        string                `json:"phone_number" bson:"phone_number"`
	Email              string                `json:"email" bson:"email"`
	Address            string                `json:"address" bson:"address"`
}

// DonationCompanyCPSRequest is used for CPS actions where logo is a URL string
type DonationCompanyCPSRequest struct {
	ID                 string `json:"id,omitempty" bson:"id,omitempty"`
	CompanyName        string `json:"company_name" bson:"company_name"`
	CompanyCode        string `json:"company_code" bson:"company_code"`
	CompanyDescription string `json:"company_description" bson:"company_description"`
	CompanyLogo        string `json:"company_logo,omitempty" bson:"company_logo,omitempty"`
	PhoneNumber        string `json:"phone_number" bson:"phone_number"`
	Email              string `json:"email" bson:"email"`
	Address            string `json:"address" bson:"address"`
}
