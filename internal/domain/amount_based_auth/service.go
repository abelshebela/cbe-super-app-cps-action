package amount_based_auth_domain

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service struct {
	repo   Repository
	logger utils.Logger
}

func NewAmountBasedAuthService(repo Repository, logger utils.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (service *Service) UpdateAuthTier(ctx context.Context, request UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {
	if err := request.Valiadate(); err != nil {
		service.logger.Errorf("validation error: %v", err)
		return nil, err
	}

	cpsActionRes, err := service.repo.UpdateAuthTier(ctx, request, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (service *Service) ApproveAuthTierApprove(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error) {
	auth, err := service.repo.ApproveAmountBasedAuth(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (service *Service) RejectAuthTier(ctx context.Context, id string, cpsAction model.RejectCPSAction) (*model.CPSAction, error) {
	if err := cpsAction.Validate(); err != nil {
		service.logger.Errorf("validation error: %v", err)
		return nil, err
	}
	auth, err := service.repo.RejectAmountBasedAuth(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
