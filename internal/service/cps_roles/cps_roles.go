package cpsroles

// import (
// 	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
// 	cps_role_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/cps_roles"
// 	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
// 	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
// 	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
// 	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
// 	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
// 	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
// 	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
// 	"context"
// 	"errors"
// 	"strings"
// 	"time"

// 	"github.com/hugokessem/coreio/core"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// )

// type cpsRoleService struct {
// 	repo       storage.CPSRolesRepository
// 	cpsService service.CPSActionService
// 	core       core.CBECoreAPIInterface
// 	logger     utils.Logger
// }

// func NewCPSRoleService(repo storage.CPSRolesRepository, cpsService service.CPSActionService, coreInterface core.CBECoreAPIInterface, logger utils.Logger) service.CPSRolesService {
// 	return &cpsRoleService{repo: repo, cpsService: cpsService, core: coreInterface, logger: logger}
// }

// func (r *cpsRoleService) Create(ctx context.Context, req cps_role_dto.CreateCPSRoleRequest) error {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)
// 	makerUser := local_util.ExtractUserFromContext(ctx)

// 	cpsRole, _ := r.repo.FindByNameOrRoleCode(ctx, req.Name, req.RoleCode)
// 	if cpsRole != nil {
// 		if strings.EqualFold(cpsRole.Name, req.Name) {
// 			log.Warnf("[CpsRoleSvc][Create] name exists: %s", req.Name)
// 			return errors.New(localization.ErrorCPSRoleNameAlreadyExists.Code)
// 		}
// 		if strings.EqualFold(cpsRole.RoleCode, req.RoleCode) {
// 			log.Warnf("[CpsRoleSvc][Create] code exists: %s", req.RoleCode)
// 			return errors.New(localization.ErrorCPSRoleCodeAlreadyExists.Code)
// 		}
// 	}

// 	enabled := true
// 	role := imodel.CPSRoles{
// 		Name:        req.Name,
// 		RoleCode:    req.RoleCode,
// 		Lable:       req.Lable,
// 		Description: req.Description,
// 		Enabled:     &enabled,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, role, string(constants.RequestCreateCpsRole), constants.CREATE)

// 	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
// 		log.Errorf("[CpsRoleSvc][Create] cps action err: %v", err)
// 		return err
// 	}

// 	log.Infof("[CpsRoleSvc][Create] request created")
// 	return nil
// }

// func (r *cpsRoleService) Update(ctx context.Context, id string, req cps_role_dto.UpdateCPSRoleRequest) error {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	makerUser := local_util.ExtractUserFromContext(ctx)

// 	existing, err := r.repo.FindById(ctx, id)
// 	if err != nil {
// 		log.Errorf("[CpsRoleSvc][Update] find err: %v", err)
// 		return err
// 	}

// 	var name, roleCode, label string
// 	if req.Name != nil {
// 		name = *req.Name
// 	}
// 	if req.RoleCode != nil {
// 		roleCode = *req.RoleCode
// 	}
// 	if req.Lable != nil {
// 		label = *req.Lable
// 	}
// 	cpsRole, err := r.repo.FindByNameOrRoleCode(ctx, name, roleCode)
// 	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
// 		log.Errorf("[CpsRoleSvc][Update] check existing err: %v", err)
// 		return err
// 	}

// 	if cpsRole != nil && cpsRole.ID != id {
// 		if strings.EqualFold(cpsRole.Name, name) && existing.Name != name {
// 			log.Warnf("[CpsRoleSvc][Update] name exists: %s", req.Name)
// 			return errors.New(localization.ErrorCPSRoleNameAlreadyExists.Code)
// 		}
// 		if cpsRole.RoleCode == roleCode && existing.RoleCode != roleCode {
// 			log.Warnf("[CpsRoleSvc][Update] code exists: %s", req.RoleCode)
// 			return errors.New(localization.ErrorCPSRoleCodeAlreadyExists.Code)
// 		}
// 		if cpsRole.Lable == label && existing.Lable != label {
// 			log.Warnf("[CpsRoleSvc][Update] label exists: %s", req.Lable)
// 			return errors.New(localization.ErrorCPSRoleLabelAlreadyExists.Code)
// 		}
// 	}

// 	updated := *existing
// 	if req.Name != nil {
// 		updated.Name = *req.Name
// 	}
// 	if req.RoleCode != nil {
// 		updated.RoleCode = *req.RoleCode
// 	}
// 	if req.Description != nil {
// 		updated.Description = *req.Description
// 	}
// 	if req.Lable != nil {
// 		updated.Lable = *req.Lable
// 	}
// 	updated.UpdatedAt = time.Now()

// 	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, string(constants.RequestUpdateCpsRole), constants.UPDATE)

// 	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
// 		log.Errorf("[CpsRoleSvc][Update] cps action err: %v", err)
// 		return err
// 	}

// 	log.Infof("[CpsRoleSvc][Update] request created")
// 	return nil
// }

// func (r *cpsRoleService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	makerUser := local_util.ExtractUserFromContext(ctx)

// 	existing, err := r.repo.FindById(ctx, id)
// 	if err != nil {
// 		log.Errorf("[CpsRoleSvc][EnableDisable] find err: %v", err)
// 		return err
// 	}

