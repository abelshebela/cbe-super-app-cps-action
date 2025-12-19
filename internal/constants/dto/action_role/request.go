package actionrole_dto

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateActionRoleRequest struct {
	ActionCode           string     `json:"action_code"`
	ActionName           string     `json:"action_name"`
	AssignedMakersRoles  []string   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
	AssignedCheckerRoles [][]string `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
	AssignedAuditorRoles []string   `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
	IsMakerOnly          bool       `json:"is_maker_only" bson:"is_maker_only"`
}

type UpdateActionRoleRequest struct {
	ActionName           string     `json:"action_name"`
	AssignedMakersRoles  []string   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
	AssignedCheckerRoles [][]string `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
	AssignedAuditorRoles []string   `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
	IsMakerOnly          bool       `json:"is_maker_only" bson:"is_maker_only"`
}

type GetActionRoleByActionCodeRes struct {
	ID                   bson.ObjectID  `json:"_id" bson:"_id"`
	ActionCode           string         `json:"action_code" bson:"action_code"`
	ActionName           string         `json:"action_name" bson:"action_name"`
	AssignedMakersRoles  []model.Role   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
	AssignedCheckerRoles [][]model.Role `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
	AssignedAuditorRoles []model.Role   `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
	Enabled              bool           `json:"enabled" bson:"enabled"`
	UpdatedAt            time.Time      `json:"updated_at" bson:"updated_at"`
	CreatedAt            time.Time      `json:"created_at" bson:"created_at"`
}
