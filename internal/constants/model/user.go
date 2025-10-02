package model

import (
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID                bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string           `json:"user_code,omitempty" bson:"user_code,omitempty"`
	FullName          string           `json:"full_name,omitempty" bson:"full_name,omitempty"`
	MotherName        string           `json:"mother_name,omitempty" bson:"mother_name,omitempty"`
	Nationality       string           `json:"nationality,omitempty" bson:"nationality,omitempty"`
	BirthDate         time.Time        `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	BranchName        string           `json:"branch_name,omitempty" bson:"branch_name,omitempty"`
	DistrictName      string           `json:"district_name,omitempty" bson:"district_name,omitempty"`
	BranchCode        string           `json:"branch_code,omitempty" bson:"branch_code,omitempty"`
	DistrictCode      string           `json:"district_code,omitempty" bson:"district_code,omitempty"`
	ResidentialStatus string           `json:"residential_status,omitempty" bson:"residential_status,omitempty"`
	IssuedDate        time.Time        `json:"issued_date,omitempty" bson:"isssued_date,omitempty"`
	PhoneNumber       string           `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Gender            constants.Gender `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfileThemeType  string           `json:"profile_theme_type,omitempty" bson:"profile_theme_type,omitempty"`

	Fayda struct {
		FaydaID          string `json:"id_number,omitempty" bson:"id_number,omitempty"`
		FaydaAccessToken string `json:"fayda_access_token,omitempty" bson:"fayda_access_token,omitempty"`
		EmploymentStatus string `json:"employment_status,omitempty" bson:"employement_status,omitempty"`
		EmployerName     string `json:"employer_name,omitempty" bson:"employer_name,omitempty"`
		IssuedBy         string `json:"issued_by,omitempty" bson:"issued_by,omitempty"`
		MonthlyIncome    uint64 `json:"monthly_incode,omitempty" bson:"monthly_incode,omitempty"`
	} `json:"fayda,omitempty" bson:"fayda,omitempty"`
	Address           types.Address           `json:"address,omitempty" bson:"address,omitempty"`
	DocumentFront     string                  `json:"document_front,omitempty" bson:"document_front,omitempty"`
	DocumentBack      string                  `json:"document_back,omitempty" bson:"document_back,omitempty"`
	Photo             string                  `json:"photo,omitempty" bson:"photo,omitempty"`
	Signature         string                  `json:"signature,omitempty" bson:"signature,omitempty"`
	Avatar            string                  `json:"avater,omitempty" bson:"avater,omitempty"`
	Email             string                  `json:"email,omitempty" bson:"email,omitempty"`
	Username          string                  `json:"username,omitempty" bson:"username,omitempty"`
	PushToken         string                  `json:"push_token,omitempty" bson:"push_token,omitempty"`
	Realm             constants.Realm         `json:"realm,omitempty" bson:"realm,omitempty"`
	IsAccountBlocked  bool                    `json:"is_account_blocked,omitempty" bson:"is_account_blocked,omitempty"`
	MainAccount       string                  `json:"main_account,omitempty" bson:"main_account,omitempty"`
	LastMainAccount   string                  `json:"last_main_account,omitempty" bson:"last_main_account,omitempty"`
	AccountLinked     bool                    `json:"account_linked,omitempty" bson:"account_linked,omitempty"`
	LastAccountLinked bool                    `json:"last_account_linked,omitempty" bson:"last_account_linked,omitempty"`
	MemberType        constants.MemberType    `json:"account_branch_type,omitempty" bson:"account_branch_type,omitempty"`
	AccountStatus     constants.AccountStatus `json:"account_status,omitempty" bson:"account_status,omitempty"`
	KYC               struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed,omitempty" bson:"kyc_reject_reason_failed,omitempty"`
		KYCStatus            constants.KYCStatus `json:"kyc_status,omitempty" bson:"kyc_status,omitempty"`
		KYCRejectReason      string              `json:"kyc_reject_reason,omitempty" bson:"kyc_reject_reason,omitempty"`
		KYCApproved          bool                `json:"kyc_approved,omitempty" bson:"kyc_approved,omitempty"`
		KYCActivityBy        map[string]struct{} `json:"kyc_activity_by,omitempty" bson:"kyc_activity_by,omitempty"`
	} `json:"kyc,omitempty" bson:"kyc,omitempty"`
	KYCLevel uint8 `json:"kyc_level" bson:"kyc_level,omitempty"`

	BranchApproved    bool      `json:"branch_approved,omitempty" bson:"branch_approved,omitempty"`
	IsVerified        bool      `json:"is_verified,omitempty" bson:"is_verified,omitempty"`
	IsSelfRegister    bool      `json:"is_self_register,omitempty" bson:"is_self_register,omitempty"`
	BlockedOnCPS      bool      `json:"blocked_on_cps,omitempty" bson:"blocked_on_cps,omitempty"` // default: false
	IsBlocked         bool      `json:"is_blocked,omitempty" bson:"is_blocked,omitempty"`         // default: false
	RegisterBy        struct{}  `json:"register_by,omitempty" bson:"register_By,omitempty"`
	LoginAttemptCount uint8     `json:"login_attempt_count,omitempty" bson:"login_attempt_count,omitempty"`
	NextLoginAttempt  time.Time `json:"next_attempt_count,omitempty" bson:"next_attempt_count,omitempty"`
	LastLoginAttempt  time.Time `json:"last_login_attempt,omitempty" bson:"last_login_attempt,omitempty"`
	LastOnlineDate    time.Time `json:"last_online_date,omitempty" bson:"last_online_date,omitempty"`
	LastLogin         time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`

	BPSStatus           constants.BPSStatus    `json:"bps_reject_status,omitempty" bson:"bps_reject_status,omitempty"` // questioned
	BPSRejectionReason  string                 `json:"bps_reject_reason,omitempty" bson:"bps_reject_reason,omitempty"`
	BPSRejectionField   []string               `json:"bps_reject_failed,omitempty" bson:"bps_reject_failed,omitempty"` // questioned
	LoginPIN            types.LoginPIN         `json:"login_pin,omitempty" bson:"login_pin,omitempty"`
	DeviceUUID          string                 `json:"device_uuid,omitempty" bson:"device_uuid,omitempty"`
	AppVersion          string                 `json:"app_version,omitempty" bson:"app_version,omitempty"`
	Platform            constants.Platform     `json:"platform,omitempty" bson:"platform,omitempty"`
	APPInstallationDate time.Time              `json:"application_installation_date,omitempty" bson:"application_installation_date,omitempty"`
	CustomerNumber      string                 `json:"customer_number,omitempty" bson:"customer_number,omitempty"`
	InitialLinkedDate   time.Time              `json:"initial_linked_date,omitempty" bson:"initiali_linked_date,omitempty"`
	LoanScore           uint16                 `json:"loan_score,omitempty" bson:"loan_score,omitempty"`
	DeviceStatus        constants.DeviceStatus `json:"device_status,omitempty" bson:"device_status,omitempty"`
	Enabled             bool                   `json:"enabled,omitempty" bson:"enabled,omitempty"`
	FirstPinSet         bool                   `json:"first_pin_set,omitempty" bson:"first_pin_set,omitempty"`
	PINChangedAt        time.Time              `json:"pin_changed_at,omitempty" bson:"pin_changed_at,omitempty"`
	OTPVerifyCount      uint8                  `json:"otp_verify_count,omitempty" bson:"otp_verify_count,omitempty"`
	InitialiLinkedAt    time.Time              `json:"initial_linked_at,omitempty" bson:"initial_linked_at,omitempty"`
	CreatedAt           time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	LastModifiedAt      time.Time              `json:"last_modified_at,omitempty" bson:"last_modified_at,omitempty"`
	IsDeleted           bool                   `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}
