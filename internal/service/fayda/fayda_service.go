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

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableFayda", "Fayda", "EnableOrDisableFayda")
	defer span.End()

	f.logger.Infof("[FaydaSvc][EnableDisable] enabled: %v", isEnabled)

	existingUser, err := f.faydaRepo.FindByUserCode(ctx, user_code)
	if err != nil {
		f.logger.Errorf("[FaydaSvc][EnableDisable] fetch user err: %v", err)
		span.AddEvent("Failed to fetch fayda user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", user_code),
		))
		return err
	}

	if isEnabled && existingUser.Enabled {
		f.logger.Errorf("[FaydaSvc][EnableDisable] already enabled")
		span.AddEvent("Fayda user account already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorFaydaUserAccountEnabled.Code),
			attribute.String("user_code", user_code),
		))
		return errors.New(localization.ErrorFaydaUserAccountEnabled.Code)
	}

	if !isEnabled && !existingUser.Enabled {
		f.logger.Errorf("[FaydaSvc][EnableDisable] already disabled")
		span.AddEvent("Fayda user account already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorFaydaUserAccountDisabled.Code),
			attribute.String("user_code", user_code),
		))
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
		f.logger.Errorf("[FaydaSvc][EnableDisable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", user_code),
		))
		return err
	}

	f.logger.Infof("[FaydaSvc][EnableDisable] request created")
	return nil
}

func (f *faydaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Fayda", "Authorize")
	defer span.End()

	f.logger.Infof("[FaydaSvc][Authorize] action: %s", cpsAction.RequestAction)

	var faydaUser *member.User
	if err := local_util.BindAction(cpsAction.CurrentAction, &faydaUser); err != nil {
		f.logger.Errorf("[FaydaSvc][Authorize] bind err: %v", err)
		span.AddEvent("Failed to bind current action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch cpsAction.RequestAction {
	case string(cps_const.RequestEnableFaydaAccount):
		err := f.faydaRepo.Update(ctx, faydaUser, true)
		if err != nil {
			f.logger.Errorf("[FaydaSvc][Authorize] enable err: %v", err)
			span.AddEvent("Failed to enable fayda account", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		f.logger.Infof("[FaydaSvc][Authorize] enabled")
	case string(cps_const.RequestDisableFaydaAccount):
		err := f.faydaRepo.Update(ctx, faydaUser, false)
		if err != nil {
			f.logger.Errorf("[FaydaSvc][Authorize] disable err: %v", err)
			span.AddEvent("Failed to disable fayda account", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		f.logger.Infof("[FaydaSvc][Authorize] disabled")
	default:
		f.logger.Errorf("[FaydaSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = faydaUser
	f.logger.Infof("[FaydaSvc][Authorize] completed: %s", cpsAction.RequestAction)
	return cpsAction, nil
}
