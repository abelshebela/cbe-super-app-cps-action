package bps_user_core

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func BPSUser_mapper(action map[string]interface{}) model.BPSUser {
	var user model.BPSUser

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
			user.UserName = s
		}
	}
	if v, ok := action["phone_number"]; ok {
		if s, ok := v.(string); ok {
			user.PhoneNumber = s
		}
	}
	if v, ok := action["branch_code"]; ok {
		if arr, ok := v.([]string); ok {
			user.BranchCode = arr
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
		if p, ok := v.(shared_types.Password); ok {
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
	if v, ok := action["permission_group"]; ok {
		if arr, ok := v.([]bson.ObjectID); ok {
			user.PermissionGroup = arr
		}
	}
	if v, ok := action["permissions"]; ok {
		if arr, ok := v.([]bson.ObjectID); ok {
			user.Permissions = arr
		}
	}
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
	if v, ok := action["is_first_time_login"]; ok {
		if b, ok := v.(bool); ok {
			user.IsFirstTimeLogin = b
		}
	}
	if v, ok := action["last_login"]; ok {
		if t, ok := v.(time.Time); ok {
			user.LastLogin = t
		}
	}

	return user
}
