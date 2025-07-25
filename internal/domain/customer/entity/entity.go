package entity

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerRespose struct {
	Page      int            `json:"page"`
	Customers []*member.User `json:"customer"`
	Limit     int            `json:"limit"`
	Total     int64          `json:"total"`
}

type User struct {
	ID                    string                       `json:"id"`
	UserCode              string                       `json:"user_code"`
	FullName              string                       `json:"full_name"`
	MotherName            string                       `json:"mother_name"`
	Nationality           string                       `json:"nationality"`
	BirthDate             time.Time                    `json:"birth_date"`
	ResidentialStatus     string                       `json:"residential_status"`
	IssuedDate            time.Time                    `json:"issued_date"`
	PhoneNumber           string                       `json:"phone_number"`
	Gender                member.Gender                `json:"gender"`
	MaritalStatus         member.MaritalStatus         `json:"marital_status"`
	Fayda                 FaydaInfo                    `json:"fayda"`
	Address               member.Address               `json:"address"`
	DocumentFront         string                       `json:"document_front"`
	DocumentBack          string                       `json:"document_back"`
	Photo                 string                       `json:"photo"`
	Signature             string                       `json:"signature"`
	Avatar                string                       `json:"avatar"`
	Email                 string                       `json:"email"`
	PushToken             string                       `json:"push_token"`
	Realm                 member.Realm                 `json:"realm"`
	PermissionGroup       []string                     `json:"permission_group"`
	Permissions           []string                     `json:"permissions"`
	IsAccountBlocked      bool                         `json:"is_account_blocked"`
	IsAccountLinked       bool                         `json:"is_account_linked"`
	MemberType            member.MemberType            `json:"account_branch_type"`
	RegistrationType      member.RegistrationType      `json:"account_type"`
	AccountStatus         member.AccountStatus         `json:"account_status"`
	KYCLevel              uint8                        `json:"kyc_level"`
	KYC                   KYCInfo                      `json:"kyc"`
	IsBranchApproved      bool                         `json:"is_branch_approved"`
	IsVerified            bool                         `json:"is_verified"`
	IsBlocked             bool                         `json:"is_blocked"`
	IsAccessRestricted    bool                         `json:"is_access_restricted"`
	BlockedAt             time.Time                    `json:"blocked_at"`
	RegisterBy            struct{}                     `json:"register_by"`
	LoginAttemptCount     uint8                        `json:"login_attempt_count"`
	LastLoginAttempt      time.Time                    `json:"last_login_attempt"`
	LastOnlineDate        time.Time                    `json:"last_online_date"`
	LastLogin             time.Time                    `json:"last_login"`
	BPSStatus             member.BPSStatus             `json:"bps_reject_status"`
	LoginPIN              member.LoginPIN              `json:"pin"`
	DeviceUUID            string                       `json:"device_uuid"`
	Device                Device                       `json:"device"`
	CustomerNumber        string                       `json:"customer_number"`
	InitialLinkedDate     time.Time                    `json:"initial_linked_date"`
	PrimaryAuthentication member.PrimaryAuthentication `json:"primary_authentication"`
	DeviceStatus          member.DeviceStatus          `json:"device_status"`
	Enabled               bool                         `json:"enabled"`
	IsDeleted             bool                         `json:"is_deleted"`
	PINChangedAt          time.Time                    `json:"pin_changed_at"`
	OTPLastTriedAt        time.Time                    `json:"otp_last_tried_at"`
	OTPLastVerifiedAt     time.Time                    `json:"otp_last_verified_at"`
	OTPVerifyCount        uint8                        `json:"otp_verify_count"`
	InitialLinkedAt       time.Time                    `json:"initial_linked_at"`
	CreatedAt             time.Time                    `json:"created_at"`
	DeletedAt             time.Time                    `json:"deleted_at"`
	LastModifiedAt        time.Time                    `json:"last_modified_at"`
}

type Fayda struct {
	FaydaID          string `json:"id_number" bson:"id_number"`
	FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
	EmploymentStatus string `json:"employment_status" bson:"employement_status"`
	EmployerName     string `json:"employer_name" bson:"employer_name"`
	IssuedBy         string `json:"issued_by" bson:"issued_by"`
	MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
}

type Device struct {
	DevicePlatform      string
	AppVersion          string
	APPInstallationDate time.Time
}

type KYCInfo struct {
	KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed"`
	KYCStatus            member.KYCStatus    `json:"kyc_status"`
	KYCRejectReason      string              `json:"kyc_reject_reason"`
	KYCIsApproved        bool                `json:"kyc_approved"`
	KYCActivityBy        map[string]struct{} `json:"kyc_activity_by"`
}

type FaydaInfo struct {
	FaydaID          string `json:"id_number"`
	FaydaAccessToken string `json:"fayda_access_token"`
	EmploymentStatus string `json:"employment_status"`
	EmployerName     string `json:"employer_name"`
	IssuedBy         string `json:"issued_by"`
	MonthlyIncome    uint64 `json:"monthly_income"`
}
