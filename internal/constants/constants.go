package constants

import "time"

type ContextKey string
type Platform string

const (
	// MESSAGE GROUP
	UpdateApp           = "Update your app"
	DeviceFound         = "Device Successfuly Found"
	OTPMessage          = "Your device lookup OTP is: %s. Valid for %d minutes."
	PhoneFound          = "Phone successfuly found"
	ImageUploadSuccess  = "image uploaded successfully"
	ProfileUploadSucess = "Profile theme set successfuly"
	PINResetSuccess     = "OTP verified for PIN reset"

	// constant
	Android                  Platform = "ANDROID"
	Ios                      Platform = "IOS"
	Prelogin                          = "PRE_LOGIN"
	OTPLength                         = 6
	DEV                               = "dev"
	UAT                               = "uat"
	Password                          = "PASSWORD"
	Login                             = "LOGIN"
	Change                            = "CHANGE"
	Checker                           = "CHECKER"
	IFBChecker                        = "IFB_CHECKER"
	Maker                             = "MAKER"
	IFBMaker                          = "IFBMAKER"
	Permanent                         = "PERMANENT"
	TokenType                         = "TOKEN_TYPE"
	Token                             = "TOKEN"
	OTP                               = "OTP"
	VerifyOtp                         = "VERIFY_OTP"
	DeviceLookUp                      = "DEVICE_LOOKUP"
	Register                          = "REGISTER"
	Empty                             = ""
	SetPin                            = "SET_PIN"
	Pin                               = "PIN"
	OTPForRegistration                = "REGISTRATION"
	Incomplete                        = "INCOMPLETE"
	ForgetPinVerifyOtp                = "FORGET_PIN_VERIFY_OTP"
	ResetPin                          = "RESET_PIN"
	Completed                         = "COMPLETED"
	ProfileTemp                       = "PROFILE-*.TMP"
	BucketUserProfilePicture          = "USER-PROFILE-PICTURES"
	OtpExpirationTime                 = 3 * time.Minute
	ActionCode                        = "action_code"
	Avatar                            = "avatar"
	DonationIcon					= "donation_icon"
	CampanyLogo 					="company_logo"
	ActionID                          = "action_id"
	ActionStatus                      = "action_status"
	IncompleteUserInfo                = "incomplete user info"
	Approved                          = "APPROVED"
	Rejected                          = "REJECTED"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
)

type Realm string

const (
	ELST_REALM     Realm = "ELST"
	BANK_REALM     Realm = "BANK"
	DISTRICT_REALM Realm = "DISTRICT"
	BRANCH_REALM   Realm = "BRANCH"
	MERCHANT_REALM Realm = "MERCHANT"
	COMPANY_REALM  Realm = "COMPANY"
	MEMBER_REALM   Realm = "MEMBER"
)

type ActionType string

const (
	ActionDelete ActionType = "DELETE"
	ActionUpdate ActionType = "UPDATE"
	ActionCreate ActionType = "CREATE"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
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

type MaritalStatus string

const (
	Single   MaritalStatus = "SINGLE"
	Married  MaritalStatus = "MARRIED"
	Divorced MaritalStatus = "DIVORCED"
	Widow    MaritalStatus = "WIDOW"
)

type DeviceStatus string

const (
	Linked   DeviceStatus = "LINKED"
	UnLinked DeviceStatus = "UNLINKED"
)

type AccountStatus string

const (
	Active   AccountStatus = "ACTIVE"
	InActive AccountStatus = "INACTIVE"
)

type MemberType string

const (
	CBT  MemberType = "CB"
	IFBT MemberType = "IFB"
)

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
	OTPForForgetPin        OTPFor = "FORGET_PIN"
)

type OTPStatus string

const (
	Pending  OTPStatus = "PENDING"
	Verified OTPStatus = "VERIFIED"
	Denied   OTPStatus = "DENIED"
)

const (
	UPDATE = "UPDATE"
	DELETE = "DELETE"
	CREATE = "CREATE"
)

type RequestAction string

