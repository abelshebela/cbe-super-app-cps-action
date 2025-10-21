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

	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type miniAppMerchantService struct {
	repo       storage.MiniAppMerchantRepository
	cpsService service.CPSActionService
	miniRepo   storage.MiniAppRepository
	logger     utils.Logger
}

func NewMiniAppMerchantService(repo storage.MiniAppMerchantRepository, cpsService service.CPSActionService, miniRepo storage.MiniAppRepository, logger utils.Logger) service.MiniAppMerchantService {
	return &miniAppMerchantService{
		repo:       repo,
		cpsService: cpsService,
		miniRepo:   miniRepo,
		logger:     logger,
	}
}

func (m *miniAppMerchantService) Create(ctx context.Context, data *model.MiniAppMerchant) (*model.MiniAppMerchant, error) {
	m.logger.Infof("Creating mini app merchant, name: %s", data.MerchantName)

	exist, err := core.CheckMerchantExists(ctx, m.repo, &model.CheckMiniAppMerchant{
		BankAccountNumber: data.BankAccountNumber,
		Email:             data.KYC.Representative.Email,
		PhoneNumber:       data.KYC.Representative.Phone,
	}, nil)
	if err != nil {
		m.logger.Errorf("Failed to check merchant existence: %v", err)
		return nil, errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
	}
	if exist {
		m.logger.Warnf("Merchant already exists with bank account: %s", data.BankAccountNumber)
		return nil, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	if data.ID.IsZero() {
		data.ID = bson.NewObjectID()
	}

	now := time.Now()
	data.Code = utils.RandomGenerator(10)
	data.CreatedAt = now
	data.LastModifiedAt = now
	data.KYC.Status = model.KYCStatusComplete
	data.Email = ""
	data.PhoneNumber = ""
	data.Enabled= true
	err = core.HandleCPSActionForMiniAppMerchant(
		ctx,
		m.cpsService,
		"",
		constants.RequestCreateMiniAppMerchant,
		data,
		nil,
		constants.ActionCreate,
	)
	if err != nil {
		m.logger.Errorf("CPS action failed: %v", err)
		return nil, err
	}

	m.logger.Infof("Mini app merchant created successfully, name: %s", data.MerchantName)
	return data, nil
}

func (m *miniAppMerchantService) Update(ctx context.Context, id string, data *model.MiniAppMerchant) (*model.MiniAppMerchant, *model.MiniAppMerchant, error) {
	m.logger.Infof("Updating mini app merchant, id: %s", id)

	old, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find merchant by ID: %s, error: %v", id, err)
		return nil, nil, err
	}

	updated := core.MergeMiniAppMerchantData(old, data)
	updated.Email = ""
	updated.PhoneNumber = ""

	var check model.CheckMiniAppMerchant

	if updated.BankAccountNumber != old.BankAccountNumber {
		check.BankAccountNumber = updated.BankAccountNumber
	}
	if updated.KYC.Representative.Email != old.KYC.Representative.Email {
		check.Email = updated.KYC.Representative.Email
	}
	if updated.KYC.Representative.Phone != old.KYC.Representative.Phone {
		check.PhoneNumber = updated.KYC.Representative.Phone
	}

	if check.BankAccountNumber != "" || check.Email != "" || check.PhoneNumber != "" {
		exist, err := core.CheckMerchantExists(ctx, m.repo, &check, &model.MiniAppMerchantExistOptions{ExcludeID: id})
		if err != nil {
			m.logger.Errorf("Failed to check merchant existence for update: %v", err)
			return nil, nil, errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
		}
		if exist {
			m.logger.Warnf("Merchant with updated data already exists, id: %s", id)
			return nil, nil, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
		}
	}

	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestUpdateMiniAppMerchant, updated, old, constants.ActionUpdate)
	if err != nil {
		m.logger.Errorf("CPS action failed for merchant update, id: %s, error: %v", id, err)
		return nil, nil, err
	}

	m.logger.Infof("Mini app merchant updated successfully, id: %s", id)
	return updated, old, nil
}

