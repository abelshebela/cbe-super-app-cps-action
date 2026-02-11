package bank

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BankMapper maps Bank model to BSON for database operations
func BankMapper(data model.Bank) bson.M {
	result := bson.M{}
	if data.Name != "" {
		result["name"] = data.Name
	}
	if data.BICCode != "" {
		result["bic_code"] = data.BICCode
	}
	if data.Type != "" {
		result["type"] = data.Type
	}
	if data.Logo != "" {
		result["logo"] = data.Logo
	}
	result["enabled"] = data.Enabled
	return result
}
