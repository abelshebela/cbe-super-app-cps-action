package cps_user

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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
