package faydaaccount

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"cbe-super-app-cps-action/internal/port/outbound"

	constant "cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FaydaAccountRepo struct {
	client      *mongo.Client
	logger      utils.Logger
	cpsDal      dal.MongoDal[entity.CPSAction, entity.CPSAction]
	customerDal dal.MongoDal[member.User, member.User]
}

var _ outbound.FaydaAccountRepository = (*FaydaAccountRepo)(nil)

func InitFaydaAccountPersistence(client *mongo.Client, database string,
	cpsCollection []string, logger utils.Logger) *FaydaAccountRepo {
	cpsDal := dal.NewMongoDal[entity.CPSAction, entity.CPSAction](client, database, cpsCollection[0])
	customerDal := dal.NewMongoDal[member.User, member.User](client, database, cpsCollection[1])
	return &FaydaAccountRepo{
		client:      client,
		logger:      logger,
		cpsDal:      cpsDal,
		customerDal: customerDal,
	}
}

func (f *FaydaAccountRepo) InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	filter := bson.M{
		"maker_user.phone_number": req.MakerUser.PhoneNumber,
		"status":                  entity.ActionPending,
		"department":              req.Department,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	faydaAccount, err := f.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		f.logger.Errorf("failed to get fayda account", err)
		err = fmt.Errorf("failed to get fayda account %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	if faydaAccount != nil {
		f.logger.Infof("pending cps action present", req.MakerUser.FullName, req.MakerUser.UserCode, req.Department)
		err = fmt.Errorf("failed to get fayda account %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "pending cps action present",
		})
		return nil, err
	}

	customerFilter := bson.M{
		"phone_number": req.ActionData.PhoneNumber,
	}

	customerProjection := bson.M{
		"is_account_blocked": 1,
		"user_code":          1,
		"full_name":          1,
		"phone_number":       1,
	}

	customer, err := f.customerDal.FindOne(ctx, customerFilter, customerProjection)
	if err != nil {
		f.logger.Errorf("failed to get customer account", err)
		err = fmt.Errorf("failed to get customer account %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	req.Status = entity.ActionPending
	req.RequestAction = entity.RequestDisableFaydaAccount
	req.ActionType = entity.ActionCreate
	req.PreviousData = map[string]any{
		"user_code":          customer.UserCode,
		"full_name":          customer.FullName,
		"phone_number":       customer.PhoneNumber,
		"is_account_blocked": customer.IsAccountBlocked,
	}
	req.CurrentData = map[string]any{
		"is_account_blocked": true,
	}

	req.MakerActionTime = time.Now()
	cpsAction, err := f.cpsDal.InsertOne(ctx, req)
	if err != nil {
		f.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cpsAction, nil
}

func (f *FaydaAccountRepo) AuthorizeFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	filter := bson.M{
		"action_data.phone_number": req.ActionData.PhoneNumber,
		"department":               req.Department,
		"status":                   entity.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"status":              entity.ActionApproved,
		"checker_action_time": time.Now(),
	}
	cpsAction, err := f.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		f.logger.Errorf("failed to update cps action", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	customerFilter := bson.M{
		"phone_number": req.ActionData.PhoneNumber,
		"kyc.level":    1,
	}
	customerUpdate := bson.M{
		"is_account_blocked": true,
	}

	_, err = f.customerDal.UpdateOne(ctx, customerFilter, customerUpdate)
	if err != nil {
		f.logger.Errorf("failed to update  action", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cpsAction, nil
}

func (f *FaydaAccountRepo) RejectFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	filter := bson.M{
		"action_data.phone_number": req.ActionData.PhoneNumber,
		"department":               req.Department,
		"status":                   entity.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"status":                 entity.ActionRejected,
		"rejected_action_reason": req.RejectedReason,
		"checker_action_time":    time.Now(),
	}

	cpsAction, err := f.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		f.logger.Errorf("failed to update fayda customer status", err)
		err = fmt.Errorf("failed to update fayda customer status %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return &cpsAction, nil
}
