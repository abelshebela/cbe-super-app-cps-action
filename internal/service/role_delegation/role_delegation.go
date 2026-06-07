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
	cpsService   service.CPSActionService
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
	case string(constants.RequestCreateRoleDelegation):
		roleDelegation.CreatedAt = time.Now()
		if err := r.repo.Create(ctx, &roleDelegation); err != nil {
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
		if err := r.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
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
func (r *roleDelegation) Create(ctx context.Context, roleDelegation imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegation/Create] Creating role delegation for user %s with job title %s", roleDelegation.UserID, roleDelegation.JobTitle)

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	// Validate user existence
	user, err := r.cpsUserRepo.FindByUserID(ctx, roleDelegation.UserID)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
		return err
	}
	if user == nil || !user.Enabled {
		r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.UserID)
		return errors.New(localization.ErrorUserNotFoundOrDisabled.Code)
	}

	jobRole, err := r.jobTitleRepo.FindByName(ctx, roleDelegation.JobTitle)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/Create] Failed to find job role: %v", err)
		return err
	}
	if jobRole == nil || !jobRole.Enabled {
		r.logger.Warnf("[RoleDelegation/Create] Job role not found: %s", roleDelegation.JobTitle)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}

	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, roleDelegation, constants.RequestCreateRoleDelegation, constants.CREATE)
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
	if update.JobTitle != "" {
		newRole.JobTitle = update.JobTitle
	}
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

func NewRoleDelegationService(repo storage.RoleDelegationRepository, jobTitleRepo storage.JobRoleRepository, cpsUserRepo storage.CpsUserRepository, cpsService service.CPSActionService, logger utils.Logger) service.RoleDelegationService {
	return &roleDelegation{
		cpsService:   cpsService,
		repo:         repo,
		jobTitleRepo: jobTitleRepo,
		cpsUserRepo:  cpsUserRepo,
		logger:       logger,
	}
}
