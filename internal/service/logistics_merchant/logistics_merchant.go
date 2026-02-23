package logistics_merchant_service

import (
	"cbe-super-app-cps-action/internal/constants"
	erp_merchant_update_dto "cbe-super-app-cps-action/internal/constants/dto/erp_merchant_update"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/logistics_merchant/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"context"
	"errors"
	"time"

	local_model "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type LogisticsMerchantService struct {
	repo                 storage.LogisticsMerchantRepository
	cpsService           service.CPSActionService
	accountLookupService account_lookup.Account
	cfg                  *config.VaultConfig
	logger               utils.Logger
}

func NewLogisticsMerchantService(repo storage.LogisticsMerchantRepository, cpsService service.CPSActionService, accountLookupService account_lookup.Account, cfg *config.VaultConfig, logger utils.Logger) service.LogisticsMerchantService {
	return &LogisticsMerchantService{
		repo:                 repo,
		cpsService:           cpsService,
		accountLookupService: accountLookupService,
		cfg:                  cfg,
		logger:               logger,
	}
}

// Authorize implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "LogisticsMerchant", "Authorize")
	defer span.End()

	merchant, err := local_util.JsonUnmarshal[local_model.LogisticsMerchant](cpsAction.CurrentAction)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorServiceUnhandledServerError.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateLogisticsMerchant):
		err = e.repo.Create(ctx, *merchant)
		if err != nil {
			span.AddEvent("Failed to create Logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		if err := core.UpdateERP(ctx, e.cfg, merchant.BankAccountNumber, merchant.MerchantID, e.logger); err != nil {
			e.logger.Errorf("[LogisMerchSvc][Authorize] ERP create err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err
		}

	case string(constants.RequestUpdateLogisticsMerchant):
		err = e.repo.Update(ctx, cpsAction.UniqueId, *merchant)
		if err != nil {
			span.AddEvent("Failed to update Logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		prevMerchant, err := local_util.JsonUnmarshal[local_model.LogisticsMerchant](cpsAction.PreviousAction)
		if err != nil {
			e.logger.Errorf("[LogisMerchSvc][Authorize] unmarshal prev err: %v", err)
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, errors.New(localization.ErrorServiceUnhandledServerError.Code)
		}
		if prevMerchant.BankAccountNumber != merchant.BankAccountNumber {
			dto := erp_merchant_update_dto.ERPUpdateRequest{
				MainAccountNumber: merchant.BankAccountNumber,
			}
			if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, e.logger); err != nil {
				e.logger.Errorf("[LogisMerchSvc][Authorize] ERP update err: %v", err)
				span.AddEvent("Failed to update ERP", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				// return nil, err
			}
		}

	case string(constants.RequestDeleteLogisticsMerchant):
		err = e.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		enabled := false
		dto := erp_merchant_update_dto.ERPUpdateRequest{
			CpsEnabled: &enabled,
		}
		if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, e.logger); err != nil {
			e.logger.Errorf("[LogisMerchSvc][Authorize] ERP delete err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err

		}
	case string(constants.RequestEnableLogisticsMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		enabled := true
		dto := erp_merchant_update_dto.ERPUpdateRequest{
			CpsEnabled: &enabled,
		}
		if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, e.logger); err != nil {
			e.logger.Errorf("[LogisMerchSvc][Authorize] ERP enable err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err

		}
	case string(constants.RequestDisableLogisticsMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		enabled := false
		dto := erp_merchant_update_dto.ERPUpdateRequest{
			CpsEnabled: &enabled,
		}
		if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, e.logger); err != nil {
			e.logger.Errorf("[LogisMerchSvc][Authorize] ERP disable err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err

		}
	default:
		e.logger.Errorf("[LogisMerchSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)

	}

	cpsAction.CurrentAction = merchant
	e.logger.Infof("[LogisMerchSvc][Authorize] done action=%s id=%s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}

// Create implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Create(ctx context.Context, LogisticsMerchant local_model.LogisticsMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "LogisticsMerchant", "Create")
	defer span.End()

	e.logger.Infof("[LogisMerchSvc][Create] name: %s", LogisticsMerchant.MerchantName)

	exist, err := core.CheckMerchantExists(ctx, e.repo, &types.CheckMiniAppMerchant{
		BankAccountNumber: LogisticsMerchant.BankAccountNumber,
		MerchantCode:      LogisticsMerchant.MerchantID,
	}, nil)
	if err != nil && err.Error() != localization.ErrorLogisticMerchantNotFound.Code {
		e.logger.Errorf("[LogisMerchSvc][Create] exist check err: %v", err)
		span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	if exist {
		e.logger.Warnf("[LogisMerchSvc][Create] already exists acct: %s", LogisticsMerchant.BankAccountNumber)
		span.AddEvent("Merchant already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
			attribute.String("bank_account_number", LogisticsMerchant.BankAccountNumber),
		))
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	if LogisticsMerchant.MerchantType == "merchant" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, LogisticsMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			e.logger.Errorf("[LogisMerchSvc][Create] acct validation err: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("bank_account_number", LogisticsMerchant.BankAccountNumber),
			))
			return err
		}
	}

	err = core.HandleCPSActionForLogisticsMerchant(ctx, e.cpsService, "", constants.RequestCreateLogisticsMerchant, LogisticsMerchant, nil, constants.ActionCreate)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][Create] cps action err: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	e.logger.Infof("[LogisMerchSvc][Create] done name: %s", LogisticsMerchant.MerchantName)
	return nil
}

// Delete implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "LogisticsMerchant", "Delete")
	defer span.End()

	e.logger.Infof("[LogisMerchSvc][Delete] id: %s", id)

	prev, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][Delete] find err id=%s: %v", id, err)
		span.AddEvent("Failed to find merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	now := time.Now()
	deletedMerchant := *prev
	deletedMerchant.IsDeleted = true
	deletedMerchant.DeletedAt = now

	err = core.HandleCPSActionForLogisticsMerchant(ctx, e.cpsService, id, constants.RequestDeleteMiniAppMerchant, deletedMerchant, *prev, constants.ActionDelete)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][Delete] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("[LogisMerchSvc][Delete] done id: %s", id)
	return nil
}

// EnableOrDisable implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "LogisticsMerchant", "EnableOrDisable")
	defer span.End()

	e.logger.Infof("[LogisMerchSvc][EnableDisable] id: %s enable: %v", id, enable)

	prevMerchant, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][EnableDisable] find err id=%s: %v", id, err)
		span.AddEvent("Merchant not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorLogisticMerchantNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorLogisticMerchantNotFound.Code)
	}

	if enable && prevMerchant.Enabled {
		e.logger.Warnf("[LogisMerchSvc][EnableDisable] already enabled id: %s", id)
		span.AddEvent("Merchant already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorLogisticMerchantEnableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorLogisticMerchantEnableFailed.Code)
	}
	if !enable && !prevMerchant.Enabled {
		e.logger.Warnf("[LogisMerchSvc][EnableDisable] already disabled id: %s", id)
		span.AddEvent("Merchant already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorLogisticMerchantDisableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorLogisticMerchantDisableFailed.Code)
	}

	updatedMerchant := *prevMerchant
	updatedMerchant.Enabled = enable
	updatedMerchant.UpdatedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableLogisticsMerchant
	} else {
		action = constants.RequestDisableLogisticsMerchant
	}

	err = core.HandleCPSActionForLogisticsMerchant(ctx, e.cpsService, id, action, updatedMerchant, *prevMerchant, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][EnableDisable] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("[LogisMerchSvc][EnableDisable] done id: %s enabled: %v", id, enable)
	return nil
}

