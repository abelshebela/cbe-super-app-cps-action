package service

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type KYCLevel string

const (
	KYCLevelZero KYCLevel = "ZERO"
	KYCLevelOne  KYCLevel = "ONE"
	KYCLevelTwo  KYCLevel = "TWO"
)

type ProductCodes struct {
	PRD    string
	VATPRD string
	SFPRD  string
	TRXN   string
}

type GLEntry struct {
	ProductAccount    string
	ProductBranchCode string
	ServiceAccount    string
	ServiceBranchCode string
	VatAccount        string
	VatBranchCode     string
}

type Tier struct {
	ID        string
	Min       uint64
	Max       uint64
	FeeAmount uint64
}

type Cap struct {
	KYCLevel  KYCLevel
	SingleCap uint64
	DailyCap  uint64
	MinAmount uint64
	MaxAmount uint64
}

type Service struct {
	ID                 string       `json:"id" bson:"id"`
	ServiceCode        string       `json:"service_code" bson:"service_code"`
	ServiceName        string       `json:"service_name" bson:"service_name"`
	ServiceType        string       `json:"service_type" bson:"service_type"`
	Key                string       `json:"key" bson:"key"`
	Cap                Cap          `json:"cap" bson:"cap"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	AboveAmount        uint64       `json:"above_amount" bson:"above_amount"`
	AboveServiceFee    uint64       `json:"above_service_fee" bson:"above_service_fee"`
	PaymentType        string       `json:"payment_type" bson:"payment_type"`
	Tiers              []Tier       `json:"tiers" bson:"tiers"`
	CBEGLEntry         GLEntry      `json:"cbe_gl_entry" bson:"cbe_gl_entry"`
	CBEIFBGLEntry      GLEntry      `json:"cbe_ifb_gl_entry" bson:"cbe_ifb_gl_entry"`
	Enabled            bool         `json:"enabled" bson:"enabled"`
	IsDeleted          bool         `json:"is_deleted" bson:"is_deleted"`
	CreatedAt          time.Time    `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time    `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt          time.Time    `json:"deleted_at" bson:"deleted_at"`
}

type User struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type ActionData struct {
	Tier []Tier
}
type RequestAction string

const (
	RequestServiceFeeUpdate RequestAction = "SERVICE_FEE_UPDATE"
	RequestServiceFeeCreate RequestAction = "SERVICE_FEE_CREATE"
)

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type CPSAction struct {
	ID                 string        `json:"id" bson:"id"`
	ActionCode         string        `json:"action_code" bson:"action_code"`
	UniqueId           string        `json:"unique_id" bson:"unique_id"`
	MakerID            string        `json:"maker_id" bson:"maker_id"`
	MakerName          string        `json:"maker_name" bson:"maker_name"`
	MakerPhoneNumber   string        `json:"maker_phone_number" bson:"maker_phone_number"`
	CheckerID          string        `json:"checker_id" bson:"checker_id"`
	CheckerName        string        `json:"checker_name" bson:"checker_name"`
	CheckerPhoneNumber string        `json:"checker_phone_number" bson:"checker_phone_number"`
	Department         string        `json:"department" bson:"department"`
	RejectionReason    string        `json:"rejection_reason" bson:"rejection_reason"`
	PreviosAction      any           `json:"previous_action" bson:"previous_action"`
	CurrentAction      any           `json:"current_action" bson:"current_action"`
	ActionStatus       ActionStatus  `json:"action_status" bson:"action_status"`
	ActionType         ActionType    `json:"action_type" bson:"action_type"`
	RequestAction      RequestAction `json:"request_action" bson:"request_action"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	MakerActionTime    time.Time     `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime  time.Time     `json:"checker_action_time" bson:"checker_action_time"`
}

type UpdateServiceDetailsResponse struct {
	ActionID string `json:"action_id"`
}

func (s Service) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ServiceCode, validation.Required.Error("service code cannot be empty")),
		validation.Field(&s.ServiceName, validation.Required.Error("service name cannot be empty")),
		validation.Field(&s.ServiceType, validation.Required.Error("service type cannot be empty")),
	)
}
