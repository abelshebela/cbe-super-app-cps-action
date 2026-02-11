package cps_actionrole_dto

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateActionRoleRequest struct {
	ActionCode           string     `json:"action_code"`
	ActionName           string     `json:"action_name"`
	PortalCardName       string     `json:"portal_card_name"`
	AssignedViewersRoles []string   `json:"assigned_viewers_roles" bson:"assigned_viewers_roles"`
	AssignedMakersRoles  []string   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
	AssignedCheckerRoles [][]string `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
	AssignedAuditorRoles [][]string `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
	IsMakerOnly          bool       `json:"is_maker_only" bson:"is_maker_only"`
	IsViewOnly           bool       `json:"is_view_only" bson:"is_view_only"`
}

type UpdateActionRoleRequest struct {
	ActionName           string     `json:"action_name"`
	PortalCardName       string     `json:"portal_card_name"`
	AssignedViewersRoles []string   `json:"assigned_viewers_roles" bson:"assigned_viewers_roles"`
	AssignedMakersRoles  []string   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
	AssignedCheckerRoles [][]string `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
	AssignedAuditorRoles [][]string `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
	IsMakerOnly          bool       `json:"is_maker_only" bson:"is_maker_only"`
	IsViewOnly           bool       `json:"is_view_only" bson:"is_view_only"`
}

type AuditorMarkRequest struct {
	Mark   string `json:"mark"`
	Reason string `json:"reason"`
}
type GetActionRoleByActionCodeRes struct {
	ID                   bson.ObjectID     `json:"_id" bson:"_id"`
	ActionCode           string            `json:"action_code" bson:"action_code"`
	ActionName           string            `json:"action_name" bson:"action_name"`
	AssignedViewersRoles []model.JobRole   `json:"assigned_viewers_roles" bson:"assigned_viewers_roles"`
	AssignedMakersRoles  []model.JobRole   `json:"assigned_makers_roles" bson:"assigned_makers_roles"`
	AssignedCheckerRoles [][]model.JobRole `json:"assigned_checkers_roles" bson:"assigned_checkers_roles"`
	AssignedAuditorRoles [][]model.JobRole `json:"assigned_auditor_roles" bson:"assigned_auditor_roles"`
	ApproverCount        int64             `json:"approver_count" bson:"approver_count"`
	IsMakerOnly          bool              `json:"is_maker_only" bson:"is_maker_only"`
	Version              int64             `json:"version" bson:"version"`
	Enabled              bool              `json:"enabled" bson:"enabled"`
	UpdatedAt            time.Time         `json:"updated_at" bson:"updated_at"`
	CreatedAt            time.Time         `json:"created_at" bson:"created_at"`
}