// FindAllWithPagination implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.LogisticsMerchant], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "LogisticsMerchant", "FindAllWithPagination")
	defer span.End()

	e.logger.Infof("[LogisMerchSvc][FindAll] filter: %+v", filterParam)
	result, err := e.repo.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("Failed to find mini app merchants", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

// FindByID implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) FindByID(ctx context.Context, id string) (*local_model.LogisticsMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "LogisticsMerchant", "FindByID")
	defer span.End()

	e.logger.Infof("[LogisMerchSvc][FindByID] id: %s", id)
	result, err := e.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find logistics merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

// Update implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Update(ctx context.Context, id string, LogisticsMerchant local_model.LogisticsMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "LogisticsMerchant", "Update")
	defer span.End()

	e.logger.Infof("[LogisMerchSvc][Update] id: %s", id)

	old, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][Update] find err id=%s: %v", id, err)
		span.AddEvent("Failed to find merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	updated := core.MergeLogisticsMerchantData(old, &LogisticsMerchant)

	var check types.CheckMiniAppMerchant

	if updated.BankAccountNumber != old.BankAccountNumber {
		check.BankAccountNumber = updated.BankAccountNumber
	}

	if check.BankAccountNumber != "" {
		exist, err := core.CheckMerchantExists(ctx, e.repo, &check, &types.MiniAppMerchantExistOptions{ExcludeID: id})
		if err != nil {
			e.logger.Errorf("[LogisMerchSvc][Update] exist check err: %v", err)
			span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
		}
		if exist {
			e.logger.Warnf("[LogisMerchSvc][Update] already exists id: %s", id)
			span.AddEvent("Merchant already exists", trace.WithAttributes(
				attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
		}
	}
	if check.BankAccountNumber != "" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, LogisticsMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			e.logger.Errorf("[LogisMerchSvc][Update] acct validation err: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	err = core.HandleCPSActionForLogisticsMerchant(ctx, e.cpsService, id, constants.RequestUpdateLogisticsMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("[LogisMerchSvc][Update] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("[LogisMerchSvc][Update] done id: %s", id)
	return nil
}
