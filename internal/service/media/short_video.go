package media

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type shortVideoService struct {
	repo     storage.ShortVideoRepository
	cache    storage.RedisRepository
	producer KafkaProducerService
	logger   utils.Logger
}

func NewShortVideoService(repo storage.ShortVideoRepository, cache storage.RedisRepository, producer KafkaProducerService, logger utils.Logger) service.ShortVideoService {
	return &shortVideoService{
		repo:     repo,
		cache:    cache,
		producer: producer,
		logger:   logger,
	}
}

func (m *shortVideoService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Media", "Authorize")
	defer span.End()

	m.logger.Infof("[ShortVideoSvc][Authorize] action: %s", cpsAction.ActionCode)

	const (
		NewsShortVideoCacheKeyPattern   = "news:shortVideo:%s"
		NewsShortVideoCacheDeleteErrMsg = "failed to delete cache for shortVideo %s: %v"
	)

	shortVideo, err := local_util.JsonUnmarshal[model.ShortVideoDetail](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateShortVideo):
		createdVideo, err := m.repo.Create(ctx, shortVideo.ToShortVideo())
		if err != nil {
			span.AddEvent("Failed to create short video", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		shortVideoID := createdVideo.ID.Hex()
		if err := m.producer.PublishShortVideoEvent(shortVideo.VideoURL, shortVideoID); err != nil {
			m.logger.Errorf("[ShortVideoSvc][Authorize] publish err: %v", err)
		} else {
			m.logger.Infof("[ShortVideoSvc][Authorize] published id: %s", shortVideoID)
		}
		shortVideo.ID = createdVideo.ID
	case string(constants.RequestUpdateShortVideo):
		err = m.repo.Update(ctx, shortVideo.ToShortVideo(), cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to update short video", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		prevVideo, err := local_util.JsonUnmarshal[model.ShortVideoDetail](cpsAction.PreviousAction)
		if err != nil {
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if shortVideo.VideoURL != prevVideo.VideoURL {
			shortVideoID := shortVideo.ID.Hex()
			if shortVideo.ID.IsZero() && cpsAction.UniqueId != "" {
				shortVideoID = cpsAction.UniqueId
			}
			if err := m.producer.PublishShortVideoEvent(shortVideo.VideoURL, shortVideoID); err != nil {
				m.logger.Errorf("[ShortVideoSvc][Authorize] publish err: %v", err)
			} else {
				m.logger.Infof("[ShortVideoSvc][Authorize] published id: %s", shortVideoID)
			}
		}

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	case string(constants.RequestDeleteShortVideo):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete short video", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	case string(constants.RequestEnableShortVideo):
		err = m.repo.PublishUnpublish(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("Failed to enable short video", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	case string(constants.RequestDisableShortVideo):
		err = m.repo.PublishUnpublish(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("Failed to disable short video", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	default:
		m.logger.Errorf("[ShortVideoSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported request action", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	cpsAction.CurrentAction = shortVideo
	m.logger.Infof("[ShortVideoSvc][Authorize] approved action: %s id: %s", cpsAction.RequestAction, shortVideo.ID)
	return cpsAction, nil

}
