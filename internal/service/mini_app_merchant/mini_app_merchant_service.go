package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/mini_app_merchant/core"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type miniAppMerchantService struct {
	repo       storage.MiniAppMerchantRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewMiniAppMerchantService(repo storage.MiniAppMerchantRepository, cpsService service.CPSActionService, logger utils.Logger) service.MiniAppMerchantService {
	return &miniAppMerchantService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

// Create a new Mini App Merchant.
func (m *miniAppMerchantService) Create(ctx context.Context, data *model.MiniAppMerchant) (*model.MiniAppMerchant, error) {
	m.logger.Debugf(">>> Entered Service.Create with data: %+v", data)

	// Validate bank account number
	if data.BankAccountNumber == "" {
		m.logger.Warnf("Bank account number is required")
		return nil, errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	// 🔎 Check if merchant exists by Bank Account Number
	exist, err := m.repo.Exists(ctx, &model.CheckMiniAppMerchant{
		BankAccountNumber: data.BankAccountNumber,
	}, nil)
	if err != nil {
		m.logger.Errorf("failed checking merchant existence: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if exist != nil {
		m.logger.Warnf("Merchant already exists with BankAccountNumber: %s", data.BankAccountNumber)
		return nil, errors.New(localization.ErrorWalletAlreadyExists.Code)
	}

	// Generate new ObjectID if not provided
	if data.ID.IsZero() {
		m.logger.Debugf("ID is empty, generating new ObjectID")
		data.ID = bson.NewObjectID()
	}
	m.logger.Debugf("New generated id: %s", data.ID.Hex())

	// Prepare fields
	now := time.Now()
	data.Code = utils.RandomGenerator(10)
	data.CreatedAt = now
	data.LastModifiedAt = now
	data.KYC.Status = model.KYCStatusComplete
	m.logger.Debugf("Prepared data for repo.Create: %+v", data)

	// Merge (create mode just uses its own values)
	miniApp := data
	m.logger.Debugf("Merged miniApp data: %+v", miniApp)

	// 🔄 CPS Action
	if err := core.HandleCPSActionForMiniAppMerchant(
		ctx,
		m.cpsService,
		data.ID.Hex(),
		constants.RequestCreateMiniAppMerchant,
		miniApp,
		nil,
		constants.ActionCreate,
	); err != nil {
		m.logger.Errorf("CPS action failed: %v", err)
		return nil, err
	}

	m.logger.Debugf("Successfully finished Service.Create")
	return miniApp, nil
}

// Update modifies an existing Mini App Merchant by ID.
func (m *miniAppMerchantService) Update(ctx context.Context, id string, data *model.MiniAppMerchant) (*model.MiniAppMerchant, *model.MiniAppMerchant, error) {
	// Fetch existing merchant
	old, err := m.repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	// Merge old and new data
	updated := core.MergeMiniAppMerchantData(old, data)
	if err := core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestUpdateMiniAppMerchant, updated, old, constants.ActionUpdate); err != nil {
		m.logger.Errorf("CPS action failed for miniapp %m: %v", updated.Code, err)
		return nil, nil, err
	}
	return updated, old, nil
}

// FindAllWithPagination returns paginated list of merchants based on filter criteria.
func (m *miniAppMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	m.logger.Infof("Finding all Mini App Merchants with filter: %+v", filterParam)
	return m.repo.FindAllWithPagination(ctx, *filterParam)
}

// FindByID fetches a merchant by its ID.
func (m *miniAppMerchantService) FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	return m.repo.FindByID(ctx, id)
}

// Delete removes a merchant (soft delete / archive).
func (m *miniAppMerchantService) Delete(ctx context.Context, id string) error {
	prev, err := m.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	deletedMiniApp := *prev
	deletedMiniApp.IsDeleted = true
	deletedMiniApp.DeletedAt = now

	if err := core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestDeleteWallet, deletedMiniApp, *prev, constants.ActionDelete); err != nil {
		m.logger.Errorf("CPS action failed for  miniapp %m: %v", deletedMiniApp.Code, err)
		return err
	}

	return nil
}

// EnableOrDisable toggles merchant active status.
func (m *miniAppMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	m.logger.Debugf(">>> EnableOrDisable called with id=%s, enable=%v", id, enable)

	// Fetch existing merchant
	prevMiniApp, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find miniAppMerchant by ID=%s: %v", id, err)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	m.logger.Debugf("Found miniAppMerchant: %+v", prevMiniApp)

	// Check current status
	if enable && prevMiniApp.Enabled {
		m.logger.Warnf("Merchant already enabled: ID=%s", id)
		return errors.New(localization.ErrorMiniAppMerchantEnableFailed.Code)
	}
	if !enable && !prevMiniApp.Enabled {
		m.logger.Warnf("Merchant already disabled: ID=%s", id)
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}

	// Prepare updated data
	updatedMiniApp := *prevMiniApp
	updatedMiniApp.Enabled = enable
	updatedMiniApp.LastModifiedAt = time.Now()
	m.logger.Debugf("Prepared updated miniAppMerchant for CPS action: %+v", updatedMiniApp)

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableMiniAppMerchant
	} else {
		action = constants.RequestDisableMiniAppMerchant
	}
	m.logger.Debugf("Selected CPS action: %s", action)

	// Handle CPS action
	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, action, updatedMiniApp, *prevMiniApp, constants.ActionUpdate)
	if err != nil {
		m.logger.Errorf("CPS action failed for miniAppMerchant ID=%s, Code=%s, err=%v", id, updatedMiniApp.Code, err)
		return err
	}

	m.logger.Infof("EnableOrDisable completed successfully for ID=%s, new enabled state=%v", id, enable)
	return nil
}

// Authorize validates and approves/rejects CPS actions related to mini-app merchants.
func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	miniAppMerchant, ok := cpsAction.CurrentAction.(*model.MiniAppMerchant)
	if !ok || miniAppMerchant == nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	var err error
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateWallet):
		_, err = m.repo.Create(ctx, miniAppMerchant)
	case string(constants.RequestUpdateWallet):
		_, err = m.repo.Update(ctx, miniAppMerchant.ID.Hex(), miniAppMerchant)
	case string(constants.RequestDeleteWallet):
		err = m.repo.Delete(ctx, miniAppMerchant.ID.Hex())
	case string(constants.RequestEnableWallet):
		err = m.repo.EnableOrDisable(ctx, miniAppMerchant.ID.Hex(), true)
	case string(constants.RequestDisableWallet):
		err = m.repo.EnableOrDisable(ctx, miniAppMerchant.ID.Hex(), false)
	default:
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	cpsAction.CurrentAction = miniAppMerchant
	return cpsAction, nil
}
