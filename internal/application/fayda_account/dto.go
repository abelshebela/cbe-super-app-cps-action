package faydaaccount

import (
	"fmt"
	"regexp"
	"time"

	Validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
	RequestDisableFaydaAccount RequestAction = "DISABLE_FAYDA_ACCOUNT"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type ActionData struct {
	UseCode     string `json:"user_code"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

func (a ActionData) Validate() error {
	return Validation.ValidateStruct(&a,
		Validation.Field(&a.UseCode, Validation.Required.Error("user code is required")),
		Validation.Field(&a.FullName, Validation.Required.Error("full name is required")),
		Validation.Field(&a.PhoneNumber, Validation.Required.Error("phone number is required"),
			Validation.By(ValidatePhone)),
	)
}

func ValidatePhone(phone interface{}) error {
	phoneStr := fmt.Sprintf("%v", phone)
	re := regexp.MustCompile(`^(?:\+?251|0)?([97]\d{8})$`)

	if matches := re.FindStringSubmatch(phoneStr); matches != nil {
		return nil
	}
	return fmt.Errorf("invalid phone number")
}

type CPSAction struct {
	ID                string        `json:"id,omitempty"`
	ActionCode        string        `json:"action_code,omitempty"`
	CheckerUser       User          `json:"checker_user"`
	MakerUser         User          `json:"maker_user"`
	RejectedReason    string        `json:"rejected_reason,omitempty"`
	Department        string        `json:"department,omitempty"`
	Status            ActionStatus  `json:"status,omitempty"`
	RequestAction     RequestAction `json:"request_action,omitempty"`
	ActionType        ActionType    `json:"action_type,omitempty"`
	ActionData        ActionData    `json:"action_data"`
	PreviousData      any           `json:"previous_action,omitempty"`
	CurrentData       any           `json:"current_action,omitempty"`
	MakerActionTime   time.Time     `json:"maker_action_time,omitzero"`
	CheckerActionTime time.Time     `json:"checker_action_time,omitzero"`
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
