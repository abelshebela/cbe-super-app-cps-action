package persistence

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/lib"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/rs/zerolog/log"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AmountBasedAuthRepo struct {
	client       *mongo.Client
	authTier     dal.MongoDal[amount_based_auth_domain.AuthTier, amount_based_auth_domain.AuthTier]
	cpsActionDal dal.MongoDal[model.CPSAction, model.CPSAction]
	timeOut      time.Duration
	logger       utils.Logger
}

func InitAmountBasedAuth(client *mongo.Client, database string, collection []string, logger utils.Logger) *AmountBasedAuthRepo {
	authTier := dal.NewMongoDal[amount_based_auth_domain.AuthTier, amount_based_auth_domain.AuthTier](client, database, collection[0])
	cpsActionDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collection[1])

	return &AmountBasedAuthRepo{
		client:       client,
		authTier:     authTier,
		cpsActionDal: cpsActionDal,
		timeOut:      5 * time.Second,
		logger:       logger,
	}
}

func (a *AmountBasedAuthRepo) GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*amount_based_auth_domain.AmountBasedAuthRespose, error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{}

	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"min_amount": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"max_amount": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"enabled": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"is_deleted": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	if filterParams.Filters != "" {
		filter["account_status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	amountBased, err := a.authTier.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("no auth tier data found", err)
			return nil, fmt.Errorf("NO_CUSTOMER_DATA_FOUND")
		}
		a.logger.Errorf("failed to get auth tier data", err)
		return nil, fmt.Errorf("FAILED_TO_GET_CUSTOMER")
	}

	total, err := a.authTier.TotalCount(ctx, bson.M{})
	if err != nil {
		a.logger.Errorf("failed to get total counts", err)
		return nil, fmt.Errorf("FAILED_TO_GET_CUSTOMER_COUNT")
	}

	return &amount_based_auth_domain.AmountBasedAuthRespose{
		Page:            1,
		AmountBasedAuth: amountBased,
		Limit:           constant.DefaultPerPage,
		Total:           total,
	}, nil
}

func (a AmountBasedAuthRepo) checkExistingAuthTier(ctx context.Context, cpsAction model.CreateCPSAction) error {
	filter := bson.M{
		"maker_phone_number": cpsAction.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsAction.Department,
		"request_action":     model.RequestAuthTier,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingAuthTier, err := a.cpsActionDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get cps action", err)
		err = fmt.Errorf(common_util.UnhandledServerError)
		return err
	} else if existingAuthTier != nil {
		err = fmt.Errorf(common_util.PendingCPSActionExists)
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return err
	}

	return nil
}

