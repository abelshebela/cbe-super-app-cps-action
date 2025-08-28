package permission

import (
	"context"
	"errors"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/permission"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/permission/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type permissionService struct {
	repo       storage.PermissionRepository
	logger     shared_utils.Logger
	cpsService service.CPSActionService
}

func InitPermissionService(repo storage.PermissionRepository, cpsService service.CPSActionService, logger shared_utils.Logger) service.PermissionService {
	return &permissionService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *permissionService) CreatePermissionGroup(ctx context.Context, req permission.CreatePermissionGroupRequest) error {
	// Validate group name
	if req.GroupName == "" {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}

	if s.repo.CheckPermissionGroupExists(req.GroupName) {
		return errors.New("PERMISSION_GROUP_ALREADY_EXISTS")
	}

	if len(req.PermissionCategoryLists) > 0 {
		validCategories, err := s.repo.ValidatePermissionCategories(ctx, req.PermissionCategoryLists)
		if err != nil {
			return err
		}
		if len(validCategories) != len(req.PermissionCategoryLists) {
			return errors.New("SOME_PERMISSION_CATEGORIES_NOT_FOUND")
		}
	}

	permissionGroup := core.PermissionGroupModel(req)

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(
		req.GroupName,
		maker,
		nil,
		permissionGroup,
		string(constants.CREATE),
		constants.CREATE,
	)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *permissionService) UpdatePermissionGroup(ctx context.Context, req permission.UpdatePermissionGroupRequest) error {

	if req.OldGroupName == "" {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}

	existingGroup, err := s.repo.GetPermissionGroup(req.OldGroupName)
	if err != nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	if req.NewGroupName != "" && req.NewGroupName != req.OldGroupName {
		if s.repo.CheckPermissionGroupExists(req.NewGroupName) {
			return errors.New("PERMISSION_GROUP_ALREADY_EXISTS")
		}
	}

	if len(req.PermissionCategoryLists) > 0 {
		validCategories, err := s.repo.ValidatePermissionCategories(ctx, req.PermissionCategoryLists)
		if err != nil {
			return err
		}
		if len(validCategories) != len(req.PermissionCategoryLists) {
			return errors.New("SOME_PERMISSION_CATEGORIES_NOT_FOUND")
		}
	}

	updatedGroup := core.PermissionGroupUpdateModel(req)

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(
		req.OldGroupName,
		maker,
		existingGroup,
		updatedGroup,
		string(constants.UPDATE),
		constants.UPDATE,
	)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *permissionService) GetPermissionGroup(groupName string) (*model.PermissionGroup, error) {
	if groupName == "" {
		return nil, errors.New(localization.ErrorPermissionGroupRequired.Code)
	}
	groupName = strings.ToUpper(groupName)

	return s.repo.GetPermissionGroup(groupName)
}

func (s *permissionService) GetPermissionGroups(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	if filterParams == nil {
		filter := types.Filter{}
		filterParams = &filter
	}

	return s.repo.FindAllWithPagination(ctx, *filterParams)
}

func (s *permissionService) GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error) {
	return s.repo.GetAllPermissionCategoriesWithPermissions(ctx)
}

func (s *permissionService) ValidatePermissionCategories(ctx context.Context, categoryIDs []string) (bool, error) {
	validCategories, err := s.repo.ValidatePermissionCategories(ctx, categoryIDs)
	if err != nil {
		s.logger.Errorf("Failed to validate permission categories: %v", err)
		return false, err
	}

	// Check if all requested categories were validated
	if len(validCategories) != len(categoryIDs) {
		s.logger.Warnf("Some permission categories were not found. Requested: %d, Valid: %d", len(categoryIDs), len(validCategories))
		return false, errors.New(localization.ErrorPermissionCategoryNotFound.Code)
	}

	return true, nil
}

func (s *permissionService) ValidatePermissionGroups(ctx context.Context, groupIDs []string) (bool, error) {
	validGroups, err := s.repo.ValidatePermissionGroups(ctx, groupIDs)
	if err != nil {
		s.logger.Errorf("Failed to validate permission groups: %v", err)
		return false, errors.New(localization.ErrorPermissionGroupValidationFailed.Code)
	}

	// Check if all requested groups were validated
	if len(validGroups) != len(groupIDs) {
		s.logger.Warnf("Some permission groups were not found. Requested: %d, Valid: %d", len(groupIDs), len(validGroups))
		// return false, errors.New("SOME_PERMISSION_GROUPS_NOT_FOUND")
		return false, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	return true, nil
}

func (s *permissionService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	switch action.ActionType {
	case string(constants.CREATE):
		cur, ok := action.CurrentAction.(model.PermissionGroup)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Create(ctx, &cur); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.UPDATE):
		upd, ok := action.CurrentAction.(model.PermissionGroup)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Update(ctx, upd.GroupName, &upd); err != nil {
			return nil, err
		}
		return action, nil

	default:
		return nil, errors.New("UNHANDLED_ACTION_TYPE")
	}
}
