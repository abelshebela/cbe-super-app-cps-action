package entities

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/enums"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/type_definition"
)

type CPSAction struct {
	ID                string                     `json:"id,omitempty"`
	ActionCode        string                     `json:"action_code" bson:"action_code"`
	CheckerUser       type_definition.User       `json:"checker_user" bson:"checker_user"`
	MakerUser         type_definition.User       `json:"maker_user" bson:"maker_user"`
	RejectedReason    string                     `json:"rejected_reason,omitempty" bson:"rejected_reason"`
	Department        string                     `json:"department,omitempty" bson:"department"`
	Status            enums.ActionStatus         `json:"status,omitempty" bson:"status"`
	RequestAction     enums.RequestAction        `json:"request_action,omitempty" bson:"request_action"`
	ActionType        enums.ActionType           `json:"action_type,omitempty" bson:"action_type"`
	ActionData        type_definition.ActionData `json:"action_data" bson:"action_data"`
	PreviousAction    any                        `json:"previous_action,omitempty"`
	CurrentAction     any                        `json:"current_action,omitempty"`
	MakerActionTime   time.Time                  `json:"maker_action_time,omitzero" bson:"maker_action_time"`
	CheckerActionTime time.Time                  `json:"checker_action_time,omitzero" bson:"checker_action_time"`
	CreatedAt         time.Time                  `json:"created_at,omitempty" bson:"created_at"`
	LastModifiedAt    time.Time                  `json:"last_modified_at,omitempty" bson:"last_modified_at"`
	UniqueID          string                     `json:"unique_id,omitempty" bson:"unique_id"`
}
