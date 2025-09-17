package dto

import (
	// "fmt"
	// "regexp"
	"time"
	// "time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	Validation "github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
)

// type ServiceDetailsResponse struct {
// 	ID                 string               `json:"id"`
// 	ServiceID          string               `json:"service_id"`
// 	ServiceCode        string               `json:"service_code"`
// 	ServiceName        string               `json:"service_name"`
// 	ServiceType        string               `json:"service_type"`
// 	Key                string               `json:"key"`
// 	Cap                service.Cap          `json:"cap"`
// 	CBEProductCodes    service.ProductCodes `json:"cbe_product_codes"`
// 	CBEIFBProductCodes service.ProductCodes `json:"cbe_ifb_product_codes"`
// 	AboveAmount        uint64               `json:"above_amount"`
// 	AboveServiceFee    uint64               `json:"above_service_fee"`
// 	PaymentType        string               `json:"payment_type"`
// 	Tiers              []service.Tier       `json:"tiers"`
// 	CBEGLEntry         service.GLEntry      `json:"cbe_gl_entry"`
// 	CBEIFBGLEntry      service.GLEntry      `json:"cbe_ifb_gl_entry"`
// 	Enabled            bool                 `json:"enabled"`
// 	IsDeleted          bool                 `json:"is_deleted"`
// 	CreatedAt          time.Time            `json:"created_at"`
// 	LastModifiedAt     time.Time            `json:"last_modified_at"`
// 	DeletedAt          time.Time            `json:"deleted_at,omitempty"`
// }

// type UpdateServiceDetailsRequest struct {
// 	ID                 string               `json:"id"`
// 	ServiceID          string               `json:"service_id"`
// 	ServiceCode        string               `json:"service_code"`
// 	ServiceName        string               `json:"service_name"`
// 	ServiceType        string               `json:"service_type"`
// 	Key                string               `json:"key"`
// 	Cap                service.Cap          `json:"cap"`
// 	CBEProductCodes    service.ProductCodes `json:"cbe_product_codes"`
// 	CBEIFBProductCodes service.ProductCodes `json:"cbe_ifb_product_codes"`
// 	AboveAmount        uint64               `json:"above_amount"`
// 	AboveServiceFee    uint64               `json:"above_service_fee"`
// 	PaymentType        string               `json:"payment_type"`
// 	Tiers              []service.Tier       `json:"tiers"`
// 	CBEGLEntry         service.GLEntry      `json:"cbe_gl_entry"`
// 	CBEIFBGLEntry      service.GLEntry      `json:"cbe_ifb_gl_entry"`
// 	Enabled            bool                 `json:"enabled"`
// 	MakerID            string               `json:"maker_id"`
// }

type UpdateServiceDetailsResponse struct {
	ActionID string `json:"action_id" bson:"action_id"`
}

type ApproveServiceDetailsRequest struct {
	ActionID        string `json:"action_id" bson:"action_id"`
	Approve         bool   `json:"approve" bson:"approve"`
	CheckerID       string `json:"checker_id" bson:"checker_id"`
	RejectionReason string `json:"rejection_reason,omitempty" bson:"rejection_reason,omitempty"`
}

type ServiceDetailsErrorResponse struct {
	Error string `json:"error" bson:"error"`
}

type SingleCapServiceRequest struct {
	ServiceId string `json:"service_id" bson:"service_id"`
	SingleCap uint64 `json:"single_cap" bson:"single_cap"`
}

type DailyCapServiceRequest struct {
	ServiceId string `json:"service_id" bson:"service_id"`
	DailyCap  uint64 `json:"daily_cap" bson:"daily_cap"`
}

type Tier struct {
	ID        string `json:"id" bson:"id"`
	Min       uint64 `json:"min" bson:"min"`
	Max       uint64 `json:"max" bson:"max"`
	FeeAmount uint64 `json:"fee_amount" bson:"fee_amount"`
}

type ServiceFeeMakerRequest struct {
	ServiceID string       `json:"service_id" bson:"service_id"`
	Tries     []model.Tier `json:"tries" bson:"tries"`
}

type TotalTransferCapRequest struct {
	ServiceId string `json:"service_id" bson:"service_id"`
	TotalCap  uint64 `json:"total_cap" bson:"total_cap"`
}

type WholeCapServiceRequest struct {
	ServiceId string `json:"service_id" bson:"service_id"`
	SingleCap uint64 `json:"single_cap" bson:"single_cap"`
	DailyCap  uint64 `json:"daily_cap" bson:"daily_cap"`
	TotalCap  uint64 `json:"total_cap" bson:"total_cap"`
}

// type ServiceFeeMakerRequest struct {
// 	ServiceId string `json:"service_id"`
// 	Tries     []service.Tier
// }

type User struct {
	UserCode    string `json:"user_code,omitempty" bson:"user_code,omitempty"`
	FullName    string `json:"full_name,omitempty" bson:"full_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
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
	Tier []model.Tier `json:"tier" bson:"tier"`
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
// 	}
// 	return fmt.Errorf("invalid phone number")
// }

type CPSAction struct {
	ID                string        `json:"id,omitempty" bson:"id,omitempty"`
	ActionCode        string        `json:"action_code,omitempty" bson:"action_code,omitempty"`
	CheckerUser       User          `json:"checker_user" bson:"checker_user"`
	MakerUser         User          `json:"maker_user" bson:"maker_user"`
	RejectedReason    string        `json:"rejected_reason,omitempty" bson:"rejected_reason,omitempty"`
	Department        string        `json:"department,omitempty" bson:"department,omitempty"`
	Status            ActionStatus  `json:"status,omitempty" bson:"status,omitempty"`
	RequestAction     RequestAction `json:"request_action,omitempty" bson:"request_action,omitempty"`
	ActionType        ActionType    `json:"action_type,omitempty" bson:"action_type,omitempty"`
	ActionData        []model.Tier  `json:"action_data" bson:"action_data"`
	PreviousData      any           `json:"previous_action,omitempty" bson:"previous_action,omitempty"`
	CurrentData       any           `json:"current_action,omitempty" bson:"current_action,omitempty"`
	MakerActionTime   time.Time     `json:"maker_action_time,omitzero" bson:"maker_action_time,omitzero"`
	CheckerActionTime time.Time     `json:"checker_action_time,omitzero" bson:"checker_action_time,omitzero"`
}

type RejectCPSAction struct {
	RejectedReason string `json:"rejected_reason,omitempty" bson:"rejected_reason,omitempty"`
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
