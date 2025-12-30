package logistics_merchant_service

import (
	"cbe-super-app-cps-action/internal/constants"
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

func NewLogisticsMerchantService(repo storage.LogisticsMerchantRepository, cpsService service.CPSActionService, cfg *config.VaultConfig, logger utils.Logger) service.LogisticsMerchantService {
	return &LogisticsMerchantService{
		repo:       repo,
		cpsService: cpsService,
		cfg:        cfg,
		logger:     logger,
	}
}

// Authorize implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "LogisticsMerchant", "Authorize")
	defer span.End()

	merchant, err := local_util.JsonUnmarshal[local_model.LogisticsMerchant](cpsAction.CurrentAction)
	if err != nil {
		e.logger.Errorf("Failed to unmarshal current action into merchant: %v", err)
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
	case string(constants.RequestUpdateLogisticsMerchant):
		err = e.repo.Update(ctx, cpsAction.UniqueId, *merchant)
		if err != nil {
			span.AddEvent("Failed to update Logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
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
	case string(constants.RequestEnableLogisticsMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable logistics merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
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
	default:
		e.logger.Errorf("Unsupported action requested, action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)

	}

	cpsAction.CurrentAction = merchant
	e.logger.Infof("Authorization completed for Logistics merchant action, action: %s, id: %s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}

// Create implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Create(ctx context.Context, LogisticsMerchant local_model.LogisticsMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "LogisticsMerchant", "Create")
	defer span.End()

	e.logger.Infof("Creating Logistics merchant, name: %s", LogisticsMerchant.MerchantName)

	exist, err := core.CheckMerchantExists(ctx, e.repo, &types.CheckMiniAppMerchant{
		BankAccountNumber: LogisticsMerchant.BankAccountNumber,
	}, nil)
	if err != nil && err.Error() != localization.ErrorLogisticMerchantNotFound.Code {
		e.logger.Errorf("Failed to check merchant existence: %v", err)
		span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	if exist {
		e.logger.Warnf("Logistics merchant already exists with bank account: %s", LogisticsMerchant.BankAccountNumber)
		span.AddEvent("Merchant already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
			attribute.String("bank_account_number", LogisticsMerchant.BankAccountNumber),
		))
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	if LogisticsMerchant.MerchantType == "merchant" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, LogisticsMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			e.logger.Errorf("Account number validation failed: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("bank_account_number", LogisticsMerchant.BankAccountNumber),
			))
			return err
		}
	}

	err = core.HandleCPSActionForLogisticsMerchant(ctx, e.cpsService, "", constants.RequestCreateLogisticsMerchant, LogisticsMerchant, nil, constants.ActionCreate)
	if err != nil {
		e.logger.Errorf("CPS action failed: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	e.logger.Infof("Logistics merchant created successfully, name: %s", LogisticsMerchant.MerchantName)
	return nil
}

// Delete implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "LogisticsMerchant", "Delete")
	defer span.End()

	e.logger.Infof("Deleting Logistics merchant, id: %s", id)

	prev, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to find merchant for deletion, id: %s, error: %v", id, err)
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
		e.logger.Errorf("CPS action failed for merchant deletion, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("Logistics merchant deleted successfully, id: %s", id)
	return nil
}

// EnableOrDisable implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "LogisticsMerchant", "EnableOrDisable")
	defer span.End()

	e.logger.Infof("EnableOrDisable Logistics merchant, id: %s, enable: %v", id, enable)

	prevMerchant, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to find merchant for enable/disable, id: %s, error: %v", id, err)
		span.AddEvent("Merchant not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorLogisticMerchantNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorLogisticMerchantNotFound.Code)
	}

	if enable && prevMerchant.Enabled {
		e.logger.Warnf("Logistics merchant already enabled, id: %s", id)
		span.AddEvent("Merchant already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorLogisticMerchantEnableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorLogisticMerchantEnableFailed.Code)
	}
	if !enable && !prevMerchant.Enabled {
		e.logger.Warnf("Logistics merchant already disabled, id: %s", id)
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
		e.logger.Errorf("CPS action failed for merchant enable/disable, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("Logistics merchant enable/disable completed successfully, id: %s, enabled: %v", id, enable)
	return nil
}

// FindAllWithPagination implements service.LogisticsMerchantService.
func (e *LogisticsMerchantService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*local_model.LogisticsMerchant], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "LogisticsMerchant", "FindAllWithPagination")
	defer span.End()

	e.logger.Infof("Finding all Logistics merchants with filter: %+v", filterParam)
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

	e.logger.Infof("Finding Logistics merchant by ID: %s", id)
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

	e.logger.Infof("Updating Logistics merchant, id: %s", id)

	old, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to find merchant by ID: %s, error: %v", id, err)
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
			e.logger.Errorf("Failed to check merchant existence for update: %v", err)
			span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
		}
		if exist {
			e.logger.Warnf("Logistics merchant with updated data already exists, id: %s", id)
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
			e.logger.Errorf("Account number validation failed: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	err = core.HandleCPSActionForLogisticsMerchant(ctx, e.cpsService, id, constants.RequestUpdateLogisticsMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for merchant update, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("Logistics merchant updated successfully, id: %s", id)
	return nil
}
