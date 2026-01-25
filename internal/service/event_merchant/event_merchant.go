package event_merchant_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/event_merchant/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type EventMerchantService struct {
	repo                 storage.EventMerchantRepository
	cpsService           service.CPSActionService
	accountLookupService account_lookup.Account
	cfg                  *config.VaultConfig
	logger               utils.Logger
}

// Authorize implements service.EventMerchantService.
func (e *EventMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "EventMerchant", "Authorize")
	defer span.End()

	merchant, err := local_util.JsonUnmarshal[model.EventMerchant](cpsAction.CurrentAction)
	if err != nil {
		e.logger.Errorf("Failed to unmarshal current action into merchant: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorServiceUnhandledServerError.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateEventMerchant):
		err = e.repo.Create(ctx, *merchant)
		if err != nil {
			span.AddEvent("Failed to create event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		if err := core.UpdateERP(ctx, e.cfg, merchant.BankAccountNumber, merchant.MerchantID, e.logger); err != nil {
			e.logger.Errorf("Failed to update ERP after creating logistics merchant: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err
		}
	case string(constants.RequestUpdateEventMerchant):
		err = e.repo.Update(ctx, cpsAction.UniqueId, *merchant)
		if err != nil {
			span.AddEvent("Failed to update event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		prevMerchant, err := local_util.JsonUnmarshal[model.EventMerchant](cpsAction.PreviousAction)
		if err != nil {
			e.logger.Errorf("Failed to unmarshal current action into merchant: %v", err)
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, errors.New(localization.ErrorServiceUnhandledServerError.Code)
		}
		if prevMerchant.BankAccountNumber != merchant.BankAccountNumber {
			if err := core.UpdateERP(ctx, e.cfg, merchant.BankAccountNumber, merchant.MerchantID, e.logger); err != nil {
				e.logger.Errorf("Failed to update ERP after creating logistics merchant: %v", err)
				span.AddEvent("Failed to update ERP", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				// return nil, err
			}
		}

	case string(constants.RequestDeleteEventMerchant):
		err = e.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestEnableEventMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableEventMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable mini app merchant", trace.WithAttributes(
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
	e.logger.Infof("Authorization completed for event merchant action, action: %s, id: %s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}

// Create implements service.EventMerchantService.
func (e *EventMerchantService) Create(ctx context.Context, eventMerchant model.EventMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "EventMerchant", "Create")
	defer span.End()

	e.logger.Infof("Creating event merchant, name: %s", eventMerchant.MerchantName)

	exist, err := core.CheckMerchantExists(ctx, e.repo, &types.CheckMiniAppMerchant{
		BankAccountNumber: eventMerchant.BankAccountNumber,
		MerchantCode:      eventMerchant.MerchantID,
		Email:             eventMerchant.Email,
		PhoneNumber:       eventMerchant.PhoneNumber,
	}, nil)
	if err != nil {
		e.logger.Errorf("Failed to check merchant existence: %v", err)
		span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	if exist {
		e.logger.Warnf("Event merchant already exists with bank account: %s", eventMerchant.BankAccountNumber)
		span.AddEvent("Merchant already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
			attribute.String("bank_account_number", eventMerchant.BankAccountNumber),
		))
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	if eventMerchant.MerchantType == "merchant" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, eventMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			e.logger.Errorf("Account number validation failed: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("bank_account_number", eventMerchant.BankAccountNumber),
			))
			return err
		}
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, "", constants.RequestCreateEventMerchant, eventMerchant, nil, constants.ActionCreate)
	if err != nil {
		e.logger.Errorf("CPS action failed: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	e.logger.Infof("Event merchant created successfully, name: %s", eventMerchant.MerchantName)
	return nil
}

// Delete implements service.EventMerchantService.
func (e *EventMerchantService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "EventMerchant", "Delete")
	defer span.End()

	e.logger.Infof("Deleting event merchant, id: %s", id)

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

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, constants.RequestDeleteMiniAppMerchant, deletedMerchant, *prev, constants.ActionDelete)
	if err != nil {
		e.logger.Errorf("CPS action failed for merchant deletion, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("Event merchant deleted successfully, id: %s", id)
	return nil
}

// EnableOrDisable implements service.EventMerchantService.
func (e *EventMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "EventMerchant", "EnableOrDisable")
	defer span.End()

	e.logger.Infof("EnableOrDisable event merchant, id: %s, enable: %v", id, enable)

	prevMerchant, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to find merchant for enable/disable, id: %s, error: %v", id, err)
		span.AddEvent("Merchant not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventMerchantNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventMerchantNotFound.Code)
	}

	if enable && prevMerchant.Enabled {
		e.logger.Warnf("Event merchant already enabled, id: %s", id)
		span.AddEvent("Merchant already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventMerchantEnableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventMerchantEnableFailed.Code)
	}
	if !enable && !prevMerchant.Enabled {
		e.logger.Warnf("Event merchant already disabled, id: %s", id)
		span.AddEvent("Merchant already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventMerchantDisableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventMerchantDisableFailed.Code)
	}

	updatedMerchant := *prevMerchant
	updatedMerchant.Enabled = enable
	updatedMerchant.UpdatedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableEventMerchant
	} else {
		action = constants.RequestDisableEventMerchant
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, action, updatedMerchant, *prevMerchant, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for merchant enable/disable, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("Event merchant enable/disable completed successfully, id: %s, enabled: %v", id, enable)
	return nil
}

// FindAllWithPagination implements service.EventMerchantService.
func (e *EventMerchantService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EventMerchant], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "EventMerchant", "FindAllWithPagination")
	defer span.End()

	e.logger.Infof("Finding all event merchants with filter: %+v", filterParam)
	result, err := e.repo.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("Failed to find mini app merchants", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

// FindByID implements service.EventMerchantService.
func (e *EventMerchantService) FindByID(ctx context.Context, id string) (*model.EventMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "EventMerchant", "FindByID")
	defer span.End()

	e.logger.Infof("Finding event merchant by ID: %s", id)
	result, err := e.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find mini app merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

// Update implements service.EventMerchantService.
func (e *EventMerchantService) Update(ctx context.Context, id string, eventMerchant model.EventMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "EventMerchant", "Update")
	defer span.End()

	e.logger.Infof("Updating event merchant, id: %s", id)

	old, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to find merchant by ID: %s, error: %v", id, err)
		span.AddEvent("Failed to find merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	updated := core.MergeEventMerchantData(old, &eventMerchant)

	var check types.CheckMiniAppMerchant

	if updated.BankAccountNumber != old.BankAccountNumber {
		check.BankAccountNumber = updated.BankAccountNumber
	}
	if updated.Email != old.Email {
		check.Email = updated.Email
	}
	if updated.PhoneNumber != old.PhoneNumber {
		check.PhoneNumber = updated.PhoneNumber
	}

	if check.BankAccountNumber != "" || check.Email != "" || check.PhoneNumber != "" {
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
			e.logger.Warnf("Event merchant with updated data already exists, id: %s", id)
			span.AddEvent("Merchant already exists", trace.WithAttributes(
				attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
		}
	}
	if check.BankAccountNumber != "" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, eventMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			e.logger.Errorf("Account number validation failed: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, constants.RequestUpdateEventMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("CPS action failed for merchant update, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("Event merchant updated successfully, id: %s", id)
	return nil
}

func NewEventMerchantService(repo storage.EventMerchantRepository, cpsService service.CPSActionService, accountLookupService account_lookup.Account, cfg *config.VaultConfig, logger utils.Logger) service.EventMerchantService {
	return &EventMerchantService{
		repo:                 repo,
		cpsService:           cpsService,
		accountLookupService: accountLookupService,
		cfg:                  cfg,
		logger:               logger,
	}
}
