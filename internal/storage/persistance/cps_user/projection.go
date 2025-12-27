package cps_user

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CPSUserUpdateMapper(u *model.CPSUser) bson.M {
	set := bson.M{}
	if u.UserName != "" {
		set["username"] = u.UserName
	}
	if u.FullName != "" {
		set["full_name"] = u.FullName
	}
	if u.PhoneNumber != "" {
		set["phone_number"] = u.PhoneNumber
	}
	if u.Gender != "" {
		set["gender"] = u.Gender
	}
	if u.Email != "" {
		set["email"] = u.Email
	}
	if u.JobTitle != "" {
		set["job_title"] = u.JobTitle
	}

	if u.JobTitle != "" {
		set["job_title"] = u.JobTitle
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

func PipelineBuilder(userCode string) mongo.Pipeline {
	return mongo.Pipeline{
		// Match the user
		bson.D{{Key: "$match", Value: bson.M{"user_code": userCode, "is_deleted": false}}},

		// Lookup roles by job_title
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "roles",
			"let":  bson.M{"job_title": "$job_title"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$eq": []interface{}{"$job_title", "$$job_title"}},
				}}},
				bson.D{{Key: "$project", Value: bson.M{
					"_id":  1,
					"role": 1,
				}}},
			},
			"as": "role_doc",
		}}},

		// Unwind role_doc
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$role_doc", "preserveNullAndEmptyArrays": true}}},

		// Lookup job_roles using role code
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"roleCode": "$role_doc.role"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$eq": []interface{}{"$code", "$$roleCode"}},
				}}},
				bson.D{{Key: "$project", Value: bson.M{
					"_id":          0,
					"portal_cards": 1,
				}}},
			},
			"as": "job_role_doc",
		}}},

		// Unwind job_role_doc
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$job_role_doc", "preserveNullAndEmptyArrays": true}}},

		// Lookup department document
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "departments",
			"localField":   "department",
			"foreignField": "_id",
			"as":           "department_doc",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$department_doc", "preserveNullAndEmptyArrays": true}}},

		// Project final user fields with populated portal_cards and department
		bson.D{{Key: "$project", Value: bson.M{
			"_id":                 0,
			"user_code":           1,
			"full_name":           1,
			"username":            1,
			"job_title":           1,
			"gender":              1,
			"phone_number":        1,
			"email":               1,
			"realm":               1,
			"enabled":             1,
			"permission_category": 1,
			"portal_cards":        "$job_role_doc.portal_cards",
			"last_modified":       1,
			"date_joined":         1,
		}}},
	}
}
