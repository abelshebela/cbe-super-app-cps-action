package faydaaccount

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/fayda_account/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/outbound"
	"fmt"
	"net/http"
	"time"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FaydaAccountRepo struct {
	client      *mongo.Client
	timeout     time.Duration
	logger      utils.Logger
	cpsDal      dal.MongoDal[entity.CPSAction, entity.CPSAction]
	customerDal dal.MongoDal[member.User, member.User]
}

var _ outbound.FaydaAccountRepository = (*FaydaAccountRepo)(nil)

func InitFaydaAccountPersistence(client *mongo.Client, database string,
	timeout time.Duration, logger utils.Logger) *FaydaAccountRepo {
	cpsDal := dal.NewMongoDal[entity.CPSAction, entity.CPSAction](client, database, "cps_actions")
	customerDal := dal.NewMongoDal[member.User, member.User](client, database, "customers")
	return &FaydaAccountRepo{
		client:      client,
		timeout:     timeout,
		logger:      logger,
		cpsDal:      cpsDal,
		customerDal: customerDal,
	}
}

func (f *FaydaAccountRepo) InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	filter := bson.M{
		"user_code":  req.MakerUser.UserCode,
		"status":     entity.ActionPending,
		"department": req.Department,
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

	req.Status = entity.ActionPending
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
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	filter := bson.M{
		"department": req.Department,
		"status":     entity.ActionPending,
	}

	update := bson.M{
		"$set": bson.M{
			"checker_user": bson.M{
				"full_name":    req.CheckerUser.FullName,
				"phone_number": req.CheckerUser.PhoneNumber,
				"user_code":    req.CheckerUser.UserCode,
			},
			"status":                 entity.ActionApproved,
			"checker_action_time":    time.Now(),
		},
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
		"user_code": req.ActionData.UseCode,
		"kyc": bson.M{
			"level": 1,
		},
	}
	customerUpdate := bson.M{
		"account_status": "DISABLED",
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
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	filter := bson.M{
		"department": req.Department,
		"status":     entity.ActionPending,
	}

	update := bson.M{
		"$set": bson.M{
			"checker_user": bson.M{
				"full_name":    req.CheckerUser.FullName,
				"phone_number": req.CheckerUser.PhoneNumber,
				"user_code":    req.CheckerUser.UserCode,
			},
			"status":                 entity.ActionRejected,
			"rejected_action_reason": req.RejectedReason,
			"checker_action_time":    time.Now(),
		},
	}
	cpsAction, err := f.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		f.logger.Errorf("failed to update customer status", err)
		err = fmt.Errorf("failed to update customer status %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return &cpsAction, nil
}
