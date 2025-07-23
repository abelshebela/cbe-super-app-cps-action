package persistence

import (
	"context"
	"errors"
	"fmt"

	// "net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/lib"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
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

func (a *AmountBasedAuthRepo) GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*amount_based_auth_domain.AuthTier], error) {
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

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

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

	meta := common_util.BuildPaginationMeta(total, page, limit)

	return &common_util.PaginatedResponse[[]*amount_based_auth_domain.AuthTier]{
		Data: amountBased,
		Meta: meta,
	}, nil
}

func (a AmountBasedAuthRepo) checkExistingAuthTier(ctx context.Context, cpsAction model.CreateCPSAction) error {
	filter := bson.M{
		"maker_phone_number": cpsAction.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsAction.Department,
		"request_action":     model.RequestUpdateAmountBasedAuth,
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

// func (a AmountBasedAuthRepo) validateAuthTier(ctx context.Context,
// 	authTier *amount_based_auth_domain.AuthTier, request amount_based_auth_domain.UpdateAmountBasedAuth) error {
// 	if request.Method == amount_based_auth_domain.OPEN {
// 		if authTier.MinAmount >= uint64(request.MaxAmount) {
// 			a.logger.Errorf("min amount cannot be greater than or equal to max amount min: %s, max: %s", authTier.MinAmount, request.MaxAmount)
// 			err := fmt.Errorf(common_util.OpenMinGEOpenMax)
// 			return err
// 		}
// 		filter := bson.M{
// 			"method":     amount_based_auth_domain.PIN,
// 			"is_deleted": false,
// 		}
// 		projection := bson.M{
// 			"max_amount": 1,
// 		}

// 		pinTier, err := a.authTier.FindOne(ctx, filter, projection)
// 		if err != nil {
// 			if errors.Is(err, mongo.ErrNoDocuments) {
// 				a.logger.Errorf("pin authier not found", err)
// 				err = fmt.Errorf(common_util.PinAuthorNotFound)
// 				return err
// 			}
// 			a.logger.Errorf("failed to get pin authier", err)
// 			err = fmt.Errorf(common_util.UnhandledServerError)
// 		}

// 		if pinTier.MaxAmount <= uint64(request.MaxAmount) {
// 			a.logger.Warnf("open tier max amount can not be greater than max amount of pin tier", pinTier.MaxAmount, request.MaxAmount)
// 			err = fmt.Errorf(common_util.OpenMaxGEPinMax)
// 			return err
// 		}
// 	} else if request.Method == amount_based_auth_domain.PIN {

// 		if request.MinAmount > 0 {
// 			if authTier.MaxAmount <= uint64(request.MinAmount) {
// 				a.logger.Errorf("min amount cannot be greater than or equal to max amount min: %s, max: %s", authTier.MinAmount, request.MaxAmount)
// 				err := fmt.Errorf(common_util.PinMinGEPinMax)
// 				return err
// 			}

// 			filter := bson.M{
// 				"method":     amount_based_auth_domain.OPEN,
// 				"is_deleted": false,
// 			}
// 			projection := bson.M{
// 				"min_amount": 1,
// 			}

// 			openTier, err := a.authTier.FindOne(ctx, filter, projection)
// 			if err != nil {
// 				if errors.Is(err, mongo.ErrNoDocuments) {
// 					a.logger.Errorf("open authier not found", err)
// 					err = fmt.Errorf(common_util.TierAuthNotFound)
// 					return err
// 				}
// 				a.logger.Errorf("failed to get open authier", err)
// 				err = fmt.Errorf(common_util.FailedToGetAuthTier)
// 			}

// 			if openTier.MinAmount >= uint64(request.MinAmount) {
// 				a.logger.Errorf("pin min amount can not be less than open min amount")
// 				err = fmt.Errorf(common_util.PinMinLEOpenMin)
// 				return err
// 			}
// 		}

// 		if request.MaxAmount > 0 {
// 			if authTier.MinAmount >= uint64(request.MaxAmount) {
// 				a.logger.Errorf("max pin amount should be greater than pin min amount")
// 				err := fmt.Errorf(common_util.PinMaxLEPinMin)
// 				return err
// 			}
// 		}
// 	} else if request.Method == amount_based_auth_domain.OTPANDPIN {
// 		if request.MinAmount > 0 {
// 			filter := bson.M{
// 				"method":     amount_based_auth_domain.PIN,
// 				"is_deleted": false,
// 			}
// 			projection := bson.M{
// 				"min_amount": 1,
// 			}

// 			pinTier, err := a.authTier.FindOne(ctx, filter, projection)
// 			if err != nil {
// 				if errors.Is(err, mongo.ErrNoDocuments) {
// 					a.logger.Errorf("pin authier not found", err)
// 					err = fmt.Errorf(common_util.PinAuthorNotFound)
// 					return err
// 				}
// 				a.logger.Errorf("failed to get open authier", err)
// 				err = fmt.Errorf(common_util.FailedToGetAuthTier)
// 			}

// 			if pinTier.MinAmount >= uint64(request.MinAmount) {
// 				a.logger.Errorf("min amount cannot be greater than or equal to pin min amount min: %s, max: %s", pinTier.MinAmount, request.MinAmount)
// 				err = fmt.Errorf(common_util.OTPMinGEPinMin)
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

func (a AmountBasedAuthRepo) UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth_domain.UpdateAmountBasedAuth,
	cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error) {

	if err := a.checkExistingAuthTier(ctx, cpsAction); err != nil {

		return nil, err
	}

	ID, err := bson.ObjectIDFromHex(request.ID)
	if err != nil {
		a.logger.Errorf("can't convert ID string to object ID : %v", err)
		return nil, fmt.Errorf("INVALID_ID")
	}
	filter := bson.M{
		"_id":        ID,
		"is_deleted": false,
	}
	projection := bson.M{}

	authTier, err := a.authTier.FindOne(ctx, filter, projection)
	fmt.Println("persistance----------------------------------------------")
	fmt.Println("authTire   :", authTier)
	fmt.Println("persistance----------------------------------------------")

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("no Tire found for given ID: %s ", request.ID, err)
			return nil, fmt.Errorf("INVALID_ID")
		}
		a.logger.Errorf("failed to fetch auth tier step 1: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}
	if authTier.Method == "OPEN" {
		filter := bson.M{
			// "tire":       Tier,
			"method":     "PIN",
			"is_deleted": false,
		}

		projection := bson.M{}
		authTierNext, err := a.authTier.FindOne(ctx, filter, projection)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				a.logger.Errorf("no auth tier found given ID: %v   error : %v", request.ID, err)
				return nil, fmt.Errorf(common_util.NotFound)
			}
			a.logger.Errorf("failed to fetch auth tier step 2: %v", err)
			return nil, fmt.Errorf(common_util.UnhandledServerError)
		}

		if authTierNext.MaxAmount <= request.MaxAmount {
			a.logger.Errorf("Max amount can't be Greater than next")
			return nil, fmt.Errorf(common_util.CurrentMaxGreaterThanNext)
		}
	}
	fmt.Println("min max check=================================================")
	fmt.Printf("priv min: %v  >= incomingMax : %v", authTier.MinAmount, request.MaxAmount)
	fmt.Println("=================================================")

	if authTier.MinAmount >= request.MaxAmount {
		a.logger.Infof("max amount can not be less that or equal to min amount")
		return nil, fmt.Errorf(common_util.MaxTireLessThanMIN)
	}

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
		PreviousAction:   authTier,
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

	return nil, nil

}

