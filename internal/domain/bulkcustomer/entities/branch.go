package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserType string

const (
	Maker   UserType = "MAKER"
	Checker UserType = "CHECKER"
)

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type MakerAndChecker struct {
	Maker   User
	Checker User
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
	RequestEnableBank               RequestAction = "ENABLE_BANK"
	RequestDisableBank              RequestAction = "DISABLE_BANK"
	RequestEnableWallet             RequestAction = "ENABLE_WALLET"
	RequestDisableWallet            RequestAction = "DISABLE_WALLET"
	RequestUpdatePasswordExpiry     RequestAction = "UPDATE_PASSWORD_EXPIRY"
	RequestCreateValidation         RequestAction = "CREATE_VALIDATION"
	RequestUpdateValidation         RequestAction = "UPDATE_VALIDATION"
	RequestDeleteValidation         RequestAction = "DELETE_VALIDATION"
	RequestUpdateArchiveExpiry      RequestAction = "UPDATE_ARCHIVE_EXPIRY"
	RequestCreateServiceFee         RequestAction = "CREATE SERVICE FEE"
	RequestUpdateServiceFee         RequestAction = "UPDATE SERVICE FEE"
	RequestDeleteServiceFee         RequestAction = "DELETE SERVICE FEE"
	RequestCreateDailyLimit         RequestAction = "CREATE DAILY LIMIT"
	RequestUpdateDailyLimit         RequestAction = "UPDATE DAILY LIMIT"
	RequestDeleteDailyLimit         RequestAction = "DELETE DAILY LIMIT"
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
	RequestUpdateBlockTime          RequestAction = "UPDATE_BLOCK_TIME"
)

type Location struct {
	Type      string    `json:"type" bson:"type"`
	Locationt []float32 `json:"locationt" bson:"locationt"`
}

type User struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

type ActionData struct {
	BranchCode    string `json:"branch_code" bson:"branch_code"`
	BranchName    string `json:"branch_name" bson:"branch_name"`
	BranchAddress string `json:"branch_address" bson:"branch_address"`
	DistrictCode  string `json:"district_code" bson:"district_code"`
	DistrictName  string `json:"district_name" bson:"district_name"`
	BranchRegion  string `json:"branch_region" bson:"branch_region"`
}

type CPSAction struct {
	ID                string        `json:"id" bson:"id"`
	ActionCode        string        `json:"action_code" bson:"action_code"`
	CheckerUser       User          `json:"checker_user" bson:"checker_user"`
	MakerUser         User          `json:"maker_user" bson:"maker_user"`
	RejectedReason    string        `json:"rejected_reason" bson:"rejected_reason"`
	Department        string        `json:"department" bson:"department"`
	Status            ActionStatus  `json:"status" bson:"status"`
	PreviousAction    interface{}   `json:"previous_action" bson:"previous_action"`
	RequestAction     RequestAction `json:"request_action" bson:"request_action"`
	ActionType        ActionType    `json:"action_type" bson:"action_type"`
	ActionData        ActionData    `json:"action_data" bson:"action_data"` // for single branch
	MakerActionTime   time.Time     `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime time.Time     `json:"checker_action_time" bson:"checker_action_time"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt    time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}
type Branch struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BranchCode    string        `json:"branchCode" bson:"branchCode"`
	BranchName    string        `json:"branchName" bson:"branchName"`
	BranchAddress string        `json:"branchAddress" bson:"branchAddress"`
	DistrictCode  string        `json:"districtCode" bson:"districtCode"`
	DistrictName  string        `json:"districtName" bson:"districtName"`
	BranchRegion  string        `json:"branchRegion" bson:"branchRegion"`
	RecordStat    string        `json:"RecordStat" bson:"RecordStat"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`
	Version       int           `json:"__v" bson:"__v"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
}
