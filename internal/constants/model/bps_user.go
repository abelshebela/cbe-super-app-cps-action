package model

import (
	// "cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSUser struct {
	ID                bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string          `json:"user_code,omitempty" bson:"user_code,omitempty"`
	FullName          string          `json:"full_name,omitempty" bson:"full_name,omitempty"`
	UserName          string          `json:"username,omitempty" bson:"username,omitempty"`
	PhoneNumber       string          `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	BranchCode        []string        `json:"branch_code,omitempty" bson:"branch_code,omitempty"`
	Email             string          `json:"email,omitempty" bson:"email,omitempty"`
	BranchName        string          `json:"branch_name,omitempty" bson:"branch_name,omitempty"`
	HomeBranch        string          `json:"home_branch,omitempty" bson:"home_branch,omitempty"`
	JobTitle          string          `json:"job_title,omitempty" bson:"job_title,omitempty"`
	Role              string          `json:"role,omitempty" bson:"role,omitempty"`
	RoleName          string          `json:"role_name,omitempty" bson:"role_name,omitempty"`
	LoginAttemptCount uint8           `json:"-" bson:"login_attempt_count"`
	Password          Password        `json:"-" bson:"password"`
	FirstPasswordSet  bool            `json:"-" bson:"first_password_set"`
	Enabled           bool            `json:"enabled" bson:"enabled"`
	IsDeleted         bool            `json:"-" bson:"is_deleted"`
	OTPVerifyCount    uint8           `json:"-" bson:"otp_verfy_count"`
	OTPLastTriedAt    time.Time       `json:"-" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt time.Time       `json:"-" bson:"otp_last_verified_at"`
	PermissionGroup   []bson.ObjectID `json:"-" bson:"permission_group"`
	Permissions       []bson.ObjectID `json:"-" bson:"permissions"`
	LastLoginAttempt  time.Time       `json:"-" bson:"last_login_attempt"`
	NextLoginAttempt  time.Time       `json:"-" bson:"next_login_attempt"`
	IsFirstTimeLogin  bool            `json:"-" bson:"is_first_time_login"`
	LastLogin         time.Time       `json:"-" bson:"last_login"`
	CreatedAt         time.Time       `json:"created_at" bson:"created_at,omitempty"`
	LastModifiedAt    time.Time       `json:"-" bson:"last_modifed_at,omitempty"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}
