package media

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type mediaTagsService struct {
	repo   storage.NewsTagsRepository
	logger utils.Logger
}

func NewMediaTagsService(repo storage.NewsTagsRepository, logger utils.Logger) service.NewsTagsService {
	return &mediaTagsService{
		repo:   repo,
		logger: logger,
	}
}

func (m *mediaTagsService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Media", "Authorize")
	defer span.End()

	m.logger.Infof("Media tags service authorizing action: %s", cpsAction.ActionCode)

	tags, err := local_util.JsonUnmarshal[model.NewsTags](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateNewsTag):
		err = m.repo.Create(ctx, tags)
		if err != nil {
			span.AddEvent("Failed to create news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateNewsTag):
		err = m.repo.Update(ctx, tags, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to update news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDeleteNewsTag):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestEnableNewsTag):
		err = m.repo.EnableDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableNewsTag):
		err = m.repo.EnableDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable news tag", trace.WithAttributes(
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

	cpsAction.CurrentAction = tags
	m.logger.Infof("Action %s approved for Article Tags %s", cpsAction.RequestAction, tags.ID)
	return cpsAction, nil

}
