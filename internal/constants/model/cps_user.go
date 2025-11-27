package model

import (
	"time"

	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUser struct {
	ID                 bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode           string          `json:"user_code,omitempty" bson:"user_code"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name"`
	Role               string          `json:"role,omitempty" bson:"role"`
	Department         bson.ObjectID   `json:"department,omitempty" bson:"department"`
	Gender             string          `json:"gender,omitempty" bson:"gender"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number"`
	Email              string          `json:"email,omitempty" bson:"email"`
	UserName           string          `json:"username,omitempty" bson:"username"`
	Realm              string          `json:"realm,omitempty" bson:"realm"`
	PermissionCategory []string `json:"permission_category,omitempty" bson:"permission_category"`
	PermissionGroup    []string `json:"permission_group,omitempty" bson:"permission_group"`

	Password                 types.Password `json:"password" bson:"password"`
	PasswordDisable          bool           `json:"password_disable,omitempty" bson:"password_disable"`
	SyncDisabled             bool           `json:"sync_disabled,omitempty" bson:"sync_disabled"`
	LoginAttemptCount        uint8          `json:"login_attempt_count,omitempty" bson:"login_attempt_count"`
	LastLoginAttempt         time.Time      `json:"last_login_attempt,omitempty" bson:"last_login_attempt"`
	NextLoginAttempt         time.Time      `json:"next_login_attempt,omitempty" bson:"next_login_attempt"`
	LastOnlineDate           time.Time      `json:"last_online_date,omitempty" bson:"last_online_date"`
	LastLogin                time.Time      `json:"last_login,omitempty" bson:"last_login"`
	LoginPassword            string         `json:"login_password,omitempty" bson:"login_password"`
	AccountAuthorizationCode string         `json:"account_authorization_code,omitempty" bson:"account_authorization_code"`
	UnlockAccountRequested   bool           `json:"unlock_account_requested,omitempty" bson:"unlock_account_requested"`

	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" bson:"password_changed_at"`
	OTPStatus         string     `json:"otp_status,omitempty" bson:"otp_status"`
	OTPLastTriedAt    *time.Time `json:"otp_last_tried_at,omitempty" bson:"otp_last_tried_at"`
	OPTLastVerifiedAt *time.Time `json:"otp_last_verified_at,omitempty" bson:"otp_last_verified_at"`
	OTPVerifyCount    int        `json:"otp_verify_count,omitempty" bson:"otp_verify_count"`
	IsFirstTimeLogin  bool       `json:"is_first_time_login" bson:"is_first_time_login"`

	Enabled      bool       `json:"enabled,omitempty" bson:"enabled"`
	IsDeleted    bool       `json:"is_deleted,omitempty" bson:"is_deleted"`
	DateJoined   *time.Time `json:"date_joined,omitempty" bson:"date_joined"`
	LastModified *time.Time `json:"last_modified,omitempty" bson:"last_modified"`

	Country string `json:"country,omitempty" bson:"country"`
	Region  string `json:"region,omitempty" bson:"region"`
}
