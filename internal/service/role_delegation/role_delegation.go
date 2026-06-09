package role_delegation_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type roleDelegation struct {
	repo         storage.RoleDelegationRepository
	jobTitleRepo storage.JobRoleRepository
	cpsUserRepo  storage.CpsUserRepository
	bpsUserRepo  storage.BPSUserRepository
	department   storage.DepartmentRepository
	cpsService   service.CPSActionService
	branchRepo   storage.AccountBlockRepository
	logger       utils.Logger
}

// Authorize implements [service.RoleDelegationService].
func (r *roleDelegation) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	var asAny any
	raw, err := json.Marshal(cpsAction.CurrentAction)
	log.Infof("[RoleDelegation Service][Authorize] raw CPS action data: %s", string(raw))
	if err != nil {
		log.Errorf("[RoleDelegation Service][Authorize] marshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if err := json.Unmarshal(raw, &asAny); err != nil {
		log.Errorf("[RoleDelegation Service][Authorize] unmarshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[RoleDelegation Service][Authorize] CPS action data as any: %v", asAny)

	raw2, _ := json.Marshal(asAny)
	var roleDelegation imodel.RoleDelegation
	if err := json.Unmarshal(raw2, &roleDelegation); err != nil {
		log.Errorf("[RoleDelegation Service][Authorize] map to RoleDelegation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[RoleDelegation Service][Authorize] CPS action data as RoleDelegation: %v", roleDelegation)
	if cpsAction.UniqueId != "" {
		if oid, err := bson.ObjectIDFromHex(cpsAction.UniqueId); err == nil {
			roleDelegation.ID = oid
		}
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateRoleDelegationForExistingUser):
		roleDelegation.CreatedAt = time.Now()
		if err := r.repo.CreateWithExistingUser(ctx, &roleDelegation); err != nil {
			return nil, err
		}
	case string(constants.RequestCreateRoleDelegationForNewUser):
		roleDelegation.CreatedAt = time.Now()
		if err := r.repo.CreateWithNewUser(ctx, &roleDelegation); err != nil {
			return nil, err
		}
	case string(constants.RequestUpdateRoleDelegation):
		roleDelegation.UpdatedAt = time.Now()
		if err := r.repo.Update(ctx, roleDelegation.ID.Hex(), &roleDelegation); err != nil {
			return nil, err
		}
	case string(constants.RequestEnableRoleDelegation):
		if err := r.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			return nil, err
		}
	case string(constants.RequestDisableRoleDelegation):
		if err := r.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			return nil, err
		}
	case string(constants.RequestDeleteRoleDelegation):
		if err := r.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			return nil, err
		}
	default:
		log.Errorf("[RoleDelegation Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = roleDelegation
	return cpsAction, nil
}

