package amount_based_auth_domain

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service struct {
	repo   AmountBasedAuthRepository
	logger utils.Logger
}

func NewAmountBasedAuthService(repo AmountBasedAuthRepository, logger utils.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (service *Service) UpdateAmountBasedAuth(ctx context.Context, request UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error) {
	if err := request.Valiadate(); err != nil {
		service.logger.Errorf("validation error: %v", err)
		return nil, fmt.Errorf(common_util.InvalidInput)
	}

	return service.repo.UpdateAmountBasedAuth(ctx, request, cpsAction)
}

func (service *Service) GetAllAmountBasedDetail(ctx context.Context, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*AuthTier], error) {

	authTiers, err := service.repo.GetAllAmountBasedDetail(ctx, filerParams)
	if err != nil {
		return nil, err
	}
	return authTiers, nil
}

func (service *Service) ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error) {
	return service.repo.ApproveAmountBasedAuth(ctx, id, cpsAction)
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
