package action

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities"
)

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

type User struct {
	UserID      string    `json:"user_id" bson:"user_id"`
	UserCode    string    `json:"user_code" bson:"user_code"`
	FullName    string    `json:"full_name" bson:"full_name"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Department  string    `json:"department" bson:"department"`
}

type CurrentAction struct {
	Id            []string `json:"id" bson:"id"`
	Action        bool     `json:"action" bson:"action"`
	RequestAction string   `json:"request_action" bson:"request_action"`
}

type ActionResponse struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	ActionId string `json:"action_id" bson:"action_id"`
}

type CPSAction struct {
	ID                 string        `json:"id" bson:"id"`
	ActionCode         string        `json:"action_code" bson:"action_code"`
	UniqueId           string        `json:"unique_id" bson:"unique_id"`
	MakerID            string        `json:"maker_id" bson:"maker_id"`
	MakerName          string        `json:"maker_name" bson:"maker_name"`
	MakerPhoneNumber   string        `json:"maker_phone_number" bson:"maker_phone_number"`
	CheckerID          string        `json:"checker_id,omitempty" bson:"checker_id,omitempty"`
	CheckerName        string        `json:"checker_name,omitempty" bson:"checker_name,omitempty"`
	CheckerPhoneNumber string        `json:"checker_phone_number,omitempty" bson:"checker_phone_number,omitempty"`
	Department         string        `json:"department" bson:"department"`
	RejectionReason    *string       `json:"rejection_reason,omitempty" bson:"rejection_reason,omitempty"`
	PreviousAction     interface{}   `json:"previos_action" bson:"previos_action"`
	CurrentAction      interface{}   `json:"current_action" bson:"current_action"`
	ActionStatus       ActionStatus  `json:"action_status" bson:"action_status"`
	ActionType         ActionType    `json:"action_type" bson:"action_type"`
	RequestAction      RequestAction `json:"request_action" bson:"request_action"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	MakerActionTime    time.Time     `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime  *time.Time    `json:"checker_action_time,omitempty" bson:"checker_action_time,omitempty"`
}

type RequestAction string

const (
	RequestUser                     RequestAction = "USER"
	RequestPermissionGroup          RequestAction = "PERMISSION_GROUP"
	RequestDepartment               RequestAction = "DEPARMTENT"
	RequestBPSUser                  RequestAction = "BPS_USER"
	RequestDisableBPSUser           RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser            RequestAction = "ENABLE_BPS_USER"
	RequestUpdateUser               RequestAction = "UPDATE_USER"
	RequestTotalDailyLimit          RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT                RequestAction = "UPDATE_VAT"
	RequestAuthTier                 RequestAction = "AUTHTIER"
	RequestCreateAdvert             RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert             RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert             RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert            RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert             RequestAction = "DELETE_ADVERT"
	RequestCreateBank               RequestAction = "CREATE_BANK"
	RequestUpdateBank               RequestAction = "UPDATE_BANK"
	RequestEnableBank               RequestAction = "ENABLE_BANK"
	RequestDisableBank              RequestAction = "DISABLE_BANK"
	RequestEnableWallet             RequestAction = "ENABLE_WALLET"
	RequestDisableWallet            RequestAction = "DISABLE_WALLET"
	RequestUpdatePasswordExpiry     RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreateValidation         RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation         RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation         RequestAction = "DELETE_VALIDATION"
	RequestUpdateArchiveExpiry      RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestCreateServiceFee         RequestAction = "CREATE_SERVICE_FEE"
	RequestUpdateServiceFee         RequestAction = "UPDATE SERVICE FEE"
	RequestDeleteServiceFee         RequestAction = "DELETE SERVICE FEE"
	RequestCreateDailyLimit         RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit         RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit         RequestAction = "DELETE DAILY LIMIT"
	RequestBudgetUpdate             RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestBudgetCreate             RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestBudgetDelete             RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestBudgetColor              RequestAction = "BUDGET_COLOR"
	RequestBudgetIcon               RequestAction = "BUDGET_ICON"
	RequestUpdateProduct            RequestAction = "UPDATE_PRODUCT"
	RequestCreatePublicNotification RequestAction = "CREATE_PUBLIC_NOTIFICATION"
	RequestArchiveUser              RequestAction = "ARCHIVE_USER"
	RequestCreatePasswordRule       RequestAction = "CREATE_PASSWORD_RULE"
	RequestUpdatePasswordRule       RequestAction = "UPDATE_PASSWORD_RULE"
	RequestUpdateMinimumService     RequestAction = "UPDATE_MINIMUM_SERVICE"
	RequestUpdateServiceRule        RequestAction = "UPDATE_SERVICE_RULE"
	RequestUpdateTotal              RequestAction = "UPDATE_TOTAL"
	RequestUpdateAccessConfig       RequestAction = "UPDATE_ACCESS_CONFIG"
	RequestEnableSingleBranch       RequestAction = "ENABLE_SINGLE_BRANCH"
	RequestDisableSingleBranch      RequestAction = "DISABLE_SINGLE_BRANCH"
	RequestEnableMultiUsers         RequestAction = "ENABLE_MULTI_USERS"
	RequestDisableMultiUsers        RequestAction = "DISABLE_MULTI_USERS"
	RequestCreateBusiness           RequestAction = "CREATE_BUSINESS"
	RequestUpdateBusiness           RequestAction = "UPDATE_BUSINESS"
	RequestCreateEvent              RequestAction = "CREATE_EVENT"
	RequestUpdateEvent              RequestAction = "UPDATE_EVENT"
	RequestCreateEventCategory      RequestAction = "CREATE_EVENT_CATEGORY"
	RequestUpdateEventCategory      RequestAction = "UPDATE_EVENT_CATEGORY"
	RequestDisableEvent             RequestAction = "DISABLE_EVENT"
	RequestCreateMiniAppMerchant    RequestAction = "CREATE_MINIAPP_MERCHANT"
	RequestUpdateMiniAppMerchant    RequestAction = "UPDATE_MINIAPP_MERCHANT"
	RequestDeleteMiniAppMerchant    RequestAction = "DELETE_MINIAPP_MERCHANT"
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
	RequestUpdateAccountValidation  RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	RequestUpdateServiceDetails     RequestAction = "UPDATE_SERVICE_DETAILS"

	RequestUpdateEevent      RequestAction = "UPDATE_EVENT"
	RequestCIFRemove         RequestAction = "CIF_REMOVE"
	RequestServiceFlagUpdate RequestAction = "SERVICE_FLAG_UPDATE"
)

type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)

type ProductCodes struct {
	PRD    string `json:"prd" bson:"prd"`
	VATPRD string `json:"vatprd" bson:"vatprd"`
	SFPRD  string `json:"sfprd" bson:"sfprd"`
	TRXN   string `json:"trxn" bson:"trxn"`
}

type GLEntry struct {
	ProductAccount    string `json:"product_account" bson:"product_account"`
	ProductBranchCode string `json:"product_branch_code" bson:"product_branch_code"`
	ServiceAccount    string `json:"service_account" bson:"service_account"`
	ServiceBranchCode string `json:"service_branch_code" bson:"service_branch_code"`
	VatAccount        string `json:"vat_account" bson:"vat_account"`
	VatBranchCode     string `json:"vat_branch_code" bson:"vat_branch_code"`
}

type Tier struct {
	ID        *string `json:"id" bson:"id"`
	Min       uint64  `json:"min" bson:"min"`
	Max       uint64  `json:"max" bson:"max"`
	FeeAmount uint64  `json:"fee_amount" bson:"fee_amount"`
}

type Cap struct {
	SingleCap uint64 `json:"single_cap" bson:"single_cap"`
	DailyCap  uint64 `json:"daily_cap" bson:"daily_cap"`
	MinAmount uint64 `json:"min_amount" bson:"min_amount"`
}

type ServiceDetails struct {
	ID                 *string              `json:"id" bson:"id"`
	ServiceCode        string               `json:"service_code" bson:"service_code"`
	ServiceName        string               `json:"service_name" bson:"service_name"`
	ServiceType        string               `json:"service_type" bson:"service_type"`
	Key                string               `json:"key" bson:"key"`
	Cap                Cap                  `json:"cap" bson:"cap"`
	CBEProductCodes    ProductCodes         `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes         `json:"cbeifb_product_codes" bson:"cbeifb_product_codes"`
	AboveAmount        uint64               `json:"above_amount" bson:"above_amount"`
	AboveServiceFee    uint64               `json:"above_service_fee" bson:"above_service_fee"`
	PaymentType        entities.PaymentType `json:"payment_type" bson:"payment_type"`
	Tiers              []Tier               `json:"tiers" bson:"tiers"`
	CBEGLEntry         GLEntry              `json:"cbegl_entry" bson:"cbegl_entry"`
	CBEIFBGLEntry      GLEntry              `json:"cbeifbgl_entry" bson:"cbeifbgl_entry"`
	Enabled            bool                 `json:"enabled" bson:"enabled"`
	IsDeleted          bool                 `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time            `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time            `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time            `json:"deleted_at" bson:"deleted_at"`
}

type RegistrationType string

const (
	RegistrationTypeNew    RegistrationType = "NEW"
	RegistrationTypeLinked RegistrationType = "LINKED"
)

type LinkedAccount struct {
	ID                *string
	UserID            *string
	CustomerNumber    string
	AccountNumber     string
	AccountHolderName string
	AccountType       string
	BranchCode        string
	LinkedStatus      bool
	LastLinkedStatus  bool
	LinkedAt          time.Time
	LinkedBranch      string
	RegistrationType  RegistrationType
	IsAccountActive   bool
	AndOrStatus       bool
	AccountBranchCode string
	CurrencyCode      string
	IsMain            bool // default: false, first account: true
	MakerAndChecker   struct {
		Linkers struct {
			Maker   string
			Checker string
		}
		Unlinkers struct {
			Maker   string
			Checker string
		}
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password"`
	PasswordChangeAt time.Time `json:"password_change_at" bson:"password_change_at"`
}

