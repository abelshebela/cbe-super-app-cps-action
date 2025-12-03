package donation_company

type DonationCompanyResponse struct {
	CompanyName       string `json:"company_name" bson:"company_name"`
	CompanyLogo       string `json:"company_logo" bson:"company_logo"`
	AccountNumber     string `json:"account_number" bson:"account_number"`
	AccountHolderName string `json:"account_holder_name" bson:"account_holder_name"`
	PhoneNumber       string `json:"phone_number" bson:"phone_number"`
	Email             string `json:"email" bson:"email"`
	Address           string `json:"address" bson:"address"`
	Enabled           bool   `json:"enabled" bson:"enabled"`
}

// DonationCompanyListResponse is used for fetching donation companies
type DonationCompanyListResponse struct {
	ID                string `json:"id" bson:"id"`
	CompanyName       string `json:"company_name" bson:"company_name"`
	CompanyLogo       string `json:"company_logo" bson:"company_logo"`
	AccountNumber     string `json:"account_number" bson:"account_number"`
	AccountHolderName string `json:"account_holder_name" bson:"account_holder_name"`
	PhoneNumber       string `json:"phone_number" bson:"phone_number"`
	Email             string `json:"email" bson:"email"`
	Address           string `json:"address" bson:"address"`
	Enabled           bool   `json:"enabled" bson:"enabled"`
	IsDeleted         bool   `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         string `json:"created_at" bson:"created_at"`
	LastModifiedAt    string `json:"last_modified_at" bson:"last_modified_at"`
}
