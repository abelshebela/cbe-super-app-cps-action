package avatar

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarApplicationService interface {
	CreateAvatar(ctx context.Context, req avatar.AvatarRequest, maker entities.User) error
	UpdateAvatar(ctx context.Context, id string, req avatar.AvatarRequest, maker entities.User) error
	DeleteAvatar(ctx context.Context, id string, maker entities.User) error
	EnableOrDisableAvatar(ctx context.Context, id string, maker entities.User, enable bool) error
	GetAvatar(ctx context.Context, id string) (*avatar.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*avatar.Avatar], error)
}

type AvatarApplication struct {
	avatarDomain avatar.AvatarDomainService
	cpsService   cps_service.CPSActionService
	logger       utils.Logger
}

func InitAvatarAPP(avatarDomain avatar.AvatarDomainService, cpsService cps_service.CPSActionService, logger utils.Logger) AvatarApplicationService {
	return &AvatarApplication{
		avatarDomain: avatarDomain,
		cpsService:   cpsService,
		logger:       logger,
	}
}

func (a *AvatarApplication) handleCPSAction(ctx context.Context, maker entities.User, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	_, err := a.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("Failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *AvatarApplication) CreateAvatar(ctx context.Context, req avatar.AvatarRequest, maker entities.User) error {
	res, err := a.avatarDomain.CreateAvatar(ctx, req)
	if err != nil {
		a.logger.Errorf("[avatar.CreateAvatar] Failed to create avatar: %v", err)
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestCreateAvatar, res, nil, cps_const.ActionCreate)
}

func (a *AvatarApplication) UpdateAvatar(ctx context.Context, id string, req avatar.AvatarRequest, maker entities.User) error {

	curAction, prevAction, err := a.avatarDomain.UpdateAvatar(ctx, id, req)
	if err != nil {
		a.logger.Errorf("[avatar.UpdateAvatar] Failed to update avatar: %v", err)
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestUpdateAvatar, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *AvatarApplication) DeleteAvatar(ctx context.Context, id string, maker entities.User) error {
	curAction, prevAction, err := a.avatarDomain.DeleteAvatar(ctx, id)
	if err != nil {
		a.logger.Errorf("[avatar.DeleteAvatar] Failed to delete avatar: %v", err)
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestDeleteAvatar, curAction, prevAction, cps_const.ActionDelete)
}

func (a *AvatarApplication) EnableOrDisableAvatar(ctx context.Context, id string, maker entities.User, enable bool) error {
	curAction, prevAction, err := a.avatarDomain.EnableOrDisableAvatar(ctx, id, enable)
	if err != nil {
		a.logger.Errorf("[avatar.EnableOrDisableAvatar] Failed to enable/disable avatar: %v", err)
		return err
	}

	var action cps_const.RequestAction
	if enable {
		action = cps_const.RequestEnableAvatar
	} else {
		action = cps_const.RequestDisableAvatar
	}

	return a.handleCPSAction(ctx, maker, action, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *AvatarApplication) GetAvatar(ctx context.Context, id string) (*avatar.Avatar, error) {
	avatar, err := a.avatarDomain.GetAvatar(ctx, id)
	if err != nil {
		a.logger.Errorf("[avatar.GetAvatar] Failed to get avatar: %v", err)
		return nil, err
	}
	return avatar, nil
}

func (a *AvatarApplication) GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*avatar.Avatar], error) {
	avatars, err := a.avatarDomain.GetAllAvatar(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("[avatar.GetAllAvatar] Failed to get avatars: %v", err)
		return nil, err
	}
	return avatars, nil
}
