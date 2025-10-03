package wallet

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToWalletDocument(wallet model.Wallet) (*model.Wallet, error) {
	if wallet.ID == bson.NilObjectID {
		wallet.ID = bson.NewObjectID()
	}

	if wallet.CreatedAt.IsZero() {
		wallet.CreatedAt = time.Now()
	}
	if wallet.LastModifiedAt.IsZero() {
		wallet.LastModifiedAt = time.Now()
	}
	if wallet.IsDeleted {
		wallet.IsDeleted = false
	}
	if wallet.Enabled {
		wallet.Enabled = false

	}

	return &wallet, nil
}

func UpdateMapper(wallet model.Wallet) bson.M {

	update := bson.M{"last_modified_at": time.Now()}
	if wallet.Name != "" {
		update["name"] = wallet.Name
	}
	if wallet.Code != "" {
		update["code"] = wallet.Code
	}
	if wallet.Avatar != "" {
		update["avatar"] = wallet.Avatar
	}

	return update
}
