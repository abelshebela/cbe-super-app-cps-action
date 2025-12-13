package miniapp

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type miniAppService struct {
	repo            storage.MiniAppRepository
	cpsService      service.CPSActionService
	merchantService service.MiniAppMerchantService
	logger          utils.Logger
}

func NewMiniAppService(repo storage.MiniAppRepository, cpsService service.CPSActionService, merchantService service.MiniAppMerchantService, logger utils.Logger) service.MiniAppService {
	return &miniAppService{
		repo:            repo,
		cpsService:      cpsService,
		merchantService: merchantService,
		logger:          logger,
	}
}

func (s *miniAppService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "MiniApp", "Authorize")
	defer span.End()

	s.logger.Infof("Authorize called, action: %s", cpsAction.RequestAction)

	miniApp, err := local_util.JsonUnmarshal[model.MiniApp](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniApp):
		err = s.ValidMerchant(miniApp.MerchantID.Hex(), ctx, s.merchantService)
		if err != nil {
			span.AddEvent("Merchant validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			break
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
		err = s.ValidMerchant(miniApp.MerchantID.Hex(), ctx, s.merchantService)
		if err != nil {
			span.AddEvent("Merchant validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			break
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
		err = s.ValidMerchant(miniApp.MerchantID.Hex(), ctx, s.merchantService)
		if err != nil {
			span.AddEvent("Merchant validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			break
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
		s.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported request action", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		s.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.ActionStatus = "APPROVED"
	cpsAction.CurrentAction = miniApp
	s.logger.Infof("Action %s approved for MiniApp %s", cpsAction.RequestAction, miniApp.AppName)

	return cpsAction, nil
}

func (s *miniAppService) ValidMerchant(MerchantID string, ctx context.Context, merchantService service.MiniAppMerchantService) error {
	merchant, err := merchantService.FindByID(ctx, MerchantID)
	if err != nil {
		s.logger.Errorf("Failed to get merchant details", "merchantID", MerchantID, "error", err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}
	if merchant.IsDeleted {
		s.logger.Errorf("Merchant is deleted", "merchantID", MerchantID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	if !merchant.Enabled {
		s.logger.Errorf("Merchant is disabled", "merchantID", MerchantID)
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}
	return nil
}
