package cps_user

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CPSUserUpdateMapper(u *model.CPSUser) bson.M {
	set := bson.M{}
	if u.FullName != "" {
		set["full_name"] = u.FullName
	}
	if u.Role != "" {
		set["role"] = u.Role
	}
	if !u.Department.IsZero() {
		set["department"] = u.Department
	}
	if u.Gender != "" {
		set["gender"] = u.Gender
	}
	if u.PhoneNumber != "" {
		set["phone_number"] = u.PhoneNumber
	}
	if u.Email != "" {
		set["email"] = u.Email
	}
	if u.UserName != "" {
		set["username"] = u.UserName
	}
	if u.PermissionCategory != nil {
		set["permission_category"] = u.PermissionCategory
	}
	if u.PermissionGroup != nil {
		set["permission_group"] = u.PermissionGroup
	}

	set["last_modified"] = time.Now()
	return set
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

func PermissionCategoryProjection(permissionColl, permissionCategoryColl string) bson.D {
	return bson.D{{Key: "$lookup", Value: bson.M{
		"from": permissionCategoryColl,
		"let":  bson.M{"categoryIds": "$permission_category"},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
				"$in": []interface{}{"$_id", IdConverter("$$categoryIds")},
			}}}},
			bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
			bson.D{{Key: "$project", Value: bson.M{
				"category_name": 1,
				"access":        1,
				"permissions":   1,
			}}},
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": permissionColl,
				"let":  bson.M{"permissionIds": "$permissions"},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
						"$in": []interface{}{"$_id", IdConverter("$$permissionIds")},
					}}}},
					bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
					bson.D{{Key: "$project", Value: bson.M{"permission_name": 1}}},
				},
				"as": "permissions_docs",
			}}},
		},
		"as": "permission_category_docs",
	}}}

}

func PipelineBuilder(userCode string, departmentColl, permissionGroupColl, permissionColl, permissionCategoryColl string) mongo.Pipeline {
	return mongo.Pipeline{
		// Match the user by user_code and ensure it's not deleted
		bson.D{{Key: "$match", Value: bson.M{"user_code": userCode, "is_deleted": false}}},

		// Lookup department information
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         departmentColl,
			"localField":   "department",
			"foreignField": "_id",
			"as":           "department_doc",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
				bson.D{{Key: "$project", Value: bson.M{
					"_id":          1,
					"department":   1,
					"portal_cards": 1,
				}}},
			},
		}}},

		// Unwind department (preserve null for users without department)
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$department_doc", "preserveNullAndEmptyArrays": true}}},

		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         permissionGroupColl,
			"localField":   "permission_group",
			"foreignField": "_id",
			"as":           "permission_groups_raw",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
				bson.D{{Key: "$project", Value: bson.M{
					"group_name":          1,
					"permission_category": 1,
				}}},
				PermissionCategoryProjection(permissionColl, permissionCategoryColl),
			},
		}}},

		bson.D{{Key: "$project", Value: bson.M{
			"_id":           1,
			"user_code":     1,
			"full_name":     1,
			"role":          1,
			"gender":        1,
			"phone_number":  1,
			"email":         1,
			"username":      1,
			"realm":         1,
			"enabled":       1,
			"date_joined":   1,
			"last_modified": 1,
			"country":       1,
			"region":        1,
			"department": bson.M{
				"$cond": bson.M{
					"if": bson.M{"$ne": []interface{}{"$department_doc", nil}},
					"then": bson.M{
						"id":           "$department_doc._id",
						"name":         "$department_doc.department",
						"portal_cards": "$department_doc.portal_cards",
					},
					"else": nil,
				},
			},
			// Permission groups with nested structure
			"permission_groups": bson.M{
				"$map": bson.M{
					"input": "$permission_groups_raw",
					"as":    "pg",
					"in": bson.M{
						"id":         "$$pg._id",
						"group_name": "$$pg.group_name",
						"permission_category": bson.M{
							"$map": bson.M{
								"input": "$$pg.permission_category_docs",
								"as":    "cat",
								"in": bson.M{
									"id":            "$$cat._id",
									"category_name": "$$cat.category_name",
									"access":        "$$cat.access",
									"permissions": bson.M{
										"$map": bson.M{
											"input": "$$cat.permissions_docs",
											"as":    "perm",
											"in": bson.M{
												"id":              "$$perm._id",
												"permission_name": "$$perm.permission_name",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}}},
	}
}
