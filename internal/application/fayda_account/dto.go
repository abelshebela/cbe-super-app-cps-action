package faydaaccount

import (
	"fmt"
	"regexp"
	"time"

	Validation "github.com/go-ozzo/ozzo-validation/v4"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"user_code"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
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
	UseCode     string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

func (a ActionData) Validate() error {
	return Validation.ValidateStruct(&a,
		Validation.Field(&a.UseCode, validation.Required.Error("user code is required")),
		validation.Field(&a.FullName, validation.Required.Error("full name is required")),
		validation.Field(&a.PhoneNumber, validation.Required.Error("phone number is required"),
			validation.By(ValidatePhone)),
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
	ID                string        `json:"id"`
	ActionCode        string        `json:"action_code"`
	CheckerUser       User          `json:"checker_user"`
	MakerUser         User          `json:"maker_user"`
	RejectedReason    string        `json:"rejected_reason"`
	Department        string        `json:"department"`
	Status            ActionStatus  `json:"status"`
	RequestAction     RequestAction `json:"request_action"`
	ActionType        ActionType    `json:"action_type"`
	ActionData        ActionData    `json:"action_data"`
	MakerActionTime   time.Time     `json:"maker_action_time"`
	CheckerActionTime time.Time     `json:"checker_action_time"`
}

func (c CPSAction) Validate(rejectOrApprove string) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ActionData),
		Validation.Field(&c.RejectedReason, validation.When(
			rejectOrApprove == "REJECT", validation.Required, validation.Min(30).Error("minimum character should be 30"),
			validation.Max(300).Error("maximum character should be 300"))),
	)
}
