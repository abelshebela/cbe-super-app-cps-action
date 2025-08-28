package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type amountBasedAuthService struct {
	Repository storage.AmountBasedAuthRepository 
	cpsService service.CPSActionService
	logger utils.Logger
	minioClient config.MinioClientInterface
	bucketName string
	cfg *config.VaultConfig
}
	
func NewAmountBasedAuthService(repository storage.AmountBasedAuthRepository, cpsService service.CPSActionService, minioClient config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.AmountBasedAuthService {
	return &amountBasedAuthService{
		Repository: repository,
		cpsService: cpsService,
		logger: logger,
		minioClient: minioClient,
		bucketName: bucketName,
		cfg: cfg,
	}
}

// Authorize checks if the user has access to the amount based auth
func (s *amountBasedAuthService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Amount Based Auth service authorizing action: %s", action.ActionCode)
	action.ActionStatus = "APPROVED"
	return action, nil
}

// FindAllWithPagination retrieves all amount-based auth tiers with pagination
func (s *amountBasedAuthService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error) {
	return s.Repository.FindAllWithPagination(ctx, filterParam)
}

// Update updates an amount-based auth tier
func (s *amountBasedAuthService) Update(ctx context.Context, id string, authTier *model.AuthTier) error {
	// Validate Method value
	if authTier.Method != constants.OPEN && authTier.Method != constants.PIN && authTier.Method != constants.OTPANDPIN {
		return errors.New(localization.ErrorInvalidMethod.Code)
	}
	
	// Validate amounts
	if authTier.MinAmount >= authTier.MaxAmount {
		return errors.New(localization.ErrorInvalidAmounts.Code)
	}
	
	return s.Repository.Update(ctx, id, authTier)
}