const (
	RequestUser           RequestAction = "USER"
	RequestCpsUserCreate  RequestAction = "CREATE_CPS_USER"
	RequestCpsUserUpdate  RequestAction = "UPDATE_CPS_USER"
	RequestCpsUserDelete  RequestAction = "DELETE_CPS_USER"
	RequestCpsUserEnable  RequestAction = "ENABLE_CPS_USER"
	RequestCpsUserDisable RequestAction = "DISABLE_CPS_USER"

	RequestPermissionGroup          RequestAction = "PERMISSION_GROUP"
	RequestCreatePermissionGroup    RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup    RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup    RequestAction = "DELETE_PERMISSION_GROUP"
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
	RequestUpdateBankLogo           RequestAction = "UPDATE_BANK_LOGO	"
	RequestDeleteBank               RequestAction = "DELETE_BANK"
	RequestEnableBank               RequestAction = "ENABLE_BANK"
	RequestEnableDisableBank        RequestAction = "ENABLE_DISABLE_BANK"
	RequestDisableBank              RequestAction = "DISABLE_BANK"
	RequestCreateDepartment         RequestAction = "CREATE_DEPARTMENT"
	RequestUpdateDepartment         RequestAction = "UPDATE_DEPARTMENT"
	RequestDeleteDepartment         RequestAction = "DELETE_DEPARTMENT"
	RequestEnableDepartment         RequestAction = "ENABLE_DEPARTMENT"
	RequestDisableDepartment        RequestAction = "DISABLE_DEPARTMENT"
	RequestEnableDisableDepartment  RequestAction = "ENABLE_DISABLE_DEPARTMENT"
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
	RequestUpdateProductCode        RequestAction = "UPDATE_PRODUCT_CODE"
	RequestCreatePublicNotification RequestAction = "CREATE_PUBLIC_NOTIFICATION"
	RequestUpdatePublicNotification RequestAction = "UPDATE_PUBLIC_NOTIFICATION"
	RequestDeleteNotification       RequestAction = "DELETE_NOTIFICATION"
	RequestEnableNotification       RequestAction = "ENABLE_NOTIFICATION"
	RequestDisableNotification      RequestAction = "DISABLE_NOTIFICATION"
	RequestMarkNotificationAsSeen   RequestAction = "MARK_NOTIFICATION_AS_SEEN"
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
	RequestEnableMiniAppMerchant    RequestAction = "ENABLE_MINIAPP_MERCHANT"
	RequestDisableMiniAppMerchant   RequestAction = "DISABLE_MINIAPP_MERCHANT"
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar             RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar             RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar             RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar            RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar             RequestAction = "UPDATE_AVATAR"
	RequestUnlinkUser               RequestAction = "UNLINK_USER"
	RequestBlockRegion              RequestAction = "BLOCK_REGION"
	RequestBlockDistrict            RequestAction = "BLOCK_DISTRICT"
	RequestBlockCity                RequestAction = "BLOCK_CITY"
	RequestBlockUser                RequestAction = "BLOCK_USER"
	RequestEnableSingleBranches     RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches    RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches      RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches     RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"

	// Newly added for block_account
	// Branch
	RequestEnableBranches  RequestAction = "REQUEST_ENABLE_BRANCHES"
	RequestDisableBranches RequestAction = "REQUEST_DISABLE_BRANCHES"

	// Region
	RequestEnableRegions  RequestAction = "REQUEST_ENABLE_REGIONS"
	RequestDisableRegions RequestAction = "REQUEST_DISABLE_REGIONS"

	// District
	RequestEnableDistricts  RequestAction = "REQUEST_ENABLE_DISTRICTS"
	RequestDisableDistricts RequestAction = "REQUEST_DISABLE_DISTRICTS"

	// City
	RequestEnableCities  RequestAction = "REQUEST_ENABLE_CITIES"
	RequestDisableCities RequestAction = "REQUEST_DISABLE_CITIES"

	RequestBulkServiceEnable  RequestAction = "ENABLE_BULK_SERVICE"
	RequestBulkServiceDisable RequestAction = "DISABLE_BULK_SERVICE"

	RequestDeleteEvent             RequestAction = "DELETE_EVENT"
	RequestEnableEvent             RequestAction = "ENABLE_EVENT"
	RequestCreateEvent             RequestAction = "CREATE_EVENT"
	RequestUpdateEvent             RequestAction = "UPDATE_EVENT"
	RequestCreateEventCategory     RequestAction = "CREATE_EVENT_CATEGORY"
	RequestUpdateEventCategory     RequestAction = "UPDATE_EVENT_CATEGORY"
	RequestDisableEvent            RequestAction = "DISABLE_EVENT"
	RequestUpdateAccountValidation RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	RequestUpdateServiceDetails    RequestAction = "UPDATE_SERVICE_DETAILS"
	RequestUpdateHQBlockTime       RequestAction = "UPDATE_HQ_BLOCK_TIME"
	RequestUpdateHQArchiveTime     RequestAction = "UPDATE_HQ_ARCHIVE_TIME"

	RequestCreateBudgetColor             RequestAction = "BUDGET_CREATE_COLOR"
	RequestUpdateBudgetColor             RequestAction = "BUDGET_UPDATE_COLOR"
	RequestDeleteBudgetColor             RequestAction = "BUDGET_DELETE_COLOR"
	RequestCreateBudgetIcon              RequestAction = "BUDGET_CREATE_ICON"
	RequestUpdateBudgetIcon              RequestAction = "BUDGET_UPDATE_ICON"
	RequestDeleteBudgetIcAdvertServiceon RequestAction = "BUDGET_DELETE_ICON"

	RequestBudgetUpdate   RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestBudgetCreate   RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestBudgetDelete   RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestCreateMiniApp  RequestAction = "CREATE_MINIAPP"
	RequestUpdateMiniApp  RequestAction = "UPDATE_MINIAPP"
	RequestDeleteMiniApp  RequestAction = "DELETE_MINIAPP"
	RequestEnableMiniApp  RequestAction = "ENABLE_MINIAPP"
	RequestDisableMiniApp RequestAction = "DISABLE_MINIAPP"

	RequestUpdateEevent      RequestAction = "UPDATE_EVENT"
	RequestCIFRemove         RequestAction = "CIF_REMOVE"
	RequestServiceFlagUpdate RequestAction = "SERVICE_FLAG_UPDATE"

	RequestDisableFaydaAccount    RequestAction = "DISABLE_FAYDA_ACCOUNT"
	RequestEnableFaydaAccount     RequestAction = "ENABLE_FAYDA_ACCOUNT"
	RequestCreateDonationCategory RequestAction = "CREATE_DONATION_CATEGORY"
	RequestUpdateDonationCategory RequestAction = "UPDATE_DONATION_CATEGORY"
	RequestCreateDonationCompany  RequestAction = "CREATE_DONATION_COMPANY"
	RequestUpdateDonationCompany  RequestAction = "UPDATE_DONATION_COMPANY"
	RequestCreateDonation         RequestAction = "CREATE_DONATION"
	RequestUpdateDonation         RequestAction = "UPDATE_DONATION"


	ActionPending  RequestAction = "PENDING"
	ActionApproved RequestAction = "APPROVED"
	ActionRejected RequestAction = "REJECTED"
)

