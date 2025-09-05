package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ArchivedUser struct {
	ID                bson.ObjectID           `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string                  `json:"user_code,omitempty"" bson:"user_code"`
	FullName          string                  `json:"full_name,omitempty"" bson:"full_name"`
	MotherName        string                  `json:"mother_name,omitempty"" bson:"mother_name"`
	Nationality       string                  `json:"nationality,omitempty"" bson:"nationality"`
	BirthDate         time.Time               `json:"birth_date,omitempty"" bson:"birth_date"`
	ResidentialStatus string                  `json:"residential_status,omitempty"" bson:"residential_status"`
	IssuedDate        time.Time               `json:"issued_date,omitempty"" bson:"isssued_date"`
	PhoneNumber       string                  `json:"phone_number,omitempty"" bson:"phone_number"`
	Gender            constants.Gender        `json:"gender,omitempty"" bson:"gender"`
	MaritalStatus     constants.MaritalStatus `json:"marital_status,omitempty"" bson:"marital_status"`
	Fayda             struct {
		FaydaID          string `json:"id_number,omitempty"" bson:"id_number"`
		FaydaAccessToken string `json:"fayda_access_token,omitempty"" bson:"fayda_access_token"`
		EmploymentStatus string `json:"employment_status,omitempty"" bson:"employement_status"`
		EmployerName     string `json:"employer_name,omitempty"" bson:"employer_name"`
		IssuedBy         string `json:"issued_by,omitempty"" bson:"issued_by"`
		MonthlyIncome    uint64 `json:"monthly_incode,omitempty"" bson:"monthly_incode"`
	} `json:"fayda,omitempty"" bson:"fayda"`
	Address          types.Address              `json:"address,omitempty"" bson:"address"`
	DocumentFront    string                     `json:"document_front,omitempty"" bson:"document_front"`
	DocumentBack     string                     `json:"document_back,omitempty"" bson:"document_back"`
	Photo            string                     `json:"photo,omitempty"" bson:"photo"`
	Signature        string                     `json:"signature,omitempty"" bson:"signature"`
	Avatar           string                     `json:"avater,omitempty"" bson:"avater"`
	Email            string                     `json:"email,omitempty"" bson:"email"`
	PushToken        string                     `json:"push_token,omitempty"" bson:"push_token"`
	Realm            constants.Realm            `json:"realm,omitempty"" bson:"realm"`
	PermissionGroup  []bson.ObjectID            `json:"permission_group,omitempty"" bson:"permission_group"`
	Permissions      []bson.ObjectID            `json:"permissions,omitempty"" bson:"persmissions"`
	IsAccountBlocked bool                       `json:"is_account_blocked,omitempty"" bson:"is_account_blocked"`
	IsAccountLinked  bool                       `json:"is_account_linked,omitempty"" bson:"is_account_linked"`
	MemberType       constants.MemberType       `json:"account_branch_type,omitempty"" bson:"account_branch_type"`
	RegistrationType constants.RegistrationType `json:"account_type,omitempty"" bson:"account_type"`
	AccountStatus    constants.AccountStatus    `json:"account_status,omitempty"" bson:"account_status"`
	KYCLevel         uint8                      `json:"kyc_level,omitempty"" bson:"kyc_level"`
	KYC              struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed,omitempty"" bson:"kyc_reject_reason_failed"`
		KYCStatus            constants.KYCStatus `json:"kyc_status,omitempty"" bson:"kyc_status"`
		KYCRejectReason      string              `json:"kyc_reject_reason,omitempty"" bson:"kyc_reject_reason"`
		KYCIsApproved        bool                `json:"kyc_approved,omitempty"" bson:"kyc_approved"`
		KYCActivityBy        map[string]struct{} `json:"kyc_activity_by,omitempty"" bson:"kyc_activity_by"`
	} `json:"kyc,omitempty"" bson:"kyc"`
	IsBranchApproved   bool      `json:"is_branch_approved,omitempty"" bson:"is_branch_approved"`
	IsVerified         bool      `json:"is_verfied,omitempty"" bson:"is_verfied"`
	IsBlocked          bool      `json:"is_blocked,omitempty"" bson:"is_blocked"`                     // default: false
	IsAccessRestricted bool      `json:"is_access_restricted,omitempty"" bson:"is_access_restricted"` // default: false
	BlockedAt          time.Time `json:"blocked_at,omitempty"" bson:"blocked_at"`
	RegisterBy         struct{}  `json:"register_by,omitempty"" bson:"register_By"`
	LoginAttemptCount  uint8     `json:"login_attempt_count,omitempty"" bson:"login_attempt_count"`
	LastLoginAttempt   time.Time `json:"last_login_attempt,omitempty"" bson:"last_login_attempt"`
	LastOnlineDate     time.Time `json:"last_online_date,omitempty"" bson:"last_online_date"`
	LastLogin          time.Time `json:"last_login,omitempty"" bson:"last_login"`

	// move BPS Status and realted field to new collection
	BPSStatus constants.BPSStatus `json:"bps_reject_status,omitempty"" bson:"bps_reject_status"`
	// BPSRejectionReason string    `json:"bps_reject_reason,omitempty"" bson:"bps_reject_reason"`
	// BPSRejectionField  []string  `json:"bps_reject_failed,omitempty"" bson:"bps_reject_failed"`

	LoginPIN   types.LoginPIN `json:"pin,omitempty"" bson:"login_pin"`
	DeviceUUID string         `json:"device_uuid,omitempty"" bson:"device_uuid"`
	Device     struct {
		DevicePlatform      string    `json:"device_platform,omitempty"" bson:"device_platform"`
		AppVersion          string    `json:"app_version,omitempty"" bson:"app_version"`
		APPInstallationDate time.Time `json:"application_installation_date,omitempty"" bson:"application_installation_date"`
	} `json:"device,omitempty"" bson:"device"`
	CustomerNumber        string                          `json:"customer_number,omitempty"" bson:"customer_number"`
	InitialLinkedDate     time.Time                       `json:"initial_linked_date,omitempty"" bson:"initiali_linked_date"`
	PrimaryAuthentication constants.PrimaryAuthentication `json:"primary_authentication,omitempty"" bson:"primary_authentication"`
	// LoanScore             uint8                 `json:"loan_score,omitempty"" bson:"loan_score"`
	DeviceStatus      constants.DeviceStatus `json:"device_status,omitempty"" bson:"device_status"`
	Enabled           bool                   `json:"enabled,omitempty"" bson:"enabled"`
	IsDeleted         bool                   `json:"is_deleted,omitempty"" bson:"is_deleted"`
	PINChangedAt      time.Time              `json:"pin_changed_at,omitempty"" bson:"pin_changed_at"`
	OTPLastTriedAt    time.Time              `json:"otp_last_tried_at,omitempty"" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt time.Time              `json:"otp_last_verified_at,omitempty"" bson:"otp_last_verified_at"`
	OTPVerifyCount    uint8                  `json:"otp_verify_count,omitempty"" bson:"otp_verify_count"`
	InitialiLinkedAt  time.Time              `json:"initial_linked_at,omitempty"" bson:"initial_linked_at"`
	CreatedAt         time.Time              `json:"created_at,omitempty"" bson:"created_at"`
	DeletedAt         time.Time              `json:"delete_at,omitempty"" bson:"deleted_at"`
	LastModifiedAt    time.Time              `json:"last_modified_at,omitempty"" bson:"last_modified_at"`
}
