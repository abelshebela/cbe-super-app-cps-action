package core

import (
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
)

// UsernameExists checks if a username already exists in the database
func UsernameExists(ctx context.Context, repo storage.CpsUserRepository, username string) (bool, error) {
	user, err := repo.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func EmailExists(ctx context.Context, repo storage.CpsUserRepository, email string) (bool, error) {
	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
func PhoneNumberExists(ctx context.Context, repo storage.CpsUserRepository, phoneNumber string) (bool, error) {
	user, err := repo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

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
		PasswordDisable:    false,
		IsFirstTimeLogin:   true,
		Enabled:            true,
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

// BindCPSUserFromAction decodes action.CurrentAction into model.CPSUser
func BindCPSUserFromAction(currentAction interface{}) (model.CPSUser, error) {
	var user model.CPSUser

	// Fast-path if already the correct type
	if v, ok := currentAction.(model.CPSUser); ok {
		return v, nil
	}

	// If stored as JSON string
	if s, ok := currentAction.(string); ok {
		if err := json.Unmarshal([]byte(s), &user); err == nil {
			return user, nil
		}
	}

	// Generic path: marshal then unmarshal
	bytes, err := json.Marshal(currentAction)
	if err != nil {
		return user, err
	}

	// Try to unmarshal as CPSUserActionPayload first
	var payload cpsuser.CPSUserActionPayload
	if err := json.Unmarshal(bytes, &payload); err == nil && payload.User != nil {
		// Extract user from payload
		userBytes, err := json.Marshal(payload.User)
		if err != nil {
			return user, err
		}
		if err := json.Unmarshal(userBytes, &user); err != nil {
			return user, err
		}
		return user, nil
	}

	// Fallback: try to unmarshal directly as CPSUser
	if err := json.Unmarshal(bytes, &user); err != nil {
		return user, err
	}
	return user, nil
}

// BindCPSUserUpdateFromAction decodes action.CurrentAction into cpsuser.UpdateUserRequest
func BindCPSUserUpdateFromAction(currentAction interface{}) (cpsuser.UpdateUserRequest, error) {
	var updateReq cpsuser.UpdateUserRequest

	// Fast-path if already the correct type
	if v, ok := currentAction.(cpsuser.UpdateUserRequest); ok {
		return v, nil
	}

	// If stored as JSON string
	if s, ok := currentAction.(string); ok {
		if err := json.Unmarshal([]byte(s), &updateReq); err == nil {
			return updateReq, nil
		}
	}

	// Generic path: marshal then unmarshal
	bytes, err := json.Marshal(currentAction)
	if err != nil {
		return updateReq, err
	}
	if err := json.Unmarshal(bytes, &updateReq); err != nil {
		return updateReq, err
	}
	return updateReq, nil
}
