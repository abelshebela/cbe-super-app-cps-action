// Package wallet provides persistence and CPS action handling for wallet entities.
package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/wallet"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Wallet struct {
	walletDal dal.MongoDal[entity.WalletDocument, entity.WalletDocument]
	cpsDal    dal.MongoDal[model.CPSAction, model.CPSAction]
	logger    utils.Logger
}

func InitWalletPersistence(client *mongo.Client, database string, collections []string, logger utils.Logger) outbound.WalletPersistence {
	return &Wallet{
		walletDal: dal.NewMongoDal[entity.WalletDocument, entity.WalletDocument](client, database, collections[0]),
		cpsDal:    dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collections[1]),
		logger:    logger,
	}
}

func (w *Wallet) CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error {
	filter := bson.M{
		"maker_phone_number": cpsReq.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsReq.Department,
		"request_action":     cpsReq.RequestAction,
	}

	existingAction, err := w.cpsDal.FindOne(ctx, filter, bson.M{"action_code": 1})
	if err != nil && err != mongo.ErrNoDocuments {
		return w.handleError("find cps action", err, error_codes.UnhandledServerError)
	}
	if existingAction != nil {
		w.logger.Infof("pending cps action present for user: %s, code: %s, dept: %s",
			cpsReq.MakerUser.FullName, cpsReq.MakerUser.UserCode, cpsReq.Department)
		return w.handleError("check cps action", nil, error_codes.PendingCPSActionExists)
	}
	return nil
}

func (w *Wallet) CreateWallet(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.cpsDal.InsertOne(ctx, w.createCPSAction(
		cpsReq,
		string(model.ActionCreate),
		string(model.RequestCreateWallet),
		nil,
		cpsReq.ActionData,
	))
	if err != nil {
		return nil, w.handleError("create cps action", err, error_codes.UnhandledServerError)
	}
	return &cpsAction, nil
}

func (w *Wallet) UpdateWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	wallet, err := w.findWallet(ctx, id, bson.M{"name": 1, "code": 1, "bic": 1})
	if err != nil {
		return nil, err
	}

	cps, err := w.cpsDal.InsertOne(ctx, w.createCPSAction(
		cpsReq,
		string(model.ActionUpdate),
		string(model.RequestUpdateWallet),
		bson.M{
			"name": wallet.Name,
			"code": wallet.Code,
		},
		cpsReq.ActionData,
	))
	if err != nil {
		return nil, w.handleError("create cps action", err, error_codes.UnhandledServerError)
	}
	return &cps, nil
}

func (w *Wallet) DeleteWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	wallet, err := w.findWallet(ctx, id, bson.M{"name": 1, "code": 1, "avatar": 1})
	if err != nil {
		return nil, err
	}

	cps, err := w.cpsDal.InsertOne(ctx, w.createCPSAction(
		cpsReq,
		string(model.ActionDelete),
		string(model.RequestDeleteWallet),
		bson.M{
			"id":         wallet.ID,
			"name":       wallet.Name,
			"code":       wallet.Code,
			"avatar":     wallet.Avatar,
			"is_deleted": wallet.IsDeleted,
		},
		bson.M{"is_deleted": true},
	))
	if err != nil {
		return nil, w.handleError("create cps action", err, error_codes.UnhandledServerError)
	}
	return &cps, nil
}

func (w *Wallet) EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	wallet, err := w.findWallet(ctx, id, bson.M{"name": 1, "code": 1, "enabled": 1})
	if err != nil {
		return nil, err
	}

	cps, err := w.cpsDal.InsertOne(ctx, w.createCPSAction(
		cpsReq,
		string(model.ActionUpdate),
		string(requestAction),
		bson.M{
			"name":    wallet.Name,
			"code":    wallet.Code,
			"enabled": wallet.Enabled,
		},
		cpsReq.ActionData,
	))
	if err != nil {
		return nil, w.handleError("create cps action", err, error_codes.UnhandledServerError)
	}
	return &cps, nil
}

