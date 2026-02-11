package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants"
	constant_lib "cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/ecommerce-merchant/core"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	merchantDto "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"
	erp_merchant_update_dto "cbe-super-app-cps-action/internal/constants/dto/erp_merchant_update"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type miniAppMerchantService struct {
	repo                 storage.EcommerceMerchantRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	accountLookupService account_lookup.Account
	merchantLookup       merchant_lookup.MerchantLookupAdapter
	cfg                  config.VaultConfig
}

func NewEcommerceMerchantService(
	repo storage.EcommerceMerchantRepository,
	cpsService service.CPSActionService,
	merchantLookup merchant_lookup.MerchantLookupAdapter,
	logger utils.Logger,
	accountLookupService account_lookup.Account,
	cfg config.VaultConfig,
) service.EcommerceMerchantService {
	return &miniAppMerchantService{
		repo:                 repo,
		cpsService:           cpsService,
		accountLookupService: accountLookupService,
		merchantLookup:       merchantLookup,
		logger:               logger,
		cfg:                  cfg,
	}
}

func (m *miniAppMerchantService) Create(ctx context.Context, req *merchantDto.EcommerceMerchant) (*model.EcommerceMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "MiniAppMerchant", "Create")
	defer span.End()

	data := core.ToMiniAppMerchantDomainFromUpdateDTO(req)

	m.logger.Infof("Creating mini app merchant, name: %s", data.MerchantName)
	exist, err := core.CheckMerchantExists(ctx, m.repo, &types.CheckMiniAppMerchant{
		// BankAccountNumber: data.BankAccountNumber,
		// Email:        data.Email,
		// PhoneNumber:  data.PhoneNumber,
		MerchantCode: data.Code,
	}, nil)
	if err != nil {
		m.logger.Errorf("Failed to check merchant existence: %v", err)
		span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	if exist {
		m.logger.Warnf("Merchant already exists with bank account: %s", data.BankAccountNumber)
		span.AddEvent("Merchant already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
			attribute.String("bank_account_number", data.BankAccountNumber),
		))
		return nil, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	if data.ID.IsZero() {
		data.ID = bson.NewObjectID()
	}

	// _, err = core.ValidateAccountNumberWithExternalAPI(ctx, data.BankAccountNumber, m.accountLookupService)
	// if err != nil {
	// 	m.logger.Errorf("Account number validation failed: %v", err)
	// 	span.AddEvent("Account number validation failed", trace.WithAttributes(
	// 		attribute.String("error", err.Error()),
	// 		attribute.String("bank_account_number", data.BankAccountNumber),
	// 	))
	// 	return nil, err
	// }

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	// data.KYC.Status = string(constants.KYCStatusComplete)
	data.Enabled = true
	err = core.HandleCPSActionForMiniAppMerchant(
		ctx,
		m.cpsService,
		"",
		constants.RequestCreateEcommerceMerchant,
		data,
		nil,
		constants.ActionCreate,
	)
	if err != nil {
		m.logger.Errorf("CPS action failed: %v", err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	m.logger.Infof("Mini app merchant created successfully, name: %s", data.MerchantName)
	return data, nil
}

func (m *miniAppMerchantService) Update(ctx context.Context, id string, req *merchantDto.EcommerceMerchant) (*model.EcommerceMerchant, *model.EcommerceMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "MiniAppMerchant", "Update")
	defer span.End()

	merchantReq := core.ToMiniAppMerchantDomainFromUpdateDTO(req)
	m.logger.Infof("Updating mini app merchant, id: %s", id)
	old, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find merchant by ID: %s, error: %v", id, err)
		span.AddEvent("Failed to find merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, nil, err
	}

	updated := core.MergeMiniAppMerchantData(old, merchantReq)

	var check types.CheckMiniAppMerchant

	if updated.BankAccountNumber != old.BankAccountNumber {
		check.BankAccountNumber = updated.BankAccountNumber
	}

	if check.BankAccountNumber != "" || check.Email != "" || check.PhoneNumber != "" {
		exist, err := core.CheckMerchantExists(ctx, m.repo, &check, &types.MiniAppMerchantExistOptions{ExcludeID: id})
		if err != nil {
			m.logger.Errorf("Failed to check merchant existence for update: %v", err)
			span.AddEvent("Failed to check merchant existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return nil, nil, errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
		}
		if exist {
			m.logger.Warnf("Merchant with updated data already exists, id: %s", id)
			span.AddEvent("Merchant already exists", trace.WithAttributes(
				attribute.String("error", localization.ErrorAccountNumberAlreadyExists.Code),
				attribute.String("id", id),
			))
			return nil, nil, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
		}
	}
	if check.BankAccountNumber != "" {
		_, err = core.ValidateAccountNumberWithExternalAPI(ctx, merchantReq.BankAccountNumber, m.accountLookupService)
		if err != nil {
			m.logger.Errorf("Account number validation failed: %v", err)
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return nil, nil, err
		}
	}

	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestUpdateEcommerceMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		m.logger.Errorf("CPS action failed for merchant update, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, nil, err
	}

	m.logger.Infof("Mini app merchant updated successfully, id: %s", id)
	return updated, old, nil
}

func (m *miniAppMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.EcommerceMerchant], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "MiniAppMerchant", "FindAllWithPagination")
	defer span.End()

	m.logger.Infof("Finding all mini app merchants with filter: %+v", filterParam)
	result, err := m.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		span.AddEvent("Failed to find mini app merchants", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

func (m *miniAppMerchantService) FindByID(ctx context.Context, id string) (*model.EcommerceMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "MiniAppMerchant", "FindByID")
	defer span.End()

	m.logger.Infof("Finding mini app merchant by ID: %s", id)
	result, err := m.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find mini app merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	// response := core.ToMiniAppMerchantResponseDTO(result)
	return result, nil
}

func (m *miniAppMerchantService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "MiniAppMerchant", "Delete")
	defer span.End()

	m.logger.Infof("Deleting mini app merchant, id: %s", id)

	prev, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find merchant for deletion, id: %s, error: %v", id, err)
		span.AddEvent("Failed to find merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	now := time.Now()
	deletedMerchant := *prev
	deletedMerchant.IsDeleted = true
	deletedMerchant.DeletedAt = &now

	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestDeleteEcommerceMerchant, deletedMerchant, *prev, constants.ActionDelete)
	if err != nil {
		m.logger.Errorf("CPS action failed for merchant deletion, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	m.logger.Infof("Mini app merchant deleted successfully, id: %s", id)
	return nil
}

func (m *miniAppMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "MiniAppMerchant", "EnableOrDisable")
	defer span.End()

	m.logger.Infof("EnableOrDisable ecommerce merchant, id: %s, enable: %v", id, enable)

	prevMerchant, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find merchant for enable/disable, id: %s, error: %v", id, err)
		span.AddEvent("Merchant not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorMiniAppMerchantNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	if enable && prevMerchant.Enabled {
		m.logger.Warnf("Merchant already enabled, id: %s", id)
		span.AddEvent("Merchant already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorMiniAppMerchantEnableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorMiniAppMerchantEnableFailed.Code)
	}
	if !enable && !prevMerchant.Enabled {
		m.logger.Warnf("Merchant already disabled, id: %s", id)
		span.AddEvent("Merchant already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorMiniAppMerchantDisableFailed.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}

	updatedMerchant := *prevMerchant
	updatedMerchant.Enabled = enable
	updatedMerchant.UpdatedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableEcommerceMerchant
	} else {
		action = constants.RequestDisableEcommerceMerchant
	}

	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, action, updatedMerchant, *prevMerchant, constants.ActionUpdate)
	if err != nil {
		m.logger.Errorf("CPS action failed for merchant enable/disable, id: %s, error: %v", id, err)
		span.AddEvent("CPS action failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	m.logger.Infof("Mini app merchant enable/disable completed successfully, id: %s, enabled: %v", id, enable)
	return nil
}

func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "MiniAppMerchant", "Authorize")
	defer span.End()

	merchant, err := local_util.JsonUnmarshal[model.EcommerceMerchant](cpsAction.CurrentAction)
	if err != nil {
		m.logger.Errorf("Failed to unmarshal current action into merchant: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	erpBranches := make([]erp_merchant_update_dto.ERPBranch, len(merchant.Branches))
	for i, b := range merchant.Branches {
		erpBranches[i] = erp_merchant_update_dto.ERPBranch{
			Merchant:         b.BranchCode,
			CPSAccountNumber: b.BranchAccountNumber,
			CpsEnabled:       nil,
		}
	}
	ERPUpdate := erp_merchant_update_dto.ERPUpdateRequest{
		MainAccountNumber: merchant.BankAccountNumber,
		Branches:          erpBranches,
		CpsEnabled:        &merchant.Enabled,
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateEcommerceMerchant):
		// err = m.updateERP(ctx, merchant)
		// if err != nil {
		// 	m.logger.Errorf("Failed to update ERP for merchant creation: %v", err)
		// }
		// Convert []model.BranchInformation to []erp_merchant_update_dto.ERPBranch
		constant_lib.PublishMerchantChangeToERP(ctx, &m.cfg, ERPUpdate, merchant.ID.Hex(), m.logger)

		_, err = m.repo.Create(ctx, merchant)
		if err != nil {
			span.AddEvent("Failed to create mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateEcommerceMerchant):
		// err = m.updateERP(ctx, merchant)
		// if err != nil {
		// 	m.logger.Errorf("Failed to update ERP for merchant creation: %v", err)
		// }
		constant_lib.PublishMerchantChangeToERP(ctx, &m.cfg, ERPUpdate, merchant.ID.Hex(), m.logger)

		err = m.repo.Update(ctx, cpsAction.UniqueId, merchant)
		if err != nil {
			span.AddEvent("Failed to update mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDeleteEcommerceMerchant):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestEnableEcommerceMerchant):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableEcommerceMerchant):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	default:
		m.logger.Errorf("Unsupported action requested, action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = merchant
	m.logger.Infof("Authorization completed for merchant action, action: %s, id: %s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}

func (m *miniAppMerchantService) updateERP(ctx context.Context, merchant *model.EcommerceMerchant) error {
	if merchant.Code == "" {
		return errors.New(localization.ErrorMerchantIDRequired.Code)
	}

	erpBranches := make([]merchantDto.ERPUpdateBranch, len(merchant.Branches))
	for i, b := range merchant.Branches {
		erpBranches[i] = merchantDto.ERPUpdateBranch{
			Merchant:         b.BranchCode,
			CPSAccountNumber: b.BranchAccountNumber,
		}
	}

	payload := merchantDto.ERPUpdateMerchantRequest{
		CPSAccountNumber: merchant.BankAccountNumber,
		Branches:         erpBranches,
	}

	// url := m.cfg.OddoEcommerceBaseUrl + "/cps/merchant/update/"
	base := "https://qaapisuperapp.cbe.com.et/api/v1/cbesuperapp/ecommerce"
	url := base + "/cps/merchant/update/"

	xAPIKey := m.cfg.ApiKey
	err := m.merchantLookup.UpdateMerchant(ctx, merchant.Code, payload, xAPIKey, url)
	if err != nil {
		m.logger.Errorf("ERP update failed for merchant %s: %v", merchant.Code, err)
		return err
	}

	return nil
}

func (m *miniAppMerchantService) DetailMiniAppByID(ctx context.Context, id string) (*model.EcommerceMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "DetailMiniAppByID", "MiniAppMerchant", "DetailMiniAppByID")
	defer span.End()

	m.logger.Infof("Getting mini app merchant details, id: %s", id)
	result, err := m.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to get mini app merchant details", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

func (m *miniAppMerchantService) MerchantLookup(ctx context.Context, merchantID string) (*merchantDto.MerchantLookUpResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "MerchantLookup", "MiniAppMerchant", "MerchantLookup")
	defer span.End()

	base := strings.TrimRight(m.cfg.OddoEcommerceBaseUrl, "/")
	// base := "https://qaapisuperapp.cbe.com.et/api/v1/cbesuperapp/ecommerce"
	url := base + "/cps/merchant/"
	xAPIKey := m.cfg.ApiKey

	merchantData, err := m.merchantLookup.LookupMerchant(ctx, merchantID, xAPIKey, url)
	if err != nil {
		m.logger.Errorf("Merchant lookup error : %v", err)
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
