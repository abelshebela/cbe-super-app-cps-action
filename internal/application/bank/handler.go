package bank

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/service"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BankHandlerService interface {
	GetAllBank(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error)
	GetOneBank(ctx context.Context, id string) (*entity.Bank, error)
	CreateOneBank(ctx context.Context, req model.CreateCPSAction) (*model.CpsAction, error)
	UpdateOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error)
	DeleteOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error)
	EnableOrDisableBank(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
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

func (b *BankHandler) CreateOneBank(ctx context.Context, req model.CreateCPSAction) (*model.CpsAction, error) {
	cpsRes, err := b.bankDomain.CreateOneBank(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsRes, nil
}

func (b *BankHandler) DeleteOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error) {
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

func (b *BankHandler) UpdateOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankDomain.UpdateOneBank(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankHandler) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankDomain.Authorize(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankHandler) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankDomain.Reject(ctx, req)
	if err != nil {
		return nil, err
	}
	return cpsAction, nil
}

func (b *BankHandler) EnableOrDisableBank(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankDomain.EnableOrDisableBank(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}