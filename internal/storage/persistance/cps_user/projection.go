package cps_user

import (
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CPSUserUpdateMapper(u *imodel.CPSUser) bson.M {
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

		// 1️⃣ Match CPS user by user_code
		bson.D{{Key: "$match", Value: bson.M{
			"user_code": userCode,
		}}},

		// 2️⃣ Lookup role by job_title
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "roles",
			"let": bson.M{
				"jobTitle": "$job_title",
			},
			"pipeline": mongo.Pipeline{

				// match roles.job_title == user.job_title
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$eq": []interface{}{"$job_title", "$$jobTitle"},
					},
				}}},

				// lookup job_roles using role code
				bson.D{{Key: "$lookup", Value: bson.M{
					"from": "job_roles",
					"let": bson.M{
						"roleCode": "$role",
					},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$code", "$$roleCode"},
							},
						}}},
					},
					"as": "job_role",
				}}},

				// unwind job_role
				bson.D{{Key: "$unwind", Value: bson.M{
					"path":                       "$job_role",
					"preserveNullAndEmptyArrays": false,
				}}},

				// project role_id from job_roles._id
				bson.D{{Key: "$project", Value: bson.M{
					"_id":     0,
					"role_id": "$job_role._id",
				}}},
			},
			"as": "role_doc",
		}}},

		// 3️⃣ Unwind role_doc
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$role_doc",
			"preserveNullAndEmptyArrays": false,
		}}},

		// 4️⃣ Final projection
		bson.D{{Key: "$project", Value: bson.M{
			"_id":          0,
			"user_code":    1,
			"role":         1,
			"full_name":    1,
			"username":     1,
			"email":        1,
			"phone_number": 1,
			"gender":       1,
			"realm":        1,
			"enabled":      1,
			"job_title":    1,
			"role_id":      "$role_doc.role_id",
		}}},
	}
}

func CPSUserMapper(u imodel.CPSUser) *model.CPSUser {
	return &model.CPSUser{
		UserCode:           u.UserCode,
		FullName:           u.FullName,
		Gender:             u.Gender,
		PhoneNumber:        u.PhoneNumber,
		Email:              u.Email,
		UserName:           u.UserName,
		Realm:              u.Realm,
		PermissionCategory: u.PermissionCategory,
		PermissionGroup:    u.PermissionGroup,
		JobTitle:           u.JobTitle,
		PasswordDisable:    u.PasswordDisable,
		SyncDisabled:       u.SyncDisabled,
		IsFirstTimeLogin:   u.IsFirstTimeLogin,
		Enabled:            true,
		IsDeleted:          u.IsDeleted,
		DateJoined:         u.DateJoined,
		LastModified:       u.LastModified,
	}
}

func PipelineBuilderWithRole(userCode, departmentColl, rolesColl, jobRolesColl string) mongo.Pipeline {
	return mongo.Pipeline{

		// 1️⃣ Match CPS user by user_code
		bson.D{{Key: "$match", Value: bson.M{
			"user_code": userCode,
		}}},

		// 2️⃣ Lookup role by job_title
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": rolesColl,
			"let": bson.M{
				"jobTitle": "$job_title",
			},
			"pipeline": mongo.Pipeline{

				// match roles.job_title == user.job_title
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$eq": []interface{}{"$job_title", "$$jobTitle"},
					},
				}}},

				// lookup job_roles using role code
				bson.D{{Key: "$lookup", Value: bson.M{
					"from": jobRolesColl,
					"let": bson.M{
						"roleCode": "$role",
					},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$code", "$$roleCode"},
							},
						}}},
					},
					"as": "job_role",
				}}},

				// unwind job_role
				bson.D{{Key: "$unwind", Value: bson.M{
					"path":                       "$job_role",
					"preserveNullAndEmptyArrays": false,
				}}},

				// project role_id and name from job_roles
				bson.D{{Key: "$project", Value: bson.M{
					"_id":     0,
					"role_id": "$job_role._id",
					"name":    "$job_role.name", // Assuming 'name' is the field in job_roles
				}}},
			},
			"as": "role_doc",
		}}},

		// 3️⃣ Unwind role_doc
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$role_doc",
			"preserveNullAndEmptyArrays": true,
		}}},

		// 4️⃣ Final projection
		bson.D{{Key: "$project", Value: bson.M{
			"_id":       1,
			"user_code": 1,
			"role": bson.M{
				"code": "$role",          // Role code from user
				"name": "$role_doc.name", // Role name from job_roles via role_doc
			},
			"full_name":    1,
			"username":     1,
			"email":        1,
			"phone_number": 1,
			"gender":       1,
			"realm":        1,
			"enabled":      1,
			"job_title":    1,
			"role_id":      "$role_doc.role_id",
		}}},
	}
}

func CPSUserPopulated(u imodel.CPSUser) *model.CPSUser {
	return &model.CPSUser{
		UserCode:           u.UserCode,
		FullName:           u.FullName,
		Gender:             u.Gender,
		PhoneNumber:        u.PhoneNumber,
		Email:              u.Email,
		UserName:           u.UserName,
		Realm:              u.Realm,
		PermissionCategory: u.PermissionCategory,
		PermissionGroup:    u.PermissionGroup,
		JobTitle:           u.JobTitle,
		PasswordDisable:    u.PasswordDisable,
		SyncDisabled:       u.SyncDisabled,
		IsFirstTimeLogin:   u.IsFirstTimeLogin,
		Enabled:            true,
		IsDeleted:          u.IsDeleted,
		DateJoined:         u.DateJoined,
		LastModified:       u.LastModified,
	}
}
