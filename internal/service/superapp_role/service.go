package superapprole

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	shared_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type superAppRoleService struct {
	repo       storage.SuperAppRoleRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewSuperAppRoleService(repo storage.SuperAppRoleRepository, cpsService service.CPSActionService, logger utils.Logger) *superAppRoleService {
	return &superAppRoleService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *superAppRoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.SuperAppRoleGroup], error) {
	return s.repo.FindAllWithPagination(ctx, filterParam)
}

func (s *superAppRoleService) EnableByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		map[string]interface{}{"superapp_role": superappRole},
		map[string]interface{}{"superapp_role": superappRole, "is_enabled": true},
		string(constants.RequestEnableSuperAppRole),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][EnableByRole] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) DisableByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		map[string]interface{}{"superapp_role": superappRole, "is_enabled": true},
		map[string]interface{}{"superapp_role": superappRole, "is_enabled": false},
		string(constants.RequestDisableSuperAppRole),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][DisableByRole] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) DeleteByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		map[string]interface{}{"superapp_role": superappRole},
		map[string]interface{}{"superapp_role": superappRole},
		string(constants.RequestDeleteSuperAppRole),
		constants.DELETE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][DeleteByRole] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) Authorize(ctx context.Context, action *shared_model.CPSAction) (*shared_model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[SuperAppRole][Authorize] action: %s, role: %s", action.RequestAction, action.UniqueId)

	role := action.UniqueId
	if role == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestEnableSuperAppRole):
		if err := s.repo.EnableByRole(ctx, role); err != nil {
			log.Errorf("[SuperAppRole][Authorize] enable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableSuperAppRole):
		if err := s.repo.DisableByRole(ctx, role); err != nil {
			log.Errorf("[SuperAppRole][Authorize] disable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDeleteSuperAppRole):
		if err := s.repo.DeleteByRole(ctx, role); err != nil {
			log.Errorf("[SuperAppRole][Authorize] delete err: %v", err)
			return nil, err
		}
	default:
		log.Errorf("[SuperAppRole][Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	return action, nil
}
