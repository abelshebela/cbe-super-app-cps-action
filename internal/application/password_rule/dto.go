package passwordrule

import (
	"cbe-super-app-cps-action/internal/domain/action"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RequestPasswordRuleUpdateDTO struct {
	Rule  action.PasswordRule `json:"rule"`
	Maker action.User         `json:"maker"`
}

func (r RequestPasswordRuleUpdateDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Rule.MinLength, validation.Required, validation.Min(1).Error("min_length must be greater than 0")),
	)
}

type ApproveOrRejectPasswordRuleActionDTO struct {
	ActionID        string      `json:"action_id"`
	Approve         bool        `json:"approve"`
	Checker         action.User `json:"checker"`
	RejectionReason *string     `json:"rejection_reason"`
}

func (r ApproveOrRejectPasswordRuleActionDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ActionID, validation.Required.Error("action_id is required")),
	)
}
