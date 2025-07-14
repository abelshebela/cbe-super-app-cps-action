package permission

import "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"

type PermissionRepository interface {
	CheckPendingRequest(userCode, action string) error
	CheckPermissionGroupExists(groupName string) bool
	CreatePermissionGroupFromAction(action entities.CPSAction) error
	UpdatePermissionGroupFromAction(action entities.CPSAction) error
}

type PermissionCategoryRepository interface {
	ValidatePermissionCategories(ids []string) ([]string, error)
}

type CPSActionRepository interface {
	CreatePermissionGroup(action entities.CPSAction) error
	ApproveActionRequest(actioncode string, action entities.CPSAction) error
	ValidateActionRequest(actionCode, department string) (entities.CPSAction, error)
	RejectActionRequest(actionCode string, action entities.CPSAction, reason string) (entities.CPSAction, error)
}
