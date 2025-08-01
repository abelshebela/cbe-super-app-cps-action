package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/v2/bson"
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
	ActionCreate  ActionType = "CREATE"
	ActionUpdate  ActionType = "UPDATE"
	ActionDelete  ActionType = "DELETE"
	ActionEnable  ActionType = "ENABLE"
	ActionDisable ActionType = "DISABLE"
)

type CurrentAction struct {
	Id     []string `bson:"id"`
	Action bool     `bson:"action"`
}
type User struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Department  string `json:"department" bson:"department"`
}

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


type CPSAction struct {
	ID                 bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ActionCode         string        `bson:"action_code" json:"action_code"`
	UniqueId           string        `bson:"unique_id" json:"unique_id,omitempty"`
	MakerID            string        `bson:"maker_id" json:"maker_id"`
	MakerName          string        `bson:"maker_name" json:"maker_name"`
	MakerPhoneNumber   string        `bson:"maker_phone_number" json:"maker_phone_number"`
	CheckerID          string        `bson:"checker_id" json:"checker_id,omitempty"`
	CheckerName        string        `bson:"checker_name" json:"checker_name,omitempty"`
	CheckerPhoneNumber string        `bson:"checker_phone_number" json:"checker_phone_number,omitempty"`
	Department         string        `bson:"department" json:"department"`
	RejectionReason    string        `bson:"rejection_reason" json:"rejection_reason,omitempty"`
	PreviousAction     interface{}   `bson:"previous_action" json:"previous_action"`
	CurrentAction      interface{}   `bson:"current_action" json:"current_action"`
	ActionStatus       string        `bson:"action_status" json:"action_status"`
	ActionType         string        `bson:"action_type" json:"action_type"`
	RequestAction      string        `bson:"request_action" json:"request_action"`
	CreatedAt          time.Time     `bson:"created_at" json:"created_at"`
	LastModifiedAt     time.Time     `bson:"last_modified_at" json:"last_modified_at"`
	MakerActionTime    time.Time     `bson:"maker_action_time" json:"maker_action_time"`
	CheckerActionTime  *time.Time    `bson:"checker_action_time" json:"checker_action_time,omitempty"`
}

type RequestAction string

const (
	RequestUser                     RequestAction = "USER"
	RequestCpsUserCreate            RequestAction = "CREATE_CPS_USER"
	RequestCpsUserUpdate            RequestAction = "UPDATE_CPS_USER"
	RequestCpsUserDelete            RequestAction = "DELETE_CPS_USER"
	RequestCpsUserEnable            RequestAction = "ENABLE_CPS_USER"
	RequestCpsUserDisable           RequestAction = "DISABLE_CPS_USER"
	RequestPermissionGroup          RequestAction = "PERMISSION_GROUP"
	RequestDepartment               RequestAction = "DEPARMTENT"
	RequestEnableUser               RequestAction = "ENABLE_USER"
	RequestDisableUser              RequestAction = "DISABLE_USER"
	RequestBPSUser                  RequestAction = "BPS_USER"
	RequestDisableBPSUser           RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser            RequestAction = "ENABLE_BPS_USER"
	RequestUpdateUser               RequestAction = "UPDATE_USER"
	RequestTotalDailyLimit          RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT                RequestAction = "UPDATE_VAT"
	RequestAuthTier                 RequestAction = "AUTHTIER"
	RequestDeleteAmountBasedAuth    RequestAction = "DELETE_AMOUNT_BASED_AUTH"
	RequestCreateAmountBasedAuth    RequestAction = "CREATE_AMOUNT_BASED_AUTH"
	RequestUpdateAmountBasedAuth    RequestAction = "UPDATE_AMOUNT_BASED_AUTH"
	RequestCreateAdvert             RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert             RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert             RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert            RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert             RequestAction = "DELETE_ADVERT"
	RequestCreateBank               RequestAction = "CREATE_BANK"
	RequestUpdateBank               RequestAction = "UPDATE_BANK"
	RequestDeleteBank               RequestAction = "DELETE_BANK"
	RequestEnableBank               RequestAction = "ENABLE_BANK"
	RequestDisableBank              RequestAction = "DISABLE_BANK"
	RequestCreateWallet             RequestAction = "CREATE_WALLET"
	RequestUpdateWallet             RequestAction = "UPDATE_WALLET"
	RequestDeleteWallet             RequestAction = "DELETE_WALLET"
	RequestEnableWallet             RequestAction = "ENABLE_WALLET"
	RequestDisableWallet            RequestAction = "DISABLE_WALLET"
	RequestUpdatePasswordExpiry     RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreateValidation         RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation         RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation         RequestAction = "DELETE_VALIDATION"
	RequestUpdateArchiveExpiry      RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestUpdateServiceSingle      RequestAction = "UPDATE_SERVICE_SINGLE_CAP"
	RequestUpdateServiceTotal       RequestAction = "UPDATE_SERVICE_TOTAL_CAP"
	RequestUpdateServiceMinCap      RequestAction = "UPDATE_SERVICE_MIN_CAP"
	RequestCreateServiceFee         RequestAction = "CREATE_SERVICE_FEE"
	RequestUpdateServiceFee         RequestAction = "UPDATE_SERVICE_FEE"
	RequestDeleteServiceFee         RequestAction = "DELETE_SERVICE_FEE"
	RequestCreateDailyLimit         RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit         RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit         RequestAction = "DELETE DAILY LIMIT"
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
	RequestCreateMiniAppMerchant    RequestAction = "CREATE_MINIAPP_MERCHANT"
	RequestUpdateMiniAppMerchant    RequestAction = "UPDATE_MINIAPP_MERCHANT"
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar             RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar             RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar             RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar            RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar             RequestAction = "UPDATE_AVATAR"
	RequestBlockRegion              RequestAction = "BLOCK_REGION"
	RequestEnableRegion             RequestAction = "ENABLE_REGION"
	RequestBlockDistrict            RequestAction = "BLOCK_DISTRICT"
	RequestEnableDistrict           RequestAction = "ENABLE_DISTRICT"
	RequestBlockCity                RequestAction = "BLOCK_CITY"
	RequestEnableCity               RequestAction = "ENABLE_CITY"
	RequestBlockUser                RequestAction = "BLOCK_USER"
	RequestEnableSingleBranches     RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches    RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches      RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches     RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"
)

