package fayda

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	cps_const "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/service/fayda/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

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
	f.logger.Infof("[EnableOrDisableFayda] processing fayda account, enabled: %v", isEnabled)

	existingUser, err := f.faydaRepo.FindByUserCode(ctx, user_code)
	if err != nil {
		f.logger.Errorf("[EnableOrDisableFayda] failed to fetch fayda user: %v", err)
		return err
	}

	if existingUser.KYCLevel != 1 {
		f.logger.Errorf("[EnableOrDisableFayda] user is not a fayda account user")
		return errors.New(localization.ErrorNotFaydaUser.Code)
	}

	if isEnabled && existingUser.Enabled {
		f.logger.Errorf("[EnableOrDisableFayda] fayda user account already enabled")
		return errors.New(localization.ErrorFaydaUserAccountEnabled.Code)
	}

	if !isEnabled && !existingUser.Enabled {
		f.logger.Errorf("[EnableOrDisableFayda] fayda user account already disabled")
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
		f.logger.Errorf("[EnableOrDisableFayda] failed to create CPS action: %v", err)
		return err
	}

	f.logger.Infof("[EnableOrDisableFayda] fayda account enable/disable request created successfully")
	return nil
}

func (f *faydaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	f.logger.Infof("[Authorize] authorizing fayda action: %s", cpsAction.RequestAction)

	var faydaUser *model.User
	if err := local_util.BindAction(cpsAction.CurrentAction, &faydaUser); err != nil {
		f.logger.Errorf("[Authorize] failed to bind current action to fayda: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch cpsAction.RequestAction {
	case string(cps_const.RequestEnableFaydaAccount):
		err := f.faydaRepo.Update(ctx, faydaUser, true)
		if err != nil {
			f.logger.Errorf("[Authorize] failed to enable fayda account: %v", err)
			return nil, err
		}
		f.logger.Infof("[Authorize] fayda account enabled successfully")
	case string(cps_const.RequestDisableFaydaAccount):
		err := f.faydaRepo.Update(ctx, faydaUser, false)
		if err != nil {
			f.logger.Errorf("[Authorize] failed to disable fayda account: %v", err)
			return nil, err
		}
		f.logger.Infof("[Authorize] fayda account disabled successfully")
	default:
		f.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = faydaUser
	f.logger.Infof("[Authorize] fayda account action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}
