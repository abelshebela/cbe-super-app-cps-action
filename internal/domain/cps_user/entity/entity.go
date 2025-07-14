package cpsuser

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUser struct {
	ID                 string          `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode           string          `json:"user_code,omitempty" bson:"user_code,omitempty"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Role               string          `json:"role,omitempty" bson:"role,omitempty"`
	Department         string          `json:"department,omitempty" bson:"department,omitempty"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty"`
	UserName           string          `json:"username,omitempty" bson:"username,omitempty"`
	Realm              string          `json:"realm,omitempty" bson:"realm,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category,omitempty"`
	PermissionGroup    []bson.ObjectID `json:"permission_group,omitempty" bson:"permission_group,omitempty"`

	Password                 Password  `json:"password" bson:"password"`
	PasswordDisable          bool      `json:"password_disable,omitempty" bson:"password_disable,omitempty"`
	SyncDisabled             bool      `json:"sync_disabled,omitempty" bson:"sync_disabled,omitempty"`
	LoginAttemptCount        uint8     `json:"login_attempt_count,omitempty" bson:"login_attempt_count,omitempty"`
	LastLoginAttempt         time.Time `json:"last_login_attempt,omitempty" bson:"last_login_attempt,omitempty"`
	NextLoginAttempt         time.Time `json:"next_login_attempt,omitempty" bson:"next_login_attempt,omitempty"`
	LastOnlineDate           time.Time `json:"last_online_date,omitempty" bson:"last_online_date,omitempty"`
	LastLogin                time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
	LoginPassword            string    `json:"login_password,omitempty" bson:"login_password,omitempty"`
	AccountAuthorizationCode string    `json:"account_authorization_code,omitempty" bson:"account_authorization_code,omitempty"`
	UnlockAccountRequested   bool      `json:"unlock_account_requested,omitempty" bson:"unlock_account_requested,omitempty"`

	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" bson:"password_changed_at,omitempty"`
	OTPStatus         string     `json:"otp_status,omitempty" bson:"otp_status,omitempty"`
	OTPLastTriedAt    *time.Time `json:"otp_last_tried_at,omitempty" bson:"otp_last_tried_at,omitempty"`
	OPTLastVerifiedAt *time.Time `json:"otp_last_verified_at,omitempty" bson:"otp_last_verified_at,omitempty"`
	OTPVerifyCount    int        `json:"otp_verify_count,omitempty" bson:"otp_verify_count,omitempty"`

	Enabled      bool       `json:"enabled,omitempty" bson:"enabled,omitempty"`
	IsDeleted    bool       `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	DateJoined   *time.Time `json:"date_joined,omitempty" bson:"date_joined,omitempty"`
	LastModified *time.Time `json:"last_modified,omitempty" bson:"last_modified,omitempty"`

	Country string `json:"country,omitempty" bson:"country,omitempty"`
	Region  string `json:"region,omitempty" bson:"region,omitempty"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}

func MapToCPSUser(userData userDTO.CreateUserRequest) (model.CPSUser, error) {
	return model.CPSUser{
		FullName:           userData.FullName,
		Role:               userData.Role,
		Department:         userData.Department,
		Gender:             userData.Gender,
		PhoneNumber:        userData.PhoneNumber,
		Email:              userData.Email,
		UserName:           userData.UserName,
		Realm:              userData.Realm,
		PermissionCategory: userData.PermissionCategory,
		PermissionGroup:    userData.PermissionGroups,

		Password: model.Password{
			Salt:             "",
			CurrentPassword:  "",
			OldPassword:      [4]string{},
			PasswordChangeAt: time.Time{},
		},

		PasswordDisable:          false,
		SyncDisabled:             false,
		LoginAttemptCount:        0,
		LastLoginAttempt:         time.Time{},
		NextLoginAttempt:         time.Time{},
		LastOnlineDate:           time.Time{},
		LastLogin:                time.Time{},
		LoginPassword:            "",
		AccountAuthorizationCode: "",
		UnlockAccountRequested:   false,
		PasswordChangedAt:        nil,
		OTPStatus:                "",
		OTPLastTriedAt:           nil,
		OPTLastVerifiedAt:        nil,
		OTPVerifyCount:           0,
		Enabled:                  true,
		IsDeleted:                false,
		DateJoined:               nil,
		LastModified:             nil,
		Country:                  "",
		Region:                   "",
	}, nil
}