type CPSUser struct {
	ID                 string   `json:"id" bson:"id"`
	UserCode           string   `json:"user_code" bson:"user_code"`
	FullName           string   `json:"full_name" bson:"full_name"`
	Role               string   `json:"role" bson:"role"`
	Department         string   `json:"department" bson:"department"`
	Gender             string   `json:"gender" bson:"gender"`
	PhoneNumber        string   `json:"phone_number" bson:"phone_number"`
	Email              string   `json:"email" bson:"email"`
	UserName           string   `json:"user_name" bson:"user_name"`
	Realm              string   `json:"realm" bson:"realm"`
	PermissionCategory []string `json:"permission_category" bson:"permission_category"`
	PermissionGroup    []string `json:"permission_group" bson:"permission_group"`

	Password                 Password  `json:"password" bson:"password"`
	PasswordDisable          bool      `json:"password_disable" bson:"password_disable"`
	SyncDisabled             bool      `json:"sync_disabled" bson:"sync_disabled"`
	LoginAttemptCount        uint8     `json:"login_attempt_count" bson:"login_attempt_count"`
	LastLoginAttempt         time.Time `json:"last_login_attempt" bson:"last_login_attempt"`
	NextLoginAttempt         time.Time `json:"next_login_attempt" bson:"next_login_attempt"`
	LastOnlineDate           time.Time `json:"last_online_date" bson:"last_online_date"`
	LastLogin                time.Time `json:"last_login" bson:"last_login"`
	LoginPassword            string    `json:"login_password" bson:"login_password"`
	AccountAuthorizationCode string    `json:"account_authorization_code" bson:"account_authorization_code"`
	UnlockAccountRequested   bool      `json:"unlock_account_requested" bson:"unlock_account_requested"`

	PasswordChangedAt *time.Time `json:"password_changed_at" bson:"password_changed_at"`
	OTPStatus         string     `json:"otp_status" bson:"otp_status"`
	OTPLastTriedAt    *time.Time `json:"otp_last_tried_at" bson:"otp_last_tried_at"`
	OPTLastVerifiedAt *time.Time `json:"otp_last_verified_at" bson:"otp_last_verified_at"`
	OTPVerifyCount    int        `json:"otp_verify_count" bson:"otp_verify_count"`

	Enabled      bool       `json:"enabled" bson:"enabled"`
	IsDeleted    bool       `json:"is_deleted" bson:"is_deleted"`
	DateJoined   *time.Time `json:"date_joined" bson:"date_joined"`
	LastModified *time.Time `json:"last_modified" bson:"last_modified"`

	Country string `json:"country" bson:"country"`
	Region  string `json:"region" bson:"region"`
}

