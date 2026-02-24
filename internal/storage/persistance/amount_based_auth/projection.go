<<<<<<< HEAD
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
=======
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
>>>>>>> 3640c5b5dde77249222ec7dd9f90a5770ab0dbc4
