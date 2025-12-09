package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ActionRole struct {
  ID                   bson.ObjectID     `json:"_id" bson:"_id"`
  ActionCode           string            `json:"action_code" bson:"action_code"`
  ActionName           string            `json:"action_name" bson:"action_name"`
  AssignedMakersRoles  []bson.ObjectID   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
  AssignedCheckerRoles [][]bson.ObjectID `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
  AssignedAuditorRoles []bson.ObjectID   `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
  ApproverCount        int32             `json:"approver_count" bson:"approver_count"`
  IsMakerOnly          bool              `json:"is_maker_only" bson:"is_maker_only"`
  Enabled              bool              `json:"enabled" bson:"enabled"`
  UpdatedAt            time.Time         `json:"updated_at" bson:"updated_at"`
  CreatedAt            time.Time         `json:"created_at" bson:"created_at"`
}