type RegistrationType string

const (
	RegistrationTypeNew    RegistrationType = "NEW"
	RegistrationTypeLinked RegistrationType = "LINKED"
)

type LinkedAccountExternal struct {
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

func (ext *LinkedAccountExternal) MapFromExternal() LinkedAccount {
	return LinkedAccount{
		AccountNumber:     ext.AccountNumber,
		AccountType:       ext.AccountType,
		AccountHolderName: ext.CustomerName,
		CustomerNumber:    ext.CustomerNumber,
		BranchCode:        ext.AccountBranchCode,
		AccountBranchCode: ext.AccountBranchCode,
		CurrencyCode:      ext.AccountCurrency,
		IsAccountActive:   ext.ActiveAccount,
	}
}

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

type CpsActionNormalized struct {
	ID                 string      `bson:"_id,omitempty" json:"id,omitempty"`
	ActionCode         string      `bson:"action_code" json:"action_code,omitempty"`
	UniqueId           string      `bson:"unique_id" json:"unique_id,omitempty"`
	MakerID            string      `bson:"maker_id" json:"maker_id"`
	MakerName          string      `bson:"maker_name" json:"maker_name"`
	MakerPhoneNumber   string      `bson:"maker_phone_number" json:"maker_phone_number"`
	CheckerID          string      `bson:"checker_id" json:"checker_id,omitempty"`
	CheckerName        string      `bson:"checker_name" json:"checker_name,omitempty"`
	CheckerPhoneNumber string      `bson:"checker_phone_number" json:"checker_phone_number,omitempty"`
	Department         string      `bson:"department" json:"department"`
	RejectionReason    string      `bson:"rejection_reason" json:"rejection_reason,omitempty"`
	PreviousAction     interface{} `bson:"previous_action" json:"previous_action"`
	CurrentAction      interface{} `bson:"current_action" json:"current_action"`
	ActionStatus       string      `bson:"action_status" json:"action_status"`
	ActionType         string      `bson:"action_type" json:"action_type"`
	RequestAction      string      `bson:"request_action" json:"request_action"`
	CreatedAt          time.Time   `bson:"created_at" json:"created_at"`
	LastModifiedAt     time.Time   `bson:"last_modified_at" json:"last_modified_at"`
	MakerActionTime    time.Time   `bson:"maker_action_time" json:"maker_action_time"`
	CheckerActionTime  time.Time   `bson:"checker_action_time" json:"checker_action_time"`
}


type Card struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CardName string        `bson:"card_name" json:"card_name"`
	SubCards []string      `bson:"sub_cards" json:"sub_cards"`
}

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
	Realm              string          `json:"realm,omitempty" bson:"realm,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category,omitempty"`
	PermissionGroup    []bson.ObjectID `json:"permission_group,omitempty" bson:"permission_group,omitempty"`

	Password                 Password  `json:"password" bson:"password,omitempty"`
	PasswordDisable          bool      `json:"password_disable,omitempty" bson:"password_disable,omitempty"`
	SyncDisabled             bool      `json:"sync_disabled,omitempty" bson:"sync_disabled,omitempty"`
	LoginAttemptCount        uint8     `json:"login_attempt_count,omitempty" bson:"login_attempt_count,omitempty"`
	LastLoginAttempt         time.Time `json:"last_login_attempt,omitempty" bson:"last_login_attempt,omitempty"`
	NextLoginAttempt         time.Time `json:"next_login_attempt,omitempty" bson:"next_login_attempt,omitempty"`
	LastOnlineDate           time.Time `json:"last_online_date,omitempty" bson:"last_online_date,omitempty"`
	LastLogin                time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
	LoginPassword            string    `json:"login_password,omitempty" bson:"login_password,omitempty"`
	AccountAuthorizationCode string    `json:"account_authorization_code,omitempty" bson:"account_authorization_code,omitempty"`
	UnlockAccountRequested   bool      `json:"unlock_account_requested,omitempty" bson:"unlock_account_requested,omitempty"`

	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" bson:"password_changed_at"`
	OTPStatus         string     `json:"otp_status,omitempty" bson:"otp_status"`
	OTPLastTriedAt    *time.Time `json:"otp_last_tried_at,omitempty" bson:"otp_last_tried_at"`
	OPTLastVerifiedAt *time.Time `json:"otp_last_verified_at,omitempty" bson:"otp_last_verified_at"`
	OTPVerifyCount    int        `json:"otp_verify_count,omitempty" bson:"otp_verify_count"`

	Enabled      bool       `json:"enabled,omitempty" bson:"enabled"`
	IsDeleted    bool       `json:"is_deleted,omitempty" bson:"is_deleted"`
	DateJoined   *time.Time `json:"date_joined,omitempty" bson:"date_joined"`
	LastModified *time.Time `json:"last_modified,omitempty" bson:"last_modified"`

	Country string `json:"country,omitempty" bson:"country,omitempty"`
	Region  string `json:"region,omitempty" bson:"region,omitempty"`
}
type Password struct {
	Salt             string    `json:"salt" bson:"salt"`
	CurrentPassword  string    `json:"current_password" bson:"current_password"`
	OldPassword      [4]string `json:"old_password" bson:"old_password,omitempty"`
	PasswordChangeAt time.Time `json:"password_changed_at" bson:"password_changed_at"`
}
type ValidationRule struct {
	ID             bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EntityType     string        `json:"entity_type" bson:"entity_type"`
	ValidationFor  string        `json:"validation_for" bson:"validation_for"`
	Identifier     string        `json:"identifier" bson:"identifier"`
	MinLength      int           `json:"min_length" bson:"min_length"`
	MaxLength      int           `json:"max_length" bson:"max_length"`
	Enabled        bool          `json:"enabled" bson:"enabled"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	ServiceID      string        `json:"service_id" bson:"service_id"`
}






type ServiceDetails struct {
	ID                 bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ServiceCode        string        `bson:"service_code"`
	ServiceName        string        `bson:"service_name"`
	ServiceType        string        `bson:"service_type"`
	Key                string        `bson:"key"`
	Cap                Cap           `bson:"cap"`
	CBEProductCodes    ProductCodes  `bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes  `bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64        `bson:"above_amount"`
	AboveServiceFee    uint64        `bson:"above_service_fee"`
	PaymentType        string        `bson:"payment_type"`
	Tiers              []Tier        `bson:"tiers"`
	CBEGLEntry         GLEntry       `bson:"cbe_gl_entry"`
	CBEIFBGLEntry      GLEntry       `bson:"cbe_ifb_gl_entry"`
	Enabled            bool          `bson:"enabled"`
	IsDeleted          bool          `bson:"is_deleted"`
	CreatedAt          time.Time     `bson:"created_at"`
	LastModifiedAt     time.Time     `bson:"last_modified_at"`
	DeletedAt          time.Time     `bson:"deleted_at"`
}

type CreateCPSAction struct {
	ActionCode      string        `json:"action_code" bson:"action_code"`
	MakerUser       User          `json:"maker_user" bson:"maker_user"`
	Department      string        `json:"department,omitempty" bson:"department"`
	Status          ActionStatus  `json:"status,omitempty" bson:"status"`
	RequestAction   RequestAction `json:"request_action,omitempty" bson:"request_action"`
	ActionType      ActionType    `json:"action_type,omitempty" bson:"action_type"`
	ActionData      any           `json:"action_data" bson:"action_data"`
	PreviousData    any           `json:"previous_action,omitempty"`
	CurrentData     any           `json:"current_action,omitempty"`
	MakerActionTime time.Time     `json:"maker_action_time,omitzero" bson:"maker_action_time"`
}

type AuthorizeCPSAction struct {
	ActionCode        string    `json:"action_code" bson:"action_code"`
	Department        string    `json:"department,omitempty" bson:"department"`
	CheckerUser       User      `json:"checker_user" bson:"checker_user"`
	CheckerActionTime time.Time `json:"checker_action_time,omitzero" bson:"checker_action_time"`
}

type RejectAuthTierCPSAction struct {
	RejectionReason string `json:"rejection_reason,omitempty" bson:"rejection_reason"`
}
type RejectCPSAction struct {
	CreateCPSAction
	CheckerUser       User      `json:"checker_user" bson:"checker_user"`
	RejectedReason    string    `json:"rejected_reason,omitempty" bson:"rejected_reason"`
	CheckerActionTime time.Time `json:"checker_action_time,omitzero" bson:"checker_action_time"`
}

func (r RejectCPSAction) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RejectedReason, validation.Required, validation.Length(30, 100)),
	)
}

func (r RejectAuthTierCPSAction) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RejectionReason, validation.Required, validation.Length(30, 100)),
	)
}

type PasswordRule struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name           string        `bson:"name" json:"name"`
	MinLength      int           `bson:"min_length" json:"min_length"`
	MaxLength      int           `bson:"max_length" json:"max_length"`
	Numbers        bool          `bson:"numbers" json:"numbers"`
	CapitalLetters bool          `bson:"capital_letters" json:"capital_letters"`
	SmallLetters   bool          `bson:"small_letters" json:"small_letters"`
	Characters     bool          `bson:"characters" json:"characters"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
}

