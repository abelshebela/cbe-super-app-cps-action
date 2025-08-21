package cpsactions

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
)

type CPSAction struct {
	ID                 string                 `json:"id"`
	ActionCode         string                 `json:"action_code"`
	UniqueID           string                 `json:"unique_id,omitempty"`
	MakerID            string                 `json:"maker_id"`
	MakerName          string                 `json:"maker_name"`
	MakerPhoneNumber   string                 `json:"maker_phone_number"`
	CheckerID          string                 `json:"checker_id,omitempty"`
	CheckerName        string                 `json:"checker_name,omitempty"`
	CheckerPhoneNumber string                 `json:"checker_phone_number,omitempty"`
	Department         string                 `json:"department"`
	RejectionReason    string                 `json:"rejection_reason,omitempty"`
	PreviousAction     any                    `json:"previous_action"`
	CurrentAction      any                    `json:"current_action"`
	ActionStatus       constant.ActionStatus  `json:"action_status"`
	ActionType         constant.ActionType    `json:"action_type"`
	RequestAction      constant.RequestAction `json:"request_action"`
	CreatedAt          time.Time              `json:"created_at"`
	LastModifiedAt     time.Time              `json:"last_modified_at"`
	MakerActionTime    time.Time              `json:"maker_action_time"`
	CheckerActionTime  *time.Time             `json:"checker_action_time"`
}

type User struct {
	UserCode    string
	FullName    string
	PhoneNumber string
	Department  string
}

type CheckCPSAction struct {
	UserCode      string
	FullName      string
	PhoneNumber   string
	Department    string
	RequestAction string
}

type AuthorizeCPSAction struct {
	ActionCode        string    `json:"action_code"`
	Department        string    `json:"department,omitempty"`
	RejectionReason   string    `json:"rejection_reason"`
	CheckerUser       User      `json:"checker_user"`
	CheckerActionTime time.Time `json:"checker_action_time,omitzero"`
}

type CreateCPSAction struct {
	ActionCode      string                 `json:"action_code"`
	MakerUser       User                   `json:"maker_user"`
	Department      string                 `json:"department,omitempty"`
	Status          constant.ActionStatus  `json:"status,omitempty"`
	RequestAction   constant.RequestAction `json:"request_action,omitempty"`
	ActionType      constant.ActionType    `json:"action_type,omitempty"`
	ActionData      any                    `json:"action_data"`
	PreviousData    any                    `json:"previous_action,omitempty"`
	CurrentData     any                    `json:"current_action,omitempty"`
	MakerActionTime time.Time              `json:"maker_action_time,omitzero"`
}

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type CreateCPSRequest struct {
	User          User
	CurData       any
	PrevData      any
	RequestAction constant.RequestAction
	ActionStatus  constant.ActionStatus
	ActionType    constant.ActionType
}

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type RequestAction string

const (
	RequestUser                     RequestAction = "USER"
	RequestPermissionGroup          RequestAction = "PERMISSION_GROUP"
	RequestDepartment               RequestAction = "DEPARMTENT"
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
	RequestEnableBank               RequestAction = "ENABLE_BANK"
	RequestDisableBank              RequestAction = "DISABLE_BANK"
	RequestEnableWallet             RequestAction = "ENABLE_WALLET"
	RequestDisableWallet            RequestAction = "DISABLE_WALLET"
	RequestUpdatePasswordExpiry     RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreateValidation         RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation         RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation         RequestAction = "DELETE_VALIDATION"
	RequestUpdateArchiveExpiry      RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestCreateServiceFee         RequestAction = "CREATE_SERVICE_FEE"
	RequestUpdateServiceFee         RequestAction = "UPDATE SERVICE FEE"
	RequestDeleteServiceFee         RequestAction = "DELETE SERVICE FEE"
	RequestCreateDailyLimit         RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit         RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit         RequestAction = "DELETE DAILY LIMIT"
	RequestBudgetUpdate             RequestAction = "UPDATE_BUDGET_CATEGORY"
	RequestBudgetCreate             RequestAction = "CREATE_BUDGET_CATEGORY"
	RequestBudgetDelete             RequestAction = "DELETE_BUDGET_CATEGORY"
	RequestBudgetColor              RequestAction = "BUDGET_COLOR"
	RequestBudgetIcon               RequestAction = "BUDGET_ICON"
	RequestUpdateProduct            RequestAction = "UPDATE_PRODUCT"
	RequestCreatePublicNotification RequestAction = "CREATE_PUBLIC_NOTIFICATION"
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
	RequestCreateEvent              RequestAction = "CREATE_EVENT"
	RequestUpdateEvent              RequestAction = "UPDATE_EVENT"
	RequestCreateEventCategory      RequestAction = "CREATE_EVENT_CATEGORY"
	RequestUpdateEventCategory      RequestAction = "UPDATE_EVENT_CATEGORY"
	RequestDisableEvent             RequestAction = "DISABLE_EVENT"
	RequestCreateMiniAppMerchant    RequestAction = "CREATE_MINIAPP_MERCHANT"
	RequestUpdateMiniAppMerchant    RequestAction = "UPDATE_MINIAPP_MERCHANT"
	RequestDeleteMiniAppMerchant    RequestAction = "DELETE_MINIAPP_MERCHANT"
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
	RequestUpdateAccountValidation  RequestAction = "UPDATE_ACCOUNT_VALIDATION"
	RequestUpdateServiceDetails     RequestAction = "UPDATE_SERVICE_DETAILS"
	RequestUpdateHQArchiveTime      RequestAction = "UPDATE_HQ_ARCHIVE_TIME"
	RequestUpdateEevent             RequestAction = "UPDATE_EVENT"
	RequestCIFRemove                RequestAction = "CIF_REMOVE"
	RequestServiceFlagUpdate        RequestAction = "SERVICE_FLAG_UPDATE"
)
