package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type miniAppMerchantService struct {
	repo   storage.MiniAppMerchantRepository
	logger utils.Logger
}

func NewMiniAppMerchantService(repo storage.MiniAppMerchantRepository, logger utils.Logger) service.MiniAppMerchantService {
	return &miniAppMerchantService{
		repo:   repo,
		logger: logger,
	}
}

// Authorize validates and approves/rejects CPS actions related to mini-app merchants.
func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	m.logger.Infof("Mini App Merchant service authorizing action: %s", cpsAction.ActionCode)

	// TODO: add real authorization logic
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// Create a new Mini App Merchant.
func (m *miniAppMerchantService) Create(ctx context.Context, data *model.MiniAppMerchant) (*model.MiniAppMerchant, error) {
	// Check if merchant already exists
	exist, err := m.repo.Exists(ctx, &model.CheckMiniAppMerchant{
		BankAccountNumber: data.BankAccountNumber,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
	}, nil)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("merchant already exists")
	}

	now := time.Now()
	data.Code = utils.RandomGenerator(10)
	data.CreatedAt = now
	data.LastModifiedAt = now
	data.KYC.Status = model.KYCStatusComplete

	merch, err := m.repo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return merch, nil
}

// Update modifies an existing Mini App Merchant by ID.
func (m *miniAppMerchantService) Update(ctx context.Context, id string, data *model.MiniAppMerchant) (*model.MiniAppMerchant, *model.MiniAppMerchant, error) {
	// Fetch existing merchant
	old, err := m.repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// Check for conflicts with other merchants
	exist, err := m.repo.Exists(ctx, &model.CheckMiniAppMerchant{
		BankAccountNumber: data.BankAccountNumber,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
	}, &model.MiniAppMerchantExistOptions{
		ExcludeID: id,
	})
	if err != nil {
		return nil, nil, err
	}
	if exist {
		return nil, nil, fmt.Errorf("miniapp merchant already exists")
	}

	now := time.Now()

	// Merge old and new data
	updated := &model.MiniAppMerchant{
		ID:                old.ID,
		Code:              old.Code,
		MerchantName:      local_util.NonEmptyString(data.MerchantName, old.MerchantName),
		MerchantType:      local_util.NonEmptyString(data.MerchantType, old.MerchantType),
		PhoneNumber:       local_util.NonEmptyString(data.PhoneNumber, old.PhoneNumber),
		Email:             local_util.NonEmptyString(data.Email, old.Email),
		BankAccountNumber: local_util.NonEmptyString(data.BankAccountNumber, old.BankAccountNumber),
		Enabled:           old.Enabled,
		IsDeleted:         old.IsDeleted,
		CreatedAt:         old.CreatedAt,
		LastModifiedAt:    now,
		KYC: model.KYC{
			Status: old.KYC.Status,
			Representative: model.KYCInformation{
				Name:  local_util.NonEmptyString(data.KYC.Representative.Name, old.KYC.Representative.Name),
				Email: local_util.NonEmptyString(data.KYC.Representative.Email, old.KYC.Representative.Email),
				Phone: local_util.NonEmptyString(data.KYC.Representative.Phone, old.KYC.Representative.Phone),
			},
		},
		Branches: old.Branches,
		MiniApps: old.MiniApps,
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
	m.logger.Infof("Deleting Mini App Merchant ID: %s", id)
	return m.repo.Delete(ctx, id)
}

// EnableOrDisable toggles merchant active status.
func (m *miniAppMerchantService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	fmt.Println("Wee are here in handler too")
	m.logger.Infof("Enabling/Disabling Mini App Merchant ID: %s, Enable: %v", id, enable)
	return m.repo.EnableOrDisable(ctx, id, enable)
}