type HQ struct {
	ID                bson.ObjectID `bson:"_id" json:"id"`
	UniqueId          string        `bson:"unique_id" json:"unique_id"`
	Name              string        `bson:"name" json:"name"`
	TotalCap          string        `bson:"total_cap" json:"total_cap"`
	BlockTime         uint          `bson:"block_time" json:"block_time"`
	BlockTimeStatus   string        `bson:"block_time_status" json:"block_time_status"`
	ArchiveTime       uint          `bson:"archive_time" json:"archive_time"`
	ArchiveTimeStatus string        `bson:"archive_time_status" json:"archive_time_status"`
	CreatedAt         time.Time     `bson:"created_at" json:"created_at"`
	LastModifiedAt    time.Time     `bson:"last_modified_at" json:"last_modified_at"`
}
type Branch struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BranchCode    string        `bson:"branch_code" json:"branch_code"`
	BranchName    string        `bson:"branch_name" json:"branch_name"`
	BranchAddress string        `bson:"branch_address" json:"branch_address"`
	DistrictCode  string        `bson:"district_code" json:"district_code"`
	DistrictName  string        `bson:"district_name" json:"district_name"`
	BranchRegion  string        `bson:"branch_region" json:"branch_region"`
	RecordStat    string        `bson:"record_stat" json:"record_stat"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
	Enabled       bool          `bson:"enabled" json:"enabled"`
}

type Region struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RegionCode    string        `json:"region_code" bson:"region_code"`
	RegionName    string        `json:"region_name" bson:"region_name"`
	RegionAddress string        `json:"region_address" bson:"region_address"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
}
type District struct {
	ID              bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DistrictCode    string        `json:"district_code" bson:"district_code"`
	DistrictName    string        `json:"district_name" bson:"district_name"`
	DistrictAddress string        `json:"district_address" bson:"district_address"`
	RegionID        string        `json:"region_id" bson:"region_id"`
	RegionName      string        `json:"region_name" bson:"region_name"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
	Enabled         bool          `json:"enabled" bson:"enabled"`
}
type City struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CityCode     string        `json:"city_code" bson:"city_code"`
	CityName     string        `json:"city_name" bson:"city_name"`
	City         string        `json:"city_address" bson:"city_address"`
	DistrictID   string        `json:"district_id" bson:"district_id"`
	DistrictName string        `json:"district_name" bson:"district_name"`
	RegionID     string        `json:"region_id" bson:"region_id"`
	RegionName   string        `json:"region_name" bson:"region_name"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
	Enabled      bool          `json:"enabled" bson:"enabled"`
}