func (w *Wallet) GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error) {
	filter := bson.M{"is_deleted": false}
	ifr := bson.M{}
	if filterParams.Filters != "" {
		filter["status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	walletsDoc, err := w.walletDal.FindAllWithPagination(ctx, filter, ifr, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, w.handleError("find wallets", err, "WALLET_NOT_FOUND")
		}
		return nil, w.handleError("find wallets", err, error_codes.UnhandledServerError)
	}

	total, err := w.walletDal.TotalCount(ctx, bson.M{})
	if err != nil {
		return nil, w.handleError("count wallets", err, error_codes.UnhandledServerError)
	}

	var wallets []*entity.Wallet
	for _, wallet := range walletsDoc {
		wallets = append(wallets, w.toDomain(*wallet))
	}

	return &entity.WalletResponse{
		Page:    filterParams.Page,
		Wallets: wallets,
		Limit:   constant.DefaultPerPage,
		Total:   total,
	}, nil
}

func (w *Wallet) GetWallet(ctx context.Context, id string) (*entity.Wallet, error) {
	return w.findWallet(ctx, id, bson.M{})
}

func (w *Wallet) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error) {
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"action_status":       model.ActionApproved,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := w.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, w.handleError("find cps action", err, "WALLET_NOT_FOUND")
		}
		return nil, w.handleError("update cps action", err, error_codes.UnhandledServerError)
	}

	var actionData entity.Wallet
	data, err := bson.Marshal(cpsAction.PreviosAction)
	if err != nil {
		return nil, w.handleError("marshal action data", err, error_codes.InvalidActionData)
	}

	if err := bson.Unmarshal(data, &actionData); err != nil {
		return nil, w.handleError("unmarshal action data", err, error_codes.InvalidActionData)
	}

	switch cpsAction.ActionType {
	case string(model.ActionCreate):
		newWallet := w.toDocument(&entity.Wallet{
			Name:      actionData.Name,
			Avatar:    actionData.Avatar,
			Code:      actionData.Code,
			CreatedAt: time.Now(),
		})

		wallet, err := w.walletDal.InsertOne(ctx, *newWallet)
		if err != nil {
			return nil, w.handleError("create wallet", err, error_codes.UnhandledServerError)
		}
		cpsAction.CurrentAction = wallet
		return &cpsAction, nil

	case string(model.ActionUpdate):
		update := bson.M{
			"last_modified_at": time.Now(),
		}
		if actionData.Name != "" {
			update["name"] = actionData.Name
		}
		if actionData.Code != "" {
			update["code"] = actionData.Code
		}
		if cpsAction.RequestAction == string(model.RequestEnableWallet) {
			update["enabled"] = true
		}
		if cpsAction.RequestAction == string(model.RequestDisableWallet) {
			update["enabled"] = false
		}

		wallet, err := w.walletDal.UpdateOne(ctx, bson.M{
			"id":         actionData.ID,
			"is_deleted": false,
		}, update)
		if err != nil {
			return nil, w.handleError("update wallet", err, error_codes.UnhandledServerError)
		}
		cpsAction.CurrentAction = wallet
		return &cpsAction, nil

	case string(model.ActionDelete):
		objectID, err := w.parseObjectID(actionData.ID)
		if err != nil {
			return nil, err
		}
		wallet, err := w.walletDal.UpdateOne(ctx, bson.M{
			"_id":        objectID,
			"is_deleted": false,
		}, bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		})
		if err != nil {
			return nil, w.handleError("delete wallet", err, error_codes.UnhandledServerError)
		}
		cpsAction.CurrentAction = wallet
		return &cpsAction, nil
	}

	return &cpsAction, nil
}

func (w *Wallet) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"action_status":       model.ActionRejected,
		"rejected_reason":     req.RejectedReason,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := w.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, w.handleError("find cps action", err, "WALLET_NOT_FOUND")
		}
		return nil, w.handleError("update cps action", err, error_codes.UnhandledServerError)
	}
	return &cpsAction, nil
}
