package service

import (
	// "context"

	"context"
	"errors"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cpsuser "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	serviceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ServiceApplication interface {
	ValidateTiers(tiers []dto.Tier, aboveAmount float64) error
	GetOneService(ctx context.Context, id bson.ObjectID) (*service.Service, error)
	InitCPSAction(user cpsuser.CPSUser, actionData map[string]interface{}, requestAction, actionType string, previousData *service.Service) action.CPSAction
	CreateAction(ctx context.Context, req action.CPSAction) error
}

type serviceApp struct {
	serviceDomain serviceDomain.Repository
	actionDomain  action.Repository
	logger        utils.Logger
}

func NewServiceApp(service serviceDomain.Repository, actions action.Repository, logger utils.Logger) ServiceApplication {
	return &serviceApp{
		serviceDomain: service,
		actionDomain:  actions,
		logger:        logger,
	}
}

// func (s *serviceApp) GetAllService(ctx context.Context) (*[]ServiceResponse, error) {

// }

func (s *serviceApp) GetOneService(ctx context.Context, id bson.ObjectID) (*service.Service, error) {
	service, err := s.serviceDomain.GetOneServiceDetail(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *serviceApp) InitCPSAction(user cpsuser.CPSUser, actionData map[string]interface{}, requestAction, actionType string, previousData *service.Service) action.CPSAction {

	return action.CPSAction{
		ActionCode:       utils.ObjectIDGenerator().String(),
		MakerID:          user.ID.String(),
		MakerName:        user.FullName,
		MakerPhoneNumber: user.PhoneNumber,
		MakerActionTime:  time.Now(),
		Department:       user.Department.String(),
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
		return errors.New("At least one tier is required")
	}

	if tiers[0].Min != 0 {
		return errors.New("The first tier's minimum must start from 0")
	}

	for i := 1; i < len(tiers); i++ {
		if tiers[i].Min != tiers[i-1].Max {
			return errors.New("Tier " + strconv.Itoa(i+1) + " minimum must equal the previous tier's maximum")
		}
		if tiers[i].Max <= tiers[i-1].Max {
			return errors.New("Tier " + strconv.Itoa(i+1) + " maximum must be greater than the previous tier's maximum")
		}
	}

	lastTierMax := tiers[len(tiers)-1].Max
	if aboveAmount != lastTierMax {
		return errors.New("Above amount must match the last tier's maximum value: " + strconv.FormatFloat(lastTierMax, 'f', -1, 64))
	}

	return nil
}

func (s *serviceApp) CreateAction(ctx context.Context, req action.CPSAction) error {
	_, err := s.actionDomain.CreateCpsAction(ctx, req)
	// err := s.serviceDomain.CreateAction(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
