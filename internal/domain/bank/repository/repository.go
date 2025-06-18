package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/entity"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type Repository interface {
	CreateBank(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	UpdateBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	DeleteBank(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	GetBank(ctx context.Context, id string) (*entity.Bank, error)
	GetAllBanks(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error)
	EnableOrDisableWallet(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
}
