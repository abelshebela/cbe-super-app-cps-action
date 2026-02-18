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
	jobRoleRepository storage.JobRoleRepository
	roleRepository    storage.RoleRepository
	cfg               config.VaultConfig
	logger            utils.Logger
}

func NewJobRoleService(jobRole storage.JobRoleRepository, roleRepo storage.RoleRepository, cpsService service.CPSActionService, cfg config.VaultConfig, logger utils.Logger) service.JobRoleService {
	return &jobRoleService{
		cpsService:        cpsService,
		jobRoleRepository: jobRole,
		roleRepository:    roleRepo,
		cfg:               cfg,
		logger:            logger,
	}
}

func (j *jobRoleService) Create(ctx context.Context, role imodel.Role) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		j.logger.Errorf("[JobRole Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if err := core.CheckRoleExistent(ctx, role.Role, j.jobRoleRepository); err != nil {
		j.logger.Errorf("[JobRole Service] the give role not found %v", err)
		return err
	}

	if err := core.JobTitleExistentChecker(ctx, constants.CREATE, "", role.JobTitle, j.roleRepository); err != nil {
		j.logger.Errorf("[JobRole Service] the give job title already exists %v", err)
		if err.Error() != localization.ErrorResourceNotFound.Code {
			j.logger.Errorf("[JobRole Service] the give job title not exists")
			return err
		}
	}
	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, role, constants.RequestCreateJobRole, constants.CREATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *jobRoleService) FindAll(ctx context.Context) (*[]imodel.Role, error) {
	return j.roleRepository.FindAll(ctx)
}

func (j *jobRoleService) Update(ctx context.Context, id string, update imodel.Role) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		j.logger.Errorf("[JobRole Service][Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prev, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if update.JobTitle != "" {
		if err := core.CheckJobTitleExistent(ctx, *prev, update.JobTitle, j.roleRepository); err != nil {
			j.logger.Errorf("[JobRole Service] the give job title not found")
			return err
		}
	}

	newRole := *prev
	if update.JobTitle != "" {
		newRole.JobTitle = update.JobTitle
	}
	if update.Role != "" {
		newRole.Role = update.Role
	}

	newRole.UpdateAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, newRole, constants.RequestUpdateJobRole, constants.UPDATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *jobRoleService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	makerUser := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerUser) {
		j.logger.Errorf("[JobRole Service][EnableOrDisable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		j.logger.Errorf("[JobRole Service][EnableOrDisable] failed to find existing job role: %v", err)
		return err
	}

	if !enable && makerUser.UserRole == existing.Role {
		j.logger.Errorf("[JobRole Service][EnableOrDisable] user cannot disable their own role")
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
		j.logger.Errorf("[JobRole Service][EnableOrDisable] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (j *jobRoleService) Delete(ctx context.Context, id string) error {
	makerUser := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerUser) {
		j.logger.Errorf("[JobRole Service][Delete] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		j.logger.Errorf("[JobRole Service][Delete] failed to find existing job role: %v", err)
		return err
	}

	updated := *existing
	updated.UpdateAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existing, updated, constants.RequestDeleteJobRole, constants.DELETE)

	if err := j.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		j.logger.Errorf("[JobRole Service][Delete] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (j *jobRoleService) FindById(ctx context.Context, id string) (*imodel.Role, error) {
	return j.roleRepository.FindByID(ctx, id)
}

func (j *jobRoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Role], error) {
	return j.roleRepository.FindAllWithPagination(ctx, filterParam)
}

func (j *jobRoleService) Authorize(ctx context.Context, cpsAction *sharedmodel.CPSAction) (*sharedmodel.CPSAction, error) {
	// Turn CurrentAction into Role, attach ID from UniqueId (if present), apply action
	var asAny any
	raw, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		j.logger.Errorf("[JobRole Service][Authorize] marshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if err := json.Unmarshal(raw, &asAny); err != nil {
		j.logger.Errorf("[JobRole Service][Authorize] unmarshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	raw2, _ := json.Marshal(asAny)
	var role imodel.Role
	if err := json.Unmarshal(raw2, &role); err != nil {
		j.logger.Errorf("[JobRole Service][Authorize] map to Role failed: %v", err)
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
		if err := j.roleRepository.Create(ctx, &role); err != nil {
			return nil, err
		}
	case string(constants.RequestUpdateJobRole):
		role.UpdateAt = time.Now()
		if err := j.roleRepository.Update(ctx, role.ID.Hex(), &role); err != nil {
			return nil, err
		}
	case constants.RequestEnableJobRole:
		if err := j.roleRepository.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			return nil, err
		}
	case constants.RequestDisableJobRole:
		if err := j.roleRepository.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			return nil, err
		}
	case constants.RequestDeleteJobRole:
		if err := j.roleRepository.SoftDelete(ctx, cpsAction.UniqueId); err != nil {
			return nil, err
		}
	default:
		j.logger.Errorf("[JobRole Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = role
	return cpsAction, nil
}
