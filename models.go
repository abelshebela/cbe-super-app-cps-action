package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// =============================================================================
// CONSTANTS
// =============================================================================

// User Types
const (
	Maker   = "MAKER"
	Checker = "CHECKER"
)

// Action Statuses
const (
	ActionPending  = "PENDING"
	ActionApproved = "APPROVED"
	ActionRejected = "REJECTED"
)

// Action Types
const (
	ActionCreate  = "CREATE"
	ActionUpdate  = "UPDATE"
	ActionDelete  = "DELETE"
	ActionEnable  = "ENABLE"
	ActionDisable = "DISABLE"
)

// Gender Types
const (
	Male   = "MALE"
	Female = "FEMALE"
)

// Marital Status
const (
	Single   = "SINGLE"
	Married  = "MARRID"
	Divorced = "DIVORCED"
	Widow    = "WIDOW"
)

// Device Status
const (
	Linked   = "LINKED"
	UnLinked = "UNLINKED"
)

// OTP Status
const (
	OTPPending  = "PENDING"
	OTPVerified = "VERIFIED"
	OTPDenied   = "DENIED"
)

// OTP For
const (
	OTPForLogin            = "LOGIN"
	OTPForAddAccount       = "ADD_ACCOUNT"
	OTPForPINSet           = "PIN_SET"
	OTPForTransfer         = "TRANSFER"
	OTPForAcctivateAccount = "ACCTIVATE_ACCOUNT"
	OTPForPINReset         = "PIN_RESET"
	OTPForSignup           = "SIGNUP"
	OTPForAccountLink      = "ACCOUNT_LINK"
	OTPForChangePhone      = "CHANGE_PHONE"
	OTPForDetachPhone      = "DETACH_PHONE"
	OTPForAttachPhone      = "ATTACH_PHONE"
	OTPForEnable           = "ENABLE"
	OTPForTransferLimit    = "TRANSFER_LIMIT"
	OTPForChangeEmail      = "CHANGE_EMAIL"
	OTPForUpgradeLimit     = "UPGRADE_LIMIT"
)

// Realm Types
const (
	ElstRealm     = "ELST"
	BankRealm     = "BANK"
	DistrictRealm = "DISTRICT"
	BranchRealm   = "BRANCH"
	MerchantRealm = "MERCHANT"
	CompanyRealm  = "COMPANY"
	MemberRealm   = "MEMBER"
)

// Advert For
const (
	Both = "ALL"
)

// Event Status
const (
	EventUpcomming = "UPCOMMING"
	EventLive      = "LIVE"
	EventClosed    = "CLOSED"
)

// Restriction Type
const (
	AgeRestriction = "AGE_RESTRICTION"
)

// Ticket Status
const (
	TicketBooked   = "BOOKED"
	TicketExpired  = "EXPIRED"
	TicketPaid     = "PAID"
	TicketRedeemed = "REDEEMED"
)

// Notification Status
const (
	StatusPending = "PENDING"
	StatusSent    = "SENT"
	StatusSeen    = "SEEN"
)

// Notification For
const (
	ForIFB = "IFB"
	ForCB  = "CB"
	ForAll = "ALL"
)

// Environment Type
const (
	UatEnvironment        = "UAT"
	DevEnvironment        = "DEV"
	TestEnvironment       = "TEST"
	ProductionEnvironment = "PRODUCTION"
)

// Branch Type
const (
	IFB = "IFB"
	CB  = "CB"
)

// App Type
const (
	URL = "URL"
)

// Stage
const (
	StageUat = "UAT"
)

// App View Type
const (
	AppViewTypeBoth = "BOTH"
	AppViewTypeCB   = "CB"
	AppViewTypeIFB  = "IFB"
)

// Method
const (
	OPEN      = "OPEN"
	PIN       = "PIN"
	OTPANDPIN = "OTP_PIN"
)

// =============================================================================
// TYPE DEFINITIONS
// =============================================================================

// UserType represents the type of user
type UserType string

// ActionStatus represents the status of an action
type ActionStatus string

// ActionType represents the type of action
type ActionType string

// Gender represents user gender
type Gender string

// MaritalStatus represents user marital status
type MaritalStatus string

// DeviceStatus represents device linking status
type DeviceStatus string

