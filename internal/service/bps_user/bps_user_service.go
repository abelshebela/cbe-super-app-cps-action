package bpsuser

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		fmt.Printf("failed to marshal CurrentAction: %v\n", err)
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}
	fmt.Printf("JSON bytes: %s\n", string(marshaled))
	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		fmt.Printf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal to interface{}: %v", err)
	}

	actionData := bps_user_core.BPSUser_mapper(actionMap.(map[string]interface{}))
	if cpsAction.UniqueId != "" {
		objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
		if err != nil {
			b.logger.Errorf("[Authorize] failed to parse unique id: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		actionData.ID = objID
	}

	actionData.LastModifiedAt = time.Now()

	switch cpsAction.RequestAction {
	case string(constants.RequestEnableBPSUser):
		actionData.Enabled = true
		b.logger.Infof("[Authorize] enabling BPS user for id: %s", cpsAction.UniqueId)
	case string(constants.RequestDisableBPSUser):
		actionData.Enabled = false
		b.logger.Infof("[Authorize] disabling BPS user for id: %s", cpsAction.UniqueId)
	default:
		b.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
	if err := b.repo.Update(ctx, &actionData); err != nil {
		b.logger.Errorf("[Authorize] failed to update BPS user: %v", err)
		return nil, err
	}
	b.logger.Infof("[Authorize] BPS user action authorized successfully for id: %s", cpsAction.UniqueId)
	return nil, nil
}

// FetchUserByUserCode implements service.BPSUserService.
func (b *bpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error) {
	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		b.logger.Errorf("[FetchUserByUserCode] failed to fetch BPS user: %v", err)
		return nil, err
	}
	b.logger.Infof("[FetchUserByUserCode] BPS user retrieved successfully for user_code: %s", userCode)
	return user, nil
}

// GetAllBPSUsers implements service.BPSUserService.
func (b *bpsUserService) GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error) {
	result, err := b.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		b.logger.Errorf("[GetAllBPSUsers] failed to fetch BPS users: %v", err)
		return nil, err
	}
	b.logger.Infof("[GetAllBPSUsers] retrieved %d BPS users", len(result.Data))
	return result, nil
}

// UpdateBpsUser implements service.BPSUserService.
func (b *bpsUserService) UpdateBpsUser(ctx context.Context, userCode string, status bool) error {
	b.logger.Infof("[UpdateBpsUser] updating BPS user status, enabled: %v", status)
	makerData := local_util.ExtractUserFromContext(ctx)
	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		b.logger.Errorf("[UpdateBpsUser] failed to find BPS user: %v", err)
		return err
	}

	if user == nil {
		b.logger.Errorf("[UpdateBpsUser] BPS user not found: %s", userCode)
		return errors.New(localization.ErrorUserNotFound.Code)
	}

	if user.Enabled && status {
		b.logger.Errorf("[UpdateBpsUser] BPS user already enabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	if !user.Enabled && !status {
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
		b.logger.Errorf("[UpdateBpsUser] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[UpdateBpsUser] BPS user update request created successfully for user_code: %s", userCode)
	return nil
}
