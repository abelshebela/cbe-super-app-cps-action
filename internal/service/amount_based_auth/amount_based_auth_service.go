package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

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

// Authorize handles persistence for amount-based auth actions
func (s *amountBasedAuthService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Authorizing amount-based auth action, action: %s", action.RequestAction)

	var authTier *model.AuthTier
	if err := local_util.BindAction(action.CurrentAction, &authTier); err != nil {
		s.logger.Errorf("Failed to bind current action to auth tier: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	var err error	
	switch action.RequestAction {
	case string(cpsaction.RequestUpdateAmountBasedAuth):
		authTier.LastModified = time.Now()
		err = s.Repository.Update(ctx, authTier.ID.Hex(), authTier)
		
	default:
		s.logger.Errorf("Unsupported action requested, action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	if err != nil {
		s.logger.Errorf("Failed to process amount-based auth action, action: %s, error: %v", action.RequestAction, err)
		return nil, err
	}

	action.CurrentAction = authTier
	s.logger.Infof("Authorization completed for action, action: %s, id: %s", action.RequestAction, authTier.ID)
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