// OTPStatus represents OTP status
type OTPStatus string

// OTPFor represents what the OTP is for
type OTPFor string

// Realm represents user realm
type Realm string

// AdvertFor represents advertisement target
type AdvertFor string

// EventStatus represents event status
type EventStatus string

// RestrictionType represents restriction type
type RestrictionType string

// TicketStatus represents ticket status
type TicketStatus string

// NotificationStatus represents notification status
type NotificationStatus string

// NotificationFor represents notification target
type NotificationFor string

// EnvironmentType represents environment type
type EnvironmentType string

// BranchType represents branch type
type BranchType string

// AppType represents app type
type AppType string

// Stage represents app stage
type Stage string

// AppViewType represents app view type
type AppViewType string

// Method represents authentication method
type Method string

type RequestAction string

const (
	REQUEST_USER                       RequestAction = "USER"
	REQUEST_PERMISSION_GROUP           RequestAction = "PERMISSION_GROUP"
	REQUEST_DEPARTMENT                 RequestAction = "DEPARMTENT"
	REQUEST_ENABLE_USER                RequestAction = "ENABLE_USER"
	REQUEST_DISABLE_USER               RequestAction = "DISABLE_USER"
	REQUEST_BPS_USER                   RequestAction = "BPS_USER"
	REQUEST_DISABLE_BPS_USER           RequestAction = "DISABLE_BPS_USER"
	REQUEST_ENABLE_BPS_USER            RequestAction = "ENABLE_BPS_USER"
	REQUEST_UPDATE_USER                RequestAction = "UPDATE_USER"
	REQUEST_TOTAL_DAILY_LIMIT          RequestAction = "TOTAL_DAILY_LIMIT"
	REQUEST_UPDATE_VAT                 RequestAction = "UPDATE_VAT"
	REQUEST_AUTH_TIER                  RequestAction = "AUTHTIER"
	REQUEST_CREATE_ADVERT              RequestAction = "CREATE_ADVERT"
	REQUEST_UPDATE_ADVERT              RequestAction = "UPDATE_ADVERT"
	REQUEST_ENABLE_ADVERT              RequestAction = "ENABLE_ADVERT"
	REQUEST_DISABLE_ADVERT             RequestAction = "DISABLE_ADVERT"
	REQUEST_DELETE_ADVERT              RequestAction = "DELETE_ADVERT"
	REQUEST_CREATE_BANK                RequestAction = "CREATE_BANK"
	REQUEST_UPDATE_BANK                RequestAction = "UPDATE_BANK"
	REQUEST_ENABLE_BANK                RequestAction = "ENABLE_BANK"
	REQUEST_DISABLE_BANK               RequestAction = "DISABLE_BANK"
	REQUEST_ENABLE_WALLET              RequestAction = "ENABLE_WALLET"
	REQUEST_DISABLE_WALLET             RequestAction = "DISABLE_WALLET"
	REQUEST_UPDATE_PASSWORD_EXPIRY     RequestAction = "UPDATE_PASSWORD_EXPIRY"
	REQUEST_CREATE_VALIDATION          RequestAction = "CREATE_VALIDATION"
	REQUEST_UPDATE_VALIDATION          RequestAction = "UPDATE_VALIDATION"
	REQUEST_DELETE_VALIDATION          RequestAction = "DELETE_VALIDATION"
	REQUEST_UPDATE_ARCHIVE_EXPIRY      RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	REQUEST_CREATE_SERVICE_FEE         RequestAction = "CREATE_SERVICE_FEE"
	REQUEST_UPDATE_SERVICE_FEE         RequestAction = "UPDATE_SERVICE_FEE"
	REQUEST_DELETE_SERVICE_FEE         RequestAction = "DELETE_SERVICE_FEE"
	REQUEST_CREATE_DAILY_LIMIT         RequestAction = "CREATE_DAILY_LIMIT"
	REQUEST_UPDATE_DAILY_LIMIT         RequestAction = "UPDATE_DAILY_LIMIT"
	REQUEST_DELETE_DAILY_LIMIT         RequestAction = "DELETE_DAILY_LIMIT"
	REQUEST_BUDGET_COLOR               RequestAction = "BUDGET_COLOR"
	REQUEST_BUDGET_ICON                RequestAction = "BUDGET_ICON"
	REQUEST_UPDATE_PRODUCT             RequestAction = "UPDATE_PRODUCT"
	REQUEST_CREATE_PUBLIC_NOTIFICATION RequestAction = "CREATE_PUBLIC_NOTIFICATION"
	REQUEST_ARCHIVE_USER               RequestAction = "ARCHIVE_USER"
	REQUEST_CREATE_PASSWORD_RULE       RequestAction = "CREATE_PASSWORD_RULE"
	REQUEST_UPDATE_PASSWORD_RULE       RequestAction = "UPDATE_PASSWORD_RULE"
	REQUEST_UPDATE_MINIMUM_SERVICE     RequestAction = "UPDATE_MINIMUM_SERVICE"
	REQUEST_UPDATE_SERVICE_RULE        RequestAction = "UPDATE_SERVICE_RULE"
	REQUEST_UPDATE_TOTAL               RequestAction = "UPDATE_TOTAL"
	REQUEST_UPDATE_ACCESS_CONFIG       RequestAction = "UPDATE_ACCESS_CONFIG"
	REQUEST_ENABLE_SINGLE_BRANCH       RequestAction = "ENABLE_SINGLE_BRANCH"
	REQUEST_DISABLE_SINGLE_BRANCH      RequestAction = "DISABLE_SINGLE_BRANCH"
	REQUEST_ENABLE_MULTI_USERS         RequestAction = "ENABLE_MULTI_USERS"
	REQUEST_DISABLE_MULTI_USERS        RequestAction = "DISABLE_MULTI_USERS"
	REQUEST_CREATE_BUSINESS            RequestAction = "CREATE_BUSINESS"
	REQUEST_UPDATE_BUSINESS            RequestAction = "UPDATE_BUSINESS"
	REQUEST_CREATE_EVENT               RequestAction = "CREATE_EVENT"
	REQUEST_UPDATE_EVENT               RequestAction = "UPDATE_EVENT"
	REQUEST_CREATE_EVENT_CATEGORY      RequestAction = "CREATE_EVENT_CATEGORY"
	REQUEST_UPDATE_EVENT_CATEGORY      RequestAction = "UPDATE_EVENT_CATEGORY"
	REQUEST_DISABLE_EVENT              RequestAction = "DISABLE_EVENT"
	REQUEST_CREATE_MINIAPP_MERCHANT    RequestAction = "CREATE_MINIAPP_MERCHANT"
	REQUEST_UPDATE_MINIAPP_MERCHANT    RequestAction = "UPDATE_MINIAPP_MERCHANT"
	REQUEST_UPDATE_BLOCK_TIME          RequestAction = "UPDATE_BLOCK_TIME"
	REQUEST_UPDATE_ACCOUNT_VALIDATION  RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	REQUEST_UPDATE_SERVICE_DETAILS     RequestAction = "UPDATE_SERVICE_DETAILS"

	// Add these for budget icon and color actions
	REQUEST_CREATE_ICON  RequestAction = "CREATE_ICON"
	REQUEST_UPDATE_ICON  RequestAction = "UPDATE_ICON"
	REQUEST_DELETE_ICON  RequestAction = "DELETE_ICON"
	REQUEST_CREATE_COLOR RequestAction = "CREATE_COLOR"
	REQUEST_UPDATE_COLOR RequestAction = "UPDATE_COLOR"
	REQUEST_DELETE_COLOR RequestAction = "DELETE_COLOR"
)

