package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountValidationService struct {
	logger         utils.Logger
	validationRule storage.ValidationRuleRepository
	cpsService     service.CPSActionService
}

func NewAccountValidationService(validationData storage.ValidationRuleRepository, cpsService service.CPSActionService, logger utils.Logger) service.AccountValidationService {
	return &accountValidationService{
		logger:         logger,
		validationRule: validationData,
		cpsService:     cpsService,
	}
}

func (s *accountValidationService) Update(ctx context.Context, id string, rule *model.ValidationRule) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	validationRule, err := s.validationRule.FindByID(ctx, id)
	if err != nil {
		return err
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, validationRule, rule, string(constants.RequestUpdateAccountValidation), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (s *accountValidationService) FindById(ctx context.Context, id string) (*model.ValidationRule, error) {
	return s.validationRule.FindByID(ctx, id)
}

func (s *accountValidationService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error) {
	return s.validationRule.FindAllWithPagination(ctx, filterParam)
}

func (f *accountValidationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	f.logger.Infof("Feedback service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as appro-ved
	if cpsAction.ActionStatus != constants.Approved {
		return nil, errors.New(localization.ErrorCPSActionFailed.Code)
	}

	validationRule, err := local_util.JsonUnmarshal[model.ValidationRule](cpsAction.CurrentAction)
	if err != nil {
		f.logger.Errorf("Failed to unmarshal current action into the model validation rule")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := f.validationRule.Update(ctx, cpsAction.UniqueId, validationRule); err != nil {

		return nil, err
	}
	return cpsAction, nil
}
