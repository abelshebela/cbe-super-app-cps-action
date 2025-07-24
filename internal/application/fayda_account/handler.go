package faydaaccount

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
)

type ApplicationService interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) error
	InitiateEnableFaydaAccount(ctx context.Context, req entities.CPSAction) error
}

type FaydaHandler struct {
	FaydaDomain service.FaydaAccount
	cpsService  cps_service.CPSActionService
	logger      utils.Logger
}

func InitFaydaHandler(faydaDomain service.FaydaAccount, cpsService cps_service.CPSActionService, logger utils.Logger) ApplicationService {
	return &FaydaHandler{
		FaydaDomain: faydaDomain,
		logger:      logger,
		cpsService:  cpsService,
	}
}

func (f *FaydaHandler) InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) error {
	updatedCps := f.FaydaDomain.InitiateDisableFaydaAccount(ctx, &req)
	_, err := f.cpsService.CreateCPSAction(ctx, updatedCps)

	if err != nil {
		return err
	}

	return nil
}

func (f *FaydaHandler) InitiateEnableFaydaAccount(ctx context.Context, req entities.CPSAction) error {
	updatedCps := f.FaydaDomain.InitiateDisableFaydaAccount(ctx, &req)
	_, err := f.cpsService.CreateCPSAction(ctx, updatedCps)

	if err != nil {
		return err
	}
	return nil
}
