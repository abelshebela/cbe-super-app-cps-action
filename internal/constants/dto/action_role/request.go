package actionrole_dto

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateActionRoleRequest struct {
	ActionCode       string     `json:"action_code"`
	ActionName       string     `json:"action_name"`
	AssignedMakers   []string   `json:"assigned_makers"`
	AssignedCheckers [][]string `json:"assigned_checkers"`
}

type UpdateActionRoleRequest struct {
	ActionName       string     `json:"action_name"`
	AssignedMakers   []string   `json:"assigned_makers"`
	AssignedCheckers [][]string `json:"assigned_checkers"`
}

type GetActionRoleByActionCodeRes struct {
	ID              bson.ObjectID  `json:"_id" bson:"_id"`
	ActionCode      string         `json:"action_code" bson:"action_code"`
	ActionName      string         `json:"action_name" bson:"action_name"`
	AssignedMakers  []model.Role   `json:"assigned_makers" bson:"assigned_makers"`
	AssignedChecker [][]model.Role `json:"assigned_checkers" bson:"assigned_checkers"`
	Enabled         bool           `json:"enabled" bson:"enabled"`
	UpdatedAt       time.Time      `json:"updated_at" bson:"updated_at"`
	CreatedAt       time.Time      `json:"created_at" bson:"created_at"`
}
