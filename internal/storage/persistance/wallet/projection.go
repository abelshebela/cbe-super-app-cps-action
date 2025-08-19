package wallet

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// WalletMapper maps a Wallet model to a bson.M for updates
func WalletMapper(wallet model.Wallet) bson.M {
	return bson.M{
		"$set": bson.M{
			"name":             wallet.Name,
			"code":             wallet.Code,
			"avatar":           wallet.Avatar,
			"enabled":          wallet.Enabled,
			"last_modified_at": time.Now(),
		},
	}
} 