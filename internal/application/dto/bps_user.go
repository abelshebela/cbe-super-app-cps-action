package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSUser struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserCode    string        `json:"user_code" bson:"user_code"`
	FullName    string        `json:"full_name" bson:"full_name"`
	Username    string        `json:"username" bson:"username"`
	PhoneNumber string        `json:"phone_number" bson:"phone_number"`
	BranchCode  []string      `json:"branch_code" bson:"branch_code"`
	BranchName  string        `json:"branch_name" bson:"branch_name"`
	HomeBranch  string        `json:"home_branch" bson:"home_branch"`
	Role        string        `json:"role" bson:"role"`
	Realm       string        `json:"realm" bson:"realm"`
	// LoginAttemptCount uint8           `json:"login_attempt_count" bson:"login_attempt_count"`
	// Password          Password        `json:"password" bson:"password"`
	// FirstPasswordSet  bool            `json:"first_password_set" bson:"first_password_set"`
	Enabled   bool `json:"enabled" bson:"enabled"`
	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`
	// OTPVerifyCount    uint8           `json:"otp_verify_count" bson:"otp_verify_count"`
	// OTPLastTriedAt    time.Time       `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	// OTPLastVerifiedAt time.Time       `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	// PermissionGroup   []bson.ObjectID `json:"permission_group" bson:"permission_group"`
	// Permissions       []bson.ObjectID `json:"permissions" bson:"permissions"`
	// LastLoginAttempt  time.Time       `json:"last_login_attempt" bson:"last_login_attempt"`
	// NextLoginAttempt  time.Time       `json:"next_login_attempt" bson:"next_login_attempt"`
	// IsFirstTimeLogin  bool            `json:"is_first_time_login" bson:"is_first_time_login"`
	// LastLogin         time.Time       `json:"last_login" bson:"last_login"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at,omitempty"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at,omitempty"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password"`
	PasswordChangeAt time.Time `json:"password_change_at" bson:"password_change_at"`
}
