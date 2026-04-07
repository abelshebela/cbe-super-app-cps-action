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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "CreatePermissionGroup", "Permission", "CreatePermissionGroup")
	defer span.End()

	// Validate group name
	if req.GroupName == "" {
		span.AddEvent("Invalid request", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidRequest.Code),
		))
		return errors.New(localization.ErrorInvalidRequest.Code)
	}

	if s.repo.CheckPermissionGroupExists(req.GroupName) {
		span.AddEvent("Permission group already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionGroupAlreadyExists.Code),
			attribute.String("group_name", req.GroupName),
		))
		return errors.New(localization.ErrorPermissionGroupAlreadyExists.Code)
	}

	if len(req.PermissionCategoryLists) > 0 {
		validCategories, err := s.repo.ValidatePermissionCategories(ctx, req.PermissionCategoryLists)
		if err != nil {
			span.AddEvent("Failed to validate permission categories", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}
		if len(validCategories) != len(req.PermissionCategoryLists) {
			span.AddEvent("Permission category not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorPermissionCatagoryNotFound.Code),
			))
			return errors.New(localization.ErrorPermissionCatagoryNotFound.Code)
		}
	}

	dept, err := s.department.FindByID(ctx, req.DepartmentID)
	if err != nil || dept == nil {
		s.logger.Errorf("[PermSvc][CreateGroup] dept not found: %s", req.DepartmentID)
		span.AddEvent("Department not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorDepartmentNotFound.Code),
			attribute.String("department_id", req.DepartmentID),
		))
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

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	return nil
}

func (s *permissionService) UpdatePermissionGroup(ctx context.Context, req permission.UpdatePermissionGroupRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdatePermissionGroup", "Permission", "UpdatePermissionGroup")
	defer span.End()

	s.logger.Infof("[PermSvc][UpdateGroup] id: %s", req.Id)
	existingGroup, err := s.repo.GetPermissionGroupById(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("[PermSvc][UpdateGroup] find err: %v", err)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("id", req.Id),
		))
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	if req.NewGroupName != "" && req.NewGroupName != existingGroup.GroupName {
		if s.repo.CheckPermissionGroupExists(req.NewGroupName) {
			s.logger.Errorf("[PermSvc][UpdateGroup] name already exists")
			span.AddEvent("Permission group already exists", trace.WithAttributes(
				attribute.String("error", localization.ErrorPermissionGroupAlreadyExists.Code),
				attribute.String("id", req.Id),
			))
			return errors.New(localization.ErrorPermissionGroupAlreadyExists.Code)
		}
	}

	if len(req.PermissionCategoryLists) > 0 {
		validCategories, err := s.repo.ValidatePermissionCategories(ctx, req.PermissionCategoryLists)
		if err != nil {
			s.logger.Errorf("[PermSvc][UpdateGroup] validate categories err: %v", err)
			span.AddEvent("Failed to validate permission categories", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", req.Id),
			))
			return err
		}

		if len(validCategories) != len(req.PermissionCategoryLists) {
			s.logger.Errorf("[PermSvc][UpdateGroup] invalid categories")
			span.AddEvent("Permission category not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorPermissionCatagoryNotFound.Code),
				attribute.String("id", req.Id),
			))
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
		s.logger.Errorf("[PermSvc][UpdateGroup] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", req.Id),
		))
		return err
	}
	s.logger.Infof("[PermSvc][UpdateGroup] request created id: %s", req.Id)
	return nil
}

func (s *permissionService) GetPermissionGroup(groupName string) (*model.PermissionGroup, error) {
	ctx := context.Background()
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPermissionGroup", "Permission", "GetPermissionGroup")
	defer span.End()

	if groupName == "" {
		s.logger.Errorf("[PermSvc][GetGroup] name empty")
		span.AddEvent("Group name is empty", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionGroupRequired.Code),
		))
		return nil, errors.New(localization.ErrorPermissionGroupRequired.Code)
	}
	groupName = strings.ToUpper(groupName)

	permissionGroup, err := s.repo.GetPermissionGroup(groupName)
	if err != nil {
		if err == mongo.ErrNoDocuments || err.Error() == "mongo: no documents in result" {
			s.logger.Errorf("[PermSvc][GetGroup] not found: %s", groupName)
			span.AddEvent("Permission group not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorPermissionGroupNotFound.Code),
				attribute.String("group_name", groupName),
			))
			return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
		}
		s.logger.Errorf("[PermSvc][GetGroup] fetch err: %v", err)
		span.AddEvent("Failed to fetch permission group", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("group_name", groupName),
		))
		return nil, err
	}
	s.logger.Infof("[PermSvc][GetGroup] retrieved: %s", groupName)
	return permissionGroup, nil
}

