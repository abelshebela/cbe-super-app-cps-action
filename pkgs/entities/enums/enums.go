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
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)
const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeTimeout      ErrorType = "TIMEOUT"
)
const (
	EventUpcomming EventStatus = "upcomming"
	EventLive      EventStatus = "live"
	EventClosed    EventStatus = "closed"
)

// RestrictionType represents the type of restriction for an event
// (from internal/domain/event/entity.go)
type RestrictionType string

const (
	AgeRestriction RestrictionType = "age_restriction"
)

// TicketStatus represents the status of a ticket
// (from internal/domain/ticket/entity.go)
type TicketStatus string

const (
	TicketBooked   TicketStatus = "booked"
	TicketExpired  TicketStatus = "expired"
	TicketPaid     TicketStatus = "paid"
	TicketRedeemed TicketStatus = "redeemed"
)

// AdvertFor represents the target audience for an advert
// (from internal/domain/ad/entity.go)
type AdvertFor string

const (
	IFBAdvertFor  AdvertFor = "ifb"
	CBAdvertFor   AdvertFor = "cb"
	BothAdvertFor AdvertFor = "all"
)

// EnvironmentType represents the environment type for miniapp
// (from internal/domain/miniapp/entity.go)
type EnvironmentType string

const (
	UatEnvironment        EnvironmentType = "uat"
	DevEnvironment        EnvironmentType = "dev"
	TestEnvironment       EnvironmentType = "test"
	ProductionEnvironment EnvironmentType = "production"
)

// BranchType represents the branch type for miniapp
// (from internal/domain/miniapp/entity.go)
type BranchType string

const (
	IFBBranchType BranchType = "ifb"
	CBBranchType  BranchType = "cb"
)

const (
	Maker   UserType = "maker"
	Checker UserType = "checker"
)

const (
	ActionPending  ActionStatus = "pending"
	ActionApproved ActionStatus = "approved"
	ActionRejected ActionStatus = "rejected"
)

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

const (
	RequestUser                     RequestAction = "user"
	RequestPermissionGroup          RequestAction = "permission_group"
	RequestDepartment               RequestAction = "deparmtent"
	RequestEnableUser               RequestAction = "enable_user"
	RequestDisableUser              RequestAction = "disable_user"
	RequestBPSUser                  RequestAction = "bps_user"
	RequestDisableBPSUser           RequestAction = "disable_bps_user"
	RequestEnableBPSUser            RequestAction = "enable_bps_user"
	RequestUpdateUser               RequestAction = "update_user"
	RequestTotalDailyLimit          RequestAction = "total_daily_limit"
	RequestUpdateVAT                RequestAction = "update_vat"
	RequestAuthTier                 RequestAction = "authtier"
	RequestCreateAdvert             RequestAction = "create_advert"
	RequestUpdateAdvert             RequestAction = "update_advert"
	RequestEnableAdvert             RequestAction = "enable_advert"
	RequestDisableAdvert            RequestAction = "disable_advert"
	RequestDeleteAdvert             RequestAction = "delete_advert"
	RequestCreateBank               RequestAction = "create_bank"
	RequestUpdateBank               RequestAction = "update_bank"
	RequestEnableBank               RequestAction = "enable_bank"
	RequestDisableBank              RequestAction = "disable_bank"
	RequestEnableWallet             RequestAction = "enable_wallet"
	RequestDisableWallet            RequestAction = "disable_wallet"
	RequestUpdatePasswordExpiry     RequestAction = "update_password_expiry"
	RequestCreateValidation         RequestAction = "create_validation"
	RequestUpdateValidation         RequestAction = "update_validation"
	RequestDeleteValidation         RequestAction = "delete_validation"
	RequestUpdateArchiveExpiry      RequestAction = "update_archive_expiry"
	RequestCreateServiceFee         RequestAction = "create_service_fee"
	RequestUpdateServiceFee         RequestAction = "update_service_fee"
	RequestDeleteServiceFee         RequestAction = "delete_service_fee"
	RequestCreateDailyLimit         RequestAction = "create_daily_limit"
	RequestUpdateDailyLimit         RequestAction = "update_daily_limit"
	RequestDeleteDailyLimit         RequestAction = "delete_daily_limit"
	RequestBudgetColor              RequestAction = "budget_color"
	RequestBudgetIcon               RequestAction = "budget_icon"
	RequestUpdateProduct            RequestAction = "update_product"
	RequestCreatePublicNotification RequestAction = "create_public_notification"
	RequestArchiveUser              RequestAction = "archive_user"
	RequestCreatePasswordRule       RequestAction = "create_password_rule"
	RequestUpdatePasswordRule       RequestAction = "update_password_rule"
	RequestUpdateMinimumService     RequestAction = "update_minimum_service"
	RequestUpdateServiceRule        RequestAction = "update_service_rule"
	RequestUpdateTotal              RequestAction = "update_total"
	RequestUpdateAccessConfig       RequestAction = "update_access_config"
	RequestEnableSingleBranch       RequestAction = "enable_single_branch"
	RequestDisableSingleBranch      RequestAction = "disable_single_branch"
	RequestEnableMultiUsers         RequestAction = "enable_multi_users"
	RequestDisableMultiUsers        RequestAction = "disable_multi_users"
	RequestCreateBusiness           RequestAction = "create_business"
	RequestUpdateBusiness           RequestAction = "update_business"
	RequestCreateEvent              RequestAction = "create_event"
	RequestUpdateEvent              RequestAction = "update_event"
	RequestCreateEventCategory      RequestAction = "create_event_category"
	RequestUpdateEventCategory      RequestAction = "update_event_category"
	RequestDisableEvent             RequestAction = "disable_event"
	RequestCreateMiniAppMerchant    RequestAction = "create_miniapp_merchant"
	RequestUpdateMiniAppMerchant    RequestAction = "update_miniapp_merchant"
	RequestUpdateBlockTime          RequestAction = "update_block_time"
	RequestUpdateAccountValidation  RequestAction = "update_account_validation"
	RequestUpdateServiceDetails     RequestAction = "update_service_details"
	RequestUpdateHQBlockTime        RequestAction = "update_hq_block_time"
	RequestUpdateHQArchiveTime      RequestAction = "update_hq_archive_time"
)
const (
	KYCLevelZero KYCLevel = "zero"
	KYCLevelOne  KYCLevel = "one"
	KYCLevelTwo  KYCLevel = "two"
)