// =============================================================================
// CORE MODELS
// =============================================================================

type ActionData struct {
	UseCode     string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

// User represents a basic user
type User struct {
	UserID      string    `json:"user_id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Timestamp   time.Time `json:"timestamp"`
}

// BPSUser represents a bank personnel user
type BPSUser struct {
	ID                bson.ObjectID `json:"id" bson:"_id,omitempty"`
	UserCode          string        `json:"user_code" bson:"user_code"`
	FullName          string        `json:"full_name" bson:"full_name"`
	UserName          string        `json:"UserName" bson:"UserName"`
	PhoneNumber       string        `json:"phone_number" bson:"phone_number"`
	BranchCode        []string      `json:"branch_code" bson:"branch_code"`
	BranchName        string        `json:"branch_name" bson:"branch_name"`
	HomeBranch        string        `json:"home_branch" bson:"home_branch"`
	Role              string        `json:"role" bson:"role"`
	Realm             Realm         `json:"realm" bson:"realm"`
	LoginAttemptCount uint8         `json:"login_attempt_count" bson:"login_attempt_count"`
	Password          Password      `json:"password" bson:"password"`
	FirstPasswordSet  bool          `json:"first_password_set" bson:"first_password_set"`
	Enabled           bool          `json:"enabled" bson:"enabled"`
	IsDeleted         bool          `json:"is_deleted" bson:"is_deleted"`
	OTPVerifyCount    uint8         `json:"otp_verfy_count" bson:"otp_verify_count"`
	OTPLastTriedAt    time.Time     `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	OTPLastVerifiedAt time.Time     `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	LastLoginAttempt  time.Time     `json:"last_login_attempt" bson:"last_login_attempt"`
	NextLoginAttempt  time.Time     `json:"next_login_attempt" bson:"next_login_attempt"`
	LastLogin         time.Time     `json:"last_login" bson:"last_login"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt    time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}

// CPSUser represents a CPS system user
type CPSUser struct {
	ID                 bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode           string          `json:"user_code,omitempty" bson:"user_code,omitempty"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Role               string          `json:"role,omitempty" bson:"role,omitempty"`
	Department         bson.ObjectID   `json:"department,omitempty" bson:"department,omitempty"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty"`
	UserName           string          `json:"username,omitempty" bson:"username,omitempty"`
	Realm              Realm           `json:"realm,omitempty" bson:"realm,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category,omitempty"`
	PermissionGroup    []bson.ObjectID `json:"permission_group,omitempty" bson:"permission_group,omitempty"`
	Password           Password        `json:"password" bson:"password"`
	PasswordDisable    bool            `json:"password_disable,omitempty" bson:"password_disable,omitempty"`
	SyncDisabled       bool            `json:"sync_disabled,omitempty" bson:"sync_disabled,omitempty"`
	LoginAttemptCount  uint8           `json:"login_attempt_count,omitempty" bson:"login_attempt_count,omitempty"`
	LastLoginAttempt   time.Time       `json:"last_login_attempt,omitempty" bson:"last_login_attempt,omitempty"`
	NextLoginAttempt   time.Time       `json:"next_login_attempt,omitempty" bson:"next_login_attempt,omitempty"`
	LastOnlineDate     time.Time       `json:"last_online_date,omitempty" bson:"last_online_date,omitempty"`
	Enabled            bool            `json:"enabled,omitempty" bson:"enabled,omitempty"`
	IsDeleted          bool            `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	CreatedAt          time.Time       `json:"created_at,omitempty" bson:"created_at,omitempty"`
	LastModifiedAt     time.Time       `json:"last_modified_at,omitempty" bson:"last_modified_at,omitempty"`
}

// Password represents password information
type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password"`
	PasswordChangeAt time.Time `json:"password_change_at" bson:"password_change_at"`
}

// OTP represents one-time password
type OTP struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber   string    `json:"phone_number" bson:"phone_number"`
	AccountNumber *string   `json:"account_number" bson:"account_number"`
	Realm         Realm     `json:"user_realm" bson:"user_realm"`
	Email         string    `json:"email" bson:"email"`
	UserCode      string    `json:"user_code" bson:"user_code"`
	OTPCode       string    `json:"otp_code" bson:"otp_code"`
	BillNo        *string   `json:"bill_no,omitempty" bson:"bill_no,omitempty"`
	DeviceUUID    *string   `json:"device_uuid,omitempty" bson:"device_uuid,omitempty"`
	OTPFor        OTPFor    `json:"otp_for" bson:"otp_for"`
	Status        OTPStatus `json:"status" bson:"status"`
	IsDeleted     bool      `json:"is_deleted" bson:"is_deleted"`
	ExpiresAt     time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	LastModified  time.Time `json:"last_modified" bson:"last_modified"`
	DeletedAt     time.Time `json:"deleted_at" bson:"deleted_at,omitempty"`
}

