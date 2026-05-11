package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserActionResponsibility string

const (
	MAKER   UserActionResponsibility = "MAKER"
	CHECKER UserActionResponsibility = "CHECKER"
	AUDITOR UserActionResponsibility = "AUDITOR"
)

type UserActionLog struct {
	ID                         bson.ObjectID            `bson:"_id,omitempty" json:"id"`
	ActionID                   bson.ObjectID            `bson:"action_id" json:"action_id"`
	ActionCode                 string                   `bson:"action_code" json:"action_code"`
	GivenActionStatus          string                   `bson:"given_action_status" json:"given_action_status"`
	GivenAuditorStatus         AuditorMark              `bson:"given_auditor_status" json:"given_auditor_status"`
	RequestAction              constants.RequestAction  `bson:"request_action" json:"request_action"`
	ActionTakenServiceName     string                   `bson:"action_taken_service_name" json:"action_taken_service_name"`
	ActionTakenServiceUniqueID string                   `bson:"action_taken_service_unique_id" json:"action_taken_service_unique_id"`
	CheckerLevel               string                   `bson:"checker_level" json:"checker_level"`
	AuditorLevel               string                   `bson:"auditor_level" json:"auditor_level"`
	UserID                     bson.ObjectID            `bson:"user_id" json:"user_id"`
	Username                   string                   `bson:"username" json:"username"`
	UserPhone                  string                   `bson:"user_phone" json:"user_phone"`
	UserRoleCode               string                   `bson:"user_role_code" json:"user_role_code"`
	UserActionResponsibilities UserActionResponsibility `bson:"user_action_responsibilities" json:"user_action_responsibilities"`
	IsDeleted                  bool                     `bson:"is_deleted" json:"is_deleted"`
	CreatedAt                  time.Time                `bson:"created_at" json:"created_at"`
	DeletedAt                  time.Time                `bson:"deleted_at" json:"deleted_at"`
}
