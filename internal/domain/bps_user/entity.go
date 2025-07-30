package bps_user

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSUser struct {
	ID                bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserCode          string        `json:"user_code" bson:"user_code"`
	FullName          string        `json:"full_name" bson:"full_name"`
	Username          string        `json:"username" bson:"username"`
	PhoneNumber       string        `json:"phone_number" bson:"phone_number"`
	BranchCode        []string      `json:"branch_code" bson:"branch_code"`
	BranchName        string        `json:"branch_name" bson:"branch_name"`
	HomeBranch        string        `json:"home_branch" bson:"home_branch"`
	Role              string        `json:"role" bson:"role"`
	Realm             string        `json:"realm" bson:"realm"`
	Password          Password      `json:"password" bson:"password"`
	Enabled           bool          `json:"enabled" bson:"enabled"`
	IsDeleted         bool          `json:"is_deleted" bson:"is_deleted"`
	// LastLoginAttempt  time.Time     `json:"last_login_attempt" bson:"last_login_attempt"`
	// NextLoginAttempt  time.Time     `json:"next_login_attempt" bson:"next_login_attempt"`
	// LastLogin         time.Time     `json:"last_login" bson:"last_login"`
	// CreatedAt         time.Time     `json:"created_at" bson:"created_at"`
	// LastModifiedAt    time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password"`
	PasswordChangeAt time.Time `json:"password_change_at" bson:"password_change_at"`
}

type UserType string

const (
	Maker   UserType = "MAKER"
	Checker UserType = "CHECKER"
)

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type ActionType string

const (
	ActionCreate  ActionType = "CREATE"
	ActionUpdate  ActionType = "UPDATE"
	ActionDelete  ActionType = "DELETE"
	ActionEnable  ActionType = "ENABLE"
	ActionDisable ActionType = "DISABLE"
)

type RequestAction string

const (
	RequestBPSUser        RequestAction = "BPS_USER"
	RequestDisableBPSUser RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser  RequestAction = "ENABLE_BPS_USER"
)

type User struct {
	UserID      string    `json:"user_id" bson:"user_id"`
	FullName    string    `json:"full_name" bson:"full_name"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
}

type MakerAndChecker struct {
	Maker   User `json:"maker" bson:"maker"`
	Checker User `json:"checker" bson:"checker"`
}