// =============================================================================
// BUSINESS ENTITY MODELS
// =============================================================================

// Bank represents bank information
type Bank struct {
	ID             string    `json:"id,omitempty" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	Logo           string    `json:"logo" bson:"logo"`
	Code           string    `json:"code" bson:"code"`
	BIC            string    `json:"bic" bson:"bic"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at,omitzero" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at,omitzero" bson:"last_modified_at"`
}

// Branch represents branch information
type Branch struct {
	ID            string `json:"id" bson:"id"`
	Code          string `json:"code" bson:"code"`
	Name          string `json:"name" bson:"name"`
	Address       string `json:"address" bson:"address"`
	Owner         string `json:"owner" bson:"owner"`
	AccountNumber string `json:"account_number" bson:"account_number"`
}

// HQ represents headquarters information
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

// LinkedAccount represents linked account information
type LinkedAccount struct {
	AccountNumber string `json:"account_number" bson:"account_number"`
	AccountType   string `json:"account_type" bson:"account_type"`
}

// =============================================================================
// PERMISSION & ACCESS MODELS
// =============================================================================

// Permission represents individual permission
type Permission struct {
	PermissionName string `bson:"permissionName" json:"permissionName"`
}

// PermissionGroup represents permission groupings
type PermissionGroup struct {
	ID                 bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GroupName          string        `bson:"group_name,omitempty" json:"group_name,omitempty"`
	Permissions        interface{}   `bson:"permissions,omitempty" json:"permissions,omitempty"`
	PermissionCategory interface{}   `bson:"permission_category,omitempty" json:"permission_category,omitempty"`
	Role               string        `bson:"role,omitempty" json:"role,omitempty"`
	Realm              string        `bson:"realm,omitempty" json:"realm,omitempty"`
	Enabled            bool          `bson:"enabled,omitempty" json:"enabled,omitempty"`
	IsDeleted          bool          `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
	CreatedAt          time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LastModified       time.Time     `bson:"last_modified,omitempty" json:"last_modified,omitempty"`
}

