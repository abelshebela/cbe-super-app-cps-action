package constants

import (
	"fmt"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Type Definitions
// ---------------------------------------------------------------------------

type (
	ContextKey               string
	Platform                 string
	AdvertFor                string
	RegistrationType         string
	Method                   string
	AccrualMethod            string
	AccrualFrequency         string
	RestrictionType          string
	EventStatus              string
	AppType                  string
	EnvironmentType          string
	AppViewType              string
	NotificationFor          string
	NotificationStatus       string
	BlockedOn                string
	VaultStatus              string
	VaultType                string
	VaultCategory            string
	BranchType               string
	CurrencyType             string
	Realm                    string
	FeeType                  string
	TransactionMethod        string
	AccountType              string
	Gender                   string
	Vendor                   string
	KYCStatus                string
	BPSStatus                string
	MaritalStatus            string
	DeviceStatus             string
	OnboardingMethod         string
	EnabledChannels          string
	AccountStatus            string
	MemberType               string
	OTPFor                   string
	OTPStatus                string
	RiskLevel                string
	SettlementMethod         string
	AuditorStatus            string
	AuditorMark              string
	CustomerStatus           string
	EmploymentStatus         string
	Level                    uint8
	ActionType               string
	RequestAction            string
	PrimaryAuthentication    string
	Stage                    string
	InviteStatus             string
	MemberStatus             string
	MemberRole               string
	Recurrence               string
	ContributionType         string
	Type                     string
	KafkaTopic               string
	FinancialInstitutionType string
	maxMemory                int64
	ImageFolderName          string
)

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

const (
	UpdateApp           = "Update your app"
	DeviceFound         = "Device Successfuly Found"
	OTPMessage          = "Your device lookup OTP is: %s. Valid for %d minutes."
	PhoneFound          = "Phone successfuly found"
	ImageUploadSuccess  = "image uploaded successfully"
	ProfileUploadSucess = "Profile theme set successfuly"
	PINResetSuccess     = "OTP verified for PIN reset"
	IncompleteUserInfo  = "incomplete user info"
	RoleBackReason      = "Roll back by the system because of the failed action"
)

// ---------------------------------------------------------------------------
// General-Purpose String Constants
// ---------------------------------------------------------------------------

const (
	Prelogin                 = "PRE_LOGIN"
	OTPLength                = 6
	DEV                      = "dev"
	UAT                      = "uat"
	Password                 = "PASSWORD"
	Login                    = "LOGIN"
	Change                   = "CHANGE"
	Permanent                = "PERMANENT"
	TokenType                = "TOKEN_TYPE"
	Token                    = "TOKEN"
	OTP                      = "OTP"
	VerifyOtp                = "VERIFY_OTP"
	DeviceLookUp             = "DEVICE_LOOKUP"
	Register                 = "REGISTER"
	Empty                    = ""
	SetPin                   = "SET_PIN"
	Pin                      = "PIN"
	OTPForRegistration       = "REGISTRATION"
	Incomplete               = "INCOMPLETE"
	ForgetPinVerifyOtp       = "FORGET_PIN_VERIFY_OTP"
	ResetPin                 = "RESET_PIN"
	Completed                = "COMPLETED"
	ProfileTemp              = "PROFILE-*.TMP"
	BucketUserProfilePicture = "USER-PROFILE-PICTURES"
	OtpExpirationTime        = 3 * time.Minute
	UssdMerchant             = "USSDM_"
	Cred                     = "Credential"
	Checker                  = "CHECKER"
	IFBChecker               = "IFB_CHECKER"
	Maker                    = "MAKER"
	IFBMaker                 = "IFBMAKER"
	ActionCode               = "action_code"
	Avatar                   = "avatar"
	DonationIcon             = "donation_icon"
	CampanyLogo              = "company_logo"
	DonationImage            = "donation_image"
	DonationCoverImage       = "donation_cover_image"
	BudgetCategoryIcon       = "budget_category_icon"
	ActionID                 = "action_id"
	ActionStatus             = "action_status"
	AuditorStatusDBFieldName = "auditor_status"
	Approved                 = "APPROVED"
	Rejected                 = "REJECTED"
	Reversed                 = "REVERSED"
	Canceled                 = "CANCELED"
	AccountNumber            = "account_number"
	PhoneNumber              = "phone_number"
	WithFayda                = "with_fayda"
	UPDATE                   = "UPDATE"
	ENABLE                   = "ENABLE"
	DISABLE                  = "DISABLE"
	DELETE                   = "DELETE"
	CREATE                   = "CREATE"
)

// ---------------------------------------------------------------------------
// Pagination Defaults
// ---------------------------------------------------------------------------

const (
	DefaultPage    = 1
	DefaultPerPage = 10
)

// ---------------------------------------------------------------------------
// BPS Topics
// ---------------------------------------------------------------------------

const (
	ContextKeyMetadata   ContextKey = "context_metadata"
	BPSApproveTopic                 = "bps.action.approve"
	BPSRejectTopic                  = "bps.action.reject"
	BPSAuditorClaimTopic            = "bps.action.auditor_claim"
	BPSAuditorMarkTopic             = "bps.action.auditor_mark"
)

// ---------------------------------------------------------------------------
// Currency & Branch
// ---------------------------------------------------------------------------

const (
	ETB CurrencyType = "ETB"
	USD CurrencyType = "USD"
)

const (
	IFB BranchType = "IFB"
	CB  BranchType = "CB"
)

// ---------------------------------------------------------------------------
// Platform
// ---------------------------------------------------------------------------

const (
	Android Platform = "ANDROID"
	Ios     Platform = "IOS"
)

// ---------------------------------------------------------------------------
// Vault Status
// ---------------------------------------------------------------------------

const (
	VaultStatusActive        VaultStatus = "ACTIVE"
	VaultStatusMatured       VaultStatus = "MATURED"
	VaultStatusUnlockedEarly VaultStatus = "UNLOCKED_EARLY"
	VaultStatusPaidOut       VaultStatus = "PAID_OUT"
	StatusOpen               VaultStatus = "OPEN"
	StatusClosed             VaultStatus = "CLOSED"
	StatusLocked             VaultStatus = "LOCKED"
	StatusActive             VaultStatus = "ACTIVE"
	StatusMatured            VaultStatus = "MATURED"
	StatusWithdrawn          VaultStatus = "WITHDRAWN"
	StatusEndedEarly         VaultStatus = "ENDED_EARLY"
	StatusDeleted            VaultStatus = "DELETED"
)

// ---------------------------------------------------------------------------
// Vault Type
// ---------------------------------------------------------------------------

const (
	VaultTypeBanking VaultType = "BANK_VAULT"
	VaultTypePrivate VaultType = "PRIVATE_VAULT"
	VaultTypeGroup   VaultType = "GROUP_VAULT"
)

// VaultCategory defines the category of the vault.
const (
	CategoryTravel  VaultCategory = "TRAVEL"
	CategoryGadgets VaultCategory = "GADGETS"
	CategoryEvents  VaultCategory = "EVENTS"
	CategoryGeneral VaultCategory = "GENERAL"
)

// ---------------------------------------------------------------------------
// Blocked On
// ---------------------------------------------------------------------------

const (
	NotBlocked BlockedOn = ""
	BPS        BlockedOn = "BPS"
	CPS        BlockedOn = "CPS"
)

// ---------------------------------------------------------------------------
// Notification
// ---------------------------------------------------------------------------

const (
	ForIFB NotificationFor = "IFB"
	ForCB  NotificationFor = "CB"
	ForAll NotificationFor = "ALL"
)

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusSeen    NotificationStatus = "SEEN"
)

