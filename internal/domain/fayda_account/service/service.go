package service

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/repository"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FaydaAccountDomain struct {
	FaydaAccountRepo repository.Repository
	logger           utils.Logger
}

var _ repository.Repository = (*FaydaAccountDomain)(nil)

func InitFaydaAccountDomain(faydaRepository repository.Repository, logger utils.Logger) *FaydaAccountDomain {
	return &FaydaAccountDomain{
		FaydaAccountRepo: faydaRepository,
		logger:           logger,
	}
}

func (f *FaydaAccountDomain) InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {

	cpsAction, err := f.FaydaAccountRepo.InitiateDisableFaydaAccount(ctx, req)

	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (f *FaydaAccountDomain) AuthorizeFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	cpsAction, err := f.FaydaAccountRepo.AuthorizeFaydaAccountDisable(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (f *FaydaAccountDomain) RejectFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	cpsAction, err := f.FaydaAccountRepo.RejectFaydaAccountDisable(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (f *FaydaAccountDomain) GetAllFaydaAccounts(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	return f.FaydaAccountRepo.GetAllFaydaAccounts(ctx, filterParams)
}
