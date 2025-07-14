package bank

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BankHandlerService interface {
	GetAllBank(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error)
	GetOneBank(ctx context.Context, id string) (*entity.Bank, error)
	CreateOneBank(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error)
	UpdateOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	DeleteOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error)
	EnableOrDisableBank(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	UpdateLogo(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
}

type BankHandler struct {
	bankDomain service.BankService
	logger     utils.Logger
}

func InitBankHanlder(bankDomin service.BankService, logger utils.Logger) BankHandlerService {
	return &BankHandler{
		bankDomain: bankDomin,
		logger:     logger,
	}
}

func (b *BankHandler) CreateOneBank(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error) {
	reqData := req.ActionData.(dto.CreateBankRequest)
	existing, err := b.CheckExstingBank(ctx, reqData.Name)
	if err != nil {
		return nil, err
	}

	if existing {
		return nil, fmt.Errorf(common_util.BankAlreadyExists)
	}
	cpsRes, err := b.bankDomain.CreateOneBank(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsRes, nil
}

func (b *BankHandler) DeleteOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := b.bankDomain.DeleteOneBank(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankHandler) GetAllBank(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error) {
	banks, err := b.bankDomain.GetAllBank(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (b *BankHandler) GetOneBank(ctx context.Context, id string) (*entity.Bank, error) {
	bank, err := b.bankDomain.GetOneBank(ctx, id)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (b *BankHandler) UpdateOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := b.bankDomain.UpdateOneBank(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankHandler) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error) {
	cpsAction, err := b.bankDomain.Authorize(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankHandler) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
	cpsAction, err := b.bankDomain.Reject(ctx, req)
	if err != nil {
		return nil, err
	}
	return cpsAction, nil
}

func (b *BankHandler) EnableOrDisableBank(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := b.bankDomain.EnableOrDisableBank(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankHandler) CheckExstingBank(ctx context.Context, name string) (bool, error) {

	bank, err := b.bankDomain.CheckExistingBank(ctx, name)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, nil
		}
		return false, err
	}
	if bank {
		return true, nil
	}
	return false, nil
}

func (b *BankHandler) UpdateLogo(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	return b.bankDomain.UpdateLogo(ctx, id, req)
}