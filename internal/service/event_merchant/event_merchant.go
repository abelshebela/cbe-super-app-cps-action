package event_merchant_service

import (
	merchantDto "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"

	"cbe-super-app-cps-action/internal/constants"
	erp_merchant_update_dto "cbe-super-app-cps-action/internal/constants/dto/erp_merchant_update"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/event_merchant/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"
	"time"

	event_merchant_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type EventMerchantService struct {
	repo                 storage.EventMerchantRepository
	cpsService           service.CPSActionService
	accountLookupService account_lookup.Account
	merchantLookup       merchant_lookup.MerchantLookupAdapter
	cfg                  *config.VaultConfig
	logger               utils.Logger
}

// Authorize implements service.EventMerchantService.
func (e *EventMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "EventMerchant", "Authorize")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	merchant, err := local_util.JsonUnmarshal[event_merchant_model.EventMerchant](cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[EventMerchSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return model.CPSAction{}, errors.New(localization.ErrorServiceUnhandledServerError.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateEventMerchant):
		err = e.repo.Create(ctx, *merchant)
		if err != nil {
			span.AddEvent("Failed to create event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return model.CPSAction{}, err
		}

		isErp, _ := ctx.Value(constants.ContextKey("is_erp")).(bool)
		if !isErp {
			enabled := true
			dto := erp_merchant_update_dto.ERPUpdateRequest{
				MainAccountNumber: merchant.BankAccountNumber,
				CpsEnabled:        &enabled,
				Branches:          []erp_merchant_update_dto.ERPBranch{},
			}
			if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, true, e.logger); err != nil {
				log.Errorf("[EventMerchSvc][Authorize] ERP create err: %v", err)
				span.AddEvent("Failed to update ERP", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				// return nil, err

			}
		}
	case string(constants.RequestUpdateEventMerchant):
		err = e.repo.Update(ctx, cpsAction.UniqueId, *merchant)
		if err != nil {
			span.AddEvent("Failed to update event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return model.CPSAction{}, err
		}

		isErp, _ := ctx.Value(constants.ContextKey("is_erp")).(bool)
		if !isErp {
			prevMerchant, err := local_util.JsonUnmarshal[event_merchant_model.EventMerchant](cpsAction.PreviousAction)
			if err != nil {
				log.Errorf("[EventMerchSvc][Authorize] unmarshal prev err: %v", err)
				span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				return model.CPSAction{}, errors.New(localization.ErrorServiceUnhandledServerError.Code)
			}
			if prevMerchant.BankAccountNumber != merchant.BankAccountNumber {
				dto := erp_merchant_update_dto.ERPUpdateRequest{
					MainAccountNumber: merchant.BankAccountNumber,
				}
				if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, true, e.logger); err != nil {
					log.Errorf("[EventMerchSvc][Authorize] ERP update err: %v", err)
					span.AddEvent("Failed to update ERP", trace.WithAttributes(
						attribute.String("error", err.Error()),
						attribute.String("unique_id", cpsAction.UniqueId),
					))
					// return nil, err
				}
			}
		}

	case string(constants.RequestDeleteEventMerchant):
		err = e.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return model.CPSAction{}, err
		}

		isErp, _ := ctx.Value(constants.ContextKey("is_erp")).(bool)
		if !isErp {
			enabled := false
			dto := erp_merchant_update_dto.ERPUpdateRequest{
				CpsEnabled: &enabled,
			}
			if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, true, e.logger); err != nil {
				log.Errorf("[EventMerchSvc][Authorize] ERP delete err: %v", err)
				span.AddEvent("Failed to update ERP", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				// return nil, err

			}
		}
	case string(constants.RequestEnableEventMerchant):
		err = e.repo.EnableOrDisable(ctx, []string{cpsAction.UniqueId}, true)
		if err != nil {
			span.AddEvent("Failed to enable event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return model.CPSAction{}, err
		}

		isErp, _ := ctx.Value(constants.ContextKey("is_erp")).(bool)
		if !isErp {
			enabled := true
			dto := erp_merchant_update_dto.ERPUpdateRequest{
				CpsEnabled: &enabled,
			}
			if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, true, e.logger); err != nil {
				log.Errorf("[EventMerchSvc][Authorize] ERP enable err: %v", err)
				span.AddEvent("Failed to update ERP", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				// return nil, err

			}
		}
	case string(constants.RequestDisableEventMerchant):
		err = e.repo.EnableOrDisable(ctx, []string{cpsAction.UniqueId}, false)
		if err != nil {
			span.AddEvent("Failed to disable event merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return model.CPSAction{}, err
		}

		isErp, _ := ctx.Value(constants.ContextKey("is_erp")).(bool)
		if !isErp {
			enabled := false
			dto := erp_merchant_update_dto.ERPUpdateRequest{
				CpsEnabled: &enabled,
			}
			if err := lib.PublishMerchantChangeToERP(ctx, e.cfg, dto, merchant.MerchantID, true, e.logger); err != nil {
				log.Errorf("[EventMerchSvc][Authorize] ERP disable err: %v", err)
				span.AddEvent("Failed to update ERP", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				// return nil, err

			}
		}
	default:
		log.Errorf("[EventMerchSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return model.CPSAction{}, errors.New(localization.ErrorUnsupportedAction.Code)

	}

	cpsAction.CurrentAction = merchant
	log.Infof("[EventMerchSvc][Authorize] done action=%s id=%s", cpsAction.RequestAction, merchant.ID)
	return *cpsAction, nil
}

// Create implements service.EventMerchantService.
func (e *EventMerchantService) Create(ctx context.Context, eventMerchant event_merchant_model.EventMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "EventMerchant", "Create")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	log.Infof("[EventMerchSvc][Create] name: %s", eventMerchant.MerchantName)

	// exist, err := core.CheckMerchantExists(ctx, e.repo, &types.CheckMerchant{
	// 	BankAccountNumber: eventMerchant.BankAccountNumber,
	// 	MerchantCode:      eventMerchant.MerchantID,
	// 	Email:             eventMerchant.Email,
	// 	PhoneNumber:       eventMerchant.PhoneNumber,
	// }, nil)

	exist, err := core.CheckEventMercahntExist(ctx, e.repo, &types.CheckMerchant{
		BankAccountNumber: eventMerchant.BankAccountNumber,
		MerchantCode:      eventMerchant.MerchantID,
		Email:             eventMerchant.Email,
		PhoneNumber:       eventMerchant.PhoneNumber,
	})
	if err != nil {
		log.Errorf("[EventMerchSvc][Create] exist check err: %v", err)
		if err.Error() != localization.ErrorResourceNotFound.Code {
			span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}

	}

	if exist {
		log.Warnf("[EventMerchSvc][Create] already exists acct: %s", eventMerchant.BankAccountNumber)
		span.AddEvent("Merchant already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
			attribute.String("bank_account_number", eventMerchant.BankAccountNumber),
		))
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	eventMerchant.CreatedAt = time.Now()
	eventMerchant.UpdatedAt = time.Now()
	if eventMerchant.MerchantType == "merchant" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, eventMerchant.BankAccountNumber, e.accountLookupService)
		if err != nil {
			log.Errorf("[EventMerchSvc][Create] acct validation err: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("bank_account_number", eventMerchant.BankAccountNumber),
			))
			return err
		}
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, "", constants.RequestCreateEventMerchant, eventMerchant, nil, constants.ActionCreate)
	if err != nil {
		log.Errorf("[EventMerchSvc][Create] cps action err: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	log.Infof("[EventMerchSvc][Create] done name: %s", eventMerchant.MerchantName)
	return nil
}

// Delete implements service.EventMerchantService.
func (e *EventMerchantService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "EventMerchant", "Delete")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	log.Infof("[EventMerchSvc][Delete] id: %s", id)

	prev, err := e.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[EventMerchSvc][Delete] find err id=%s: %v", id, err)
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

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, constants.RequestDeleteEventMerchant, deletedMerchant, *prev, constants.ActionDelete)
	if err != nil {
		log.Errorf("[EventMerchSvc][Delete] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	log.Infof("[EventMerchSvc][Delete] done id: %s", id)
	return nil
}

// EnableOrDisable implements service.EventMerchantService.
func (e *EventMerchantService) EnableOrDisable(ctx context.Context, ids []string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "EventMerchant", "EnableOrDisable")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableEventMerchant
	} else {
		action = constants.RequestDisableEventMerchant
	}

	for _, id := range ids {
		log.Infof("[EventMerchSvc][EnableDisable] id: %s enable: %v", id, enable)

		prevMerchant, err := e.repo.FindByID(ctx, id)
		if err != nil {
			log.Errorf("[EventMerchSvc][EnableDisable] find err id=%s: %v", id, err)
			span.AddEvent("Merchant not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorEventMerchantNotFound.Code),
				attribute.String("id", id),
			))
			return errors.New(localization.ErrorEventMerchantNotFound.Code)
		}

		if enable && prevMerchant.Enabled {
			log.Warnf("[EventMerchSvc][EnableDisable] already enabled id: %s", id)
			span.AddEvent("Merchant already enabled", trace.WithAttributes(
				attribute.String("error", localization.ErrorEventMerchantAlreadyEnabled.Code),
				attribute.String("id", id),
			))
			return localization.ErrorEventMerchantAlreadyEnabled
		}
		if !enable && !prevMerchant.Enabled {
			log.Warnf("[EventMerchSvc][EnableDisable] already disabled id: %s", id)
			span.AddEvent("Merchant already disabled", trace.WithAttributes(
				attribute.String("error", localization.ErrorEventMerchantAlreadyDisabled.Code),
				attribute.String("id", id),
			))
			return localization.ErrorEventMerchantAlreadyDisabled
		}

		updatedMerchant := *prevMerchant
		updatedMerchant.Enabled = enable
		updatedMerchant.UpdatedAt = time.Now()

		err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, action, updatedMerchant, *prevMerchant, constants.ActionUpdate)
		if err != nil {
			log.Errorf("[EventMerchSvc][EnableDisable] cps action err id=%s: %v", id, err)
			span.AddEvent("CPS action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	log.Infof("[EventMerchSvc][EnableDisable] done ids: %v enabled: %v", ids, enable)
	return nil
}

// FindAllWithPagination implements service.EventMerchantService.
func (e *EventMerchantService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]event_merchant_model.EventMerchant], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "EventMerchant", "FindAllWithPagination")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	log.Infof("[EventMerchSvc][FindAll] filter: %+v", filterParam)
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
func (e *EventMerchantService) FindByID(ctx context.Context, id string) (*event_merchant_model.EventMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "EventMerchant", "FindByID")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	log.Infof("[EventMerchSvc][FindByID] id: %s", id)
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
func (e *EventMerchantService) Update(ctx context.Context, id string, eventMerchant event_merchant_model.EventMerchant) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "EventMerchant", "Update")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	log.Infof("[EventMerchSvc][Update] id: %s", id)

	old, err := e.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[EventMerchSvc][Update] find err id=%s: %v", id, err)
		span.AddEvent("Failed to find merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if old.BankAccountNumber == eventMerchant.BankAccountNumber &&
		old.MerchantName == eventMerchant.MerchantName &&
		old.Email == eventMerchant.Email &&
		old.PhoneNumber == eventMerchant.PhoneNumber &&
		old.MerchantType == eventMerchant.MerchantType &&
		old.SettlementMethod == eventMerchant.SettlementMethod {
		log.Infof("[EventMerchSvc][Update] no changes id: %s", id)
		return errors.New(localization.ErrorNoChangesDetected.Code)
	}

	updated := core.MergeEventMerchantData(old, &eventMerchant)

	var check types.CheckMerchant

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
		exist, err := core.CheckEventMercahntExist(ctx, e.repo, &check)
		if err != nil {
			log.Errorf("[EventMerchSvc][Update] exist check err: %v", err)
			span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
		if exist {
			log.Warnf("[EventMerchSvc][Update] already exists id: %s", id)
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
			log.Errorf("[EventMerchSvc][Update] acct validation err: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	err = core.HandleCPSActionForEventMerchant(ctx, e.cpsService, id, constants.RequestUpdateEventMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		log.Errorf("[EventMerchSvc][Update] cps action err id=%s: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	log.Infof("[EventMerchSvc][Update] done id: %s", id)
	return nil
}

func (e *EventMerchantService) EventMerchantLookup(ctx context.Context, merchantID string) (*merchantDto.MerchantLookUpResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "EventMerchantLookup", "MiniAppMerchant", "EventMerchantLookup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, e.logger)

	base := strings.TrimRight(e.cfg.OddoEcommerceBaseUrl, "/")
	// base := "https://qaapisuperapp.cbe.com.et/api/v1/cbesuperapp/ecommerce"
	url := base + "/cps/event/merchant/"
	xAPIKey := e.cfg.ApiKey

	merchantData, err := e.merchantLookup.LookupMerchant(ctx, merchantID, xAPIKey, url)
	if err != nil {
		log.Errorf("[EventMerchSvc][EventMerchantLookup] err: %v", err)
		span.AddEvent("Merchant lookup failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("merchant_id", merchantID),
		))
		span.AddEvent("Merchant lookup failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("merchant_id", merchantID),
		))
		return nil, err
	}

	return &merchantData, nil
}
func NewEventMerchantService(repo storage.EventMerchantRepository, cpsService service.CPSActionService, accountLookupService account_lookup.Account, merchantLookupAdaptor merchant_lookup.MerchantLookupAdapter, cfg *config.VaultConfig, logger utils.Logger) service.EventMerchantService {
	return &EventMerchantService{
		repo:                 repo,
		cpsService:           cpsService,
		accountLookupService: accountLookupService,
		merchantLookup:       merchantLookupAdaptor,
		cfg:                  cfg,
		logger:               logger,
	}
}
