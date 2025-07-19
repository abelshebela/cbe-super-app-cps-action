package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (w *Wallet) parseObjectID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, fmt.Errorf(error_codes.InvalidID)
	}
	return objectID, nil
}

func (w *Wallet) handleError(operation string, err error, errorCode string) error {
	w.logger.Errorf("%s: %v", operation, err)
	return errors.New(errorCode)
}

func (w *Wallet) toDomain(wallet entity.WalletDocument) *entity.Wallet {
	return &entity.Wallet{
		ID:             wallet.ID.Hex(),
		Name:           wallet.Name,
		Code:           wallet.Code,
		Avatar:         wallet.Avatar,
		Enabled:        wallet.Enabled,
		IsDeleted:      wallet.IsDeleted,
		CreatedAt:      wallet.CreatedAt,
		LastModifiedAt: wallet.LastModifiedAt,
		DeletedAt:      wallet.DeletedAt,
	}
}

func (w *Wallet) toDocument(wallet *entity.Wallet) *entity.WalletDocument {
	return &entity.WalletDocument{
		ID:             bson.NewObjectID(),
		Name:           wallet.Name,
		Code:           wallet.Code,
		Avatar:         wallet.Avatar,
		Enabled:        wallet.Enabled,
		IsDeleted:      wallet.IsDeleted,
		CreatedAt:      wallet.CreatedAt,
		LastModifiedAt: wallet.LastModifiedAt,
		DeletedAt:      wallet.DeletedAt,
	}
}

// createCPSAction creates a new CPS action with common fields
func (w *Wallet) createCPSAction(cpsReq model.CreateCPSAction, actionType, requestAction string, previous, current any) model.CPSAction {
	return model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.FullName,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       actionType,
		RequestAction:    requestAction,
		PreviousAction:   previous,
		CurrentAction:    current,
		MakerActionTime:  time.Now(),
	}
}

// findWallet retrieves a wallet with standardized error handling
func (w *Wallet) findWallet(ctx context.Context, id string, projection bson.M) (*entity.Wallet, error) {
	objectId, err := w.parseObjectID(id)
	if err != nil {
		return nil, w.handleError("parse object ID", err, error_codes.InvalidID)
	}
	filter := bson.M{
		"_id":        objectId,
		"is_deleted": false,
	}
	wallet, err := w.walletDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, w.handleError("find wallet", err, error_codes.WalletNotFound)
		}
		return nil, w.handleError("find wallet", err, error_codes.UnhandledServerError)
	}
	return w.toDomain(*wallet), nil
}
