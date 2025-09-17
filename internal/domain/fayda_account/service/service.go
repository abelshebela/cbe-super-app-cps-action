package service

import (
	"context"
	"fmt"
	// "strings"
	"time"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	// entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/repository"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FaydaAccountDomain struct {
	FaydaAccountRepo repo.FaydaRepository
	logger           utils.Logger
}

type FaydaAccount interface {
	InitiateDisableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error)
	InitiateEnableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

func InitFaydaAccountDomain(faydaRepository repo.FaydaRepository, logger utils.Logger) FaydaAccount {
	return &FaydaAccountDomain{
		FaydaAccountRepo: faydaRepository,
		logger:           logger,
	}
}

func (f *FaydaAccountDomain) InitiateDisableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error) {
	req.ActionCode = utils.RandomGenerator(20)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.ActionType = cps_const.ActionUpdate
	req.RequestAction = cps_const.RequestDisableFaydaAccount
	req.ActionStatus = cps_const.ActionPending

	action_code, err := f.FaydaAccountRepo.DisableFaydaAccount(ctx, req, userID)
	if err != nil {
		return "", err
	}
	return action_code, nil
}

func (f *FaydaAccountDomain) InitiateEnableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error) {
	req.ActionCode = utils.RandomGenerator(20)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.ActionType = cps_const.ActionUpdate
	req.RequestAction = cps_const.RequestEnableFaydaAccount
	req.ActionStatus = cps_const.ActionPending
	action_code, err := f.FaydaAccountRepo.EnableFaydaAccount(ctx, req, userID)

	if err != nil {
		return "", err
	}

	return action_code, nil

}

func (s *FaydaAccountDomain) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	reqAction := string(action.RequestAction)

	switch reqAction {
	case string(cps_const.RequestDisableFaydaAccount):

		CPSAction, err := s.FaydaAccountRepo.AuthorizeBlockFaydaUser(ctx, action)
		if err != nil {
			return nil, err
		}
		return &entities.CPSAction{ActionCode: CPSAction.ActionCode}, nil
	case string(cps_const.RequestEnableFaydaAccount):

		CPSAction, err := s.FaydaAccountRepo.AuthorizeEnableFaydaUser(ctx, action)
		if err != nil {
			return nil, err
		}
		return &entities.CPSAction{ActionCode: CPSAction.ActionCode}, nil
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}
