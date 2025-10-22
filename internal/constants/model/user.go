package model

import (
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID                bson.ObjectID    `json:"id" bson:"_id,omitempty"`
	UserCode          string           `json:"user_code" bson:"user_code,omitempty"`
	FullName          string           `json:"full_name" bson:"full_name,omitempty"`
	MotherName        string           `json:"mother_name" bson:"mother_name,omitempty"`
	Nationality       string           `json:"nationality" bson:"nationality,omitempty"`
	BirthDate         time.Time        `json:"birth_date" bson:"birth_date,omitempty"`
	BranchName        string           `json:"branch_name" bson:"branch_name,omitempty"`
	DistrictName      string           `json:"district_name" bson:"district_name,omitempty"`
	BranchCode        string           `json:"branch_code" bson:"branch_code,omitempty"`
	DistrictCode      string           `json:"district_code" bson:"district_code,omitempty"`
	ResidentialStatus string           `json:"residential_status" bson:"residential_status,omitempty"`
	IssuedDate        time.Time        `json:"issued_date" bson:"isssued_date,omitempty"`
	PhoneNumber       string           `json:"phone_number" bson:"phone_number,omitempty"`
	Gender            constants.Gender `json:"gender" bson:"gender,omitempty"`
	ProfileThemeType  string           `json:"profile_theme_type" bson:"profile_theme_type,omitempty"`

	Fayda struct {
		FaydaID          string `json:"id_number" bson:"id_number,omitempty"`
		FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token,omitempty"`
		EmploymentStatus string `json:"employment_status" bson:"employement_status,omitempty"`
		EmployerName     string `json:"employer_name" bson:"employer_name,omitempty"`
		IssuedBy         string `json:"issued_by" bson:"issued_by,omitempty"`
		MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode,omitempty"`
	} `json:"fayda" bson:"fayda,omitempty"`

	Address           types.Address           `json:"address" bson:"address,omitempty"`
	DocumentFront     string                  `json:"document_front" bson:"document_front,omitempty"`
	DocumentBack      string                  `json:"document_back" bson:"document_back,omitempty"`
	Photo             string                  `json:"photo" bson:"photo,omitempty"`
	Signature         string                  `json:"signature" bson:"signature,omitempty"`
	Avatar            string                  `json:"avater" bson:"avater,omitempty"`
	Email             string                  `json:"email" bson:"email,omitempty"`
	Username          string                  `json:"username" bson:"username,omitempty"`
	PushToken         string                  `json:"push_token" bson:"push_token,omitempty"`
	Realm             constants.Realm         `json:"realm" bson:"realm,omitempty"`
	IsAccountBlocked  bool                    `json:"is_account_blocked" bson:"is_account_blocked,omitempty"`
	MainAccount       string                  `json:"main_account" bson:"main_account,omitempty"`
	LastMainAccount   string                  `json:"last_main_account" bson:"last_main_account,omitempty"`
	AccountLinked     bool                    `json:"account_linked" bson:"account_linked,omitempty"`
	LastAccountLinked bool                    `json:"last_account_linked" bson:"last_account_linked,omitempty"`
	MemberType        constants.MemberType    `json:"account_branch_type" bson:"account_branch_type,omitempty"`
	AccountStatus     constants.AccountStatus `json:"account_status" bson:"account_status,omitempty"`

	KYC struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed,omitempty"`
		KYCStatus            constants.KYCStatus `json:"kyc_status" bson:"kyc_status,omitempty"`
		KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason,omitempty"`
		KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved,omitempty"`
		KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by,omitempty"`
	} `json:"kyc" bson:"kyc,omitempty"`

	KYCLevel uint8 `json:"kyc_level" bson:"kyc_level,omitempty"`

	BranchApproved    bool      `json:"branch_approved" bson:"branch_approved,omitempty"`
	IsVerified        bool      `json:"is_verified" bson:"is_verified,omitempty"`
	IsSelfRegister    bool      `json:"is_self_register" bson:"is_self_register,omitempty"`
	BlockedOnCPS      bool      `json:"blocked_on_cps" bson:"blocked_on_cps,omitempty"`
	IsBlocked         bool      `json:"is_blocked" bson:"is_blocked,omitempty"`
	RegisterBy        struct{}  `json:"register_by" bson:"register_By,omitempty"`
	LoginAttemptCount uint8     `json:"login_attempt_count" bson:"login_attempt_count,omitempty"`
	NextLoginAttempt  time.Time `json:"next_attempt_count" bson:"next_attempt_count,omitempty"`
	LastLoginAttempt  time.Time `json:"last_login_attempt" bson:"last_login_attempt,omitempty"`
	LastOnlineDate    time.Time `json:"last_online_date" bson:"last_online_date,omitempty"`
	LastLogin         time.Time `json:"last_login" bson:"last_login,omitempty"`

	BPSStatus           constants.BPSStatus    `json:"bps_reject_status" bson:"bps_reject_status,omitempty"`
	BPSRejectionReason  string                 `json:"bps_reject_reason" bson:"bps_reject_reason,omitempty"`
	BPSRejectionField   []string               `json:"bps_reject_failed" bson:"bps_reject_failed,omitempty"`
	LoginPIN            types.LoginPIN         `json:"login_pin" bson:"login_pin,omitempty"`
	DeviceUUID          string                 `json:"device_uuid" bson:"device_uuid,omitempty"`
	AppVersion          string                 `json:"app_version" bson:"app_version,omitempty"`
	Platform            constants.Platform     `json:"platform" bson:"platform,omitempty"`
	APPInstallationDate time.Time              `json:"application_installation_date" bson:"application_installation_date,omitempty"`
	CustomerNumber      string                 `json:"customer_number" bson:"customer_number,omitempty"`
	InitialLinkedDate   time.Time              `json:"initial_linked_date" bson:"initiali_linked_date,omitempty"`
	LoanScore           uint16                 `json:"loan_score" bson:"loan_score,omitempty"`
	DeviceStatus        constants.DeviceStatus `json:"device_status" bson:"device_status,omitempty"`
	Enabled             bool                   `json:"enabled" bson:"enabled"`
	FirstPinSet         bool                   `json:"first_pin_set" bson:"first_pin_set,omitempty"`
	PINChangedAt        time.Time              `json:"pin_changed_at" bson:"pin_changed_at,omitempty"`
	OTPVerifyCount      uint8                  `json:"otp_verify_count" bson:"otp_verify_count,omitempty"`
	InitialiLinkedAt    time.Time              `json:"initial_linked_at" bson:"initial_linked_at,omitempty"`
	CreatedAt           time.Time              `json:"created_at" bson:"created_at,omitempty"`
	LastModifiedAt      time.Time              `json:"last_modified_at" bson:"last_modified_at,omitempty"`
	IsDeleted           bool                   `json:"is_deleted" bson:"is_deleted,omitempty"`
}