// ---------------------------------------------------------------------------
// App View / Environment / App Type
// ---------------------------------------------------------------------------

const (
	AppViewTypeBoth AppViewType = "BOTH"
	AppViewTypeCB   AppViewType = "CB"
	AppViewTypeIFB  AppViewType = "IFB"
)

const (
	UatEnvironment        EnvironmentType = "UAT"
	DevEnvironment        EnvironmentType = "DEV"
	TestEnvironment       EnvironmentType = "TEST"
	ProductionEnvironment EnvironmentType = "PRODUCTION"
)

const (
	URL        AppType = "URL"
	UATApp     AppType = "UAT"
	Production AppType = "PRODUCATION"
	Test       AppType = "TEST"
	Dev        AppType = "DEV"
)

// ---------------------------------------------------------------------------
// Event Status
// ---------------------------------------------------------------------------

const (
	EventUpcomming EventStatus = "UPCOMMING"
	EventLive      EventStatus = "LIVE"
	EventClosed    EventStatus = "CLOSED"
)

// ---------------------------------------------------------------------------
// Restriction
// ---------------------------------------------------------------------------

const (
	AgeRestriction RestrictionType = "AGE_RESTRICTION"
)

// ---------------------------------------------------------------------------
// Accrual
// ---------------------------------------------------------------------------

const (
	AccrualFreqDaily     AccrualFrequency = "DAILY"
	AccrualFreqMonthly   AccrualFrequency = "MONTHLY"
	AccrualFreqQuarterly AccrualFrequency = "QUARTERLY"
	AccrualFreqAnnually  AccrualFrequency = "ANNUALLY"
)

const (
	AccrualMethodCompound AccrualMethod = "COMPOUND"
	AccrualMethodSimple   AccrualMethod = "SIMPLE"
)

// ---------------------------------------------------------------------------
// Method
// ---------------------------------------------------------------------------

const (
	OPEN      Method = "OPEN"
	PIN       Method = "PIN"
	OTPANDPIN Method = "OTP_PIN"
)

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

const (
	RegistrationTypeNew    RegistrationType = "NEW"
	RegistrationTypeLinked RegistrationType = "LINKED"
)

// ---------------------------------------------------------------------------
// Advert
// ---------------------------------------------------------------------------

const (
	IFB_ADVERT_FOR  AdvertFor = "IFB"
	CB_ADVERT_FOR   AdvertFor = "CB"
	BOTH_ADVERT_FOR AdvertFor = "ALL"
)

// ---------------------------------------------------------------------------
// Realm
// ---------------------------------------------------------------------------

const (
	ELST_REALM     Realm = "ELST"
	BANK_REALM     Realm = "BANK"
	DISTRICT_REALM Realm = "DISTRICT"
	BRANCH_REALM   Realm = "BRANCH"
	MERCHANT_REALM Realm = "MERCHANT"
	COMPANY_REALM  Realm = "COMPANY"
	MEMBER_REALM   Realm = "MEMBER"
)

// ---------------------------------------------------------------------------
// Fee Type
// ---------------------------------------------------------------------------

const (
	Flat    FeeType = "FLAT"
	Percent FeeType = "PERCENT"
)

// type DeviceType string

// const (
// 	Android DeviceType = "Android"
// 	IOS     DeviceType = "IOS"
// )

// ---------------------------------------------------------------------------
// Transaction Method
// ---------------------------------------------------------------------------

const (
	Free      TransactionMethod = "FREE"
	Otp       TransactionMethod = "OTP"
	OtpAndPin TransactionMethod = "OTP_PIN"
)

// ---------------------------------------------------------------------------
// Account Type
// ---------------------------------------------------------------------------

const (
	AccountTypeIFB AccountType = "IFB"
	AccountTypeCB  AccountType = "CB"
)

// ---------------------------------------------------------------------------
// Gender
// ---------------------------------------------------------------------------

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

// ---------------------------------------------------------------------------
// Vendor
// ---------------------------------------------------------------------------

const (
	ByVerigram Vendor = "VERIGRAM"
	ByFayda    Vendor = "FAYDA"
)

// ---------------------------------------------------------------------------
// KYC Status
// ---------------------------------------------------------------------------

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusApproved KYCStatus = "APPROVED"
	KYCStatusComplete KYCStatus = "COMPLETE"
	KYCStatusRejected KYCStatus = "REJECTED"
)

// ---------------------------------------------------------------------------
// BPS Status
// ---------------------------------------------------------------------------

const (
	BPSStatusAuthorized BPSStatus = "AUTHORIZED"
	BPSStatusDenied     BPSStatus = "DENIED"
	BPSStatusPending    BPSStatus = "PENDING"
	BPSStatusInitiated  BPSStatus = "INITIATED"
)

// ---------------------------------------------------------------------------
// Marital Status
// ---------------------------------------------------------------------------

const (
	Single   MaritalStatus = "SINGLE"
	Married  MaritalStatus = "MARRIED"
	Divorced MaritalStatus = "DIVORCED"
	Widow    MaritalStatus = "WIDOW"
)

// ---------------------------------------------------------------------------
// Device Status
// ---------------------------------------------------------------------------

const (
	Linked   DeviceStatus = "LINKED"
	UnLinked DeviceStatus = "UNLINKED"
)

