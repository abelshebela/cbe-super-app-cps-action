package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
)

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

const (
	// Customer KYC
	// RequestCreateCustomerKYC  constants.RequestAction = "CREATE_CUSTOMER_KYC"
	RequestApproveCustomerKYC constants.RequestAction = "APPROVE_CUSTOMER_KYC"
	RequestRejectCustomerKYC  constants.RequestAction = "REJECT_CUSTOMER_KYC"
	// RequestUpdateCustomerKYC  constants.RequestAction = "UPDATE_CUSTOMER_KYC"
	// RequestDeleteCustomerKYC  constants.RequestAction = "DELETE_CUSTOMER_KYC"

	RequestCreateCpsRole     constants.RequestAction = "CREATE_CPS_ROLE"
	RequestUpdateCpsRole     constants.RequestAction = "UPDATE_CPS_ROLE"
	RequestDeleteCpsRole     constants.RequestAction = "DELETE_CPS_ROLE"
	RequestEnableCpsRole     constants.RequestAction = "ENABLE_CPS_ROLE"
	RequestDisableCpsRole    constants.RequestAction = "DISABLE_CPS_ROLE"
	RequestDeleteServiceList constants.RequestAction = "DELETE_SERVICE_LIST"

	RequestCreateUssdMerchant  constants.RequestAction = "CREATE_USSD_MERCHANT"
	RequestUpdateUssdMerchant  constants.RequestAction = "UPDATE_USSD_MERCHANT"
	RequestDeleteUssdMerchant  constants.RequestAction = "DELETE_USSD_MERCHANT"
	RequestEnableUssdMerchant  constants.RequestAction = "ENABLE_USSD_MERCHANT"
	RequestDisableUssdMerchant constants.RequestAction = "DISABLE_USSD_MERCHANT"

	RequestCreateCustomerGroup  constants.RequestAction = "CREATE_CUSTOMER_GROUP"
	RequestUpdateCustomerGroup  constants.RequestAction = "UPDATE_CUSTOMER_GROUP"
	RequestEnableCustomerGroup  constants.RequestAction = "ENABLE_CUSTOMER_GROUP"
	RequestDisableCustomerGroup constants.RequestAction = "DISABLE_CUSTOMER_GROUP"
	RequestDeleteCustomerGroup  constants.RequestAction = "DELETE_CUSTOMER_GROUP"

	RequestEnableSuperAppRole  constants.RequestAction = "ENABLE_SUPERAPP_ROLE"
	RequestDisableSuperAppRole constants.RequestAction = "DISABLE_SUPERAPP_ROLE"
	RequestDeleteSuperAppRole  constants.RequestAction = "DELETE_SUPERAPP_ROLE"

	RequestCreateCustomerSegmentation  constants.RequestAction = "CREATE_CUSTOMER_SEGMENTATION"
	RequestUpdateCustomerSegmentation  constants.RequestAction = "UPDATE_CUSTOMER_SEGMENTATION"
	RequestEnableCustomerSegmentation  constants.RequestAction = "ENABLE_CUSTOMER_SEGMENTATION"
	RequestDisableCustomerSegmentation constants.RequestAction = "DISABLE_CUSTOMER_SEGMENTATION"
	RequestDeleteCustomerSegmentation  constants.RequestAction = "DELETE_CUSTOMER_SEGMENTATION"

	RequestAccessListCreateCustomerSegmentation  constants.RequestAction = "CREATE_ACCESS_LIST_CUSTOMER_SEGMENTATION"
	RequestAccessListUpdateCustomerSegmentation  constants.RequestAction = "UPDATE_ACCESS_LIST_CUSTOMER_SEGMENTATION"
	RequestAccessListEnableCustomerSegmentation  constants.RequestAction = "ENABLE_ACCESS_LIST_CUSTOMER_SEGMENTATION"
	RequestAccessListDisableCustomerSegmentation constants.RequestAction = "DISABLE_ACCESS_LIST_CUSTOMER_SEGMENTATION"
	RequestAccessListDeleteCustomerSegmentation  constants.RequestAction = "DELETE_ACCESS_LIST_CUSTOMER_SEGMENTATION"

	RequestCreateMiniappProductCode  constants.RequestAction = "CREATE_MINI_APP_PRODUCT_CODE"
	RequestUpdateMiniappProductCode  constants.RequestAction = "UPDATE_MINI_APP_PRODUCT_CODE"
	RequestDeleteMiniappProductCode  constants.RequestAction = "DELETE_MINI_APP_PRODUCT_CODE"
	RequestEnableMiniappProductCode  constants.RequestAction = "ENABLE_MINI_APP_PRODUCT_CODE"
	RequestDisableMiniappProductCode constants.RequestAction = "DISABLE_MINI_APP_PRODUCT_CODE"

	RequestCreateJobRole  constants.RequestAction = "CREATE_JOB_ROLE"
	RequestUpdateJobRole  constants.RequestAction = "UPDATE_JOB_ROLE"
	RequestDeleteJobRole  constants.RequestAction = "DELETE_JOB_ROLE"
	RequestEnableJobRole  constants.RequestAction = "ENABLE_JOB_ROLE"
	RequestDisableJobRole constants.RequestAction = "DISABLE_JOB_ROLE"

	RequestCreateRole  constants.RequestAction = "CREATE_ROLE"
	RequestUpdateRole  constants.RequestAction = "UPDATE_ROLE"
	RequestDeleteRole  constants.RequestAction = "DELETE_ROLE"
	RequestEnableRole  constants.RequestAction = "ENABLE_ROLE"
	RequestDisableRole constants.RequestAction = "DISABLE_ROLE"

	RequestDeleteAmountBasedAuth constants.RequestAction = "DELETE_AMOUNT_BASED_AUTH"
	RequestCreateAmountBasedAuth constants.RequestAction = "CREATE_AMOUNT_BASED_AUTH"
	RequestUpdateAmountBasedAuth constants.RequestAction = "UPDATE_AMOUNT_BASED_AUTH"
	RequestResetAmountBasedAuth  constants.RequestAction = "RESET_AMOUNT_BASED_AUTH"
	RequestUser                  constants.RequestAction = "USER"
	RequestCpsUserCreate         constants.RequestAction = "CREATE_CPS_USER"
	RequestCpsUserUpdate         constants.RequestAction = "UPDATE_CPS_USER"
	RequestCpsUserDelete         constants.RequestAction = "DELETE_CPS_USER"
	RequestCpsUserEnable         constants.RequestAction = "ENABLE_CPS_USER"
	RequestCpsUserDisable        constants.RequestAction = "DISABLE_CPS_USER"

	RequestBpsUserCreate  constants.RequestAction = "CREATE_BPS_USER"
	RequestBpsUserUpdate  constants.RequestAction = "UPDATE_BPS_USER"
	RequestBpsUserDelete  constants.RequestAction = "DELETE_BPS_USER"
	RequestBpsUserEnable  constants.RequestAction = "ENABLE_BPS_USER"
	RequestBpsUserDisable constants.RequestAction = "DISABLE_BPS_USER"

	RequestPermissionGroup    constants.RequestAction = "PERMISSION_GROUP"
	RequestBulkServiceEnable  constants.RequestAction = "ENABLE_BULK_SERVICE"
	RequestBulkServiceDisable constants.RequestAction = "DISABLE_BULK_SERVICE"

	// RequestDepartment               constants.RequestAction = "DEPARMTENT"
	RequestCreateDepartment           constants.RequestAction = "CREATE_DEPARTMENT"
	RequestUpdateDepartment           constants.RequestAction = "UPDATE_DEPARTMENT"
	RequestDeleteDepartment           constants.RequestAction = "DELETE_DEPARTMENT"
	RequestEnableDisableDepartment    constants.RequestAction = "ENABLE_DISABLE_DEPARTMENT"
	RequestEnableUser                 constants.RequestAction = "ENABLE_USER"
	RequestDisableUser                constants.RequestAction = "DISABLE_USER"
	RequestBPSUser                    constants.RequestAction = "BPS_USER"
	RequestDisableBPSUser             constants.RequestAction = "DISABLE_BPS_USER"
	RequestEnableBPSUser              constants.RequestAction = "ENABLE_BPS_USER"
	RequestUpdateUser                 constants.RequestAction = "UPDATE_USER"
	RequestTotalDailyLimit            constants.RequestAction = "TOTAL_DAILY_LIMIT"
	RequestUpdateVAT                  constants.RequestAction = "UPDATE_VAT"
	RequestAuthTier                   constants.RequestAction = "AUTHTIER"
	RequestCreateDeviceVersion        constants.RequestAction = "CREATE_DEVICE_VERSION"
	RequestUpdateDeviceVersion        constants.RequestAction = "UPDATE_DEVICE_VERSION"
	RequestEnableDeviceVersion        constants.RequestAction = "ENABLE_DEVICE_VERSION"
	RequestDisableDeviceVersion       constants.RequestAction = "DISABLE_DEVICE_VERSION"
	RequestDeleteDeviceVersion        constants.RequestAction = "DELETE_DEVICE_VERSION"
	RequestEnableDisableDeviceVersion constants.RequestAction = "ENABLE_DISABLE_DEVICE_VERSION"
	RequestCreateAdvert               constants.RequestAction = "CREATE_ADVERT"
	RequestUpdateAdvert               constants.RequestAction = "UPDATE_ADVERT"
	RequestEnableAdvert               constants.RequestAction = "ENABLE_ADVERT"
	RequestDisableAdvert              constants.RequestAction = "DISABLE_ADVERT"
	RequestDeleteAdvert               constants.RequestAction = "DELETE_ADVERT"
	RequestCreateBank                 constants.RequestAction = "CREATE_BANK"
	RequestUpdateBank                 constants.RequestAction = "UPDATE_BANK"
	RequestDeleteBank                 constants.RequestAction = "DELETE_BANK"
	RequestUpdateBankLogo             constants.RequestAction = "UPDATE_BANK_LOGO"
	RequestEnableDisableBank          constants.RequestAction = "ENABLE_DISABLE_BANK"
	RequestEnableBank                 constants.RequestAction = "ENABLE_BANK"
	RequestDisableBank                constants.RequestAction = "DISABLE_BANK"
	RequestCreateWallet               constants.RequestAction = "CREATE_WALLET"
	RequestUpdateWallet               constants.RequestAction = "UPDATE_WALLET"
	RequestDeleteWallet               constants.RequestAction = "DELETE_WALLET"
	RequestEnableWallet               constants.RequestAction = "ENABLE_WALLET"
	RequestDisableWallet              constants.RequestAction = "DISABLE_WALLET"
	RequestEnableWalletService        constants.RequestAction = "ENABLE_WALLET_SERVICE"
	RequestDisableWalletService       constants.RequestAction = "DISABLE_WALLET_SERVICE"

	RequestCreateEcommerceMerchant        constants.RequestAction = "CREATE_ECOMMERCE_MERCHANT"
	RequestUpdateEcommerceMerchant        constants.RequestAction = "UPDATE_ECOMMERCE_MERCHANT"
	RequestEnableEcommerceMerchant        constants.RequestAction = "ENABLE_ECOMMERCE_MERCHANT"
	RequestDisableEcommerceMerchant       constants.RequestAction = "DISABLE_ECOMMERCE_MERCHANT"
	RequestDeleteEcommerceMerchant        constants.RequestAction = "DELETE_ECOMMERCE_MERCHANT"
	RequestDeleteEcommerceMerchantBranch  constants.RequestAction = "DELETE_ECOMMERCE_MERCHANT_BRANCH"
	RequestEnableEcommerceMerchantBranch  constants.RequestAction = "ENABLE_ECOMMERCE_MERCHANT_BRANCH"
	RequestDisableEcommerceMerchantBranch constants.RequestAction = "DISABLE_ECOMMERCE_MERCHANT_BRANCH"

	// Services catalog (model.Services)
	RequestCreateService      constants.RequestAction = "CREATE_SERVICE"
	RequestUpdateService      constants.RequestAction = "UPDATE_SERVICE"
	RequestEnableService      constants.RequestAction = "ENABLE_SERVICE"
	RequestDisableService     constants.RequestAction = "DISABLE_SERVICE"
	RequestDeleteService      constants.RequestAction = "DELETE_SERVICE"
	RequestCreateServiceList  constants.RequestAction = "CREATE_SERVICE_LIST"
	RequestUpdateServiceList  constants.RequestAction = "UPDATE_SERVICE_LIST"
	RequestEnableServiceList  constants.RequestAction = "ENABLE_SERVICE_LIST"
	RequestDisableServiceList constants.RequestAction = "DISABLE_SERVICE_LIST"

	RequestCreateTopup  constants.RequestAction = "CREATE_TOPUP"
	RequestUpdateTopup  constants.RequestAction = "UPDATE_TOPUP"
	RequestDeleteTopup  constants.RequestAction = "DELETE_TOPUP"
	RequestEnableTopup  constants.RequestAction = "ENABLE_TOPUP"
	RequestDisableTopup constants.RequestAction = "DISABLE_TOPUP"

	RequestUpdatePasswordExpiry constants.RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreateValidation     constants.RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation     constants.RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation     constants.RequestAction = "DELETE_VALIDATION"
	RequestUpdateArchiveExpiry  constants.RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestUpdateServiceSingle  constants.RequestAction = "UPDATE_SERVICE_SINGLE_CAP"
	RequestUpdateServiceTotal   constants.RequestAction = "UPDATE_SERVICE_TOTAL_CAP"
	RequestUpdateServiceMinCap  constants.RequestAction = "UPDATE_SERVICE_MIN_CAP"
	RequestBudgetColor          constants.RequestAction = "BUDGET_COLOR"
	RequestBudgetIcon           constants.RequestAction = "BUDGET_ICON"
	RequestUpdateProduct        constants.RequestAction = "UPDATE_PRODUCT"

	// Notification
	RequestCreatePublicNotification constants.RequestAction = "CREATE_PUBLIC_NOTIFICATION"
	RequestUpdatePublicNotification constants.RequestAction = "UPDATE_PUBLIC_NOTIFICATION"
	RequestDeleteNotification       constants.RequestAction = "DELETE_NOTIFICATION"
	RequestEnableNotification       constants.RequestAction = "ENABLE_NOTIFICATION"
	RequestDisableNotification      constants.RequestAction = "DISABLE_NOTIFICATION"
	RequestMarkNotificationAsSeen   constants.RequestAction = "MARK_NOTIFICATION_AS_SEEN"

	RequestArchiveUser             constants.RequestAction = "ARCHIVE_USER"
	RequestUpdatePasswordRule      constants.RequestAction = "UPDATE_PASSWORD_RULE"
	RequestUpdateMinimumService    constants.RequestAction = "UPDATE_MINIMUM_SERVICE"
	RequestUpdateServiceRule       constants.RequestAction = "UPDATE_SERVICE_RULE"
	RequestUpdateTotal             constants.RequestAction = "UPDATE_TOTAL"
	RequestUpdateAccessConfig      constants.RequestAction = "UPDATE_ACCESS_CONFIG"
	RequestEnableSingleBranch      constants.RequestAction = "ENABLE_SINGLE_BRANCH"
	RequestDisableSingleBranch     constants.RequestAction = "DISABLE_SINGLE_BRANCH"
	RequestUpdateAccountValidation constants.RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	RequestEnableMultiUsers        constants.RequestAction = "ENABLE_MULTI_USERS"
	RequestDisableMultiUsers       constants.RequestAction = "DISABLE_MULTI_USERS"
	RequestCreateBusiness          constants.RequestAction = "CREATE_BUSINESS"
	RequestUpdateBusiness          constants.RequestAction = "UPDATE_BUSINESS"
	RequestCreateEvent             constants.RequestAction = "CREATE_EVENT"
	RequestUpdateEvent             constants.RequestAction = "UPDATE_EVENT"
	RequestDeleteEvent             constants.RequestAction = "DELETE_EVENT"
	RequestEnableEvent             constants.RequestAction = "ENABLE_EVENT"
	RequestDisableEvent            constants.RequestAction = "DISABLE_EVENT"

	RequestUpdateBlockTime constants.RequestAction = "UPDATE_BLOCK_TIME"
	RequestCreateAvatar    constants.RequestAction = "CREATE_AVATAR"
	RequestDeleteAvatar    constants.RequestAction = "DELETE_AVATAR"
	RequestEnableAvatar    constants.RequestAction = "ENABLE_AVATAR"
	RequestDisableAvatar   constants.RequestAction = "DISABLE_AVATAR"
	RequestUpdateAvatar    constants.RequestAction = "UPDATE_AVATAR"
	RequestBlockRegion     constants.RequestAction = "BLOCK_REGION"
	// RequestEnableRegion          constants.RequestAction = "ENABLE_REGION"
	RequestBlockDistrict constants.RequestAction = "BLOCK_DISTRICT"
	// RequestEnableDistrict        constants.RequestAction = "ENABLE_DISTRICT"
	// RequestEnableCity            constants.RequestAction = "ENABLE_CITY"
	RequestBlockCity constants.RequestAction = "BLOCK_CITY"

	// Branch
	RequestEnableBranches  constants.RequestAction = "REQUEST_ENABLE_BRANCHES"
	RequestDisableBranches constants.RequestAction = "REQUEST_DISABLE_BRANCHES"

	// Region
	RequestEnableRegions  constants.RequestAction = "REQUEST_ENABLE_REGIONS"
	RequestDisableRegions constants.RequestAction = "REQUEST_DISABLE_REGIONS"

	// District
	RequestEnableDistricts  constants.RequestAction = "REQUEST_ENABLE_DISTRICTS"
	RequestDisableDistricts constants.RequestAction = "REQUEST_DISABLE_DISTRICTS"

	// City
	RequestEnableCities  constants.RequestAction = "REQUEST_ENABLE_CITIES"
	RequestDisableCities constants.RequestAction = "REQUEST_DISABLE_CITIES"

	RequestCreateEventCategory   constants.RequestAction = "CREATE_EVENT_CATEGORY"
	RequestUpdateEventCategory   constants.RequestAction = "UPDATE_EVENT_CATEGORY"
	RequestBlockUser             constants.RequestAction = "BLOCK_USER"
	RequestEnableSingleBranches  constants.RequestAction = "REQUEST_ENABLE_SINGLE_BRANCHES"
	RequestDisableSingleBranches constants.RequestAction = "REQUEST_DISABLE_SINGLE_BRANCHES"
	RequestEnableMultiBranches   constants.RequestAction = "REQUEST_ENABLE_MULTI_BRANCHES"
	RequestDisableMultiBranches  constants.RequestAction = "REQUEST_DISABLE_MULTI_BRANCHES"
	RequestDisableFaydaAccount   constants.RequestAction = "DISABLE_FAYDA_ACCOUNT"
	RequestEnableFaydaAccount    constants.RequestAction = "ENABLE_FAYDA_ACCOUNT"

	RequestAccountUpdate         constants.RequestAction = "REQUEST_ACCOUNT_UPDATE"
	RequestCreatePermissionGroup constants.RequestAction = "CREATE_PERMISSION_GROUP"
	RequestUpdatePermissionGroup constants.RequestAction = "UPDATE_PERMISSION_GROUP"
	RequestDeletePermissionGroup constants.RequestAction = "DELETE_PERMISSION_GROUP"

	RequestCreateBudgetCategory  constants.RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestDeleteBudgetCategory  constants.RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestUpdateBudgetCategory  constants.RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestDisableBudgetCategory constants.RequestAction = "DISABLE_BUDGET_CATEGORY"
	RequestEnableBudgetCategory  constants.RequestAction = "ENABLE_BUDGET_CATEGORY"

	RequestUnlinkDevice           constants.RequestAction = "UNLINK_DEVICE"
	RequestUnlinkUser             constants.RequestAction = "UNLINK_USER"
	RequestUpdateHQBlockTime      constants.RequestAction = "UPDATE_HQ_BLOCK_TIME"
	RequestUpdateHQArchiveTime    constants.RequestAction = "UPDATE_HQ_ARCHIVE_TIME"
	RequestCreateMiniAppMerchant  constants.RequestAction = "CREATE_MINI_APP_MERCHANT"
	RequestUpdateMiniAppMerchant  constants.RequestAction = "UPDATE_MINI_APP_MERCHANT"
	RequestDeleteMiniAppMerchant  constants.RequestAction = "DELETE_MINI_APP_MERCHANT"
	RequestEnableMiniAppMerchant  constants.RequestAction = "ENABLE_MINI_APP_MERCHANT"
	RequestDisableMiniAppMerchant constants.RequestAction = "DISABLE_MINI_APP_MERCHANT"

	RequestCreateMiniApp  constants.RequestAction = "CREATE_MINI_APP"
	RequestUpdateMiniApp  constants.RequestAction = "UPDATE_MINI_APP"
	RequestDeleteMiniApp  constants.RequestAction = "DELETE_MINI_APP"
	RequestEnableMiniApp  constants.RequestAction = "ENABLE_MINI_APP"
	RequestDisableMiniApp constants.RequestAction = "DISABLE_MINI_APP"

	RequestCreateBudgetColor constants.RequestAction = "BUDGET_CREATE_COLOR"
	RequestUpdateBudgetColor constants.RequestAction = "BUDGET_UPDATE_COLOR"
	RequestDeleteBudgetColor constants.RequestAction = "BUDGET_DELETE_COLOR"
	RequestCreateBudgetIcon  constants.RequestAction = "BUDGET_CREATE_ICON"
	RequestUpdateBudgetIcon  constants.RequestAction = "BUDGET_UPDATE_ICON"
	RequestDeleteBudgetIcon  constants.RequestAction = "BUDGET_DELETE_ICON"

	RequestCreateServiceFee constants.RequestAction = "CREATE_SERVICE_FEE"
	RequestUpdateServiceFee constants.RequestAction = "UPDATE_SERVICE_FEE"
	RequestDeleteServiceFee constants.RequestAction = "DELETE_SERVICE_FEE"
	RequestCreateDailyLimit constants.RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit constants.RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit constants.RequestAction = "DELETE DAILY LIMIT"

	// RequestCreateNotification     constants.RequestAction = "CREATE_NOTIFICATION"
	// RequestUpdateNotification     constants.RequestAction = "UPDATE_NOTIFICATION"
	// RequestDeleteNotification     constants.RequestAction = "DELETE_NOTIFICATION"
	// RequestEnableNotification     constants.RequestAction = "ENABLE_NOTIFICATION"
	// RequestDisableNotification    constants.RequestAction = "DISABLE_NOTIFICATION"
	// RequestMarkNotificationAsSeen constants.RequestAction = "MARK_NOTIFICATION_AS_SEEN"
	RequestUpdateProductCode       constants.RequestAction = "UPDATE_PRODUCT_CODE"
	RequestCreateDonationCategory  constants.RequestAction = "CREATE_DONATION_CATEGORY"
	RequestUpdateDonationCategory  constants.RequestAction = "UPDATE_DONATION_CATEGORY"
	RequestEnableDonationCategory  constants.RequestAction = "ENABLE_DONATION_CATEGORY"
	RequestDeleteDonationCategory  constants.RequestAction = "DELETE_DONATION_CATEGORY"
	RequestDisableDonationCategory constants.RequestAction = "DISABLE_DONATION_CATEGORY"
	RequestEnableDonationCompany   constants.RequestAction = "ENABLE_DONATION_COMPANY"
	RequestDeleteDonationCompany   constants.RequestAction = "DELETE_DONATION_COMPANY"
	RequestDisableDonationCompany  constants.RequestAction = "DISABLE_DONATION_COMPANY"
	RequestCreateDonationCompany   constants.RequestAction = "CREATE_DONATION_COMPANY"
	RequestUpdateDonationCompany   constants.RequestAction = "UPDATE_DONATION_COMPANY"
	RequestCreateDonation          constants.RequestAction = "CREATE_DONATION"
	RequestUpdateDonation          constants.RequestAction = "UPDATE_DONATION"
	RequestUpdateDonationImage     constants.RequestAction = "UPDATE_DONATION_IMAGE"
	RequestDeleteDonationImage     constants.RequestAction = "DELETE_DONATION_IMAGE"
	RequestAddDonationImage        constants.RequestAction = "ADD_DONATION_IMAGE"
	RequestEnableDonation          constants.RequestAction = "ENABLE_DONATION"
	RequestDisableDonation         constants.RequestAction = "DISABLE_DONATION"
	RequestDeleteDonation          constants.RequestAction = "DELETE_DONATION"
	// for bankvault
	RequestCreateBankVault  constants.RequestAction = "CREATE_VAULT_BANK"
	RequestUpdateBankVault  constants.RequestAction = "UPDATE_VAULT_BANK"
	RequestDeleteBankVault  constants.RequestAction = "DELETE_VAULT_BANK"
	RequestEnableBankVault  constants.RequestAction = "ENABLE_VAULT_BANK"
	RequestDisAbleBankVault constants.RequestAction = "DISABLE_VAULT_BANK"

	// for vault group category
	RequestCreateVaultCategory  constants.RequestAction = "CREATE_VAULT_CATEGORY"
	RequestUpdateVaultCategory  constants.RequestAction = "UPDATE_VAULT_CATEGORY"
	RequestDeleteVaultCategory  constants.RequestAction = "DELETE_VAULT_CATEGORY"
	RequestEnableVaultCategory  constants.RequestAction = "ENABLE_VAULT_CATEGORY"
	RequestDisAbleVaultCategory constants.RequestAction = "DISABLE_VAULT_CATEGORY"
	RequestUnlockDeadlock       constants.RequestAction = "UNLOCK_DEADLOCK_REQUEST"

	RequestCreateVaultAmountTier  constants.RequestAction = "CREATE_VAULT_AMOUNT_TIER"
	RequestUpdateVaultAmountTier  constants.RequestAction = "UPDATE_VAULT_AMOUNT_TIER"
	RequestDeleteVaultAmountTier  constants.RequestAction = "DELETE_VAULT_AMOUNT_TIER"
	RequestEnableVaultAmountTier  constants.RequestAction = "ENABLE_VAULT_AMOUNT_TIER"
	RequestDisAbleVaultAmountTier constants.RequestAction = "DISABLE_VAULT_AMOUNT_TIER"

	RequestUpdateKYCVerifier constants.RequestAction = "UPDATE_KYC"
	RequestApproveKYC        constants.RequestAction = "APPROVE_KYC"
	// for article
	RequestCreateArticle  constants.RequestAction = "CREATE_ARTICLE"
	RequestUpdateArticle  constants.RequestAction = "UPDATE_ARTICLE"
	RequestEnableArticle  constants.RequestAction = "ENABLE_ARTICLE"
	RequestDisableArticle constants.RequestAction = "DISABLE_ARTICLE"
	RequestDeleteArticle  constants.RequestAction = "DELETE_ARTICLE"

	// for article category
	RequestCreateArticleCategory  constants.RequestAction = "CREATE_ARTICLE_CATEGORY"
	RequestUpdateArticleCategory  constants.RequestAction = "UPDATE_ARTICLE_CATEGORY"
	RequestDeleteArticleCategory  constants.RequestAction = "DELETE_ARTICLE_CATEGORY"
	RequestEnableArticleCategory  constants.RequestAction = "ENABLE_ARTICLE_CATEGORY"
	RequestDisableArticleCategory constants.RequestAction = "DISABLE_ARTICLE_CATEGORY"

	// for tags
	RequestCreateNewsTag  constants.RequestAction = "CREATE_NEWS_TAG"
	RequestUpdateNewsTag  constants.RequestAction = "UPDATE_NEWS_TAG"
	RequestEnableNewsTag  constants.RequestAction = "ENABLE_NEWS_TAG"
	RequestDisableNewsTag constants.RequestAction = "DISABLE_NEWS_TAG"
	RequestDeleteNewsTag  constants.RequestAction = "DELETE_NEWS_TAG"

	// Action Role Mapper
	RequestCreateActionRole  constants.RequestAction = "CREATE_ACTION_ROLE"
	RequestUpdateActionRole  constants.RequestAction = "UPDATE_ACTION_ROLE"
	RequestEnableActionRole  constants.RequestAction = "ENABLE_ACTION_ROLE"
	RequestDisableActionRole constants.RequestAction = "DISABLE_ACTION_ROLE"
	RequestDeleteActionRole  constants.RequestAction = "DELETE_ACTION_ROLE"

	RequestCreateCpsActionRole  constants.RequestAction = "CREATE_CPS_ACTION_ROLE"
	RequestUpdateCpsActionRole  constants.RequestAction = "UPDATE_CPS_ACTION_ROLE"
	RequestEnableCpsActionRole  constants.RequestAction = "ENABLE_CPS_ACTION_ROLE"
	RequestDisableCpsActionRole constants.RequestAction = "DISABLE_CPS_ACTION_ROLE"
	RequestDeleteCpsActionRole  constants.RequestAction = "DELETE_CPS_ACTION_ROLE"

	RequestCreateShortVideo  constants.RequestAction = "CREATE_SHORT_VIDEO"
	RequestUpdateShortVideo  constants.RequestAction = "UPDATE_SHORT_VIDEO"
	RequestEnableShortVideo  constants.RequestAction = "ENABLE_SHORT_VIDEO"
	RequestDisableShortVideo constants.RequestAction = "DISABLE_SHORT_VIDEO"
	RequestDeleteShortVideo  constants.RequestAction = "DELETE_SHORT_VIDEO"

	RequestApproveFaydaCustomer constants.RequestAction = "APPROVE_FAYDA_CUSTOMER"

	RequestEnableDisableCustomer constants.RequestAction = "ENABLE_DISABLE_CUSTOMER"

	RequestCreateMiniAppCategory  constants.RequestAction = "CREATE_MINI_APP_CATEGORY"
	RequestUpdateMiniAppCategory  constants.RequestAction = "UPDATE_MINI_APP_CATEGORY"
	RequestDeleteMiniAppCategory  constants.RequestAction = "DELETE_MINI_APP_CATEGORY"
	RequestEnableMiniAppCategory  constants.RequestAction = "ENABLE_MINI_APP_CATEGORY"
	RequestDisableMiniAppCategory constants.RequestAction = "DISABLE_MINI_APP_CATEGORY"

	RequestCreateEventMerchant  constants.RequestAction = "CREATE_EVENT_MERCHANT"
	RequestUpdateEventMerchant  constants.RequestAction = "UPDATE_EVENT_MERCHANT"
	RequestDeleteEventMerchant  constants.RequestAction = "DELETE_EVENT_MERCHANT"
	RequestEnableEventMerchant  constants.RequestAction = "ENABLE_EVENT_MERCHANT"
	RequestDisableEventMerchant constants.RequestAction = "DISABLE_EVENT_MERCHANT"

	RequestCreateLogisticsMerchant  constants.RequestAction = "CREATE_LOGISTICS_MERCHANT"
	RequestUpdateLogisticsMerchant  constants.RequestAction = "UPDATE_LOGISTICS_MERCHANT"
	RequestDeleteLogisticsMerchant  constants.RequestAction = "DELETE_LOGISTICS_MERCHANT"
	RequestEnableLogisticsMerchant  constants.RequestAction = "ENABLE_LOGISTICS_MERCHANT"
	RequestDisableLogisticsMerchant constants.RequestAction = "DISABLE_LOGISTICS_MERCHANT"

	RequestCreateAccessListSegmentation        constants.RequestAction = "CREATE_ACCESS_LIST_SEGMENTATION_BY_GEOGRAPHIC_LOCATION"
	RequestUpdateAccessListSegmentation        constants.RequestAction = "UPDATE_ACCESS_LIST_SEGMENTATION_BY_GEOGRAPHIC_LOCATION"
	RequestEnableDisableAccessListSegmentation constants.RequestAction = "ENABLE_DISABLE_ACCESS_LIST_SEGMENTATION_BY_GEOGRAPHIC_LOCATION"
	RequestEnableAccessListSegmentation        constants.RequestAction = "ENABLE_ACCESS_LIST_SEGMENTATION_BY_GEOGRAPHIC_LOCATION"
	RequestDisableAccessListSegmentation       constants.RequestAction = "DISABLE_ACCESS_LIST_SEGMENTATION_BY_GEOGRAPHIC_LOCATION"
	RequestDeleteServiceKey                    constants.RequestAction = "DELETE_SERVICE_KEY"
)

