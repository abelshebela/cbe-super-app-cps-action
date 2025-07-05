package users

import (
	"time"

	"cbe-super-app-member-users/pkgs/entities/enums"
	"cbe-super-app-member-users/pkgs/entities/type_definition"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FullName struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type LinkedAccountDetail struct {
	AccountNumber     string `json:"account_number"`
	AccountBranchCode string `json:"account_branch_code"`
	LinkedBranch      string `json:"linked_branch"`
	IsAccountActive   bool   `json:"is_account_active"`
	LinkedStatus      bool   `json:"linked_status"`
	CurrencyCode      string `json:"currency_code"`
}

type DeviceInfo struct {
	DeviceUUID string `json:"device_uuid"`
	AppVersion string `json:"app_version"`
}

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}

type User struct {
	ID                bson.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string              `json:"user_code" bson:"user_code"`
	FullName          string              `json:"full_name" bson:"full_name"`
	MotherName        string              `json:"mother_name" bson:"mother_name"`
	Nationality       string              `json:"nationality" bson:"nationality"`
	BirthDate         time.Time           `json:"birth_date" bson:"birth_date"`
	BranchName        string              `json:"branch_name" bson:"branch_name"`
	DistrictName      string              `json:"district_name" bson:"district_name"`
	BranchCode        string              `json:"branch_code" bson:"branch_code"`
	DistrictCode      string              `json:"district_code" bson:"district_code"`
	ResidentialStatus string              `json:"residential_status" bson:"residential_status"`
	IssuedDate        time.Time           `json:"issued_date" bson:"isssued_date"`
	PhoneNumber       string              `json:"phone_number" bson:"phone_number"`
	Gender            enums.Gender        `json:"gender" bson:"gender"`
	MartialStatus     enums.MartialStatus `json:"marital_status" bson:"marital_status"`
	Fayda             struct {
		FaydaID          string `json:"id_number" bson:"id_number"`
		FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
		EmploymentStatus string `json:"employment_status" bson:"employement_status"`
		EmployerName     string `json:"employer_name" bson:"employer_name"`
		IssuedBy         string `json:"issued_by" bson:"issued_by"`
		MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
	} `json:"fayda" bson:"fayda"`
	Address           type_definition.Address `json:"address" bson:"address"`
	AndOrCustomerName []string                `json:"and_or_customer_name" bson:"and_or_customer_name"`
	DocumentFront     string                  `json:"document_front" bson:"document_front"`
	DocumentBack      string                  `json:"document_back" bson:"document_back"`
	Photo             string                  `json:"photo" bson:"photo"`
	Signature         string                  `json:"signature" bson:"signature"`
	Avatar            string                  `json:"avater" bson:"avater"`
	Email             string                  `json:"email" bson:"email"`
	UserName          string                  `json:"UserName" bson:"UserName"`
	AndOrStatus       bool                    `json:"and_or_status" bson:"and_or_status"`
	PushToken         string                  `json:"push_token" bson:"push_token"`
	Realm             enums.Realm             `json:"realm" bson:"realm"`
	PermissionGroup   []bson.ObjectID         `json:"permission_group" bson:"permission_group"`
	Permissions       []bson.ObjectID         `json:"permissions" bson:"persmissions"`
	IsAccountBlocked  bool                    `json:"is_account_blocked" bson:"is_account_blocked"`
	MainAccount       string                  `json:"main_account" bson:"main_account"`
	LastMainAccount   string                  `json:"last_main_account" bson:"last_main_account"`
	AccountLinked     bool                    `json:"account_linked" bson:"account_linked"`
	LastAccountLinked bool                    `json:"last_account_linked" bson:"last_account_linked"`
	MemberType        enums.MemberType        `json:"account_branch_type" bson:"account_branch_type"`
	RegistrationType  enums.RegistrationType  `json:"account_type" bson:"account_type"`
	AccountStatus     enums.AccountStatus     `json:"account_status" bson:"account_status"`
	KYC               struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
		KYCStatus            enums.KYCStatus     `json:"kyc_status" bson:"kyc_status"`
		KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
		KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
		KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
		KYCLevel             uint8               `json:"level" bson:"level"`
	} `json:"kyc" bson:"kyc"`
	BranchApproved    bool      `json:"branch_approved" bson:"branch_approved"`
	IsVerified        bool      `json:"is_verfied" bson:"is_verfied"`
	IsSelfRegister    bool      `json:"is_self_register" bson:"is_self_register"`
	BlockedOnCPS      bool      `json:"blocked_on_cps" bson:"blocked_on_cps"` // default: false
	IsBlocked         bool      `json:"is_blocked" bson:"is_blocked"`         // default: false
	IsDetached        bool      `json:"is_detached" bson:"is_detached"`       // default: false
	DetachedAt        time.Time `json:"detached_at" bson:"detached_at"`
	RegisterBy        struct{}  `json:"register_by" bson:"register_By"`
	LoginAttemptCount uint8     `json:"login_attempt_count" bson:"login_attempt_count"`
	LastLoginAttempt  time.Time `json:"last_login_attempt" bson:"last_login_attempt"`
	LastOnlineDate    time.Time `json:"last_online_date" bson:"last_online_date"`
	LastLogin         time.Time `json:"last_login" bson:"last_login"`

	// move BPS Status and realted field to new collection
	BPSStatus          enums.BPSStatus          `json:"bps_reject_status" bson:"bps_reject_status"`
	BPSRejectionReason string                   `json:"bps_reject_reason" bson:"bps_reject_reason"`
	BPSRejectionField  []string                 `json:"bps_reject_failed" bson:"bps_reject_failed"`
	LoginPIN           type_definition.LoginPIN `json:"login_pin" bson:"login_pin"`
	Device             struct {
		DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
		AppVersion string `json:"app_version" bson:"app_version"`
	} `json:"device" bson:"device"`
	APPInstallationDate   time.Time                   `json:"application_installation_date" bson:"application_installation_date"`
	CustomerNumber        string                      `json:"customer_number" bson:"customer_number"`
	InitialLinkedDate     time.Time                   `json:"initial_linked_date" bson:"initiali_linked_date"`
	PrimaryAuthentication enums.PrimaryAuthentication `json:"primary_authentication" bson:"primary_authentication"`
	LoanScore             uint16                      `json:"loan_score" bson:"loan_score"`
	DeviceStatus          enums.DeviceStatus          `json:"device_status" bson:"device_status"`
	SessionExpirsOn       time.Time                   `json:"session_expires_on" bson:"session_expires_on"`
	PinStatus             string                      `json:"pin_status" bson:"pin_status"`
	Enabled               bool                        `json:"enabled" bson:"enabled"`
	IsDeleted             bool                        `json:"is_deleted" bson:"is_deleted"`
	PINChangedAt          time.Time                   `json:"pin_changed_at" bson:"pin_changed_at"`
	OTPLastTriedAt        time.Time                   `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt     time.Time                   `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	OTPVerifyCount        uint8                       `json:"otp_verify_count" bson:"otp_verify_count"`
	PINHistory            []string                    `json:"pin_history" bson:"pin_history"`
	InitialiLinkedAt      time.Time                   `json:"initial_linked_at" bson:"initial_linked_at"`
	CreatedAt             time.Time                   `json:"created_at" bson:"created_at"`
	DeletedAt             time.Time                   `json:"delete_at" bson:"deleted_at"`
	LastModifiedAt        time.Time                   `json:"last_modified_at" bson:"last_modified_at"`
}

