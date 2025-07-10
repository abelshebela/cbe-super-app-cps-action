package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cpsuser "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user"
	serviceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ServiceApplication interface {
	ValidateTiers(tiers []dto.Tier, aboveAmount float64) error
	GetOneService(ctx context.Context, id string) (*serviceDomain.Service, error)
	InitCPSAction(user cpsuser.CPSUser, actionData map[string]any, requestAction, actionType string, previousData *serviceDomain.Service) action.CPSAction
	CreateAction(ctx context.Context, req action.CPSAction) error
}

type serviceApp struct {
	serviceDomain serviceDomain.ServiceInterface
	actionDomain  action.ServiceImpl
	logger        utils.Logger
}

func NewServiceApp(service serviceDomain.ServiceInterface, actions action.ServiceImpl, logger utils.Logger) ServiceApplication {
	return &serviceApp{
		serviceDomain: service,
		actionDomain:  actions,
		logger:        logger,
	}
}

func (s *serviceApp) GetOneService(ctx context.Context, id string) (*serviceDomain.Service, error) {
	service, err := s.serviceDomain.GetServiceDetailsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *serviceApp) InitCPSAction(user cpsuser.CPSUser, actionData map[string]any, requestAction, actionType string, previousData *serviceDomain.Service) action.CPSAction {
	return action.CPSAction{
		MakerID:          user.ID,
		MakerName:        user.FullName,
		MakerPhoneNumber: user.PhoneNumber,
		MakerActionTime:  time.Now(),
		Department:       user.Department,
		ActionStatus:     "PENDING",
		RequestAction:    action.RequestAction(requestAction),
		ActionType:       action.ActionType(actionType),
		CurrentAction:    actionData,
		PreviosAction:    previousData,
		CreatedAt:        time.Now(),
	}
}

func (s *serviceApp) ValidateTiers(tiers []dto.Tier, aboveAmount float64) error {
	if len(tiers) == 0 {
		return errors.New(error_codes.TiersRequired)
	}

	if tiers[0].Min != 0 {
		return errors.New(error_codes.TiersFirstMinZero)
	}

	for i := 1; i < len(tiers); i++ {
		if tiers[i].Min != tiers[i-1].Max {
			return fmt.Errorf(error_codes.TiersMinMustEqualPrevMax)
		}
		if tiers[i].Max <= tiers[i-1].Max {
			return fmt.Errorf(error_codes.TiersMaxMustIncrease)
		}
	}

	lastTierMax := tiers[len(tiers)-1].Max
	if aboveAmount != lastTierMax {
		return fmt.Errorf(error_codes.TiersAboveAmountMismatch)
	}

	return nil
}

func (s *serviceApp) CreateAction(ctx context.Context, req action.CPSAction) error {
	_, err := s.actionDomain.CreateCpsAction(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
