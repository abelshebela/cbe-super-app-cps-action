package cpsroles

import (
	"cbe-super-app-cps-action/internal/constants"
	cps_role_dto "cbe-super-app-cps-action/internal/constants/dto/cps_roles"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type cpsRoleService struct {
	repo       storage.CPSRolesRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewCPSRoleService(repo storage.CPSRolesRepository, cpsService service.CPSActionService, logger utils.Logger) service.CPSRolesService {
	return &cpsRoleService{repo: repo, cpsService: cpsService, logger: logger}
}

func (r *cpsRoleService) Create(ctx context.Context, req cps_role_dto.CreateCPSRoleRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	enabled := true
	role := model.CPSRoles{
		Name:      req.Name,
		Enabled:   &enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, role, string(constants.RequestCreateCpsRole), constants.CREATE)

	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		r.logger.Errorf("[Create] failed to create CPS action: %v", err)
		return err
	}

	r.logger.Infof("[Create] cps role creation request created successfully")
	return nil
}

func (r *cpsRoleService) Update(ctx context.Context, id string, req cps_role_dto.UpdateCPSRoleRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := r.repo.FindById(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("[Update] Resource not found")
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Update] failed to find existing cps role: %v", err)
		return err
	}

	updated := *existing
	if req.Name != "" {
		updated.Name = req.Name
	}
	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, string(constants.RequestUpdateCpsRole), constants.UPDATE)

	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		r.logger.Errorf("[Update] failed to create CPS action: %v", err)
		return err
	}

	r.logger.Infof("[Update] cps role update request created successfully")
	return nil
}

func (r *cpsRoleService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := r.repo.FindById(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("[EnableOrDisable] Resource not found")
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[Enable] failed to find existing cps role: %v", err)
		return err
	}

	if enable && *existing.Enabled {
		r.logger.Warnf("CPS Role already enabled, id: %s", id)
		return errors.New(localization.ErrorCPSRoleAlreadyEnabled.Code)
	}
	if !enable && !*existing.Enabled {
		r.logger.Warnf("CPS Role already disabled, id: %s", id)
		return errors.New(localization.ErrorCPSRoleAlreadyDisabled.Code)
	}

	updated := *existing
	updated.UpdatedAt = time.Now()
	updated.Enabled = &enable

	var requestType string
	if enable {
		requestType = string(constants.RequestEnableCpsRole)
	} else {
		requestType = string(constants.RequestDisableCpsRole)
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, requestType, constants.UPDATE)

	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		r.logger.Errorf("[Enable] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (r *cpsRoleService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.CPSRoles], error) {
	return r.repo.FindAllWithPagination(ctx, filterParam)
}

func (r *cpsRoleService) FindById(ctx context.Context, id string) (*model.CPSRoles, error) {
	return r.repo.FindById(ctx, id)
}

func (r *cpsRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	r.logger.Infof("[Authorize] authorizing cps role action: %s", action.RequestAction)

	var err error
	role, marshal_err := local_util.JsonUnmarshal[model.CPSRoles](action.CurrentAction)
	if marshal_err != nil || role == nil {
		r.logger.Errorf("[Authorize] failed to unmarshal current action: %v", marshal_err)
		return nil, marshal_err
	}

	switch action.RequestAction {
	case string(constants.RequestCreateCpsRole):
		err = r.repo.Create(ctx, *role)
	case string(constants.RequestUpdateCpsRole):
		err = r.repo.Update(ctx, action.UniqueId, *role)
	case string(constants.RequestEnableCpsRole):
		err = r.repo.EnableOrDisable(ctx, action.UniqueId, true)
	case string(constants.RequestDisableCpsRole):
		err = r.repo.EnableOrDisable(ctx, action.UniqueId, false)
	default:
		r.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New("unsupported action")
	}

	if err != nil {
		r.logger.Errorf("[Authorize] operation failed: %v", err)
		return nil, err
	}

	r.logger.Infof("[Authorize] cps role action %s completed successfully", action.RequestAction)
	return action, nil
}