// ---------------------------------------------------------------------------
// Onboarding Method
// ---------------------------------------------------------------------------

const (
	Ussd   OnboardingMethod = "USSD"
	Branch OnboardingMethod = "BRANCH"
	OldApp OnboardingMethod = "OLD_APP"
	Fayda  OnboardingMethod = "FAYDA"
)

// ---------------------------------------------------------------------------
// Enabled Channels
// ---------------------------------------------------------------------------

const (
	UssdChannel     EnabledChannels = "USSD"
	SuperappChannel EnabledChannels = "SUPERAPP"
)

// ---------------------------------------------------------------------------
// Account Status
// ---------------------------------------------------------------------------

const (
	Active   AccountStatus = "ACTIVE"
	InActive AccountStatus = "INACTIVE"
)

// ---------------------------------------------------------------------------
// Member Type
// ---------------------------------------------------------------------------

const (
	CBT  MemberType = "CB"
	IFBT MemberType = "IFB"
)

// ---------------------------------------------------------------------------
// OTP For
// ---------------------------------------------------------------------------

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
	OTPForForgetPin        OTPFor = "FORGET_PIN"
)

// ---------------------------------------------------------------------------
// OTP Status
// ---------------------------------------------------------------------------

const (
	Pending  OTPStatus = "PENDING"
	Verified OTPStatus = "VERIFIED"
	Denied   OTPStatus = "DENIED"
)

// ---------------------------------------------------------------------------
// Risk Level
// ---------------------------------------------------------------------------

const (
	None   RiskLevel = ""
	High   RiskLevel = "HIGH"
	Midium RiskLevel = "MIDIUM"
	Low    RiskLevel = "LOW"
)

// ---------------------------------------------------------------------------
// Customer Status
// ---------------------------------------------------------------------------

const (
	CustomerActive  CustomerStatus = "ACTIVE"
	CustomerPending CustomerStatus = "PENDING"
	CustomerExpired CustomerStatus = "EXPIRED"
)

const (
	ConventionalAccount AccountType = "CONVENTIONAL"
	CBENoor             AccountType = "CBENOOR"
)

const (
	MaritalMarried  MaritalStatus = "MARRIED"
	MaritalDivorced MaritalStatus = "DIVORCED"
	MaritalSingle   MaritalStatus = "SINGLE"
	MaritalWidowed  MaritalStatus = "WIDOWED"
)

// ---------------------------------------------------------------------------
// Employment Status
// ---------------------------------------------------------------------------

const (
	EmploymentStatusAgent      EmploymentStatus = "AGENT"
	EmploymentStatusEmployed   EmploymentStatus = "EMPLOYED"
	EmploymentStatusForeigner  EmploymentStatus = "FOREIGNER"
	EmploymentStatusMinor      EmploymentStatus = "MINOR"
	EmploymentStatusPensioner  EmploymentStatus = "PENSIONER"
	EmploymentStatusCBEStaff   EmploymentStatus = "STAFF_OF_CBE"
	EmploymentStatusUnemployed EmploymentStatus = "UNEMPLOYED"
	EmploymentStatusOther      EmploymentStatus = "OTHER_INDIVIDUALS"
)

// ---------------------------------------------------------------------------
// Auditor
// ---------------------------------------------------------------------------

const (
	MARKEDASRIGHT AuditorMark = "MARKEDASRIGHT"
	MARKEDASWRONG AuditorMark = "MARKEDASWRONG"
)

const (
	AUDITORNOTCHECKED AuditorStatus = "NOTCHECKED"
	AUDITORINPROGRESS AuditorStatus = "INPROGRESS"
	AUDITORCHECKED    AuditorStatus = "CHECKED"
)

// ---------------------------------------------------------------------------
// Settlement Method
// ---------------------------------------------------------------------------

const (
	SettlementMethodDirect       SettlementMethod = "DIRECT"
	SettlementMethodGL           SettlementMethod = "GL"
	SettlementMethodMultiAccount SettlementMethod = "MULTI_ACCOUNT"
)

// ---------------------------------------------------------------------------
// Level
// ---------------------------------------------------------------------------

const (
	ZERO Level = 0
	ONE  Level = 1
	TWO  Level = 2
)

// ---------------------------------------------------------------------------
// Action Type
// ---------------------------------------------------------------------------

const (
	ActionDelete ActionType = "DELETE"
	ActionUpdate ActionType = "UPDATE"
	ActionCreate ActionType = "CREATE"
)

// ---------------------------------------------------------------------------
// Request Actions
// ---------------------------------------------------------------------------