func (m *miniAppMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	m.logger.Infof("Finding all mini app merchants with filter: %+v", filterParam)
	return m.repo.FindAllWithPagination(ctx, *filterParam)
}

func (m *miniAppMerchantService) FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	m.logger.Infof("Finding mini app merchant by ID: %s", id)
	return m.repo.FindByID(ctx, id)
}

func (m *miniAppMerchantService) Delete(ctx context.Context, id string) error {
	m.logger.Infof("Deleting mini app merchant, id: %s", id)

	prev, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find merchant for deletion, id: %s, error: %v", id, err)
		return err
	}

	now := time.Now()
	deletedMerchant := *prev
	deletedMerchant.IsDeleted = true
	deletedMerchant.DeletedAt = now

	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, constants.RequestDeleteMiniAppMerchant, deletedMerchant, *prev, constants.ActionDelete)
	if err != nil {
		m.logger.Errorf("CPS action failed for merchant deletion, id: %s, error: %v", id, err)
		return err
	}

	m.logger.Infof("Mini app merchant deleted successfully, id: %s", id)
	return nil
}

func (m *miniAppMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	m.logger.Infof("EnableOrDisable mini app merchant, id: %s, enable: %v", id, enable)

	prevMerchant, err := m.repo.FindByID(ctx, id)
	if err != nil {
		m.logger.Errorf("Failed to find merchant for enable/disable, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	if enable && prevMerchant.Enabled {
		m.logger.Warnf("Merchant already enabled, id: %s", id)
		return errors.New(localization.ErrorMiniAppMerchantEnableFailed.Code)
	}
	if !enable && !prevMerchant.Enabled {
		m.logger.Warnf("Merchant already disabled, id: %s", id)
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}

	updatedMerchant := *prevMerchant
	updatedMerchant.Enabled = enable
	updatedMerchant.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableMiniAppMerchant
	} else {
		action = constants.RequestDisableMiniAppMerchant
	}

	err = core.HandleCPSActionForMiniAppMerchant(ctx, m.cpsService, id, action, updatedMerchant, *prevMerchant, constants.ActionUpdate)
	if err != nil {
		m.logger.Errorf("CPS action failed for merchant enable/disable, id: %s, error: %v", id, err)
		return err
	}

	m.logger.Infof("Mini app merchant enable/disable completed successfully, id: %s, enabled: %v", id, enable)
	return nil
}

func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	merchant, err := local_util.JsonUnmarshal[model.MiniAppMerchant](cpsAction.CurrentAction)
	if err != nil {
		m.logger.Errorf("Failed to unmarshal current action into merchant: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	merchant.Email = ""
	merchant.PhoneNumber = ""

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniAppMerchant):
		_, err = m.repo.Create(ctx, merchant)
	case string(constants.RequestUpdateMiniAppMerchant):
		err = m.repo.Update(ctx, cpsAction.UniqueId, merchant)
	case string(constants.RequestDeleteMiniAppMerchant):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
		if err == nil {
			_ = core.CascadeDeleteMiniApps(ctx, m.miniRepo, cpsAction.UniqueId)
		}
	case string(constants.RequestEnableMiniAppMerchant):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableMiniAppMerchant):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err == nil {
			_ = core.CascadeEnableDisableMiniApps(ctx, m.miniRepo, cpsAction.UniqueId, false)
		}
	default:
		m.logger.Errorf("Unsupported action requested, action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)

	}

	if err != nil {
		m.logger.Errorf("Failed to process merchant action, action: %s, error: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = merchant
	m.logger.Infof("Authorization completed for merchant action, action: %s, id: %s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}

func (m *miniAppMerchantService) DetailMiniAppByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	m.logger.Infof("Getting mini app merchant details, id: %s", id)
	return m.repo.FindByID(ctx, id)
}
