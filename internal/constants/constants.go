package constants

import (
	"fmt"
	"strings"
	"time"
)

type ContextKey string
type Platform string
type SettlementMethod string
type AuditorStatus string
type AuditorMark string

const (
	MARKEDASRIGHT AuditorMark = "MARKEDASRIGHT"
	MARKEDASWRONG AuditorMark = "MARKEDASWRONG"
)

const (
	AUDITORNOTCHECKED AuditorStatus = "NOTCHECKED"
	AUDITORINPROGRESS AuditorStatus = "INPROGRESS"
	AUDITORCHECKED    AuditorStatus = "CHECKED"
)
const (
	SettlementMethodDirect       SettlementMethod = "DIRECT"
	SettlementMethodGL           SettlementMethod = "GL"
	SettlementMethodMultiAccount SettlementMethod = "MULTI_ACCOUNT"
)

const (
	// Unique Separator

	UssdMerchant = "USSDM_"
	// MESSAGE GROUP
	UpdateApp           = "Update your app"
	DeviceFound         = "Device Successfuly Found"
	OTPMessage          = "Your device lookup OTP is: %s. Valid for %d minutes."
	PhoneFound          = "Phone successfuly found"
	ImageUploadSuccess  = "image uploaded successfully"
	ProfileUploadSucess = "Profile theme set successfuly"
	PINResetSuccess     = "OTP verified for PIN reset"

	// constant
	Android                  Platform   = "ANDROID"
	Ios                      Platform   = "IOS"
	Prelogin                            = "PRE_LOGIN"
	OTPLength                           = 6
	DEV                                 = "dev"
	UAT                                 = "uat"
	Cred                                = "Credential"
	Password                            = "PASSWORD"
	Login                               = "LOGIN"
	Change                              = "CHANGE"
	Checker                             = "CHECKER"
	IFBChecker                          = "IFB_CHECKER"
	Maker                               = "MAKER"
	IFBMaker                            = "IFBMAKER"
	Permanent                           = "PERMANENT"
	TokenType                           = "TOKEN_TYPE"
	Token                               = "TOKEN"
	OTP                                 = "OTP"
	VerifyOtp                           = "VERIFY_OTP"
	DeviceLookUp                        = "DEVICE_LOOKUP"
	Register                            = "REGISTER"
	Empty                               = ""
	SetPin                              = "SET_PIN"
	Pin                                 = "PIN"
	OTPForRegistration                  = "REGISTRATION"
	Incomplete                          = "INCOMPLETE"
	ForgetPinVerifyOtp                  = "FORGET_PIN_VERIFY_OTP"
	ResetPin                            = "RESET_PIN"
	Completed                           = "COMPLETED"
	ProfileTemp                         = "PROFILE-*.TMP"
	BucketUserProfilePicture            = "USER-PROFILE-PICTURES"
	OtpExpirationTime                   = 3 * time.Minute
	ActionCode                          = "action_code"
	Avatar                              = "avatar"
	DonationIcon                        = "donation_icon"
	CampanyLogo                         = "company_logo"
	DonationImage                       = "donation_image"
	DonationCoverImage                  = "donation_cover_image"
	BudgetCategoryIcon                  = "budget_category_icon"
	ActionID                            = "action_id"
	ActionStatus                        = "action_status"
	IncompleteUserInfo                  = "incomplete user info"
	Approved                            = "APPROVED"
	Rejected                            = "REJECTED"
	Reversed                            = "REVERSED"
	Canceled                            = "CANCELED"
	ContextKeyMetadata       ContextKey = "context_metadata"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
)

type Level uint8

const (
	ZERO Level = 0
	ONE  Level = 1
	TWO  Level = 2
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
	KYCStatusComplete KYCStatus = "COMPLETE"
	KYCStatusRejected KYCStatus = "REJECTED"
)

type BPSStatus string

const (
	BPSStatusAuthorized BPSStatus = "AUTHORIZED"
	BPSStatusDenied     BPSStatus = "DENIED"
	BPSStatusPending    BPSStatus = "PENDING"
	BPSStatusInitiated  BPSStatus = "INITIATED"
)

type Vendor string

const (
	Fayda    Vendor = "FAYDA"
	Verigram Vendor = "Verigram"
)

type BlockedOn string

