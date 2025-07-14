package permission

import (
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type PermissionRepository interface {
	CheckPendingRequest(userCode, action string) error
	CheckPermissionGroupExists(groupName string) bool
	CreatePermissionGroupFromAction(action model.CPSAction) error
	UpdatePermissionGroupFromAction(action model.CPSAction) error
}

type PermissionCategoryRepository interface {
	ValidatePermissionCategories(ids []string) ([]string, error)
}

type CPSActionRepository interface {
	CreatePermissionGroup(action model.CPSAction) error
	ApproveActionRequest(actioncode string, action model.CPSAction) error
	ValidateActionRequest(actionCode, department string) (model.CPSAction, error)
	RejectActionRequest(actionCode string, action model.CPSAction, reason string) (model.CPSAction, error)
}
