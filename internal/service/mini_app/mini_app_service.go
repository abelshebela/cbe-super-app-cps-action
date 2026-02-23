package miniapp

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type miniAppService struct {
	repo            storage.MiniAppRepository
	merchantService service.MiniAppMerchant
	logger          utils.Logger
}

func NewMiniAppService(repo storage.MiniAppRepository, merchantService service.MiniAppMerchant, logger utils.Logger) service.MiniAppService {
	return &miniAppService{
		repo:            repo,
		merchantService: merchantService,
		logger:          logger,
	}
}

func (s *miniAppService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "MiniApp", "Authorize")
	defer span.End()

	const financialAPPtype = "FINANCIAL"

	miniApp, err := local_util.JsonUnmarshal[mini_model.MiniApp](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniApp):
		if miniApp.AppType == financialAPPtype {
			err = s.ValidMerchant(miniApp.MerchantID.Hex(), ctx, s.merchantService)
			if err != nil {
				span.AddEvent("Merchant validation failed", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				break
			}

		}
		err = s.repo.Create(ctx, miniApp)
		if err != nil {
			span.AddEvent("Failed to create mini app", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateMiniApp):
		if miniApp.AppType == financialAPPtype {
			err = s.ValidMerchant(miniApp.MerchantID.Hex(), ctx, s.merchantService)
			if err != nil {
				span.AddEvent("Merchant validation failed", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				break
			}
		}
		err = s.repo.Update(ctx, cpsAction.UniqueId, miniApp)
		if err != nil {
			span.AddEvent("Failed to update mini app", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDeleteMiniApp):
		err = s.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete mini app", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestEnableMiniApp):
		if miniApp.AppType == financialAPPtype {
			err = s.ValidMerchant(miniApp.MerchantID.Hex(), ctx, s.merchantService)
			if err != nil {
				span.AddEvent("Merchant validation failed", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", cpsAction.UniqueId),
				))
				break
			}
		}
		err = s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable mini app", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDisableMiniApp):
		err = s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable mini app", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	default:
		s.logger.Errorf("[MiniAppSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported request action", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		s.logger.Errorf("[MiniAppSvc][Authorize] process err: %v", err)
		return nil, err
	}

	cpsAction.ActionStatus = "APPROVED"
	cpsAction.CurrentAction = miniApp
	s.logger.Infof("[MiniAppSvc][Authorize] approved action: %s app: %s", cpsAction.RequestAction, miniApp.AppName)

	return cpsAction, nil
}

func (s *miniAppService) ValidMerchant(MerchantID string, ctx context.Context, merchantService service.MiniAppMerchant) error {
	merchant, err := merchantService.FindByID(ctx, MerchantID)
	if err != nil {
		s.logger.Errorf("[MiniAppSvc][ValidMerchant] find err: %v", err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}
	if merchant.IsDeleted {
		s.logger.Errorf("[MiniAppSvc][ValidMerchant] deleted id: %s", MerchantID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	if !merchant.Enabled {
		s.logger.Errorf("[MiniAppSvc][ValidMerchant] disabled id: %s", MerchantID)
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}
	return nil
}
