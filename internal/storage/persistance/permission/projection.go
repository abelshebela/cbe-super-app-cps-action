package permission

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func PermissionGroupUpdateMapper(permissionGroup *model.PermissionGroup) bson.M {
	set := bson.M{}
	if permissionGroup.GroupName != "" {
		set["group_name"] = permissionGroup.GroupName
	}
	if permissionGroup.Role != "" {
		set["role"] = permissionGroup.Role
	}
	if permissionGroup.PermissionCategory != nil {
		set["permission_category_list"] = permissionGroup.PermissionCategory
	}
	set["last_modified"] = time.Now()
	return bson.M{"$set": set}
}
func PermissionCategoryUpdateMapper(permissionCategory *model.PermissionCategory) bson.M {
	set := bson.M{}
	if permissionCategory.CategoryName != "" {
		set["category_name"] = permissionCategory.CategoryName
	}
	if permissionCategory.Access != "" {
		set["access"] = permissionCategory.Access
	}
	if permissionCategory.Permissions != nil {
		set["permissions"] = permissionCategory.Permissions
	}
	set["updatedAt"] = time.Now()
	return bson.M{"$set": set}
}
