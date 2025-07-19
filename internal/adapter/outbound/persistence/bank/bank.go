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
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func shortActionCode(code string) string {
	if len(code) <= 6 {
		return code
	}
	return code[:3] + "***" + code[len(code)-3:]
}

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
		b.logger.Errorf("failed to get bank: %v", err)
		return fmt.Errorf(error_codes.UnhandledServerError)
	} else if existingBank != nil {
		b.logger.Infof("pending cps action present for user_code: %s, department: %s", cpsReq.MakerUser.UserCode, cpsReq.Department)
		return fmt.Errorf(error_codes.PendingRequestExists)
	}

	b.logger.Infof("no pending cps action found for user_code: %s, department: %s", cpsReq.MakerUser.UserCode, cpsReq.Department)
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

	b.logger.Infof("cps action created for bank creation, action_code: %s, maker_code: %s", shortActionCode(result.ActionCode), cpsReq.MakerUser.UserCode)
	return &result, nil
}

func (b *Bank) DeleteBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	objectID, err := b.parseObjectID(id)
	if err != nil {
		b.logger.Errorf("invalid bank id provided: %s, error: %v", id, err)
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
		b.logger.Errorf("failed to get bank for deletion, id: %s, error: %v", id, err)
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
		b.logger.Errorf("failed to create cps action for bank deletion, bank_id: %s, error: %v", id, err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("cps action created for bank deletion, action_code: %s, bank_id: %s, maker_code: %s", shortActionCode(cpsRes.ActionCode), id, cpsReq.MakerUser.UserCode)
	return &cpsRes, nil
}

func (b *Bank) GetAllBanks(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Bank], error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["status"] = filterParams.Filters
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	banksDocs, err := b.bankDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		b.logger.Errorf("failed to get bank data, page: %d, per_page: %d, error: %v",
			filterParams.Page, filterParams.PerPage, err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	var banks []*entity.Bank
	for _, doc := range banksDocs {
		banks = append(banks, b.toDomain(doc))
	}

	total, err := b.bankDal.TotalCount(ctx, bson.M{})
	if err != nil {
		b.logger.Errorf("failed to get bank total counts, error: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("retrieved banks, page: %d, count: %d, total: %d", filterParams.Page, len(banks), total)
	meta := common_util.BuildPaginationMeta(total, page, limit)

	return &common_util.PaginatedResponse[[]*entity.Bank]{
		Data: banks,
		Meta: meta,
	}, nil
}

func (b *Bank) GetBank(ctx context.Context, id string) (*entity.Bank, error) {
	objectID, err := b.parseObjectID(id)
	if err != nil {
		b.logger.Errorf("invalid bank id provided: %s, error: %v", id, err)
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
			b.logger.Errorf("bank not found, id: %s", id)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to get bank, id: %s, error: %v", id, err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("bank retrieved, id: %s", id)
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
		b.logger.Errorf("failed to create cps action for bank update, bank_id: %s, error: %v", id, err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("cps action created for bank update, action_code: %s, bank_id: %s, maker_code: %s", shortActionCode(result.ActionCode), id, cpsReq.MakerUser.UserCode)
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
			b.logger.Errorf("cps action not found for authorization, action_code: %s", shortActionCode(req.ActionCode))
			return nil, fmt.Errorf(error_codes.ActionNotFound)
		}
		b.logger.Errorf("failed to update cps action for authorization, action_code: %s, error: %v", shortActionCode(req.ActionCode), err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	var actionData entity.BankDocument
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		b.logger.Errorf("failed to unmarshal action data for authorization, action_code: %s", shortActionCode(req.ActionCode))
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	// handle based on action
	switch cpsAction.RequestAction {
	case string(model.RequestCreateBank):
		doc, err := b.toDocument(&actionData)
		if err != nil {
			b.logger.Errorf("failed to convert action data to document for bank creation, action_code: %s", shortActionCode(req.ActionCode))
			return nil, fmt.Errorf(error_codes.InvalidActionData)
		}
		bank, err := b.bankDal.InsertOne(ctx, *doc)
		if err != nil {
			b.logger.Errorf("failed to insert bank after authorization, action_code: %s", shortActionCode(req.ActionCode))
			return nil, fmt.Errorf(error_codes.UnhandledServerError)
		}
		cpsAction.CurrentAction = bank
		b.logger.Infof("bank created after authorization, action_code: %s, checker_code: %s", shortActionCode(req.ActionCode), req.CheckerUser.UserCode)

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
		bank, err := b.updateBankByID(ctx, actionData.ID.Hex(), update)
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank
		b.logger.Infof("bank updated after authorization, action_code: %s, checker_code: %s", shortActionCode(req.ActionCode), req.CheckerUser.UserCode)

	case string(model.RequestDeleteBank):
		update := bson.M{"is_deleted": true, "deleted_at": time.Now()}
		bank, err := b.updateBankByID(ctx, actionData.ID.Hex(), update)
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank
		b.logger.Infof("bank deleted after authorization, action_code: %s, checker_code: %s", shortActionCode(req.ActionCode), req.CheckerUser.UserCode)

	case string(model.RequestEnableBank):
		bank, err := b.updateBankByID(ctx, actionData.ID.Hex(), bson.M{"enabled": true})
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank
		b.logger.Infof("bank enabled after authorization, action_code: %s, checker_code: %s", shortActionCode(req.ActionCode), req.CheckerUser.UserCode)

	case string(model.RequestDisableBank):
		bank, err := b.updateBankByID(ctx, actionData.ID.Hex(), bson.M{"enabled": false})
		if err != nil {
			return nil, err
		}
		cpsAction.CurrentAction = bank
		b.logger.Infof("bank disabled after authorization, action_code: %s, checker_code: %s", shortActionCode(req.ActionCode), req.CheckerUser.UserCode)
	}

	return &cpsAction, nil
}

func (b *Bank) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
	checkerUser := contexts.ExtractContext(ctx)
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_id":           checkerUser.UserCode,
		"checker_name":         checkerUser.FullName,
		"checker_phone_number": checkerUser.PhoneNumber,

		"action_status":       model.ActionRejected,
		"rejected_reason":     req.RejectedReason,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("cps action not found for rejection, action_code: %s", shortActionCode(req.ActionCode))
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to update cps action for rejection, action_code: %s, error: %v", shortActionCode(req.ActionCode), err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("cps action rejected, action_code: %s, checker_code: %s", shortActionCode(req.ActionCode), checkerUser.UserCode)
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
		b.logger.Errorf("failed to create cps action for bank enable/disable, bank_id: %s, action: %s, error: %v",
			id, requestAction, err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("cps action created for bank enable/disable, action_code: %s, bank_id: %s, action: %s, maker_code: %s", shortActionCode(result.ActionCode), id, requestAction, cpsReq.MakerUser.UserCode)
	return &result, nil
}

func (b *Bank) CheckExistingBank(ctx context.Context, name string) (bool, error) {
	bank, err := b.bankDal.FindOne(ctx, bson.M{
		"name":       name,
		"is_deleted": false,
	}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Infof("bank name available, name: [REDACTED]")
			return false, nil
		}
		b.logger.Errorf("failed to check existing bank, name: [REDACTED], error: %v", err)
		return false, fmt.Errorf(error_codes.UnhandledServerError)
	}
	if bank != nil {
		b.logger.Infof("bank name already exists, name: [REDACTED]")
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
		b.logger.Errorf("failed to create cps action for logo update, bank_id: %s, error: %v", id, err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	b.logger.Infof("cps action created for logo update, action_code: %s, bank_id: %s, maker_code: %s", shortActionCode(result.ActionCode), id, cpsReq.MakerUser.UserCode)
	return &result, nil
}
