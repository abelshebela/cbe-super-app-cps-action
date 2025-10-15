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

type mediaService struct {
	repo   storage.ArticleRepository
	logger utils.Logger
}

func NewMediaService(repo storage.ArticleRepository, logger utils.Logger) service.ArticleService {
	return &mediaService{
		repo:   repo,
		logger: logger,
	}
}

func (m *mediaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	m.logger.Infof("Media service authorizing action: %s", cpsAction.ActionCode)

	article, err := local_util.JsonUnmarshal[model.NewsArticle](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateArticle):
		err = m.repo.CreateArticle(ctx, article)
	case string(constants.RequestUpdateArticle):
		err = m.repo.UpdateArticle(ctx, article)
	case string(constants.RequestDeleteArticle):
		err = m.repo.DeleteArticle(ctx, cpsAction.UniqueId)
	case string(constants.RequestEnableArticle):
		err = m.repo.PublishUnpublishArticle(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableArticle):
		err = m.repo.PublishUnpublishArticle(ctx, cpsAction.UniqueId, false)
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		m.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = article
	m.logger.Infof("Action %s approved for Article %s", cpsAction.RequestAction, article.ID)

	return cpsAction, nil
}
