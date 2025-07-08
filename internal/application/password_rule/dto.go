package passwordrule

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RequestPasswordRuleUpdateDTO struct {
	ID   string              `json:"id"`
	Rule action.PasswordRule `json:"rule"`
}

func (r RequestPasswordRuleUpdateDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.Rule, validation.Required, validation.By(func(value interface{}) error {
			if v, ok := value.(action.PasswordRule); ok {
				if v.MinLength < 1 {
					return validation.NewError("min_length", "min_length must be greater than 0")
				}
				if v.PasswordID == "" {
					return validation.NewError("password_id", "password_id is required")
				}
				
				
			}
			return nil
		})),
	)
}

type ApproveOrRejectPasswordRuleActionDTO struct {
	ActionID        string      `json:"action_code"`
	Decision        string      `json:"decison"`
	
	RejectionReason *string     `json:"rejection_reason"`
	// Department      string      `json:"department"`
}

func (r ApproveOrRejectPasswordRuleActionDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ActionID, validation.Required.Error("action_id is required")),
		validation.Field(&r.Decision, validation.Required.Error("decision is required")),

		// validation.Field(&r.Department, validation.Required.Error("department is required")),
	)
}
