package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSUser struct {
	ID                bson.ObjectID   `json:"id" bson:"_id,omitempty"`
	UserCode          string          `json:"user_code" bson:"user_code"`
	FullName          string          `json:"full_name" bson:"full_name"`
	UserName          string          `json:"username" bson:"username"`
	PhoneNumber       string          `json:"phone_number" bson:"phone_number"`
	BranchCode        []string        `json:"branch_code" bson:"branch_code"`
	BranchName        string          `json:"branch_name" bson:"branch_name"`
	HomeBranch        string          `json:"home_branch" bson:"home_branch"`
	Role              string          `json:"role" bson:"role"`
	LoginAttemptCount uint8           `json:"-" bson:"login_attempt_count"`
	Password          types.Password  `json:"-" bson:"password"`
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
