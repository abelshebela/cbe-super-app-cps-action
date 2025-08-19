package donation_company

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DonationCompanyMapper maps DonationCompany model to BSON for database operations
func DonationCompanyMapper(data model.DonationCompany) bson.M {
	result := bson.M{}
	if data.CompanyName != "" {
		result["company_name"] = data.CompanyName
	}
	if data.CompanyLogo != "" {
		result["company_logo"] = data.CompanyLogo
	}
	if data.AccountNumber != "" {
		result["account_number"] = data.AccountNumber
	}

	return result
}
