package dto

import (
	"time"

	"github.com/elst-bank/elst-bank-backend/internal/constants"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Address struct {
	Zone        string `json:"zone" bson:"zone"`
	Wereda      string `json:"wereda" bson:"wereda"`
	Kebele      string `json:"kebele" bson:"kebele"`
	Region      string `json:"region" bson:"region"`
	City        string `json:"city" bson:"city"`
	SubCity     string `json:"sub_city" bson:"sub_city"`
	StreetName  string `json:"street_name" bson:"street_name"`
	HouseNumber string `json:"house_number" bson:"house_number"`
}

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}

type DeviceLinkHistroy struct {
	ID              bson.ObjectID        `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          bson.ObjectID        `json:"user_id" bson:"user_id"`
	DeviceType      constants.DeviceType `json:"device_type" bson:"device_type"`
	DeviceUUID      string               `json:"device_uuid" bson:"device_uuid"`
	DeviceName      string               `json:"device_name" bson:"device_name"`
	DeviceOSVersion string               `json:"device_os_version" bson:"device_os_version"`
	LinkedAt        time.Time            `json:"linked_at" bson:"linked_at"`
	UnLinkedAt      time.Time            `json:"unlinked_at" bson:"unlinked_at"`
}

type User struct {
	ID                bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string           `json:"user_code" bson:"user_code"`
	FullName          string           `json:"full_name" bson:"full_name"`
	MotherName        string           `json:"mother_name" bson:"mother_name"`
	Nationality       string           `json:"nationality" bson:"nationality"`
	BirthDate         time.Time        `json:"birth_date" bson:"birth_date"`
	BranchName        string           `json:"branch_name" bson:"branch_name"`
	DistrictName      string           `json:"district_name" bson:"district_name"`
	BranchCode        string           `json:"branch_code" bson:"branch_code"`
	DistrictCode      string           `json:"district_code" bson:"district_code"`
	ResidentialStatus string           `json:"residential_status" bson:"residential_status"`
	IssuedDate        time.Time        `json:"issued_date" bson:"isssued_date"`
	PhoneNumber       string           `json:"phone_number" bson:"phone_number"`
	Gender            constants.Gender `json:"gender" bson:"gender"`
	ProfileThemeType  string           `json:"profile_theme_type" bson:"profile_theme_type"`

	Fayda struct {
		FaydaID          string `json:"id_number" bson:"id_number"`
		FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
		EmploymentStatus string `json:"employment_status" bson:"employement_status"`
		EmployerName     string `json:"employer_name" bson:"employer_name"`
		IssuedBy         string `json:"issued_by" bson:"issued_by"`
		MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
	} `json:"fayda" bson:"fayda"`
	Address Address `json:"address" bson:"address"`
	// AndOrCustomerName []string         `json:"and_or_customer_name" bson:"and_or_customer_name"`
	DocumentFront string `json:"document_front" bson:"document_front"`
	DocumentBack  string `json:"document_back" bson:"document_back"`
	Photo         string `json:"photo" bson:"photo"`
	Signature     string `json:"signature" bson:"signature"`
	Avatar        string `json:"avater" bson:"avater"`
	Email         string `json:"email" bson:"email"`
	Username      string `json:"username" bson:"username"`
	// AndOrStatus       bool             `json:"and_or_status" bson:"and_or_status"`
	PushToken string          `json:"push_token" bson:"push_token"`
	Realm     constants.Realm `json:"realm" bson:"realm"`
	// PermissionGroup   []bson.ObjectID  `json:"permission_group" bson:"permission_group"`
	// Permissions       []bson.ObjectID  `json:"permissions" bson:"persmissions"`
	IsAccountBlocked  bool                 `json:"is_account_blocked" bson:"is_account_blocked"`
	MainAccount       string               `json:"main_account" bson:"main_account"`
	LastMainAccount   string               `json:"last_main_account" bson:"last_main_account"`
	AccountLinked     bool                 `json:"account_linked" bson:"account_linked"`
	LastAccountLinked bool                 `json:"last_account_linked" bson:"last_account_linked"`
	MemberType        constants.MemberType `json:"account_branch_type" bson:"account_branch_type"`
	// RegistrationType  RegistrationType `json:"account_type" bson:"account_type"`
	AccountStatus constants.AccountStatus `json:"account_status" bson:"account_status"`
	KYC           struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
		KYCStatus            constants.KYCStatus `json:"kyc_status" bson:"kyc_status"`
		KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
		KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
		KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
	} `json:"kyc" bson:"kyc"`
	KYCLevel uint8 `json:"level" bson:"level"`

	BranchApproved bool `json:"branch_approved" bson:"branch_approved"`
	IsVerified     bool `json:"is_verfied" bson:"is_verfied"`
	IsSelfRegister bool `json:"is_self_register" bson:"is_self_register"`
	BlockedOnCPS   bool `json:"blocked_on_cps" bson:"blocked_on_cps"` // default: false
	IsBlocked      bool `json:"is_blocked" bson:"is_blocked"`         // default: false
	// IsDetached        bool      `json:"is_detached" bson:"is_detached"`       // default: false
	// DetachedAt        time.Time `json:"detached_at" bson:"detached_at"`
	RegisterBy        struct{}  `json:"register_by" bson:"register_By"`
	LoginAttemptCount uint8     `json:"login_attempt_count" bson:"login_attempt_count"`
	NextLoginAttempt  time.Time `json:"next_attempt_count" bson:"next_attempt_count"`
	LastLoginAttempt  time.Time `json:"last_login_attempt" bson:"last_login_attempt"`
	LastOnlineDate    time.Time `json:"last_online_date" bson:"last_online_date"`
	LastLogin         time.Time `json:"last_login" bson:"last_login"`

	// move BPS Status and realted field to new collection
	BPSStatus           constants.BPSStatus  `json:"bps_reject_status" bson:"bps_reject_status"` // questioned
	BPSRejectionReason  string               `json:"bps_reject_reason" bson:"bps_reject_reason"`
	BPSRejectionField   []string             `json:"bps_reject_failed" bson:"bps_reject_failed"` // questioned
	LoginPIN            constants.LoginPIN   `json:"login_pin" bson:"login_pin"`
	DeviceUUID          string               `json:"device_uuid" bson:"device_uuid"`
	AppVersion          string               `json:"app_version" bson:"app_version"`
	Platform            constants.DeviceType `json:"platform" bson:"platform"`
	APPInstallationDate time.Time            `json:"application_installation_date" bson:"application_installation_date"`
	CustomerNumber      string               `json:"customer_number" bson:"customer_number"`
	InitialLinkedDate   time.Time            `json:"initial_linked_date" bson:"initiali_linked_date"`
	// PrimaryAuthentication PrimaryAuthentication `json:"primary_authentication" bson:"primary_authentication"`
	LoanScore    uint16                 `json:"loan_score" bson:"loan_score"`
	DeviceStatus constants.DeviceStatus `json:"device_status" bson:"device_status"`
	// SessionExpirsOn       time.Time             `json:"session_expires_on" bson:"session_expires_on"`
	Enabled     bool `json:"enabled" bson:"enabled"`
	FirstPinSet bool `json:"first_pin_set" bson:"first_pin_set"`
	// IsDeleted             bool                  `json:"is_deleted" bson:"is_deleted"`
	PINChangedAt time.Time `json:"pin_changed_at" bson:"pin_changed_at"`
	// OTPLastTriedAt        time.Time             `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	// OTPLastVerifiedAt     time.Time             `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	OTPVerifyCount uint8 `json:"otp_verify_count" bson:"otp_verify_count"`
	// PINHistory            []string              `json:"pin_history" bson:"pin_history"`
	InitialiLinkedAt time.Time `json:"initial_linked_at" bson:"initial_linked_at"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	// DeletedAt             time.Time             `json:"delete_at" bson:"deleted_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
	// LinkedAccount         []LinkedAccount       `json:"linked_account" bson:"linked_account"`
	// OrganizationID        []LinkedAccount       `json:"organization_id" bson:"organization_id"`
}
