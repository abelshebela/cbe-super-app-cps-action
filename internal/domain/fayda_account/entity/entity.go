package entity

import "time"

type User struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
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
	RejectionReason    *string       `json:"rejection_reason" bson:"rejection_reason"`
	PreviosAction      interface{}   `json:"previos_action" bson:"previos_action"`
	CurrentAction      interface{}   `json:"current_action" bson:"current_action"`
	ActionStatus       ActionStatus  `json:"action_status" bson:"action_status"`
	ActionType         ActionType    `json:"action_type" bson:"action_type"`
	RequestAction      RequestAction `json:"request_action" bson:"request_action"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	MakerActionTime    time.Time     `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime  time.Time     `json:"checker_action_time" bson:"checker_action_time"`
}
