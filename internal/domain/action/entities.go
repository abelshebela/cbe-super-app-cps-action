package action

import (
	"time"

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
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type User struct {
	UserID      string    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserCode    string    `bson:"user_code,omitempty" json:"user_code,omitempty"`
	FullName    string    `bson:"full_name,omitempty" json:"full_name,omitempty"`
	PhoneNumber string    `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Timestamp   time.Time `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

type CurrentAction struct {
	Id     []string
	Action bool
}
type CPSAction struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ActionCode      string        `bson:"action_code,omitempty" json:"action_code,omitempty"`
	UniqueId        string        `bson:"unique_id,omitempty" json:"unique_id,omitempty"`
	Maker           User          `bson:"maker,omitempty" json:"maker,omitempty"`
	Checker         User          `bson:"checker,omitempty" json:"checker,omitempty"`
	Department      string        `bson:"department,omitempty" json:"department,omitempty"`
	RejectionReason *string       `bson:"rejection_reason,omitempty" json:"rejection_reason,omitempty"`
	PreviosAction   interface{}   `bson:"previous_action,omitempty" json:"previous_action,omitempty"`
	CurrentAction   interface{}   `bson:"current_action,omitempty" json:"current_action,omitempty"`
	ActionStatus    ActionStatus  `bson:"action_status,omitempty" json:"action_status,omitempty"`
	ActionType      ActionType    `bson:"action_type,omitempty" json:"action_type,omitempty"`
	RequestAction   RequestAction `bson:"request_action,omitempty" json:"request_action,omitempty"`
	CreatedAt       time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LastModifiedAt  time.Time     `bson:"last_modified_at,omitempty" json:"last_modified_at,omitempty"`
	IsDeleted       bool          `bson:"is_deleted" json:"is_deleted,omitempty"`
}

type RequestAction string

const (
	RequestUser                     RequestAction = "USER"
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
	RequestCreateServiceFee         RequestAction = "CREATE SERVICE FEE"
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
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
	RequestUpdateAccountValidation  RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	RequestUpdateServiceDetails     RequestAction = "UPDATE_SERVICE_DETAILS"
	RequestUpdateHQBlockTime        RequestAction = "UPDATE_HQ_BLOCK_TIME"
	RequestUpdateHQArchiveTime      RequestAction = "UPDATE_HQ_ARCHIVE_TIME"
)

type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)

type ProductCodes struct {
	PRD    string
	VATPRD string
	SFPRD  string
	TRXN   string
}

type GLEntry struct {
	ProductAccount    string
	ProductBranchCode string
	ServiceAccount    string
	ServiceBranchCode string
	VatAccount        string
	VatBranchCode     string
}

type Tier struct {
	ID        *string
	Min       uint64
	Max       uint64
	FeeAmount uint64
}

type Cap struct {
	KYCLevel  KYCLevel
	SingleCap uint64
	DailyCap  uint64
	MinAmount uint64
}

type ServiceDetails struct {
	ID                 *string
	ServiceCode        string
	ServiceName        string
	ServiceType        string
	Key                string
	Cap                Cap
	CBEProductCodes    ProductCodes
	CBEIFBProductCodes ProductCodes
	AboveAmount        uint64
	AboveServiceFee    uint64
	PaymentType        string
	Tiers              []Tier
	CBEGLEntry         GLEntry
	CBEIFBGLEntry      GLEntry
	Enabled            bool
	IsDeleted          bool
	CreatedAt          time.Time
	LastModifiedAt     time.Time
	DeletedAt          time.Time
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
	Salt             string
	CurrentPassword  string
	OldPassword      [4]string
	PasswordChangeAt time.Time
}
type CPSUser struct {
	ID                 string
	UserCode           string
	FullName           string
	Role               string
	Department         string
	Gender             string
	PhoneNumber        string
	Email              string
	UserName           string
	Realm              string
	PermissionCategory []string
	PermissionGroup    []string

	Password                 Password
	PasswordDisable          bool
	SyncDisabled             bool
	LoginAttemptCount        uint8
	LastLoginAttempt         time.Time
	NextLoginAttempt         time.Time
	LastOnlineDate           time.Time
	LastLogin                time.Time
	LoginPassword            string
	AccountAuthorizationCode string
	UnlockAccountRequested   bool

	PasswordChangedAt *time.Time
	OTPStatus         string
	OTPLastTriedAt    *time.Time
	OPTLastVerifiedAt *time.Time
	OTPVerifyCount    int

	Enabled      bool
	IsDeleted    bool
	DateJoined   *time.Time
	LastModified *time.Time

	Country string
	Region  string
}
type PasswordRule struct {
	ID             string
	PasswordID     string
	Name           string
	MinLength      int
	MaxLength      int
	Numbers        bool
	CapitalLetters bool
	SmallLetters   bool
	Characters     bool
	CreatedAt      time.Time
}
type Branch struct {
	ID            string
	BranchCode    string
	BranchName    string
	BranchAddress string
	DistrictCode  string
	DistrictName  string
	BranchRegion  string
	RecordStat    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Version       int
	Enabled       bool
}

type Region struct {
	ID            string
	RegionCode    string
	RegionName    string
	RegionAddress string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Enabled       bool
}
type District struct {
	ID              string
	DistrictCode    string
	DistrictName    string
	DistrictAddress string
	RegionID        string
	RegionName      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Enabled         bool
}
type City struct {
	ID           string
	CityCode     string
	CityName     string
	CityAddress  string
	DistrictID   string
	DistrictName string
	RegionID     string
	RegionName   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Enabled      bool
}
