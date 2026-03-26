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
)

var validRequestActions = map[RequestAction]struct{}{
	LinkAccount:                     {},
	LinkOtherAccount:                {},
	UnlinkAccount:                   {},
	ResetPin:                        {},
	ActivateAccount:                 {},
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
}

func IsValidRequestAction(requestAction string) bool {
	_, ok := validRequestActions[NormalizeRequestAction(requestAction)]
	return ok
}

var RequestActionGroups = map[string][]RequestAction{
	"LINK_ACCOUNT":                       {LinkAccount},
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
