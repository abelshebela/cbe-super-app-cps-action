package core

import (
	"strings"

	"cbe-super-app-cps-action/internal/constants/dto/permission"
	"cbe-super-app-cps-action/internal/constants/model"
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
