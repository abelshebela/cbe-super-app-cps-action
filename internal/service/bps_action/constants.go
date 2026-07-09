package bps_action

import "strings"

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
	case ActionCreate, ActionUpdate, ActionDelete, ActionEnable, ActionDisable:
		return true
	default:
		return false
	}
}

func NormalizeRequestAction(requestAction string) RequestAction {
	return RequestAction(strings.ToUpper(strings.TrimSpace(requestAction)))
}

type RequestAction string

const (
	LinkAccount                     RequestAction = "LINK_ACCOUNT"
	LinkOtherAccount                RequestAction = "LINK_OTHER_ACCOUNT"
	UnlinkAccount                   RequestAction = "UNLINK_ACCOUNT"
	ResetPin                        RequestAction = "RESET_PIN"
	ActivateAccount                 RequestAction = "ACTIVATE_ACCOUNT"
	ActivateSuperApp                RequestAction = "ACTIVATE_SUPERAPP"
	UnlinkDevice                    RequestAction = "UNLINK_DEVICE"
	ChangeName                      RequestAction = "CHANGE_NAME"
	LinkAndorAccount                RequestAction = "LINK_ANDOR_ACCOUNT"
	RejectAccpimt                   RequestAction = "REJECT_ACCOUNT"
	DetachPhoneNumber               RequestAction = "DETACH_PHONE_NUMBER"
	Reactivate                      RequestAction = "REACTIVATE"
	Terminate                       RequestAction = "TERMINATION"
	ChangePhoneNumber               RequestAction = "CHANGE_PHONE_NUMBER"
	AddAccount                      RequestAction = "ADD_ACCOUNT"
	AttachPhoneNumber               RequestAction = "ATTACH_PHONE_NUMBER"
	EnableBlocked                   RequestAction = "ENABLE_BLOCKED"
	DisableBlocked                  RequestAction = "DISABLE_BLOCKED"
	DisableUSSD                     RequestAction = "DISABLE_USSD"
	DisableInAPP                    RequestAction = "DISABLE_INAPP"
	DisableBoth                     RequestAction = "DISABLE_BOTH"
	UpdateOneLimit                  RequestAction = "UPDATE_ONE_LIMIT"
	UpdateLimit                     RequestAction = "UPDATE_LIMIT"
	EnableUSSD                      RequestAction = "ENABLE_USSD"
	Renewal                         RequestAction = "RENEWAL"
	EnableInAPP                     RequestAction = "ENABLE_INAPP"
	EnableBoth                      RequestAction = "ENABLE_BOTH"
	ChangeEmail                     RequestAction = "CHANGE_EMAIL"
	UpgradeAccount                  RequestAction = "UPGRADE_ACCOUNT"
	AccessControl                   RequestAction = "ACCESS_CONTROL"
	CreateMerchant                  RequestAction = "CREATE_MERCHANT"
	UpdateMerchant                  RequestAction = "UPDATE_MERCHANT"
	EnableDisableUser               RequestAction = "ENABLE_DISABLE_USER"
	LimitTransfer                   RequestAction = "LIMIT_TRANSFER"
	UpdateTransferLimits            RequestAction = "UPDATE_TRANSFER_LIMIT"
	CreateTransferLimits            RequestAction = "CREATE_TRANSFER_LIMIT"
	ResetTransferLimits             RequestAction = "RESET_TRANSFER_LIMIT"
	Reversal                        RequestAction = "REVERSAL"
	Renew                           RequestAction = "RENEW"
	Member                          RequestAction = "MEMBER"
	Blocked                         RequestAction = "BLOCKED"
	BpsAction                       RequestAction = "BPS_ACTION"
	CreateChequeAuthorizationAction RequestAction = "CREATE_CHEQUE_AUTHORIZATION_ACTION"
	ActivateResetPin                RequestAction = "ACTIVATE_RESET_PIN"
	UnlockPin                       RequestAction = "UNLOCK_PIN"
	ActivateDeactivated             RequestAction = "ACTIVATE_DEACTIVATED"
	ActivatePinUnlock               RequestAction = "ACTIVATE_PIN_UNLOCK"
	EnableDisableUSSD               RequestAction = "ENABLE_DISABLE_USSD"
	EnableDisableInApp              RequestAction = "ENABLE_DISABLE_INAPP"
	EnableDisableBoth               RequestAction = "ENABLE_DISABLE_BOTH"
	ActivateTerminated              RequestAction = "ACTIVATE_TERMINATED"
	ActivateChangePhoneNumber       RequestAction = "ACTIVATE_CHANGE_PHONE_NUMBER"

	//==================================================
	User                     RequestAction = "user"
	PermissionGroup          RequestAction = "permission_group"
	Department               RequestAction = "deparmtent"
	stEnableUser             RequestAction = "ENABLE_USER"
	DisableUser              RequestAction = "DISABLE_USER"
	BPSUser                  RequestAction = "bps_user"
	DisableBPSUser           RequestAction = "disable_bps_user"
	EnableBPSUser            RequestAction = "enable_bps_user"
	UpdateUser               RequestAction = "update_user"
	TotalDailyLimit          RequestAction = "total_daily_limit"
	UpdateVAT                RequestAction = "update_vat"
	AuthTier                 RequestAction = "authtier"
	CreateAdvert             RequestAction = "create_advert"
	UpdateAdvert             RequestAction = "update_advert"
	EnableAdvert             RequestAction = "enable_advert"
	DisableAdvert            RequestAction = "disable_advert"
	DeleteAdvert             RequestAction = "delete_advert"
	CreateBank               RequestAction = "create_bank"
	UpdateBank               RequestAction = "update_bank"
	EnableBank               RequestAction = "enable_bank"
	DisableBank              RequestAction = "disable_bank"
	EnableWallet             RequestAction = "enable_wallet"
	DisableWallet            RequestAction = "disable_wallet"
	UpdatePasswordExpiry     RequestAction = "update_password_expiry"
	CreateValidation         RequestAction = "create_validation"
	UpdateValidation         RequestAction = "update_validation"
	DeleteValidation         RequestAction = "delete_validation"
	UpdateArchiveExpiry      RequestAction = "update_archive_expiry"
	CreateServiceFee         RequestAction = "create_service_fee"
	UpdateServiceFee         RequestAction = "update_service_fee"
	DeleteServiceFee         RequestAction = "delete_service_fee"
	CreateDailyLimit         RequestAction = "create_daily_limit"
	UpdateDailyLimit         RequestAction = "update_daily_limit"
	DeleteDailyLimit         RequestAction = "delete_daily_limit"
	BudgetColor              RequestAction = "budget_color"
	BudgetIcon               RequestAction = "budget_icon"
	UpdateProduct            RequestAction = "update_product"
	CreatePublicNotification RequestAction = "create_public_notification"
	ArchiveUser              RequestAction = "archive_user"
	CreatePasswordRule       RequestAction = "create_password_rule"
	UpdatePasswordRule       RequestAction = "update_password_rule"
	UpdateMinimumService     RequestAction = "update_minimum_service"
	UpdateServiceRule        RequestAction = "update_service_rule"
	UpdateTotal              RequestAction = "update_total"
	UpdateAccessConfig       RequestAction = "update_access_config"
	EnableSingleBranch       RequestAction = "enable_single_branch"
	DisableSingleBranch      RequestAction = "disable_single_branch"
	EnableMultiUsers         RequestAction = "enable_multi_users"
	DisableMultiUsers        RequestAction = "disable_multi_users"
	CreateBusiness           RequestAction = "create_business"
	UpdateBusiness           RequestAction = "update_business"
	CreateEvent              RequestAction = "create_event"
	UpdateEvent              RequestAction = "update_event"
	CreateEventCategory      RequestAction = "create_event_category"
	UpdateEventCategory      RequestAction = "update_event_category"
	DisableEvent             RequestAction = "disable_event"
	CreateMiniAppMerchant    RequestAction = "create_miniapp_merchant"
	UpdateMiniAppMerchant    RequestAction = "update_miniapp_merchant"
	UpdateBlockTime          RequestAction = "update_block_time"
	UpdateAccountValidation  RequestAction = "update_account_validation"
	UpdateServiceDetails     RequestAction = "update_service_details"
	UpdateHQBlockTime        RequestAction = "update_hq_block_time"
	UpdateHQArchiveTime      RequestAction = "update_hq_archive_time"
	ResetAccessControl       RequestAction = "RESET_ACCESS_CONTROL"
	UpgradeKYCLevel          RequestAction = "UPGRADE_KYC_LEVEL"
	EnableUssdSupperapp      RequestAction = "ENABLE_USSD_SUPERAPP"
	DisableUssdSupperapp     RequestAction = "DISABLE_USSD_SUPERAPP"
)

