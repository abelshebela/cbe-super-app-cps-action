package miniapp

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"

	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type miniAppMerchantService struct {
	repo           storage.MiniAppMerchant
	miniRepo       storage.MiniAppRepository
	logger         utils.Logger
	merchantLookup merchant_lookup.MerchantLookupAdapter
}

func NewMiniAppMerchantService(
	repo storage.MiniAppMerchant,
	miniRepo storage.MiniAppRepository,
	merchantLookup merchant_lookup.MerchantLookupAdapter,
	logger utils.Logger,
) service.MiniAppMerchant {
	return &miniAppMerchantService{
		repo:           repo,
		miniRepo:       miniRepo,
		merchantLookup: merchantLookup,
		logger:         logger,
	}
}

func (m *miniAppMerchantService) FindByID(ctx context.Context, id string) (*mini_model.MiniAppMerchant, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "MiniAppMerchant", "FindByID")
	defer span.End()

	m.logger.Infof("[MiniMerchSvc][FindByID] id: %s", id)
	result, err := m.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find mini app merchant", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

func (m *miniAppMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "MiniAppMerchant", "Authorize")
	defer span.End()

	merchant, err := local_util.JsonUnmarshal[mini_model.MiniAppMerchant](cpsAction.CurrentAction)
	if err != nil {
		m.logger.Errorf("[MiniMerchSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniAppMerchant):
		_, err = m.repo.Create(ctx, merchant)
		if err != nil {
			span.AddEvent("Failed to create mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateMiniAppMerchant):
		err = m.repo.Update(ctx, cpsAction.UniqueId, merchant)
		if err != nil {
			span.AddEvent("Failed to update mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDeleteMiniAppMerchant):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		if err == nil {
			err := m.miniRepo.DeleteManyByMerchantIDs(ctx, cpsAction.UniqueId)
			if err != nil {
				m.logger.Errorf("[MiniMerchSvc][Authorize] cascade delete err id: %s: %v", cpsAction.UniqueId, err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		}
	case string(constants.RequestEnableMiniAppMerchant):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableMiniAppMerchant):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable mini app merchant", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		if err == nil {
			err := m.miniRepo.DisableManyByMerchantIDs(ctx, cpsAction.UniqueId)
			if err != nil {
				m.logger.Errorf("[MiniMerchSvc][Authorize] cascade disable err id: %s: %v", cpsAction.UniqueId, err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
		}
	default:
		m.logger.Errorf("[MiniMerchSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	if err != nil {
		m.logger.Errorf("[MiniMerchSvc][Authorize] process err: %v", err)
		return nil, err
	}

	cpsAction.CurrentAction = merchant
	m.logger.Infof("[MiniMerchSvc][Authorize] completed action: %s id: %s", cpsAction.RequestAction, merchant.ID)
	return cpsAction, nil
}
