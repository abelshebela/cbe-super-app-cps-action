package amount_based_auth_domain

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service struct {
	repo   Repository
	ctx    context.Context
	logger utils.Logger
}

func NewAmountBasedAuthService(repo Repository, ctx context.Context) *Service {
	return &Service{
		repo: repo,
		ctx:  ctx,
	}
}

func (service Service) BuildAuthTierRequest(request AmountBasedAuthRequest) (*AuthTier, error) {
	tier, err := service.repo.UpdateAuthTier(request, service.ctx, service.logger)
	if err != nil {
		return nil, err
	}

	return &tier, nil
}