func (a AmountBasedAuthRepo) validateAuthTier(ctx context.Context,
	authTier *amount_based_auth_domain.AuthTier, request amount_based_auth_domain.UpdateAmountBasedAuth) error {
	if request.Method == amount_based_auth_domain.OPEN {
		if authTier.MinAmount >= uint64(request.MaxAmount) {
			a.logger.Errorf("min amount cannot be greater than or equal to max amount min: %s, max: %s", authTier.MinAmount, request.MaxAmount)
			err := fmt.Errorf(common_util.OpenMinGEOpenMax)
			return err
		}
		filter := bson.M{
			"method":     amount_based_auth_domain.PIN,
			"is_deleted": false,
		}
		projection := bson.M{
			"max_amount": 1,
		}

		pinTier, err := a.authTier.FindOne(ctx, filter, projection)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				a.logger.Errorf("pin authier not found", err)
				err = fmt.Errorf(common_util.PinAuthorNotFound)
				return err
			}
			a.logger.Errorf("failed to get pin authier", err)
			err = fmt.Errorf(common_util.UnhandledServerError)
		}

		if pinTier.MaxAmount <= uint64(request.MaxAmount) {
			a.logger.Warnf("open tier max amount can not be greater than max amount of pin tier", pinTier.MaxAmount, request.MaxAmount)
			err = fmt.Errorf(common_util.OpenMaxGEPinMax)
			return err
		}
	} else if request.Method == amount_based_auth_domain.PIN {

		if request.MinAmount > 0 {
			if authTier.MaxAmount <= uint64(request.MinAmount) {
				a.logger.Errorf("min amount cannot be greater than or equal to max amount min: %s, max: %s", authTier.MinAmount, request.MaxAmount)
				err := fmt.Errorf(common_util.PinMinGEPinMax)
				return err
			}

			filter := bson.M{
				"method":     amount_based_auth_domain.OPEN,
				"is_deleted": false,
			}
			projection := bson.M{
				"min_amount": 1,
			}

			openTier, err := a.authTier.FindOne(ctx, filter, projection)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					a.logger.Errorf("open authier not found", err)
					err = fmt.Errorf(common_util.TierAuthNotFound)
					return err
				}
				a.logger.Errorf("failed to get open authier", err)
				err = fmt.Errorf(common_util.FailedToGetAuthTier)
			}

			if openTier.MinAmount >= uint64(request.MinAmount) {
				a.logger.Errorf("pin min amount can not be less than open min amount")
				err = fmt.Errorf(common_util.PinMinLEOpenMin)
				return err
			}
		}

		if request.MaxAmount > 0 {
			if authTier.MinAmount >= uint64(request.MaxAmount) {
				a.logger.Errorf("max pin amount should be greater than pin min amount")
				err := fmt.Errorf(common_util.PinMaxLEPinMin)
				return err
			}
		}
	} else if request.Method == amount_based_auth_domain.OTPANDPIN {
		if request.MinAmount > 0 {
			filter := bson.M{
				"method":     amount_based_auth_domain.PIN,
				"is_deleted": false,
			}
			projection := bson.M{
				"min_amount": 1,
			}

			pinTier, err := a.authTier.FindOne(ctx, filter, projection)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					a.logger.Errorf("pin authier not found", err)
					err = fmt.Errorf(common_util.PinAuthorNotFound)
					return err
				}
				a.logger.Errorf("failed to get open authier", err)
				err = fmt.Errorf(common_util.FailedToGetAuthTier)
			}

			if pinTier.MinAmount >= uint64(request.MinAmount) {
				a.logger.Errorf("min amount cannot be greater than or equal to pin min amount min: %s, max: %s", pinTier.MinAmount, request.MinAmount)
				err = fmt.Errorf(common_util.OTPMinGEPinMin)
				return err
			}
		}
	}

	return nil
}

