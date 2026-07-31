package roles

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service/roles/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleService struct {
	cpsService           service.CPSActionService
	portalCardRepo       storage.PortalCardRepository
	roleRepository       storage.RoleRepository
	approveIndexRepo     storage.CPSActionApproveIndexRepository
	jobRepo              storage.JobRoleRepository
	bpsApproverIndexRepo storage.BPSActionApproveIndexRepository
	cfg                  config.VaultConfig
	logger               utils.Logger
}

func NewRoleService(roleRepo storage.RoleRepository, portalCard storage.PortalCardRepository, approveIndexRepo storage.CPSActionApproveIndexRepository, jobRepo storage.JobRoleRepository, bpsAprroverIndexRepo storage.BPSActionApproveIndexRepository, cpsService service.CPSActionService, cfg config.VaultConfig, logger utils.Logger) service.RoleService {
	return &RoleService{
		cpsService:           cpsService,
		portalCardRepo:       portalCard,
		roleRepository:       roleRepo,
		approveIndexRepo:     approveIndexRepo,
		jobRepo:              jobRepo,
		bpsApproverIndexRepo: bpsAprroverIndexRepo,
		cfg:                  cfg,
		logger:               logger,
	}
}

func (j *RoleService) Create(ctx context.Context, role imodel.Role) error {
	log := local_util.LoggerFromCtx(ctx, j.logger)

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if err := core.RoleExistenChecker(ctx, constants.CREATE, "", role, j.roleRepository); err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return errors.New(err.Error())
		}
	}
	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, role, constants.RequestCreateRole, constants.CREATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *RoleService) FindAll(ctx context.Context) (*[]imodel.Role, error) {
	return j.roleRepository.FindAll(ctx)
}

func (j *RoleService) Update(ctx context.Context, id string, update imodel.Role) error {
	log := local_util.LoggerFromCtx(ctx, j.logger)

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	log.Infof("[Role Service][Update] Update role existence check started--------------------")
	if err := core.RoleExistenChecker(ctx, constants.UPDATE, id, update, j.roleRepository); err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	log.Infof("[Role Service][Update] Update role existence check existing pass--------------------")

	prev, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	log.Infof("[Role Service][Update] Update role existence check existing prev data--------------------")

	log.Infof("[Role][Update] Update role existence check passed")
	newRole := *prev
	if update.Name != "" {
		newRole.Name = update.Name
	}

	if update.Code != "" {
		newRole.Code = update.Code
	}
	if update.Type != "" {
		newRole.Type = update.Type
	}
	if update.Description != "" {
		newRole.Description = update.Description
	}
	newRole.UpdatedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, newRole, constants.RequestUpdateRole, constants.UPDATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *RoleService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, j.logger)

	makerUser := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[Role Service][EnableOrDisable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[Role Service][EnableOrDisable] failed to find existing role: %v", err)
		return err
	}

	if !enable && makerUser.UserRole == existing.Code {
		log.Errorf("[Role Service][EnableOrDisable] user cannot disable their own role")
		return errors.New(localization.ErrorCannotDisableOwnRole.Code)
	}

	if existing.Enable == enable && enable {
		log.Errorf("[Role Service][EnableOrDisable] role is already %v", enable)
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	} else if existing.Enable == enable && !enable {
		log.Errorf("[Role Service][EnableOrDisable] role is already %v", enable)
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}

	if !enable {
		hasActive, err := j.approveIndexRepo.HasActiveActionRoles(ctx, existing.Code)
		if err != nil {
			log.Errorf("[Role Service][EnableOrDisable] failed to check active action roles: %v", err)
			return err
		}
		if hasActive {
			log.Errorf("[Role Service][EnableOrDisable] role %s has active jobs, cannot disable", existing.Code)
			return errors.New(localization.ErrorRoleHasActiveJobs.Code)
		}
	}

	updated := *existing
	updated.Enable = enable
	updated.UpdatedAt = time.Now()

	var requestType string
	if enable {
		requestType = string(constants.RequestEnableRole)
	} else {
		requestType = string(constants.RequestDisableRole)
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, requestType, constants.UPDATE)

	if err := j.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[Role Service][EnableOrDisable] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (j *RoleService) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, j.logger)

	makerUser := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[Role Service][Delete] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[Role Service][Delete] failed to find existing role: %v", err)
		return err
	}

	// Check if role has active jobs before allowing deletion
	hasActive, err := j.approveIndexRepo.HasActiveActionRoles(ctx, existing.Code)
	if err != nil {
		log.Errorf("[Role Service][Delete] failed to check active action roles: %v", err)
		return err
	}
	if hasActive {
		log.Errorf("[Role Service][Delete] role %s has active jobs, cannot delete", existing.Code)
		return errors.New(localization.ErrorRoleHasActiveJobs.Code)
	}

	hasActiveJobRole, err := j.jobRepo.HasActiveJobRole(ctx, existing.Code)
	if err != nil {
		log.Errorf("[Role Service][Delete] failed to check active job roles: %v", err)
		return err
	}
	if hasActiveJobRole {
		log.Errorf("[Role Service][Delete] role %s has active jobs, cannot delete", existing.Code)
		return errors.New(localization.ErrorRoleHasActiveJobs.Code)
	}

	updated := *existing
	updated.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, string(constants.RequestDeleteRole), constants.DELETE)

	if err := j.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[Role Service][Delete] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (j *RoleService) FindById(ctx context.Context, id string) (*imodel.Role, error) {
	log := local_util.LoggerFromCtx(ctx, j.logger)

	role, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var viewerActions, makerActions, checkerActions, auditorActions []string
	if strings.TrimSpace(role.Type) == "CPS" {
		viewerActions, makerActions, checkerActions, auditorActions, _, err = j.approveIndexRepo.PopulateUserApproverAllocations(ctx, role.Code)
		if err != nil {
			log.Errorf("[Role Service][FindById] failed to populate CPS approver allocations: %v", err)
			return role, nil
		}
	} else if strings.TrimSpace(role.Type) == "BPS" {
		viewerActions, makerActions, checkerActions, auditorActions, err = j.bpsApproverIndexRepo.PopulateUserApproverAllocations(ctx, role.Code)
		if err != nil {
			log.Errorf("[Role Service][FindById] failed to populate BPS approver allocations: %v", err)
			return role, nil
		}
	}

	role.ViewerActions = viewerActions
	role.MakerActions = makerActions
	role.CheckerActions = checkerActions
	role.AuditorActions = auditorActions

	log.Infof("[Role Service][FindById] role: %v", role)
	return role, nil
}

