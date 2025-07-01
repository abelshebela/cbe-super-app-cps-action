package passwordrule

import "cbe-super-app-cps-action/internal/domain/action"

type RequestPasswordRuleUpdateDTO struct {
	Rule  action.PasswordRule `json:"rule"`
	Maker action.User         `json:"maker"`
}

type ApproveOrRejectPasswordRuleActionDTO struct {
	ActionID        string      `json:"action_id"`
	Approve         bool        `json:"approve"`
	Checker         action.User `json:"checker"`
	RejectionReason *string     `json:"rejection_reason"`
}
