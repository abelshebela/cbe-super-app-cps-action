package avatar

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarApplication struct {
	avatarDomain avatar.AvatarDomainService
	logger       utils.Logger
}

type AvatarApplicationService interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CpsAction, error)
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
