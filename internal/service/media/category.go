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

type mediaCategoryService struct {
	repo   storage.ArticleCategoryRepository
	logger utils.Logger
}

func NewMediaCategoryService(repo storage.ArticleCategoryRepository, logger utils.Logger) service.ArticleCategoryService {
	return &mediaCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (m *mediaCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Media", "Authorize")
	defer span.End()

	m.logger.Infof("Media category service authorizing action: %s", cpsAction.ActionCode)

	category, err := local_util.JsonUnmarshal[model.NewsCategoryModel](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateArticleCategory):
		err = m.repo.CreateArticleCategory(ctx, category)
		if err != nil {
			span.AddEvent("Failed to create article category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateArticleCategory):
		err = m.repo.UpdateArticleCategory(ctx, category, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to update article category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDeleteArticleCategory):
		err = m.repo.DeleteArticleCategory(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete article category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestEnableArticleCategory):
		err = m.repo.EnableOrDisableArticleCategory(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable article category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableArticleCategory):
		err = m.repo.EnableOrDisableArticleCategory(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable article category", trace.WithAttributes(
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
	m.logger.Infof("Action %s approved for Article Category %s", cpsAction.RequestAction, category.ID)
	return cpsAction, nil

}
