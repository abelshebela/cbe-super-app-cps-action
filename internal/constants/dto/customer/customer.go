package customer

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CustomerEnableDTO struct {
	UserOTP string `json:"user_otp"`
}
type CustomerDisableDTO struct {
	IsTemporary   *bool  `json:"is_temporary"`
	DisableReason string `json:"disable_reason"`
}

type CustomerEnableSessionResponse struct {
	Otp string `json:"otp"`
}

type FaydaApproveRequest struct {
	RiskLevel constants.RiskLevel `json:"risk_level"`
}
type SearchCustomerByCIRequest struct {
	CifOrAccountNumber string `json:"cif_or_account_number"`
}
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

type CustomerDetailRespons struct {
	ID             string   `json:"id" bson:"_id"`
	CustomerCode   string   `json:"customer_code" bson:"user_code"`
	FirstName      string   `json:"first_name" bson:"first_name"`
	MiddleName     string   `json:"middle_name" bson:"middle_name"`
	LastName       string   `json:"last_name" bson:"last_name"`
	FullName       string   `json:"full_name" bson:"full_name"`
	PhoneNumber    string   `json:"phone_number" bson:"phone_number"`
	Gender         string   `json:"gender" bson:"gender"`
	AccountNumbers []string `json:"account_numbers" bson:"account_numbers"`
	BranchName     string   `json:"branch_name" bson:"branch_name"`
	DistrictName   string   `json:"district_name" bson:"district_name"`
	MothersName    string   `json:"mothers_name" bson:"mothers_name"`
	Nationality    string   `json:"nationality" bson:"nationality"`
	BirthDate      string   `json:"birth_date" bson:"birth_date"`
	Address        Address  `json:"address" bson:"address"`
	MonthlyIncome  string   `json:"monthly_income" bson:"monthly_income"`
}

type LinkedAccount struct {
	AccountNumber     string `json:"account_number" bson:"account_number"`           // linked_account
	AccountHolderName string `json:"account_holder_name" bson:"account_holder_name"` // linked_account
	AccountType       string `json:"account_type" bson:"account_type"`               // linked_account
	AccountBranchName string `json:"account_branch_name" bson:"account_branch_name"` // members
	AccountBranchCode string `json:"account_branch_code" bson:"account_branch_code"` //linked_account
	IsActive          bool   `json:"is_active" bson:"is_active"`                     /// linked_account
}

type PersonalInfo struct {
	FullName       string `json:"full_name" bson:"full_name"`             // customer_kyc.kyc_data
	Gender         string `json:"gender" bson:"gender"`                   //customer_kyc.kyc_data
	PhoneNumber    string `json:"phone_number" bson:"phone_number"`       //customer_kyc.kyc_data
	Email          string `json:"email " bson:"email"`                    // members
	CustomerNumber string `json:"customer_number" bson:"customer_number"` // members
	DateOfBirth    string `json:"date_of_birth" bson:"date_of_birth"`     // customer_kyc.kyc_data
}

type CustomerDetailResponse struct {
	ID            string          `json:"id" bson:"_id"`
	LinkedAccount []LinkedAccount `json:"linked_account" bson:"linked_account"`
	PersonalInfo  PersonalInfo    `json:"personal_info" bson:"personal_info"`
}

type CustomerListResponse struct {
	ID             string `json:"id" bson:"_id"`
	UserCode       string `json:"user_code" bson:"user_code"`
	UserID         string `json:"user_id" bson:"user_id"`
	FullName       string `json:"full_name" bson:"full_name"`
	PhoneNumber    string `json:"phone_number" bson:"phone_number"`
	CustomerNumber string `json:"customer_number" bson:"customer_number"`
	BranchCode     string `json:"branch_code" bson:"branch_code"`
	Gender         string `json:"gender" bson:"gender"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	Avatar         string `json:"avatar" bson:"avatar"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	IsBlocked      bool   `json:"is_blocked" bson:"is_blocked"`
}

