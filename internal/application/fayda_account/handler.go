package faydaaccount

import (
	"context"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	// entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
)

type ApplicationService interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction, userID string) (string, error)
	InitiateEnableFaydaAccount(ctx context.Context, req entities.CPSAction, userID string) (string, error)
}

type FaydaHandler struct {
	FaydaDomain service.FaydaAccount
	cpsService  cps_service.CPSActionService
	logger      utils.Logger
}

func InitFaydaHandler(faydaDomain service.FaydaAccount, logger utils.Logger) ApplicationService {
	return &FaydaHandler{
		FaydaDomain: faydaDomain,
		logger:      logger,
	}
}

func (f *FaydaHandler) InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction, userID string) (string, error) {
	actionCode, err := f.FaydaDomain.InitiateDisableFaydaAccount(ctx, &req, userID)

	if err != nil {
		return "", err
	}

	return actionCode, nil
}

func (f *FaydaHandler) InitiateEnableFaydaAccount(ctx context.Context, req entities.CPSAction, userID string) (string, error) {

	ActionCode, err := f.FaydaDomain.InitiateEnableFaydaAccount(ctx, &req, userID)

	if err != nil {
		return "", err
	}
	return ActionCode, nil
}
