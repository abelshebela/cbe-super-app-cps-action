package core

import (
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

// ConvertToDTO converts a CPSUser model to CPSUserDTO
func ConvertToDTO(user *model.CPSUser) *cpsuser.CPSUserDTO {
	return &cpsuser.CPSUserDTO{
		ID:                 user.ID,
		UserCode:           user.UserCode,
		FullName:           user.FullName,
		Role:               user.Role,
		Department:         user.Department,
		Gender:             user.Gender,
		PhoneNumber:        user.PhoneNumber,
		Email:              user.Email,
		UserName:           user.UserName,
		Realm:              user.Realm,
		PermissionCategory: user.PermissionCategory,
		PermissionGroup:    user.PermissionGroup,
		Enabled:            user.Enabled,
		DateJoined:         user.DateJoined,
		LastModified:       user.LastModified,
		Country:            user.Country,
		Region:             user.Region,
	}
}

func CPSUModel(req cpsuser.CreateUserRequest) model.CPSUser {
	return model.CPSUser{
		UserCode:           local_util.GenerateCPSUserCode(),
		UserName:           req.UserName,
		FullName:           req.FullName,
		Department:         req.Department,
		PhoneNumber:        req.PhoneNumber,
		Role:               req.Role,
		Gender:             req.Gender,
		Email:              req.Email,
		PermissionCategory: req.PermissionCategory,
		PermissionGroup:    req.PermissionGroups,
		PasswordDisable:    true,
	}
}

func CPSUUpdateModel(req cpsuser.UpdateUserRequest) *model.CPSUser {
	return &model.CPSUser{
		FullName:           req.FullName,
		Role:               req.Role,
		Department:         req.Department,
		Gender:             req.Gender,
		PhoneNumber:        req.PhoneNumber,
		Email:              req.Email,
		UserName:           req.UserName,
		PermissionCategory: req.PermissionCategory,
		PermissionGroup:    req.PermissionGroups,
	}
}