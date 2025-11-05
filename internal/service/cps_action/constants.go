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
	RequestCreateDepartment        RequestAction = "CREATE_DEPARTMENT"
	RequestUpdateDepartment        RequestAction = "UPDATE_DEPARTMENT"
	RequestEnableDisableDepartment RequestAction = "ENABLE_DISABLE_DEPARTMENT"
	RequestEnableUser              RequestAction = "ENABLE_USER"
	RequestDisableUser             RequestAction = "DISABLE_USER"
	RequestBPSUser                 RequestAction = "BPS_USER"
	RequestDisableBPSUser          RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser           RequestAction = "ENABLE_BPS_USER"
	RequestUpdateUser              RequestAction = "UPDATE_USER"
	RequestTotalDailyLimit         RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT               RequestAction = "UPDATE_VAT"
	RequestAuthTier                RequestAction = "AUTHTIER"
	RequestCreateAdvert            RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert            RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert            RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert           RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert            RequestAction = "DELETE_ADVERT"
	RequestCreateBank              RequestAction = "CREATE_BANK"
	RequestUpdateBank              RequestAction = "UPDATE_BANK"
	RequestDeleteBank              RequestAction = "DELETE_BANK"
	RequestUpdateBankLogo          RequestAction = "UPDATE_BANK_LOGO"
	RequestEnableDisableBank       RequestAction = "ENABLE_DISABLE_BANK"
	RequestEnableBank              RequestAction = "ENABLE_BANK"
	RequestDisableBank             RequestAction = "DISABLE_BANK"
	RequestCreateWallet            RequestAction = "CREATE_WALLET"
	RequestUpdateWallet            RequestAction = "UPDATE_WALLET"
	RequestDeleteWallet            RequestAction = "DELETE_WALLET"
	RequestEnableWallet            RequestAction = "ENABLE_WALLET"
	RequestDisableWallet           RequestAction = "DISABLE_WALLET"

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

	RequestAccountUpdate          RequestAction = "REQUEST_ACCOUNT_UPDATE"
	RequestCreatePermissionGroup  RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup  RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup  RequestAction = "DELETE_PERMISSION_GROUP"
	RequestCreateBudgetCategory   RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestDeleteBudgetCategory   RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestUpdateBudgetCategory   RequestAction = "UPDATE_BUDGET_CATEGORY"
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
	RequestCreateBankVault  RequestAction = "CREATE VAULT BANK"
	RequestUpdateBankVault  RequestAction = "UPDATE VAULT BANK"
	RequestDeleteBankVault  RequestAction = "DELETE VAULT BANK"
	RequestEnableBankVault  RequestAction = "ENABLE VAULT BANK"
	RequestDisAbleBankVault RequestAction = "DISABLE VAULT BANK"

	// for vault group category
	RequestCreateVaultGroupCategory  RequestAction = "CREATE VAULT GROUP CATEGORY"
	RequestUpdateVaultGroupCategory  RequestAction = "UPDATE VAULT GROUP CATEGORY"
	RequestDeleteVaultGroupCategory  RequestAction = "DELETE VAULT GROUP CATEGORY"
	RequestEnableVaultGroupCategory  RequestAction = "ENABLE VAULT GROUP CATEGORY"
	RequestDisAbleVaultGroupCategory RequestAction = "DISABLE VAULT GROUP CATEGORY"

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

	RequestCreateShortVideo  RequestAction = "CREATE_SHORT_VIDEO"
	RequestUpdateShortVideo  RequestAction = "UPDATE_SHORT_VIDEO"
	RequestEnableShortVideo  RequestAction = "ENABLE_SHORT_VIDEO"
	RequestDisableShortVideo RequestAction = "DISABLE_SHORT_VIDEO"
	RequestDeleteShortVideo  RequestAction = "DELETE_SHORT_VIDEO"

	RequestApproveFaydaCustomer RequestAction = "APPROVE_FAYDA_CUSTOMER"

	RequestEnableDisableCustomer RequestAction = "ENABLE_DISABLE_CUSTOMER"
)

var validRequestActions = map[RequestAction]struct{}{
	RequestCreateBankVault:  {},
	RequestUpdateBankVault:  {},
	RequestDeleteBankVault:  {},
	RequestEnableBankVault:  {},
	RequestDisAbleBankVault: {},

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

	RequestCpsUserCreate:            {},
	RequestCpsUserUpdate:            {},
	RequestCpsUserDelete:            {},
	RequestPermissionGroup:          {},
	RequestCreateDepartment:         {},
	RequestUpdateDepartment:         {},
	RequestEnableUser:               {},
	RequestDisableUser:              {},
	RequestBPSUser:                  {},
	RequestDisableBPSUser:           {},
	RequestEnableBPSUser:            {},
	RequestUpdateUser:               {},
	RequestTotalDailyLimit:          {},
	RequestUpdateVAT:                {},
	RequestAuthTier:                 {},
	RequestCreateAdvert:             {},
	RequestUpdateAdvert:             {},
	RequestEnableAdvert:             {},
	RequestDisableAdvert:            {},
	RequestDeleteAdvert:             {},
	RequestCreateBank:               {},
	RequestUpdateBank:               {},
	RequestUpdateBankLogo:           {},
	RequestEnableWallet:             {},
	RequestDisableWallet:            {},
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

	// RequestCreateNotification:      {},
	// RequestUpdateNotification:      {},
	// RequestDeleteNotification:      {},
	// RequestEnableNotification:      {},
	// RequestDisableNotification:     {},
	// RequestMarkNotificationAsSeen:  {},
	RequestUpdateProductCode:       {},
	RequestEnableDisableDepartment: {},
	RequestEnableBranches:          {},
	RequestCpsUserEnable:           {},
	RequestCpsUserDisable:          {},
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
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[RequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]RequestAction{
	"Service": {
		RequestUpdateServiceSingle,
		RequestUpdateServiceTotal,
		RequestUpdateServiceMinCap,
		RequestCreateServiceFee,
		RequestUpdateServiceFee,
		RequestDeleteServiceFee,
		RequestCreateDailyLimit,
		RequestUpdateDailyLimit,
		RequestDeleteDailyLimit,
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
		RequestAction("CREATE_BUDGET_CATEGORY"),
		RequestAction("UPDATE_BUDGET_CATEGORY"),
		RequestAction("DELETE_BUDGET_CATEGORY"),
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
		RequestAction("CREATE_NEWS_TAG"),
		RequestAction("UPDATE_NEWS_TAG"),
		RequestAction("DELETE_NEWS_TAG"),
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
