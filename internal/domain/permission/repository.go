package permission

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"
)

type PermissionGroupRepository interface {
	CheckPermissionGroupExists(groupName string) bool
	CreatePermissionGroupFromAction(action model.CPSAction) error
	UpdatePermissionGroupFromAction(action model.CPSAction) error
}

type PermissionCategoryRepository interface {
	ValidatePermissionCategories(ids []string) ([]string, error)
}

type CPSActionRepository interface {
	CheckPendingRequest(makerId string, actionStatus model.ActionStatus, requestAction model.RequestAction) error
	CreatePermissionGroup(action model.CPSAction) (model.CPSAction, error)
	ApproveActionRequest(actionCode string, action model.CPSAction) error
	ValidateActionRequest(actionCode, department string) (model.CPSAction, error)
}
