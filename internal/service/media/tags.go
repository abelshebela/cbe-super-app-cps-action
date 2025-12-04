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
	m.logger.Infof("Media tags service authorizing action: %s", cpsAction.ActionCode)

	tags, err := local_util.JsonUnmarshal[model.NewsTags](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateNewsTag):
		err = m.repo.Create(ctx, tags)
	case string(constants.RequestUpdateNewsTag):
		err = m.repo.Update(ctx, tags, cpsAction.UniqueId)
	case string(constants.RequestDeleteNewsTag):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
	case string(constants.RequestEnableNewsTag):
		err = m.repo.EnableDisable(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableNewsTag):
		err = m.repo.EnableDisable(ctx, cpsAction.UniqueId, false)
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		m.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = tags
	m.logger.Infof("Action %s approved for Article Tags %s", cpsAction.RequestAction, tags.ID)
	return cpsAction, nil

}
