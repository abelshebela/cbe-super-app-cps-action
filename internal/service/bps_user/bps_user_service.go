package bpsuser

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	bps_user_core "cbe-super-app-cps-action/internal/service/bps_user/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bpsUserService struct {
	cpsService service.CPSActionService
	repo       storage.BPSUserRepository
	logger     utils.Logger
}

func NewBPSUserService(repo storage.BPSUserRepository, cpsService service.CPSActionService, logger utils.Logger) service.BPSUserService {
	return &bpsUserService{
		cpsService: cpsService,
		repo:       repo,
		logger:     logger,
	}
}

// Authorize implements service.BPSUserService.
func (b *bpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "BPS User", "Authorize")
	defer span.End()

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		fmt.Printf("failed to marshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	fmt.Printf("JSON bytes: %s\n", string(marshaled))
	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		fmt.Printf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	actionData := bps_user_core.BPSUser_mapper(actionMap.(map[string]interface{}))
	if cpsAction.UniqueId != "" {
		objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] failed to parse unique id", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] failed to parse unique id: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		actionData.ID = objID
	}

	actionData.LastModifiedAt = time.Now()

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateBPSUser):
		if err := b.repo.Create(ctx, actionData); err != nil {
			span.AddEvent("[Authorize] failed to create BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] failed to create BPS user: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] created BPS user for id: %s", cpsAction.UniqueId)
		return nil, nil

	case string(constants.RequestEnableBPSUser):
		actionData.Enabled = true
		b.logger.Infof("[Authorize] enabling BPS user for id: %s", cpsAction.UniqueId)
	case string(constants.RequestDisableBPSUser):
		actionData.Enabled = false
		b.logger.Infof("[Authorize] disabling BPS user for id: %s", cpsAction.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		b.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
	if err := b.repo.Update(ctx, &actionData); err != nil {
		span.AddEvent("[Authorize] failed to update BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		b.logger.Errorf("[Authorize] failed to update BPS user: %v", err)
		return nil, err
	}
	b.logger.Infof("[Authorize] BPS user action authorized successfully for id: %s", cpsAction.UniqueId)
	return nil, nil
}

// FetchUserByUserCode implements service.BPSUserService.
func (b *bpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchUserByUserCode", "BPS User", "FetchUserByUserCode")
	defer span.End()

	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("[FetchUserByUserCode] failed to fetch BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[FetchUserByUserCode] failed to fetch BPS user: %v", err)
		return nil, err
	}
	b.logger.Infof("[FetchUserByUserCode] BPS user retrieved successfully for user_code: %s", userCode)
	return user, nil
}

// GetAllBPSUsers implements service.BPSUserService.
func (b *bpsUserService) GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllBPSUsers", "BPS User", "GetAllBPSUsers")
	defer span.End()

	result, err := b.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[GetAllBPSUsers] failed to fetch BPS users", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		b.logger.Errorf("[GetAllBPSUsers] failed to fetch BPS users: %v", err)
		return nil, err
	}
	b.logger.Infof("[GetAllBPSUsers] retrieved %d BPS users", len(result.Data))
	return result, nil
}

// UpdateBpsUser implements service.BPSUserService.
func (b *bpsUserService) UpdateBpsUser(ctx context.Context, userCode string, status bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBpsUser", "BPS User", "UpdateBpsUser")
	defer span.End()

	b.logger.Infof("[UpdateBpsUser] updating BPS user status, enabled: %v", status)
	makerData := local_util.ExtractUserFromContext(ctx)
	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("[UpdateBpsUser] failed to find BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[UpdateBpsUser] failed to find BPS user: %v", err)
		return err
	}

	if user == nil {
		span.AddEvent("[UpdateBpsUser] BPS user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[UpdateBpsUser] BPS user not found: %s", userCode)
		return errors.New(localization.ErrorUserNotFound.Code)
	}

	if user.Enabled && status {
		span.AddEvent("[UpdateBpsUser] BPS user already enabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[UpdateBpsUser] BPS user already enabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	if !user.Enabled && !status {
		span.AddEvent("[UpdateBpsUser] BPS user already disabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[UpdateBpsUser] BPS user already disabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updatedUser := *user
	updatedUser.Enabled = status

	var requestAction string
	if status {
		requestAction = string(constants.RequestEnableBPSUser)
	} else {
		requestAction = string(constants.RequestDisableBPSUser)

	}

	cpsActionData := lib.CpsModelBuilder(user.ID.Hex(), makerData, user, updatedUser, requestAction, constants.UPDATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		span.AddEvent("[UpdateBpsUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[UpdateBpsUser] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[UpdateBpsUser] BPS user update request created successfully for user_code: %s", userCode)
	return nil
}
func (b *bpsUserService) CreateBPSUser(ctx context.Context, req model.BPSUser) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBPSUser", "BPS User", "CreateBPSUser")
	defer span.End()
	makerData := local_util.ExtractUserFromContext(ctx)

	b.logger.Infof("[CreateBPSUser] creating BPS user with user_code: %s", req.UserCode)

	// Check if user already exists by user_code
	existingUser, err := b.repo.GetByUserCode(ctx, req.UserCode)
	if err == nil && existingUser != nil {
		b.logger.Errorf("[CreateBPSUser] user already exists: %s", req.UserCode)
		return errors.New(localization.ErrorUserAlreadyExists.Code)
	}

	// Check if email already exists
	// dev in using shated 46 so i dont validate eamil for now
	// if req.Email != "" {
	// 	emailUser, err := b.repo.FindByFilterKey(ctx, "email", req.Email)
	// 	if err == nil && emailUser != nil {
	// 		b.logger.Errorf("[CreateBPSUser] email already exists: %s", req.Email)
	// 		return errors.New(localization.ErrorExistEmail.Code)
	// 	}
	// }

	// Check if phone number already exists
	if req.PhoneNumber != "" {
		phoneUser, err := b.repo.FindByFilterKey(ctx, "phone_number", req.PhoneNumber)
		if err == nil && phoneUser != nil {
			b.logger.Errorf("[CreateBPSUser] phone number already exists: %s", req.PhoneNumber)
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
	}

	// Build CPS action model for create
	cpsActionModel := lib.CpsModelBuilder(
		"",                                     // unique id
		makerData,                              // maker data
		nil,                                    // old data (nil for create)
		req,                                    // new data
		string(constants.RequestCreateBPSUser), // request action
		constants.CREATE,                       // action type
	)

	// Create CPS action
	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[CreateBPSUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", req.UserCode),
		))
		b.logger.Errorf("[CreateBPSUser] failed to create CPS action: %v", err)
		return err
	}

	b.logger.Infof("[CreateBPSUser] CPS action created successfully for user_code: %s", req.UserCode)
	return nil
}
