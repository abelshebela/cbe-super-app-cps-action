package dto

import (
	// "fmt"
	// "regexp"
	"time"
	// "time"

	Validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
)

type ServiceDetailsResponse struct {
	ID                 string               `json:"id"`
	ServiceID          string               `json:"service_id"`
	ServiceCode        string               `json:"service_code"`
	ServiceName        string               `json:"service_name"`
	ServiceType        string               `json:"service_type"`
	Key                string               `json:"key"`
	Cap                service.Cap          `json:"cap"`
	CBEProductCodes    service.ProductCodes `json:"cbe_product_codes"`
	CBEIFBProductCodes service.ProductCodes `json:"cbe_ifb_product_codes"`
	AboveAmount        uint64               `json:"above_amount"`
	AboveServiceFee    uint64               `json:"above_service_fee"`
	PaymentType        string               `json:"payment_type"`
	Tiers              []service.Tier       `json:"tiers"`
	CBEGLEntry         service.GLEntry      `json:"cbe_gl_entry"`
	CBEIFBGLEntry      service.GLEntry      `json:"cbe_ifb_gl_entry"`
	Enabled            bool                 `json:"enabled"`
	IsDeleted          bool                 `json:"is_deleted"`
	CreatedAt          time.Time            `json:"created_at"`
	LastModifiedAt     time.Time            `json:"last_modified_at"`
	DeletedAt          time.Time            `json:"deleted_at,omitempty"`
}

type UpdateServiceDetailsRequest struct {
	ID                 string               `json:"id"`
	ServiceID          string               `json:"service_id"`
	ServiceCode        string               `json:"service_code"`
	ServiceName        string               `json:"service_name"`
	ServiceType        string               `json:"service_type"`
	Key                string               `json:"key"`
	Cap                service.Cap          `json:"cap"`
	CBEProductCodes    service.ProductCodes `json:"cbe_product_codes"`
	CBEIFBProductCodes service.ProductCodes `json:"cbe_ifb_product_codes"`
	AboveAmount        uint64               `json:"above_amount"`
	AboveServiceFee    uint64               `json:"above_service_fee"`
	PaymentType        string               `json:"payment_type"`
	Tiers              []service.Tier       `json:"tiers"`
	CBEGLEntry         service.GLEntry      `json:"cbe_gl_entry"`
	CBEIFBGLEntry      service.GLEntry      `json:"cbe_ifb_gl_entry"`
	Enabled            bool                 `json:"enabled"`
	MakerID            string               `json:"maker_id"`
}

type UpdateServiceDetailsResponse struct {
	ActionID string `json:"action_id"`
}

type ApproveServiceDetailsRequest struct {
	ActionID        string `json:"action_id"`
	Approve         bool   `json:"approve"`
	CheckerID       string `json:"checker_id"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

type ServiceDetailsErrorResponse struct {
	Error string `json:"error"`
}

type SingleCapServiceRequest struct {
	ServiceId string `json:"service_id"`
	SingleCap uint64 `json:"single_cap"`
}

type DailyCapServiceRequest struct {
	ServiceId string `json:"service_id"`
	DailyCap  uint64 `json:"daily_cap"`
}

type TotalTransferCapRequest struct {
	ServiceId string `json:"service_id"`
	TotalCap  uint64 `json:"total_cap"`
}

type WholeCapServiceRequest struct {
	ServiceId string `json:"service_id"`
	SingleCap uint64 `json:"single_cap"`
	DailyCap  uint64 `json:"daily_cap"`
	TotalCap  uint64 `json:"total_cap"`
}

type ServiceFeeMakerRequest struct {
	ServiceId string `json:"service_id"`
	Tries     []service.Tier
}

type User struct {
	UserCode    string `json:"user_code,omitempty" bson:"user_code"`
	FullName    string `json:"full_name,omitempty" bson:"user_code"`
	PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number"`
}

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type RequestAction string

const (
	RequestServiceFeeUpdate RequestAction = "SERVICE_FEE_UPDATE"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type ActionData struct {
	Tier []service.Tier
}

// func (a ActionData) Validate() error {
// 	return Validation.ValidateStruct(&a,
// 		Validation.Field(&a.UseCode, Validation.Required.Error("user code is required")),
// 		Validation.Field(&a.FullName, Validation.Required.Error("full name is required")),
// 		Validation.Field(&a.PhoneNumber, Validation.Required.Error("phone number is required"),
// 			Validation.By(ValidatePhone)),
// 	)
// }

// func ValidatePhone(phone interface{}) error {
// 	phoneStr := fmt.Sprintf("%v", phone)
// 	re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

// 	if matches := re.FindStringSubmatch(phoneStr); matches != nil {
// 		return nil
// 	}
// 	return fmt.Errorf("invalid phone number")
// }

type CPSAction struct {
	ID                string         `json:"id,omitempty"`
	ActionCode        string         `json:"action_code,omitempty"`
	CheckerUser       User           `json:"checker_user"`
	MakerUser         User           `json:"maker_user"`
	RejectedReason    string         `json:"rejected_reason,omitempty"`
	Department        string         `json:"department,omitempty"`
	Status            ActionStatus   `json:"status,omitempty"`
	RequestAction     RequestAction  `json:"request_action,omitempty"`
	ActionType        ActionType     `json:"action_type,omitempty"`
	ActionData        []service.Tier `json:"action_data"`
	PreviousData      any            `json:"previous_action,omitempty"`
	CurrentData       any            `json:"current_action,omitempty"`
	MakerActionTime   time.Time      `json:"maker_action_time,omitzero"`
	CheckerActionTime time.Time      `json:"checker_action_time,omitzero"`
}

type RejectCPSAction struct {
	RejectedReason string `json:"rejected_reason,omitempty"`
	ActionData
}

func (c CPSAction) Validate(rejectOrApprove string) error {
	return Validation.ValidateStruct(&c,
		Validation.Field(&c.ActionData),
		Validation.Field(&c.RejectedReason, Validation.When(
			rejectOrApprove == "REJECT", Validation.Required, Validation.Length(30, 300).Error("length of reason should be between 30 and 300 character"),
		)),
	)
}