// Create implements [service.RoleDelegationService].
func (r *roleDelegation) CreateWithExistingUser(ctx context.Context, roleDelegation imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegation/Create] Creating role delegation for user %s with job title %s", roleDelegation.DelegatedUserID, roleDelegation.DelegatedUserJobTitle)

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	if roleDelegation.DelegatedUserUserType == "BPS" {

		user, err := r.bpsUserRepo.GetByUserCode(ctx, roleDelegation.DelegatedUserID)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find BPS user: %v", err)
			return err
		}
		if user == nil || !user.Enabled {
			r.logger.Warnf("[RoleDelegation/Create] BPS User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorUserNotFoundOrDisabled.Code)
		}
		roleDelegation.DelegatedUserFullName = user.FullName
		roleDelegation.DelegatedUserDepartmentOrBranch = user.BranchName
		roleDelegation.DelegatedUserJobTitle = user.JobTitle
		roleDelegation.DelegatedUserExistingRole = user.Role
		roleDelegation.DelegatedUserPhoneNumber = user.PhoneNumber
		roleDelegation.DelegatedUserEmail = user.Email
	} else {
		// Validate user existence
		user, err := r.cpsUserRepo.FindByUsername(ctx, roleDelegation.DelegatedUserID)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return err
		}

		if user == nil || !user.Enabled {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorUserNotFoundOrDisabled.Code)
		}

		roleDelegation.DelegatedUserFullName = user.FullName
		roleDelegation.DelegatedUserDepartmentOrBranch = user.DelegationID.Hex()
		roleDelegation.DelegatedUserJobTitle = user.JobTitle
		roleDelegation.DelegatedUserExistingRole = user.Role
		roleDelegation.DelegatedUserPhoneNumber = user.PhoneNumber
		roleDelegation.DelegatedUserEmail = user.Email
	}

	jobRole, err := r.jobTitleRepo.FindByRole(ctx, roleDelegation.NewRoleID)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/Create] Failed to find job role: %v", err)
		return err
	}
	if jobRole == nil || !jobRole.Enabled {
		r.logger.Errorf("[RoleDelegation/Create] Job role not found: %s", roleDelegation.NewRoleID)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}

	roleDelegation.DelegatedUserExistingRole = jobRole.Role
	if roleDelegation.DelegatorUserUserType == "BPS" {
		if branches, err := r.branchRepo.GetBranchesByIds(ctx, []string{roleDelegation.NewDepartmentOrBranch}); err != nil || branches == nil || len(branches) == 0 {
			r.logger.Errorf("[RoleDelegation/Create] Branch not found: %s", roleDelegation.NewDepartmentOrBranch)
			return localization.ErrorInvalidDelegationDepartment
		}
	} else {
		department, err := r.department.FindByID(ctx, roleDelegation.NewDepartmentOrBranch)
		if department == nil || !department.Enabled {
			r.logger.Errorf("[RoleDelegation/Create] Department not found: %s err: %v", roleDelegation.DelegatedUserDepartmentOrBranch, err.Error())
			return localization.ErrorInvalidDelegationDepartment
		}
	}
	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, roleDelegation, constants.RequestCreateRoleDelegationForExistingUser, constants.CREATE)
	return r.cpsService.CreateCPSAction(ctx, &cpsModel)

}

// Create implements [service.RoleDelegationService].
func (r *roleDelegation) CreateWithNewUser(ctx context.Context, roleDelegation imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegation/Create] Creating role delegation for user %s with job title %s", roleDelegation.DelegatedUserID, roleDelegation.DelegatedUserJobTitle)

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	if roleDelegation.DelegatedUserUserType == "BPS" {
		user, err := r.bpsUserRepo.FindByOr(ctx, roleDelegation.DelegatedUserPhoneNumber, roleDelegation.DelegatedUserEmail, roleDelegation.DelegatedUserID)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find BPS user: %v", err)
			return localization.ErrorUnexpectedError
		}
		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] BPS User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorBpsUserAlreadyExists.Code)
		}
	} else {
		// Validate user existence
		user, err := r.cpsUserRepo.FindByUsername(ctx, roleDelegation.DelegatedUserID)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return localization.ErrorUnexpectedError
		}

		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorCpsUserAlreadyExists.Code)
		}
		// Validate user existence
		user, err = r.cpsUserRepo.FindByPhoneNumber(ctx, roleDelegation.DelegatedUserPhoneNumber)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return localization.ErrorUnexpectedError
		}

		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorCpsUserAlreadyExists.Code)
		}
		// Validate user existence
		user, err = r.cpsUserRepo.FindByEmail(ctx, roleDelegation.DelegatedUserEmail)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return localization.ErrorUnexpectedError
		}

		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorCpsUserAlreadyExists.Code)
		}
	}

	jobRole, err := r.jobTitleRepo.FindByRole(ctx, roleDelegation.NewRoleID)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/Create] Failed to find job role: %v", err)
		return err
	}

	if jobRole == nil || !jobRole.Enabled {
		r.logger.Errorf("[RoleDelegation/Create] Job role not found: %s", roleDelegation.NewRoleID)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	if roleDelegation.DelegatorUserUserType == "BPS" {
		if branches, err := r.branchRepo.GetBranchesByIds(ctx, []string{roleDelegation.NewDepartmentOrBranch}); err != nil || branches == nil || len(branches) == 0 {
			r.logger.Errorf("[RoleDelegation/Create] Branch not found: %s", roleDelegation.NewDepartmentOrBranch)
			return localization.ErrorInvalidDelegationDepartment
		}
	} else {
		department, err := r.department.FindByID(ctx, roleDelegation.NewDepartmentOrBranch)
		if department == nil || !department.Enabled {
			r.logger.Errorf("[RoleDelegation/Create] Department not found: %s err: %v", roleDelegation.DelegatedUserDepartmentOrBranch, err.Error())
			return localization.ErrorInvalidDelegationDepartment
		}
		roleDelegation.DelegatedUserDepartmentOrBranch = department.ID.Hex()
	}

	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, roleDelegation, constants.RequestCreateRoleDelegationForNewUser, constants.CREATE)
	return r.cpsService.CreateCPSAction(ctx, &cpsModel)

}

