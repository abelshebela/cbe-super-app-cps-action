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
	ID                 string        `json:"id"`
	ActionCode         string        `json:"action_code"`
	UniqueId           string        `json:"unique_id"`
	MakerID            string        `json:"maker_id"`
	MakerName          string        `json:"maker_name"`
	MakerPhoneNumber   string        `json:"maker_phone_number"`
	CheckerID          string        `json:"checker_id"`
	CheckerName        string        `json:"checker_name"`
	CheckerPhoneNumber string        `json:"checker_phone_number"`
	Department         string        `json:"department"`
	RejectionReason    *string       `json:"rejection_reason"`
	PreviousAction     interface{}   `json:"previos_action"`
	CurrentAction      interface{}   `json:"current_action"`
	ActionStatus       ActionStatus  `json:"action_status"`
	ActionType         ActionType    `json:"action_type"`
	RequestAction      RequestAction `json:"request_action"`
	CreatedAt          time.Time     `json:"created_at"`
	LastModifiedAt     time.Time     `json:"last_modified_at"`
	MakerActionTime    time.Time     `json:"maker_action_time"`
	CheckerActionTime  time.Time     `json:"checker_action_time"`
}

type RejectCPSAction struct {
	RejectedReason string `json:"rejected_reason,omitempty"`
}

func (c CPSAction) Validate(rejectOrApprove string) error {
	return Validation.ValidateStruct(&c,
		Validation.Field(&c.CurrentAction),
		Validation.Field(&c.RejectionReason, Validation.When(
			rejectOrApprove == "REJECT", Validation.Required, Validation.Length(30, 300).Error("length of reason should be between 30 and 300 character"),
		)),
	)
}
