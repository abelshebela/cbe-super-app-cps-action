package model

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
type User struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}
type CPSAction struct {
	ID                 bson.ObjectID   `bson:"_id,omitempty"`
	ActionCode         string          `bson:"action_code"` // Generated
	MakerUser          User            `json:"maker_user" bson:"maker_user"`
	CheckerUser        User            `json:"checker_user" bson:"checker_user"`
	Unique_ID          string          `bson:"unique_id"`
	CheckerID          *string         `bson:"checker_id,omitempty"`
	CheckerName        *string         `bson:"checker_name,omitempty"`
	CheckerPhoneNumber *string         `bson:"checker_phone_number,omitempty"`
	Department         string          `bson:"department"`
	RejectionReason    *string         `bson:"rejection_reason,omitempty"`
	PreviosAction      json.RawMessage `bson:"previos_action"`
	CurrentAction      json.RawMessage `bson:"current_action"`
	ActionStatus       ActionStatus    `bson:"action_status"`
	ActionType         ActionType      `bson:"action_type"`
	RequestAction      RequestAction   `bson:"request_action"`
	CreatedAt          time.Time       `bson:"created_at"`
	LastModifiedAt     time.Time       `bson:"last_modified_at"`
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
	RequestDeleteBank               RequestAction = "DELETE_BANK"
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
)

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

type MiniApp struct {
	ID                string                  `bson:"id"`
	AppName           string                  `bson:"app_name"`
	AppIcon           string                  `bson:"app_icon"`
	CommisonGLAccount string                  `bson:"commison_gl_account"`
	AppType           AppType                 `bson:"app_type"`
	MerchantID        string                  `bson:"merchant_id"`
	ProductCode       []ProductCode           `bson:"product_code"`
	Credential        []CredentialInformation `bson:"credential"`
	IsEventMiniApp    bool                    `bson:"is_event_mini_app"`
	IsThreeClick      bool                    `bson:"is_three_click"`
	Enabled           bool                    `bson:"enabled"`
	IsDeleted         bool                    `bson:"is_deleted"`
	CreatedAt         time.Time               `bson:"created_at"`
	LastModifiedAt    time.Time               `bson:"last_modified_at"`
	DeletedAt         time.Time               `bson:"deleted_at"`
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

	Password                 Password  `json:"password" bson:"password"`
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

	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" bson:"password_changed_at,omitempty"`
	OTPStatus         string     `json:"otp_status,omitempty" bson:"otp_status,omitempty"`
	OTPLastTriedAt    *time.Time `json:"otp_last_tried_at,omitempty" bson:"otp_last_tried_at,omitempty"`
	OPTLastVerifiedAt *time.Time `json:"otp_last_verified_at,omitempty" bson:"otp_last_verified_at,omitempty"`
	OTPVerifyCount    int        `json:"otp_verify_count,omitempty" bson:"otp_verify_count,omitempty"`

	Enabled      bool       `json:"enabled,omitempty" bson:"enabled,omitempty"`
	IsDeleted    bool       `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	DateJoined   *time.Time `json:"date_joined,omitempty" bson:"date_joined,omitempty"`
	LastModified *time.Time `json:"last_modified,omitempty" bson:"last_modified,omitempty"`

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
type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)

type ProductCodes struct {
	PRD    string `bson:"prd"`
	VATPRD string `bson:"vatprd"`
	SFPRD  string `bson:"sfprd"`
	TRXN   string `bson:"trxn"`
}

type GLEntry struct {
	ProductAccount    string `bson:"product_account"`
	ProductBranchCode string `bson:"product_branch_code"`
	ServiceAccount    string `bson:"service_account"`
	ServiceBranchCode string `bson:"service_branch_code"`
	VatAccount        string `bson:"vat_account"`
	VatBranchCode     string `bson:"vat_branch_code"`
}

type Tier struct {
	ID        bson.ObjectID `bson:"id"`
	Min       uint64        `bson:"min"`
	Max       uint64        `bson:"max"`
	FeeAmount uint64        `bson:"fee_amount"`
}

type Cap struct {
	KYCLevel  KYCLevel `bson:"kyc_level"`
	SingleCap uint64   `bson:"single_cap"`
	DailyCap  uint64   `bson:"daily_cap"`
	MinAmount uint64   `bson:"min_amount"`
}

type ServiceDetails struct {
	ID                 bson.ObjectID `bson:"_id,omitempty"`
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

type CpsAction struct {
	ID                string        `json:"id" bson:"id"`
	ActionCode        string        `json:"action_code" bson:"action_code"`
	MakerUser         User          `json:"maker_user" bson:"maker_user"`
	CheckerUser       User          `json:"checker_user,omitempty" bson:"checker_user"`
	RejectedReason    string        `json:"rejected_reason,omitempty" bson:"rejected_reason"`
	Department        string        `json:"department" bson:"department"`
	Status            ActionStatus  `json:"status" bson:"status"`
	RequestAction     RequestAction `json:"request_action" bson:"request_action"`
	ActionType        ActionType    `json:"action_type" bson:"action_type"`
	ActionData        any           `json:"action_data" bson:"action_data"`
	PreviousData      any           `json:"previous_action,omitempty"`
	CurrentData       any           `json:"current_action,omitempty"`
	MakerActionTime   time.Time     `json:"maker_action_time,omitzero" bson:"maker_action_time"`
	CheckerActionTime time.Time     `json:"checker_action_time,omitzero" bson:"checker_action_time"`
}
