package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func AuthTierMapper(authTier model.AuthTier) bson.M {
	return bson.M{
		"$set": bson.M{
			"min_amount":    authTier.MinAmount,
			"max_amount":    authTier.MaxAmount,
			"method":        authTier.Method,
			"last_modified": time.Now(),
		},
	}
}