const (
	ActionPending  RequestAction = "PENDING"
	ActionApproved RequestAction = "APPROVED"
	ActionRejected RequestAction = "REJECTED"

	// Mini App Product Code
	RequestCreateMiniappProductCode  string = "CREATE_MINI_APP_PRODUCT_CODE"
	RequestUpdateMiniappProductCode  string = "UPDATE_MINI_APP_PRODUCT_CODE"
	RequestDeleteMiniappProductCode  string = "DELETE_MINI_APP_PRODUCT_CODE"
	RequestEnableMiniappProductCode  string = "ENABLE_MINI_APP_PRODUCT_CODE"
	RequestDisableMiniappProductCode string = "DISABLE_MINI_APP_PRODUCT_CODE"

	// Job Role
	RequestCreateJobRole  string = "CREATE_JOB_ROLE"
	RequestUpdateJobRole  string = "UPDATE_JOB_ROLE"
	RequestDeleteJobRole  string = "DELETE_JOB_ROLE"
	RequestEnableJobRole  string = "ENABLE_JOB_ROLE"
	RequestDisableJobRole string = "DISABLE_JOB_ROLE"

	// Role
	RequestCreateRole  string = "CREATE_ROLE"
	RequestUpdateRole  string = "UPDATE_ROLE"
	RequestDeleteRole  string = "DELETE_ROLE"
	RequestEnableRole  string = "ENABLE_ROLE"
	RequestDisableRole string = "DISABLE_ROLE"

	// USSD Merchant
	RequestCreateUssdMerchant  string = "CREATE_USSD_MERCHANT"
	RequestUpdateUssdMerchant  string = "UPDATE_USSD_MERCHANT"
	RequestDeleteUssdMerchant  string = "DELETE_USSD_MERCHANT"
	RequestEnableUssdMerchant  string = "ENABLE_USSD_MERCHANT"
	RequestDisableUssdMerchant string = "DISABLE_USSD_MERCHANT"

	// Ecommerce Merchant
	RequestCreateEcommerceMerchant  RequestAction = "CREATE_ECOMMERCE_MERCHANT"
	RequestUpdateEcommerceMerchant  RequestAction = "UPDATE_ECOMMERCE_MERCHANT"
	RequestEnableEcommerceMerchant  RequestAction = "ENABLE_ECOMMERCE_MERCHANT"
	RequestDisableEcommerceMerchant RequestAction = "DISABLE_ECOMMERCE_MERCHANT"
	RequestDeleteEcommerceMerchant  RequestAction = "DELETE_ECOMMERCE_MERCHANT"

	// Action Role
	RequestUser                 RequestAction = "USER"
	RequestActionRole           RequestAction = "ACTION_ROLE"
	RequestCreateActionRole     RequestAction = "CREATE_ACTION_ROLE"
	RequestUpdateActionRole     RequestAction = "UPDATE_ACTION_ROLE"
	RequestEnableActionRole     RequestAction = "ENABLE_ACTION_ROLE"
	RequestDisableActionRole    RequestAction = "DISABLE_ACTION_ROLE"
	RequestDeleteActionRole     RequestAction = "DELETE_ACTION_ROLE"
	RequestCreateCpsActionRole  RequestAction = "CREATE_CPS_ACTION_ROLE"
	RequestUpdateCpsActionRole  RequestAction = "UPDATE_CPS_ACTION_ROLE"
	RequestEnableCpsActionRole  RequestAction = "ENABLE_CPS_ACTION_ROLE"
	RequestDisableCpsActionRole RequestAction = "DISABLE_CPS_ACTION_ROLE"
	RequestDeleteCpsActionRole  RequestAction = "DELETE_CPS_ACTION_ROLE"

	// CPS User
	RequestCpsUserCreate  RequestAction = "CREATE_CPS_USER"
	RequestCpsUserUpdate  RequestAction = "UPDATE_CPS_USER"
	RequestCpsUserDelete  RequestAction = "DELETE_CPS_USER"
	RequestCpsUserEnable  RequestAction = "ENABLE_CPS_USER"
	RequestCpsUserDisable RequestAction = "DISABLE_CPS_USER"

	// KYC
	RequestUpdateKYC  RequestAction = "UPDATE_KYC"
	RequestApproveKYC RequestAction = "APPROVE_KYC"

	// Permission Group
	RequestPermissionGroup       RequestAction = "PERMISSION_GROUP"
	RequestCreatePermissionGroup RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup RequestAction = "DELETE_PERMISSION_GROUP"

	// Department
	RequestDepartment              RequestAction = "DEPARMTENT"
	RequestCreateDepartment        RequestAction = "CREATE_DEPARTMENT"
	RequestUpdateDepartment        RequestAction = "UPDATE_DEPARTMENT"
	RequestDeleteDepartment        RequestAction = "DELETE_DEPARTMENT"
	RequestEnableDepartment        RequestAction = "ENABLE_DEPARTMENT"
	RequestDisableDepartment       RequestAction = "DISABLE_DEPARTMENT"
	RequestEnableDisableDepartment RequestAction = "ENABLE_DISABLE_DEPARTMENT"

	// User
	RequestEnableUser  RequestAction = "ENABLE_USER"
	RequestDisableUser RequestAction = "DISABLE_USER"
	RequestUpdateUser  RequestAction = "UPDATE_USER"
	RequestArchiveUser RequestAction = "ARCHIVE_USER"
	RequestUnlinkUser  RequestAction = "UNLINK_USER"

	// BPS User
	RequestBPSUser RequestAction = "BPS_USER"

	RequestBpsUserCreate  RequestAction = "CREATE_BPS_USER"
	RequestBpsUserDelete  RequestAction = "DELETE_BPS_USER"
	RequestBpsUserEnable  RequestAction = "ENABLE_BPS_USER"
	RequestBpsUserDisable RequestAction = "DISABLE_BPS_USER"

	RequestBpsUserUpdate  RequestAction = "UPDATE_BPS_USER"
	RequestCreateBPSUser  RequestAction = "CREATE_BPS_USER"
	RequestDisableBPSUser RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser  RequestAction = "ENABLE_BPS_USER"

	// Limits & Auth
	RequestTotalDailyLimit       RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT             RequestAction = "UPDATE_VAT"
	RequestAuthTier              RequestAction = "AUTHTIER"
	RequestDeleteAmountBasedAuth RequestAction = "DELETE_AMOUNT_BASED_AUTH"
	RequestCreateAmountBasedAuth RequestAction = "CREATE_AMOUNT_BASED_AUTH"
	RequestUpdateAmountBasedAuth RequestAction = "UPDATE_AMOUNT_BASED_AUTH"

	// Advert
	RequestCreateAdvert  RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert  RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert  RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert  RequestAction = "DELETE_ADVERT"

	// Bank
	RequestCreateBank        RequestAction = "CREATE_BANK"
	RequestUpdateBank        RequestAction = "UPDATE_BANK"
	RequestUpdateBankLogo    RequestAction = "UPDATE_BANK_LOGO"
	RequestDeleteBank        RequestAction = "DELETE_BANK"
	RequestEnableDisableBank RequestAction = "ENABLE_DISABLE_BANK"
	RequestEnableBank        RequestAction = "ENABLE_BANK"
	RequestDisableBank       RequestAction = "DISABLE_BANK"

	// Device Version
	RequestCreateDeviceVersion        RequestAction = "CREATE_DEVICE_VERSION"
	RequestUpdateDeviceVersion        RequestAction = "UPDATE_DEVICE_VERSION"
	RequestEnableDeviceVersion        RequestAction = "ENABLE_DEVICE_VERSION"
	RequestDisableDeviceVersion       RequestAction = "DISABLE_DEVICE_VERSION"
	RequestDeleteDeviceVersion        RequestAction = "DELETE_DEVICE_VERSION"
	RequestEnableDisableDeviceVersion RequestAction = "ENABLE_DISABLE_DEVICE_VERSION"

	// Customer
	RequestEnableDisableCustomer RequestAction = "ENABLE_DISABLE_CUSTOMER"
	RequestApproveFaydaCustomer  RequestAction = "APPROVE_FAYDA_CUSTOMER"
	RequestDisableFaydaAccount   RequestAction = "DISABLE_FAYDA_ACCOUNT"
	RequestEnableFaydaAccount    RequestAction = "ENABLE_FAYDA_ACCOUNT"

	// Wallet
	RequestCreateWallet  RequestAction = "CREATE_WALLET"
	RequestUpdateWallet  RequestAction = "UPDATE_WALLET"
	RequestDeleteWallet  RequestAction = "DELETE_WALLET"
	RequestEnableWallet  RequestAction = "ENABLE_WALLET"
	RequestDisableWallet RequestAction = "DISABLE_WALLET"

	// Services catalog (model.Services)
	RequestCreateServiceList RequestAction = "CREATE_SERVICE_LIST"
	RequestUpdateServiceList RequestAction = "UPDATE_SERVICE_LIST"
	// Service
	RequestCreateService        RequestAction = "CREATE_SERVICE"
	RequestUpdateService        RequestAction = "UPDATE_SERVICE"
	RequestEnableService        RequestAction = "ENABLE_SERVICE"
	RequestDisableService       RequestAction = "DISABLE_SERVICE"
	RequestUpdateServiceSingle  RequestAction = "UPDATE_SERVICE_SINGLE_CAP"
	RequestUpdateServiceTotal   RequestAction = "UPDATE_SERVICE_TOTAL_CAP"
	RequestUpdateServiceMinCap  RequestAction = "UPDATE_SERVICE_MIN_CAP"
	RequestUpdateServiceDetails RequestAction = "UPDATE_SERVICE_DETAILS"
	RequestUpdateServiceRule    RequestAction = "UPDATE_SERVICE_RULE"
	RequestUpdateMinimumService RequestAction = "UPDATE_MINIMUM_SERVICE"
	RequestBulkServiceEnable    RequestAction = "ENABLE_BULK_SERVICE"
	RequestBulkServiceDisable   RequestAction = "DISABLE_BULK_SERVICE"

	// Service Fee
	RequestCreateServiceFee RequestAction = "CREATE_SERVICE_FEE"
	RequestUpdateServiceFee RequestAction = "UPDATE_SERVICE_FEE"
	RequestDeleteServiceFee RequestAction = "DELETE_SERVICE_FEE"

	// Daily Limit
	RequestCreateDailyLimit RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit RequestAction = "DELETE DAILY LIMIT"

	// Validation
	RequestCreateValidation RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation RequestAction = "DELETE_VALIDATION"

	// Password & Access
	RequestUpdatePasswordExpiry RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreatePasswordRule   RequestAction = "CREATE_PASSWORD_RULE"
	RequestUpdatePasswordRule   RequestAction = "UPDATE_PASSWORD_RULE"
	RequestUpdateArchiveExpiry  RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestUpdateAccessConfig   RequestAction = "UPDATE_ACCESS_CONFIG"
	RequestUpdateBlockTime      RequestAction = "UPDATE_BLOCK_TIME"
	RequestUpdateHQBlockTime    RequestAction = "UPDATE_HQ_BLOCK_TIME"
	RequestUpdateHQArchiveTime  RequestAction = "UPDATE_HQ_ARCHIVE_TIME"
	RequestUpdateTotal          RequestAction = "UPDATE_TOTAL"

	// Budget
	RequestBudgetColor           RequestAction = "BUDGET_COLOR"
	RequestBudgetIcon            RequestAction = "BUDGET_ICON"
	RequestCreateBudgetCategory  RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestUpdateBudgetCategory  RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestDeleteBudgetCategory  RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestEnableBudgetCategory  RequestAction = "ENABLE_BUDGET_CATEGORY"
	RequestDisableBudgetCategory RequestAction = "DISABLE_BUDGET_CATEGORY"

	// Product
	RequestUpdateProduct     RequestAction = "UPDATE_PRODUCT"
	RequestUpdateProductCode RequestAction = "UPDATE_PRODUCT_CODE"

	// Notification
	RequestCreatePublicNotification RequestAction = "CREATE_PUBLIC_NOTIFICATION"
	RequestUpdatePublicNotification RequestAction = "UPDATE_PUBLIC_NOTIFICATION"
	RequestDeleteNotification       RequestAction = "DELETE_NOTIFICATION"
	RequestEnableNotification       RequestAction = "ENABLE_NOTIFICATION"
	RequestDisableNotification      RequestAction = "DISABLE_NOTIFICATION"
	RequestMarkNotificationAsSeen   RequestAction = "MARK_NOTIFICATION_AS_SEEN"

	// Account Validation
	RequestUpdateAccountValidation RequestAction = "UPDATE_ACCOUNT_VALIDATION"

	// Branch Enable/Disable
	RequestEnableSingleBranch    RequestAction = "ENABLE_SINGLE_BRANCH"
	RequestDisableSingleBranch   RequestAction = "DISABLE_SINGLE_BRANCH"
	RequestEnableMultiUsers      RequestAction = "ENABLE_MULTI_USERS"
	RequestDisableMultiUsers     RequestAction = "DISABLE_MULTI_USERS"
	RequestEnableSingleBranches  RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches   RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches  RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"
	RequestEnableBranches        RequestAction = "REQUEST_ENABLE_BRANCHES"
	RequestDisableBranches       RequestAction = "REQUEST_DISABLE_BRANCHES"

	// Region
	RequestEnableRegions  RequestAction = "REQUEST_ENABLE_REGIONS"
	RequestDisableRegions RequestAction = "REQUEST_DISABLE_REGIONS"
	RequestBlockRegion    RequestAction = "BLOCK_REGION"

	// District
	RequestEnableDistricts  RequestAction = "REQUEST_ENABLE_DISTRICTS"
	RequestDisableDistricts RequestAction = "REQUEST_DISABLE_DISTRICTS"
	RequestBlockDistrict    RequestAction = "BLOCK_DISTRICT"

	// City
	RequestEnableCities  RequestAction = "REQUEST_ENABLE_CITIES"
	RequestDisableCities RequestAction = "REQUEST_DISABLE_CITIES"
	RequestBlockCity     RequestAction = "BLOCK_CITY"
	RequestBlockUser     RequestAction = "BLOCK_USER"

	// Business
	RequestCreateBusiness RequestAction = "CREATE_BUSINESS"
	RequestUpdateBusiness RequestAction = "UPDATE_BUSINESS"

	// Mini App Merchant
	RequestCreateMiniAppMerchant  RequestAction = "CREATE_MINI_APP_MERCHANT"
	RequestUpdateMiniAppMerchant  RequestAction = "UPDATE_MINI_APP_MERCHANT"
	RequestEnableMiniAppMerchant  RequestAction = "ENABLE_MINI_APP_MERCHANT"
	RequestDeleteMiniAppMerchant  RequestAction = "DELETE_MINI_APP_MERCHANT"
	RequestDisableMiniAppMerchant RequestAction = "DISABLE_MINI_APP_MERCHANT"

	// Mini App
	RequestCreateMiniApp  RequestAction = "CREATE_MINI_APP"
	RequestUpdateMiniApp  RequestAction = "UPDATE_MINI_APP"
	RequestDeleteMiniApp  RequestAction = "DELETE_MINI_APP"
	RequestEnableMiniApp  RequestAction = "ENABLE_MINI_APP"
	RequestDisableMiniApp RequestAction = "DISABLE_MINI_APP"

	// Mini App Category
	RequestCreateMiniAppCategory  RequestAction = "CREATE_MINI_APP_CATEGORY"
	RequestUpdateMiniAppCategory  RequestAction = "UPDATE_MINI_APP_CATEGORY"
	RequestDeleteMiniAppCategory  RequestAction = "DELETE_MINI_APP_CATEGORY"
	RequestEnableMiniAppCategory  RequestAction = "ENABLE_MINI_APP_CATEGORY"
	RequestDisableMiniAppCategory RequestAction = "DISABLE_MINI_APP_CATEGORY"

	// Avatar
	RequestCreateAvatar  RequestAction = "CREATE_AVATAR"
	RequestUpdateAvatar  RequestAction = "UPDATE_AVATAR"
	RequestDeleteAvatar  RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar  RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar RequestAction = "DISABLE_AVATAR"

	// Bank Vault
	RequestCreateBankVault  RequestAction = "CREATE_VAULT_BANK"
	RequestUpdateBankVault  RequestAction = "UPDATE_VAULT_BANK"
	RequestDeleteBankVault  RequestAction = "DELETE_VAULT_BANK"
	RequestEnableBankVault  RequestAction = "ENABLE_VAULT_BANK"
	RequestDisAbleBankVault RequestAction = "DISABLE_VAULT_BANK"

	// Vault Amount Tier
	RequestCreateVaultAmountTier  RequestAction = "CREATE_VAULT_AMOUNT_TIER"
	RequestUpdateVaultAmountTier  RequestAction = "UPDATE_VAULT_AMOUNT_TIER"
	RequestDeleteVaultAmountTier  RequestAction = "DELETE_VAULT_AMOUNT_TIER"
	RequestEnableVaultAmountTier  RequestAction = "ENABLE_VAULT_AMOUNT_TIER"
	RequestDisAbleVaultAmountTier RequestAction = "DISABLE_VAULT_AMOUNT_TIER"

	// Vault Category
	RequestCreateVaultCategory  RequestAction = "CREATE_VAULT_CATEGORY"
	RequestUpdateVaultCategory  RequestAction = "UPDATE_VAULT_CATEGORY"
	RequestDeleteVaultCategory  RequestAction = "DELETE_VAULT_CATEGORY"
	RequestEnableVaultCategory  RequestAction = "ENABLE_VAULT_CATEGORY"
	RequestDisAbleVaultCategory RequestAction = "DISABLE_VAULT_CATEGORY"

	// Vault Withdrawal
	RequestCreateWithdrawal RequestAction = "CREATE_LOCKED_VAULT_WITHDRAWAL"
	RequestUpdateWithdrawal RequestAction = "UPDATE_LOCKED_VAULT_WITHDRAWAL"

	// Event
	RequestCreateEvent         RequestAction = "CREATE_EVENT"
	RequestUpdateEvent         RequestAction = "UPDATE_EVENT"
	RequestDeleteEvent         RequestAction = "DELETE_EVENT"
	RequestEnableEvent         RequestAction = "ENABLE_EVENT"
	RequestDisableEvent        RequestAction = "DISABLE_EVENT"
	RequestCreateEventCategory RequestAction = "CREATE_EVENT_CATEGORY"
	RequestUpdateEventCategory RequestAction = "UPDATE_EVENT_CATEGORY"

	// Event Merchant
	RequestCreateEventMerchant  RequestAction = "CREATE_EVENT_MERCHANT"
	RequestUpdateEventMerchant  RequestAction = "UPDATE_EVENT_MERCHANT"
	RequestDeleteEventMerchant  RequestAction = "DELETE_EVENT_MERCHANT"
	RequestEnableEventMerchant  RequestAction = "ENABLE_EVENT_MERCHANT"
	RequestDisableEventMerchant RequestAction = "DISABLE_EVENT_MERCHANT"

	// Logistics Merchant
	RequestCreateLogisticsMerchant  RequestAction = "CREATE_LOGISTICS_MERCHANT"
	RequestUpdateLogisticsMerchant  RequestAction = "UPDATE_LOGISTICS_MERCHANT"
	RequestDeleteLogisticsMerchant  RequestAction = "DELETE_LOGISTICS_MERCHANT"
	RequestEnableLogisticsMerchant  RequestAction = "ENABLE_LOGISTICS_MERCHANT"
	RequestDisableLogisticsMerchant RequestAction = "DISABLE_LOGISTICS_MERCHANT"

	// Customer Segmentation
	RequestCreateCustomerSegmentation  RequestAction = "CREATE_CUSTOMER_SEGMENTATION"
	RequestUpdateCustomerSegmentation  RequestAction = "UPDATE_CUSTOMER_SEGMENTATION"
	RequestEnableCustomerSegmentation  RequestAction = "ENABLE_CUSTOMER_SEGMENTATION"
	RequestDisableCustomerSegmentation RequestAction = "DISABLE_CUSTOMER_SEGMENTATION"
	RequestDeleteCustomerSegmentation  RequestAction = "DELETE_CUSTOMER_SEGMENTATION"

	// Access List Segmentation
	RequestCreateAccessListSegmentation        RequestAction = "CREATE_ACCESS_LIST_SEGMENTATION"
	RequestUpdateAccessListSegmentation        RequestAction = "UPDATE_ACCESS_LIST_SEGMENTATION"
	RequestEnableDisableAccessListSegmentation RequestAction = "ENABLE_DISABLE_ACCESS_LIST_SEGMENTATION"

	// CPS Role
	RequestCreateCpsRole  RequestAction = "CREATE_CPS_ROLE"
	RequestUpdateCpsRole  RequestAction = "UPDATE_CPS_ROLE"
	RequestDeleteCpsRole  RequestAction = "DELETE_CPS_ROLE"
	RequestEnableCpsRole  RequestAction = "ENABLE_CPS_ROLE"
	RequestDisableCpsRole RequestAction = "DISABLE_CPS_ROLE"

	RequestUpdateEevent      RequestAction = "UPDATE_EVENT"
	RequestCIFRemove         RequestAction = "CIF_REMOVE"
	RequestServiceFlagUpdate RequestAction = "SERVICE_FLAG_UPDATE"

	// Donation Category
	RequestCreateDonationCategory  RequestAction = "CREATE_DONATION_CATEGORY"
	RequestUpdateDonationCategory  RequestAction = "UPDATE_DONATION_CATEGORY"
	RequestEnableDonationCategory  RequestAction = "ENABLE_DONATION_CATEGORY"
	RequestDisableDonationCategory RequestAction = "DISABLE_DONATION_CATEGORY"

	// Donation Company
	RequestCreateDonationCompany  RequestAction = "CREATE_DONATION_COMPANY"
	RequestUpdateDonationCompany  RequestAction = "UPDATE_DONATION_COMPANY"
	RequestEnableDonationCompany  RequestAction = "ENABLE_DONATION_COMPANY"
	RequestDisableDonationCompany RequestAction = "DISABLE_DONATION_COMPANY"

	// Donation
	RequestCreateDonation      RequestAction = "CREATE_DONATION"
	RequestUpdateDonation      RequestAction = "UPDATE_DONATION"
	RequestUpdateDonationImage RequestAction = "UPDATE_DONATION_IMAGE"
	RequestDeleteDonationImage RequestAction = "DELETE_DONATION_IMAGE"
	RequestAddDonationImage    RequestAction = "ADD_DONATION_IMAGE"
	RequestEnableDonation      RequestAction = "ENABLE_DONATION"
	RequestDisableDonation     RequestAction = "DISABLE_DONATION"

	// Article
	RequestCreateArticle  RequestAction = "CREATE_ARTICLE"
	RequestUpdateArticle  RequestAction = "UPDATE_ARTICLE"
	RequestEnableArticle  RequestAction = "ENABLE_ARTICLE"
	RequestDisableArticle RequestAction = "DISABLE_ARTICLE"
	RequestDeleteArticle  RequestAction = "DELETE_ARTICLE"

	// Article Category
	RequestCreateArticleCategory  RequestAction = "CREATE_ARTICLE_CATEGORY"
	RequestUpdateArticleCategory  RequestAction = "UPDATE_ARTICLE_CATEGORY"
	RequestEnableArticleCategory  RequestAction = "ENABLE_ARTICLE_CATEGORY"
	RequestDisableArticleCategory RequestAction = "DISABLE_ARTICLE_CATEGORY"
	RequestDeleteArticleCategory  RequestAction = "DELETE_ARTICLE_CATEGORY"

	// News Tag
	RequestCreateNewsTag  RequestAction = "CREATE_NEWS_TAG"
	RequestUpdateNewsTag  RequestAction = "UPDATE_NEWS_TAG"
	RequestEnableNewsTag  RequestAction = "ENABLE_NEWS_TAG"
	RequestDisableNewsTag RequestAction = "DISABLE_NEWS_TAG"
	RequestDeleteNewsTag  RequestAction = "DELETE_NEWS_TAG"

	// News Category
	RequestCreateNewsCategory RequestAction = "CREATE_NEWS_CATEGORY"
	RequestUpdateNewsCategory RequestAction = "UPDATE_NEWS_CATEGORY"
	RequestDeleteNewsCategory RequestAction = "DELETE_NEWS_CATEGORY"

	// Topup
	RequestCreateTopup  RequestAction = "CREATE_TOPUP"
	RequestUpdateTopup  RequestAction = "UPDATE_TOPUP"
	RequestDeleteTopup  RequestAction = "DELETE_TOPUP"
	RequestEnableTopup  RequestAction = "ENABLE_TOPUP"
	RequestDisableTopup RequestAction = "DISABLE_TOPUP"

	// Short Video
	RequestCreateShortVideo  string = "CREATE_SHORT_VIDEO"
	RequestUpdateShortVideo  string = "UPDATE_SHORT_VIDEO"
	RequestEnableShortVideo  string = "ENABLE_SHORT_VIDEO"
	RequestDisableShortVideo string = "DISABLE_SHORT_VIDEO"
	RequestDeleteShortVideo  string = "DELETE_SHORT_VIDEO"

	// Customer KYC
	RequestCreateCustomerKYC RequestAction = "CREATE_CUSTOMER_KYC"
	RequestUpdateCustomerKYC RequestAction = "UPDATE_CUSTOMER_KYC"
	RequestDeleteCustomerKYC RequestAction = "DELETE_CUSTOMER_KYC"
)

