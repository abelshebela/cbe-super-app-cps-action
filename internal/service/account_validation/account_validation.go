package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountValidationService struct {
	logger         utils.Logger
	validationRule storage.ValidationRuleRepository
}

func NewAccountValidationService(validationData storage.ValidationRuleRepository, logger utils.Logger) service.AccountValidationService {
	return &accountValidationService{
		logger:         logger,
		validationRule: validationData,
	}
}

func (s *accountValidationService) UpdateAccountValidation(ctx context.Context, id string, rule *model.ValidationRule) error {
	return s.validationRule.UpdateAccountValidation(ctx, id, rule)
}

func (s *accountValidationService) GetAccountValidation(ctx context.Context, id string) (*model.ValidationRule, error) {
	return s.validationRule.GetAccountValidationByID(ctx, id)
}

func (s *accountValidationService) GetAllAccountValidation(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error) {
	return s.validationRule.GetAllAccountValidation(ctx, filterParam)
}

func (f *accountValidationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	f.logger.Infof("Feedback service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
