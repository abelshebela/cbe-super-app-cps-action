package permission

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
		set["permission_category"] = permissionGroup.PermissionCategory
	}
	set["last_modified"] = time.Now()
	return set
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
	return set
}

func PermissionGroupsPipeline(filter bson.M, skip int64, limit int64) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: 1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	return pipeline
}
