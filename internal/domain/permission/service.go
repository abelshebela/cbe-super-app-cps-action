package permission

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PermissionDomainService interface {
	CreatePermissionGroup(ctx context.Context, oldGroupName, groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error)
	Authorize(ctx context.Context, action *cps_entities.CPSAction) (*cps_entities.CPSAction, error)
	UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error)
	GetPermissionGroup(groupName string) (entities.PermissionGroup, error)
	GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error)
	ValidatePermissionCategories(ctx context.Context, ids []string) ([]string, error)
	ValidatePermissionGroups(ctx context.Context, ids []string) ([]string, error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*entities.PermissionCategory, error)
}

type Service struct {
	cpsActionRepo          CPSActionRepository
	permissionGroupRepo    PermissionGroupRepository
	permissionCategoryRepo PermissionCategoryRepository
	logger                 utils.Logger
}

func InitPermissionDomain(cpsActionRepo CPSActionRepository, permissionGroupRepo PermissionGroupRepository, permissionCategoryRepo PermissionCategoryRepository, logger utils.Logger) Service {
	return Service{
		cpsActionRepo:          cpsActionRepo,
		permissionGroupRepo:    permissionGroupRepo,
		permissionCategoryRepo: permissionCategoryRepo,
		logger:                 logger,
	}
}

func (s Service) CreatePermissionGroup(ctx context.Context, oldGroupName, groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error) {
	oldGroupName = strings.ToUpper(oldGroupName)
	groupName = strings.ToUpper(groupName)
	if err := s.cpsActionRepo.CheckPendingRequest(cpsAction.MakerID, model.ActionStatus(cpsAction.ActionStatus), model.RequestAction(cpsAction.RequestAction)); err != nil {
		s.logger.Warnf("Pending request check failed for user %s: %v", cpsAction.MakerID, err)
		return model.CPSAction{}, fmt.Errorf("PENDING_REQUEST_EXISTS")
	}

	if exists := s.permissionGroupRepo.CheckPermissionGroupExists(groupName); exists {
		s.logger.Warnf("Permission group already exists with name: %s", groupName)
		return model.CPSAction{}, fmt.Errorf("PERMISSION_GROUP_ALREADY_EXIXTS")
	}

	// For UPDATE operations, only validate permission categories if they are being changed
	var validCategories []string
	var err error
	if cpsAction.ActionType == string(entities.ActionUpdate) && oldGroupName != "" {
		// This is an update operation - check if permission categories are being changed
		if len(permissionCategoryLists) > 0 {
			// Permission categories are being updated, so validate them
			s.logger.Infof("Validating permission categories for update: %+v", permissionCategoryLists)
			validCategories, err = s.permissionCategoryRepo.ValidatePermissionCategories(ctx, permissionCategoryLists)
			if err != nil {
				s.logger.Errorf("Failed to validate permission categories: %v", err)
				return model.CPSAction{}, err
			}
			s.logger.Infof("Permission categories validated successfully: %+v", validCategories)
		} else {
			// Permission categories are not being changed, use existing ones without validation
			s.logger.Infof("No permission categories provided for update, keeping existing ones unchanged")
			validCategories = permissionCategoryLists // This will be empty, indicating no change
		}
	} else {
		// This is a create operation - always validate permission categories
		s.logger.Infof("Validating permission categories for create: %+v", permissionCategoryLists)
		validCategories, err = s.permissionCategoryRepo.ValidatePermissionCategories(ctx, permissionCategoryLists)
		if err != nil {
			s.logger.Errorf("Failed to validate permission categories: %v", err)
			return model.CPSAction{}, err
		}
		s.logger.Infof("Permission categories validated successfully: %+v", validCategories)
	}

	cpsAction.ActionCode = utils.RandomGenerator(20)
	if cpsAction.ActionType == string(entities.ActionUpdate) {
		cpsAction.CurrentAction = map[string]interface{}{
			"group_name":          groupName,
			"permission_category": validCategories,
			"role":                role,
			"realm":               "bank",
			"old_group":           oldGroupName,
		}
	} else {
		cpsAction.CurrentAction = map[string]interface{}{
			"group_name":          groupName,
			"permission_category": validCategories,
			"role":                role,
			"realm":               "bank",
		}
	}

	cpsAction.MakerActionTime = time.Now()
	cpsAction.CreatedAt = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	cpsAction, err = s.cpsActionRepo.CreatePermissionGroup(cpsAction)
	if err != nil {
		s.logger.Errorf("Failed to create CPSAction: %v", err)
		return model.CPSAction{}, err
	}

	s.logger.Infof("Successfully created permission group request with ActionCode: %s", cpsAction.ActionCode)
	return cpsAction, nil
}

func (s Service) Authorize(ctx context.Context, action *cps_entities.CPSAction) (*cps_entities.CPSAction, error) {

	switch action.ActionType {
	case "CREATE":
		return s.permissionGroupRepo.CreatePermissionGroupFromAction(ctx, action)
	case "UPDATE":
		return s.permissionGroupRepo.UpdatePermissionGroupFromAction(ctx, action)
	default:
		return nil, errors.New("invalid action type")
	}

}

func (s Service) UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error) {
	groupName = strings.ToUpper(groupName)
	return s.permissionGroupRepo.UpdatePermissionGroup(groupName, permissionCategoryLists)
}

func (s Service) GetPermissionGroup(groupName string) (entities.PermissionGroup, error) {
	groupName = strings.ToUpper(groupName)
	return s.permissionGroupRepo.GetPermissionGroup(groupName)
}

func (s Service) GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error) {
	return s.permissionGroupRepo.GetPermissionGroups(ctx, filterParams)
}

func (s Service) ValidatePermissionCategories(ctx context.Context, ids []string) ([]string, error) {
	return s.permissionCategoryRepo.ValidatePermissionCategories(ctx, ids)
}

func (s Service) ValidatePermissionGroups(ctx context.Context, ids []string) ([]string, error) {
	return s.permissionGroupRepo.ValidatePermissionGroups(ctx, ids)
}

func (s Service) GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*entities.PermissionCategory, error) {
	return s.permissionCategoryRepo.GetAllPermissionCategoriesWithPermissions(ctx)
}
