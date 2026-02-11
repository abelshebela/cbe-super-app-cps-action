package auth_tier

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AuthTierMapper maps AuthTier model to BSON for database operations
func AuthTierMapper(data model.AuthTier) bson.M {
	result := bson.M{}

	if data.MinAmount != 0 {
		result["min_amount"] = data.MinAmount
	}
	if data.MaxAmount != 0 {
		result["max_amount"] = data.MaxAmount
	}
	return result
}
