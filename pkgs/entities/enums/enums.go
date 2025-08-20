package enums

type BPSStatus string
type MemberType string
type MartialStatus string
type Realm string
type Gender string
type RegistrationType string
type OTPFor string
type KYCStatus string
type LDAPStatus string
type PrimaryAuthentication string
type DeviceStatus string
type OTPStatus string
type MerchantRole string
type PoolSource string
type AccountBranchType string
type AccountType string
type AccountStatus string
type KYCLevel string
type RequestAction string
type ActionType string
type ActionStatus string
type UserType string
type EventStatus string
type ErrorType string
type HealthStatus string

const (
	HEALTH_STATUS_HEALTHY   HealthStatus = "HEALTHY"
	HEALTH_STATUS_UNHEALTHY HealthStatus = "UNHEALTHY"
	HEALTH_STATUS_DEGRADED  HealthStatus = "DEGRADED"
)
const (
	ERROR_TYPE_VALIDATION   ErrorType = "VALIDATION_ERROR"
	ERROR_TYPE_NOT_FOUND    ErrorType = "NOT_FOUND"
	ERROR_TYPE_UNAUTHORIZED ErrorType = "UNAUTHORIZED"
	ERROR_TYPE_FORBIDDEN    ErrorType = "FORBIDDEN"
	ERROR_TYPE_CONFLICT     ErrorType = "CONFLICT"
	ERROR_TYPE_INTERNAL     ErrorType = "INTERNAL_ERROR"
	ERROR_TYPE_BAD_REQUEST  ErrorType = "BAD_REQUEST"
	ERROR_TYPE_TIMEOUT      ErrorType = "TIMEOUT"
)
const (
	EVENT_UPCOMMING EventStatus = "UPCOMMING"
	EVENT_LIVE      EventStatus = "LIVE"
	EVENT_CLOSED    EventStatus = "CLOSED"
)

// RestrictionType represents the type of restriction for an event
// (from internal/domain/event/entity.go)
type RestrictionType string

const (
	AGE_RESTRICTION RestrictionType = "AGE_RESTRICTION"
)

// TicketStatus represents the status of a ticket
// (from internal/domain/ticket/entity.go)
type TicketStatus string

const (
	TICKET_BOOKED   TicketStatus = "BOOKED"
	TICKET_EXPIRED  TicketStatus = "EXPIRED"
	TICKET_PAID     TicketStatus = "PAID"
	TICKET_REDEEMED TicketStatus = "REDEEMED"
)

// AdvertFor represents the target audience for an advert
// (from internal/domain/ad/entity.go)
type AdvertFor string

const (
	IFB_ADVERT_FOR  AdvertFor = "IFB"
	CB_ADVERT_FOR   AdvertFor = "CB"
	BOTH_ADVERT_FOR AdvertFor = "ALL"
)

// EnvironmentType represents the environment type for miniapp
// (from internal/domain/miniapp/entity.go)
type EnvironmentType string

const (
	UAT_ENVIRONMENT        EnvironmentType = "UAT"
	DEV_ENVIRONMENT        EnvironmentType = "DEV"
	TEST_ENVIRONMENT       EnvironmentType = "TEST"
	PRODUCTION_ENVIRONMENT EnvironmentType = "PRODUCTION"
)

// BranchType represents the branch type for miniapp
// (from internal/domain/miniapp/entity.go)
type BranchType string

const (
	IFB_BRANCH_TYPE BranchType = "IFB"
	CB_BRANCH_TYPE  BranchType = "CB"
)

const (
	USER_TYPE_MAKER   UserType = "MAKER"
	USER_TYPE_CHECKER UserType = "CHECKER"
)

const (
	ACTION_PENDING  ActionStatus = "PENDING"
	ACTION_APPROVED ActionStatus = "APPROVED"
	ACTION_REJECTED ActionStatus = "REJECTED"
)

const (
	ACTION_CREATE ActionType = "CREATE"
	ACTION_UPDATE ActionType = "UPDATE"
	ACTION_DELETE ActionType = "DELETE"
)

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
const (
	KYC_LEVEL_ZERO KYCLevel = "ZERO"
	KYC_LEVEL_ONE  KYCLevel = "ONE"
	KYC_LEVEL_TWO  KYCLevel = "TWO"
)

const (
	BPS_STATUS_AUTHORIZED BPSStatus = "AUTHORIZED"
	BPS_STATUS_DENIED     BPSStatus = "DENIED"
	BPS_STATUS_PENDING    BPSStatus = "PENDING"
	BPS_STATUS_INITIATED  BPSStatus = "INITIATED"
)

const (
	MEMBER_TYPE_CB  MemberType = "CB"
	MEMBER_TYPE_IFB MemberType = "IFB"
)