const (
	UpdateAction ActionType = "UPDATE"
	CreateAction ActionType = "CREATE"
	DeleteAction ActionType = "DELETE"
)

// ---------------------------------------------------------------------------
// Primary Authentication
// ---------------------------------------------------------------------------

const (
	PrimaryAuthenticationByPhoneNumber         PrimaryAuthentication = "PHONE_NUMBER"
	PrimaryAuthenticationByEmail               PrimaryAuthentication = "EMAIL"
	PrimaryAuthenticationByEmailAndPhoneNumber PrimaryAuthentication = "EMAIL_AND_PHONE_NUMBER"
)

// ---------------------------------------------------------------------------
// Stage
// ---------------------------------------------------------------------------

const (
	StageUat Stage = "UAT"
)

// ---------------------------------------------------------------------------
// Invite Status
// ---------------------------------------------------------------------------

const (
	InvitePending   InviteStatus = "PENDING"
	InviteAccepted  InviteStatus = "ACCEPTED"
	InviteRejected  InviteStatus = "REJECTED"
	InviteWithdrawn InviteStatus = "WITHDRAWN"
	InviteExpired   InviteStatus = "EXPIRED"
)

// ---------------------------------------------------------------------------
// Member Status & Role
// ---------------------------------------------------------------------------

const (
	StatusPaid    MemberStatus = "PAID"
	StatusWaiting MemberStatus = "WAITING"
	StatusOverdue MemberStatus = "OVERDUE"
)

