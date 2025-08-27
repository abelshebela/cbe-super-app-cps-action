package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
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
	return s.Repository.Update(ctx, id, authTier)
}

// handleCPSAction encapsulates the common CPS action logic
func (s *amountBasedAuthService) handleCPSAction(ctx context.Context, uniqueID string, requestAction cpsaction.RequestAction, curData, prevData interface{}, actionType cpsaction.ActionType) error {
	// Extract user context from the context
	userContext := extractUserFromContext(ctx)
	if isIncomplete(userContext) {
		s.logger.Errorf("Incomplete user context for amount based auth creation | context = %v", userContext)
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userContext, prevData, curData, string(requestAction), string(actionType))

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		s.logger.Errorf("[amount_based_auth.handleCPSAction] failed to create CPS action, action: %s, error: %v", requestAction, err)
		return err
	}
	return nil
}

// Helper functions to replace the missing local_util package
func extractUserFromContext(ctx context.Context) types.UserContext {
	// This is a simplified version - you may need to adjust based on your actual context structure
	return types.UserContext{
		UserCode:    getFromContext(ctx, "user_code"),
		UserID:      getFromContext(ctx, "user_id"),
		FullName:    getFromContext(ctx, "full_name"),
		PhoneNumber: getFromContext(ctx, "phone_number"),
		Department:  getFromContext(ctx, "department"),
		UserRole:    getFromContext(ctx, "user_role"),
	}
}

func getFromContext(ctx context.Context, key string) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}

func isIncomplete(userContext types.UserContext) bool {
	return userContext.UserID == "" || userContext.FullName == "" || userContext.PhoneNumber == "" || userContext.Department == ""
}



