package model

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	ID             bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key            string        `json:"key" bson:"key"`
	ServiceName    string        `json:"service_name" bson:"service_name"`
	SingleCap      float32       `json:"single_cap" bson:"single_cap"`
	MinAmount      float32       `json:"min_amount" bson:"min_amount"`
	Flag           bool          `json:"flag" bson:"flag"`
	DailyCap       float32       `json:"daily_cap" bson:"daily_cap"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modifed_at" bson:"last_modifed_at"`
}

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
	Maker   bps.BPSUser `bson:"maker"`
	Checker bps.BPSUser `bson:"checker"`
}

type CurrentAction struct {
	Id     string `bson:"id"`
	Action bool   `bson:"action"`
}

type CPSAction struct {
	ID                 bson.ObjectID `bson:"id"`
	ActionCode         string        `bson:"action_code"` // Generated
	MakerID            string        `bson:"maker_id"`
	MakerName          string        `bson:"maker_name"`
	MakerPhoneNumber   string        `bson:"maker_phone_number"`
	Unique_ID          string        `bson:"unique_id"`
	CheckerID          *string       `bson:"checker_id,omitempty"`
	CheckerName        *string       `bson:"checker_name,omitempty"`
	CheckerPhoneNumber *string       `bson:"checker_phone_number,omitempty"`
	Department         string        `bson:"department"`
	RejectionReason    *string       `bson:"rejection_reason,omitempty"`
	PreviosAction      struct{}      `bson:"previos_action"`
	CurrentAction      CurrentAction `bson:"current_action"`
	ActionStatus       ActionStatus  `bson:"action_status"`
	ActionType         ActionType    `bson:"action_type"`
	RequestAction      RequestAction `bson:"request_action"`
	CreatedAt          time.Time     `bson:"created_at"`
	LastModifiedAt     time.Time     `bson:"last_modified_at"`
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

type RegistrationType string

const (
	RegistrationTypeNew    RegistrationType = "NEW"
	RegistrationTypeLinked RegistrationType = "LINKED"
)

type LinkedAccount struct {
	ID                bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID            bson.ObjectID    `json:"user_id" bson:"user_id"`
	CustomerNumber    string           `json:"customer_number" bson:"customer_number"` // cif
	AccountNumber     string           `json:"account_number" bson:"account_number"`
	AccountHolderName string           `json:"account_holder_name" bson:"account_holder_name"`
	AccountType       string           `json:"account_type" bson:"account_type"`
	BranchCode        string           `json:"branch_code" bson:"branch_code"`
	LinkedStatus      bool             `json:"linked_status" bson:"linked_status"`
	LastLinkedStatus  bool             `json:"last_linked_status" bson:"last_linked_status"`
	LinkedAt          time.Time        `json:"linked_at" bson:"linked_at"`
	LinkedBranch      string           `json:"linker_branch" bson:"linker_branch"`
	RegistrationType  RegistrationType `json:"registration_type" bson:"registration_type"`
	IsAccountActive   bool             `json:"is_account_active" bson:"is_account_active"`
	AndOrStatus       bool             `json:"and_or_status" bson:"and_or_status"`
	AccountBranchCode string           `json:"account_branch_code" bson:"account_branch_code"`
	CurrencyCode      string           `json:"currency" bson:"curreny"`
	IsMain            bool             `json:"is_main" bson:"is_main"` // default: false, first account: true
	MakerAndChecker   struct {
		Linkers struct {
			Maker   string `json:"maker" bson:"maker"`
			Checker string `json:"checker" bson:"checker"`
		} `json:"linkers" bson:"linkers"`
		Unlinkers struct {
			Maker   string `json:"maker" bson:"maker"`
			Checker string `json:"checker" bson:"checker"`
		} `json:"unlinkers" bson:"unlinkers,omitempty"`
	} `json:"maker_and_checker" bson:"maker_and_checker"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at"  bson:"updated_at"`
}
