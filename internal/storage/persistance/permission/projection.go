package permission

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

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
	if permissionGroup.DepartmentID != "" {
		set["department_id"] = permissionGroup.DepartmentID
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

		// ===== HYDRATE PERMISSION CATEGORIES =====
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "permission_categories",
				"localField":   "permission_category",
				"foreignField": "_id",
				"as":           "permission_category_docs",
			}},
		},
		{
			{Key: "$addFields", Value: bson.M{
				"permission_category": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$gt": bson.A{
								bson.M{"$size": "$permission_category_docs"},
								0,
							},
						},
						"then": "$permission_category_docs",
						"else": "$permission_category",
					},
				},
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"permission_category_docs": 0,
			}},
		},
		// ===== END HYDRATION =====

		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: 1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	return pipeline
}

func IdConverter(fieldName string) bson.M {
	return bson.M{
		"$map": bson.M{
			"input": fieldName,
			"as":    "id",
			"in": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$$id"}, "string"}},
					"then": bson.M{"$toObjectId": "$$id"},
					"else": "$$id",
				},
			},
		},
	}
}
