package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSActionApproveIndex struct {
	ID            bson.ObjectID `json:"_id" bson:"_id"`
	RoleId        string        `json:"role_id" bson:"role_id"`
	ActionName    string        `json:"action_name" bson:"action_name"`
	ViewerIndex   *int64        `json:"viewer_index" bson:"viewer_index"`
	MakerIndex    *int64        `json:"maker_index" bson:"maker_index"`
	CheckerIndex  *float64      `json:"checker_index" bson:"checker_index"`
	AuditorIndex  *float64      `json:"auditor_index" bson:"auditor_index"`
	ApproverCount int8          `json:"approver_count" bson:"approver_count"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}
