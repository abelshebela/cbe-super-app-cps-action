package miniapp

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
	m.logger.Infof("Mini app category service authorizing action: %s", cpsAction.ActionCode)

	category, err := local_util.JsonUnmarshal[model.MiniAppCategory](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniAppCategory):
		err = m.repo.Create(ctx, category)
	case string(constants.RequestUpdateMiniAppCategory):
		err = m.repo.Update(ctx, category, cpsAction.UniqueId)
	case string(constants.RequestDeleteMiniAppCategory):
		err = m.repo.Delete(ctx, cpsAction.UniqueId)
	case string(constants.RequestEnableMiniAppCategory):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableMiniAppCategory):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
	default:
		m.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		m.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.CurrentAction = category
	m.logger.Infof("Action %s approved for mini app Category %s", cpsAction.RequestAction, category.ID)
	return cpsAction, nil

}
