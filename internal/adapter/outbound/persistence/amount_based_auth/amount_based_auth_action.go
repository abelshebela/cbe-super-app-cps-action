package persistence

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/lib"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
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

func normalizePhone(search string) bson.M {
	re := regexp.MustCompile(`^\+?251[79]\d{8}$`)

	trimmedSearch := strings.TrimSpace(search)

	if re.MatchString(trimmedSearch) {
		return bson.M{
			"phone_number": bson.M{
				"$regex":   trimmedSearch,
				"$options": "i",
			},
		}
	}

	return bson.M{
		"phone_number.number": bson.M{
			"$regex":   "^$",
			"$options": "i",
		},
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
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	} else if existingAuthTier != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
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
			err := fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusBadRequest,
				Message: "open min amount cannot be greater than or equal to open max amount",
			})
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
				err = fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusNotFound,
					Message: "authier not found",
				})
				return err
			}
			a.logger.Errorf("failed to get pin authier", err)
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}

		if pinTier.MaxAmount <= uint64(request.MaxAmount) {
			a.logger.Warnf("open tier max amount can not be greater than max amount of pin tier", pinTier.MaxAmount, request.MaxAmount)
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusBadRequest,
				Message: "open tier max amount can not be greater than pin max amount",
			})
			return err
		}
	} else if request.Method == amount_based_auth_domain.PIN {

		if request.MinAmount > 0 {
			if authTier.MaxAmount <= uint64(request.MinAmount) {
				a.logger.Errorf("min amount cannot be greater than or equal to max amount min: %s, max: %s", authTier.MinAmount, request.MaxAmount)
				err := fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusBadRequest,
					Message: "min amount cannot be greater than or equal to max amount",
				})
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
					err = fmt.Errorf("%w", constant.ErrorDefinition{
						Code:    http.StatusNotFound,
						Message: "authier not found",
					})
					return err
				}
				a.logger.Errorf("failed to get open authier", err)
				err = fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				})
			}

			if openTier.MinAmount >= uint64(request.MinAmount) {
				a.logger.Errorf("pin min amount can not be less than open min amount")
				err = fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusBadRequest,
					Message: "pin min amount cannot be less than or equal to open min amount",
				})
				return err
			}
		}

		if request.MaxAmount > 0 {
			if authTier.MinAmount >= uint64(request.MaxAmount) {
				a.logger.Errorf("max pin amount should be greater than pin min amount")
				err := fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusBadRequest,
					Message: "max pin amount should be greater than pin min amount",
				})
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
					err = fmt.Errorf("%w", constant.ErrorDefinition{
						Code:    http.StatusNotFound,
						Message: "authier not found",
					})
					return err
				}
				a.logger.Errorf("failed to get open authier", err)
				err = fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				})
			}

			if pinTier.MinAmount >= uint64(request.MinAmount) {
				a.logger.Errorf("min amount cannot be greater than or equal to pin min amount min: %s, max: %s", pinTier.MinAmount, request.MinAmount)
				err = fmt.Errorf("%w", constant.ErrorDefinition{
					Code:    http.StatusBadRequest,
					Message: "min amount cannot be greater than or equal to pin min amount",
				})
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
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid object ID format",
		})
		return nil, err
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
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "no auth tier found for the provided ID",
			})
			return nil, err
		}
		a.logger.Errorf("failed to fetch auth tier: %v", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	if err := a.validateAuthTier(ctx, authTier, request); err != nil {
		return nil, err
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
		PreviosAction:    authTier,
		CurrentAction:    request,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
	})
	if err != nil {
		a.logger.Errorf("Failed to insert cps action: %v", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
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
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "No action found for the provided ID",
			})
			return nil, err
		}
		a.logger.Errorf("Failed to fetch action: %v", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return nil, err
	}

	var actionData amount_based_auth_domain.UpdateAmountBasedAuth

	// Decode CurrentAction
	rawDoc, err := bson.Marshal(savedAction.CurrentAction)
	if err != nil {
		a.logger.Errorf("failed to marshal current action", err)
		return nil, fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	err = bson.Unmarshal(rawDoc, &actionData)
	if err != nil {
		a.logger.Errorf("failed to unmarshal current action", err)
		return nil, fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
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
			"max_amount": maxAmount,
			"min_amount": minAmount,
		}

		fmt.Println("Filter----------", openFilter)
		actionData, err := a.authTier.UpdateOne(ctx, openFilter, updateOpenFilter)
		if err != nil {
			a.logger.Errorf("Failed to update open auth tier: %v", err)
			return nil, fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}

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
		"rejected_reason":      cpsAction.RejectedReason,
		"checker_action_time":  time.Now(),
		"last_modified":        time.Now(),
	}

	savedAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("No action found ")
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "No action found for the provided action code",
			})
			return nil, err
		}
		a.logger.Errorf("Failed to fetch action: %v", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return nil, err
	}

	return lib.MapCPSAction(savedAction), nil

}
