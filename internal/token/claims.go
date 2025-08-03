package token

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
)

func ClaimBuilder(user *model.User, sessionExpiry interface{}, action interface{}, additional map[string]interface{}) map[string]interface{} {
	claims := map[string]interface{}{
		"user_id":        user.ID,
		"user_code":      user.UserCode,
		"full_name":      user.FullName,
		"phone_number":   user.PhoneNumber,
		"user_email":     user.Email,
		"user_realm":     user.Realm,
		"ifb_member":     user.MemberType == constants.IFBT,
		"device_uuid":    user.DeviceUUID,
		"session_expiry": sessionExpiry,
		"action":         action,
	}

	if additional != nil {
		for k, v := range additional {
			claims[k] = v
		}
	}

	return claims
}

func PermanentClaimBuilder(user *model.User, sessionExpiry interface{}, additional map[string]interface{}) map[string]interface{} {
	claims := map[string]interface{}{
		"user_id":                 user.ID,
		"user_code":               user.UserCode,
		"full_name":               user.FullName,
		"phone_number":            user.PhoneNumber,
		"user_email":              user.Email,
		"user_realm":              user.Realm,
		"ifb_member":              user.MemberType == "ifb",
		"device_uuid":             user.DeviceUUID,
		"kyc_level":               user.KYCLevel,
		"user_device_linked_date": "",
		"session_expiry":          sessionExpiry,
	}

	if additional != nil {
		for k, v := range additional {
			claims[k] = v
		}
	}

	return claims
}
