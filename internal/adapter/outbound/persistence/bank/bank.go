package bank

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/entity"
	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/bank"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Bank struct {
	bankDal dal.MongoDal[entity.Bank, entity.Bank]
	cpsDal  dal.MongoDal[model.CpsAction, model.CpsAction]
	logger  utils.Logger
}

func InitBank(client *mongo.Client, database string, collections []string, logger utils.Logger) outbound.BankPersistence {
	bankDal := dal.NewMongoDal[entity.Bank, entity.Bank](client, database, collections[0])
	cpsDal := dal.NewMongoDal[model.CpsAction, model.CpsAction](client, database, collections[1])
	return &Bank{
		bankDal: bankDal,
		cpsDal:  cpsDal,
		logger:  logger,
	}
}

func (b *Bank) CreateBank(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingBank, err := b.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingBank != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		b.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	cpsAction, err := b.cpsDal.InsertOne(ctx, model.CpsAction{
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
		b.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cpsAction, nil
}

func (b *Bank) DeleteBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingBank, err := b.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingBank != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		b.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	bankFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	bankProjection := bson.M{
		"name": 1,
		"code": 1,
		"bic":  1,
		"logo": 1,
	}

	bank, err := b.bankDal.FindOne(ctx, bankFilter, bankProjection)
	if err != nil {
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cpsRes, err := b.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		RequestAction: model.RequestDeleteBank,
		ActionType:    model.ActionDelete,
		ActionData:    cpsReq.ActionData,
		PreviousData: map[string]any{
			"name":       bank.Name,
			"code":       bank.Code,
			"bic":        bank.BIC,
			"logo":       bank.Logo,
			"is_deleted": bank.IsDeleted,
		},
		CurrentData: map[string]any{
			"is_deleted": true,
		},
		MakerActionTime: time.Now(),
	})

	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cpsRes, nil
}

func (b *Bank) GetAllBanks(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	banks, err := b.bankDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("no bank data found", err)
			err = fmt.Errorf("banks not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "banks data not found",
			})
			return nil, err
		}
		b.logger.Errorf("failed to get bank data", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	total, err := b.bankDal.TotalCount(ctx, bson.M{})
	if err != nil {
		b.logger.Errorf("failed to get bank total counts", err)
		err := fmt.Errorf("failed to get bank total counts %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &entity.BankResponse{
		Page:  filterParams.Page,
		Banks: banks,
		Limit: constant.DefaultPerPage,
		Total: total,
	}, nil
}

func (b *Bank) GetBank(ctx context.Context, id string) (*entity.Bank, error) {
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}
	projection := bson.M{}

	bank, err := b.bankDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found", err)
			err = fmt.Errorf("bank not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "bank not found",
			})
			return nil, err
		}
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return bank, nil
}

func (b *Bank) UpdateBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingAd, err := b.bankDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		b.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	bankFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	bankProjection := bson.M{
		"name": 1,
		"code": 1,
		"bic":  1,
	}

	bank, err := b.bankDal.FindOne(ctx, bankFilter, bankProjection)
	if err != nil {
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cps, err := b.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		ActionType:    model.ActionUpdate,
		ActionData:    cpsReq.ActionData,
		RequestAction: model.RequestUpdateBank,
		PreviousData: map[string]any{
			"name": bank.Name,
			"code": bank.Code,
			"bic":  bank.BIC,
		},
		CurrentData:     cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (b *Bank) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error) {
	var bank entity.Bank
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

	cpsAction, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to update cps action", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	actionData, ok := cpsAction.ActionData.(entity.Bank)
	if !ok {
		b.logger.Errorf("failed to cast action data to bank entity")
		return nil, fmt.Errorf("failed to cast: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if cpsAction.ActionType == model.ActionCreate {
		req := entity.Bank{
			ID:        bson.NewObjectID().Hex(),
			Name:      actionData.Name,
			Logo:      actionData.Logo,
			Code:      actionData.Code,
			BIC:       actionData.BIC,
			CreatedAt: time.Now(),
		}

		bank, err = b.bankDal.InsertOne(ctx, req)
		if err != nil {
			b.logger.Errorf("failed to create bank", err)
			err = fmt.Errorf("failed to create bank %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}

		cpsAction.ActionData = bank

		return &cpsAction, nil

	}

	if cpsAction.ActionType == model.ActionUpdate {
		filter := bson.M{"id": actionData.ID}
		update := bson.M{}

		if actionData.Name != "" {
			update["name"] = actionData.Name
		}

		if actionData.Code != "" {
			update["code"] = actionData.Code
		}

		if actionData.BIC != "" {
			update["bic"] = actionData.BIC
		}

		update["last_modified_at"] = time.Now()

		bank, err = b.bankDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to update bank", err)
			err = fmt.Errorf("failed to update bank %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}

		cpsAction.ActionData = bank

		return &cpsAction, nil
	}

	if cpsAction.ActionType == model.ActionDelete {
		filter := bson.M{
			"id": actionData.ID,
		}

		update := bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}

		bank, err = b.bankDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to update bank", err)
			err = fmt.Errorf("failed to update bank %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}
		cpsAction.ActionData = bank

		return &cpsAction, nil
	}

	return &cpsAction, nil
}

func (b *Bank) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error) {
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

	cpsAction, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to update bank status", err)
		err = fmt.Errorf("failed to update bank status %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return &cpsAction, nil
}

func (b *Bank) EnableOrDisableWallet(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingAd, err := b.bankDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("failed to get bank", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		b.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return nil, err
	}

	bankFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	bankProjection := bson.M{
		"name":    1,
		"code":    1,
		"bic":     1,
		"enabled": 1,
	}

	bank, err := b.bankDal.FindOne(ctx, bankFilter, bankProjection)
	if err != nil {
		b.logger.Errorf("failed to get wallet", err)
		err = fmt.Errorf("failed to get wallet %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cps, err := b.cpsDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		ActionType:    model.ActionUpdate,
		ActionData:    cpsReq.ActionData,
		RequestAction: requestAction,
		PreviousData: map[string]any{
			"name":    bank.Name,
			"code":    bank.Code,
			"bic":     bank.BIC,
			"enabled": bank.Enabled,
		},
		CurrentData:     cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}
