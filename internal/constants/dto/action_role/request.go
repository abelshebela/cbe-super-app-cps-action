package actionrole_dto

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
