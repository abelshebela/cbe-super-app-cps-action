package cpsaction

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

func IsValidActionStatus(status string) bool {
	switch ActionStatus(status) {
	case ActionPending, ActionApproved, ActionRejected:
		return true
	default:
		return false
	}
}

type ActionType string

const (
	ActionCreate  ActionType = "CREATE"
	ActionUpdate  ActionType = "UPDATE"
	ActionDelete  ActionType = "DELETE"
	ActionEnable  ActionType = "ENABLE"
	ActionDisable ActionType = "DISABLE"
)

func IsValidActionType(actionType string) bool {
	switch ActionType(actionType) {
	case ActionCreate, ActionUpdate, ActionDelete:
		return true
	default:
		return false
	}
}

type RequestAction string

const (
	RequestCreateCustomerSegmentation RequestAction = "CREATE_CUSTOMER_SEGMENTATION"
	RequestUpdateCustomerSegmentation RequestAction = "UPDATE_CUSTOMER_SEGMENTATION"
	RequestDeleteCustomerSegmentation RequestAction = "DELETE_CUSTOMER_SEGMENTATION"

	RequestCreateMiniappProductCode  RequestAction = "CREATE_MINI_APP_PRODUCT_CODE"
	RequestUpdateMiniappProductCode  RequestAction = "UPDATE_MINI_APP_PRODUCT_CODE"
	RequestDeleteMiniappProductCode  RequestAction = "DELETE_MINI_APP_PRODUCT_CODE"
	RequestEnableMiniappProductCode  RequestAction = "ENABLE_MINI_APP_PRODUCT_CODE"
	RequestDisableMiniappProductCode RequestAction = "DISABLE_MINI_APP_PRODUCT_CODE"

	RequestCreateJobRole  RequestAction = "CREATE_JOB_ROLE"
	RequestUpdateJobRole  RequestAction = "UPDATE_JOB_ROLE"
	RequestDeleteJobRole  RequestAction = "DELETE_JOB_ROLE"
	RequestEnableJobRole  RequestAction = "ENABLE_JOB_ROLE"
	RequestDisableJobRole RequestAction = "DISABLE_JOB_ROLE"

	RequestCreateRole  RequestAction = "CREATE_ROLE"
	RequestUpdateRole  RequestAction = "UPDATE_ROLE"
	RequestDeleteRole  RequestAction = "DELETE_ROLE"
	RequestEnableRole  RequestAction = "ENABLE_ROLE"
	RequestDisableRole RequestAction = "DISABLE_ROLE"

	RequestDeleteAmountBasedAuth RequestAction = "DELETE_AMOUNT_BASED_AUTH"
	RequestCreateAmountBasedAuth RequestAction = "CREATE_AMOUNT_BASED_AUTH"
	RequestUpdateAmountBasedAuth RequestAction = "UPDATE_AMOUNT_BASED_AUTH"
	RequestUser                  RequestAction = "USER"
	RequestCpsUserCreate         RequestAction = "CREATE_CPS_USER"
	RequestCpsUserUpdate         RequestAction = "UPDATE_CPS_USER"
	RequestCpsUserDelete         RequestAction = "DELETE_CPS_USER"
	RequestCpsUserEnable         RequestAction = "ENABLE_CPS_USER"
	RequestCpsUserDisable        RequestAction = "DISABLE_CPS_USER"
	RequestPermissionGroup       RequestAction = "PERMISSION_GROUP"
	RequestBulkServiceEnable     RequestAction = "ENABLE_BULK_SERVICE"
	RequestBulkServiceDisable    RequestAction = "DISABLE_BULK_SERVICE"
	// RequestDepartment               RequestAction = "DEPARMTENT"
	RequestCreateDepartment           RequestAction = "CREATE_DEPARTMENT"
	RequestUpdateDepartment           RequestAction = "UPDATE_DEPARTMENT"
	RequestEnableDisableDepartment    RequestAction = "ENABLE_DISABLE_DEPARTMENT"
	RequestEnableUser                 RequestAction = "ENABLE_USER"
	RequestDisableUser                RequestAction = "DISABLE_USER"
	RequestBPSUser                    RequestAction = "BPS_USER"
	RequestDisableBPSUser             RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser              RequestAction = "ENABLE_BPS_USER"
	RequestUpdateUser                 RequestAction = "UPDATE_USER"
	RequestTotalDailyLimit            RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT                  RequestAction = "UPDATE_VAT"
	RequestAuthTier                   RequestAction = "AUTHTIER"
	RequestCreateDeviceVersion        RequestAction = "CREATE_DEVICE_VERSION"
	RequestUpdateDeviceVersion        RequestAction = "UPDATE_DEVICE_VERSION"
	RequestEnableDeviceVersion        RequestAction = "ENABLE_DEVICE_VERSION"
	RequestDisableDeviceVersion       RequestAction = "DISABLE_DEVICE_VERSION"
	RequestDeleteDeviceVersion        RequestAction = "DELETE_DEVICE_VERSION"
	RequestEnableDisableDeviceVersion RequestAction = "ENABLE_DISABLE_DEVICE_VERSION"
	RequestCreateAdvert               RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert               RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert               RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert              RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert               RequestAction = "DELETE_ADVERT"
	RequestCreateBank                 RequestAction = "CREATE_BANK"
	RequestUpdateBank                 RequestAction = "UPDATE_BANK"
	RequestDeleteBank                 RequestAction = "DELETE_BANK"
	RequestUpdateBankLogo             RequestAction = "UPDATE_BANK_LOGO"
	RequestEnableDisableBank          RequestAction = "ENABLE_DISABLE_BANK"
	RequestEnableBank                 RequestAction = "ENABLE_BANK"
	RequestDisableBank                RequestAction = "DISABLE_BANK"
	RequestCreateWallet               RequestAction = "CREATE_WALLET"
	RequestUpdateWallet               RequestAction = "UPDATE_WALLET"
	RequestDeleteWallet               RequestAction = "DELETE_WALLET"
	RequestEnableWallet               RequestAction = "ENABLE_WALLET"
	RequestDisableWallet              RequestAction = "DISABLE_WALLET"

	// Services catalog (model.Services)
	RequestCreateService  RequestAction = "CREATE_SERVICE"
	RequestUpdateService  RequestAction = "UPDATE_SERVICE"
	RequestEnableService  RequestAction = "ENABLE_SERVICE"
	RequestDisableService RequestAction = "DISABLE_SERVICE"

	RequestCreateTopup  RequestAction = "CREATE_TOPUP"
	RequestUpdateTopup  RequestAction = "UPDATE_TOPUP"
	RequestDeleteTopup  RequestAction = "DELETE_TOPUP"
	RequestEnableTopup  RequestAction = "ENABLE_TOPUP"
	RequestDisableTopup RequestAction = "DISABLE_TOPUP"

	RequestUpdatePasswordExpiry RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreateValidation     RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation     RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation     RequestAction = "DELETE_VALIDATION"
	RequestUpdateArchiveExpiry  RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestUpdateServiceSingle  RequestAction = "UPDATE_SERVICE_SINGLE_CAP"
	RequestUpdateServiceTotal   RequestAction = "UPDATE_SERVICE_TOTAL_CAP"
	RequestUpdateServiceMinCap  RequestAction = "UPDATE_SERVICE_MIN_CAP"
	RequestBudgetColor          RequestAction = "BUDGET_COLOR"
	RequestBudgetIcon           RequestAction = "BUDGET_ICON"
	RequestUpdateProduct        RequestAction = "UPDATE_PRODUCT"

	// Notification
	RequestCreatePublicNotification RequestAction = "CREATE_PUBLIC_NOTIFICATION"
	RequestUpdatePublicNotification RequestAction = "UPDATE_PUBLIC_NOTIFICATION"
	RequestDeleteNotification       RequestAction = "DELETE_NOTIFICATION"
	RequestEnableNotification       RequestAction = "ENABLE_NOTIFICATION"
	RequestDisableNotification      RequestAction = "DISABLE_NOTIFICATION"
	RequestMarkNotificationAsSeen   RequestAction = "MARK_NOTIFICATION_AS_SEEN"

	RequestArchiveUser             RequestAction = "ARCHIVE_USER"
	RequestUpdatePasswordRule      RequestAction = "UPDATE_PASSWORD_RULE"
	RequestUpdateMinimumService    RequestAction = "UPDATE_MINIMUM_SERVICE"
	RequestUpdateServiceRule       RequestAction = "UPDATE_SERVICE_RULE"
	RequestUpdateTotal             RequestAction = "UPDATE_TOTAL"
	RequestUpdateAccessConfig      RequestAction = "UPDATE_ACCESS_CONFIG"
	RequestEnableSingleBranch      RequestAction = "ENABLE_SINGLE_BRANCH"
	RequestDisableSingleBranch     RequestAction = "DISABLE_SINGLE_BRANCH"
	RequestUpdateAccountValidation RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	RequestEnableMultiUsers        RequestAction = "ENABLE_MULTI_USERS"
	RequestDisableMultiUsers       RequestAction = "DISABLE_MULTI_USERS"
	RequestCreateBusiness          RequestAction = "CREATE_BUSINESS"
	RequestUpdateBusiness          RequestAction = "UPDATE_BUSINESS"
	RequestCreateEvent             RequestAction = "CREATE_EVENT"
	RequestUpdateEvent             RequestAction = "UPDATE_EVENT"
	RequestDeleteEvent             RequestAction = "DELETE_EVENT"
	RequestEnableEvent             RequestAction = "ENABLE_EVENT"
	RequestDisableEvent            RequestAction = "DISABLE_EVENT"

	RequestUpdateBlockTime RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar    RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar    RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar    RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar   RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar    RequestAction = "UPDATE_AVATAR"
	RequestBlockRegion     RequestAction = "BLOCK_REGION"
	// RequestEnableRegion          RequestAction = "ENABLE_REGION"
	RequestBlockDistrict RequestAction = "BLOCK_DISTRICT"
	// RequestEnableDistrict        RequestAction = "ENABLE_DISTRICT"
	// RequestEnableCity            RequestAction = "ENABLE_CITY"
	RequestBlockCity RequestAction = "BLOCK_CITY"

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

	RequestCreateEventCategory   RequestAction = "CREATE_EVENT_CATEGORY"
	RequestUpdateEventCategory   RequestAction = "UPDATE_EVENT_CATEGORY"
	RequestBlockUser             RequestAction = "BLOCK_USER"
	RequestEnableSingleBranches  RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches   RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches  RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"
	RequestDisableFaydaAccount   RequestAction = "DISABLE_FAYDA_ACCOUNT"
	RequestEnableFaydaAccount    RequestAction = "ENABLE_FAYDA_ACCOUNT"

	RequestAccountUpdate         RequestAction = "REQUEST_ACCOUNT_UPDATE"
	RequestCreatePermissionGroup RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup RequestAction = "DELETE_PERMISSION_GROUP"

	RequestCreateBudgetCategory  RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestDeleteBudgetCategory  RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestUpdateBudgetCategory  RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestDisableBudgetCategory RequestAction = "DISABLE_BUDGET_CATEGORY"
	RequestEnableBudgetCategory  RequestAction = "ENABLE_BUDGET_CATEGORY"

	RequestUnlinkDevice           RequestAction = "UNLINK_DEVICE"
	RequestUnlinkUser             RequestAction = "UNLINK_USER"
	RequestUpdateHQBlockTime      RequestAction = "UPDATE_HQ_BLOCK_TIME"
	RequestUpdateHQArchiveTime    RequestAction = "UPDATE_HQ_ARCHIVE_TIME"
	RequestCreateMiniAppMerchant  RequestAction = "CREATE_MINI_APP_MERCHANT"
	RequestUpdateMiniAppMerchant  RequestAction = "UPDATE_MINI_APP_MERCHANT"
	RequestDeleteMiniAppMerchant  RequestAction = "DELETE_MINI_APP_MERCHANT"
	RequestEnableMiniAppMerchant  RequestAction = "ENABLE_MINI_APP_MERCHANT"
	RequestDisableMiniAppMerchant RequestAction = "DISABLE_MINI_APP_MERCHANT"

	RequestCreateMiniApp  RequestAction = "CREATE_MINI_APP"
	RequestUpdateMiniApp  RequestAction = "UPDATE_MINI_APP"
	RequestDeleteMiniApp  RequestAction = "DELETE_MINI_APP"
	RequestEnableMiniApp  RequestAction = "ENABLE_MINI_APP"
	RequestDisableMiniApp RequestAction = "DISABLE_MINI_APP"

	RequestCreateBudgetColor RequestAction = "BUDGET_CREATE_COLOR"
	RequestUpdateBudgetColor RequestAction = "BUDGET_UPDATE_COLOR"
	RequestDeleteBudgetColor RequestAction = "BUDGET_DELETE_COLOR"
	RequestCreateBudgetIcon  RequestAction = "BUDGET_CREATE_ICON"
	RequestUpdateBudgetIcon  RequestAction = "BUDGET_UPDATE_ICON"
	RequestDeleteBudgetIcon  RequestAction = "BUDGET_DELETE_ICON"

	RequestCreateServiceFee RequestAction = "CREATE_SERVICE_FEE"
	RequestUpdateServiceFee RequestAction = "UPDATE_SERVICE_FEE"
	RequestDeleteServiceFee RequestAction = "DELETE_SERVICE_FEE"
	RequestCreateDailyLimit RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit RequestAction = "DELETE DAILY LIMIT"

	// RequestCreateNotification     RequestAction = "CREATE_NOTIFICATION"
	// RequestUpdateNotification     RequestAction = "UPDATE_NOTIFICATION"
	// RequestDeleteNotification     RequestAction = "DELETE_NOTIFICATION"
	// RequestEnableNotification     RequestAction = "ENABLE_NOTIFICATION"
	// RequestDisableNotification    RequestAction = "DISABLE_NOTIFICATION"
	// RequestMarkNotificationAsSeen RequestAction = "MARK_NOTIFICATION_AS_SEEN"
	RequestUpdateProductCode       RequestAction = "UPDATE_PRODUCT_CODE"
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
	// for bankvault
	RequestCreateBankVault  RequestAction = "CREATE_VAULT_BANK"
	RequestUpdateBankVault  RequestAction = "UPDATE_VAULT_BANK"
	RequestDeleteBankVault  RequestAction = "DELETE_VAULT_BANK"
	RequestEnableBankVault  RequestAction = "ENABLE_VAULT_BANK"
	RequestDisAbleBankVault RequestAction = "DISABLE_VAULT_BANK"

	// for vault group category
	RequestCreateVaultGroupCategory  RequestAction = "CREATE_VAULT_GROUP_CATEGORY"
	RequestUpdateVaultGroupCategory  RequestAction = "UPDATE_VAULT_GROUP_CATEGORY"
	RequestDeleteVaultGroupCategory  RequestAction = "DELETE_VAULT_GROUP_CATEGORY"
	RequestEnableVaultGroupCategory  RequestAction = "ENABLE_VAULT_GROUP_CATEGORY"
	RequestDisAbleVaultGroupCategory RequestAction = "DISABLE_VAULT_GROUP_CATEGORY"

	RequestCreateVaultAmountTier  RequestAction = "CREATE_VAULT_AMOUNT_TIER"
	RequestUpdateVaultAmountTier  RequestAction = "UPDATE_VAULT_AMOUNT_TIER"
	RequestDeleteVaultAmountTier  RequestAction = "DELETE_VAULT_AMOUNT_TIER"
	RequestEnableVaultAmountTier  RequestAction = "ENABLE_VAULT_AMOUNT_TIER"
	RequestDisAbleVaultAmountTier RequestAction = "DISABLE_VAULT_AMOUNT_TIER"

	RequestUpdateKYCVerifier RequestAction = "UPDATE_KYC"
	RequestApproveKYC        RequestAction = "APPROVE_KYC"
	// for article
	RequestCreateArticle  RequestAction = "CREATE_ARTICLE"
	RequestUpdateArticle  RequestAction = "UPDATE_ARTICLE"
	RequestEnableArticle  RequestAction = "ENABLE_ARTICLE"
	RequestDisableArticle RequestAction = "DISABLE_ARTICLE"
	RequestDeleteArticle  RequestAction = "DELETE_ARTICLE"

	// for article category
	RequestCreateArticleCategory  RequestAction = "CREATE_ARTICLE_CATEGORY"
	RequestUpdateArticleCategory  RequestAction = "UPDATE_ARTICLE_CATEGORY"
	RequestDeleteArticleCategory  RequestAction = "DELETE_ARTICLE_CATEGORY"
	RequestEnableArticleCategory  RequestAction = "ENABLE_ARTICLE_CATEGORY"
	RequestDisableArticleCategory RequestAction = "DISABLE_ARTICLE_CATEGORY"

	// for tags
	RequestCreateNewsTag  RequestAction = "CREATE_NEWS_TAG"
	RequestUpdateNewsTag  RequestAction = "UPDATE_NEWS_TAG"
	RequestEnableNewsTag  RequestAction = "ENABLE_NEWS_TAG"
	RequestDisableNewsTag RequestAction = "DISABLE_NEWS_TAG"
	RequestDeleteNewsTag  RequestAction = "DELETE_NEWS_TAG"

	// Action Role Mapper
	RequestCreateActionRole  RequestAction = "CREATE_ACTION_ROLE"
	RequestUpdateActionRole  RequestAction = "UPDATE_ACTION_ROLE"
	RequestEnableActionRole  RequestAction = "ENABLE_ACTION_ROLE"
	RequestDisableActionRole RequestAction = "DISABLE_ACTION_ROLE"

	RequestCreateCpsActionRole  RequestAction = "CREATE_CPS_ACTION_ROLE"
	RequestUpdateCpsActionRole  RequestAction = "UPDATE_CPS_ACTION_ROLE"
	RequestEnableCpsActionRole  RequestAction = "ENABLE_CPS_ACTION_ROLE"
	RequestDisableCpsActionRole RequestAction = "DISABLE_CPS_ACTION_ROLE"

	RequestCreateShortVideo  RequestAction = "CREATE_SHORT_VIDEO"
	RequestUpdateShortVideo  RequestAction = "UPDATE_SHORT_VIDEO"
	RequestEnableShortVideo  RequestAction = "ENABLE_SHORT_VIDEO"
	RequestDisableShortVideo RequestAction = "DISABLE_SHORT_VIDEO"
	RequestDeleteShortVideo  RequestAction = "DELETE_SHORT_VIDEO"

	RequestApproveFaydaCustomer RequestAction = "APPROVE_FAYDA_CUSTOMER"

	RequestEnableDisableCustomer RequestAction = "ENABLE_DISABLE_CUSTOMER"

	RequestCreateMiniAppCategory  RequestAction = "CREATE_MINI_APP_CATEGORY"
	RequestUpdateMiniAppCategory  RequestAction = "UPDATE_MINI_APP_CATEGORY"
	RequestDeleteMiniAppCategory  RequestAction = "DELETE_MINI_APP_CATEGORY"
	RequestEnableMiniAppCategory  RequestAction = "ENABLE_MINI_APP_CATEGORY"
	RequestDisableMiniAppCategory RequestAction = "DISABLE_MINI_APP_CATEGORY"

	RequestCreateEventMerchant  RequestAction = "CREATE_EVENT_MERCHANT"
	RequestUpdateEventMerchant  RequestAction = "UPDATE_EVENT_MERCHANT"
	RequestDeleteEventMerchant  RequestAction = "DELETE_EVENT_MERCHANT"
	RequestEnableEventMerchant  RequestAction = "ENABLE_EVENT_MERCHANT"
	RequestDisableEventMerchant RequestAction = "DISABLE_EVENT_MERCHANT"

	RequestCreateAccessListSegmentation        RequestAction = "CREATE_ACCESS_LIST_SEGMENTATION"
	RequestUpdateAccessListSegmentation        RequestAction = "UPDATE_ACCESS_LIST_SEGMENTATION"
	RequestEnableDisableAccessListSegmentation RequestAction = "ENABLE_DISABLE_ACCESS_LIST_SEGMENTATION"
)

