package faydaaccount

import (
	"context"
	"encoding/json"
	"time"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) error
	InitiateEnableFaydaAccount(ctx context.Context, req entities.CPSAction) error
}

type FaydaHandler struct {
	FaydaDomain *service.FaydaAccountDomain
	logger      utils.Logger
}

func InitFaydaHandler(faydaDomain *service.FaydaAccountDomain, logger utils.Logger) ApplicationService {
	return FaydaHandler{
		FaydaDomain: faydaDomain,
		logger:      logger,
	}
}

func (f FaydaHandler) InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) error {
	var actionData ActionDisableData
	jsonBytes, err := json.Marshal(req.CurrentAction)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonBytes, &actionData); err != nil {
		return err
	}

	cpsReq := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerID,
		MakerName:        req.MakerName,
		MakerPhoneNumber: req.MakerPhoneNumber,
		RequestAction:    cps_const.RequestDisableFaydaAccount,
		Department:       req.Department,
		CurrentAction:    actionData,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err = f.FaydaDomain.InitiateDisableFaydaAccount(ctx, cpsReq)
	if err != nil {
		return err
	}

	return nil
}

func (f FaydaHandler) InitiateEnableFaydaAccount(ctx context.Context, req entities.CPSAction) error {
	var actionData ActionEnableData
	jsonBytes, err := json.Marshal(req.CurrentAction)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonBytes, &actionData); err != nil {
		return err
	}

	cpsReq := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerID,
		MakerName:        req.MakerName,
		MakerPhoneNumber: req.MakerPhoneNumber,
		RequestAction:    cps_const.RequestEnableFaydaAccount,
		Department:       req.Department,
		CurrentAction:    actionData,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err = f.FaydaDomain.InitiateDisableFaydaAccount(ctx, cpsReq)
	if err != nil {
		return err
	}

	return nil
}
