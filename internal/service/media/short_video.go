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
)

type shortVideoService struct {
	repo   storage.ShortVideoRepository
	cache  storage.RedisRepository
	logger utils.Logger
}

func NewShortVideoService(repo storage.ShortVideoRepository, cache storage.RedisRepository, logger utils.Logger) service.ShortVideoService {
	return &shortVideoService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (m *shortVideoService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	m.logger.Infof("Media short video service authorizing action: %s", cpsAction.ActionCode)

	const (
		NewsShortVideoCacheKeyPattern   = "news:shortVideo:%s"
		NewsShortVideoCacheDeleteErrMsg = "failed to delete cache for shortVideo %s: %v"
	)

	shortVideo, err := local_util.JsonUnmarshal[model.ShortVideoDetail](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateShortVideo):
		err = m.repo.Create(ctx, shortVideo.ToShortVideo())
	case string(constants.RequestUpdateShortVideo):
		err = m.repo.Update(ctx, shortVideo.ToShortVideo(), cpsAction.UniqueId)

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	case string(constants.RequestDeleteShortVideo):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	case string(constants.RequestEnableShortVideo):
		err = m.repo.PublishUnpublish(ctx, cpsAction.UniqueId, true)

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	case string(constants.RequestDisableShortVideo):
		err = m.repo.PublishUnpublish(ctx, cpsAction.UniqueId, false)

		cacheKey := fmt.Sprintf(NewsShortVideoCacheKeyPattern, cpsAction.UniqueId)
		if err := m.cache.Delete(ctx, cacheKey); err != nil {
			m.logger.Warnf(NewsShortVideoCacheDeleteErrMsg, shortVideo.ID, err)
		}
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		m.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = shortVideo
	m.logger.Infof("Action %s approved for shortVideo shortVideo %s", cpsAction.RequestAction, shortVideo.ID)
	return cpsAction, nil

}
