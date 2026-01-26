package bps_user_core

import (
	"errors"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	// shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	// "go.mongodb.org/mongo-driver/v2/bson"
	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	// local_model "cbe-super-app-cps-action/internal/constants/model"
)

func ExistingIdentifier(existing *bps_model.BPSUser, req bps_model.BPSUser) error {
	if existing == nil {
		return nil
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	normalizedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	normalizedUsername := strings.TrimSpace(req.Username)

	// Email duplicate
	if normalizedEmail != "" && strings.EqualFold(strings.TrimSpace(existing.Email), normalizedEmail) {
		return errors.New(localization.ErrorEmailAlreadyExist.Code)
	}
	// Phone duplicate
	if normalizedPhone != "" {
		// Normalize stored phone too, just in case
		storedPhone := local_util.FormatPhoneNumber(existing.PhoneNumber)
		if storedPhone == normalizedPhone {
			return errors.New(localization.ErrorPhonenumberAlreadyExist.Code)
		}
	}
	// Account number duplicate
	if strings.EqualFold(strings.TrimSpace(existing.Username), normalizedUsername) {
		return errors.New(localization.ErrorUsernameAlreadyExist.Code)
	}
	return nil
}

func ExistingIdentifierForUpdate(existing bps_model.BPSUser, id string, req bps_model.BPSUser) error {
	if &existing == nil {
		return nil
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	normalizedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	normalizedUsername := strings.TrimSpace(req.Username)

	existingID := local_util.FirstHex24(existing.ID.String())
	// Email duplicate
	if normalizedEmail != "" && strings.EqualFold(existing.Email, normalizedEmail) && existingID != id {
		return errors.New(localization.ErrorEmailAlreadyExist.Code)
	}
	// Phone duplicate
	if normalizedPhone != "" && existingID != id {
		// Normalize stored phone too, just in case
		storedPhone := local_util.FormatPhoneNumber(existing.PhoneNumber)
		if storedPhone == normalizedPhone {
			return errors.New(localization.ErrorPhonenumberAlreadyExist.Code)
		}
	}
	// Account number duplicate
	if normalizedUsername != "" && strings.EqualFold(strings.TrimSpace(existing.Username), normalizedUsername) && existingID != id {
		return errors.New(localization.ErrorUsernameAlreadyExist.Code)
	}
	return nil
}
func BPSUser_mapper(action map[string]interface{}) bps_model.BPSUser {
	var user bps_model.BPSUser

	if v, ok := action["user_code"]; ok {
		if s, ok := v.(string); ok {
			user.UserCode = s
		}
	}
	if v, ok := action["full_name"]; ok {
		if s, ok := v.(string); ok {
			user.FullName = s
		}
	}
	if v, ok := action["username"]; ok {
		if s, ok := v.(string); ok {
			user.Username = s
		}
	}

	if v, ok := action["job_title"]; ok {
		if s, ok := v.(string); ok {
			user.JobTitle = s
		}
	}

	if v, ok := action["phone_number"]; ok {
		if s, ok := v.(string); ok {
			user.PhoneNumber = s
		}
	}
	if v, ok := action["email"]; ok {
		if s, ok := v.(string); ok {
			user.Email = s
		}
	}
	if v, ok := action["branch_code"]; ok {
		switch val := v.(type) {
		case []string:
			user.BranchCode = val
		case []interface{}:
			var branchCodes []string
			for _, v := range val {
				if str, ok := v.(string); ok {
					branchCodes = append(branchCodes, str)
				}
			}
			user.BranchCode = branchCodes
		case string:
			user.BranchCode = []string{val}
		}
	}
	if v, ok := action["branch_name"]; ok {
		if s, ok := v.(string); ok {
			user.BranchName = s
		}
	}
	if v, ok := action["home_branch"]; ok {
		if s, ok := v.(string); ok {
			user.HomeBranch = s
		}
	}
	if v, ok := action["role"]; ok {
		if s, ok := v.(string); ok {
			user.Role = s
		}
	}
	if v, ok := action["login_attempt_count"]; ok {
		if n, ok := v.(uint8); ok {
			user.LoginAttemptCount = n
		}
	}
	if v, ok := action["password"]; ok {
		if p, ok := v.(bps_model.Password); ok {
			user.Password = p
		}
	}
	if v, ok := action["first_password_set"]; ok {
		if b, ok := v.(bool); ok {
			user.FirstPasswordSet = b
		}
	}
	if v, ok := action["enabled"]; ok {
		if b, ok := v.(bool); ok {
			user.Enabled = b
		}
	}
	if v, ok := action["is_deleted"]; ok {
		if b, ok := v.(bool); ok {
			user.IsDeleted = b
		}
	}
	if v, ok := action["otp_verify_count"]; ok {
		if n, ok := v.(uint8); ok {
			user.OTPVerifyCount = n
		}
	}
	if v, ok := action["otp_last_tried_at"]; ok {
		if t, ok := v.(time.Time); ok {
			user.OTPLastTriedAt = t
		}
	}
	if v, ok := action["otp_last_verified_at"]; ok {
		if t, ok := v.(time.Time); ok {
			user.OTPLastVerifiedAt = t
		}
	}
	// if v, ok := action["permission_group"]; ok {
	// 	if arr, ok := v.([]bson.ObjectID); ok {
	// 		user.PermissionGroup = arr
	// 	}
	// }
	// if v, ok := action["permissions"]; ok {
	// 	if arr, ok := v.([]bson.ObjectID); ok {
	// 		user.Permissions = arr
	// 	}
	// }
	if v, ok := action["last_login_attempt"]; ok {
		if t, ok := v.(time.Time); ok {
			user.LastLoginAttempt = t
		}
	}
	if v, ok := action["next_login_attempt"]; ok {
		if t, ok := v.(time.Time); ok {
			user.NextLoginAttempt = t
		}
	}
	// if v, ok := action["is_first_time_login"]; ok {
	// 	if b, ok := v.(bool); ok {
	// 		user.IsFirstTimeLogin = b
	// 	}
	// }
	if v, ok := action["last_login"]; ok {
		if t, ok := v.(time.Time); ok {
			user.LastLogin = t
		}
	}

	return user
}

// Helper to map Password struct field-by-field
func mapPassword(src types.Password) local_model.Password {
	return local_model.Password{
		Salt:             src.Salt,
		CurrentPassword:  src.CurrentPassword,
		OldPassword:      src.OldPassword,
		PasswordChangeAt: src.PasswordChangeAt,
	}
}
func mapPasswordToShared(src local_model.Password) types.Password {
	return types.Password{
		Salt:             src.Salt,
		CurrentPassword:  src.CurrentPassword,
		OldPassword:      src.OldPassword,
		PasswordChangeAt: src.PasswordChangeAt,
	}
}

// MapBPSUserToWithJobTitle maps model.BPSUser to BPSUserWithJobTitle
func MapBPSUserToWithJobTitle(u model.BPSUser, jobTitle string) local_model.BPSUser {
	return local_model.BPSUser{
		ID:                u.ID,
		UserCode:          u.UserCode,
		FullName:          u.FullName,
		UserName:          u.UserName,
		PhoneNumber:       u.PhoneNumber,
		BranchCode:        u.BranchCode,
		BranchName:        u.BranchName,
		HomeBranch:        u.HomeBranch,
		JobTitle:          jobTitle,
		Role:              u.Role,
		LoginAttemptCount: u.LoginAttemptCount,
		Password:          mapPassword(u.Password),
		FirstPasswordSet:  u.FirstPasswordSet,
		Enabled:           u.Enabled,
		IsDeleted:         u.IsDeleted,
		OTPVerifyCount:    u.OTPVerifyCount,
		OTPLastTriedAt:    u.OTPLastTriedAt,
		OTPLastVerifiedAt: u.OTPLastVerifiedAt,
		PermissionGroup:   u.PermissionGroup,
		Permissions:       u.Permissions,
		LastLoginAttempt:  u.LastLoginAttempt,
		NextLoginAttempt:  u.NextLoginAttempt,
		IsFirstTimeLogin:  u.IsFirstTimeLogin,
		LastLogin:         u.LastLogin,
		CreatedAt:         u.CreatedAt,
		LastModifiedAt:    u.LastModifiedAt,
	}
}

// MapWithJobTitleToBPSUser maps BPSUserWithJobTitle to model.BPSUser (drops JobTitle)
func MapWithJobTitleToBPSUser(u local_model.BPSUser) model.BPSUser {
	return model.BPSUser{
		ID:                u.ID,
		UserCode:          u.UserCode,
		FullName:          u.FullName,
		UserName:          u.UserName,
		PhoneNumber:       u.PhoneNumber,
		BranchCode:        u.BranchCode,
		BranchName:        u.BranchName,
		HomeBranch:        u.HomeBranch,
		Role:              u.Role,
		LoginAttemptCount: u.LoginAttemptCount,
		Password:          mapPasswordToShared(u.Password),
		FirstPasswordSet:  u.FirstPasswordSet,
		Enabled:           u.Enabled,
		IsDeleted:         u.IsDeleted,
		OTPVerifyCount:    u.OTPVerifyCount,
		OTPLastTriedAt:    u.OTPLastTriedAt,
		OTPLastVerifiedAt: u.OTPLastVerifiedAt,
		PermissionGroup:   u.PermissionGroup,
		Permissions:       u.Permissions,
		LastLoginAttempt:  u.LastLoginAttempt,
		NextLoginAttempt:  u.NextLoginAttempt,
		IsFirstTimeLogin:  u.IsFirstTimeLogin,
		LastLogin:         u.LastLogin,
		CreatedAt:         u.CreatedAt,
		LastModifiedAt:    u.LastModifiedAt,
	}
}
