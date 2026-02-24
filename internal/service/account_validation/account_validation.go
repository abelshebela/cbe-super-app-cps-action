package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "Account Validation", "Update")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("[Update] incomplete user information", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[AccValSvc][Update] incomplete user info")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	validationRule, err := s.validationRule.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[Update] failed to find validation rule", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AccValSvc][Update] find err id %s: %v", id, err)
		return err
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, validationRule, rule, string(constants.RequestUpdateAccountValidation), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("[Update] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AccValSvc][Update] cps action err: %v", err)
		return err
	}

	s.logger.Infof("[AccValSvc][Update] request created id: %s", id)
	return nil
}

func (s *accountValidationService) FindById(ctx context.Context, id string) (*model.ValidationRule, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindById", "Account Validation", "FindById")
	defer span.End()

	rule, err := s.validationRule.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[FindById] failed to find validation rule", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[AccValSvc][FindById] find err id %s: %v", id, err)
		return nil, err
	}
	return rule, nil
}

func (s *accountValidationService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ValidationRule], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "Account Validation", "FindAllWithPagination")
	defer span.End()

	result, err := s.validationRule.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("[FindAllWithPagination] failed to fetch validation rules", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[AccValSvc][FindAll] fetch err: %v", err)
		return nil, err
	}
	return result, nil
}

func (f *accountValidationService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Account Validation", "Authorize")
	defer span.End()

	f.logger.Infof("[AccValSvc][Authorize] action: %s", cpsAction.ActionCode)

	// For now, return the action as appro-ved
	if cpsAction.ActionStatus != constants.Approved {
		span.AddEvent("[Authorize] action status is not approved", trace.WithAttributes(
			attribute.String("status", string(cpsAction.ActionStatus)),
			attribute.String("action_code", cpsAction.ActionCode),
		))
		f.logger.Errorf("[AccValSvc][Authorize] invalid status: %s", cpsAction.ActionStatus)
		return nil, errors.New(localization.ErrorCPSActionFailed.Code)
	}

	validationRule, err := local_util.JsonUnmarshal[model.ValidationRule](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("[Authorize] failed to unmarshal current action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		f.logger.Errorf("[AccValSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := f.validationRule.Update(ctx, cpsAction.UniqueId, validationRule); err != nil {
		span.AddEvent("[Authorize] failed to update validation rule", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		f.logger.Errorf("[AccValSvc][Authorize] update err: %v", err)
		return nil, err
	}
	f.logger.Infof("[AccValSvc][Authorize] authorized id: %s", cpsAction.UniqueId)
	return cpsAction, nil
}
