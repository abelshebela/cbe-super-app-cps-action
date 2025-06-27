package model

import (
	// "encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	ID             bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key            string        `json:"key" bson:"key"`
	ServiceName    string        `json:"service_name" bson:"service_name"`
	SingleCap      float32       `json:"single_cap" bson:"single_cap"`
	MinAmount      float32       `json:"min_amount" bson:"min_amount"`
	Flag           bool          `json:"flag" bson:"flag"`
	DailyCap       float32       `json:"daily_cap" bson:"daily_cap"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modifed_at" bson:"last_modifed_at"`
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
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type CurrentAction struct {
	Id     []string `bson:"id"`
	Action bool     `bson:"action"`
}

// type User struct {
// 	UserCode    string `json:"user_code" bson:"user_code"`
// 	FullName    string `json:"full_name" bson:"full_name"`
// 	PhoneNumber string `json:"phone_number" bson:"phone_number"`
// }

type RegistrationType string

const (
	RegistrationTypeNew    RegistrationType = "NEW"
	RegistrationTypeLinked RegistrationType = "LINKED"
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
type HQ struct {
	ID                   string          `bson:"_id" json:"id"`
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

type LinkedAccountResponse struct {
	ID                 string `json:"id"`
	AccountBranchType  string `json:"account_branchtype"`
	AccountBranchCode  string `json:"account_branchcode"`
	AccountNumber      string `json:"account_number"`
	CustomerNumber     string `json:"customer_number"`
	CustomerName       string `json:"customer_name"`
	AccountDescription string `json:"account_description"`
	PhoneNumber        string `json:"phone_number"`
	CustomerAddress    string `json:"customer_address"`
	DebitAllowed       bool   `json:"debit_allowed"`
	CreditAllowed      bool   `json:"credit_allowed"`
	AccountType        string `json:"account_type"`
	AccountFrozen      bool   `json:"account_frozen"`
	AccountDormant     bool   `json:"account_dormant"`
	ActiveAccount      bool   `json:"active_account"`
	AccountCurrency    string `json:"account_currency"`
}

type EnvironmentType string

const (
	UatEnvironment        EnvironmentType = "UAT"
	DevEnvironment        EnvironmentType = "DEV"
	TestEnvironment       EnvironmentType = "TEST"
	ProductionEnvironment EnvironmentType = "PRODUCTION"
)

type BranchType string

const (
	IFB BranchType = "IFB"
	CB  BranchType = "CB"
)

type ProductCode struct {
	ID          string     `bson:"id"`
	BranchType  BranchType `bson:"branch_type"`
	ProductCode string     `bson:"product_code"`
}

type AppType struct {
	UAT        string `bson:"uat"`
	Production string `bson:"production"`
	Test       string `bson:"test"`
	Dev        string `bson:"dev"`
}

type CredentialInformation struct {
	ID            string          `bson:"id"`
	Environment   EnvironmentType `bson:"environment"`
	MerchantAppID string          `bson:"merchant_app_id"`
	FabricAppID   string          `bson:"fabric_app_id"`
	ShortCode     string          `bson:"short_code"`
	AppSecret     string          `bson:"app_secret"`
	PrivateKey    string          `bson:"private_key"`
	PublicKey     string          `bson:"public_key"`
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}

///

type OTPFor string

const (
	OTPForLogin            OTPFor = "LOGIN"
	OTPForAddAccount       OTPFor = "ADD_ACCOUNT"
	OTPForPINSet           OTPFor = "PIN_SET"
	OTPForTransfer         OTPFor = "TRANSFER"
	OTPForAcctivateAccount OTPFor = "ACCTIVATE_ACCOUNT"
	OTPForPINReset         OTPFor = "PIN_RESET"
	OTPForSignup           OTPFor = "SIGNUP"
	OTPForAccountLink      OTPFor = "ACCOUNT_LINK"
	OTPForChangePhone      OTPFor = "CHANGE_PHONE"
	OTPForDetachPhone      OTPFor = "DETACH_PHONE"
	OTPForAttachPhone      OTPFor = "ATTACH_PHONE"
	OTPForEnable           OTPFor = "ENABLE"
	OTPForTransferLimit    OTPFor = "TRANSFER_LIMIT"
	OTPForChangeEmail      OTPFor = "CHANGE_EMAIL"
	OTPForUpgradeLimit     OTPFor = "UPGRADE_LIMIT"
)

type Realm string

const (
	ElstRealm     Realm = "ELST"
	BankRealm     Realm = "BANK"
	DistrictRealm Realm = "DISTRICT"
	BranchRealm   Realm = "BRANCH"
	MerchantRealm Realm = "MERCHANT"
	CompanyRealm  Realm = "COMPANY"
	MemberRealm   Realm = "MEMBER"
)

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusApproved KYCStatus = "APPROVED"
	KYCStatusRejected KYCStatus = "REJECTED"
)

type BPSStatus string

const (
	BPSStatusAuthorized BPSStatus = "AUTHORIZED"
	BPSStatusDenied     BPSStatus = "DENIED"
	BPSStatusPending    BPSStatus = "PENDING"
	BPSStatusInitiated  BPSStatus = "INITIATED"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type MaritalStatus string

const (
	Single   MaritalStatus = "SINGLE"
	Married  MaritalStatus = "MARRID"
	Divorced MaritalStatus = "DIVORCED"
	Widow    MaritalStatus = "WIDOW"
)

type PrimaryAuthentication string

const (
	PrimaryAuthenticationByPhoneNumber         PrimaryAuthentication = "PHONE_NUMBER"
	PrimaryAuthenticationByEmail               PrimaryAuthentication = "EMAIL"
	PrimaryAuthenticationByEmailAndPhoneNumber PrimaryAuthentication = "EMAIL_AND_PHONE_NUMBER"
)

type DeviceStatus string

const (
	Linked   DeviceStatus = "LINKED"
	UnLinked DeviceStatus = "UNLINKED"
)

type DeviceType string

const (
	Android DeviceType = "Android"
	IOS     DeviceType = "IOS"
)

type Device struct {
	Type    DeviceType
	Version string
}
type OTPStatus string

const (
	Pending  OTPStatus = "PENDING"
	Verified OTPStatus = "VERIFIED"
	Denied   OTPStatus = "DENIED"
)

type MerchantRole string

const (
	MerchantRoleOwner MerchantRole = "OWNER"
	MerchantRoleAgent MerchantRole = "AGENT"
)

type PoolSource string

const (
	PoolSourcePortal PoolSource = "PORTAL"
	PoolSourceApp    PoolSource = "APP"
	PoolSourceAgent  PoolSource = "AGENT"
)

type MemberType string

type AccountStatus string

const (
	Active   AccountStatus = "ACTIVE"
	InActive AccountStatus = "IN_ACTIVE"
)

type User struct {
	ID                bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode          string        `json:"user_code" bson:"user_code"`
	FullName          string        `json:"full_name" bson:"full_name"`
	MotherName        string        `json:"mother_name" bson:"mother_name"`
	Nationality       string        `json:"nationality" bson:"nationality"`
	BirthDate         time.Time     `json:"birth_date" bson:"birth_date"`
	BranchName        string        `json:"branch_name" bson:"branch_name"`
	DistrictName      string        `json:"district_name" bson:"district_name"`
	BranchCode        string        `json:"branch_code" bson:"branch_code"`
	DistrictCode      string        `json:"district_code" bson:"district_code"`
	ResidentialStatus string        `json:"residential_status" bson:"residential_status"`
	IssuedDate        time.Time     `json:"issued_date" bson:"isssued_date"`
	PhoneNumber       string        `json:"phone_number" bson:"phone_number"`
	Gender            Gender        `json:"gender" bson:"gender"`

	Fayda struct {
		FaydaID          string `json:"id_number" bson:"id_number"`
		FaydaAccessToken string `json:"fayda_access_token" bson:"fayda_access_token"`
		EmploymentStatus string `json:"employment_status" bson:"employement_status"`
		EmployerName     string `json:"employer_name" bson:"employer_name"`
		IssuedBy         string `json:"issued_by" bson:"issued_by"`
		MonthlyIncome    uint64 `json:"monthly_incode" bson:"monthly_incode"`
	} `json:"fayda" bson:"fayda"`
	Address           Address          `json:"address" bson:"address"`
	AndOrCustomerName []string         `json:"and_or_customer_name" bson:"and_or_customer_name"`
	DocumentFront     string           `json:"document_front" bson:"document_front"`
	DocumentBack      string           `json:"document_back" bson:"document_back"`
	Photo             string           `json:"photo" bson:"photo"`
	Signature         string           `json:"signature" bson:"signature"`
	Avatar            string           `json:"avater" bson:"avater"`
	Email             string           `json:"email" bson:"email"`
	Username          string           `json:"username" bson:"username"`
	AndOrStatus       bool             `json:"and_or_status" bson:"and_or_status"`
	PushToken         string           `json:"push_token" bson:"push_token"`
	Realm             Realm            `json:"realm" bson:"realm"`
	PermissionGroup   []bson.ObjectID  `json:"permission_group" bson:"permission_group"`
	Permissions       []bson.ObjectID  `json:"permissions" bson:"persmissions"`
	IsAccountBlocked  bool             `json:"is_account_blocked" bson:"is_account_blocked"`
	MainAccount       string           `json:"main_account" bson:"main_account"`
	LastMainAccount   string           `json:"last_main_account" bson:"last_main_account"`
	AccountLinked     bool             `json:"account_linked" bson:"account_linked"`
	LastAccountLinked bool             `json:"last_account_linked" bson:"last_account_linked"`
	MemberType        MemberType       `json:"account_branch_type" bson:"account_branch_type"`
	RegistrationType  RegistrationType `json:"account_type" bson:"account_type"`
	AccountStatus     AccountStatus    `json:"account_status" bson:"account_status"`
	KYC               struct {
		KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
		KYCStatus            KYCStatus           `json:"kyc_status" bson:"kyc_status"`
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
	NextLoginAttempt  time.Time `json:"next_attempt_count" bson:"next_attempt_count"`
	LastLoginAttempt  time.Time `json:"last_login_attempt" bson:"last_login_attempt"`
	LastOnlineDate    time.Time `json:"last_online_date" bson:"last_online_date"`
	LastLogin         time.Time `json:"last_login" bson:"last_login"`

	// move BPS Status and realted field to new collection
	BPSStatus          BPSStatus `json:"bps_reject_status" bson:"bps_reject_status"`
	BPSRejectionReason string    `json:"bps_reject_reason" bson:"bps_reject_reason"`
	BPSRejectionField  []string  `json:"bps_reject_failed" bson:"bps_reject_failed"`
	LoginPIN           LoginPIN  `json:"login_pin" bson:"login_pin"`
	Device             struct {
		DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
		AppVersion string `json:"app_version" bson:"app_version"`
	} `json:"device" bson:"device"`
	APPInstallationDate   time.Time             `json:"application_installation_date" bson:"application_installation_date"`
	CustomerNumber        string                `json:"customer_number" bson:"customer_number"`
	InitialLinkedDate     time.Time             `json:"initial_linked_date" bson:"initiali_linked_date"`
	PrimaryAuthentication PrimaryAuthentication `json:"primary_authentication" bson:"primary_authentication"`
	LoanScore             uint16                `json:"loan_score" bson:"loan_score"`
	DeviceStatus          DeviceStatus          `json:"device_status" bson:"device_status"`
	SessionExpirsOn       time.Time             `json:"session_expires_on" bson:"session_expires_on"`
	Enabled               bool                  `json:"enabled" bson:"enabled"`
	FirstPinSet           bool                  `json:"first_pin_set" bson:"first_pin_set"`
	IsDeleted             bool                  `json:"is_deleted" bson:"is_deleted"`
	PINChangedAt          time.Time             `json:"pin_changed_at" bson:"pin_changed_at"`
	OTPLastTriedAt        time.Time             `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt     time.Time             `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	OTPVerifyCount        uint8                 `json:"otp_verify_count" bson:"otp_verify_count"`
	PINHistory            []string              `json:"pin_history" bson:"pin_history"`
	InitialiLinkedAt      time.Time             `json:"initial_linked_at" bson:"initial_linked_at"`
	CreatedAt             time.Time             `json:"created_at" bson:"created_at"`
	DeletedAt             time.Time             `json:"delete_at" bson:"deleted_at"`
	LastModifiedAt        time.Time             `json:"last_modified_at" bson:"last_modified_at"`
	LinkedAccount         []LinkedAccount       `json:"linked_account" bson:"linked_account"`
	OrganizationID        []LinkedAccount       `json:"organization_id" bson:"organization_id"`
}

// Optionally Redis
type OTP struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber   string        `json:"phone_number" bson:"phone_number"`
	AccountNumber *string       `json:"account_number" bson:"account_number"`
	UserRealm     Realm         `json:"user_realm" bson:"user_realm"`
	Email         string        `json:"email" bson:"email"`
	UserCode      string        `json:"user_code" bson:"user_code"`
	OTPCode       string        `json:"otp_code" bson:"otp_code"`
	BillNo        *string       `json:"bill_no,omitempty" bson:"bill_no,omitempty"`
	DeviceUUID    *string       `json:"device_uuid,omitempty" bson:"device_uuid,omitempty"`
	OTPFor        OTPFor        `json:"otp_for" bson:"otp_for"`
	Status        OTPStatus     `json:"status" bson:"status"`
	ExpiresAt     time.Time     `json:"expires_at" bson:"expires_at"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	LastModified  time.Time     `json:"last_modified" bson:"last_modified"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time     `json:"deleted_at" bson:"deleted_at,omitempty"`
}

type DeviceLinkHistroy struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     bson.ObjectID `json:"user_id" bson:"user_id"`
	Device     Device        `json:"device" bson:"device"`
	LinkedAt   time.Time     `json:"linked_at" bson:"linked_at"`
	UnLinkedAt time.Time     `json:"unlinked_at" bson:"unlinked_at"`
}

type FullName struct {
	FirstName  string `json:"first_name" bson:"first_name"`
	MiddleName string `json:"middle_name" bson:"middle_name"`
	LastName   string `json:"last_name" bson:"last_name"`
}

type PhoneNumber struct {
	Code       string `json:"code" bson:"code"`
	Number     string `json:"number" bson:"number"`
	IsVerified *bool  `json:"is_verfied,omitempty" bson:"is_verfied,omitempty"`
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
type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}