// PermissionCategory represents permission categories
type PermissionCategory struct {
	ID           bson.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	CategoryName string          `bson:"categoryName" json:"categoryName"`
	Access       string          `bson:"access" json:"access"`
	Permissions  []bson.ObjectID `bson:"permissions" json:"permissions"`
	Enabled      bool            `bson:"enabled,omitempty" json:"enabled,omitempty"`
	IsDeleted    bool            `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	CreatedAt    time.Time       `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt    time.Time       `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// PortalCard represents portal access cards
type Card struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CardName string        `bson:"card_name" json:"card_name"`
	SubCards []string      `bson:"sub_cards" json:"sub_cards"`
}

// =============================================================================
// SERVICE & PRODUCT MODELS
// =============================================================================

// Service represents service definitions
type Service struct {
	ID                 bson.ObjectID `json:"_id" bson:"_id"`
	ServiceCode        string        `json:"service_code" bson:"service_code"`
	ServiceName        string        `json:"service_name" bson:"service_name"`
	ServiceType        string        `json:"service_type" bson:"service_type"`
	Key                string        `json:"key" bson:"key"`
	Cap                Cap           `json:"cap" bson:"cap"`
	CBEProductCodes    ProductCodes  `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes  `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64        `json:"above_amount" bson:"above_amount"`
	AboveServiceFee    uint64        `json:"above_service_fee" bson:"above_service_fee"`
	PaymentType        string        `json:"payment_type" bson:"payment_type"`
	Tiers              []Tier        `json:"tiers" bson:"tiers"`
	CBEGLEntry         GLEntry       `json:"cbe_gl_entry" bson:"cbe_gl_entry"`
	CBEIFBGLEntry      GLEntry       `json:"cbe_ifb_gl_entry" bson:"cbe_ifb_gl_entry"`
	Enabled            bool          `json:"enabled" bson:"enabled"`
	IsDeleted          bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time     `json:"deleted_at" bson:"deleted_at"`
}

// ProductCode represents product codes
type ProductCode struct {
	ID             string     `bson:"id" json:"id,omitempty"`
	BranchType     BranchType `bson:"branch_type" json:"branch_type"`
	ProductCode    string     `bson:"product_code" json:"product_code"`
	VATCode        string     `bson:"vat_code" json:"vat_code"`
	ServiceFeeCode string     `bson:"service_fee_code" json:"service_fee_code"`
}

// ProductCodes represents product code collections
type ProductCodes struct {
	PRD    string `bson:"prd"`
	VATPRD string `bson:"vatprd"`
	SFPRD  string `bson:"sfprd"`
	TRXN   string `bson:"trxn"`
}

// MiniApp represents mini application
type MiniApp struct {
	ID                  string                `bson:"_id" json:"id,omitempty"`
	AppName             string                `bson:"app_name" json:"app_name"`
	AppIcon             string                `bson:"app_icon" json:"app_icon"`
	BannerImage         string                `bson:"banner_image" json:"banner_image"`
	CommissionGLAccount string                `bson:"commison_gl_account" json:"commison_gl_account,omitempty"`
	AppType             AppType               `bson:"app_type" json:"app_type"`
	MerchantID          string                `bson:"merchant_id" json:"merchant_id"`
	ProductCode         []ProductCode         `bson:"product_code" json:"product_code"`
	Credential          CredentialInformation `bson:"credential" json:"credential,omitempty"`
	URL                 string                `bson:"url" json:"url,omitempty"`
	AppViewType         AppViewType           `bson:"app_view_type" json:"app_view_type"`
	Stage               Stage                 `bson:"stage" json:"stage"`
	IsEventMiniApp      bool                  `bson:"is_event_mini_app" json:"is_event_mini_app"`
	IsThreeClick        bool                  `bson:"is_three_click" json:"is_three_click"`
	Enabled             bool                  `bson:"enabled" json:"enabled"`
	IsDeleted           bool                  `bson:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time             `bson:"created_at" json:"created_at"`
	LastModifiedAt      time.Time             `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt           time.Time             `bson:"deleted_at" json:"deleted_at"`
}

// CredentialInformation represents app credentials
type CredentialInformation struct {
	ID            string          `bson:"id" json:"id,omitempty"`
	Environment   EnvironmentType `bson:"environment" json:"environment"`
	MerchantAppID string          `bson:"merchant_app_id" json:"merchant_app_id"`
	FabricAppID   string          `bson:"fabric_app_id" json:"fabric_app_id"`
	ShortCode     string          `bson:"short_code" json:"short_code"`
	MiniAppCode   string          `bson:"mini_app_code" json:"mini_app_code"`
	AppSecret     string          `bson:"app_secret" json:"app_secret"`
	PrivateKey    string          `bson:"private_key" json:"-"`
	PublicKey     string          `bson:"public_key" json:"public_key"`
	Signature     string          `bson:"signature" json:"-"`
	Timestamp     time.Time       `bson:"timestamp" json:"-"`
}

// =============================================================================
// FINANCIAL MODELS
// =============================================================================

// Wallet represents wallet management
type Wallet struct {
	ID             string     `json:"id,omitempty" bson:"id"`
	Name           string     `json:"name" bson:"name"`
	Code           string     `json:"code" bson:"code"`
	Avatar         string     `json:"avatar" bson:"avatar"`
	Enabled        bool       `json:"enabled" bson:"enabled"`
	IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deletd_at,omitempty" bson:"deleted_at"`
}

