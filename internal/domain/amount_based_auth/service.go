package amount_based_auth_domain

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

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

func (service *Service) UpdateAmountBasedAuth(ctx context.Context, request UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error) {
	if err := request.Valiadate(); err != nil {
		service.logger.Errorf("validation error: %v", err)
		return nil, err
	}

	cpsActionRes, err := service.repo.UpdateAmountBasedAuth(ctx, request, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (service *Service) GetAllAmountBasedDetail(ctx context.Context, filerParams *constant.Filter) (*AmountBasedAuthRespose, error) {

	authTiers, err := service.repo.GetAllAmountBasedDetail(ctx, filerParams)
	if err != nil {
		return nil, err
	}
	return authTiers, nil
}

func (service *Service) ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error) {
	auth, err := service.repo.ApproveAmountBasedAuth(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (service *Service) RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error) {
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