var validRequestActions = map[constants.RequestAction]struct{}{
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

	RequestCreateCpsRole:  {},
	RequestUpdateCpsRole:  {},
	RequestDeleteCpsRole:  {},
	RequestEnableCpsRole:  {},
	RequestDisableCpsRole: {},

	// for vault group category
	RequestCreateVaultCategory:  {},
	RequestUpdateVaultCategory:  {},
	RequestDeleteVaultCategory:  {},
	RequestEnableVaultCategory:  {},
	RequestDisAbleVaultCategory: {},
	RequestUnlockDeadlock:       {},

	RequestBpsUserCreate:  {},
	RequestBpsUserUpdate:  {},
	RequestBpsUserDelete:  {},
	RequestBpsUserEnable:  {},
	RequestBpsUserDisable: {},

	RequestCreateEcommerceMerchant:        {},
	RequestUpdateEcommerceMerchant:        {},
	RequestEnableEcommerceMerchant:        {},
	RequestDisableEcommerceMerchant:       {},
	RequestDeleteEcommerceMerchant:        {},
	RequestDeleteEcommerceMerchantBranch:  {},
	RequestEnableEcommerceMerchantBranch:  {},
	RequestDisableEcommerceMerchantBranch: {},

	RequestCreateDonationCategory: {},
	RequestUpdateDonationCategory: {},
	RequestCreateDonationCompany:  {},

	RequestDisableDonationCategory: {},
	RequestEnableDonationCompany:   {},
	RequestDisableDonationCompany:  {},
	RequestEnableDonationCategory:  {},
	RequestDeleteDonationCategory:  {},
	RequestDeleteDonationCompany:   {},

	RequestUpdateDonationCompany: {},
	RequestCreateDonation:        {},
	RequestUpdateDonation:        {},
	RequestUpdateDonationImage:   {},
	RequestDeleteDonationImage:   {},
	RequestAddDonationImage:      {},
	RequestEnableDonation:        {},
	RequestDeleteDonation:        {},
	RequestDisableDonation:       {},
	RequestAccountUpdate:         {},
	RequestDeleteAmountBasedAuth: {},
	RequestCreateAmountBasedAuth: {},
	RequestUpdateAmountBasedAuth: {},
	RequestResetAmountBasedAuth:  {},
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
	// RequestDisableBPSUser:   {},
	// RequestEnableBPSUser:    {},
	RequestUpdateUser:           {},
	RequestTotalDailyLimit:      {},
	RequestUpdateVAT:            {},
	RequestAuthTier:             {},
	RequestCreateAdvert:         {},
	RequestUpdateAdvert:         {},
	RequestEnableAdvert:         {},
	RequestDisableAdvert:        {},
	RequestDeleteAdvert:         {},
	RequestCreateBank:           {},
	RequestUpdateBank:           {},
	RequestUpdateBankLogo:       {},
	RequestEnableWallet:         {},
	RequestDisableWallet:        {},
	RequestEnableWalletService:  {},
	RequestDisableWalletService: {},

	// Services catalog
	RequestCreateService:      {},
	RequestUpdateService:      {},
	RequestEnableService:      {},
	RequestDisableService:     {},
	RequestDeleteService:      {},
	RequestCreateServiceList:  {},
	RequestUpdateServiceList:  {},
	RequestEnableServiceList:  {},
	RequestDisableServiceList: {},
	RequestDeleteServiceList:  {},
	RequestDeleteServiceKey:   {},

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

	RequestCreateCustomerGroup:  {},
	RequestUpdateCustomerGroup:  {},
	RequestEnableCustomerGroup:  {},
	RequestDisableCustomerGroup: {},
	RequestDeleteCustomerGroup:  {},

	RequestEnableSuperAppRole:  {},
	RequestDisableSuperAppRole: {},
	RequestDeleteSuperAppRole:  {},

	RequestCreateCustomerSegmentation:    {},
	RequestUpdateCustomerSegmentation:    {},
	RequestDisableAccessListSegmentation: {},
	RequestEnableCustomerSegmentation:    {},
	RequestDisableCustomerSegmentation:   {},
	RequestDeleteCustomerSegmentation:    {},

	RequestAccessListCreateCustomerSegmentation:  {},
	RequestAccessListUpdateCustomerSegmentation:  {},
	RequestAccessListEnableCustomerSegmentation:  {},
	RequestAccessListDisableCustomerSegmentation: {},
	RequestAccessListDeleteCustomerSegmentation:  {},

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
	RequestDeleteActionRole:     {},
	RequestCreateCpsActionRole:  {},
	RequestUpdateCpsActionRole:  {},
	RequestEnableCpsActionRole:  {},
	RequestDisableCpsActionRole: {},
	RequestDeleteCpsActionRole:  {},

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

	RequestCreateLogisticsMerchant:  {},
	RequestUpdateLogisticsMerchant:  {},
	RequestDeleteLogisticsMerchant:  {},
	RequestEnableLogisticsMerchant:  {},
	RequestDisableLogisticsMerchant: {},

	RequestCreateAccessListSegmentation:        {},
	RequestUpdateAccessListSegmentation:        {},
	RequestEnableDisableAccessListSegmentation: {},

	RequestApproveCustomerKYC: {},
	RequestRejectCustomerKYC:  {},
	// RequestCreateCustomerKYC: {},
	// RequestUpdateCustomerKYC: {},
	// RequestDeleteCustomerKYC: {},

	RequestEnableDisableBank: {},
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[constants.RequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]constants.RequestAction{
	// Canonical
	"NOTIFICATION": {
		RequestCreatePublicNotification,
		RequestUpdatePublicNotification,
		RequestDeleteNotification,
		RequestEnableNotification,
		RequestDisableNotification,
		RequestMarkNotificationAsSeen,
	},
	"ACCOUNTBLOCK": {
		RequestBlockUser,
		// RequestDisableSingleBranch,
		// RequestEnableSingleBranch,
		// RequestDisableMultiBranches,
		// RequestEnableMultiBranches,
		// RequestEnableBranches,
		// RequestDisableBranches,
		// RequestEnableRegions,
		// RequestDisableRegions,
		// RequestEnableDistricts,
		// RequestDisableDistricts,
		// RequestEnableCities,
		// RequestDisableCities,
	},
	"SINGLEBRANCHENABLEACCOUNTBLOCK": {
		RequestEnableSingleBranch,
		RequestEnableBranches,
	},
	"SINGLEBRANCHDISABLEACCOUNTBLOCK": {
		RequestDisableSingleBranch,
		RequestDisableBranches,
	},
	"MULTIBRANCHENABLEACCOUNTBLOCK": {
		RequestEnableRegions,
		RequestEnableDistricts,
		RequestEnableCities,
	},
	"MULTIBRANCHDISABLEACCOUNTBLOCK": {
		RequestDisableMultiBranches,
		RequestDisableRegions,
		RequestDisableDistricts,
		RequestDisableCities,
	},
	"ACCOUNTVALIDATION": {
		RequestUser,
		RequestUpdateAccountValidation,
		RequestEnableUser,
		RequestDisableUser,
		RequestUpdateUser,
		RequestArchiveUser,
	},
	"BPSACTIONROLE": {
		RequestCreateActionRole,
		RequestUpdateActionRole,
		RequestEnableActionRole,
		RequestDisableActionRole,
		RequestDeleteActionRole,
	},
	"JOBROLE": {
		RequestCreateJobRole,
		RequestUpdateJobRole,
		RequestDeleteJobRole,
		RequestEnableJobRole,
		RequestDisableJobRole,
	},
	"PASSWORDRULE": {
		RequestUpdatePasswordRule,
	},
	"VAULTCATEGORIES": {
		RequestCreateVaultCategory,
		RequestUpdateVaultCategory,
		RequestDeleteVaultCategory,
		RequestEnableVaultCategory,
		RequestDisAbleVaultCategory,
	},
	"VAULTEMERGENCYDEADLOCKREQUEST": {
		RequestUnlockDeadlock,
	},

	// Legacy/operational modules (kept as requested)
	"SERVICE": {
		RequestCreateService,
		RequestUpdateService,
		RequestEnableService,
		RequestDisableService,
		RequestDeleteService,
		RequestCreateServiceList,
		RequestUpdateServiceList,
		RequestEnableServiceList,
		RequestDisableServiceList,
		RequestDeleteServiceList,
		RequestDeleteServiceKey,
	},
	"USSDMERCHANT": {
		RequestCreateUssdMerchant,
		RequestUpdateUssdMerchant,
		RequestDeleteUssdMerchant,
		RequestEnableUssdMerchant,
		RequestDisableUssdMerchant,
	},
	"EVENT": {
		RequestCreateEvent,
		RequestDeleteEvent,
		RequestDisableEvent,
		RequestEnableEvent,
		RequestUpdateEvent,
	},
	"AMOUNTBASEDAUTH": {
		RequestCreateAmountBasedAuth,
		RequestUpdateAmountBasedAuth,
		RequestDeleteAmountBasedAuth,
		RequestResetAmountBasedAuth,
		RequestAuthTier,
	},
	"USER": {
		RequestUser,
		RequestEnableUser,
		RequestDisableUser,
		RequestUpdateUser,
		RequestArchiveUser,
	},
	"BPSUSER": {
		RequestBPSUser,
		RequestEnableBPSUser,
		RequestDisableBPSUser,
		RequestBpsUserCreate,
		RequestBpsUserUpdate,
		RequestBpsUserDelete,
		RequestBpsUserEnable,
		RequestBpsUserDisable,
	},
	"PERMISSIONGROUP": {
		RequestPermissionGroup,
	},
	"DEPARTMENT": {
		RequestCreateDepartment,
		RequestUpdateDepartment,
		RequestEnableDisableDepartment,
		RequestDeleteDepartment,
	},
	"SERVICEFEE": {
		RequestCreateServiceFee,
		RequestUpdateServiceFee,
		RequestDeleteServiceFee,
	},
	"DAILYLIMIT": {
		RequestCreateDailyLimit,
		RequestUpdateDailyLimit,
		RequestDeleteDailyLimit,
		RequestTotalDailyLimit,
	},
	"VAT": {RequestUpdateVAT},
	// "AUTHTIER":       {RequestAuthTier},
	"ARCHIVE":        {RequestUpdateArchiveExpiry},
	"MINIMUMSERVICE": {RequestUpdateMinimumService},
	"SERVICERULE":    {RequestUpdateServiceRule},
	"TOTAL":          {RequestUpdateTotal},
	"ACCESSCONFIG":   {RequestUpdateAccessConfig},
	"BRANCH": {
		RequestEnableSingleBranch,
		RequestEnableMultiUsers,
		RequestDisableMultiUsers,
		RequestEnableSingleBranches,
		RequestEnableMultiBranches,
	},
	"BUSINESS": {
		RequestCreateBusiness,
		RequestUpdateBusiness,
	},
	"EVENTCATEGORY": {
		RequestCreateEventCategory,
		RequestUpdateEventCategory,
	},
	"MINIAPPMERCHANT": {
		RequestCreateMiniAppMerchant,
		RequestUpdateMiniAppMerchant,
		RequestDeleteMiniAppMerchant,
		RequestEnableMiniAppMerchant,
		RequestDisableMiniAppMerchant,
	},
	"BLOCKTIME": {RequestUpdateBlockTime},
	// "PASSWORDRULE": {RequestUpdatePasswordRule},
	"PERMISSION": {
		RequestCreatePermissionGroup,
		RequestDeletePermissionGroup,
		RequestUpdatePermissionGroup,
	},
	"AVATAR": {
		RequestCreateAvatar,
		RequestUpdateAvatar,
		RequestEnableAvatar,
		RequestDisableAvatar,
		RequestDeleteAvatar,
	},
	"BUDGET": {
		RequestCreateBudgetColor,
		RequestUpdateBudgetColor,
		RequestDeleteBudgetColor,
		RequestCreateBudgetIcon,
		RequestUpdateBudgetIcon,
		RequestDeleteBudgetIcon,
	},
	"ADVERT": {
		RequestCreateAdvert,
		RequestUpdateAdvert,
		RequestEnableAdvert,
		RequestDisableAdvert,
		RequestDeleteAdvert,
	},
	"BANK": {
		RequestCreateBank,
		RequestUpdateBank,
		RequestDeleteBank,
		RequestEnableDisableBank,
		RequestUpdateBankLogo,
		RequestEnableBank,
		RequestDisableBank,
	},
	"WALLET": {
		RequestCreateWallet,
		RequestUpdateWallet,
		RequestDeleteWallet,
		RequestEnableWallet,
		RequestDisableWallet,
		RequestEnableWalletService,
		RequestDisableWalletService,
	},
	"TOPUP": {
		RequestCreateTopup,
		RequestUpdateTopup,
		RequestDeleteTopup,
		RequestEnableTopup,
		RequestDisableTopup,
	},
	"VALIDATION": {
		RequestCreateValidation,
		RequestUpdateValidation,
		RequestDeleteValidation,
	},
	// "JOBROLE": {
	// 	RequestCreateJobRole,
	// 	RequestUpdateJobRole,
	// 	RequestDeleteJobRole,
	// 	RequestEnableJobRole,
	// 	RequestDisableJobRole,
	// },
	"ROLE": {
		RequestCreateRole,
		RequestUpdateRole,
		RequestDeleteRole,
		RequestEnableRole,
		RequestDisableRole,
	},
	// "BLOCK": {
	// 	RequestBlockUser,
	// 	RequestDisableSingleBranch,
	// 	RequestEnableSingleBranch,
	// 	RequestDisableMultiBranches,
	// 	RequestEnableMultiBranches,
	// 	RequestEnableBranches,
	// 	RequestDisableBranches,
	// 	RequestEnableRegions,
	// 	RequestDisableRegions,
	// 	RequestEnableDistricts,
	// 	RequestDisableDistricts,
	// 	RequestEnableCities,
	// 	RequestDisableCities,
	// },
	"BUDGETCATEGORY": {
		RequestCreateBudgetCategory,
		RequestUpdateBudgetCategory,
		RequestDeleteBudgetCategory,
		RequestDisableBudgetCategory,
		RequestEnableBudgetCategory,
	},
	"UNLINKDEVICE": {
		RequestUnlinkDevice,
		RequestUnlinkUser,
	},
	"HQ": {
		RequestUpdateHQBlockTime,
		RequestUpdateHQArchiveTime,
		RequestUpdatePasswordExpiry,
	},
	"FAYDA": {
		RequestDisableFaydaAccount,
		RequestEnableFaydaAccount,
	},
	"CPSUSER": {
		RequestCpsUserCreate,
		RequestCpsUserUpdate,
		RequestCpsUserDelete,
		RequestCpsUserEnable,
		RequestCpsUserDisable,
	},
	"BULKSERVICEALLUSER": {
		RequestBulkServiceEnable,
		RequestBulkServiceDisable,
	},
	"MINIAPP": {
		RequestCreateMiniApp,
		RequestUpdateMiniApp,
		RequestDeleteMiniApp,
		RequestEnableMiniApp,
		RequestDisableMiniApp,
	},
	// "NOTIFICATION": {
	// 	RequestCreatePublicNotification,
	// 	RequestUpdatePublicNotification,
	// 	RequestDeleteNotification,
	// 	RequestEnableNotification,
	// 	RequestDisableNotification,
	// 	RequestMarkNotificationAsSeen,
	// },
	"BANKVAULT": {
		RequestCreateBankVault,
		RequestUpdateBankVault,
		RequestDeleteBankVault,
		RequestEnableBankVault,
		RequestDisAbleBankVault,
	},
	"PRODUCTCODE": {
		RequestUpdateProductCode,
	},
	"DONATION": {
		RequestCreateDonation,
		RequestUpdateDonation,
		RequestDisableDonation,
		RequestAddDonationImage,
		RequestUpdateDonationImage,
		RequestDeleteDonationImage,
		RequestEnableDonation,
		RequestDeleteDonation,
	},
	"DONATIONCATEGORY": {
		RequestCreateDonationCategory,
		RequestUpdateDonationCategory,
		RequestDisableDonationCategory,
		RequestEnableDonationCategory,
		RequestDeleteDonationCategory,
	},
	"DONATIONCOMPANY": {
		RequestCreateDonationCompany,
		RequestUpdateDonationCompany,
		RequestEnableDonationCompany,
		RequestDisableDonationCompany,
		RequestDeleteDonationCompany,
	},
	"KYCVERIFIER": {
		RequestUpdateKYCVerifier,
		RequestApproveKYC,
	},
	"ARTICLE": {
		RequestCreateArticle,
		RequestUpdateArticle,
		RequestDeleteArticle,
		RequestEnableArticle,
		RequestDisableArticle,
	},
	"ARTICLECATEGORY": {
		RequestCreateArticleCategory,
		RequestUpdateArticleCategory,
		RequestDeleteArticleCategory,
		RequestEnableArticleCategory,
		RequestDisableArticleCategory,
	},
	"SHORTVIDEO": {
		RequestCreateShortVideo,
		RequestUpdateShortVideo,
		RequestDeleteShortVideo,
		RequestEnableShortVideo,
		RequestDisableShortVideo,
	},
	"CUSTOMER": {
		RequestEnableDisableCustomer,
		RequestApproveFaydaCustomer,
	},
	"NEWSCATEGORY": {
		constants.RequestAction("CREATE_NEWS_CATEGORY"),
		constants.RequestAction("UPDATE_NEWS_CATEGORY"),
		constants.RequestAction("DELETE_NEWS_CATEGORY"),
	},
	"NEWSTAG": {
		RequestCreateNewsTag,
		RequestUpdateNewsTag,
		RequestEnableNewsTag,
		RequestDisableNewsTag,
		RequestDeleteNewsTag,
	},
	// "ACTIONROLE": {
	// 	RequestCreateActionRole,
	// 	RequestUpdateActionRole,
	// 	RequestEnableActionRole,
	// 	RequestDisableActionRole,
	// },
	"CPSACTIONROLE": {
		RequestCreateCpsActionRole,
		RequestUpdateCpsActionRole,
		RequestEnableCpsActionRole,
		RequestDisableCpsActionRole,
		RequestDeleteCpsActionRole,
	},
	"DEVICEVERSION": {
		RequestCreateDeviceVersion,
		RequestUpdateDeviceVersion,
		RequestEnableDeviceVersion,
		RequestDisableDeviceVersion,
		RequestDeleteDeviceVersion,
		RequestEnableDisableDeviceVersion,
	},
	"MINIAPPCATEGORY": {
		RequestCreateMiniAppCategory,
		RequestUpdateMiniAppCategory,
		RequestDeleteMiniAppCategory,
		RequestEnableMiniAppCategory,
		RequestDisableMiniAppCategory,
	},
	"EVENTMERCHANT": {
		RequestCreateEventMerchant,
		RequestUpdateEventMerchant,
		RequestDeleteEventMerchant,
		RequestEnableEventMerchant,
		RequestDisableEventMerchant,
	},
	"LOGISTICSMERCHANT": {
		RequestCreateLogisticsMerchant,
		RequestUpdateLogisticsMerchant,
		RequestDeleteLogisticsMerchant,
		RequestEnableLogisticsMerchant,
		RequestDisableLogisticsMerchant,
	},
	"VAULTAMOUNTTIER": {
		RequestCreateVaultAmountTier,
		RequestUpdateVaultAmountTier,
		RequestDeleteVaultAmountTier,
		RequestEnableVaultAmountTier,
		RequestDisAbleVaultAmountTier,
	},
	"MINIAPPPRODUCTCODE": {
		RequestCreateMiniappProductCode,
		RequestUpdateMiniappProductCode,
		RequestDeleteMiniappProductCode,
		RequestEnableMiniappProductCode,
		RequestDisableMiniappProductCode,
	},
	"ACCESSLISTSEGMENTATION": {
		RequestCreateAccessListSegmentation,
		RequestUpdateAccessListSegmentation,
		RequestEnableDisableAccessListSegmentation,
		RequestEnableAccessListSegmentation,
		RequestDisableAccessListSegmentation,
	},
	// this is made as alternative to "CUSTOMERGROUP"
	"CUSTOMERSEGMENTATIONS": {
		RequestCreateCustomerGroup,
		RequestUpdateCustomerGroup,
		RequestEnableCustomerGroup,
		RequestDisableCustomerGroup,
		RequestDeleteCustomerGroup,
	},
	// this is made as alternative to "SUPPERROLE"
	"CPSROLE": {
		RequestEnableSuperAppRole,
		RequestDisableSuperAppRole,
		RequestDeleteSuperAppRole,
	},
	"BULKSERVICECUSTOMERSEGMENT": {
		// "CUSTOMERSEGMENTATIONS": {
		RequestCreateCustomerSegmentation,
		RequestUpdateCustomerSegmentation,
		RequestEnableCustomerSegmentation,
		RequestDisableCustomerSegmentation,
		RequestDeleteCustomerSegmentation,
	},
	"ACCESSLISTCUSTOMERSEGMENT": {
		// "CUSTOMERSEGMENTATIONS": {RequestEnableAccessListSegmentation
		// RequestAccessListCreateCustomerSegmentation,
		RequestAccessListUpdateCustomerSegmentation,
		RequestAccessListEnableCustomerSegmentation,
		RequestAccessListDisableCustomerSegmentation,
		RequestAccessListDeleteCustomerSegmentation,
	},
	"ECOMMERCEMERCHANT": {
		RequestCreateEcommerceMerchant,
		RequestUpdateEcommerceMerchant,
		RequestEnableEcommerceMerchant,
		RequestDisableEcommerceMerchant,
		RequestDeleteEcommerceMerchant,
		RequestDeleteEcommerceMerchantBranch,
		RequestEnableEcommerceMerchantBranch,
		RequestDisableEcommerceMerchantBranch,
	},
	// "CPSROLE": {
	// 	RequestCreateCpsRole,
	// 	RequestUpdateCpsRole,
	// 	RequestDeleteCpsRole,
	// 	RequestEnableCpsRole,
	// 	RequestDisableCpsRole,
	// },
	"CUSTOMERKYC": {
		RequestApproveCustomerKYC,
		RequestRejectCustomerKYC,
		// RequestCreateCustomerKYC,
		// RequestUpdateCustomerKYC,
		// RequestDeleteCustomerKYC,
	},
}

func IsActionInGroup(action constants.RequestAction, group string) bool {
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
