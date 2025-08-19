package donation_category

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DonationCategoryMapper maps DonationCategory model to BSON for database operations
func DonationCategoryMapper(data model.DonationCategory) bson.M {
	result := bson.M{}
	if data.CategoryName != "" {
		result["category_name"] = data.CategoryName
	}
	if data.Icon != "" {
		result["description"] = data.Icon
	}

	return result
}
