package media

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type mediaService struct {
	repo   storage.ArticleRepository
	cache  storage.RedisRepository
	logger utils.Logger
}

func NewMediaService(repo storage.ArticleRepository, cache storage.RedisRepository, logger utils.Logger) service.ArticleService {
	return &mediaService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (m *mediaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Media", "Authorize")
	defer span.End()

	const (
		NewsArticleCacheKeyPattern   = "news:article:%s"
		NewsArticleCacheDeleteErrMsg = "failed to delete cache for article %s: %v"
	)
	m.logger.Infof("Media service authorizing action: %s", cpsAction.ActionCode)

	article, err := local_util.JsonUnmarshal[model.NewsArticleDetail](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateArticle):
		err = m.repo.CreateArticle(ctx, article.ToNewsArticle())
		if err != nil {
			span.AddEvent("Failed to create article", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateArticle):
		err = m.repo.UpdateArticle(ctx, article.ToNewsArticle(), cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to update article", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsArticleCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsArticleCacheDeleteErrMsg, article.ID, err)
		}
	case string(constants.RequestDeleteArticle):
		err = m.repo.DeleteArticle(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete article", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsArticleCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsArticleCacheDeleteErrMsg, article.ID, err)
		}
	case string(constants.RequestEnableArticle):
		err = m.repo.PublishUnpublishArticle(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable article", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsArticleCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsArticleCacheDeleteErrMsg, article.ID, err)
		}
	case string(constants.RequestDisableArticle):
		err = m.repo.PublishUnpublishArticle(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable article", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsArticleCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsArticleCacheDeleteErrMsg, article.ID, err)
		}
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported request action", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	cpsAction.CurrentAction = article
	m.logger.Infof("Action %s approved for Article %s", cpsAction.RequestAction, article.ID)

	return cpsAction, nil
}