type PasswordRule struct {
	ID             string    `json:"id" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	MinLength      int       `json:"min_length" bson:"min_length"`
	MaxLength      int       `json:"max_length" bson:"max_length"`
	Numbers        bool      `json:"numbers" bson:"numbers"`
	CapitalLetters bool      `json:"capital_letters" bson:"capital_letters"`
	SmallLetters   bool      `json:"small_letters" bson:"small_letters"`
	Characters     bool      `json:"characters" bson:"characters"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

type Branch struct {
	ID            string    `json:"id" bson:"id"`
	BranchCode    string    `json:"branch_code" bson:"branch_code"`
	BranchName    string    `json:"branch_name" bson:"branch_name"`
	BranchAddress string    `json:"branch_address" bson:"branch_address"`
	DistrictCode  string    `json:"district_code" bson:"district_code"`
	DistrictName  string    `json:"district_name" bson:"district_name"`
	BranchRegion  string    `json:"branch_region" bson:"branch_region"`
	RecordStat    string    `json:"record_stat" bson:"record_stat"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
	Version       int       `json:"version" bson:"version"`
	Enabled       bool      `json:"enabled" bson:"enabled"`
}

type Region struct {
	ID            string    `json:"id" bson:"id"`
	RegionCode    string    `json:"region_code" bson:"region_code"`
	RegionName    string    `json:"region_name" bson:"region_name"`
	RegionAddress string    `json:"region_address" bson:"region_address"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
	Enabled       bool      `json:"enabled" bson:"enabled"`
}

type District struct {
	ID              string    `json:"id" bson:"id"`
	DistrictCode    string    `json:"district_code" bson:"district_code"`
	DistrictName    string    `json:"district_name" bson:"district_name"`
	DistrictAddress string    `json:"district_address" bson:"district_address"`
	RegionID        string    `json:"region_id" bson:"region_id"`
	RegionName      string    `json:"region_name" bson:"region_name"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
	Enabled         bool      `json:"enabled" bson:"enabled"`
}

type City struct {
	ID           string    `json:"id" bson:"id"`
	CityCode     string    `json:"city_code" bson:"city_code"`
	CityName     string    `json:"city_name" bson:"city_name"`
	CityAddress  string    `json:"city_address" bson:"city_address"`
	DistrictID   string    `json:"district_id" bson:"district_id"`
	DistrictName string    `json:"district_name" bson:"district_name"`
	RegionID     string    `json:"region_id" bson:"region_id"`
	RegionName   string    `json:"region_name" bson:"region_name"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	Enabled      bool      `json:"enabled" bson:"enabled"`
}

type CreateCPSAction struct {
	Department    string        `json:"department,omitempty" bson:"department"`
	RequestAction RequestAction `json:"request_action,omitempty" bson:"request_action"`
	MakerUser     User          `json:"maker_user" bson:"maker_user"`
}
