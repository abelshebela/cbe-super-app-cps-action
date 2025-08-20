package permission

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type PermissionGroupRepository interface {
	CheckPermissionGroupExists(groupName string) bool
	CreatePermissionGroupFromAction(ctx context.Context, action *cps_entities.CPSAction) (*cps_entities.CPSAction, error)
	UpdatePermissionGroupFromAction(ctx context.Context, action *cps_entities.CPSAction) (*cps_entities.CPSAction, error)
	UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error)
	GetPermissionGroup(groupName string) (entities.PermissionGroup, error)
	GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error)
	ValidatePermissionGroups(ctx context.Context, ids []string) ([]string, error)
}

type PermissionCategoryRepository interface {
	ValidatePermissionCategories(ctx context.Context, ids []string) ([]string, error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*entities.PermissionCategory, error)
}

type CPSActionRepository interface {
	CheckPermissionGroupExists(groupName string) bool
	CheckPendingRequest(makerId string, actionStatus model.ActionStatus, requestAction model.RequestAction) error
	CreatePermissionGroup(action model.CPSAction) (model.CPSAction, error)
	// ApproveActionRequest(actionCode string, action model.CPSAction) (model.CPSAction, error)
	ValidateActionRequest(actionCode, department string) (model.CPSAction, error)
	RejectActionRequest(actionCode string, action model.CPSAction, reason string) (model.CPSAction, error)
	UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error)
	GetPermissionGroup(groupName string) (entities.PermissionGroup, error)
	GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error)
}
