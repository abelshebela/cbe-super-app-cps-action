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
    ActionDisable ActionType = "DISABLE"
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
    ID                 bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
    ActionCode         string          `json:"action_code" bson:"action_code"`
    MakerID            string          `json:"maker_id" bson:"maker_id"`
    MakerName          string          `json:"maker_name" bson:"maker_name"`
    MakerPhoneNumber   string          `json:"maker_phone_number" bson:"maker_phone_number"`
    CheckerID          *string         `json:"checker_id,omitempty" bson:"checker_id,omitempty"`
    CheckerName        *string         `json:"checker_name,omitempty" bson:"checker_name,omitempty"`
    CheckerPhoneNumber *string         `json:"checker_phone_number,omitempty" bson:"checker_phone_number,omitempty"`
    Department         string          `json:"department" bson:"department"`
    RejectionReason    *string         `json:"rejection_reason,omitempty" bson:"rejection_reason,omitempty"`
    PreviosAction      []bson.M        `json:"previos_action" bson:"previos_action"`     // array of JSON objects
    CurrentAction      []bson.M        `json:"current_action" bson:"current_action"`     // array of JSON objects
    ActionStatus       ActionStatus    `json:"action_status" bson:"action_status"`
    ActionType         ActionType      `json:"action_type" bson:"action_type"`
    RequestAction      RequestAction   `json:"request_action" bson:"request_action"`
    BranchCode         string          `json:"branch_code" bson:"branch_code"`
    CreatedAt          time.Time       `json:"created_at" bson:"created_at"`
    LastModifiedAt     time.Time       `json:"last_modified_at" bson:"last_modified_at"`
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
	Type      string
	Locationt []float32
}

type Branch struct {
    ID             bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    BranchCode     string        `json:"branchCode" bson:"branchCode"`
    BranchName     string        `json:"branchName" bson:"branchName"`
    BranchAddress  string        `json:"branchAddress" bson:"branchAddress"`
    DistrictCode   string        `json:"districtCode" bson:"districtCode"`
    DistrictName   string        `json:"districtName" bson:"districtName"`
    BranchRegion   string        `json:"branchRegion" bson:"branchRegion"`
    Enable         bool          `json:"enable" bson:"enable"`
    RecordStat     string        `json:"RecordStat" bson:"RecordStat"`
    CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
    UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
    Version        int           `json:"__v" bson:"__v"`
    Enabled        bool          `json:"enabled" bson:"enabled"`
}