const (
	RoleAdmin       MemberRole = "ADMIN"
	RoleParticipant MemberRole = "PARTICIPANT"
)

// Recurrence defines the frequency of recurring payments.
const (
	RecurrenceWeekly   Recurrence = "WEEKLY"
	RecurrenceBiweekly Recurrence = "BIWEEKLY"
	RecurrenceMonthly  Recurrence = "MONTHLY"
)

// ContributionType defines the type of contribution made by a member.
const (
	ContributionTypeDeposit    ContributionType = "DEPOSIT"
	ContributionTypeWithdrawal ContributionType = "WITHDRAWAL"
	// ContributionTypeFee        ContributionType = "FEE"
	// ContributionTypeDisbursement ContributionType = "DISBURSEMENT"
)

// ---------------------------------------------------------------------------
// Type (Vault Transaction Type)
// ---------------------------------------------------------------------------

const (
	Deposit    Type = "DEPOSIT"
	Withdrawal Type = "WITHDRAWAL"
	AutoFund   Type = "AUTO_FUND"
	Reminder   Type = "REMINDER"
	Fee        Type = "FEE"
	Creation   Type = "CREATION"
	Fund       Type = "FUND"
	Update     Type = "UPDATE"
	Deletion   Type = "DELETION"
	End        Type = "END"
)

