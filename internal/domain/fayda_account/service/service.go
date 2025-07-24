package service

import (
	"context"
	"time"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/repository"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FaydaAccountDomain struct {
	FaydaAccountRepo repo.FaydaRepository
	logger           utils.Logger
}

type FaydaAccount interface {
	InitiateDisableFaydaAccount(ctx context.Context, req *entities.CPSAction) *entities.CPSAction
	InitiateEnableFaydaAccount(ctx context.Context, req *entities.CPSAction) *entities.CPSAction
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

func InitFaydaAccountDomain(faydaRepository repo.FaydaRepository, logger utils.Logger) FaydaAccount {
	return &FaydaAccountDomain{
		FaydaAccountRepo: faydaRepository,
		logger:           logger,
	}
}

func (f *FaydaAccountDomain) InitiateDisableFaydaAccount(ctx context.Context, req *entities.CPSAction) *entities.CPSAction {
	req.ActionCode = utils.RandomGenerator(20)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.ActionType = cps_const.ActionUpdate
	req.RequestAction = cps_const.RequestDisableFaydaAccount
	req.ActionStatus = cps_const.ActionPending

	return req
}

func (f *FaydaAccountDomain) InitiateEnableFaydaAccount(ctx context.Context, req *entities.CPSAction) *entities.CPSAction {
	req.ActionCode = utils.RandomGenerator(20)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.ActionType = cps_const.ActionUpdate
	req.RequestAction = cps_const.RequestEnableFaydaAccount
	req.ActionStatus = cps_const.ActionPending
	return req
}

func (f *FaydaAccountDomain) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	return f.FaydaAccountRepo.AuthorizeFaydaAccountEnableDisable(ctx, action)
}