func (a AmountBasedAuthRepo) UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth_domain.UpdateAmountBasedAuth,
	cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error) {

	if err := a.checkExistingAuthTier(ctx, cpsAction); err != nil {
		return nil, err
	}

	objectID, err := bson.ObjectIDFromHex(request.Id)
	if err != nil {
		a.logger.Errorf("Failed to parse object id: %v", err)
		return nil, fmt.Errorf(common_util.InvalidInput)
	}

	filter := bson.M{
		"_id":        objectID,
		"method":     request.Method,
		"is_deleted": false,
	}
	projection := bson.M{}

	authTier, err := a.authTier.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("no auth tier found for ID: %s", request.Id, err)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		a.logger.Errorf("failed to fetch auth tier: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	if err := a.validateAuthTier(ctx, authTier, request); err != nil {
		return nil, err
	}

	authTier.CreatedAt = time.Now()
	authTier.LastModified = time.Now()

	actionInsert, err := a.cpsActionDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		RequestAction:    string(model.RequestAuthTier),
		PreviosAction:    authTier,
		CurrentAction:    request,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	})
	if err != nil {
		a.logger.Errorf("Failed to insert cps action: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	log.Printf("actionInsert %v", actionInsert)

	return lib.MapCPSAction(actionInsert), nil

}

func (a AmountBasedAuthRepo) ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error) {

	cpsUser := contexts.ExtractContext(ctx)
	filter := bson.M{
		"action_code":   id,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_name":         cpsUser.FullName,
		"checker_id":           cpsUser.UserCode,
		"checker_phone_number": cpsUser.PhoneNumber,
		"action_status":        model.ActionApproved,
		"checker_action_time":  time.Now(),
		"last_modified":        time.Now(),
	}

	// Fetch the action document
	savedAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("No action found for ID: %s", id)
			return nil, fmt.Errorf(common_util.AccountNotFound)
		}
		a.logger.Errorf("Failed to fetch action: %v", err)

		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	var actionData amount_based_auth_domain.UpdateAmountBasedAuth

	// Decode CurrentAction
	rawDoc, err := bson.Marshal(savedAction.CurrentAction)
	if err != nil {
		a.logger.Errorf("failed to marshal current action", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	err = bson.Unmarshal(rawDoc, &actionData)
	if err != nil {
		a.logger.Errorf("failed to unmarshal current action", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	// Update the document in the DB
	if actionData.Method == amount_based_auth_domain.OPEN {
		objectId, err := bson.ObjectIDFromHex(actionData.Id)
		var maxAmount, minAmount uint64
		openFilter := bson.M{
			"_id": objectId,
		}
		if actionData.MaxAmount != 0 {
			maxAmount = uint64(actionData.MaxAmount)
		}
		if actionData.MinAmount != 0 {
			minAmount = uint64(actionData.MinAmount)
		}
		updateOpenFilter := bson.M{
			"max_amount":       maxAmount,
			"min_amount":       minAmount,
			"last_modified_at": time.Now(),
		}

		actionData, err := a.authTier.UpdateOne(ctx, openFilter, updateOpenFilter)
		if err != nil {
			a.logger.Errorf("Failed to update open auth tier: %v", err)
			return nil, fmt.Errorf(common_util.UnhandledServerError)
		}
		actionData.LastModified = time.Now()
		actionData.CreatedAt = time.Now()
		savedAction.CurrentAction = actionData

		return lib.MapCPSAction(savedAction), nil
	}

	if actionData.Method == amount_based_auth_domain.PIN {

		objectId, err := bson.ObjectIDFromHex(actionData.Id)
		pinFilter := bson.M{
			"_id": objectId,
		}

		var minAmount, maxAmount uint64

		if actionData.MinAmount > 0 {
			minAmount = uint64(actionData.MinAmount)
		}
		if actionData.MaxAmount > 0 {
			maxAmount = uint64(actionData.MaxAmount)
		}

		updatePinFilter := bson.M{
			"min_amount": minAmount,
			"max_amount": maxAmount,
		}

		actionData, err := a.authTier.UpdateOne(ctx, pinFilter, updatePinFilter)
		if err != nil {
			a.logger.Errorf("Failed to update pin auth tier: %v", err)
			return nil, fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}

		savedAction.CurrentAction = actionData
		return lib.MapCPSAction(savedAction), nil

	}

	if actionData.Method == amount_based_auth_domain.OTPANDPIN {
		objectId, err := bson.ObjectIDFromHex(actionData.Id)
		filter := bson.M{
			"_id": objectId,
		}

		var minAmount, maxAmount uint64

		if actionData.MinAmount > 0 {
			minAmount = uint64(actionData.MinAmount)
		}
		if actionData.MaxAmount > 0 {
			maxAmount = uint64(actionData.MaxAmount)
		}

		updateFilter := bson.M{
			"min_amount": minAmount,
			"max_amount": maxAmount,
		}

		actionData, err := a.authTier.UpdateOne(ctx, filter, updateFilter)
		if err != nil {
			a.logger.Errorf("Failed to update OTP and PIN auth tier: %v", err)
			return nil, fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}

		savedAction.CurrentAction = actionData
		return lib.MapCPSAction(savedAction), nil
	}

	return lib.MapCPSAction(savedAction), nil
}

func (a AmountBasedAuthRepo) RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error) {

	checkerUser := contexts.ExtractContext(ctx)

	filter := bson.M{
		"action_code":   id,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_name":         checkerUser.FullName,
		"checker_id":           checkerUser.UserCode,
		"checker_phone_number": checkerUser.PhoneNumber,
		"action_status":        model.ActionRejected,
		"rejection_reason":     cpsAction.RejectionReason,
		"checker_action_time":  time.Now(),
		"last_modified":        time.Now(),
	}

	savedAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("No action found ")
			return nil, fmt.Errorf(common_util.ActionNotFound)
		}
		a.logger.Errorf("Failed to fetch action: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	return lib.MapCPSAction(savedAction), nil

}