// ---------------------------------------------------------------------------
// Kafka Topics
// ---------------------------------------------------------------------------

const (
	ClientOrchestrationMemberTopic               KafkaTopic = "member.sync.cps"
	ClientOrchestrationKycTopic                  KafkaTopic = "kyc.sync.cps"
	ClientOrchestrationAccessControlTopic        KafkaTopic = "access_control.sync.cps"
	BulkServiceTopic                             KafkaTopic = "bulk_service"
	ClientOrchestrationLinkedAccountTopic        KafkaTopic = "linked_account.sync.cps"
	ClientOrchestrationAccountBlockTopic         KafkaTopic = "account_block.sync.cps"
	ClientOrchestrationAccountValidationTopic    KafkaTopic = "account_validation.sync.cps"
	ClientOrchestrationDeviceVersionControlTopic KafkaTopic = "device_version_control.sync.cps"
	DeviceVersionControlTopic                    KafkaTopic = "device_version_control"
	ClientOrchestrationDonationTopic             KafkaTopic = "donation.sync.cps"
	ClientOrchestrationDonationCategoryTopic     KafkaTopic = "donation_category.sync.cps"
	ClientOrchestrationDonationCompanyTopic      KafkaTopic = "donation_company.sync.cps"
	ClientOrchestrationBudgetCategoryTopic       KafkaTopic = "budget_category.sync.cps"
	ClientOrchestrationArticleTopic              KafkaTopic = "news_article.sync.cps"
	ClientOrchestrationShortVideoTopic           KafkaTopic = "short_video.sync.cps"
	ClientOrchestrationNewsCategoryTopic         KafkaTopic = "news_category.sync.cps"
	ClientOrchestrationNewsTagTopic              KafkaTopic = "news_tag.sync.cps"
	ClientOrchestrationNotificationTopic         KafkaTopic = "notification.sync.cps"
	ClientOrchestrationServicesTopic             KafkaTopic = "services.sync.cps"
	AccessListSegmentationTopic                  KafkaTopic = "customer_segmentation"
	InAppNotificationsTopic                      KafkaTopic = "inapp.notifications"
	CustomerRoleUpdatedTopic                     KafkaTopic = "cps-customer-role-updated"
	CustomerSegmentationUpdatedTopic             KafkaTopic = "cps-customer-segmentation-updated"
)

