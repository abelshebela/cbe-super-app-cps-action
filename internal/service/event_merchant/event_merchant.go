package event_merchant_service

import (
	"cbe-super-app-cps-action/internal/constants"
	erp_merchant_update_dto "cbe-super-app-cps-action/internal/constants/dto/erp_merchant_update"
	"cbe-super-app-cps-action/internal/constants/lib"
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
		e.logger.Errorf("[EventMerchSvc][Authorize] unmarshal err: %v", err)
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
		enabled := true
		dto := erp_merchant_update_dto.ERPUpdateRequest{
			MainAccountNumber: merchant.BankAccountNumber,
			CpsEnabled:        &enabled,
			Branches:          []erp_merchant_update_dto.ERPBranch{},
		}
		if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, e.logger); err != nil {
			e.logger.Errorf("[EventMerchSvc][Authorize] ERP create err: %v", err)
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
			e.logger.Errorf("[EventMerchSvc][Authorize] unmarshal prev err: %v", err)
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
				e.logger.Errorf("[EventMerchSvc][Authorize] ERP update err: %v", err)
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
			span.AddEvent("Failed to delete event merchant", trace.WithAttributes(
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
			e.logger.Errorf("[EventMerchSvc][Authorize] ERP delete err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err

		}
	case string(constants.RequestEnableEventMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable event merchant", trace.WithAttributes(
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
			e.logger.Errorf("[EventMerchSvc][Authorize] ERP enable err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err

		}
	case string(constants.RequestDisableEventMerchant):
		err = e.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable event merchant", trace.WithAttributes(
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
			e.logger.Errorf("[EventMerchSvc][Authorize] ERP disable err: %v", err)
			span.AddEvent("Failed to update ERP", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			// return nil, err

		}
	default:
		e.logger.Errorf("[EventMerchSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)

	}

	cpsAction.CurrentAction = merchant
	e.logger.Infof("[EventMerchSvc][Authorize] done action=%s id=%s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}

// Create implements service.EventMerchantService.
func (e *EventMerchantService) Create(ctx context.Context, eventMerchant model.EventMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "EventMerchant", "Create")
	defer span.End()

	e.logger.Infof("[EventMerchSvc][Create] name: %s", eventMerchant.MerchantName)

	exist, err := core.CheckMerchantExists(ctx, e.repo, &types.CheckMiniAppMerchant{
		BankAccountNumber: eventMerchant.BankAccountNumber,
		MerchantCode:      eventMerchant.MerchantID,
		Email:             eventMerchant.Email,
		PhoneNumber:       eventMerchant.PhoneNumber,
	}, nil)
	if err != nil {
		e.logger.Errorf("[EventMerchSvc][Create] exist check err: %v", err)
		span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	if exist {
		e.logger.Warnf("[EventMerchSvc][Create] already exists acct: %s", eventMerchant.BankAccountNumber)
		span.AddEvent("Merchant already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
			attribute.String("bank_account_number", eventMerchant.BankAccountNumber),
		))
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	if eventMerchant.MerchantType == "merchant" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, eventMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			e.logger.Errorf("[EventMerchSvc][Create] acct validation err: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("bank_account_number", eventMerchant.BankAccountNumber),
			))
			return err
		}
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, "", constants.RequestCreateEventMerchant, eventMerchant, nil, constants.ActionCreate)
	if err != nil {
		e.logger.Errorf("[EventMerchSvc][Create] cps action err: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	e.logger.Infof("[EventMerchSvc][Create] done name: %s", eventMerchant.MerchantName)
	return nil
}

// Delete implements service.EventMerchantService.
func (e *EventMerchantService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "EventMerchant", "Delete")
	defer span.End()

	e.logger.Infof("[EventMerchSvc][Delete] id: %s", id)

	prev, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("[EventMerchSvc][Delete] find err id=%s: %v", id, err)
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
		e.logger.Errorf("[EventMerchSvc][Delete] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("[EventMerchSvc][Delete] done id: %s", id)
	return nil
}

// EnableOrDisable implements service.EventMerchantService.
func (e *EventMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "EventMerchant", "EnableOrDisable")
	defer span.End()

	e.logger.Infof("[EventMerchSvc][EnableDisable] id: %s enable: %v", id, enable)

	prevMerchant, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("[EventMerchSvc][EnableDisable] find err id=%s: %v", id, err)
		span.AddEvent("Merchant not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventMerchantNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventMerchantNotFound.Code)
	}

	if enable && prevMerchant.Enabled {
		e.logger.Warnf("[EventMerchSvc][EnableDisable] already enabled id: %s", id)
		span.AddEvent("Merchant already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorEventMerchantEnableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorEventMerchantEnableFailed.Code)
	}
	if !enable && !prevMerchant.Enabled {
		e.logger.Warnf("[EventMerchSvc][EnableDisable] already disabled id: %s", id)
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
		e.logger.Errorf("[EventMerchSvc][EnableDisable] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("[EventMerchSvc][EnableDisable] done id: %s enabled: %v", id, enable)
	return nil
}

// FindAllWithPagination implements service.EventMerchantService.
func (e *EventMerchantService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EventMerchant], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "EventMerchant", "FindAllWithPagination")
	defer span.End()

	e.logger.Infof("[EventMerchSvc][FindAll] filter: %+v", filterParam)
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

	e.logger.Infof("[EventMerchSvc][FindByID] id: %s", id)
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

	e.logger.Infof("[EventMerchSvc][Update] id: %s", id)

	old, err := e.repo.FindByID(ctx, id)
	if err != nil {
		e.logger.Errorf("[EventMerchSvc][Update] find err id=%s: %v", id, err)
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
			e.logger.Errorf("[EventMerchSvc][Update] exist check err: %v", err)
			span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
		}
		if exist {
			e.logger.Warnf("[EventMerchSvc][Update] already exists id: %s", id)
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
			e.logger.Errorf("[EventMerchSvc][Update] acct validation err: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, constants.RequestUpdateEventMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		e.logger.Errorf("[EventMerchSvc][Update] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	e.logger.Infof("[EventMerchSvc][Update] done id: %s", id)
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
