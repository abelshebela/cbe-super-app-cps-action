package wallet

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToWallet(domain *model.Wallet) (*model.Wallet, error) {
	var objectID bson.ObjectID
	if !domain.ID.IsZero() {
		objectID = domain.ID
	} else {
		objectID = bson.NewObjectID()
	}

	deletedAt := time.Time{}
	if !domain.DeletedAt.IsZero() {
		deletedAt = domain.DeletedAt
	}

	return &model.Wallet{
		ID:             objectID,
		Name:           domain.Name,
		Code:           domain.Code,
		Avatar:         domain.Avatar,
		Enabled:        domain.Enabled,
		IsDeleted:      domain.IsDeleted,
		CreatedAt:      domain.CreatedAt,
		LastModifiedAt: domain.LastModifiedAt,
		DeletedAt:      deletedAt,
	}, nil
}

func ToWalletDocument(wallet model.Wallet) (*model.Wallet, error) {
	var deletedAt *time.Time
	if !wallet.DeletedAt.IsZero() {
		deletedAt = &wallet.DeletedAt
	}

	return &model.Wallet{
		ID:             wallet.ID,
		Name:           wallet.Name,
		Code:           wallet.Code,
		Avatar:         wallet.Avatar,
		Enabled:        wallet.Enabled,
		IsDeleted:      wallet.IsDeleted,
		CreatedAt:      wallet.CreatedAt,
		LastModifiedAt: wallet.LastModifiedAt,
		DeletedAt: func() time.Time {
			if deletedAt != nil {
				return *deletedAt
			}
			return time.Time{}
		}(),
	}, nil
}
