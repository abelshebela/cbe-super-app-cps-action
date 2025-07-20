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
	RequestCreateServiceFee         RequestAction = "SERVICE_FEE_CREATE"
	RequestUpdateServiceFee         RequestAction = "SERVICE_FEE_UPDATE"
	RequestDeleteServiceFee         RequestAction = "SERVICE_FEE_DELETE"
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
	RequestCreateMiniAppMerchant    RequestAction = "CREATE_MINIAPP_MERCHANT"
	RequestUpdateMiniAppMerchant    RequestAction = "UPDATE_MINIAPP_MERCHANT"
	RequestDeleteMiniAppMerchant    RequestAction = "DELETE_MINI_APP"
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar             RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar             RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar             RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar            RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar             RequestAction = "UPDATE_AVATAR"
	RequestBlockRegion              RequestAction = "BLOCK_REGION"
	RequestBlockDistrict            RequestAction = "BLOCK_DISTRICT"
	RequestBlockCity                RequestAction = "BLOCK_CITY"
	RequestBlockUser                RequestAction = "BLOCK_USER"
	RequestEnableSingleBranches     RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches    RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches      RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches     RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"
	RequestDisableFaydaAccount      RequestAction = "DISABLE_FAYDA_ACCOUNT"
	RequestCreatePermissionGroup    RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup    RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup    RequestAction = "DELETE_PERMISSION_GROUP"
	RequestCreateBudgetCategory     RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestDeleteBudgetCategory     RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestUpdateBudgetCategory     RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestUnlinkDevice             RequestAction = "UNLINK_DEVICE"
)

var validRequestActions = map[RequestAction]struct{}{
	RequestUser:                     {},
	RequestPermissionGroup:          {},
	RequestDepartment:               {},
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
	RequestCreateMiniAppMerchant:    {},
	RequestUpdateMiniAppMerchant:    {},
	RequestDeleteMiniAppMerchant:    {},
	RequestUpdateBlockTime:          {},
	RequestCreatePermissionGroup:    {},
	RequestUpdatePermissionGroup:    {},
	RequestDeletePermissionGroup:    {},
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[RequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]RequestAction{
	"AccessConfig": {
		RequestUpdateAccessConfig,
	},
	"Advert": {
		RequestCreateAdvert,
		RequestDeleteAdvert,
		RequestDisableAdvert,
		RequestEnableAdvert,
		RequestUpdateAdvert,
	},
	"Archive": {
		RequestUpdateArchiveExpiry,
	},
	"AuthTier": {
		RequestAuthTier,
	},
	"Avatar": {
		RequestCreateAvatar,
		RequestDeleteAvatar,
		RequestDisableAvatar,
		RequestEnableAvatar,
		RequestUpdateAvatar,
	},
	"Bank": {
		RequestCreateBank,
		RequestDeleteBank,
		RequestDisableBank,
		RequestEnableBank,
		RequestUpdateBank,
	},
	"Block": {
		RequestBlockCity,
		RequestBlockDistrict,
		RequestBlockRegion,
		RequestBlockUser,
		RequestDisableMultiBranches,
		RequestDisableSingleBranch,
	},
	"BlockTime": {
		RequestUpdateBlockTime,
	},
	"BPSUser": {
		RequestBPSUser,
		RequestDisableBPSUser,
		RequestEnableBPSUser,
	},
	"Branch": {
		RequestDisableMultiUsers,
		RequestEnableMultiBranches,
		RequestEnableMultiUsers,
		RequestEnableSingleBranch,
		RequestEnableSingleBranches,
	},
	"Budget": {
		RequestBudgetColor,
		RequestBudgetIcon,
	},
	"BudgetCategory": {
		RequestCreateBudgetCategory,
		RequestUpdateBudgetCategory,
		RequestDeleteBudgetCategory,
	},
	"Business": {
		RequestCreateBusiness,
		RequestUpdateBusiness,
	},
	"DailyLimit": {
		RequestCreateDailyLimit,
		RequestDeleteDailyLimit,
		RequestTotalDailyLimit,
		RequestUpdateDailyLimit,
	},
	"Department": {
		RequestDepartment,
	},
	"Event": {
		RequestCreateEvent,
		RequestDisableEvent,
		RequestUpdateEvent,
	},
	"EventCategory": {
		RequestCreateEventCategory,
		RequestUpdateEventCategory,
	},
	"MiniAppMerchant": {
		RequestCreateMiniAppMerchant,
		RequestDeleteMiniAppMerchant,
		RequestUpdateMiniAppMerchant,
	},
	"MinimumService": {
		RequestUpdateMinimumService,
	},
	"Password": {
		RequestUpdatePasswordExpiry,
		RequestUpdatePasswordRule,
	},
	"Permission": {
		RequestCreatePermissionGroup,
		RequestDeletePermissionGroup,
		RequestUpdatePermissionGroup,
	},
	"PermissionGroup": {
		RequestPermissionGroup,
	},
	"Product": {
		RequestUpdateProduct,
	},
	"PublicNotification": {
		RequestCreatePublicNotification,
	},
	"ServiceFee": {
		RequestCreateServiceFee,
		RequestDeleteServiceFee,
		RequestUpdateServiceFee,
	},
	"ServiceRule": {
		RequestUpdateServiceRule,
	},
	"Total": {
		RequestUpdateTotal,
	},
	"User": {
		RequestArchiveUser,
		RequestDisableUser,
		RequestEnableUser,
		RequestUpdateUser,
		RequestUser,
	},
	"Validation": {
		RequestCreateValidation,
		RequestDeleteValidation,
		RequestUpdateValidation,
	},
	"VAT": {
		RequestUpdateVAT,
	},
	"Wallet": {
		RequestCreateWallet,
		RequestDeleteWallet,
		RequestDisableWallet,
		RequestEnableWallet,
		RequestUpdateWallet,
	},
	"UnlinkDevice": {
		RequestUnlinkDevice,
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
