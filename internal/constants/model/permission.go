package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserType string

const (
	Maker   UserType = "MAKER"
	Checker UserType = "CHECKER"
)

type ActionStatus string

const (
	ActionPending  ActionStatus = "PENDING"
	ActionApproved ActionStatus = "APPROVED"
	ActionRejected ActionStatus = "REJECTED"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type MakerAndChecker struct {
	Maker   User
	Checker User
}

type RequestAction string

type Permission struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PermissionName string        `bson:"permission_name" json:"permission_name"`
	CreatedAt      time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