var validRequestActions = map[RequestAction]struct{}{
	RequestCreateMiniappProductCode:  {},
	RequestUpdateMiniappProductCode:  {},
	RequestDeleteMiniappProductCode:  {},
	RequestEnableMiniappProductCode:  {},
	RequestDisableMiniappProductCode: {},
	RequestCreateBankVault:           {},
	RequestUpdateBankVault:           {},
	RequestDeleteBankVault:           {},
	RequestEnableBankVault:           {},
	RequestDisAbleBankVault:          {},

	// for vault group category
	RequestCreateVaultGroupCategory:  {},
	RequestUpdateVaultGroupCategory:  {},
	RequestDeleteVaultGroupCategory:  {},
	RequestEnableVaultGroupCategory:  {},
	RequestDisAbleVaultGroupCategory: {},

	RequestCreateDonationCategory: {},
	RequestUpdateDonationCategory: {},
	RequestCreateDonationCompany:  {},

	RequestDisableDonationCategory: {},
	RequestEnableDonationCompany:   {},
	RequestDisableDonationCompany:  {},
	RequestEnableDonationCategory:  {},

	RequestUpdateDonationCompany: {},
	RequestCreateDonation:        {},
	RequestUpdateDonation:        {},
	RequestUpdateDonationImage:   {},
	RequestDeleteDonationImage:   {},
	RequestAddDonationImage:      {},
	RequestEnableDonation:        {},
	RequestDisableDonation:       {},
	RequestAccountUpdate:         {},
	RequestDeleteAmountBasedAuth: {},
	RequestCreateAmountBasedAuth: {},
	RequestUpdateAmountBasedAuth: {},
	RequestUser:                  {},

	RequestUpdateAccountValidation: {},
	RequestCreateBudgetColor:       {},
	RequestUpdateBudgetColor:       {},
	RequestDeleteBudgetColor:       {},
	RequestCreateBudgetIcon:        {},
	RequestUpdateBudgetIcon:        {},
	RequestDeleteBudgetIcon:        {},

	RequestCpsUserCreate:    {},
	RequestCpsUserUpdate:    {},
	RequestCpsUserDelete:    {},
	RequestPermissionGroup:  {},
	RequestCreateDepartment: {},
	RequestUpdateDepartment: {},
	RequestEnableUser:       {},
	RequestDisableUser:      {},
	RequestBPSUser:          {},
	RequestDisableBPSUser:   {},
	RequestEnableBPSUser:    {},
	RequestUpdateUser:       {},
	RequestTotalDailyLimit:  {},
	RequestUpdateVAT:        {},
	RequestAuthTier:         {},
	RequestCreateAdvert:     {},
	RequestUpdateAdvert:     {},
	RequestEnableAdvert:     {},
	RequestDisableAdvert:    {},
	RequestDeleteAdvert:     {},
	RequestCreateBank:       {},
	RequestUpdateBank:       {},
	RequestUpdateBankLogo:   {},
	RequestEnableWallet:     {},
	RequestDisableWallet:    {},

	// Services catalog
	RequestCreateService:            {},
	RequestUpdateService:            {},
	RequestEnableService:            {},
	RequestDisableService:           {},
	RequestUpdatePasswordExpiry:     {},
	RequestCreateValidation:         {},
	RequestUpdateValidation:         {},
	RequestDeleteValidation:         {},
	RequestUpdateArchiveExpiry:      {},
	RequestBudgetColor:              {},
	RequestBudgetIcon:               {},
	RequestUpdateProduct:            {},
	RequestCreatePublicNotification: {},
	RequestUpdatePublicNotification: {},
	RequestDeleteNotification:       {},
	RequestEnableNotification:       {},
	RequestDisableNotification:      {},
	RequestMarkNotificationAsSeen:   {},
	RequestArchiveUser:              {},
	RequestUpdatePasswordRule:       {},
	RequestUpdateMinimumService:     {},
	RequestUpdateServiceRule:        {},
	RequestUpdateTotal:              {},
	RequestUpdateAccessConfig:       {},

	RequestEnableSingleBranch:  {},
	RequestDisableSingleBranch: {},

	RequestEnableMultiUsers:  {},
	RequestDisableMultiUsers: {},

	RequestBlockUser: {},

	RequestDisableMultiBranches: {},
	RequestEnableMultiBranches:  {},

	// Branch
	RequestEnableBranches:  {},
	RequestDisableBranches: {},

	RequestCreateCustomerSegmentation: {},
	RequestUpdateCustomerSegmentation: {},
	RequestDeleteCustomerSegmentation: {},

	// Region
	RequestEnableRegions:  {},
	RequestDisableRegions: {},

	// District
	RequestEnableDistricts:  {},
	RequestDisableDistricts: {},

	// City
	RequestEnableCities:  {},
	RequestDisableCities: {},

	RequestCreateBusiness:         {},
	RequestUpdateBusiness:         {},
	RequestCreateEvent:            {},
	RequestUpdateEvent:            {},
	RequestCreateEventCategory:    {},
	RequestUpdateEventCategory:    {},
	RequestDisableEvent:           {},
	RequestCreateMiniAppMerchant:  {},
	RequestUpdateMiniAppMerchant:  {},
	RequestUpdateBlockTime:        {},
	RequestDisableFaydaAccount:    {},
	RequestEnableFaydaAccount:     {},
	RequestCreateAvatar:           {},
	RequestDeleteAvatar:           {},
	RequestDisableAvatar:          {},
	RequestEnableAvatar:           {},
	RequestUpdateAvatar:           {},
	RequestEnableMiniAppMerchant:  {},
	RequestDisableMiniAppMerchant: {},
	RequestDeleteMiniAppMerchant:  {},
	RequestCreateMiniApp:          {},
	RequestUpdateMiniApp:          {},
	RequestEnableMiniApp:          {},
	RequestDisableMiniApp:         {},
	RequestDeleteMiniApp:          {},
	RequestDeleteEvent:            {},
	RequestEnableEvent:            {},

	RequestUpdateServiceSingle: {},
	RequestUpdateServiceTotal:  {},
	RequestUpdateServiceMinCap: {},
	RequestCreateServiceFee:    {},
	RequestUpdateServiceFee:    {},
	RequestDeleteServiceFee:    {},
	RequestCreateDailyLimit:    {},
	RequestUpdateDailyLimit:    {},
	RequestDeleteDailyLimit:    {},
	RequestCreateWallet:        {},
	RequestUpdateWallet:        {},
	RequestDeleteWallet:        {},

	RequestCreateTopup:  {},
	RequestUpdateTopup:  {},
	RequestDeleteTopup:  {},
	RequestEnableTopup:  {},
	RequestDisableTopup: {},

	RequestCreateJobRole:  {},
	RequestUpdateJobRole:  {},
	RequestDeleteJobRole:  {},
	RequestEnableJobRole:  {},
	RequestDisableJobRole: {},

	RequestCreateRole:  {},
	RequestUpdateRole:  {},
	RequestDeleteRole:  {},
	RequestEnableRole:  {},
	RequestDisableRole: {},

	// RequestCreateNotification:      {},
	// RequestUpdateNotification:      {},
	// RequestDeleteNotification:      {},
	// RequestEnableNotification:      {},
	// RequestDisableNotification:     {},
	// RequestMarkNotificationAsSeen:  {},
	RequestUpdateProductCode:       {},
	RequestEnableDisableDepartment: {},
	// RequestEnableBranches:          {},
	RequestCpsUserEnable:  {},
	RequestCpsUserDisable: {},
	// RequestCpsUserDelete:{},
	// RequestCpsUserCreate:{},

	// for article
	RequestCreateArticle:          {},
	RequestUpdateArticle:          {},
	RequestEnableArticle:          {},
	RequestDisableArticle:         {},
	RequestDeleteArticle:          {},
	RequestCreateArticleCategory:  {},
	RequestUpdateArticleCategory:  {},
	RequestDeleteArticleCategory:  {},
	RequestEnableArticleCategory:  {},
	RequestDisableArticleCategory: {},

	// for customer
	RequestEnableDisableCustomer: {},

	RequestCreateNewsTag:  {},
	RequestUpdateNewsTag:  {},
	RequestEnableNewsTag:  {},
	RequestDisableNewsTag: {},
	RequestDeleteNewsTag:  {},

	// Action Role Mapper
	RequestCreateActionRole:     {},
	RequestUpdateActionRole:     {},
	RequestEnableActionRole:     {},
	RequestDisableActionRole:    {},
	RequestCreateCpsActionRole:  {},
	RequestUpdateCpsActionRole:  {},
	RequestEnableCpsActionRole:  {},
	RequestDisableCpsActionRole: {},

	RequestCreateDeviceVersion:        {},
	RequestUpdateDeviceVersion:        {},
	RequestEnableDeviceVersion:        {},
	RequestDisableDeviceVersion:       {},
	RequestDeleteDeviceVersion:        {},
	RequestEnableDisableDeviceVersion: {},

	RequestCreateMiniAppCategory:  {},
	RequestUpdateMiniAppCategory:  {},
	RequestDeleteMiniAppCategory:  {},
	RequestEnableMiniAppCategory:  {},
	RequestDisableMiniAppCategory: {},

	RequestCreateEventMerchant:  {},
	RequestUpdateEventMerchant:  {},
	RequestDeleteEventMerchant:  {},
	RequestEnableEventMerchant:  {},
	RequestDisableEventMerchant: {},

	RequestCreateAccessListSegmentation:        {},
	RequestUpdateAccessListSegmentation:        {},
	RequestEnableDisableAccessListSegmentation: {},
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[RequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]RequestAction{
	// "Service": {
	// 	RequestUpdateServiceSingle,
	// 	RequestUpdateServiceTotal,
	// 	RequestUpdateServiceMinCap,
	// 	RequestCreateServiceFee,
	// 	RequestUpdateServiceFee,
	// 	RequestDeleteServiceFee,
	// 	RequestCreateDailyLimit,
	// 	RequestUpdateDailyLimit,
	// 	RequestDeleteDailyLimit,
	// },
	"Service": {
		RequestCreateService,
		RequestUpdateService,
		RequestEnableService,
		RequestDisableService,
	},
	"Account": {
		RequestUser,
		RequestUpdateAccountValidation,
		RequestEnableUser,
		RequestDisableUser,
		RequestUpdateUser,
		RequestArchiveUser,
	},
	"Event": {
		RequestCreateEvent,
		RequestDeleteEvent,
		RequestDisableEvent,
		RequestEnableEvent,
		RequestUpdateEvent,
	},
	"AmountBasedAuth": {
		RequestCreateAmountBasedAuth,
		RequestUpdateAmountBasedAuth,
		RequestDeleteAmountBasedAuth,
		RequestAuthTier,
	},
	"User": {
		RequestUser,
		RequestEnableUser,
		RequestDisableUser,
		RequestUpdateUser,
		RequestArchiveUser,
	},
	"BPSUser": {
		RequestBPSUser,
		RequestEnableBPSUser,
		RequestDisableBPSUser,
	},
	"PermissionGroup": {
		RequestPermissionGroup,
	},
	"Department": {
		RequestCreateDepartment,
		RequestUpdateDepartment,
		RequestEnableDisableDepartment,
		// RequestDeleteDepartment,
	},
	"ServiceFee": {
		RequestCreateServiceFee,
		RequestUpdateServiceFee,
		RequestDeleteServiceFee,
	},
	"DailyLimit": {
		RequestCreateDailyLimit,
		RequestUpdateDailyLimit,
		RequestDeleteDailyLimit,
		RequestTotalDailyLimit,
	},
	"VAT": {
		RequestUpdateVAT,
	},
	"AuthTier": {
		RequestAuthTier,
	},
	"Archive": {
		RequestUpdateArchiveExpiry,
	},
	"MinimumService": {
		RequestUpdateMinimumService,
	},
	"ServiceRule": {
		RequestUpdateServiceRule,
	},
	"Total": {
		RequestUpdateTotal,
	},
	"AccessConfig": {
		RequestUpdateAccessConfig,
	},
	"Branch": {
		RequestEnableSingleBranch,
		RequestEnableMultiUsers,
		RequestDisableMultiUsers,
		RequestEnableSingleBranches,
		RequestEnableMultiBranches,
	},
	"Business": {
		RequestCreateBusiness,
		RequestUpdateBusiness,
	},
	"EventCategory": {
		RequestCreateEventCategory,
		RequestUpdateEventCategory,
	},
	"MiniAppMerchant": {
		RequestCreateMiniAppMerchant,
		RequestUpdateMiniAppMerchant,
		RequestDeleteMiniAppMerchant,
		RequestEnableMiniAppMerchant,
		RequestDisableMiniAppMerchant,
	},
	"BlockTime": {
		RequestUpdateBlockTime,
	},
	"Password": {

		RequestUpdatePasswordRule,
	},
	"Permission": {
		RequestCreatePermissionGroup,
		RequestDeletePermissionGroup,
		RequestUpdatePermissionGroup,
	},
	"Avatar": {
		RequestCreateAvatar,
		RequestUpdateAvatar,
		RequestEnableAvatar,
		RequestDisableAvatar,
		RequestDeleteAvatar,
	},
	"Budget": {
		RequestCreateBudgetColor,
		RequestUpdateBudgetColor,
		RequestDeleteBudgetColor,
		RequestCreateBudgetIcon,
		RequestUpdateBudgetIcon,
		RequestDeleteBudgetIcon,
	},

	"Advert": {
		RequestCreateAdvert,
		RequestUpdateAdvert,
		RequestEnableAdvert,
		RequestDisableAdvert,
		RequestDeleteAdvert,
	},
	"Bank": {
		RequestCreateBank,
		RequestUpdateBank,
		RequestDeleteBank,
		RequestEnableDisableBank,
		RequestUpdateBankLogo,
		RequestEnableBank,
		RequestDisableBank,
	},
	"Wallet": {
		RequestCreateWallet,
		RequestUpdateWallet,
		RequestDeleteWallet,
		RequestEnableWallet,
		RequestDisableWallet,
	},
	"Topup": {
		RequestCreateTopup,
		RequestUpdateTopup,
		RequestDeleteTopup,
		RequestEnableTopup,
		RequestDisableTopup,
	},
	"Validation": {
		RequestCreateValidation,
		RequestUpdateValidation,
		RequestDeleteValidation,
	},
	"JobRole": {
		RequestCreateJobRole,
		RequestUpdateJobRole,
		RequestDeleteJobRole,
		RequestEnableJobRole,
		RequestDisableJobRole,
	},
	"Role": {
		RequestCreateRole,
		RequestUpdateRole,
		RequestDeleteRole,
		RequestEnableRole,
		RequestDisableRole,
	},
	"Block": {
		RequestBlockUser,
		RequestDisableSingleBranch,
		RequestEnableSingleBranch,

		RequestDisableMultiBranches,
		RequestEnableMultiBranches,

		// Branch
		RequestEnableBranches,
		RequestDisableBranches,

		// Region
		RequestEnableRegions,
		RequestDisableRegions,

		// District
		RequestEnableDistricts,
		RequestDisableDistricts,

		// City
		RequestEnableCities,
		RequestDisableCities,
	},
	"BudgetCategory": {
		RequestCreateBudgetCategory,
		RequestUpdateBudgetCategory,
		RequestDeleteBudgetCategory,
		RequestDisableBudgetCategory,
		RequestEnableBudgetCategory,
	},
	"UnlinkDevice": {
		RequestUnlinkDevice,
		RequestUnlinkUser,
	},
	"HQ": {
		RequestUpdateHQBlockTime,
		RequestUpdateHQArchiveTime,
		RequestUpdatePasswordExpiry,
	},
	"Fayda": {
		RequestDisableFaydaAccount,
		RequestEnableFaydaAccount,
	},
	"CPSUser": {
		RequestCpsUserCreate,
		RequestCpsUserUpdate,
		RequestCpsUserDelete,
		RequestCpsUserEnable,
		RequestCpsUserDisable,
	},
	"BulkService": {
		RequestBulkServiceEnable,
		RequestBulkServiceDisable,
	},
	"MiniApp": {
		RequestCreateMiniApp,
		RequestUpdateMiniApp,
		RequestDeleteMiniApp,
		RequestEnableMiniApp,
		RequestDisableMiniApp,
	},
	"Notification": {
		// RequestCreateNotification,
		// RequestUpdateNotification,
		RequestCreatePublicNotification,
		RequestUpdatePublicNotification,
		RequestDeleteNotification,
		RequestEnableNotification,
		RequestDisableNotification,
		RequestMarkNotificationAsSeen,
	},
	"BankVault": {
		RequestCreateBankVault,
		RequestUpdateBankVault,
		RequestDeleteBankVault,
		RequestEnableBankVault,
		RequestDisAbleBankVault,
	},
	"ProductCode": {
		RequestUpdateProductCode,
	},
	"Donation": {
		RequestCreateDonation,
		RequestUpdateDonation,
		RequestDisableDonation,
		RequestAddDonationImage,
		RequestUpdateDonationImage,
		RequestDeleteDonationImage,
		RequestEnableDonation,
	},
	"donationCategory": {
		RequestCreateDonationCategory,
		RequestUpdateDonationCategory,
		RequestDisableDonationCategory,
		RequestEnableDonationCategory,
	},
	"donationCompany": {
		RequestCreateDonationCompany,
		RequestUpdateDonationCompany,
		RequestEnableDonationCompany,
		RequestDisableDonationCompany,
	},
	"VaultGroupCategory": {
		RequestCreateVaultGroupCategory,
		RequestUpdateVaultGroupCategory,
		RequestDeleteVaultGroupCategory,
		RequestEnableVaultGroupCategory,
		RequestDisAbleVaultGroupCategory,
	},
	"KYCVerifier": {
		RequestUpdateKYCVerifier,
		RequestApproveKYC,
	},
	"article": {
		RequestCreateArticle,
		RequestUpdateArticle,
		RequestDeleteArticle,
		RequestEnableArticle,
		RequestDisableArticle,
	},
	"articleCategory": {
		RequestCreateArticleCategory,
		RequestUpdateArticleCategory,
		RequestDeleteArticleCategory,
		RequestEnableArticleCategory,
		RequestDisableArticleCategory,
	},
	"short_video": {
		RequestCreateShortVideo,
		RequestUpdateShortVideo,
		RequestDeleteShortVideo,
		RequestEnableShortVideo,
		RequestDisableShortVideo,
	},
	"customer": {
		RequestEnableDisableCustomer,
		RequestApproveFaydaCustomer,
	},

	"news_category": {
		RequestAction("CREATE_NEWS_CATEGORY"),
		RequestAction("UPDATE_NEWS_CATEGORY"),
		RequestAction("DELETE_NEWS_CATEGORY"),
	},
	"news_tag": {
		RequestCreateNewsTag,
		RequestUpdateNewsTag,
		RequestEnableNewsTag,
		RequestDisableNewsTag,
		RequestDeleteNewsTag,
	},
	"ActionRole": {
		RequestCreateActionRole,
		RequestUpdateActionRole,
		RequestEnableActionRole,
		RequestDisableActionRole,
	},
	"CpsActionRole": {
		RequestCreateCpsActionRole,
		RequestUpdateCpsActionRole,
		RequestEnableCpsActionRole,
		RequestDisableCpsActionRole,
	},
	"DeviceVersion": {
		RequestCreateDeviceVersion,
		RequestUpdateDeviceVersion,
		RequestEnableDeviceVersion,
		RequestDisableDeviceVersion,
		RequestDeleteDeviceVersion,
		RequestEnableDisableDeviceVersion,
	},
	"MiniAppCategory": {
		RequestCreateMiniAppCategory,
		RequestUpdateMiniAppCategory,
		RequestDeleteMiniAppCategory,
		RequestEnableMiniAppCategory,
		RequestDisableMiniAppCategory,
	},

	"eventMerchant": {
		RequestCreateEventMerchant,
		RequestUpdateEventMerchant,
		RequestDeleteEventMerchant,
		RequestEnableEventMerchant,
		RequestDisableEventMerchant,
	},

	"VaultAmountTier": {
		RequestCreateVaultAmountTier,
		RequestUpdateVaultAmountTier,
		RequestDeleteVaultAmountTier,
		RequestEnableVaultAmountTier,
		RequestDisAbleVaultAmountTier,
	},
	"MiniAppProductCode": {
		RequestCreateMiniappProductCode,
		RequestUpdateMiniappProductCode,
		RequestDeleteMiniappProductCode,
		RequestEnableMiniappProductCode,
		RequestDisableMiniappProductCode,
	},

	"AccessListSegmentation": {
		RequestCreateAccessListSegmentation,
		RequestUpdateAccessListSegmentation,
		RequestEnableDisableAccessListSegmentation,
	},
	"CustomerSegmentations": {
		RequestCreateCustomerSegmentation,
		RequestUpdateCustomerSegmentation,
		RequestDeleteCustomerSegmentation,
	},
}

func IsActionInGroup(action RequestAction, group string) bool {
	actions, exists := RequestActionGroups[group]
	if !exists {
		return false
	}

	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}
