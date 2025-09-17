package bank

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BankMapper maps Bank model to BSON for database operations
func BankMapper(data model.Bank) bson.M {
	result := bson.M{}
	if data.Name != "" {
		result["name"] = data.Name
	}
	if data.Code != "" {
		result["code"] = data.Code
	}
	if data.BIC != "" {
		result["description"] = data.BIC
	}
	if data.Logo != "" {
		result["logo"] = data.Logo
	}

	result["enabled"] = data.Enabled

	return result
}