// ---------------------------------------------------------------------------
// Financial Institution
// ---------------------------------------------------------------------------

const (
	Bank   FinancialInstitutionType = "BANK"
	Wallet FinancialInstitutionType = "WALLET"
	MFI    FinancialInstitutionType = "MFI"
)

// ---------------------------------------------------------------------------
// Upload
// ---------------------------------------------------------------------------

const (
	MaxMemoryForUpload maxMemory = 15 << 20 // 15 MB
)

// ---------------------------------------------------------------------------
// Image Folder Names
// ---------------------------------------------------------------------------

const (
	BankFolderName             ImageFolderName = "banks"
	UssdMerchantFolderName     ImageFolderName = "ussd_merchant"
	WalletFolderName           ImageFolderName = "wallets"
	AdFolderName               ImageFolderName = "ads"
	AvatarFolderName           ImageFolderName = "avatars"
	BudgetCategoryFolderName   ImageFolderName = "budget_categories"
	DonationFolderName         ImageFolderName = "donations"
	DonationCategoryFolderName ImageFolderName = "donation_categories"
	DonationCompanyFolderName  ImageFolderName = "donation_companies"
	EventFolderName            ImageFolderName = "events"
	TopupFolderName            ImageFolderName = "topups"
	VaultCategoryFolderName    ImageFolderName = "vault_categories"
	CustomerKYCFolderName      ImageFolderName = "customer_kyc"
)

// ---------------------------------------------------------------------------
// Redis Key Prefixes
// ---------------------------------------------------------------------------

const (
	RedisCPSUserDeviceIDPrefix = "cps:auth:device"
)

// ---------------------------------------------------------------------------
// Parser Functions
// ---------------------------------------------------------------------------

func ParseVaultStatus(input string) (VaultStatus, error) {
	switch strings.ToUpper(input) {
	case "OPEN":
		return StatusOpen, nil
	case "CLOSED":
		return StatusClosed, nil
	case "LOCKED":
		return StatusLocked, nil
	default:
		return "", fmt.Errorf("invalid VaultStatus: %s", input)
	}
}

func ParseVaultType(input string) (VaultType, error) {
	switch strings.ToUpper(input) {
	case "BANK_VAULT":
		return VaultTypeBanking, nil
	case "PRIVATE_VAULT":
		return VaultTypePrivate, nil
	case "GROUP_VAULT":
		return VaultTypeGroup, nil
	default:
		return "", fmt.Errorf("invalid VaultType: %s", input)
	}
}

func ParseVaultCategory(input string) (VaultCategory, error) {
	switch strings.ToUpper(input) {
	case "TRAVEL":
		return CategoryTravel, nil
	case "GADGETS":
		return CategoryGadgets, nil
	case "EVENTS":
		return CategoryEvents, nil
	case "GENERAL":
		return CategoryGeneral, nil
	default:
		return "", fmt.Errorf("invalid VaultCategory: %s", input)
	}
}