// BudgetCategory represents budget categories
type BudgetCategory struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string        `bson:"name,omitempty" json:"name,omitempty"`
	Icon          string        `bson:"icon,omitempty" json:"icon,omitempty"`
	Description   string        `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt     time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LastUpdatedAt time.Time     `bson:"last_updated_at,omitempty" json:"last_updated_at,omitempty"`
	IsDeleted     bool          `bson:"is_deleted" json:"is_deleted,omitempty"`
}

// AmountBasedAuth represents amount-based authentication
type AmountBasedAuth struct {
	ID           bson.ObjectID `json:"id" bson:"_id"`
	MinAmount    uint64        `json:"min_amount" bson:"min_amount"`
	MaxAmount    uint64        `json:"max_amount" bson:"max_amount"`
	Method       Method        `json:"method" bson:"method"`
	Enabled      bool          `json:"enabled" bson:"enabled"`
	IsDeleted    bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	LastModified time.Time     `json:"last_modified" bson:"last_modified"`
}

// AuthTier represents authentication tiers
type AuthTier struct {
	ID        *string `json:"id" bson:"id"`
	Min       uint64  `json:"min" bson:"min"`
	Max       uint64  `json:"max" bson:"max"`
	FeeAmount uint64  `json:"fee_amount" bson:"fee_amount"`
}

// Cap represents capacity limits
type Cap struct {
	KYCLevel  string `json:"kyc_level" bson:"kyc_level"`
	SingleCap uint64 `json:"single_cap" bson:"single_cap"`
	DailyCap  uint64 `json:"daily_cap" bson:"daily_cap"`
	MinAmount uint64 `json:"min_amount" bson:"min_amount"`
}

// Tier represents fee tiers
type Tier struct {
	ID        *string `json:"id" bson:"id"`
	Min       uint64  `json:"min" bson:"min"`
	Max       uint64  `json:"max" bson:"max"`
	FeeAmount uint64  `json:"fee_amount" bson:"fee_amount"`
}

// GLEntry represents general ledger entries
type GLEntry struct {
	ProductAccount    string `json:"product_account" bson:"product_account"`
	ProductBranchCode string `json:"product_branch_code" bson:"product_branch_code"`
	ServiceAccount    string `json:"service_account" bson:"service_account"`
	ServiceBranchCode string `json:"service_branch_code" bson:"service_branch_code"`
	VatAccount        string `json:"vat_account" bson:"vat_account"`
	VatBranchCode     string `json:"vat_branch_code" bson:"vat_branch_code"`
}

// =============================================================================
// CONTENT & MEDIA MODELS
// =============================================================================

// Advert represents advertisement
type Advert struct {
	ID            string     `json:"_id,omitempty" bson:"_id"`
	Title         string     `json:"title" bson:"title"`
	Description   string     `json:"description" bson:"description"`
	BannerImage   string     `json:"banner_image" bson:"banner_image"`
	AdvertFor     AdvertFor  `json:"advert_for" bson:"advert_for"`
	Date          AdvertDate `json:"advert_date" bson:"advert_date"`
	Enabled       bool       `json:"enabled" bson:"enabled"`
	IsDeleted     bool       `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time  `json:"deleted_at" bson:"deleted_at"`
	CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at" bson:"last_updated_at"`
}