const (
	MARTIAL_STATUS_SINGLE   MartialStatus = "SINGLE"
	MARTIAL_STATUS_MARRIED  MartialStatus = "MARRID"
	MARTIAL_STATUS_DIVORCED MartialStatus = "DIVORCED"
	MARTIAL_STATUS_WIDOW    MartialStatus = "WIDOW"
)

const (
	ELST_REALM     Realm = "ELST"
	BANK_REALM     Realm = "BANK"
	DISTRICT_REALM Realm = "DISTRICT"
	BRANCH_REALM   Realm = "BRANCH"
	MERCHANT_REALM Realm = "MERCHANT"
	COMPANY_REALM  Realm = "COMPANY"
	MEMBER_REALM   Realm = "MEMBER"
)

const (
	GENDER_MALE   Gender = "MALE"
	GENDER_FEMALE Gender = "FEMALE"
)

const (
	REGISTRATION_TYPE_NEW    RegistrationType = "NEW"
	REGISTRATION_TYPE_LINKED RegistrationType = "LINKED"
)

const (
	OTP_FOR_LOGIN             OTPFor = "LOGIN"
	OTP_FOR_ADD_ACCOUNT       OTPFor = "ADD_ACCOUNT"
	OTP_FOR_PIN_SET           OTPFor = "PIN_SET"
	OTP_FOR_TRANSFER          OTPFor = "TRANSFER"
	OTP_FOR_ACCTIVATE_ACCOUNT OTPFor = "ACCTIVATE_ACCOUNT"
	OTP_FOR_PIN_RESET         OTPFor = "PIN_RESET"
	OTP_FOR_SIGNUP            OTPFor = "SIGNUP"
	OTP_FOR_ACCOUNT_LINK      OTPFor = "ACCOUNT_LINK"
	OTP_FOR_CHANGE_PHONE      OTPFor = "CHANGE_PHONE"
	OTP_FOR_DETACH_PHONE      OTPFor = "DETACH_PHONE"
	OTP_FOR_ATTACH_PHONE      OTPFor = "ATTACH_PHONE"
	OTP_FOR_ENABLE            OTPFor = "ENABLE"
	OTP_FOR_TRANSFER_LIMIT    OTPFor = "TRANSFER_LIMIT"
	OTP_FOR_CHANGE_EMAIL      OTPFor = "CHANGE_EMAIL"
	OTP_FOR_UPGRADE_LIMIT     OTPFor = "UPGRADE_LIMIT"
)

const (
	KYC_STATUS_PENDING  KYCStatus = "PENDING"
	KYC_STATUS_APPROVED KYCStatus = "APPROVED"
	KYC_STATUS_REJECTED KYCStatus = "REJECTED"
)

const (
	LDAP_STATUS_AUTHORIZED LDAPStatus = "AUTHORIZED"
	LDAP_STATUS_DENIED     LDAPStatus = "DENIED"
	LDAP_STATUS_PENDING    LDAPStatus = "PENDING"
	LDAP_STATUS_INITIATED  LDAPStatus = "INITIATED"
)

const (
	PRIMARY_AUTHENTICATION_BY_PHONE_NUMBER           PrimaryAuthentication = "PHONE_NUMBER"
	PRIMARY_AUTHENTICATION_BY_EMAIL                  PrimaryAuthentication = "EMAIL"
	PRIMARY_AUTHENTICATION_BY_EMAIL_AND_PHONE_NUMBER PrimaryAuthentication = "EMAIL_AND_PHONE_NUMBER"
)

const (
	DEVICE_STATUS_LINKED   DeviceStatus = "LINKED"
	DEVICE_STATUS_UNLINKED DeviceStatus = "UNLINKED"
)

const (
	OTP_STATUS_PENDING  OTPStatus = "PENDING"
	OTP_STATUS_VERIFIED OTPStatus = "VERIFIED"
	OTP_STATUS_DENIED   OTPStatus = "DENIED"
)

const (
	MERCHANT_ROLE_OWNER MerchantRole = "OWNER"
	MERCHANT_ROLE_AGENT MerchantRole = "AGENT"
)

const (
	POOL_SOURCE_PORTAL PoolSource = "PORTAL"
	POOL_SOURCE_APP    PoolSource = "APP"
	POOL_SOURCE_AGENT  PoolSource = "AGENT"
)

const (
	ACCOUNT_BRANCH_TYPE_CB  AccountBranchType = "CB"
	ACCOUNT_BRANCH_TYPE_IFB AccountBranchType = "IFB"
)

const (
	ACCOUNT_TYPE_NEW    AccountType = "NEW"
	ACCOUNT_TYPE_LINKED AccountType = "LINKED"
)

const (
	ACCOUNT_STATUS_ACTIVE   AccountStatus = "ACTIVE"
	ACCOUNT_STATUS_INACTIVE AccountStatus = "IN_ACTIVE"
)
