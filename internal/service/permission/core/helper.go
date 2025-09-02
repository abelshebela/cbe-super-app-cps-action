package core

import (
	"strings"

	"cbe-super-app-cps-action/internal/constants/dto/permission"
	"cbe-super-app-cps-action/internal/constants/model"
	"encoding/json"
)

// PermissionGroupModel builds model.PermissionGroup from request
func PermissionGroupModel(req permission.CreatePermissionGroupRequest) model.PermissionGroup {
	return model.PermissionGroup{
		GroupName:          strings.ToUpper(req.GroupName),
		Role:               req.Role,
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
		Role:               req.Role,
		PermissionCategory: req.PermissionCategoryLists,
		Realm:              "bank",
		Enabled:            true,
		IsDeleted:          false,
	}
}

// UnmarshalActionToType marshals and unmarshals action.CurrentAction to the target type
// func UnmarshalActionToType[T any](currentAction interface{}) (T, error) {
// 	var result T
// 	bytes, err := json.Marshal(currentAction)
// 	if err != nil {
// 		return result, err
// 	}
// 	if err := json.Unmarshal(bytes, &result); err != nil {
// 		return result, err
// 	}
// 	return result, nil

// }

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