// 	if enable && *existing.Enabled {
// 		log.Warnf("[CpsRoleSvc][EnableDisable] already enabled id: %s", id)
// 		return errors.New(localization.ErrorCPSRoleAlreadyEnabled.Code)
// 	}
// 	if !enable && !*existing.Enabled {
// 		log.Warnf("[CpsRoleSvc][EnableDisable] already disabled id: %s", id)
// 		return errors.New(localization.ErrorCPSRoleAlreadyDisabled.Code)
// 	}

// 	updated := *existing
// 	updated.UpdatedAt = time.Now()
// 	updated.Enabled = &enable

// 	var requestType string
// 	if enable {
// 		requestType = string(constants.RequestEnableCpsRole)
// 	} else {
// 		requestType = string(constants.RequestDisableCpsRole)
// 	}

// 	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, requestType, constants.UPDATE)

// 	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
// 		log.Errorf("[CpsRoleSvc][EnableDisable] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (r *cpsRoleService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CPSRoles], error) {
// 	return r.repo.FindAllWithPagination(ctx, filterParam)
// }

// func (r *cpsRoleService) FindById(ctx context.Context, id string) (*imodel.CPSRoles, error) {
// 	return r.repo.FindById(ctx, id)
// }

// func (r *cpsRoleService) EnableServiceAccess(ctx context.Context, roleID string, req cps_role_dto.ToggleServiceAccessRequest) error {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	if err := r.repo.EnableServiceAccess(ctx, roleID, req.AccessListKeys); err != nil {
// 		log.Errorf("[CpsRoleSvc][EnableAccess] err for role %s: %v", roleID, err)
// 		return err
// 	}
// 	log.Infof("[CpsRoleSvc][EnableAccess] enabled role: %s, keys: %v", roleID, req.AccessListKeys)
// 	return nil
// }

// func (r *cpsRoleService) DisableServiceAccess(ctx context.Context, roleID string, req cps_role_dto.ToggleServiceAccessRequest) error {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	if err := r.repo.DisableServiceAccess(ctx, roleID, req.AccessListKeys); err != nil {
// 		log.Errorf("[CpsRoleSvc][DisableAccess] err for role %s: %v", roleID, err)
// 		return err
// 	}
// 	log.Infof("[CpsRoleSvc][DisableAccess] disabled role: %s, keys: %v", roleID, req.AccessListKeys)
// 	return nil
// }

// func (r *cpsRoleService) Delete(ctx context.Context, id string) error {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	makerUser := local_util.ExtractUserFromContext(ctx)

// 	existing, err := r.repo.FindById(ctx, id)
// 	if err != nil {
// 		log.Errorf("[CpsRoleSvc][Delete] find err: %v", err)
// 		return err
// 	}

// 	err = r.repo.CheckUserExistence(ctx, existing.RoleCode)
// 	if err != nil {
// 		log.Errorf("[CpsRoleSvc][Delete] check user existence err: %v", err)
// 		return err
// 	}

// 	// Check Access list by superapp role
// 	sar, err := r.repo.FindSupperAppRoleByAccessList(ctx, id)
// 	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
// 		return err
// 	}
// 	if sar {
// 		return errors.New("There is an active customer segmentation with this access list")
// 	}

// 	updated := *existing
// 	now := time.Now()
// 	updated.IsDeleted = true
// 	updated.DeletedAt = now
// 	updated.UpdatedAt = now

// 	requestType := string(constants.RequestDeleteCpsRole)

// 	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, requestType, constants.DELETE)

// 	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
// 		log.Errorf("[CpsRoleSvc][Delete] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (r *cpsRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	log.Infof("[CpsRoleSvc][Authorize] action: %s", action.RequestAction)

// 	var err error

// 	switch action.RequestAction {
// 	case string(constants.RequestDeleteCpsRole):
// 		err = r.repo.Delete(ctx, action.UniqueId)
// 		if err != nil {
// 			log.Errorf("[CpsRoleSvc][Authorize] operation err: %v", err)
// 			return nil, err
// 		}
// 		log.Infof("[CpsRoleSvc][Authorize] completed: %s", action.RequestAction)
// 		return action, nil
// 	}

// 	role, marshal_err := local_util.JsonUnmarshal[imodel.CPSRoles](action.CurrentAction)
// 	if marshal_err != nil || role == nil {
// 		log.Errorf("[CpsRoleSvc][Authorize] unmarshal err: %v", marshal_err)
// 		return nil, marshal_err
// 	}

// 	switch action.RequestAction {
// 	case string(constants.RequestCreateCpsRole):
// 		err = r.repo.Create(ctx, *role)
// 	case string(constants.RequestUpdateCpsRole):
// 		err = r.repo.Update(ctx, action.UniqueId, *role)
// 	case string(constants.RequestEnableCpsRole):
// 		err = r.repo.EnableOrDisable(ctx, action.UniqueId, true)
// 	case string(constants.RequestDisableCpsRole):
// 		err = r.repo.EnableOrDisable(ctx, action.UniqueId, false)
// 	default:
// 		log.Errorf("[CpsRoleSvc][Authorize] unsupported action: %s", action.RequestAction)
// 		return nil, errors.New("unsupported action")
// 	}

// 	if err != nil {
// 		log.Errorf("[CpsRoleSvc][Authorize] operation err: %v", err)
// 		return nil, err
// 	}

// 	log.Infof("[CpsRoleSvc][Authorize] completed: %s", action.RequestAction)
// 	return action, nil
// }
