// Package bank provides persistence and CPS action management for bank entities.
package bank

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bank"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Bank struct {
	bankDal dal.MongoDal[entity.BankDocument, entity.BankDocument]
	cpsDal  dal.MongoDal[model.CPSAction, model.CPSAction]
	logger  utils.Logger
}

func InitBank(client *mongo.Client, database string, collections []string, logger utils.Logger) outbound.BankPersistence {
	bankDal := dal.NewMongoDal[entity.BankDocument, entity.BankDocument](client, database, collections[0])
	cpsDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collections[1])
	return &Bank{
		bankDal: bankDal,
		cpsDal:  cpsDal,
		logger:  logger,
	}
}

func (b *Bank) CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error {
	filter := bson.M{
		"maker_phone_number": cpsReq.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsReq.Department,
		"request_action":     cpsReq.RequestAction,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingBank, err := b.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("failed to get bank", err)
		return fmt.Errorf(error_codes.UnhandledServerError)
	} else if existingBank != nil {
		b.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return fmt.Errorf(error_codes.PendingRequestExists)
	}

	return nil
}

func (b *Bank) CreateBank(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cps := b.buildCPSAction(cpsReq, model.RequestCreateBank, nil)
	cps.ActionType = string(model.ActionCreate)

	result, err := b.cpsDal.InsertOne(ctx, cps)
	if err != nil {
		b.logger.Errorf("failed to create cps action: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &result, nil
}

func (b *Bank) DeleteBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	objectID, err := b.parseObjectID(id)
	if err != nil {
		return nil, err
	}
	bankFilter := bson.M{
		"_id":        objectID,
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
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	cpsRes, err := b.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.FullName,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		RequestAction:    string(model.RequestDeleteBank),
		ActionType:       string(model.ActionDelete),
		CurrentAction:    cpsReq.ActionData,
		PreviosAction: map[string]any{
			"name":       bank.Name,
			"code":       bank.Code,
			"bic":        bank.BIC,
			"logo":       bank.Logo,
			"is_deleted": bank.IsDeleted,
		},
		MakerActionTime: time.Now(),
	})

	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
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

	banksDocs, err := b.bankDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		b.logger.Errorf("failed to get bank data", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	var banks []*entity.Bank
	for _, doc := range banksDocs {
		banks = append(banks, b.toDomain(doc))
	}

	total, err := b.bankDal.TotalCount(ctx, bson.M{})
	if err != nil {
		b.logger.Errorf("failed to get bank total counts", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &entity.BankResponse{
		Page:  filterParams.Page,
		Banks: banks,
		Limit: constant.DefaultPerPage,
		Total: total,
	}, nil
}

func (b *Bank) GetBank(ctx context.Context, id string) (*entity.Bank, error) {
	objectID, err := b.parseObjectID(id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%T\n", objectID)
	filter := bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}
	projection := bson.M{}

	bank, err := b.bankDal.FindOne(ctx, filter, projection)
	fmt.Println(err, "err")
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to get bank", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}
	return b.toDomain(bank), nil
}

func (b *Bank) UpdateBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	projection := bson.M{"name": 1, "code": 1, "bic": 1}

	bank, err := b.findBankByID(ctx, id, projection)
	if err != nil {
		return nil, err
	}

	prev := map[string]any{
		"name": bank.Name,
		"code": bank.Code,
		"bic":  bank.BIC,
	}

	cps := b.buildCPSAction(cpsReq, model.RequestUpdateBank, prev)

	result, err := b.cpsDal.InsertOne(ctx, cps)
	if err != nil {
		b.logger.Errorf("failed to create cps action: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &result, nil
}

func (b *Bank) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error) {
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}
	update := bson.M{
		"checker_id":           req.CheckerUser.UserCode,
		"checker_name":         req.CheckerUser.FullName,
		"checker_phone_number": req.CheckerUser.PhoneNumber,
		"action_status":        model.ActionApproved,
		"checker_action_time":  time.Now(),
	}

	cpsAction, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(error_codes.ActionNotFound)
		}
		b.logger.Errorf("failed to update cps action: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	var actionData entity.Bank
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		b.logger.Errorf("failed to unmarshal action data: %v", err)
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	// handle based on action
	switch cpsAction.RequestAction {
	case string(model.RequestCreateBank):
		doc, err := b.toDocument(&actionData)
		if err != nil {
			return nil, fmt.Errorf(error_codes.InvalidActionData)
		}
		bank, err := b.bankDal.InsertOne(ctx, *doc)
		if err != nil {
			return nil, fmt.Errorf(error_codes.UnhandledServerError)
		}
		cpsAction.CurrentAction = bank

	case string(model.RequestUpdateBank):
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
		bank, err := b.updateBankByID(ctx, actionData.ID, update)
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank

	case string(model.RequestDeleteBank):
		update := bson.M{"is_deleted": true, "deleted_at": time.Now()}
		bank, err := b.updateBankByID(ctx, actionData.ID, update)
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank

	case string(model.RequestEnableBank):
		bank, err := b.updateBankByID(ctx, actionData.ID, bson.M{"enabled": true})
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank

	case string(model.RequestDisableBank):
		bank, err := b.updateBankByID(ctx, actionData.ID, bson.M{"enabled": false})
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank
	}

	return &cpsAction, nil
}

func (b *Bank) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_id":           req.CheckerUser.UserCode,
		"checker_name":         req.CheckerUser.FullName,
		"checker_phone_number": req.CheckerUser.PhoneNumber,

		"action_status":       model.ActionRejected,
		"rejected_reason":     req.RejectedReason,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to update bank status", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}
	return &cpsAction, nil
}

func (b *Bank) EnableOrDisableBank(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	projection := bson.M{"name": 1, "code": 1, "bic": 1, "enabled": 1}

	bank, err := b.findBankByID(ctx, id, projection)
	if err != nil {
		return nil, err
	}

	prev := map[string]any{
		"name":    bank.Name,
		"code":    bank.Code,
		"bic":     bank.BIC,
		"enabled": bank.Enabled,
	}

	cps := b.buildCPSAction(cpsReq, requestAction, prev)

	result, err := b.cpsDal.InsertOne(ctx, cps)
	if err != nil {
		b.logger.Errorf("failed to create cps action: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &result, nil
}

func (b *Bank) CheckExistingBank(ctx context.Context, name string) (bool, error) {
	bank, err := b.bankDal.FindOne(ctx, bson.M{
		"name":       name,
		"is_deleted": false,
	}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank name not taken", err)
			return false, nil
		}
		b.logger.Errorf("failed to get bank", err)
		return false, fmt.Errorf(error_codes.UnhandledServerError)
	}
	if bank != nil {
		return true, nil
	}
	return false, nil
}

func (b *Bank) UpdateLogo(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	projection := bson.M{"name": 1, "code": 1, "bic": 1, "logo": 1}

	bank, err := b.findBankByID(ctx, id, projection)
	if err != nil {
		return nil, err
	}

	prev := map[string]any{
		"name": bank.Name,
		"code": bank.Code,
		"bic":  bank.BIC,
		"logo": bank.Logo,
	}

	cps := b.buildCPSAction(cpsReq, model.RequestUpdateBank, prev)

	result, err := b.cpsDal.InsertOne(ctx, cps)
	if err != nil {
		b.logger.Errorf("failed to create cps action: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &result, nil
}
