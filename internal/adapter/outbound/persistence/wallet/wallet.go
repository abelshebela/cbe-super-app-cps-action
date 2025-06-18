package wallet

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/wallet/entity"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/wallet"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type Wallet struct {
	walletDal dal.MongoDal[entity.Wallet, entity.Wallet]
	cpsDal    dal.MongoDal[model.CpsAction, model.CpsAction]
	logger    utils.Logger
}

func InitWallet(client *mongo.Client, database string, collections []string, logger utils.Logger) outbound.WalletPersistence {
	walletDal := dal.NewMongoDal[entity.Wallet, entity.Wallet](client, database, collections[0])
	cpsDal := dal.NewMongoDal[model.CpsAction, model.CpsAction](client, database, collections[1])
	return &Wallet{
		walletDal: walletDal,
		cpsDal:    cpsDal,
		logger:    logger,
	}
}

func (w *Wallet) CreateWallet(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingBank, err := w.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingBank != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		w.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	cpsAction, err := w.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:              bson.NewObjectID().Hex(),
		ActionCode:      utils.RandomGenerator(20),
		MakerUser:       cpsReq.MakerUser,
		Department:      cpsReq.Department,
		Status:          model.ActionPending,
		RequestAction:   model.RequestCreateBank,
		ActionType:      model.ActionCreate,
		ActionData:      cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})

	if err != nil {
		w.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cpsAction, nil
}

func (w *Wallet) UpdateWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingAd, err := w.walletDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		w.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	walletFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	walletProjection := bson.M{
		"name": 1,
		"code": 1,
		"bic":  1,
	}

	wallet, err := w.walletDal.FindOne(ctx, walletFilter, walletProjection)
	if err != nil {
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cps, err := w.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		ActionType:    model.ActionUpdate,
		ActionData:    cpsReq.ActionData,
		RequestAction: model.RequestUpdateBank,
		PreviousData: map[string]any{
			"name": wallet.Name,
			"code": wallet.Code,
		},
		CurrentData:     cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		w.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (w *Wallet) DeleteWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingBank, err := w.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingBank != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		w.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	walletFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	walletProjection := bson.M{
		"name": 1,
		"code": 1,
		"bic":  1,
		"logo": 1,
	}

	wallet, err := w.walletDal.FindOne(ctx, walletFilter, walletProjection)
	if err != nil {
		w.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cpsRes, err := w.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		RequestAction: model.RequestDeleteBank,
		ActionType:    model.ActionDelete,
		ActionData:    cpsReq.ActionData,
		PreviousData: map[string]any{
			"name":       wallet.Name,
			"code":       wallet.Code,
			"avatar":     wallet.Avatar,
			"is_deleted": wallet.IsDeleted,
		},
		CurrentData: map[string]any{
			"is_deleted": true,
		},
		MakerActionTime: time.Now(),
	})

	if err != nil {
		w.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cpsRes, nil
}

func (w *Wallet) EnableOrDisableWallet(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingAd, err := w.walletDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		w.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	walletFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	walletProjection := bson.M{
		"name": 1,
		"code": 1,
		"enabled":1,
	}

	wallet, err := w.walletDal.FindOne(ctx, walletFilter, walletProjection)
	if err != nil {
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cps, err := w.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		ActionType:    model.ActionUpdate,
		ActionData:    cpsReq.ActionData,
		RequestAction: requestAction,
		PreviousData: map[string]any{
			"name":    wallet.Name,
			"code":    wallet.Code,
			"enabled": wallet.Enabled,
		},
		CurrentData:     cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		w.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (w *Wallet) GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	wallets, err := w.walletDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Errorf("no wallet data found", err)
			err = fmt.Errorf("wallets not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "wallets data not found",
			})
			return nil, err
		}
		w.logger.Errorf("failed to get bank data", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	total, err := w.walletDal.TotalCount(ctx, bson.M{})
	if err != nil {
		w.logger.Errorf("failed to get bank total counts", err)
		err := fmt.Errorf("failed to get bank total counts %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &entity.WalletResponse{
		Page:    filterParams.Page,
		Wallets: wallets,
		Limit:   constant.DefaultPerPage,
		Total:   total,
	}, nil
}

func (w *Wallet) GetWallet(ctx context.Context, id string) (*entity.Wallet, error) {
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}
	projection := bson.M{}

	wallet, err := w.walletDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Errorf("wallet not found", err)
			err = fmt.Errorf("wallet not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "wallet not found",
			})
			return nil, err
		}
		w.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return wallet, nil
}

func (w *Wallet) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error) {
	var wallet entity.Wallet
	var err error

	filter := bson.M{
		"action_code": req.ActionCode,
		"department":  req.Department,
		"status":      model.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"status":              model.ActionApproved,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := w.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		w.logger.Errorf("failed to update cps action", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	actionData, ok := cpsAction.ActionData.(entity.Wallet)
	if !ok {
		w.logger.Errorf("failed to cast action data to wallet entity")
		return nil, fmt.Errorf("failed to cast: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if cpsAction.ActionType == model.ActionCreate {
		req := entity.Wallet{
			ID:        bson.NewObjectID().Hex(),
			Name:      actionData.Name,
			Avatar:    actionData.Avatar,
			Code:      actionData.Code,
			CreatedAt: time.Now(),
		}

		wallet, err = w.walletDal.InsertOne(ctx, req)
		if err != nil {
			w.logger.Errorf("failed to create wallet", err)
			err = fmt.Errorf("failed to create wallet %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}

		cpsAction.ActionData = wallet

		return &cpsAction, nil

	}

	if cpsAction.ActionType == model.ActionUpdate {
		filter := bson.M{
			"id":         actionData.ID,
			"is_deleted": false,
		}
		update := bson.M{}

		if actionData.Name != "" {
			update["name"] = actionData.Name
		}

		if actionData.Code != "" {
			update["code"] = actionData.Code
		}

		if cpsAction.RequestAction == model.RequestEnableWallet {
			update["enabled"] = true
		}

		if cpsAction.RequestAction == model.RequestDisableWallet {
			update["enabled"] = false
		}

		update["last_modified_at"] = time.Now()

		wallet, err = w.walletDal.UpdateOne(ctx, filter, update)
		if err != nil {
			w.logger.Errorf("failed to update wallet", err)
			err = fmt.Errorf("failed to update wallet %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}

		cpsAction.ActionData = wallet

		return &cpsAction, nil
	}

	if cpsAction.ActionType == model.ActionDelete {
		filter := bson.M{
			"id":         actionData.ID,
			"is_deleted": false,
		}

		update := bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}

		wallet, err = w.walletDal.UpdateOne(ctx, filter, update)
		if err != nil {
			w.logger.Errorf("failed to update wallet", err)
			err = fmt.Errorf("failed to update wallet %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}
		cpsAction.ActionData = wallet

		return &cpsAction, nil
	}

	return &cpsAction, nil
}

func (w *Wallet) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"action_code": req.ActionCode,
		"department":  req.Department,
		"status":      model.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"status":              model.ActionRejected,
		"rejected_reason":     req.RejectedReason,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := w.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		w.logger.Errorf("failed to update wallet status", err)
		err = fmt.Errorf("failed to update wallet status %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return &cpsAction, nil
}
