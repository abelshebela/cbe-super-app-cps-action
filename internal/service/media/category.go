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
	m.logger.Infof("Media category service authorizing action: %s", cpsAction.ActionCode)

	category, err := local_util.JsonUnmarshal[model.NewsCategoryModel](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateArticleCategory):
		err = m.repo.CreateArticleCategory(ctx, category)
	case string(constants.RequestUpdateArticleCategory):
		err = m.repo.UpdateArticleCategory(ctx, category, cpsAction.UniqueId)
	case string(constants.RequestDeleteArticleCategory):
		err = m.repo.DeleteArticleCategory(ctx, cpsAction.UniqueId)
	case string(constants.RequestEnableArticleCategory):
		err = m.repo.EnableOrDisableArticleCategory(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableArticleCategory):
		err = m.repo.EnableOrDisableArticleCategory(ctx, cpsAction.UniqueId, false)
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		m.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = category
	m.logger.Infof("Action %s approved for Article Category %s", cpsAction.RequestAction, category.ID)
	return cpsAction, nil

}