// Delete implements [service.RoleDelegationService].
func (r *roleDelegation) Delete(ctx context.Context, id string) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[RoleDelegation/Delete] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := r.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[RoleDelegation/Delete] failed to find existing role delegation: %v", err)
		return err
	}

	cpsActionData := lib.CpsModelBuilder(id, maker, existing, nil, constants.RequestDeleteRoleDelegation, constants.DELETE)
	return r.cpsService.CreateCPSAction(ctx, &cpsActionData)
}

// EnableOrDisable implements [service.RoleDelegationService].
func (r *roleDelegation) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[RoleDelegation/EnableOrDisable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := r.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[RoleDelegation/EnableOrDisable] failed to find existing role delegation: %v", err)
		return err
	}

	if enable && existing.Enable {
		log.Warnf("[RoleDelegation/EnableOrDisable] role delegation %s is already enabled", id)
		return errors.New(localization.ErrorRoleDelegationAlreadyEnabled.Code)
	}

	if !enable && !existing.Enable {
		log.Warnf("[RoleDelegation/EnableOrDisable] role delegation %s is already disabled", id)
		return errors.New(localization.ErrorRoleDelegationAlreadyDisabled.Code)
	}

	updated := *existing
	updated.UpdatedAt = time.Now()

	var requestType string
	if enable {
		requestType = constants.RequestEnableRoleDelegation
	} else {
		requestType = constants.RequestDisableRoleDelegation
	}

	cpsActionData := lib.CpsModelBuilder(id, maker, existing, updated, requestType, constants.UPDATE)
	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[RoleDelegation/EnableOrDisable] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

// FindAll implements [service.RoleDelegationService].
func (r *roleDelegation) FindAll(ctx context.Context) (*[]imodel.RoleDelegation, error) {
	return r.repo.FindAll(ctx)
}

// FindAllWithPagination implements [service.RoleDelegationService].
func (r *roleDelegation) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	return r.repo.FindAllWithPagination(ctx, filterParam)
}

// FindById implements [service.RoleDelegationService].
func (r *roleDelegation) FindById(ctx context.Context, id string) (*imodel.RoleDelegation, error) {
	return r.repo.FindByID(ctx, id)
}

// FindByUsername implements [service.RoleDelegationService].
func (r *roleDelegation) FindByUsername(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	return r.repo.FindByUsername(ctx, id, filterParam)
}

// Update implements [service.RoleDelegationService].
func (r *roleDelegation) Update(ctx context.Context, id string, update imodel.RoleDelegation) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[RoleDelegation/Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prev, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	newRole := *prev

	if !update.StartAt.IsZero() {
		newRole.StartAt = update.StartAt
	}
	if !update.EndAt.IsZero() {
		newRole.EndAt = update.EndAt
	}

	newRole.UpdatedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, newRole, constants.RequestUpdateRoleDelegation, constants.UPDATE)
	return r.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func NewRoleDelegationService(repo storage.RoleDelegationRepository, jobTitleRepo storage.JobRoleRepository, cpsUserRepo storage.CpsUserRepository, bpsUserRepo storage.BPSUserRepository, department storage.DepartmentRepository, branchRepo storage.AccountBlockRepository, cpsService service.CPSActionService, logger utils.Logger) service.RoleDelegationService {
	return &roleDelegation{
		cpsService:   cpsService,
		repo:         repo,
		jobTitleRepo: jobTitleRepo,
		cpsUserRepo:  cpsUserRepo,
		bpsUserRepo:  bpsUserRepo,
		branchRepo:   branchRepo,
		department:   department,
		logger:       logger,
	}
}
