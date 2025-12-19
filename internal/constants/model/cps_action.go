package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSAction struct {
	ID                  bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	ActionCode          string          `bson:"action_code" json:"action_code"`
	UniqueId            string          `bson:"unique_id" json:"unique_id,omitempty"`
	MakerID             string          `bson:"maker_id" json:"maker_id"`
	MakerName           string          `bson:"maker_name" json:"maker_name"`
	MakerPhoneNumber    string          `bson:"maker_phone_number" json:"maker_phone_number"`
	CheckerUsers        []types.Checker `bson:"checker_users" json:"checker_users"`
	CheckerCount        int32           `bson:"checker_count" json:"checker_count"`
	CurrentCheckerIndex float32         `bson:"current_checker_index" json:"current_checker_index"`
	Department          string          `bson:"department" json:"department"`
	RejectionReason     string          `bson:"rejection_reason" json:"rejection_reason,omitempty"`
	PreviousAction      interface{}     `bson:"previous_action" json:"previous_action,omitempty"`
	CurrentAction       interface{}     `bson:"current_action" json:"current_action,omitempty"`
	ActionStatus        string          `bson:"action_status" json:"action_status,omitempty"`
	ActionType          string          `bson:"action_type" json:"action_type,omitempty"`
	IsDeleted           bool            `bson:"is_deleted" json:"is_deleted,omitempty"`
	RequestAction       string          `bson:"request_action" json:"request_action"`
	ReversedByRoleID    string          `bson:"reversed_by_role_id" json:"reversed_by_role_id,omitempty"`
	ReversedByID        string          `bson:"reversed_by_id" json:"reversed_by_id,omitempty"`
	ReversedByName      string          `bson:"reversed_by_name" json:"reversed_by_name,omitempty"`
	ReversedAt          time.Time       `bson:"reversed_at" json:"reversed_at,omitempty"`
	CreatedAt           time.Time       `bson:"created_at" json:"created_at,omitempty"`
	LastModifiedAt      time.Time       `bson:"last_modified_at" json:"last_modified_at,omitempty"`
	MakerActionTime     time.Time       `bson:"maker_action_time" json:"maker_action_time,omitempty"`
	CheckerActionTime   *time.Time      `bson:"checker_action_time" json:"checker_action_time,omitempty"`
}

type CheckCPSAction struct {
	UserCode      string
	FullName      string
	PhoneNumber   string
	Department    string
	RequestAction string
}
