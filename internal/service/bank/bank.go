package bank

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BankService struct {
	cpsService service.CPSActionService
	logger     utils.Logger
	repo       storage.BankRepository
}

// Authorize implements service.BankService.
func (b *BankService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	panic("unimplemented")
}

// CreateOneBank implements service.BankService.
func (b *BankService) CreateOneBank(ctx context.Context, bank_request bank_dto.CreateBankRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)

	// bankRequest.Name = r.FormValue("name")
	// 	bankRequest.Code = r.FormValue("code")
	// 	bankRequest.BIC = r.FormValue("bic")
	// 	bankRequest.Logo = fileHeader

	bank := model.Bank{
		Name: bank_request.Name,
		BIC:  bank_request.BIC,
		Code: bank_request.Code,
	}

	action := lib.CpsModelBuilder("", makerData, nil, bank, string(constants.RequestCreateBank), constants.CREATE)

	err := b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneBank implements service.BankService.
func (b *BankService) DeleteOneBank(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)

	action := lib.CpsModelBuilder("", makerData, nil, req, string(constants.RequestCreateBank), constants.CREATE)

	err := b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

// EnableOrDisableBank implements service.BankService.
func (b *BankService) EnableOrDisableBank(ctx context.Context, id string, enableDisable bool) error {
	panic("unimplemented")
}

// GetAllBank implements service.BankService.
func (b *BankService) GetAllBank(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Bank], error) {
	panic("unimplemented")
}

// GetOneBank implements service.BankService.
func (b *BankService) GetOneBank(ctx context.Context, id string) (*model.Bank, error) {
	panic("unimplemented")
}

// UpdateLogo implements service.BankService.
func (b *BankService) UpdateLogo(ctx context.Context, id string, logo bank_dto.UpdateLogo) error {
	panic("unimplemented")
}

// UpdateOneBank implements service.BankService.
func (b *BankService) UpdateOneBank(ctx context.Context, id string, req bank_dto.UpdateBankRequest) error {
	panic("unimplemented")
}

func NewBankService(logger utils.Logger, repo storage.BankRepository, cpsService service.CPSActionService) service.BankService {
	return &BankService{
		logger:     logger,
		repo:       repo,
		cpsService: cpsService,
	}
}