func (a AmountBasedAuthRepo) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	fmt.Println("____________________AM hear in amount basedauth persistance Authorize")
	fmt.Printf("action data  %v", cpsAction.CurrentAction)
	fmt.Println("____________________AM hear in amount basedauth persistance Authorize")

	var actionData amount_based_auth_domain.UpdateAmountBasedAuth

	// Decode CurrentAction
	rawDoc, err := bson.Marshal(cpsAction.CurrentAction)
	if err != nil {
		a.logger.Errorf("failed to marshal current action", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	err = bson.Unmarshal(rawDoc, &actionData)
	if err != nil {
		a.logger.Errorf("failed to unmarshal current action", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	maxAmount := actionData.MaxAmount
	ObjID, err := bson.ObjectIDFromHex(actionData.ID)
	if err != nil {
		a.logger.Errorf("[amount_based_auth_persistance]can't convert stiring ID to object id")
		return nil, fmt.Errorf("INVALID_ID")
	}
	openFilter := bson.M{
		"_id":        ObjID,
		"is_deleted": false,
	}

	updateFilter := bson.M{
		"max_amount":       maxAmount,
		"last_modified_at": time.Now(),
	}

	actionDataUpdated, err := a.authTier.UpdateOne(ctx, openFilter, updateFilter)
	if err != nil {
		a.logger.Errorf("Failed to update open auth tier: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}
	if actionDataUpdated.Method == "OPEN" {
		nextMIN := maxAmount
		openFilterNext := bson.M{
			"method":     "PIN",
			"is_deleted": false,
		}

		updateFilterNext := bson.M{

			"min_amount": nextMIN,
			// "last_modified_at": time.Now(),
		}

		_, err := a.authTier.UpdateOne(ctx, openFilterNext, updateFilterNext)
		if err != nil {
			a.logger.Errorf("Failed to update open auth tier: %v", err)
			return nil, fmt.Errorf(common_util.UnhandledServerError)
		}
	} else {
		nextMIN := maxAmount
		openFilterNext := bson.M{
			"method":     "OTP_PIN",
			"is_deleted": false,
		}

		updateFilterNext := bson.M{

			"min_amount": nextMIN,
			// "last_modified_at": time.Now(),
		}

		_, err := a.authTier.UpdateOne(ctx, openFilterNext, updateFilterNext)
		if err != nil {
			a.logger.Errorf("Failed to update open auth tier: %v", err)
			return nil, fmt.Errorf(common_util.UnhandledServerError)
		}

	}

	cpsAction.CurrentAction = actionDataUpdated
	return cpsAction, nil
}

// 	if actionData.Method == amount_based_auth_domain.PIN {

// 		objectId, err := bson.ObjectIDFromHex(actionData.Id)
// 		pinFilter := bson.M{
// 			"_id": objectId,
// 		}

// 		var minAmount, maxAmount uint64

// 		if actionData.MinAmount > 0 {
// 			minAmount = uint64(actionData.MinAmount)
// 		}
// 		if actionData.MaxAmount > 0 {
// 			maxAmount = uint64(actionData.MaxAmount)
// 		}

// 		updatePinFilter := bson.M{
// 			"min_amount": minAmount,
// 			"max_amount": maxAmount,
// 		}

// 		actionData, err := a.authTier.UpdateOne(ctx, pinFilter, updatePinFilter)
// 		if err != nil {
// 			a.logger.Errorf("Failed to update pin auth tier: %v", err)
// 			return nil, fmt.Errorf("%w", constant.ErrorDefinition{
// 				Code:    http.StatusInternalServerError,
// 				Message: "internal server error",
// 			})
// 		}

// 		cpsAction.CurrentAction = actionData
// 		return cpsAction, nil

// 	}

// 	if actionData.Method == amount_based_auth_domain.OTPANDPIN {
// 		objectId, err := bson.ObjectIDFromHex(actionData.Id)
// 		filter := bson.M{
// 			"_id": objectId,
// 		}

// 		var minAmount, maxAmount uint64

// 		if actionData.MinAmount > 0 {
// 			minAmount = uint64(actionData.MinAmount)
// 		}
// 		if actionData.MaxAmount > 0 {
// 			maxAmount = uint64(actionData.MaxAmount)
// 		}

// 		updateFilter := bson.M{
// 			"min_amount": minAmount,
// 			"max_amount": maxAmount,
// 		}

// 		actionData, err := a.authTier.UpdateOne(ctx, filter, updateFilter)
// 		if err != nil {
// 			a.logger.Errorf("Failed to update OTP and PIN auth tier: %v", err)
// 			return nil, fmt.Errorf("%w", constant.ErrorDefinition{
// 				Code:    http.StatusInternalServerError,
// 				Message: "internal server error",
// 			})
// 		}

// 		cpsAction.CurrentAction = actionData
// 		return cpsAction, nil
// 	}

// 	return cpsAction, nil
// }

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
