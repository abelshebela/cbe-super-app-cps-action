package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ArchivedUser struct {
	ID                bson.ObjectID           `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string                  `json:"user_code" bson:"user_code"`
	FullName          string                  `json:"full_name" bson:"full_name"`
	MotherName        string                  `json:"mother_name" bson:"mother_name"`
	Nationality       string                  `json:"nationality" bson:"nationality"`
	BirthDate         time.Time               `json:"birth_date" bson:"birth_date"`
	ResidentialStatus string                  `json:"residential_status" bson:"residential_status"`
	IssuedDate        time.Time               `json:"issued_date" bson:"isssued_date"`
	PhoneNumber       string                  `json:"phone_number" bson:"phone_number"`
	Gender            constants.Gender        `json:"gender" bson:"gender"`
	MaritalStatus     constants.MaritalStatus `json:"marital_status" bson:"marital_status"`
	Fayda             struct {
		FaydaID          string `json:"id_number" bson:"id_number"`
		FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
		EmploymentStatus string `json:"employment_status" bson:"employement_status"`
		EmployerName     string `json:"employer_name" bson:"employer_name"`
		IssuedBy         string `json:"issued_by" bson:"issued_by"`
		MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
	} `json:"fayda" bson:"fayda"`
	Address          types.Address              `json:"address" bson:"address"`
	DocumentFront    string                     `json:"document_front" bson:"document_front"`
	DocumentBack     string                     `json:"document_back" bson:"document_back"`
	Photo            string                     `json:"photo" bson:"photo"`
	Signature        string                     `json:"signature" bson:"signature"`
	Avatar           string                     `json:"avater" bson:"avater"`
	Email            string                     `json:"email" bson:"email"`
	PushToken        string                     `json:"push_token" bson:"push_token"`
	Realm            constants.Realm            `json:"realm" bson:"realm"`
	PermissionGroup  []bson.ObjectID            `json:"permission_group" bson:"permission_group"`
	Permissions      []bson.ObjectID            `json:"permissions" bson:"persmissions"`
	IsAccountBlocked bool                       `json:"is_account_blocked" bson:"is_account_blocked"`
	IsAccountLinked  bool                       `json:"is_account_linked" bson:"is_account_linked"`
	MemberType       constants.MemberType       `json:"account_branch_type" bson:"account_branch_type"`
	RegistrationType constants.RegistrationType `json:"account_type" bson:"account_type"`
	AccountStatus    constants.AccountStatus    `json:"account_status" bson:"account_status"`
	KYCLevel         uint8                      `json:"kyc_level" bson:"kyc_level"`
	KYC              struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
		KYCStatus            constants.KYCStatus `json:"kyc_status" bson:"kyc_status"`
		KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
		KYCIsApproved        bool                `json:"kyc_approved" bson:"kyc_approved"`
		KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
	} `json:"kyc" bson:"kyc"`
	IsBranchApproved   bool      `json:"is_branch_approved" bson:"is_branch_approved"`
	IsVerified         bool      `json:"is_verfied" bson:"is_verfied"`
	IsBlocked          bool      `json:"is_blocked" bson:"is_blocked"`                     // default: false
	IsAccessRestricted bool      `json:"is_access_restricted" bson:"is_access_restricted"` // default: false
	BlockedAt          time.Time `json:"blocked_at" bson:"blocked_at"`
	RegisterBy         struct{}  `json:"register_by" bson:"register_By"`
	LoginAttemptCount  uint8     `json:"login_attempt_count" bson:"login_attempt_count"`
	LastLoginAttempt   time.Time `json:"last_login_attempt" bson:"last_login_attempt"`
	LastOnlineDate     time.Time `json:"last_online_date" bson:"last_online_date"`
	LastLogin          time.Time `json:"last_login" bson:"last_login"`

	// move BPS Status and realted field to new collection
	BPSStatus constants.BPSStatus `json:"bps_reject_status" bson:"bps_reject_status"`
	// BPSRejectionReason string    `json:"bps_reject_reason" bson:"bps_reject_reason"`
	// BPSRejectionField  []string  `json:"bps_reject_failed" bson:"bps_reject_failed"`

	LoginPIN   types.LoginPIN `json:"pin" bson:"login_pin"`
	DeviceUUID string         `json:"device_uuid" bson:"device_uuid"`
	Device     struct {
		DevicePlatform      string    `json:"device_platform" bson:"device_platform"`
		AppVersion          string    `json:"app_version" bson:"app_version"`
		APPInstallationDate time.Time `json:"application_installation_date" bson:"application_installation_date"`
	} `json:"device" bson:"device"`
	CustomerNumber        string                          `json:"customer_number" bson:"customer_number"`
	InitialLinkedDate     time.Time                       `json:"initial_linked_date" bson:"initiali_linked_date"`
	PrimaryAuthentication constants.PrimaryAuthentication `json:"primary_authentication" bson:"primary_authentication"`
	// LoanScore             uint8                 `json:"loan_score" bson:"loan_score"`
	DeviceStatus      constants.DeviceStatus `json:"device_status" bson:"device_status"`
	Enabled           bool                   `json:"enabled" bson:"enabled"`
	IsDeleted         bool                   `json:"is_deleted" bson:"is_deleted"`
	PINChangedAt      time.Time              `json:"pin_changed_at" bson:"pin_changed_at"`
	OTPLastTriedAt    time.Time              `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt time.Time              `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	OTPVerifyCount    uint8                  `json:"otp_verify_count" bson:"otp_verify_count"`
	InitialiLinkedAt  time.Time              `json:"initial_linked_at" bson:"initial_linked_at"`
	CreatedAt         time.Time              `json:"created_at" bson:"created_at"`
	DeletedAt         time.Time              `json:"delete_at" bson:"deleted_at"`
	LastModifiedAt    time.Time              `json:"last_modified_at" bson:"last_modified_at"`
}
