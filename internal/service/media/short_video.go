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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type shortVideoService struct {
	repo   storage.ShortVideoRepository
	logger utils.Logger
}

func NewShortVideoService(repo storage.ShortVideoRepository, logger utils.Logger) service.ShortVideoService {
	return &shortVideoService{
		repo:   repo,
		logger: logger,
	}
}

func (m *shortVideoService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	m.logger.Infof("Media short video service authorizing action: %s", cpsAction.ActionCode)

	shortVideo, err := local_util.JsonUnmarshal[model.ShortVideo](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateShortVideo):
		err = m.repo.Create(ctx, shortVideo)
	case string(constants.RequestUpdateShortVideo):
		err = m.repo.Update(ctx, shortVideo, cpsAction.UniqueId)
	case string(constants.RequestDeleteShortVideo):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
	case string(constants.RequestEnableShortVideo):
		err = m.repo.PublishUnpublish(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableShortVideo):
		err = m.repo.PublishUnpublish(ctx, cpsAction.UniqueId, false)
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		m.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = shortVideo
	m.logger.Infof("Action %s approved for Article shortVideo %s", cpsAction.RequestAction, shortVideo.ID)
	return cpsAction, nil

}