func (s *permissionService) GetPermissionGroupById(ctx context.Context, id string) (*model.PermissionGroup, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPermissionGroupById", "Permission", "GetPermissionGroupById")
	defer span.End()

	if id == "" {
		s.logger.Errorf("[PermSvc][GetGroupById] id empty")
		span.AddEvent("Id is empty", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionGroupRequired.Code),
		))
		return nil, errors.New(localization.ErrorPermissionGroupRequired.Code)
	}

	group, err := s.repo.GetPermissionGroupById(ctx, id)
	if err != nil {
		s.logger.Errorf("[PermSvc][GetGroupById] fetch err: %v", err)
		span.AddEvent("Failed to fetch permission group", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	s.logger.Infof("[PermSvc][GetGroupById] retrieved id: %s", id)
	return group, nil
}

func (s *permissionService) GetPermissionGroups(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPermissionGroups", "Permission", "GetPermissionGroups")
	defer span.End()

	if filterParams == nil {
		filter := types.Filter{}
		filterParams = &filter
	}

	result, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("[PermSvc][GetGroups] fetch err: %v", err)
		span.AddEvent("Failed to fetch permission groups", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	s.logger.Infof("[PermSvc][GetGroups] retrieved %d", len(result.Data))
	return result, nil
}

func (s *permissionService) GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllPermissionCategoriesWithPermissions", "Permission", "GetAllPermissionCategoriesWithPermissions")
	defer span.End()

	result, err := s.repo.GetAllPermissionCategoriesWithPermissions(ctx)
	if err != nil {
		span.AddEvent("Failed to get permission categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

func (s *permissionService) ValidatePermissionCategories(ctx context.Context, categoryIDs []string) (bool, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "ValidatePermissionCategories", "Permission", "ValidatePermissionCategories")
	defer span.End()

	validCategories, err := s.repo.ValidatePermissionCategories(ctx, categoryIDs)
	if err != nil {
		s.logger.Errorf("[PermSvc][ValidateCategories] validate err: %v", err)
		span.AddEvent("Failed to validate permission categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return false, err
	}

	// Check if all requested categories were validated
	if len(validCategories) != len(categoryIDs) {
		s.logger.Warnf("[PermSvc][ValidateCategories] mismatch requested: %d, valid: %d", len(categoryIDs), len(validCategories))
		span.AddEvent("Permission category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionCategoryNotFound.Code),
		))
		return false, errors.New(localization.ErrorPermissionCategoryNotFound.Code)
	}

	return true, nil
}

func (s *permissionService) ValidatePermissionGroups(ctx context.Context, groupIDs []string) (bool, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "ValidatePermissionGroups", "Permission", "ValidatePermissionGroups")
	defer span.End()

	validGroups, err := s.repo.ValidatePermissionGroups(ctx, groupIDs)
	if err != nil {
		s.logger.Errorf("[PermSvc][ValidateGroups] validate err: %v", err)
		span.AddEvent("Failed to validate permission groups", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionGroupValidationFailed.Code),
		))
		return false, errors.New(localization.ErrorPermissionGroupValidationFailed.Code)
	}

	if len(validGroups) != len(groupIDs) {
		s.logger.Warnf("[PermSvc][ValidateGroups] mismatch requested: %d, valid: %d", len(groupIDs), len(validGroups))
		span.AddEvent("Permission group not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionGroupNotFound.Code),
		))
		return false, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	return true, nil
}

