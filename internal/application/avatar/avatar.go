package avatar

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarApplication struct {
	avatarDomain avatar.AvatarDomainService
	logger       utils.Logger
	cpsService   cps_service.CPSActionService
}

type AvatarApplicationService interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (*entities.CPSAction, error)
	DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*entities.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*avatar.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*avatar.Avatar], error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error)
}

func InitAvatarAPP(avatar_domain avatar.AvatarDomainService,
	cpsService cps_service.CPSActionService,
	logger utils.Logger) AvatarApplicationService {
	return &AvatarApplication{
		avatarDomain: avatar_domain,
		logger:       logger,
		cpsService:   cpsService,
	}
}

func (a *AvatarApplication) CreateAvatar(ctx context.Context, req model.CreateCPSAction) (*entities.CPSAction, error) {
	return a.handleAvatarMakerAction(ctx, func() (*entities.CPSAction, error) {
		return a.avatarDomain.CreateAvatar(ctx, req)
	}, "[avatar.CreateAvatar]")
}

func (a *AvatarApplication) DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error) {
	return a.handleAvatarMakerAction(ctx, func() (*entities.CPSAction, error) {
		return a.avatarDomain.DeleteAvatar(ctx, id, req)
	}, "[avatar.DeleteAvatar]")
}

func (a *AvatarApplication) EnableOrDisableAvatar(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*entities.CPSAction, error) {
	return a.handleAvatarMakerAction(ctx, func() (*entities.CPSAction, error) {
		return a.avatarDomain.EnableOrDisableAvatar(ctx, id, requestAction, cpsReq)
	}, "[avatar.EnableOrDisableAvatar]")
}

func (a *AvatarApplication) UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error) {
	return a.handleAvatarMakerAction(ctx, func() (*entities.CPSAction, error) {
		return a.avatarDomain.UpdateAvatar(ctx, id, req)
	}, "[avatar.UpdateAvatar]")
}

func (a *AvatarApplication) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*avatar.Avatar], error) {
	avatars, err := a.avatarDomain.GetAllAvatar(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return avatars, nil
}

func (a *AvatarApplication) GetAvatar(ctx context.Context, id string) (*avatar.Avatar, error) {
	avatar, err := a.avatarDomain.GetAvatar(ctx, id)
	if err != nil {
		return nil, err
	}

	return avatar, nil
}

func (a *AvatarApplication) handleAvatarMakerAction(
	ctx context.Context,
	buildAction func() (*entities.CPSAction, error),
	logPrefix string,
) (*entities.CPSAction, error) {
	action, err := buildAction()
	if err != nil {
		a.logger.Errorf("%s build action error: %v", logPrefix, err)
		return nil, err
	}

	_, err = a.cpsService.CPSActionExists(ctx, entities.CheckCPSAction{
		UserCode:      action.MakerID,
		FullName:      action.MakerName,
		Department:    action.Department,
		PhoneNumber:   action.MakerPhoneNumber,
		RequestAction: string(action.RequestAction),
	})
	if err != nil {
		a.logger.Errorf("%s CPSActionExists error: %v", logPrefix, err)
		return nil, err
	}

	return a.cpsService.CreateCPSAction(ctx, action)
}
