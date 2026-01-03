package amount_based_auth

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func AuthTierMapper(authTier model.AuthTier) bson.M {
	return bson.M{
		"min_amount":    authTier.MinAmount,
		"max_amount":    authTier.MaxAmount,
		"method":        authTier.Method,
		"last_modified": authTier.LastModified,
	}
}
