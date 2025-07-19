package cpsactions

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
)

type CPSAction struct {
	ID                 string                 `json:"id"`
	ActionCode         string                 `json:"action_code"`
	UniqueID           string                 `json:"unique_id,omitempty"`
	MakerID            string                 `json:"maker_id"`
	MakerName          string                 `json:"maker_name"`
	MakerPhoneNumber   string                 `json:"maker_phone_number"`
	CheckerID          string                 `json:"checker_id,omitempty"`
	CheckerName        string                 `json:"checker_name,omitempty"`
	CheckerPhoneNumber string                 `json:"checker_phone_number,omitempty"`
	Department         string                 `json:"department"`
	RejectionReason    string                 `json:"rejection_reason,omitempty"`
	PreviousAction     any                    `json:"previous_action"`
	CurrentAction      any                    `json:"current_action"`
	ActionStatus       constant.ActionStatus  `json:"action_status"`
	ActionType         constant.ActionType    `json:"action_type"`
	RequestAction      constant.RequestAction `json:"request_action"`
	CreatedAt          time.Time              `json:"created_at"`
	LastModifiedAt     time.Time              `json:"last_modified_at"`
	MakerActionTime    time.Time              `json:"maker_action_time"`
	CheckerActionTime  time.Time              `json:"checker_action_time"`
}

type User struct {
	UserCode    string
	FullName    string
	PhoneNumber string
}

type AuthorizeCPSAction struct {
	ActionCode        string    `json:"action_code"`
	Department        string    `json:"department,omitempty"`
	RejectionReason   string    `json:"rejection_reason"`
	CheckerUser       User      `json:"checker_user"`
	CheckerActionTime time.Time `json:"checker_action_time,omitzero"`
}
