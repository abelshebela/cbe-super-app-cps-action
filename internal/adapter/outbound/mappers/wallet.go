package mappers

import (
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToWalletDocument(domain *entity.Wallet) (*model.WalletDocument, error) {
	var objectID bson.ObjectID
	if domain.ID != "" {
		objID, err := bson.ObjectIDFromHex(domain.ID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_ID")
		}
		objectID = objID
	} else {
		objectID = bson.NewObjectID()
	}

	deletedAt := time.Time{}
	if domain.DeletedAt != nil {
		deletedAt = *domain.DeletedAt
	}

	return &model.WalletDocument{
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

func ToDomainWallet(model model.WalletDocument) *entity.Wallet {
	var deletedAt *time.Time
	if !model.DeletedAt.IsZero() {
		deletedAt = &model.DeletedAt
	}

	return &entity.Wallet{
		ID:             model.ID.Hex(),
		Name:           model.Name,
		Code:           model.Code,
		Avatar:         model.Avatar,
		Enabled:        model.Enabled,
		IsDeleted:      model.IsDeleted,
		CreatedAt:      model.CreatedAt,
		LastModifiedAt: model.LastModifiedAt,
		DeletedAt:      deletedAt,
	}
}