type FindCustomerByIDResponse struct {
	ID                   bson.ObjectID              `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode             string                     `json:"user_code" bson:"user_code"`
	FullName             string                     `json:"full_name" bson:"full_name"`
	BranchCode           string                     `json:"branch_code" bson:"branch_code"`
	ActivationBranchCode string                     `json:"activation_branch_code" bson:"activation_branch_code"`
	PhoneNumber          string                     `json:"phone_number" bson:"phone_number"`
	Language             string                     `json:"language" bson:"language"`
	Avatar               string                     `json:"avatar" bson:"avatar"`
	Email                string                     `json:"email" bson:"email"`
	PushToken            string                     `json:"push_token" bson:"push_token"`
	CustomerNumber       string                     `json:"customer_number" bson:"customer_number"`
	UserCategory         string                     `json:"user_category" bson:"user_category"`
	Industry             string                     `json:"industry" bson:"industry"`
	Sector               string                     `json:"sector" bson:"sector"`
	Username             string                     `json:"username" bson:"username"`
	Ownership            string                     `json:"ownership" bson:"ownership"`
	CustomerSegment      string                     `json:"customer_segment" bson:"customer_segment"`
	BlockedReason        string                     `json:"blocked_reason" bson:"blocked_reason,omitempty"` // -- optional
	DeviceUUID           string                     `json:"device_uuid" bson:"device_uuid"`
	AppVersion           string                     `json:"app_version" bson:"app_version"`
	Gender               constants.Gender           `json:"gender" bson:"gender"`
	MemberType           constants.MemberType       `json:"account_branch_type" bson:"account_branch_type"`
	Platform             constants.Platform         `json:"platform" bson:"platform"`
	DeviceStatus         constants.DeviceStatus     `json:"device_status" bson:"device_status"`
	OnboardingMethod     constants.OnboardingMethod `json:"onboarding_method" bson:"onboarding_method"`
	KYCLevel             uint8                      `json:"kyc_level" bson:"kyc_level"`
	BlockedOn            member.BlockedOn           `json:"blocked_on" bson:"blocked_on,omitempty"` // -- optional
	IsUSSDEnabled        bool                       `json:"is_ussd_enabled" bson:"is_ussd_enabled"`
	ISuperappEnabled     bool                       `json:"is_superapp_enabled" bson:"is_superapp_enabled"`
	LastLogin            time.Time                  `json:"last_login" bson:"last_login"`
	APPInstallationDate  time.Time                  `json:"application_installation_date" bson:"application_installation_date"`
	Config               member.Config              `json:"config" bson:"config"`
	CreatedAt            time.Time                  `json:"created_at" bson:"created_at"`
	ExpiryAt             time.Time                  `json:"expiry_at" bson:"expiry_at"`
	LastModifiedAt       time.Time                  `json:"last_modified_at" bson:"last_modified_at"`
	IsBlocked            bool                       `json:"is_blocked" bson:"is_blocked"` // default: false
	IsLocked             bool                       `json:"is_locked" bson:"is_locked"`   // default: false
	Enabled              bool                       `json:"enabled" bson:"enabled"`
	FirstPinSet          bool                       `json:"first_pin_set" bson:"first_pin_set"`
	IsActivated          bool                       `json:"is_activated" bson:"is_activated"`
}

type CustomerActionLogResponse struct {
	ActionCode   string `json:"action_code" bson:"action_code"`
	MakerName    string `json:"maker_name" bson:"maker_name"`
	ActionReason struct {
		ActionType string `json:"action_type" bson:"action_type"` //  reject, enable, disable
		ActionNote string `json:"action_note" bson:"action_note"`
	} `json:"action_reason" bson:"action_reason"`

	RequestAction  string    `json:"request_action" bson:"request_action"` //action type
	ServiceName    string    `json:"service_name" bson:"service_name"`     //request action
	Status         string    `json:"status" bson:"status"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
}
