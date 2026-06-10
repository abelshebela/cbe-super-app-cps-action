package jobrole

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/job_role/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type jobRoleService struct {
	cpsService        service.CPSActionService
	JobRoleRepository storage.JobRoleRepository
	RoleRepository    storage.RoleRepository
	cpsUserRepo       storage.CpsUserRepository
	cfg               config.VaultConfig
	logger            utils.Logger
}

func NewJobRoleService(jobRole storage.JobRoleRepository, roleRepo storage.RoleRepository, cpsService service.CPSActionService, cpsUserRepo storage.CpsUserRepository, cfg config.VaultConfig, logger utils.Logger) service.JobRoleService {
	return &jobRoleService{
		cpsService:        cpsService,
		JobRoleRepository: jobRole,
		RoleRepository:    roleRepo,
		cpsUserRepo:       cpsUserRepo,
		cfg:               cfg,
		logger:            logger,
	}
}

func (j *jobRoleService) Create(ctx context.Context, role imodel.JobRole) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, j.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[JobRole Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if err := core.CheckRoleExistent(ctx, role.Role, j.RoleRepository); err != nil {
		log.Errorf("[JobRole Service] the give role not found %v", err)
		return err
	}

	if err := core.JobTitleExistentChecker(ctx, constants.CREATE, "", role.JobTitle, j.JobRoleRepository); err != nil {
		log.Errorf("[JobRole Service] the give job title already exists %v", err)
		if err.Error() != localization.ErrorResourceNotFound.Code {
			log.Errorf("[JobRole Service] the give job title not exists")
			return err
		}
	}

	if err := core.CodeExistentChecker(ctx, constants.CREATE, "", role.Code, j.JobRoleRepository); err != nil {
		log.Errorf("[JobRole Service] the given code already exists %v", err)
		return err
	}

	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, role, constants.RequestCreateJobRole, constants.CREATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

// FindAll returns enabled rows from the roles collection with job_roles fields (type, role_code, role_name) via RoleRepository aggregation. Used by GET /job_roles/all.
func (j *jobRoleService) FindAll(ctx context.Context) (*[]imodel.JobRole, error) {
	return j.JobRoleRepository.FindAll(ctx)
}

func (j *jobRoleService) Update(ctx context.Context, id string, update imodel.JobRole) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, j.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[JobRole Service][Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prev, err := j.JobRoleRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if update.JobTitle != "" {
		if err := core.CheckJobTitleExistent(ctx, *prev, update.JobTitle, j.JobRoleRepository); err != nil {
			log.Errorf("[JobRole Service] the give job title not found")
			return err
		}
	}

	newRole := *prev
	if update.Code != "" {
		if err := core.CodeExistentChecker(ctx, constants.UPDATE, id, update.Code, j.JobRoleRepository); err != nil {
			log.Errorf("[JobRole Service] the given code already exists %v", err)
			return err
		}
		newRole.Code = update.Code
	}
	if update.JobTitle != "" {
		newRole.JobTitle = update.JobTitle
	}
	if update.Role != "" {
		newRole.Role = update.Role
	}

	if update.Code != "" {
		newRole.Code = update.Code
	}

	newRole.UpdateAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, newRole, constants.RequestUpdateJobRole, constants.UPDATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *jobRoleService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	makerUser := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, j.logger)
	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[JobRole Service][EnableOrDisable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := j.JobRoleRepository.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[JobRole Service][EnableOrDisable] failed to find existing job role: %v", err)
		return err
	}

	if !enable && makerUser.UserRole == existing.Role {
		log.Errorf("[JobRole Service][EnableOrDisable] user cannot disable their own role")
		return errors.New(localization.ErrorCannotDisableOwnJobTitle.Code)
	}

	updated := *existing
	updated.UpdateAt = time.Now()

	var requestType string
	if enable {
		requestType = constants.RequestEnableJobRole
	} else {
		requestType = constants.RequestDisableJobRole
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, requestType, constants.UPDATE)

	if err := j.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[JobRole Service][EnableOrDisable] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (j *jobRoleService) Delete(ctx context.Context, id string) error {
	makerUser := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, j.logger)
	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[JobRole Service][Delete] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	hasActive, err := j.JobRoleRepository.CheckIfExists(ctx, id)
	if err != nil {
		log.Errorf("[JobRole Service][Delete] failed to find existing job role: %v", err)
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
	}

	// lookup cps user by jobtitle
	cpsUser, err := j.cpsUserRepo.GetUserByJobTitle(ctx, hasActive.JobTitle)
	if err != nil {
		log.Errorf("[JobRole Service][Delete] failed to find CPS user by job title: %v", err)
		return err
	}
	if cpsUser != nil {
		log.Errorf("[JobRole Service][Delete] cannot delete job role with active CPS user")
		return localization.ErrorJobTitleHasActiveUsers
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, hasActive, nil, constants.RequestDeleteJobRole, constants.DELETE)

	if err := j.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[JobRole Service][Delete] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

// FindById returns one roles document by id with job_roles join fields. Used by GET /job_roles/{id}.
func (j *jobRoleService) FindById(ctx context.Context, id string) (*imodel.JobRole, error) {
	return j.JobRoleRepository.FindByID(ctx, id)
}

// FindAllWithPagination lists roles collection rows with filters and job_roles enrichment. Used by GET /job_roles.
func (j *jobRoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.JobRole], error) {
	return j.JobRoleRepository.FindAllWithPagination(ctx, filterParam)
}

func (j *jobRoleService) Authorize(ctx context.Context, cpsAction *sharedmodel.CPSAction) (*sharedmodel.CPSAction, error) {
	// Turn CurrentAction into Role, attach ID from UniqueId (if present), apply action
	var asAny any
	log := local_util.LoggerFromCtx(ctx, j.logger)
	raw, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[JobRole Service][Authorize] marshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if err := json.Unmarshal(raw, &asAny); err != nil {
		log.Errorf("[JobRole Service][Authorize] unmarshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	raw2, _ := json.Marshal(asAny)
	var role imodel.JobRole
	if err := json.Unmarshal(raw2, &role); err != nil {
		log.Errorf("[JobRole Service][Authorize] map to Role failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if cpsAction.UniqueId != "" {
		if oid, err := bson.ObjectIDFromHex(cpsAction.UniqueId); err == nil {
			role.ID = oid
		}
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateJobRole):
		role.CreatedAt = time.Now()
		if err := j.JobRoleRepository.Create(ctx, &role); err != nil {
			log.Errorf("[JobRole Service][Authorize] failed to create job role: %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateJobRole):
		role.UpdateAt = time.Now()
		if err := j.JobRoleRepository.Update(ctx, role.ID.Hex(), &role); err != nil {
			log.Errorf("[JobRole Service][Authorize] failed to update job role: %v", err)
			return nil, err
		}
	case constants.RequestEnableJobRole:
		if err := j.JobRoleRepository.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			log.Errorf("[JobRole Service][Authorize] failed to enable job role: %v", err)
			return nil, err
		}
	case constants.RequestDisableJobRole:
		if err := j.JobRoleRepository.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			log.Errorf("[JobRole Service][Authorize] failed to disable job role: %v", err)
			return nil, err
		}
	case constants.RequestDeleteJobRole:
		if err := j.JobRoleRepository.SoftDelete(ctx, cpsAction.UniqueId); err != nil {
			log.Errorf("[JobRole Service][Authorize] failed to delete job role: %v", err)
			return nil, err
		}
	default:
		log.Errorf("[JobRole Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = role
	return cpsAction, nil
}
