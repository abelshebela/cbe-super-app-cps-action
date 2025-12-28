package jobrole

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
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

func (j *jobRoleService) Create(ctx context.Context, role sharedmodel.Role) error {
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

func (j *jobRoleService) Update(ctx context.Context, id string, update sharedmodel.Role) error {
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
	if update.BranchGrade != "" {
		newRole.BranchGrade = update.BranchGrade
	}
	if update.Department != "" {
		newRole.Department = update.Department
	}
	if update.Position != "" {
		newRole.Position = update.Position
	}
	newRole.UpdatedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, newRole, constants.RequestUpdateJobRole, constants.UPDATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *jobRoleService) FindById(ctx context.Context, id string) (*sharedmodel.Role, error) {
	return j.roleRepository.FindByID(ctx, id)
}

func (j *jobRoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*sharedmodel.Role], error) {
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
	var role sharedmodel.Role
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
		role.UpdatedAt = time.Now()
		if err := j.roleRepository.Update(ctx, role.ID.Hex(), &role); err != nil {
			return nil, err
		}
	default:
		j.logger.Errorf("[JobRole Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = role
	return cpsAction, nil
}