func (s *permissionService) GetPermissionCategoriesByDepartment(ctx context.Context, departmentId string) (map[string][]*model.PermissionCategory, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPermissionCategoriesByDepartment", "Permission", "GetPermissionCategoriesByDepartment")
	defer span.End()

	department, err := s.department.FindByID(ctx, departmentId)
	if err != nil || department == nil {
		s.logger.Errorf("[PermSvc][GetCatsByDept] dept not found: %s", departmentId)
		span.AddEvent("Department not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("department_id", departmentId),
		))
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	portal_cards := department.PortalCards
	cardsWithPermission := make(map[string][]*model.PermissionCategory)

	if len(portal_cards) > 0 {
		for _, card := range portal_cards {
			categories, err := s.repo.GetAllPermissionCategories(ctx, card)
			if err != nil {
				s.logger.Errorf("[PermSvc][GetCatsByDept] category not found for dept: %s", departmentId)
				span.AddEvent("Permission category not found", trace.WithAttributes(
					attribute.String("error", localization.ErrorResourceNotFound.Code),
					attribute.String("department_id", departmentId),
				))
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
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPermissionGroupsByDepartment", "Permission", "GetPermissionGroupsByDepartment")
	defer span.End()

	department, err := s.repo.FindAllGroupsWithPagination(ctx, departmentId, filterParam)
	if err != nil || department == nil {
		s.logger.Errorf("[PermSvc][GetGroupsByDept] dept not found: %s", departmentId)
		span.AddEvent("Department not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("department_id", departmentId),
		))
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	return department, nil
}

func (s *permissionService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Permission", "Authorize")
	defer span.End()

	switch action.ActionType {
	case string(constants.CREATE):
		cur, err := core.BindPermissionGroupFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind permission group from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Create(ctx, &cur); err != nil {
			span.AddEvent("Failed to create permission group", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		return action, nil

	case string(constants.UPDATE):
		upd, err := core.BindPermissionGroupFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind permission group from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		existingGroup, err := s.repo.GetPermissionGroupById(ctx, action.UniqueId)
		if err != nil {
			span.AddEvent("Resource not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorResourceNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}

		// Update using the existing group's ObjectID
		err = s.repo.Update(ctx, existingGroup.ID.Hex(), &upd)
		if err != nil {
			span.AddEvent("Failed to update permission group", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		return action, nil

	default:
		span.AddEvent("Unhandled action type", trace.WithAttributes(
			attribute.String("error", "UNHANDLED_ACTION_TYPE"),
			attribute.String("action_type", string(action.ActionType)),
		))
		return nil, errors.New("UNHANDLED_ACTION_TYPE")
	}
}

func (s *permissionService) GetPopulatedPermissionCategories(ctx context.Context, categoryIDsObject []bson.ObjectID) ([]cps_user_dto.PermissionCategoryResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPopulatedPermissionCategories", "Permission", "GetPopulatedPermissionCategories")
	defer span.End()

	if len(categoryIDsObject) == 0 {
		return []cps_user_dto.PermissionCategoryResponse{}, nil
	}
	var categoryIDs []string
	for _, id := range categoryIDsObject {
		categoryIDs = append(categoryIDs, id.Hex())
	}

	validCategories, err := s.repo.ValidatePermissionCategories(ctx, categoryIDs)
	if err != nil {
		s.logger.Errorf("[PermSvc][GetPopulatedCats] validate err: %v", err)
		span.AddEvent("Failed to validate permission categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, errors.New(localization.ErrorPermissionCategoryNotFound.Code)
	}

	if len(validCategories) != len(categoryIDs) {
		s.logger.Warnf("[PermSvc][GetPopulatedCats] mismatch requested: %d, valid: %d", len(categoryIDs), len(validCategories))
		span.AddEvent("Permission category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionCategoryNotFound.Code),
		))
		return nil, errors.New(localization.ErrorPermissionCategoryNotFound.Code)
	}

	result, err := s.repo.GetPopulatedPermissionCategories(ctx, categoryIDs)
	if err != nil {
		span.AddEvent("Failed to get populated permission categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

func (s *permissionService) GetPopulatedPermissionGroups(ctx context.Context, groupIDsObject []bson.ObjectID) ([]cps_user_dto.PermissionGroupResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPopulatedPermissionGroups", "Permission", "GetPopulatedPermissionGroups")
	defer span.End()

	if len(groupIDsObject) == 0 {
		return []cps_user_dto.PermissionGroupResponse{}, nil
	}
	var groupIDs []string
	for _, id := range groupIDsObject {
		groupIDs = append(groupIDs, id.Hex())
	}

	validGroups, err := s.repo.ValidatePermissionGroups(ctx, groupIDs)
	if err != nil {
		s.logger.Errorf("[PermSvc][GetPopulatedGroups] validate err: %v", err)
		span.AddEvent("Failed to validate permission groups", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	if len(validGroups) != len(groupIDs) {
		s.logger.Warnf("[PermSvc][GetPopulatedGroups] mismatch requested: %d, valid: %d", len(groupIDs), len(validGroups))
		span.AddEvent("Permission group not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorPermissionGroupNotFound.Code),
		))
		return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}

	result, err := s.repo.GetPopulatedPermissionGroups(ctx, groupIDs)
	if err != nil {
		span.AddEvent("Failed to get populated permission groups", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}
