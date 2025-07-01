package avatar

import (
	"context"

	"cbe-super-app-cps-action/internal/adapter/outbound/model"
	"cbe-super-app-cps-action/internal/domain/avatar"
	constant "cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarApplication struct {
	avatarDomain avatar.AvatarDomainService
	logger       utils.Logger
}

type AvatarApplicationService interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CpsAction, error)
	DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (model.CpsAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	GetAvatar(ctx context.Context, id string) (*avatar.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (avatar.AvatarResponse, error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error)
}

func InitAvatarAPP(avatar_domain avatar.AvatarDomainService, logger utils.Logger) AvatarApplicationService {
	return &AvatarApplication{
		avatarDomain: avatar_domain,
		logger:       logger,
	}
}

func (a *AvatarApplication) CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CpsAction, error) {
	cpsRes, err := a.avatarDomain.CreateAvatar(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsRes, nil
}

func (a *AvatarApplication) DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.avatarDomain.DeleteAvatar(ctx, id, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarApplication) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.avatarDomain.Authorize(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarApplication) Reject(ctx context.Context, req model.RejectCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.avatarDomain.Reject(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}
	return cpsAction, nil
}

func (a *AvatarApplication) EnableOrDisableAvatar(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	cpsAction, err := a.avatarDomain.EnableOrDisableAvatar(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (a *AvatarApplication) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (avatar.AvatarResponse, error) {
	avatars, err := a.avatarDomain.GetAllAvatar(ctx, filterParams)
	if err != nil {
		return avatar.AvatarResponse{}, err
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

func (a *AvatarApplication) UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.avatarDomain.UpdateAvatar(ctx, id, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}
