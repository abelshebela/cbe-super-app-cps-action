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
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category"`
	PermissionGroup    []bson.ObjectID `json:"permission_group,omitempty" bson:"permission_group"`
	JobTitle           string          `json:"job_title" bson:"job_title"`

	Password                 types.Password `json:"password" bson:"password"`
	PasswordDisable          bool           `json:"password_disable,omitempty" bson:"password_disable"`
	SyncDisabled             bool           `json:"sync_disabled,omitempty" bson:"sync_disabled"`
	LoginAttemptCount        uint8          `json:"login_attempt_count,omitempty" bson:"login_attempt_count"`
	LastLoginAttempt         time.Time      `json:"last_login_attempt,omitempty" bson:"last_login_attempt"`
	NextLoginAttempt         time.Time      `json:"next_login_attempt,omitempty" bson:"next_login_attempt"`
	LastOnlineDate           time.Time      `json:"last_online_date,omitempty" bson:"last_online_date"`
	LastLogin                time.Time      `json:"last_login" bson:"last_login"`
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
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`
	CreatedBy    string     `json:"created_by,omitempty" bson:"created_by,omitempty"`
	ChangedBy    string     `json:"changed_by,omitempty" bson:"changed_by,omitempty"`

	Country string `json:"country,omitempty" bson:"country"`
	Region  string `json:"region,omitempty" bson:"region"`

	IsDelegationActive bool          `json:"is_delegation_active" bson:"is_delegation_active"`
	DelegationID       bson.ObjectID `json:"delegation_id" bson:"delegation_id"`
	DelegatedRoleCode  string        `json:"delegated_role_code" bson:"delegated_role_code"`
}

type ExportCPSUser struct {
	// find from cps_user collection
	FirstName   string    `json:"first_name,omitempty" bson:"first_name"`
	MiddleName  string    `json:"middle_name,omitempty" bson:"middle_name"`
	LastName    string    `json:"last_name,omitempty" bson:"last_name"`
	PhoneNumber string    `json:"phone_number,omitempty" bson:"phone_number"`
	Email       string    `json:"email,omitempty" bson:"email"`
	Department  string    `json:"department,omitempty" bson:"department"`
	JobTitle    string    `json:"job_title" bson:"job_title"`
	Role        string    `json:"role,omitempty" bson:"role"`
	UserName    string    `json:"username,omitempty" bson:"username"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	Enabled     bool      `json:"enabled,omitempty" bson:"enabled"`
	LastLogin   time.Time `json:"last_login" bson:"last_login"`

	// from roles_delegation collection by using the returned user id from cps_user collection filter by delegated_user_id == current users id and is_active == true plus the statat and endat for the roles_delegation  then we will get the endat for ExpiryDateForDelegation
	// if roles_delegation is found then the UserType will be "Delegation" if not found then the UserType will be "permanent" while keeping the ExpiryDateForDelegation empty
	ExpiryDateForDelegation time.Time `json:"expiry_date_for_delegation" bson:"expiry_date_for_delegation"`
	UserType                string    `json:"user_type,omitempty" bson:"user_type"`

	// from cps_actions collection filter by request_action between "CREATE_CPS_USER" "UPDATE_CPS_USER" "DELETE_CPS_USER" "ENABLE_CPS_USER" "DISABLE_CPS_USER" thenwe will get the last action and last modified date for that user
	// cps_actions.previous_action.id == current users id => then we will get the last modification action(which will be cps_actions.request_action) and last modified date(which will be cps_actions.created_at)
	LastModificationAction string    `json:"last_modification_action" bson:"last_modification_action"`
	LastModified           time.Time `json:"last_modified,omitempty" bson:"last_modified"`

	// from cps_actions collection filter by request_action ==  "CREATE_CPS_USER" and action_status == "Approved" then we find for cps_actions.current_action.user_code == current users user_code => then we will get the created by user which will be cps_actions.maker_id and approved by user which will be cps_actions.checker_users[0].checker_id
	// if maker_id is empty replay with "SSO" and approved by with ""
	CreatedBy  string `json:"created_by,omitempty" bson:"created_by"`
	ApprovedBy string `json:"approved_by,omitempty" bson:"approved_by"`
}
