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
	"fmt"
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

	// Prepare fields
	now := time.Now()
	data.Code = utils.RandomGenerator(10)
	data.CreatedAt = now
	data.LastModifiedAt = now
	data.KYC.Status = model.KYCStatusComplete

	// Merge (create mode just uses its own values)
	miniApp := data
	// 🔄 CPS Action
	fmt.Println("data.ID.Hex(),", data.ID)
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

	return miniApp, nil
}

// Update modifies an existing Mini App Merchant by ID.
func (m *miniAppMerchantService) Update(ctx context.Context, id string, data *model.MiniAppMerchant) (*model.MiniAppMerchant, *model.MiniAppMerchant, error) {

	// Fetch existing merchant
	old, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Error fetching old MiniAppMerchant with ID %s: %v", id, err)
		return nil, nil, err
	}
	m.logger.Infof("Fetched old MiniAppMerchant: %+v", old)

	// Merge old and new data
	// updated := core.MergeMiniAppMerchantData(old, data)
	updated := &model.MiniAppMerchant{
		ID:                old.ID,
		Code:              old.Code,
		MerchantName:      data.MerchantName,
		MerchantType:      data.MerchantType,
		KYC:               data.KYC,
		BankAccountNumber: data.BankAccountNumber,
		Branches:          data.Branches,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
		MiniApps:          data.MiniApps,
		Enabled:           data.Enabled,
		IsDeleted:         data.IsDeleted,
		CreatedAt:         old.CreatedAt,
		LastModifiedAt:    time.Now(),
		DeletedAt:         data.DeletedAt,
	}
	m.logger.Infof("Merged MiniAppMerchant data: %+v", updated)

	// Create CPS action
	if err := core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestUpdateMiniAppMerchant, updated, old, constants.ActionUpdate); err != nil {
		m.logger.Errorf("CPS action failed for MiniAppMerchant %s: %v", updated.Code, err)
		return nil, nil, err
	}
	m.logger.Infof("CPS action successfully created for MiniAppMerchant: %s", updated.Code)

	return updated, old, nil
}

// FindAllWithPagination returns paginated list of merchants based on filter criteria.
func (m *miniAppMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	m.logger.Infof("Finding all Mini App Merchants with filter: %+v", filterParam)
	return m.repo.FindAllWithPagination(ctx, *filterParam)
}

// FindByID fetches a merchant by its ID.
func (m *miniAppMerchantService) FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	data, err := m.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return data, nil
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

	if err := core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestDeleteMiniAppMerchant, deletedMiniApp, *prev, constants.ActionDelete); err != nil {
		m.logger.Errorf("CPS action failed for  miniapp %m: %v", deletedMiniApp.Code, err)
		return err
	}

	return nil
}

// EnableOrDisable toggles merchant active status.
func (m *miniAppMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {

	// Fetch existing merchant
	prevMiniApp, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find miniAppMerchant by ID=%s: %v", id, err)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	// Check current status
	if enable && prevMiniApp.Enabled {
		m.logger.Warnf("Merchant already enabled: ID=%s", id)
		return errors.New(localization.ErrorMiniAppMerchantAlredyEnabled.Code)
	}
	if !enable && !prevMiniApp.Enabled {
		m.logger.Warnf("Merchant already disabled: ID=%s", id)
		return errors.New(localization.ErrorMiniAppMerchantAlredyDisabled.Code)
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
	fmt.Printf("We are here there in authorization")
	var current *model.MiniAppMerchant
	err := core.BindAction(cpsAction.CurrentAction, &current)
	if err != nil {
		m.logger.Errorf("failed to bind current action to wallet: %v", err)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	id := cpsAction.UniqueId
	fmt.Printf("Id for delete is here for test too %s", id)

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniAppMerchant):
		_, err = m.repo.Create(ctx, current)

	case string(constants.RequestUpdateMiniAppMerchant):
		err = m.repo.Update(ctx, id, current)

	case string(constants.RequestDeleteMiniAppMerchant):
		fmt.Println("deleteting id", id)
		err = m.repo.Delete(ctx, id)

	case string(constants.RequestEnableMiniAppMerchant):
		err = m.repo.EnableOrDisable(ctx, id, true)

	case string(constants.RequestDisableMiniAppMerchant):
		err = m.repo.EnableOrDisable(ctx, id, false)

	default:
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	// put the merged struct back into cpsAction
	cpsAction.CurrentAction = current
	return cpsAction, nil
}

func (m *miniAppMerchantService) DetailMiniAppByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {

	m.logger.Infof("Mini App Merchant service authorizing action: %s", id)

	return m.repo.FindByID(ctx, id)
}

func (s *miniAppMerchantService) AddMiniApp(ctx context.Context, merchantID string, miniApp model.MiniApps) error {
	return s.repo.AddMiniApp(ctx, merchantID, miniApp)
}
func (s *miniAppMerchantService) UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error {
	return s.repo.UpdateMiniAppEnabledState(ctx, merchantID, miniAppID, enabled)
}
func (s *miniAppMerchantService) SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error {
	return s.repo.SoftDeleteMiniApp(ctx, merchantID, miniAppID)
}
