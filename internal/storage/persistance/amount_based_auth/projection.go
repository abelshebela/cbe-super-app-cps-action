package amount_based_auth

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func AuthTierMapper(authTier local_model.AuthTier) bson.M {
	return bson.M{
		"currency":      authTier.Currency,
		"min_amount":    authTier.MinAmount,
		"max_amount":    authTier.MaxAmount,
		"method":        authTier.Method,
		"last_modified": authTier.LastModified,
	}
}
