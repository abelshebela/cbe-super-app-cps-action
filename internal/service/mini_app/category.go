package miniapp

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type mediaCategoryService struct {
	repo   storage.MiniAppCategoryRepository
	logger utils.Logger
}

func NewMiniAppCategoryService(repo storage.MiniAppCategoryRepository, logger utils.Logger) service.MiniAppCategoryService {
	return &mediaCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (m *mediaCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "MiniApp", "Authorize")
	defer span.End()

	m.logger.Infof("Mini app category service authorizing action: %s", cpsAction.ActionCode)

	category, err := local_util.JsonUnmarshal[model.MiniAppCategory](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniAppCategory):
		err = m.repo.Create(ctx, category)
		if err != nil {
			span.AddEvent("Failed to create mini app category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateMiniAppCategory):
		err = m.repo.Update(ctx, category, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to update mini app category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDeleteMiniAppCategory):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete mini app category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestEnableMiniAppCategory):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable mini app category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableMiniAppCategory):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable mini app category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported request action", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	cpsAction.CurrentAction = category
	m.logger.Infof("Action %s approved for mini app Category %s", cpsAction.RequestAction, category.ID)
	return cpsAction, nil

}
