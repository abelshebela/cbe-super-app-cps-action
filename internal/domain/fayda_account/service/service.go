package service

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FaydaAccountDomain struct {
	FaydaAccountRepo outbound.FaydaAccountRepository
	logger           utils.Logger
}

func InitFaydaAccountDomain(faydaRepository outbound.FaydaAccountRepository, logger utils.Logger) *FaydaAccountDomain {
	return &FaydaAccountDomain{
		FaydaAccountRepo: faydaRepository,
		logger:           logger,
	}
}

func (f *FaydaAccountDomain) InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction, err := f.FaydaAccountRepo.InitiateDisableFaydaAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (f *FaydaAccountDomain) GetAllFaydaAccounts(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	return f.FaydaAccountRepo.GetAllFaydaAccounts(ctx, filterParams)
}

func (f *FaydaAccountDomain) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	return f.FaydaAccountRepo.AuthorizeFaydaAccountDisable(ctx, action)
}