type RegistrationType string

const (
	RegistrationTypeNew    RegistrationType = "NEW"
	RegistrationTypeLinked RegistrationType = "LINKED"
)

type PrimaryAuthentication string

const (
	PrimaryAuthenticationByPhoneNumber         PrimaryAuthentication = "PHONE_NUMBER"
	PrimaryAuthenticationByEmail               PrimaryAuthentication = "EMAIL"
	PrimaryAuthenticationByEmailAndPhoneNumber PrimaryAuthentication = "EMAIL_AND_PHONE_NUMBER"
)

type AdvertFor string

const (
	IFB_ADVERT_FOR  AdvertFor = "IFB"
	CB_ADVERT_FOR   AdvertFor = "CB"
	BOTH_ADVERT_FOR AdvertFor = "ALL"
)

type RestrictionType string

const (
	AgeRestriction RestrictionType = "AGE_RESTRICTION"
)

type EventStatus string

const (
	EventUpcomming EventStatus = "UPCOMMING"
	EventLive      EventStatus = "LIVE"
	EventClosed    EventStatus = "CLOSED"
)

type BranchType string

const (
	IFB BranchType = "IFB"
	CB  BranchType = "CB"
)

type AppType string

const (
	URL AppType = "URL"
)
const (
	UATApp     AppType = "UAT"
	Production AppType = "PRODUCATION"
	Test       AppType = "TEST"
	Dev        AppType = "DEV"
)

type AppViewType string

const (
	AppViewTypeBoth AppViewType = "BOTH"
	AppViewTypeCB   AppViewType = "CB"
	AppViewTypeIFB  AppViewType = "IFB"
)

type Stage string

const (
	StageUat Stage = "UAT"
)

type EnvironmentType string

const (
	UatEnvironment        EnvironmentType = "UAT"
	DevEnvironment        EnvironmentType = "DEV"
	TestEnvironment       EnvironmentType = "TEST"
	ProductionEnvironment EnvironmentType = "PRODUCTION"
)

type Method string

const (
	OPEN      Method = "OPEN"
	PIN       Method = "PIN"
	OTPANDPIN Method = "OTP_PIN"
)

type NotificationFor string

const (
	ForIFB NotificationFor = "IFB"
	ForCB  NotificationFor = "CB"
	ForAll NotificationFor = "ALL"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusSeen    NotificationStatus = "SEEN"
)

const (
	UpdateAction ActionType = "UPDATE"
	CreateAction ActionType = "CREATE"
	DeleteAction ActionType = "DELETE"
)