func (j *RoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Role], error) {
	return j.roleRepository.FindAllWithPagination(ctx, filterParam)
}

func (j *RoleService) Authorize(ctx context.Context, cpsAction *sharedmodel.CPSAction) (*sharedmodel.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, j.logger)

	var asAny any
	raw, err := json.Marshal(cpsAction.CurrentAction)
	log.Infof("[Role Service][Authorize] raw CPS action data: %s", string(raw))
	if err != nil {
		log.Errorf("[JobRole Service][Authorize] marshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if err := json.Unmarshal(raw, &asAny); err != nil {
		log.Errorf("[JobRole Service][Authorize] unmarshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[Role Service][Authorize] CPS action data as any: %v", asAny)

	raw2, _ := json.Marshal(asAny)
	var role imodel.Role
	if err := json.Unmarshal(raw2, &role); err != nil {
		log.Errorf("[JobRole Service][Authorize] map to Role failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[Role Service][Authorize] CPS action data as JobRole: %v", role)
	if cpsAction.UniqueId != "" {
		if oid, err := bson.ObjectIDFromHex(cpsAction.UniqueId); err == nil {
			role.ID = oid
		}
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateRole):
		role.CreatedAt = time.Now()
		if err := j.roleRepository.Create(ctx, &role); err != nil {
			return nil, err
		}
	case string(constants.RequestUpdateRole):
		role.UpdatedAt = time.Now()
		if err := j.roleRepository.Update(ctx, role.ID.Hex(), &role); err != nil {
			return nil, err
		}
	case string(constants.RequestEnableRole):
		if err := j.roleRepository.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			return nil, err
		}
	case string(constants.RequestDisableRole):
		if err := j.roleRepository.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			return nil, err
		}
	case string(constants.RequestDeleteRole):
		if err := j.roleRepository.SoftDelete(ctx, cpsAction.UniqueId); err != nil {
			return nil, err
		}
	default:
		log.Errorf("[JobRole Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = role
	return cpsAction, nil
}
