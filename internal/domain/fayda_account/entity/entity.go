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
	ID                string        `json:"id" bson:"id"`
	ActionCode        string        `json:"action_code" bson:"action_code"`
	CheckerUser       User          `json:"checker_user" bson:"checker_user"`
	MakerUser         User          `json:"maker_user" bson:"maker_user"`
	RejectedReason    string        `json:"rejected_reason" bson:"rejected_reason"`
	Department        string        `json:"department" bson:"department"`
	Status            ActionStatus  `json:"status" bson:"status"`
	PreviousData      any           `json:"previous_action" bson:"previous_action"`
	CurrentData       any           `json:"current_action" bson:"current_action"`
	RequestAction     RequestAction `json:"request_action" bson:"request_action"`
	ActionType        ActionType    `json:"action_type" bson:"action_type"`
	ActionData        ActionData    `json:"action_data" bson:"action_data"`
	MakerActionTime   time.Time     `json:"maker_action_time" bson:"maker_action_time"`
	CheckerActionTime time.Time     `json:"checker_action_time" bson:"checker_action_time"`
}
