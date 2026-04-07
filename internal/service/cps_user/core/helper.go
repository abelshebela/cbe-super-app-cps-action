package core

import (
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UsernameExists(ctx context.Context, userCode string, repo storage.CpsUserRepository, username string) (bool, error) {
	user, err := repo.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	if userCode != "" && userCode == user.UserCode {
		return false, nil
	}
	return user != nil, nil
}

func EmailExists(ctx context.Context, user_code string, repo storage.CpsUserRepository, email string) (bool, error) {
	user, err := repo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	if user_code != "" && user_code == user.UserCode {
		return false, nil
	}

	return user != nil, nil
}

func PhoneNumberExists(ctx context.Context, user_code string, repo storage.CpsUserRepository, phoneNumber string) (bool, error) {
	user, err := repo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	if user_code != "" && user_code == user.UserCode {
		return false, nil
	}
	return user != nil, nil
}

// ConvertToDTO converts a CPSUser model to CPSUserDTO
func ConvertToDTO(portalCard []string, user *cpsuser.CpsUserResponse, makerAlloc, checkerAlloc, auditorAlloc []string) *cpsuser.CPSUserResponse {
	return &cpsuser.CPSUserResponse{
		ID:                 user.ID.Hex(),
		UserCode:           user.UserCode,
		FullName:           user.FullName,
		Role:               user.Role,
		Gender:             user.Gender,
		PhoneNumber:        user.PhoneNumber,
		Email:              user.Email,
		UserName:           user.UserName,
		Realm:              user.Realm,
		JobTitle:           user.JobTitle,
		MakerAllocations:   makerAlloc,
		CheckerAllocations: checkerAlloc,
		AuditorAllocations: auditorAlloc,
		PortalCards:        portalCard,
		Enabled:            user.Enabled,
		DateJoined:         &user.DateJoined,
		LastModified:       &user.LastModified,
		Country:            user.Country,
		Region:             user.Region,
	}
}

func ConvertToResponseDTO(portalCard []string, user *cpsuser.CpsUserPopulatedResponse, makerAlloc, checkerAlloc, auditorAlloc []string) *cpsuser.CpsUserPopulatedResponse {
	return &cpsuser.CpsUserPopulatedResponse{
		ID:       user.ID,
		UserCode: user.UserCode,
		FullName: user.FullName,
		Role: cpsuser.RoleResponse{
			Code: user.Role.Code,
			Name: user.Role.Name,
		},
		Department:         user.Department,
		RoleCode:           user.RoleCode,
		Gender:             user.Gender,
		PhoneNumber:        user.PhoneNumber,
		Email:              user.Email,
		UserName:           user.UserName,
		Realm:              user.Realm,
		JobTitle:           user.JobTitle,
		MakerAllocations:   makerAlloc,
		CheckerAllocations: checkerAlloc,
		AuditorAllocations: auditorAlloc,
		PortalCards:        portalCard,
		Enabled:            user.Enabled,
		DateJoined:         user.DateJoined,
		LastModified:       user.LastModified,
		LastLogin:          user.LastLogin,
		Country:            user.Country,
		Region:             user.Region,
	}
}

func CPSUModel(req cpsuser.CreateUserRequest) imodel.CPSUser {
	depID, err := bson.ObjectIDFromHex(req.Department)
	if err != nil {
		// Handle error appropriately
	}
	return imodel.CPSUser{
		UserCode:         local_util.GenerateCPSUserCode(),
		UserName:         req.UserName,
		FullName:         req.FullName,
		PhoneNumber:      req.PhoneNumber,
		Department:       depID,
		JobTitle:         req.JobTitle,
		Gender:           req.Gender,
		Email:            req.Email,
		PasswordDisable:  false,
		IsFirstTimeLogin: true,
		Enabled:          true,
		CreatedAt:        time.Now(),
	}
}

func CPSUUpdateModel(req cpsuser.UpdateUserRequest) *model.CPSUser {
	return &model.CPSUser{
		UserName:    req.UserName,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
		Email:       req.Email,
		JobTitle:    req.JobTitle,
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

	if err := json.Unmarshal(bytes, &user); err != nil {
		return user, err
	}
	return user, nil
}

func BindCPSUserUpdateFromAction(currentAction interface{}) (cpsuser.UpdateUserRequest, error) {
	var updateReq cpsuser.UpdateUserRequest

	if v, ok := currentAction.(cpsuser.UpdateUserRequest); ok {
		return v, nil
	}

	if s, ok := currentAction.(string); ok {
		if err := json.Unmarshal([]byte(s), &updateReq); err == nil {
			return updateReq, nil
		}
	}

	bytes, err := json.Marshal(currentAction)
	if err != nil {
		return updateReq, err
	}
	if err := json.Unmarshal(bytes, &updateReq); err != nil {
		return updateReq, err
	}
	return updateReq, nil
}

func MapForActionWithDepartment(user imodel.CPSUser, department *model.Department) cpsuser.CpsUserPopulatedResponse {
	var deptResp *cpsuser.DepartmentResponse
	if department != nil {
		deptResp = &cpsuser.DepartmentResponse{
			ID:   department.ID,
			Name: department.Department,
		}
	}

	return cpsuser.CpsUserPopulatedResponse{
		ID:               user.ID,
		UserCode:         user.UserCode,
		FullName:         user.FullName,
		Role:             cpsuser.RoleResponse{Name: user.Role},
		RoleCode:         user.Role,
		Department:       deptResp,
		JobTitle:         user.JobTitle,
		Gender:           user.Gender,
		PhoneNumber:      user.PhoneNumber,
		Email:            user.Email,
		UserName:         user.UserName,
		Realm:            user.Realm,
		Enabled:          user.Enabled,
		DateJoined:       derefTime(user.DateJoined),
		LastModified:     derefTime(user.LastModified),
		Country:          user.Country,
		Region:           user.Region,
		LastLogin:        user.LastLogin,
		PasswordDisable:  false,
		IsFirstTimeLogin: true,
		CreatedAt:        time.Now(),
	}
}

// Helper to safely dereference *time.Time to time.Time (zero if nil)
func derefTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

// MapFromPopulatedResponse maps CpsUserPopulatedResponse and DepartmentResponse to imodel.CPSUser
func MapFromPopulatedResponse(resp *cpsuser.CpsUserPopulatedResponse) *imodel.CPSUser {
	var deptID bson.ObjectID
	if resp.Department != nil {
		deptID = resp.Department.ID
	}
	return &imodel.CPSUser{
		ID:               resp.ID,
		UserCode:         resp.UserCode,
		FullName:         resp.FullName,
		Role:             resp.Role.Name,
		Department:       deptID,
		JobTitle:         resp.JobTitle,
		Gender:           resp.Gender,
		PhoneNumber:      resp.PhoneNumber,
		Email:            resp.Email,
		UserName:         resp.UserName,
		Realm:            resp.Realm,
		Enabled:          resp.Enabled,
		DateJoined:       &resp.DateJoined,
		LastModified:     &resp.LastModified,
		Country:          resp.Country,
		Region:           resp.Region,
		LastLogin:        resp.LastLogin,
		PasswordDisable:  resp.PasswordDisable,
		IsFirstTimeLogin: resp.IsFirstTimeLogin,
		CreatedAt:        resp.CreatedAt,
	}
}
