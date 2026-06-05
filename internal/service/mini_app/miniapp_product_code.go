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

type miniappProductCodeService struct {
	repo   storage.MiniAppProductCodeRepository
	logger utils.Logger
}

func NewMiniAppProductCodeService(repo storage.MiniAppProductCodeRepository, logger utils.Logger) service.MiniappProductCodeService {
	return &miniappProductCodeService{
		repo:   repo,
		logger: logger,
	}
}

func (m *miniappProductCodeService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	log.Infof("[MiniProdCodeSvc][Authorize] action: %s", cpsAction.ActionCode)

	miniappProductCode, err := local_util.JsonUnmarshal[model.MiniAppProductCode](cpsAction.CurrentAction)
	if err != nil {
		return model.CPSAction{}, errors.New(localization.ErrorInvalidRequest.Code)
	}
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniappProductCode):
		err = m.repo.Create(ctx, miniappProductCode)
	case string(constants.RequestUpdateMiniappProductCode):
		err = m.repo.Update(ctx, miniappProductCode, cpsAction.UniqueId)
	case string(constants.RequestEnableMiniappProductCode):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
	case string(constants.RequestDisableMiniappProductCode):
		err = m.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
	default:
		log.Errorf("[MiniProdCodeSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return model.CPSAction{}, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		log.Errorf("[MiniProdCodeSvc][Authorize] process err: %v", err)
		return model.CPSAction{}, err
	}

	cpsAction.CurrentAction = miniappProductCode
	log.Infof("[MiniProdCodeSvc][Authorize] approved action: %s id: %s", cpsAction.RequestAction, miniappProductCode.ID)
	return *cpsAction, nil

}
