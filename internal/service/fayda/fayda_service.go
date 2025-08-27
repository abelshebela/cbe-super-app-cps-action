package fayda

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	cps_const "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/service/fayda/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type faydaService struct {
	faydaRepo  storage.FaydaRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewFaydaService(faydaRepo storage.FaydaRepository, cpsService service.CPSActionService, logger utils.Logger) service.FaydaAccountService {
	return &faydaService{
		faydaRepo:  faydaRepo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (f *faydaService) EnableOrDisableFayda(ctx context.Context, user_code string, isEnabled bool) error {
	f.logger.Infof("Enabling fayda, user_code: %s", user_code)

	existingUser, err := f.faydaRepo.FindByUserCode(ctx, user_code)
	if err != nil {
		f.logger.Errorf("Failed to fetch fayda user by user_code: %s | Error: $%v", user_code, err)
		return err
	}

	if existingUser.KYCLevel != 1 {
		f.logger.Errorf("User with user_code: %s is not a fayda account user | Error: %v", user_code, err)
		return errors.New(localization.ErrorNotFaydaUser.Code)
	}

	if isEnabled && existingUser.Enabled {
		f.logger.Errorf("Fayda user acccount already enabled | user_code: %s", user_code)
		return errors.New(localization.ErrorFaydaUserAccountEnabled.Code)
	}

	if !isEnabled && !existingUser.Enabled {
		f.logger.Errorf("Fayda user acccount already disabled | user_code: %s", user_code)
		return errors.New(localization.ErrorFaydaUserAccountDisabled.Code)
	}

	currUser := *existingUser
	currUser.Enabled = isEnabled
	currUser.LastModifiedAt = time.Now()

	var requestAction constants.RequestAction
	if isEnabled {
		requestAction = constants.RequestEnableFaydaAccount
	} else {
		requestAction = constants.RequestDisableFaydaAccount
	}

	err = core.HandleCPSAction(ctx, f.cpsService, existingUser.ID.Hex(), requestAction, currUser, existingUser, constants.ActionUpdate)
	if err != nil {
		f.logger.Errorf("Failed to create CPS action for fayda account enable/disable | user_code: %s and Error: %v", user_code, err)
		return err
	}

	f.logger.Infof("Fayda Account with user_code: %s submitted to be enable/disable", user_code)
	return nil
}

func (f *faydaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	f.logger.Infof("Fayda service authorizing action: %s", cpsAction.ActionCode)

	var faydaUser *model.User
	if err := local_util.BindAction(cpsAction.CurrentAction, &faydaUser); err != nil {
		f.logger.Errorf("Failed to bind current action to fayda: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	var err error
	switch cpsAction.RequestAction {
	case string(cps_const.RequestEnableFaydaAccount):
		err = f.faydaRepo.AuthorizeEnableOrDisableFaydaUser(ctx, faydaUser.UserCode, true)
	case string(cps_const.RequestDisableFaydaAccount):
		err = f.faydaRepo.AuthorizeEnableOrDisableFaydaUser(ctx, faydaUser.UserCode, false)
	default:
		f.logger.Errorf("Unsupported action request: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	if err != nil {
		f.logger.Errorf("Failed to process fayda service with request action: %s and error: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = faydaUser
	f.logger.Infof("Fayda account authorization completed for request action: %s and user: %v", cpsAction.RequestAction, faydaUser.ID)
	return cpsAction, nil
}
