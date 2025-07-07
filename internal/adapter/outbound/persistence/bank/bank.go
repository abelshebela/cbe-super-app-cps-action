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
	bankDal dal.MongoDal[entity.Bank, entity.Bank]
	cpsDal  dal.MongoDal[model.CPSAction, model.CPSAction]
	logger  utils.Logger
}

func InitBank(client *mongo.Client, database string, collections []string, logger utils.Logger) outbound.BankPersistence {
	bankDal := dal.NewMongoDal[entity.Bank, entity.Bank](client, database, collections[0])
	cpsDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collections[1])
	return &Bank{
		bankDal: bankDal,
		cpsDal:  cpsDal,
		logger:  logger,
	}
}

func (b *Bank) CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
		"request_action":          cpsReq.RequestAction,
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
	cpsAction, err := b.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID().Hex(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.FullName,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		RequestAction:    string(model.RequestCreateBank),
		ActionType:       string(model.ActionCreate),
		CurrentAction:    cpsReq.ActionData,
		MakerActionTime:  time.Now(),
	})

	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &cpsAction, nil
}

func (b *Bank) DeleteBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {

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
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	cpsRes, err := b.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID().Hex(),
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

	banks, err := b.bankDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("no bank data found", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to get bank data", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
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
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}
	projection := bson.M{}

	bank, err := b.bankDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to get bank", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}
	return bank, nil
}

func (b *Bank) UpdateBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {

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
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	cps, err := b.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID().Hex(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.FullName,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		CurrentAction:    cpsReq.ActionData,
		RequestAction:    string(model.RequestUpdateBank),
		PreviosAction: map[string]any{
			"name": bank.Name,
			"code": bank.Code,
			"bic":  bank.BIC,
		},
		MakerActionTime: time.Now(),
	})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}

		b.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &cps, nil
}

func (b *Bank) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error) {
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
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	var actionData entity.Bank
	data, err := bson.Marshal(cpsAction.CurrentAction)
	if err != nil {
		b.logger.Errorf("failed to marshal bson: %v", err)
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	if err := bson.Unmarshal([]byte(data), &actionData); err != nil {
		b.logger.Errorf("failed to unmarshal into Bank: %v", err)
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	if cpsAction.ActionType == string(model.ActionCreate) {
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
			return nil, fmt.Errorf(error_codes.UnhandledServerError)
		}

		cpsAction.CurrentAction = bank

		return &cpsAction, nil

	}

	if cpsAction.ActionType == string(model.ActionUpdate) {
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

		if cpsAction.RequestAction == string(model.RequestEnableBank) {
			update["enabled"] = true
		}

		if cpsAction.RequestAction == string(model.RequestDisableBank) {
			update["enabled"] = false
		}

		update["last_modified_at"] = time.Now()

		bank, err = b.bankDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to update bank", err)
			return nil, fmt.Errorf(error_codes.UnhandledServerError)
		}

		cpsAction.CurrentAction = bank

		return &cpsAction, nil
	}

	if cpsAction.ActionType == string(model.ActionDelete) {
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
			return nil, fmt.Errorf(error_codes.UnhandledServerError)
		}
		cpsAction.CurrentAction = bank

		return &cpsAction, nil
	}

	return &cpsAction, nil
}

func (b *Bank) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
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
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}
	return &cpsAction, nil
}

func (b *Bank) EnableOrDisableBank(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {

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
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to get bank", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	cps, err := b.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID().Hex(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.FullName,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		CurrentAction:    cpsReq.ActionData,
		RequestAction:    string(requestAction),
		PreviosAction: map[string]any{
			"name":    bank.Name,
			"code":    bank.Code,
			"bic":     bank.BIC,
			"enabled": bank.Enabled,
		},
		MakerActionTime: time.Now(),
	})
	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &cps, nil
}
