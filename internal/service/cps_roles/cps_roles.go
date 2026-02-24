package cpsroles

import (
	"cbe-super-app-cps-action/internal/constants"
	cps_role_dto "cbe-super-app-cps-action/internal/constants/dto/cps_roles"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

	cpsRole, _ := r.repo.FindByNameOrRoleCode(ctx, req.Name, req.RoleCode)
	if cpsRole != nil {
		if cpsRole.Name == req.Name {
			r.logger.Warnf("[CpsRoleSvc][Create] name exists: %s", req.Name)
			return errors.New(localization.ErrorCPSRoleNameAlreadyExists.Code)
		}
		if cpsRole.RoleCode == req.RoleCode {
			r.logger.Warnf("[CpsRoleSvc][Create] code exists: %s", req.RoleCode)
			return errors.New(localization.ErrorCPSRoleCodeAlreadyExists.Code)
		}
	}

	enabled := true
	role := imodel.CPSRoles{
		Name:        req.Name,
		RoleCode:    req.RoleCode,
		Description: req.Description,
		Enabled:     &enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, role, string(constants.RequestCreateCpsRole), constants.CREATE)

	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		r.logger.Errorf("[CpsRoleSvc][Create] cps action err: %v", err)
		return err
	}

	r.logger.Infof("[CpsRoleSvc][Create] request created")
	return nil
}

func (r *cpsRoleService) Update(ctx context.Context, id string, req cps_role_dto.UpdateCPSRoleRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := r.repo.FindById(ctx, id)
	if err != nil {
		r.logger.Errorf("[CpsRoleSvc][Update] find err: %v", err)
		return err
	}

	var name, roleCode string
	if req.Name != nil {
		name = *req.Name
	}
	if req.RoleCode != nil {
		roleCode = *req.RoleCode
	}

	cpsRole, err := r.repo.FindByNameOrRoleCode(ctx, name, roleCode)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		r.logger.Errorf("[CpsRoleSvc][Update] check existing err: %v", err)
		return err
	}

	if cpsRole != nil && cpsRole.ID.Hex() != id {
		if strings.EqualFold(cpsRole.Name, name) && existing.Name != name {
			r.logger.Warnf("[CpsRoleSvc][Update] name exists: %s", req.Name)
			return errors.New(localization.ErrorCPSRoleNameAlreadyExists.Code)
		}
		if cpsRole.RoleCode == roleCode && existing.RoleCode != roleCode {
			r.logger.Warnf("[CpsRoleSvc][Update] code exists: %s", req.RoleCode)
			return errors.New(localization.ErrorCPSRoleCodeAlreadyExists.Code)
		}
	}

	updated := *existing
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.RoleCode != nil {
		updated.RoleCode = *req.RoleCode
	}
	if req.Description != nil {
		updated.Description = *req.Description
	}
	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, string(constants.RequestUpdateCpsRole), constants.UPDATE)

	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		r.logger.Errorf("[CpsRoleSvc][Update] cps action err: %v", err)
		return err
	}

	r.logger.Infof("[CpsRoleSvc][Update] request created")
	return nil
}

func (r *cpsRoleService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := r.repo.FindById(ctx, id)
	if err != nil {
		r.logger.Errorf("[CpsRoleSvc][EnableDisable] find err: %v", err)
		return err
	}

	if enable && *existing.Enabled {
		r.logger.Warnf("[CpsRoleSvc][EnableDisable] already enabled id: %s", id)
		return err
	}
	if !enable && !*existing.Enabled {
		r.logger.Warnf("[CpsRoleSvc][EnableDisable] already disabled id: %s", id)
		return err
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
		r.logger.Errorf("[CpsRoleSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

func (r *cpsRoleService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CPSRoles], error) {
	return r.repo.FindAllWithPagination(ctx, filterParam)
}

func (r *cpsRoleService) FindById(ctx context.Context, id string) (*imodel.CPSRoles, error) {
	return r.repo.FindById(ctx, id)
}

func (r *cpsRoleService) EnableServiceAccess(ctx context.Context, roleID string, req cps_role_dto.ToggleServiceAccessRequest) error {
	if err := r.repo.EnableServiceAccess(ctx, roleID, req.AccessListKeys); err != nil {
		r.logger.Errorf("[CpsRoleSvc][EnableAccess] err for role %s: %v", roleID, err)
		return err
	}
	r.logger.Infof("[CpsRoleSvc][EnableAccess] enabled role: %s, keys: %v", roleID, req.AccessListKeys)
	return nil
}

func (r *cpsRoleService) DisableServiceAccess(ctx context.Context, roleID string, req cps_role_dto.ToggleServiceAccessRequest) error {
	if err := r.repo.DisableServiceAccess(ctx, roleID, req.AccessListKeys); err != nil {
		r.logger.Errorf("[CpsRoleSvc][DisableAccess] err for role %s: %v", roleID, err)
		return err
	}
	r.logger.Infof("[CpsRoleSvc][DisableAccess] disabled role: %s, keys: %v", roleID, req.AccessListKeys)
	return nil
}

func (r *cpsRoleService) Delete(ctx context.Context, id string) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := r.repo.FindById(ctx, id)
	if err != nil {
		r.logger.Errorf("[CpsRoleSvc][Delete] find err: %v", err)
		return err
	}

	updated := *existing
	updated.DeletedAt = time.Now()

	requestType := string(constants.RequestDeleteCpsRole)

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, requestType, constants.DELETE)

	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		r.logger.Errorf("[CpsRoleSvc][Delete] cps action err: %v", err)
		return err
	}

	return nil
}

func (r *cpsRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	r.logger.Infof("[CpsRoleSvc][Authorize] action: %s", action.RequestAction)

	var err error
	role, marshal_err := local_util.JsonUnmarshal[imodel.CPSRoles](action.CurrentAction)
	if marshal_err != nil || role == nil {
		r.logger.Errorf("[CpsRoleSvc][Authorize] unmarshal err: %v", marshal_err)
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
	case string(constants.RequestDeleteCpsRole):
		err = r.repo.Delete(ctx, action.UniqueId)
	default:
		r.logger.Errorf("[CpsRoleSvc][Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New("unsupported action")
	}

	if err != nil {
		r.logger.Errorf("[CpsRoleSvc][Authorize] operation err: %v", err)
		return nil, err
	}

	r.logger.Infof("[CpsRoleSvc][Authorize] completed: %s", action.RequestAction)
	return action, nil
}
