package constant

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
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
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
	RequestPermissionGroup       RequestAction = "PERMISSION_GROUP"
	// RequestDepartment               RequestAction = "DEPARMTENT"
	RequestCreateDepartment         RequestAction = "CREATE_DEPARTMENT"
	RequestUpdateDepartment         RequestAction = "UPDATE_DEPARTMENT"
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

	RequestUpdateBlockTime       RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar          RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar          RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar          RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar         RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar          RequestAction = "UPDATE_AVATAR"
	RequestBlockRegion           RequestAction = "BLOCK_REGION"
	RequestBlockDistrict         RequestAction = "BLOCK_DISTRICT"
	RequestBlockCity             RequestAction = "BLOCK_CITY"
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
)

var validRequestActions = map[RequestAction]struct{}{
	RequestAccountUpdate:         {},
	RequestDeleteAmountBasedAuth: {},
	RequestCreateAmountBasedAuth: {},
	RequestUpdateAmountBasedAuth: {},
	RequestUser:                  {},

	RequestCreateBudgetColor: {},
	RequestUpdateBudgetColor: {},
	RequestDeleteBudgetColor: {},
	RequestCreateBudgetIcon:  {},
	RequestUpdateBudgetIcon:  {},
	RequestDeleteBudgetIcon:  {},

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
	RequestEnableBank:               {},
	RequestDisableBank:              {},
	RequestEnableWallet:             {},
	RequestDisableWallet:            {},
	RequestUpdatePasswordExpiry:     {},
	RequestCreateValidation:         {},
	RequestUpdateValidation:         {},
	RequestDeleteValidation:         {},
	RequestUpdateArchiveExpiry:      {},
	RequestCreateServiceFee:         {},
	RequestUpdateServiceFee:         {},
	RequestDeleteServiceFee:         {},
	RequestCreateDailyLimit:         {},
	RequestUpdateDailyLimit:         {},
	RequestDeleteDailyLimit:         {},
	RequestBudgetColor:              {},
	RequestBudgetIcon:               {},
	RequestUpdateProduct:            {},
	RequestCreatePublicNotification: {},
	RequestArchiveUser:              {},
	RequestUpdatePasswordRule:       {},
	RequestUpdateMinimumService:     {},
	RequestUpdateServiceRule:        {},
	RequestUpdateTotal:              {},
	RequestUpdateAccessConfig:       {},
	RequestEnableSingleBranch:       {},
	RequestDisableSingleBranch:      {},
	RequestEnableMultiUsers:         {},
	RequestDisableMultiUsers:        {},
	RequestCreateBusiness:           {},
	RequestUpdateBusiness:           {},
	RequestCreateEvent:              {},
	RequestUpdateEvent:              {},
	RequestCreateEventCategory:      {},
	RequestUpdateEventCategory:      {},
	RequestDisableEvent:             {},
	RequestUpdateBlockTime:          {},
	RequestDisableFaydaAccount:      {},
	RequestEnableFaydaAccount:       {},
	RequestCreateAvatar:             {},
	RequestDeleteAvatar:             {},
	RequestDisableAvatar:            {},
	RequestEnableAvatar:             {},
	RequestUpdateAvatar:             {},
	RequestCreateMiniAppMerchant:    {},
	RequestUpdateMiniAppMerchant:    {},
	RequestEnableMiniAppMerchant:    {},
	RequestDisableMiniAppMerchant:   {},
	RequestDeleteMiniAppMerchant:    {},
	RequestCreateMiniApp:            {},
	RequestUpdateMiniApp:            {},
	RequestEnableMiniApp:            {},
	RequestDisableMiniApp:           {},
	RequestDeleteMiniApp:            {},
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[RequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]RequestAction{
	"Account": {
		RequestUser,
		RequestEnableUser,
		RequestDisableUser,
		RequestUpdateUser,
		RequestArchiveUser,
	},
	"A": {
		RequestCreateAvatar,
		RequestDeleteAvatar,
		RequestDisableAvatar,
		RequestEnableAvatar,
		RequestUpdateAvatar,
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
	"Event": {
		RequestCreateEvent,
		RequestUpdateEvent,
		RequestDisableEvent,
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
		RequestBlockRegion,
		RequestBlockDistrict,
		RequestBlockCity,
	},
	"BudgetCategory": {
		RequestAction("CREATE_BUDGET_CATEGORY"),
		RequestAction("UPDATE_BUDGET_CATEGORY"),
		RequestAction("DELETE_BUDGET_CATEGORY"),
	},
	"UnlinkDevice": {
		RequestUnlinkDevice,
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
	},
	"MiniApp": {
		RequestCreateMiniApp,
		RequestUpdateMiniApp,
		RequestDeleteMiniApp,
		RequestEnableMiniApp,
		RequestDisableMiniApp,
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
