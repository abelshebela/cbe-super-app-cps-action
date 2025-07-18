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

type User struct {
	UserID      string
	FullName    string
	PhoneNumber string
	Timestamp   time.Time
}

type MakerAndChecker struct {
	Maker   User
	Checker User
}

type CPSAction struct {
	ID                 bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ActionCode         string        ` bson:"action_code" json:"action_code"`
	MakerID            string        `bson:"maker_id" json:"maker_id"`
	MakerName          string        `bson:"maker_name" json:"maker_name"`
	MakerPhoneNumber   string        `bson:"maker_phone_number" json:"maker_phone_number"`
	CheckerID          string        `bson:"checker_id" json:"checker_id,omitempty"`
	CheckerName        string        `bson:"checker_name" json:"checker_name,omitempty"`
	CheckerPhoneNumber string        `bson:"checker_phone_number" json:"checker_phone_number,omitempty"`
	Department         string        `bson:"department" json:"department"`
	RejectionReason    *string       `bson:"rejection_reason" json:"rejection_reason,omitempty"`
	PreviousAction     interface{}   `bson:"previous_action" json:"previous_action"`
	CurrentAction      interface{}   `bson:"current_action" json:"current_action"`
	ActionStatus       ActionStatus  `bson:"action_status" json:"action_status"`
	ActionType         ActionType    `bson:"action_type" json:"action_type"`
	RequestAction      RequestAction `bson:"request_action" json:"request_action"`
	CreatedAt          time.Time     `bson:"created_at" json:"created_at"`
	LastModifiedAt     time.Time     `bson:"last_modified_at" json:"last_modified_at"`
	MakerActionTime    time.Time     `bson:"maker_action_time" json:"maker_action_time"`
	CheckerActionTime  time.Time     `bson:"checker_action_time" json:"checker_action_time"`
}

type RequestAction string

type Department struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DepartmentCode    string        `bson:"department_code" json:"department_code"`
	Department        string        `bson:"department" json:"department"`
	PermissionGroupID bson.ObjectID `bson:"permission_group_id" json:"permission_group_id"`
	PortalCards       []string      `bson:"portal_cards" json:"portal_cards"`
	Enabled           bool          `bson:"enabled" json:"enabled"`
	IsDeleted         bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt         time.Time     `bson:"created_at" json:"created_at"`
	LastModified      time.Time     `bson:"last_modified" json:"last_modified"`
}

type Color struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Color     string        `bson:"color" json:"color"`
	Enabled   bool          `bson:"enabled" json:"enabled"`
	IsDeleted bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type Icon struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Icon         string        `bson:"icon" json:"icon"`
	Enabled      bool          `bson:"enabled" json:"enabled"`
	IsDeleted    bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	LastModified time.Time     `bson:"last_modified" json:"last_modified"`
}

type FetchIconResponse struct {
	Icons []*Icon `json:"icons"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int64  `json:"total"`
}

const (
	RequestUser                     RequestAction = "USER"
	RequestPermissionGroup          RequestAction = "PERMISSION_GROUP"
	RequestDepartment               RequestAction = "DEPARTMENT"
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
