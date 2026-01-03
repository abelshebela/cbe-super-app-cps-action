package core

import (
	"strings"

	"cbe-super-app-cps-action/internal/constants/dto/permission"
	"encoding/json"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// PermissionGroupModel builds model.PermissionGroup from request
func PermissionGroupModel(req permission.CreatePermissionGroupRequest) model.PermissionGroup {
	return model.PermissionGroup{
		GroupName:          strings.ToUpper(req.GroupName),
		Role:               strings.ToUpper(req.Role),
		DepartmentID:       req.DepartmentID,
		PermissionCategory: req.PermissionCategoryLists,
		Realm:              "bank",
		Enabled:            true,
		IsDeleted:          false,
	}
}

// PermissionGroupUpdateModel builds model.PermissionGroup for updates
func PermissionGroupUpdateModel(req permission.UpdatePermissionGroupRequest) model.PermissionGroup {
	return model.PermissionGroup{
		GroupName:          strings.ToUpper(req.NewGroupName),
		Role:               strings.ToUpper(req.Role),
		DepartmentID:       req.DepartmentID,
		PermissionCategory: req.PermissionCategoryLists,
		Realm:              "bank",
		Enabled:            true,
		IsDeleted:          false,
	}
}

// BindPermissionGroupFromAction decodes action.CurrentAction into model.PermissionGroup
func BindPermissionGroupFromAction(currentAction interface{}) (model.PermissionGroup, error) {
	var pg model.PermissionGroup

	// Fast-path if already the correct type
	if v, ok := currentAction.(model.PermissionGroup); ok {
		return v, nil
	}

	// If stored as JSON string
	if s, ok := currentAction.(string); ok {
		if err := json.Unmarshal([]byte(s), &pg); err == nil {
			return pg, nil
		}
	}

	// Generic path: marshal then unmarshal
	bytes, err := json.Marshal(currentAction)
	if err != nil {
		return pg, err
	}
	if err := json.Unmarshal(bytes, &pg); err != nil {
		return pg, err
	}
	return pg, nil
}
