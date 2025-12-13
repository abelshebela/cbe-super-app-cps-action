package permission

import (
	"cbe-super-app-cps-action/internal/constants"
	cps_user_dto "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/dto/permission"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/permission/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type permissionService struct {
	repo       storage.PermissionRepository
	department storage.DepartmentRepository
	logger     shared_utils.Logger
	cpsService service.CPSActionService
}

func InitPermissionService(repo storage.PermissionRepository, dept storage.DepartmentRepository, cpsService service.CPSActionService, logger shared_utils.Logger) service.PermissionService {
	return &permissionService{
		repo:       repo,
		department: dept,
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
		return errors.New(localization.ErrorPermissionGroupAlreadyExists.Code)
	}

	if len(req.PermissionCategoryLists) > 0 {
		validCategories, err := s.repo.ValidatePermissionCategories(ctx, req.PermissionCategoryLists)
		if err != nil {
			return err
		}
		if len(validCategories) != len(req.PermissionCategoryLists) {

			return errors.New(localization.ErrorPermissionCatagoryNotFound.Code)
		}
	}

	dept, err := s.department.FindByID(ctx, req.DepartmentID)
	if err != nil || dept == nil {
		s.logger.Errorf("Department not found with ID: %s", req.DepartmentID)
		return errors.New(localization.ErrorDepartmentNotFound.Code)
	}

	permissionGroup := core.PermissionGroupModel(req)

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(
		req.GroupName,
		maker,
		nil,
		permissionGroup,
		string(constants.RequestCreatePermissionGroup),
		constants.CREATE,
	)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *permissionService) UpdatePermissionGroup(ctx context.Context, req permission.UpdatePermissionGroupRequest) error {
	s.logger.Infof("[UpdatePermissionGroup] updating permission group for id: %s", req.Id)
	existingGroup, err := s.repo.GetPermissionGroupById(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("[UpdatePermissionGroup] failed to find permission group: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	if req.NewGroupName != "" && req.NewGroupName != existingGroup.GroupName {
		if s.repo.CheckPermissionGroupExists(req.NewGroupName) {
			s.logger.Errorf("[UpdatePermissionGroup] permission group with new name already exists")
			return errors.New(localization.ErrorPermissionGroupAlreadyExists.Code)
		}
	}

	if len(req.PermissionCategoryLists) > 0 {
		validCategories, err := s.repo.ValidatePermissionCategories(ctx, req.PermissionCategoryLists)
		if err != nil {
			s.logger.Errorf("[UpdatePermissionGroup] failed to validate permission categories: %v", err)
			return err
		}

		if len(validCategories) != len(req.PermissionCategoryLists) {
			s.logger.Errorf("[UpdatePermissionGroup] invalid permission categories provided")
			return errors.New(localization.ErrorPermissionCatagoryNotFound.Code)
		}
	}

	updatedGroup := core.PermissionGroupUpdateModel(req)

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(
		req.Id,
		maker,
		existingGroup,
		updatedGroup,
		string(constants.RequestUpdatePermissionGroup),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[UpdatePermissionGroup] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[UpdatePermissionGroup] permission group update request created successfully for id: %s", req.Id)
	return nil
}

func (s *permissionService) GetPermissionGroup(groupName string) (*model.PermissionGroup, error) {
	if groupName == "" {
		s.logger.Errorf("[GetPermissionGroup] group name is empty")
		return nil, errors.New(localization.ErrorPermissionGroupRequired.Code)
	}
	groupName = strings.ToUpper(groupName)

	permissionGroup, err := s.repo.GetPermissionGroup(groupName)
	if err != nil {
		if err == mongo.ErrNoDocuments || err.Error() == "mongo: no documents in result" {
			s.logger.Errorf("[GetPermissionGroup] permission group not found: %s", groupName)
			return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
		}
		s.logger.Errorf("[GetPermissionGroup] failed to fetch permission group: %v", err)
		return nil, err
	}
	s.logger.Infof("[GetPermissionGroup] permission group retrieved successfully for name: %s", groupName)
	return permissionGroup, nil
}

func (s *permissionService) GetPermissionGroupById(ctx context.Context, id string) (*model.PermissionGroup, error) {
	if id == "" {
		s.logger.Errorf("[GetPermissionGroupById] id is empty")
		return nil, errors.New(localization.ErrorPermissionGroupRequired.Code)
	}

	group, err := s.repo.GetPermissionGroupById(ctx, id)
	if err != nil {
		s.logger.Errorf("[GetPermissionGroupById] failed to fetch permission group: %v", err)
		return nil, err
	}
	s.logger.Infof("[GetPermissionGroupById] permission group retrieved successfully for id: %s", id)
	return group, nil
}

func (s *permissionService) GetPermissionGroups(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	if filterParams == nil {
		filter := types.Filter{}
		filterParams = &filter
	}

	result, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("[GetPermissionGroups] failed to fetch permission groups: %v", err)
		return nil, err
	}
	s.logger.Infof("[GetPermissionGroups] retrieved %d permission groups", len(result.Data))
	return result, nil
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

	if len(validGroups) != len(groupIDs) {
		s.logger.Warnf("Some permission groups were not found. Requested: %d, Valid: %d", len(groupIDs), len(validGroups))
		return false, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	return true, nil
}

func (s *permissionService) GetPermissionCategoriesByDepartment(ctx context.Context, departmentId string) (map[string][]*model.PermissionCategory, error) {
	department, err := s.department.FindByID(ctx, departmentId)
	if err != nil || department == nil {
		s.logger.Errorf("Department not found with ID: %s", departmentId)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	portal_cards := department.PortalCards
	cardsWithPermission := make(map[string][]*model.PermissionCategory)

	if len(portal_cards) > 0 {
		for _, card := range portal_cards {
			categories, err := s.repo.GetAllPermissionCategories(ctx, card)
			if err != nil {
				s.logger.Errorf("Permission category can't be found with department ID")
				return nil, errors.New(localization.ErrorResourceNotFound.Code)
			}

			if len(categories) != 0 {
				cardsWithPermission[strings.ToLower(card)] = categories
			}
		}
	}

	return cardsWithPermission, nil
}

func (s *permissionService) GetPermissionGroupsByDepartment(ctx context.Context, departmentId string, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	department, err := s.repo.FindAllGroupsWithPagination(ctx, departmentId, filterParam)
	if err != nil || department == nil {
		s.logger.Errorf("Department not found with ID: %s", departmentId)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	return department, nil
}

func (s *permissionService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	switch action.ActionType {
	case string(constants.CREATE):
		cur, err := core.BindPermissionGroupFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Create(ctx, &cur); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.UPDATE):
		upd, err := core.BindPermissionGroupFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		existingGroup, err := s.repo.GetPermissionGroupById(ctx, action.UniqueId)
		if err != nil {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}

		// Update using the existing group's ObjectID
		err = s.repo.Update(ctx, existingGroup.ID.Hex(), &upd)
		if err != nil {
			return nil, err
		}
		return action, nil

	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *permissionService) GetPopulatedPermissionCategories(ctx context.Context, categoryIDsObject []bson.ObjectID) ([]cps_user_dto.PermissionCategoryResponse, error) {
	if len(categoryIDsObject) == 0 {
		return []cps_user_dto.PermissionCategoryResponse{}, nil
	}
	var categoryIDs []string
	for _, id := range categoryIDsObject {
		categoryIDs = append(categoryIDs, id.Hex())
	}

	validCategories, err := s.repo.ValidatePermissionCategories(ctx, categoryIDs)
	if err != nil {
		s.logger.Errorf("Failed to validate permission categories: %v", err)
		return nil, errors.New(localization.ErrorPermissionCategoryNotFound.Code)
	}

	if len(validCategories) != len(categoryIDs) {
		s.logger.Warnf("Some permission categories were not found. Requested: %d, Valid: %d", len(categoryIDs), len(validCategories))
		return nil, errors.New(localization.ErrorPermissionCategoryNotFound.Code)
	}

	return s.repo.GetPopulatedPermissionCategories(ctx, categoryIDs)
}

func (s *permissionService) GetPopulatedPermissionGroups(ctx context.Context, groupIDsObject []bson.ObjectID) ([]cps_user_dto.PermissionGroupResponse, error) {
	if len(groupIDsObject) == 0 {
		return []cps_user_dto.PermissionGroupResponse{}, nil
	}
	var groupIDs []string
	for _, id := range groupIDsObject {
		groupIDs = append(groupIDs, id.Hex())
	}

	validGroups, err := s.repo.ValidatePermissionGroups(ctx, groupIDs)
	if err != nil {
		s.logger.Errorf("Failed to validate permission groups: %v", err)
		return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	if len(validGroups) != len(groupIDs) {
		s.logger.Warnf("Some permission groups were not found. Requested: %d, Valid: %d", len(groupIDs), len(validGroups))
		return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	return s.repo.GetPopulatedPermissionGroups(ctx, groupIDs)
}