var validRequestActions = map[RequestAction]struct{}{
	LinkAccount:                     {},
	LinkOtherAccount:                {},
	UnlinkAccount:                   {},
	ResetPin:                        {},
	ActivateAccount:                 {},
	Terminate:                       {},
	ActivateSuperApp:                {},
	UnlinkDevice:                    {},
	ChangeName:                      {},
	LinkAndorAccount:                {},
	RejectAccpimt:                   {},
	DetachPhoneNumber:               {},
	Reactivate:                      {},
	ChangePhoneNumber:               {},
	AddAccount:                      {},
	AttachPhoneNumber:               {},
	EnableBlocked:                   {},
	DisableBlocked:                  {},
	DisableUSSD:                     {},
	DisableInAPP:                    {},
	DisableBoth:                     {},
	UpdateOneLimit:                  {},
	UpdateLimit:                     {},
	EnableUSSD:                      {},
	Renewal:                         {},
	EnableInAPP:                     {},
	EnableBoth:                      {},
	ChangeEmail:                     {},
	UpgradeAccount:                  {},
	AccessControl:                   {},
	CreateMerchant:                  {},
	UpdateMerchant:                  {},
	EnableDisableUser:               {},
	LimitTransfer:                   {},
	UpdateTransferLimits:            {},
	CreateTransferLimits:            {},
	ResetTransferLimits:             {},
	Reversal:                        {},
	Renew:                           {},
	Member:                          {},
	Blocked:                         {},
	BpsAction:                       {},
	CreateChequeAuthorizationAction: {},
	User:                            {},
	PermissionGroup:                 {},
	Department:                      {},
	stEnableUser:                    {},
	DisableUser:                     {},
	BPSUser:                         {},
	DisableBPSUser:                  {},
	EnableBPSUser:                   {},
	UpdateUser:                      {},
	TotalDailyLimit:                 {},
	UpdateVAT:                       {},
	AuthTier:                        {},
	CreateAdvert:                    {},
	UpdateAdvert:                    {},
	EnableAdvert:                    {},
	DisableAdvert:                   {},
	DeleteAdvert:                    {},
	CreateBank:                      {},
	UpdateBank:                      {},
	EnableBank:                      {},
	DisableBank:                     {},
	EnableWallet:                    {},
	DisableWallet:                   {},
	UpdatePasswordExpiry:            {},
	CreateValidation:                {},
	UpdateValidation:                {},
	DeleteValidation:                {},
	UpdateArchiveExpiry:             {},
	CreateServiceFee:                {},
	UpdateServiceFee:                {},
	DeleteServiceFee:                {},
	CreateDailyLimit:                {},
	UpdateDailyLimit:                {},
	DeleteDailyLimit:                {},
	BudgetColor:                     {},
	BudgetIcon:                      {},
	UpdateProduct:                   {},
	CreatePublicNotification:        {},
	ArchiveUser:                     {},
	CreatePasswordRule:              {},
	UpdatePasswordRule:              {},
	UpdateMinimumService:            {},
	UpdateServiceRule:               {},
	UpdateTotal:                     {},
	UpdateAccessConfig:              {},
	EnableSingleBranch:              {},
	DisableSingleBranch:             {},
	EnableMultiUsers:                {},
	DisableMultiUsers:               {},
	CreateBusiness:                  {},
	UpdateBusiness:                  {},
	CreateEvent:                     {},
	UpdateEvent:                     {},
	CreateEventCategory:             {},
	UpdateEventCategory:             {},
	DisableEvent:                    {},
	CreateMiniAppMerchant:           {},
	UpdateMiniAppMerchant:           {},
	UpdateBlockTime:                 {},
	UpdateAccountValidation:         {},
	UpdateServiceDetails:            {},
	UpdateHQBlockTime:               {},
	UpdateHQArchiveTime:             {},
	ResetAccessControl:              {},
	UpgradeKYCLevel:                 {},

	ActivateResetPin:          {},
	UnlockPin:                 {},
	ActivateDeactivated:       {},
	ActivatePinUnlock:         {},
	EnableDisableUSSD:         {},
	EnableDisableInApp:        {},
	EnableDisableBoth:         {},
	ActivateTerminated:        {},
	ActivateChangePhoneNumber: {},
	DisableUssdSupperapp:      {},
	EnableUssdSupperapp:       {},
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[NormalizeRequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]RequestAction{
	"LINK_ACCOUNT":                       {LinkAccount},
	"TERMINATION":                        {Terminate},
	"LINK_OTHER_ACCOUNT":                 {LinkOtherAccount},
	"UNLINK_ACCOUNT":                     {UnlinkAccount},
	"RESET_PIN":                          {ResetPin},
	"ACTIVATE_ACCOUNT":                   {ActivateAccount},
	"ACTIVATE_SUPERAPP":                  {ActivateSuperApp},
	"UNLINK_DEVICE":                      {UnlinkDevice},
	"CHANGE_NAME":                        {ChangeName},
	"LINK_ANDOR_ACCOUNT":                 {LinkAndorAccount},
	"REJECT_ACCOUNT":                     {RejectAccpimt},
	"DETACH_PHONE_NUMBER":                {DetachPhoneNumber},
	"REACTIVATE":                         {Reactivate},
	"CHANGE_PHONE_NUMBER":                {ChangePhoneNumber},
	"ADD_ACCOUNT":                        {AddAccount},
	"ATTACH_PHONE_NUMBER":                {AttachPhoneNumber},
	"ENABLE_BLOCKED":                     {EnableBlocked},
	"DISABLE_BLOCKED":                    {DisableBlocked},
	"DISABLE_USSD":                       {DisableUSSD},
	"DISABLE_INAPP":                      {DisableInAPP},
	"DISABLE_BOTH":                       {DisableBoth},
	"UPDATE_ONE_LIMIT":                   {UpdateOneLimit},
	"UPDATE_LIMIT":                       {UpdateLimit},
	"ENABLE_USSD":                        {EnableUSSD},
	"RENEWAL":                            {Renewal},
	"ENABLE_INAPP":                       {EnableInAPP},
	"ENABLE_BOTH":                        {EnableBoth},
	"CHANGE_EMAIL":                       {ChangeEmail},
	"UPGRADE_ACCOUNT":                    {UpgradeAccount},
	"ACCESS_CONTROL":                     {AccessControl},
	"CREATE_MERCHANT":                    {CreateMerchant},
	"UPDATE_MERCHANT":                    {UpdateMerchant},
	"ENABLE_DISABLE_USER":                {EnableDisableUser},
	"LIMIT_TRANSFER":                     {LimitTransfer},
	"UPDATE_TRANSFER_LIMIT":              {UpdateTransferLimits},
	"CREATE_TRANSFER_LIMIT":              {CreateTransferLimits},
	"RESET_TRANSFER_LIMIT":               {ResetTransferLimits},
	"REVERSAL":                           {Reversal},
	"RENEW":                              {Renew},
	"MEMBER":                             {Member},
	"BLOCKED":                            {Blocked},
	"BPS_ACTION":                         {BpsAction},
	"CREATE_CHEQUE_AUTHORIZATION_ACTION": {CreateChequeAuthorizationAction},
	"user":                               {User},
	"permission_group":                   {PermissionGroup},
	"deparmtent":                         {Department},
	"ENABLE_USER":                        {stEnableUser},
	"DISABLE_USER":                       {DisableUser},
	"bps_user":                           {BPSUser},
	"disable_bps_user":                   {DisableBPSUser},
	"enable_bps_user":                    {EnableBPSUser},
	"update_user":                        {UpdateUser},
	"total_daily_limit":                  {TotalDailyLimit},
	"update_vat":                         {UpdateVAT},
	"authtier":                           {AuthTier},
	"create_advert":                      {CreateAdvert},
	"update_advert":                      {UpdateAdvert},
	"enable_advert":                      {EnableAdvert},
	"disable_advert":                     {DisableAdvert},
	"delete_advert":                      {DeleteAdvert},
	"create_bank":                        {CreateBank},
	"update_bank":                        {UpdateBank},
	"enable_bank":                        {EnableBank},
	"disable_bank":                       {DisableBank},
	"enable_wallet":                      {EnableWallet},
	"disable_wallet":                     {DisableWallet},
	"update_password_expiry":             {UpdatePasswordExpiry},
	"create_validation":                  {CreateValidation},
	"update_validation":                  {UpdateValidation},
	"delete_validation":                  {DeleteValidation},
	"update_archive_expiry":              {UpdateArchiveExpiry},
	"create_service_fee":                 {CreateServiceFee},
	"update_service_fee":                 {UpdateServiceFee},
	"delete_service_fee":                 {DeleteServiceFee},
	"create_daily_limit":                 {CreateDailyLimit},
	"update_daily_limit":                 {UpdateDailyLimit},
	"delete_daily_limit":                 {DeleteDailyLimit},
	"budget_color":                       {BudgetColor},
	"budget_icon":                        {BudgetIcon},
	"update_product":                     {UpdateProduct},
	"create_public_notification":         {CreatePublicNotification},
	"archive_user":                       {ArchiveUser},
	// OtherLinkAccount = "LINK_OTHER_ACCOUNT" — same value as LinkOtherAccount, already covered above
	"create_password_rule":      {CreatePasswordRule},
	"update_password_rule":      {UpdatePasswordRule},
	"update_minimum_service":    {UpdateMinimumService},
	"update_service_rule":       {UpdateServiceRule},
	"update_total":              {UpdateTotal},
	"update_access_config":      {UpdateAccessConfig},
	"enable_single_branch":      {EnableSingleBranch},
	"disable_single_branch":     {DisableSingleBranch},
	"enable_multi_users":        {EnableMultiUsers},
	"disable_multi_users":       {DisableMultiUsers},
	"create_business":           {CreateBusiness},
	"update_business":           {UpdateBusiness},
	"create_event":              {CreateEvent},
	"update_event":              {UpdateEvent},
	"create_event_category":     {CreateEventCategory},
	"update_event_category":     {UpdateEventCategory},
	"disable_event":             {DisableEvent},
	"create_miniapp_merchant":   {CreateMiniAppMerchant},
	"update_miniapp_merchant":   {UpdateMiniAppMerchant},
	"update_block_time":         {UpdateBlockTime},
	"update_account_validation": {UpdateAccountValidation},
	"update_service_details":    {UpdateServiceDetails},
	"update_hq_block_time":      {UpdateHQBlockTime},
	"update_hq_archive_time":    {UpdateHQArchiveTime},
	"RESET_ACCESS_CONTROL":      {ResetAccessControl},
	"UPGRADE_KYC_LEVEL":         {UpgradeKYCLevel},

	"ACTIVATE_RESET_PIN":           {ActivateResetPin},
	"UNLOCK_PIN":                   {UnlockPin},
	"ACTIVATE_DEACTIVATED":         {ActivateDeactivated},
	"ACTIVATE_PIN_UNLOCK":          {ActivatePinUnlock},
	"ENABLE_DISABLE_USSD":          {EnableDisableUSSD},
	"ENABLE_DISABLE_INAPP":         {EnableDisableInApp},
	"ENABLE_DISABLE_BOTH":          {EnableDisableBoth},
	"ACTIVATE_TERMINATED":          {ActivateTerminated},
	"ACTIVATE_CHANGE_PHONE_NUMBER": {ActivateChangePhoneNumber},
	"ENABLE_USSD_SUPERAPP":         {EnableUssdSupperapp},
	"DISABLE_USSD_SUPERAPP":        {DisableUssdSupperapp},
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
