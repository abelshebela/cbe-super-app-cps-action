package entities

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/enums"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/type_definition"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSUser struct {
	ID                bson.ObjectID            `json:"id" bson:"_id,omitempty"`
	UserCode          string                   `json:"user_code" bons:"user_code"` // generated
	FullName          string                   `json:"full_name" bson:"full_name"`
	UserName          string                   `json:"UserName" bson:"UserName"`
	PhoneNumber       string                   `json:"phone_number" bson:"phone_number"`
	BranchCode        []string                 `json:"branch_code" bson:"branch_code"` // IFB, CB, HomeBranch(CB == Homebranch) ice versa
	BranchName        string                   `json:"branch_name" bson:"branch_name"`
	HomeBranch        string                   `json:"home_branch" bson:"home_branch"`
	Role              string                   `json:"role" bson:"role"`   // enum: maker, checker, aduditer
	Realm             enums.Realm              `json:"realm" bson:"realm"` // bank
	LoginAttemptCount uint8                    `json:"login_attempt_count" bson:"login_attempt_count"`
	Password          type_definition.Password `json:"password" bson:"password"`
	FirstPasswordSet  bool                     `json:"first_password_set" bson:"first_password_set"`
	Enabled           bool                     `json:"enabled" bson:"enabled"`
	IsDeleted         bool                     `json:"is_deleted" bson:"is_deleted"`
	OTPVerifyCount    uint8                    `json:"otp_verfy_count" bson:"otp_verify_count"`
	OTPLastTriedAt    time.Time                `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt time.Time                `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	PermissionGroup   []bson.ObjectID          `json:"permission_group" bson:"permission_group"`
	Permissions       []bson.ObjectID          `json:"permissions" bson:"permissions"`
	LastLoginAttempt  time.Time                `json:"last_login_attempt" bson:"last_login_attempt"`
	NextLoginAttempt  time.Time                `json:"next_login_attempt" bson:"next_login_attempt"`
	IsFirstTimeLogin  bool                     `json:"is_first_time_login" bson:"is_first_time_login"`
	LastLogin         time.Time                `json:"last_login" bson:"last_login"`
	CreatedAt         time.Time                `json:"created_at" bson:"created_at,omitempty"`
	LastModifiedAt    time.Time                `json:"last_modifed_at" bson:"last_modifed_at,omitempty"`
}

// 888888888888888888888888888888888888888888888888888888888888888888888888888888888
