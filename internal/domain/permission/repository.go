package permission

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/permission/entities"
)

type PermissionGroupRepository interface {
	CheckPermissionGroupExists(groupName string) bool
	CreatePermissionGroupFromAction(action entities.CPSAction) error
	UpdatePermissionGroupFromAction(action entities.CPSAction) error
}

type PermissionCategoryRepository interface {
	ValidatePermissionCategories(ids []string) ([]string, error)
}

type CPSActionRepository interface {
	CheckPendingRequest(makerId string, actionStatus entities.ActionStatus, requestAction entities.RequestAction) error
	CreatePermissionGroup(action entities.CPSAction) error
	ApproveActionRequest(actionCode string, action entities.CPSAction) error
	ValidateActionRequest(actionCode, department string) (entities.CPSAction, error)
}