const (
	BPSStatusAuthorized BPSStatus = "authorized"
	BPSStatusDenied     BPSStatus = "denied"
	BPSStatusPending    BPSStatus = "pending"
	BPSStatusInitiated  BPSStatus = "initiated"
)

const (
	CBT  MemberType = "cb"
	IFBT MemberType = "ifb"
)

const (
	Single   MartialStatus = "single"
	Married  MartialStatus = "marrid"
	Divorced MartialStatus = "divorced"
	Widow    MartialStatus = "widow"
)

const (
	ElstRealm     Realm = "elst"
	BankRealm     Realm = "bank"
	DistrictRealm Realm = "district"
	BranchRealm   Realm = "branch"
	MerchantRealm Realm = "merchant"
	CompanyRealm  Realm = "company"
	MemberRealm   Realm = "member"
)

const (
	Male   Gender = "male"
	Female Gender = "female"
)

const (
	RegistrationTypeNew    RegistrationType = "new"
	RegistrationTypeLinked RegistrationType = "linked"
)

const (
	OTPForLogin            OTPFor = "login"
	OTPForAddAccount       OTPFor = "add_account"
	OTPForPINSet           OTPFor = "pin_set"
	OTPForTransfer         OTPFor = "transfer"
	OTPForAcctivateAccount OTPFor = "acctivate_account"
	OTPForPINReset         OTPFor = "pin_reset"
	OTPForSignup           OTPFor = "signup"
	OTPForAccountLink      OTPFor = "account_link"
	OTPForChangePhone      OTPFor = "change_phone"
	OTPForDetachPhone      OTPFor = "detach_phone"
	OTPForAttachPhone      OTPFor = "attach_phone"
	OTPForEnable           OTPFor = "enable"
	OTPForTransferLimit    OTPFor = "transfer_limit"
	OTPForChangeEmail      OTPFor = "change_email"
	OTPForUpgradeLimit     OTPFor = "upgrade_limit"
)

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusApproved KYCStatus = "approved"
	KYCStatusRejected KYCStatus = "rejected"
)

const (
	LDAPStatusAuthorized LDAPStatus = "authorized"
	LDAPStatusDenied     LDAPStatus = "denied"
	LDAPStatusPending    LDAPStatus = "pending"
	LDAPStatusInitiated  LDAPStatus = "initiated"
)

const (
	PrimaryAuthenticationByPhoneNumber         PrimaryAuthentication = "phone_number"
	PrimaryAuthenticationByEmail               PrimaryAuthentication = "email"
	PrimaryAuthenticationByEmailAndPhoneNumber PrimaryAuthentication = "email_and_phone_number"
)

const (
	Linked   DeviceStatus = "linked"
	UnLinked DeviceStatus = "unlinked"
)

const (
	Pending  OTPStatus = "PENDING"
	Verified OTPStatus = "verified"
	Denied   OTPStatus = "denied"
)

const (
	MerchantRoleOwner MerchantRole = "owner"
	MerchantRoleAgent MerchantRole = "agent"
)

const (
	PoolSourcePortal PoolSource = "portal"
	PoolSourceApp    PoolSource = "app"
	PoolSourceAgent  PoolSource = "agent"
)

const (
	CB  AccountBranchType = "cb"
	IFB AccountBranchType = "ifb"
)

const (
	AccountTypeNew    AccountType = "new"
	AccountTypeLinked AccountType = "linked"
)

const (
	Active   AccountStatus = "active"
	InActive AccountStatus = "in_active"
)
