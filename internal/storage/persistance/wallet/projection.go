package wallet

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

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
	wallet.IsDeleted = false
	wallet.Enabled = false
	return &wallet, nil
}

func UpdateMapper(wallet model.Wallet) bson.M {

	update := bson.M{"last_modified_at": time.Now()}
	if wallet.Name != "" {
		update["name"] = wallet.Name
	}
	if wallet.UniqueCode != "" {
		update["unique_code"] = wallet.UniqueCode
	}
	if wallet.Avatar != "" {
		update["avatar"] = wallet.Avatar
	}

	update["services.self"] = wallet.Services.Self
	update["services.other"] = wallet.Services.Other
	update["services.agent"] = wallet.Services.Agent
	return update
}