type UserEmail struct {
	ID    string
	Email string `json:"email" bson:"email"`
}

type LinkedAccountResponse struct {
	UserID         string                `json:"user_id"`
	FullName       string                `json:"full_name"`
	LinkedAccounts []LinkedAccountDetail `json:"linked_accounts"`
}

type GenerateOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type OTPRequest struct {
	UserID string
	Email  string
}

type OTPVerification struct {
	UserID string
	Email  string
	OTP    string
}

type OTPRecord struct {
	ID         string `json:"id" bson:"_id"`
	UserCode   string
	UserID     string
	Email      string
	OTP        string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	UserRealm  string
	DeviceUUID string
	OTPFor     string
}

type ChangePinRequest struct {
	UserID string `json:"user_id"`
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

type HQ struct {
	ID                   bson.ObjectID   `bson:"_id" json:"id"`
	UniqueID             string          `bson:"unique_id" json:"unique_id"`
	Name                 string          `bson:"name" json:"name"`
	Address              string          `bson:"address" json:"address"`
	PhoneNumber          string          `bson:"phoneNumber" json:"phoneNumber"`
	Email                string          `bson:"email" json:"email"`
	LinkedAccounts       []LinkedAccount `bson:"linkedAccounts" json:"linkedAccounts"`
	LatestiOSVersion     string          `bson:"latestiOSVersion" json:"latestiOSVersion"`
	LatestAndroidVersion string          `bson:"latestAndroidVersion" json:"latestAndroidVersion"`
	ArchiveExpiry        uint            `json:"archive_expiry" bson:"archive_expiry"`
	BlockTime            uint            `json:"block_time" bson:"block_time"`
	BlockTimeStatus      string          `bson:"block_time_status" json:"block_time_status"`
	ArchiveTime          uint            `bson:"archive_time" json:"archive_time"`
	ArchiveTimeStatus    string          `bson:"archive_time_status" json:"archive_time_status"`
	Enabled              bool            `bson:"enabled" json:"enabled"`
	IsDeleted            bool            `bson:"isDeleted" json:"isDeleted"`
	CreatedAt            time.Time       `bson:"createdAt" json:"createdAt"`
	LastModified         time.Time       `bson:"lastModified" json:"lastModified"`
}

type RegistrationType string

const (
	RegistrationTypeNew    RegistrationType = "new"
	RegistrationTypeLinked RegistrationType = "linked"
)

type LinkedAccount struct {
	ID                bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID            bson.ObjectID    `json:"user_id" bson:"user_id"`
	CustomerNumber    string           `json:"customer_number" bson:"customer_number"` // cif
	AccountNumber     string           `json:"account_number" bson:"account_number"`
	AccountHolderName string           `json:"account_holder_name" bson:"account_holder_name"`
	AccountType       string           `json:"account_type" bson:"account_type"`
	BranchCode        string           `json:"branch_code" bson:"branch_code"`
	LinkedStatus      bool             `json:"linked_status" bson:"linked_status"`
	LastLinkedStatus  bool             `json:"last_linked_status" bson:"last_linked_status"`
	LinkedAt          time.Time        `json:"linked_at" bson:"linked_at"`
	LinkedBranch      string           `json:"linker_branch" bson:"linker_branch"`
	RegistrationType  RegistrationType `json:"registration_type" bson:"registration_type"`
	IsAccountActive   bool             `json:"is_account_active" bson:"is_account_active"`
	AndOrStatus       bool             `json:"and_or_status" bson:"and_or_status"`
	AccountBranchCode string           `json:"account_branch_code" bson:"account_branch_code"`
	CurrencyCode      string           `json:"currency" bson:"curreny"`
	IsMain            bool             `json:"is_main" bson:"is_main"` // default: false, first account: true
	MakerAndChecker   struct {
		Linkers struct {
			Maker   string `json:"maker" bson:"maker"`
			Checker string `json:"checker" bson:"checker"`
		} `json:"linkers" bson:"linkers"`
		Unlinkers struct {
			Maker   string `json:"maker" bson:"maker"`
			Checker string `json:"checker" bson:"checker"`
		} `json:"unlinkers" bson:"unlinkers,omitempty"`
	} `json:"maker_and_checker" bson:"maker_and_checker"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at"  bson:"updated_at"`
}

// RegistrationRecord represents a pending user registration
type RegistrationRecord struct {
	ID          string    `bson:"_id" json:"id"`
	PhoneNumber string    `bson:"phone_number" json:"phone_number"`
	DeviceUUID  string    `bson:"device_uuid" json:"device_uuid"`
	Platform    string    `bson:"platform" json:"platform"`
	OTP         string    `bson:"otp" json:"otp"`
	OTPFor      string    `bson:"otp_for" json:"otp_for"`
	Status      string    `bson:"status" json:"status"`
	ExpiresAt   time.Time `bson:"expires_at" json:"expires_at"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	Attempts    int       `bson:"attempts" json:"attempts"`
	MaxAttempts int       `bson:"max_attempts" json:"max_attempts"`
}

type PinResetSession struct {
	ID               string    `bson:"_id" json:"id"`
	UserID           string    `bson:"user_id" json:"user_id"`
	PhoneNumber      string    `bson:"phone_number" json:"phone_number"`
	DeviceUUID       string    `bson:"device_uuid" json:"device_uuid"`
	OTP              string    `bson:"otp" json:"otp"`
	OTPFor           string    `bson:"otp_for" json:"otp_for"`
	Status           string    `bson:"status" json:"status"` // pending, verified, completed, expired
	ExpiresAt        time.Time `bson:"expires_at" json:"expires_at"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	Attempts         int       `bson:"attempts" json:"attempts"`
	MaxAttempts      int       `bson:"max_attempts" json:"max_attempts"`
	VerifiedAt       time.Time `bson:"verified_at" json:"verified_at"`
	CompletedAt      time.Time `bson:"completed_at" json:"completed_at"`
	AccessRestricted bool      `bson:"access_restricted" json:"access_restricted"`
	Restrictions     []string  `bson:"restrictions" json:"restrictions"`
}