const (
	NotBlocked BlockedOn = ""
	BPS        BlockedOn = "BPS"
	CPS        BlockedOn = "CPS"
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

type RiskLevel string

const (
	None   RiskLevel = ""
	High   RiskLevel = "HIGH"
	Midium RiskLevel = "MIDIUM"
	Low    RiskLevel = "LOW"
)

const (
	AccountNumber = "account_number"
	PhoneNumber   = "phone_number"
	WithFayda     = "with_fayda"
)

// type KYCStatus string

// const (
// 	KYCStatusPending  KYCStatus = "PENDING"
// 	KYCStatusApproved KYCStatus = "APPROVED"
// KYCStatusComplete KYCStatus = "COMPLETE"
// 	KYCStatusRejected KYCStatus = "REJECTED"
// )

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
	RequestCreateMiniappProductCode  string = "CREATE_MINI_APP_PRODUCT_CODE"
	RequestUpdateMiniappProductCode  string = "UPDATE_MINI_APP_PRODUCT_CODE"
	RequestDeleteMiniappProductCode  string = "DELETE_MINI_APP_PRODUCT_CODE"
	RequestEnableMiniappProductCode  string = "ENABLE_MINI_APP_PRODUCT_CODE"
	RequestDisableMiniappProductCode string = "DISABLE_MINI_APP_PRODUCT_CODE"

	RequestCreateJobRole  string = "CREATE_JOB_ROLE"
	RequestUpdateJobRole  string = "UPDATE_JOB_ROLE"
	RequestDeleteJobRole  string = "DELETE_JOB_ROLE"
	RequestEnableJobRole  string = "ENABLE_JOB_ROLE"
	RequestDisableJobRole string = "DISABLE_JOB_ROLE"

	RequestCreateRole  string = "CREATE_ROLE"
	RequestUpdateRole  string = "UPDATE_ROLE"
	RequestDeleteRole  string = "DELETE_ROLE"
	RequestEnableRole  string = "ENABLE_ROLE"
	RequestDisableRole string = "DISABLE_ROLE"

	// USSD
	RequestCreateUssdMerchant  string = "CREATE_USSD_MERCHANT"
	RequestUpdateUssdMerchant  string = "UPDATE_USSD_MERCHANT"
	RequestDeleteUssdMerchant  string = "DELETE_USSD_MERCHANT"
	RequestEnableUssdMerchant  string = "ENABLE_USSD_MERCHANT"
	RequestDisableUssdMerchant string = "DISABLE_USSD_MERCHANT"

	// Ecommerce
	RequestCreateEcommerceMerchant  RequestAction = "CREATE_ECOMMERCE_MERCHANT"
	RequestUpdateEcommerceMerchant  RequestAction = "UPDATE_ECOMMERCE_MERCHANT"
	RequestEnableEcommerceMerchant  RequestAction = "ENABLE_ECOMMERCE_MERCHANT"
	RequestDisableEcommerceMerchant RequestAction = "DISABLE_ECOMMERCE_MERCHANT"
	RequestDeleteEcommerceMerchant  RequestAction = "DELETE_ECOMMERCE_MERCHANT"

	RequestUser                  RequestAction = "USER"
	RequestActionRole            RequestAction = "ACTION_ROLE"
	RequestCreateActionRole      RequestAction = "CREATE_ACTION_ROLE"
	RequestUpdateActionRole      RequestAction = "UPDATE_ACTION_ROLE"
	RequestEnableActionRole      RequestAction = "ENABLE_ACTION_ROLE"
	RequestDisableActionRole     RequestAction = "DISABLE_ACTION_ROLE"
	RequestDeleteActionRole      RequestAction = "DELETE_ACTION_ROLE"
	RequestCreateCpsActionRole   RequestAction = "CREATE_CPS_ACTION_ROLE"
	RequestUpdateCpsActionRole   RequestAction = "UPDATE_CPS_ACTION_ROLE"
	RequestEnableCpsActionRole   RequestAction = "ENABLE_CPS_ACTION_ROLE"
	RequestDisableCpsActionRole  RequestAction = "DISABLE_CPS_ACTION_ROLE"
	RequestDeleteCpsActionRole   RequestAction = "DELETE_CPS_ACTION_ROLE"
	RequestCpsUserCreate         RequestAction = "CREATE_CPS_USER"
	RequestCpsUserUpdate         RequestAction = "UPDATE_CPS_USER"
	RequestCpsUserDelete         RequestAction = "DELETE_CPS_USER"
	RequestCpsUserEnable         RequestAction = "ENABLE_CPS_USER"
	RequestCpsUserDisable        RequestAction = "DISABLE_CPS_USER"
	RequestUpdateKYC             RequestAction = "UPDATE_KYC"
	RequestApproveKYC            RequestAction = "APPROVE_KYC"
	RequestPermissionGroup       RequestAction = "PERMISSION_GROUP"
	RequestCreatePermissionGroup RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup RequestAction = "DELETE_PERMISSION_GROUP"
	RequestDepartment            RequestAction = "DEPARMTENT"
	RequestEnableUser            RequestAction = "ENABLE_USER"
	RequestDisableUser           RequestAction = "DISABLE_USER"
	RequestBPSUser               RequestAction = "BPS_USER"

	// RequestBpsUserCreate  RequestAction = "CREATE_BPS_USER"
	// RequestBpsUserUpdate  RequestAction = "UPDATE_BPS_USER"
	// RequestBpsUserDelete  RequestAction = "DELETE_BPS_USER"
	// RequestBpsUserEnable  RequestAction = "ENABLE_BPS_USER"
	// RequestBpsUserDisable RequestAction = "DISABLE_BPS_USER"

	RequestBpsUserUpdate         RequestAction = "UPDATE_BPS_USER"
	RequestCreateBPSUser         RequestAction = "CREATE_BPS_USER"
	RequestDisableBPSUser        RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser         RequestAction = "ENABLE_BPS_USER"
	RequestUpdateUser            RequestAction = "UPDATE_USER"
	RequestTotalDailyLimit       RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT             RequestAction = "UPDATE_VAT"
	RequestAuthTier              RequestAction = "AUTHTIER"
	RequestDeleteAmountBasedAuth RequestAction = "DELETE_AMOUNT_BASED_AUTH"
	RequestCreateAmountBasedAuth RequestAction = "CREATE_AMOUNT_BASED_AUTH"
	RequestUpdateAmountBasedAuth RequestAction = "UPDATE_AMOUNT_BASED_AUTH"
	RequestCreateAdvert          RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert          RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert          RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert         RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert          RequestAction = "DELETE_ADVERT"
	RequestCreateBank            RequestAction = "CREATE_BANK"
	RequestUpdateBank            RequestAction = "UPDATE_BANK"
	RequestUpdateBankLogo        RequestAction = "UPDATE_BANK_LOGO"
	RequestDeleteBank            RequestAction = "DELETE_BANK"
	RequestEnableDisableBank     RequestAction = "ENABLE_DISABLE_BANK"
	RequestEnableBank            RequestAction = "ENABLE_BANK"
	RequestDisableBank           RequestAction = "DISABLE_BANK"
	RequestCreateDepartment      RequestAction = "CREATE_DEPARTMENT"

	RequestCreateDeviceVersion        RequestAction = "CREATE_DEVICE_VERSION"
	RequestUpdateDeviceVersion        RequestAction = "UPDATE_DEVICE_VERSION"
	RequestEnableDeviceVersion        RequestAction = "ENABLE_DEVICE_VERSION"
	RequestDisableDeviceVersion       RequestAction = "DISABLE_DEVICE_VERSION"
	RequestDeleteDeviceVersion        RequestAction = "DELETE_DEVICE_VERSION"
	RequestEnableDisableDeviceVersion RequestAction = "ENABLE_DISABLE_DEVICE_VERSION"

	RequestUpdateDepartment        RequestAction = "UPDATE_DEPARTMENT"
	RequestDeleteDepartment        RequestAction = "DELETE_DEPARTMENT"
	RequestEnableDepartment        RequestAction = "ENABLE_DEPARTMENT"
	RequestDisableDepartment       RequestAction = "DISABLE_DEPARTMENT"
	RequestEnableDisableDepartment RequestAction = "ENABLE_DISABLE_DEPARTMENT"
	RequestEnableDisableCustomer   RequestAction = "ENABLE_DISABLE_CUSTOMER"
	RequestApproveFaydaCustomer    RequestAction = "APPROVE_FAYDA_CUSTOMER"
	RequestCreateWallet            RequestAction = "CREATE_WALLET"
	RequestUpdateWallet            RequestAction = "UPDATE_WALLET"
	RequestDeleteWallet            RequestAction = "DELETE_WALLET"
	RequestEnableWallet            RequestAction = "ENABLE_WALLET"
	RequestDisableWallet           RequestAction = "DISABLE_WALLET"

	// Services catalog (model.Services)
	RequestCreateService            RequestAction = "CREATE_SERVICE"
	RequestUpdateService            RequestAction = "UPDATE_SERVICE"
	RequestEnableService            RequestAction = "ENABLE_SERVICE"
	RequestDisableService           RequestAction = "DISABLE_SERVICE"
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

	RequestArchiveUser            RequestAction = "ARCHIVE_USER"
	RequestCreatePasswordRule     RequestAction = "CREATE_PASSWORD_RULE"
	RequestUpdatePasswordRule     RequestAction = "UPDATE_PASSWORD_RULE"
	RequestUpdateMinimumService   RequestAction = "UPDATE_MINIMUM_SERVICE"
	RequestUpdateServiceRule      RequestAction = "UPDATE_SERVICE_RULE"
	RequestUpdateTotal            RequestAction = "UPDATE_TOTAL"
	RequestUpdateAccessConfig     RequestAction = "UPDATE_ACCESS_CONFIG"
	RequestEnableSingleBranch     RequestAction = "ENABLE_SINGLE_BRANCH"
	RequestDisableSingleBranch    RequestAction = "DISABLE_SINGLE_BRANCH"
	RequestEnableMultiUsers       RequestAction = "ENABLE_MULTI_USERS"
	RequestDisableMultiUsers      RequestAction = "DISABLE_MULTI_USERS"
	RequestCreateBusiness         RequestAction = "CREATE_BUSINESS"
	RequestUpdateBusiness         RequestAction = "UPDATE_BUSINESS"
	RequestCreateMiniAppMerchant  RequestAction = "CREATE_MINI_APP_MERCHANT"
	RequestUpdateMiniAppMerchant  RequestAction = "UPDATE_MINI_APP_MERCHANT"
	RequestEnableMiniAppMerchant  RequestAction = "ENABLE_MINI_APP_MERCHANT"
	RequestDeleteMiniAppMerchant  RequestAction = "DELETE_MINI_APP_MERCHANT"
	RequestDisableMiniAppMerchant RequestAction = "DISABLE_MINI_APP_MERCHANT"
	RequestUpdateBlockTime        RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar           RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar           RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar           RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar          RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar           RequestAction = "UPDATE_AVATAR"
	RequestUnlinkUser             RequestAction = "UNLINK_USER"
	RequestBlockRegion            RequestAction = "BLOCK_REGION"
	RequestBlockDistrict          RequestAction = "BLOCK_DISTRICT"
	RequestBlockCity              RequestAction = "BLOCK_CITY"
	RequestBlockUser              RequestAction = "BLOCK_USER"
	RequestEnableSingleBranches   RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches  RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches    RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches   RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"

	// for bankvault
	RequestCreateBankVault  RequestAction = "CREATE_VAULT_BANK"
	RequestUpdateBankVault  RequestAction = "UPDATE_VAULT_BANK"
	RequestDeleteBankVault  RequestAction = "DELETE_VAULT_BANK"
	RequestEnableBankVault  RequestAction = "ENABLE_VAULT_BANK"
	RequestDisAbleBankVault RequestAction = "DISABLE_VAULT_BANK"

	RequestCreateVaultAmountTier  RequestAction = "CREATE_VAULT_AMOUNT_TIER"
	RequestUpdateVaultAmountTier  RequestAction = "UPDATE_VAULT_AMOUNT_TIER"
	RequestDeleteVaultAmountTier  RequestAction = "DELETE_VAULT_AMOUNT_TIER"
	RequestEnableVaultAmountTier  RequestAction = "ENABLE_VAULT_AMOUNT_TIER"
	RequestDisAbleVaultAmountTier RequestAction = "DISABLE_VAULT_AMOUNT_TIER"

	// for vault group category
	RequestCreateVaultGroupCategory  RequestAction = "CREATE_VAULT_GROUP_CATEGORY"
	RequestUpdateVaultGroupCategory  RequestAction = "UPDATE_VAULT_GROUP_CATEGORY"
	RequestDeleteVaultGroupCategory  RequestAction = "DELETE_VAULT_GROUP_CATEGORY"
	RequestEnableVaultGroupCategory  RequestAction = "ENABLE_VAULT_GROUP_CATEGORY"
	RequestDisAbleVaultGroupCategory RequestAction = "DISABLE_VAULT_GROUP_CATEGORY"

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

	RequestCreateBudgetCategory  RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestDeleteBudgetCategory  RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestUpdateBudgetCategory  RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestDisableBudgetCategory RequestAction = "DISABLE_BUDGET_CATEGORY"
	RequestEnableBudgetCategory  RequestAction = "ENABLE_BUDGET_CATEGORY"

	RequestCreateCustomerSegmentation  RequestAction = "CREATE_CUSTOMER_SEGMENTATION"
	RequestUpdateCustomerSegmentation  RequestAction = "UPDATE_CUSTOMER_SEGMENTATION"
	RequestEnableCustomerSegmentation  RequestAction = "ENABLE_CUSTOMER_SEGMENTATION"
	RequestDisableCustomerSegmentation RequestAction = "DISABLE_CUSTOMER_SEGMENTATION"
	RequestDeleteCustomerSegmentation  RequestAction = "DELETE_CUSTOMER_SEGMENTATION"

	RequestCreateCpsRole  RequestAction = "CREATE_CPS_ROLE"
	RequestUpdateCpsRole  RequestAction = "UPDATE_CPS_ROLE"
	RequestDeleteCpsRole  RequestAction = "DELETE_CPS_ROLE"
	RequestEnableCpsRole  RequestAction = "ENABLE_CPS_ROLE"
	RequestDisableCpsRole RequestAction = "DISABLE_CPS_ROLE"

	RequestCreateMiniApp  RequestAction = "CREATE_MINI_APP"
	RequestUpdateMiniApp  RequestAction = "UPDATE_MINI_APP"
	RequestDeleteMiniApp  RequestAction = "DELETE_MINI_APP"
	RequestEnableMiniApp  RequestAction = "ENABLE_MINI_APP"
	RequestDisableMiniApp RequestAction = "DISABLE_MINI_APP"

	RequestUpdateEevent      RequestAction = "UPDATE_EVENT"
	RequestCIFRemove         RequestAction = "CIF_REMOVE"
	RequestServiceFlagUpdate RequestAction = "SERVICE_FLAG_UPDATE"

	RequestDisableFaydaAccount     RequestAction = "DISABLE_FAYDA_ACCOUNT"
	RequestEnableFaydaAccount      RequestAction = "ENABLE_FAYDA_ACCOUNT"
	RequestCreateDonationCategory  RequestAction = "CREATE_DONATION_CATEGORY"
	RequestUpdateDonationCategory  RequestAction = "UPDATE_DONATION_CATEGORY"
	RequestEnableDonationCategory  RequestAction = "ENABLE_DONATION_CATEGORY"
	RequestDisableDonationCategory RequestAction = "DISABLE_DONATION_CATEGORY"
	RequestEnableDonationCompany   RequestAction = "ENABLE_DONATION_COMPANY"
	RequestDisableDonationCompany  RequestAction = "DISABLE_DONATION_COMPANY"
	RequestCreateDonationCompany   RequestAction = "CREATE_DONATION_COMPANY"
	RequestUpdateDonationCompany   RequestAction = "UPDATE_DONATION_COMPANY"
	RequestCreateDonation          RequestAction = "CREATE_DONATION"
	RequestUpdateDonation          RequestAction = "UPDATE_DONATION"
	RequestUpdateDonationImage     RequestAction = "UPDATE_DONATION_IMAGE"
	RequestDeleteDonationImage     RequestAction = "DELETE_DONATION_IMAGE"
	RequestAddDonationImage        RequestAction = "ADD_DONATION_IMAGE"
	RequestEnableDonation          RequestAction = "ENABLE_DONATION"
	RequestDisableDonation         RequestAction = "DISABLE_DONATION"

	RequestCreateArticle  RequestAction = "CREATE_ARTICLE"
	RequestUpdateArticle  RequestAction = "UPDATE_ARTICLE"
	RequestEnableArticle  RequestAction = "ENABLE_ARTICLE"
	RequestDisableArticle RequestAction = "DISABLE_ARTICLE"
	RequestDeleteArticle  RequestAction = "DELETE_ARTICLE"

	// article categories
	RequestCreateArticleCategory  RequestAction = "CREATE_ARTICLE_CATEGORY"
	RequestUpdateArticleCategory  RequestAction = "UPDATE_ARTICLE_CATEGORY"
	RequestEnableArticleCategory  RequestAction = "ENABLE_ARTICLE_CATEGORY"
	RequestDisableArticleCategory RequestAction = "DISABLE_ARTICLE_CATEGORY"
	RequestDeleteArticleCategory  RequestAction = "DELETE_ARTICLE_CATEGORY"

	// article tags
	RequestCreateNewsTag  RequestAction = "CREATE_NEWS_TAG"
	RequestUpdateNewsTag  RequestAction = "UPDATE_NEWS_TAG"
	RequestEnableNewsTag  RequestAction = "ENABLE_NEWS_TAG"
	RequestDisableNewsTag RequestAction = "DISABLE_NEWS_TAG"
	RequestDeleteNewsTag  RequestAction = "DELETE_NEWS_TAG"

	RequestCreateTopup  RequestAction = "CREATE_TOPUP"
	RequestUpdateTopup  RequestAction = "UPDATE_TOPUP"
	RequestDeleteTopup  RequestAction = "DELETE_TOPUP"
	RequestEnableTopup  RequestAction = "ENABLE_TOPUP"
	RequestDisableTopup RequestAction = "DISABLE_TOPUP"

	RequestCreateShortVideo  string = "CREATE_SHORT_VIDEO"
	RequestUpdateShortVideo  string = "UPDATE_SHORT_VIDEO"
	RequestEnableShortVideo  string = "ENABLE_SHORT_VIDEO"
	RequestDisableShortVideo string = "DISABLE_SHORT_VIDEO"
	RequestDeleteShortVideo  string = "DELETE_SHORT_VIDEO"

	ActionPending  RequestAction = "PENDING"
	ActionApproved RequestAction = "APPROVED"
	ActionRejected RequestAction = "REJECTED"

	RequestCreateNewsCategory RequestAction = "CREATE_NEWS_CATEGORY"
	RequestUpdateNewsCategory RequestAction = "UPDATE_NEWS_CATEGORY"
	RequestDeleteNewsCategory RequestAction = "DELETE_NEWS_CATEGORY"

	RequestCreateMiniAppCategory  RequestAction = "CREATE_MINI_APP_CATEGORY"
	RequestUpdateMiniAppCategory  RequestAction = "UPDATE_MINI_APP_CATEGORY"
	RequestDeleteMiniAppCategory  RequestAction = "DELETE_MINI_APP_CATEGORY"
	RequestEnableMiniAppCategory  RequestAction = "ENABLE_MINI_APP_CATEGORY"
	RequestDisableMiniAppCategory RequestAction = "DISABLE_MINI_APP_CATEGORY"

	RequestEnableEventMerchant  RequestAction = "ENABLE_EVENT_MERCHANT"
	RequestDisableEventMerchant RequestAction = "DISABLE_EVENT_MERCHANT"
	RequestCreateEventMerchant  RequestAction = "CREATE_EVENT_MERCHANT"
	RequestUpdateEventMerchant  RequestAction = "UPDATE_EVENT_MERCHANT"
	RequestDeleteEventMerchant  RequestAction = "DELETE_EVENT_MERCHANT"

	RequestEnableLogisticsMerchant  RequestAction = "ENABLE_LOGISTICS_MERCHANT"
	RequestDisableLogisticsMerchant RequestAction = "DISABLE_LOGISTICS_MERCHANT"
	RequestCreateLogisticsMerchant  RequestAction = "CREATE_LOGISTICS_MERCHANT"
	RequestUpdateLogisticsMerchant  RequestAction = "UPDATE_LOGISTICS_MERCHANT"
	RequestDeleteLogisticsMerchant  RequestAction = "DELETE_LOGISTICS_MERCHANT"

	RequestCreateAccessListSegmentation        RequestAction = "CREATE_ACCESS_LIST_SEGMENTATION"
	RequestUpdateAccessListSegmentation        RequestAction = "UPDATE_ACCESS_LIST_SEGMENTATION"
	RequestEnableDisableAccessListSegmentation RequestAction = "ENABLE_DISABLE_ACCESS_LIST_SEGMENTATION"

	RequestCreateCustomerKYC RequestAction = "CREATE_CUSTOMER_KYC"
	RequestUpdateCustomerKYC RequestAction = "UPDATE_CUSTOMER_KYC"
	RequestDeleteCustomerKYC RequestAction = "DELETE_CUSTOMER_KYC"
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

type VaultStatus string

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

type VaultType string

const (
	VaultTypeBanking VaultType = "BANK_VAULT"
	VaultTypePrivate VaultType = "PRIVATE_VAULT"
	VaultTypeGroup   VaultType = "GROUP_VAULT"
)

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

type AccountType string

const (
	AccountTypeIFB AccountType = "IFB"
	AccountTypeCB  AccountType = "CB"
)

type InviteStatus string

const (
	InvitePending   InviteStatus = "PENDING"
	InviteAccepted  InviteStatus = "ACCEPTED"
	InviteRejected  InviteStatus = "REJECTED"
	InviteWithdrawn InviteStatus = "WITHDRAWN"
	InviteExpired   InviteStatus = "EXPIRED"
)

type MemberStatus string

const (
	StatusPaid    MemberStatus = "PAID"
	StatusWaiting MemberStatus = "WAITING"
	StatusOverdue MemberStatus = "OVERDUE"
)

type MemberRole string

const (
	RoleAdmin       MemberRole = "ADMIN"
	RoleParticipant MemberRole = "PARTICIPANT"
)

// VaultCategory defines the category of the vault.
type VaultCategory string

const (
	CategoryTravel  VaultCategory = "TRAVEL"
	CategoryGadgets VaultCategory = "GADGETS"
	CategoryEvents  VaultCategory = "EVENTS"
	CategoryGeneral VaultCategory = "GENERAL"
)

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

// Recurrence defines the frequency of recurring payments.
type Recurrence string

const (
	RecurrenceWeekly   Recurrence = "WEEKLY"
	RecurrenceBiweekly Recurrence = "BIWEEKLY"
	RecurrenceMonthly  Recurrence = "MONTHLY"
)

// ContributionType defines the type of contribution made by a member.
type ContributionType string

const (
	ContributionTypeDeposit    ContributionType = "DEPOSIT"
	ContributionTypeWithdrawal ContributionType = "WITHDRAWAL"
	// ContributionTypeFee        ContributionType = "FEE"
	// ContributionTypeDisbursement ContributionType = "DISBURSEMENT"
)

type Type string

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

type KafkaTopic string

const (
	ClientOrchestrationMemberTopic               KafkaTopic = "member.sync.cps"
	ClientOrchestrationKycTopic                  KafkaTopic = "kyc.sync.cps"
	ClientOrchestrationAccessControlTopic        KafkaTopic = "access_control.sync.cps"
	ClientOrchestrationLinkedAccountTopic        KafkaTopic = "linked_account.sync.cps"
	ClientOrchestrationAccountBlockTopic         KafkaTopic = "account_block.sync.cps"
	ClientOrchestrationAccountValidationTopic    KafkaTopic = "account_validation.sync.cps"
	ClientOrchestrationDeviceVersionControlTopic KafkaTopic = "device_version_control.sync.cps"
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
)

type FinancialInstitutionType string

const (
	Bank   FinancialInstitutionType = "BANK"
	Wallet FinancialInstitutionType = "WALLET"
	MFI    FinancialInstitutionType = "MFI"
)

type maxMemory int64

const (
	MaxMemoryForUpload maxMemory = 15 << 20 // 15 MB
)

type ImageFolderName string

const (
	BankFolderName               ImageFolderName = "banks"
	UssdMerchantFolderName       ImageFolderName = "ussd_merchant"
	WalletFolderName             ImageFolderName = "wallets"
	AdFolderName                 ImageFolderName = "ads"
	AvatarFolderName             ImageFolderName = "avatars"
	BudgetCategoryFolderName     ImageFolderName = "budget_categories"
	DonationFolderName           ImageFolderName = "donations"
	DonationCategoryFolderName   ImageFolderName = "donation_categories"
	DonationCompanyFolderName    ImageFolderName = "donation_companies"
	EventFolderName              ImageFolderName = "events"
	TopupFolderName              ImageFolderName = "topups"
	VaultGroupCategoryFolderName ImageFolderName = "vault_group_categories"
	CustomerKYCFolderName        ImageFolderName = "customer_kyc"
)

// redis key prefixes
const (
	RedisCPSUserDeviceIDPrefix = "cps:auth:device"
)