// AdvertDate represents advertisement date range
type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

// Avatar represents user avatars
type Avatar struct {
	ID             string     `json:"id,omitempty" bson:"_id"`
	Avatar         string     `json:"avatar" bson:"avatar"`
	Label          string     `json:"label" bson:"label"`
	Enable         bool       `json:"enable" bson:"enable"`
	IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// Event represents event management
type Event struct {
	ID                  string              `json:"id,omitempty" bson:"_id,omitempty"`
	EventCode           string              `json:"event_code" bson:"event_code"`
	EventName           string              `json:"event_name" bson:"event_name"`
	EventCity           string              `json:"event_city" bson:"event_city"`
	AccountNumber       string              `json:"account_number" bson:"account_number"`
	EventVenue          string              `json:"event_venue" bson:"event_venue"`
	RefundPolicy        []string            `json:"refund_policy" bson:"refund_policy"`
	MICSInfo            []string            `json:"mics_info" bson:"mics_info"`
	Restriction         Restriction         `json:"restriction" bson:"restriction"`
	Ticket              []Ticket            `json:"ticket" bson:"ticket"`
	Status              EventStatus         `json:"status" bson:"status"`
	EventInformation    EventInformation    `json:"event_information" bson:"event_information"`
	TicketStatistics    TicketStatistics    `json:"ticket_statistics" bson:"ticket_statistics"`
	TicketInformation   TicketInformation   `json:"ticket_information" bson:"ticket_information"`
	MerchantInformation MerchantInformation `json:"merchant_information" bson:"merchant_information"`
	Enabled             bool                `json:"enabled" bson:"enabled"`
	IsDeleted           bool                `json:"is_deleted" bson:"is_deleted"`
	HasRestriction      bool                `json:"has_restriction" bson:"has_restriction"`
	CreatedAt           time.Time           `json:"created_at" bson:"created_at"`
	DeletedAt           time.Time           `json:"deleted_at" bson:"deleted_at"`
	LastModifiedAt      time.Time           `json:"last_modified_at" bson:"last_modified_at"`
}

// EventInformation represents event details
type EventInformation struct {
	StartDate   time.Time `json:"start_date" bson:"start_date"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Description string    `json:"description" bson:"description"`
	Cover       string    `json:"cover" bson:"cover"`
	VideoLink   string    `json:"video_link" bson:"video_link"`
}

// Restriction represents event restrictions
type Restriction struct {
	Type        RestrictionType `json:"type" bson:"type"`
	Description string          `json:"description" bson:"description"`
}

// Ticket represents event tickets
type Ticket struct {
	Name           string `json:"name" bson:"name"`
	Category       string `json:"category" bson:"category"`
	Type           string `json:"type" bson:"type"`
	Price          uint64 `json:"price" bson:"price"`
	NumberOfTicker uint8  `json:"number_of_ticker" bson:"number_of_ticker"`
}

// TicketStatistics represents ticket statistics
type TicketStatistics struct {
	Category           string `json:"category" bson:"category"`
	Revenue            uint64 `json:"revenue" bson:"revenue"`
	NumberOfSoldTicket uint64 `json:"number_of_sold_ticket" bson:"number_of_sold_ticket"`
}

// TicketInformation represents ticket information
type TicketInformation struct {
	TotalNumberOfTicket          uint64 `json:"total_number_of_ticket" bson:"total_number_of_ticket"`
	TotalNumberOfAvailableTicket uint64 `json:"total_number_of_available_ticket" bson:"total_number_of_available_ticket"`
	TotalNumberOFUnsoldTicket    uint64 `json:"total_number_of_unsold_ticket" bson:"total_number_of_unsold_ticket"`
}

// MerchantInformation represents merchant details
type MerchantInformation struct {
	MerchantID          string `json:"merchant_id" bson:"merchant_id"`
	MercahntName        string `json:"merchant_name" bson:"merchant_name"`
	MerchantPhoneNumber string `json:"merchant_phone_number" bson:"merchant_phone_number"`
	MerchantEmail       string `json:"merchant_email" bson:"merchant_email"`
}

// =============================================================================
// SYSTEM MODELS
// =============================================================================

// Notification represents system notifications
type Notification struct {
	ID                string             `json:"id,omitempty" bson:"_id"`
	Title             string             `json:"title" bson:"title"`
	NotificationType  string             `json:"notification_type" bson:"notification_type"`
	NotificationBody  string             `json:"notification_body" bson:"notification_body"`
	IsPublic          bool               `json:"is_public" bson:"is_public"`
	For               NotificationFor    `json:"for" bson:"for"`
	CreatedBy         string             `json:"created_by" bson:"created_by"`
	NotificationParts any                `json:"notification_parts" bson:"notification_parts"`
	Seen              bool               `json:"seen" bson:"seen"`
	Enabled           bool               `json:"enabled" bson:"enabled"`
	Status            NotificationStatus `json:"status" bson:"status"`
	IsDeleted         bool               `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	LastModified      time.Time          `json:"last_modified" bson:"last_modified"`
	DeletedAt         time.Time          `json:"deleted_at" bson:"deleted_at"`
}

type CPSAction struct {
	ID                string        `json:"id,omitempty"`
	ActionCode        string        `json:"action_code" bson:"action_code"`
	CheckerUser       User          `json:"checker_user" bson:"checker_user"`
	MakerUser         User          `json:"maker_user" bson:"maker_user"`
	RejectedReason    string        `json:"rejected_reason,omitempty" bson:"rejected_reason"`
	Department        string        `json:"department,omitempty" bson:"department"`
	Status            ActionStatus  `json:"status,omitempty" bson:"status"`
	RequestAction     RequestAction `json:"request_action,omitempty" bson:"request_action"`
	ActionType        ActionType    `json:"action_type,omitempty" bson:"action_type"`
	ActionData        ActionData    `json:"action_data" bson:"action_data"`
	PreviousAction    any           `json:"previous_action,omitempty"`
	CurrentAction     any           `json:"current_action,omitempty"`
	MakerActionTime   time.Time     `json:"maker_action_time,omitzero" bson:"maker_action_time"`
	CheckerActionTime time.Time     `json:"checker_action_time,omitzero" bson:"checker_action_time"`
	CreatedAt         time.Time     `json:"created_at,omitempty" bson:"created_at"`
	LastModifiedAt    time.Time     `json:"last_modified_at,omitempty" bson:"last_modified_at"`
	UniqueID          string        `json:"unique_id,omitempty" bson:"unique_id"`
}
