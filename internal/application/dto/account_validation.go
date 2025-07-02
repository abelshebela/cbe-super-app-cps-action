package dto

import (
	"cbe-super-app-cps-action/internal/domain/account_validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetAccountValidationResponse struct {
	Validation account_validation.ValidationRule `json:"validation"`
}

type UpdateAccountValidationRequest struct {
	ID         string                            `json:"id"`
	Validation account_validation.ValidationRule `json:"validation"`
}

func (r UpdateAccountValidationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r.Validation.Identifier, validation.Required.Error("identifier is required")),
	)
}

type UpdateAccountValidationResponse struct {
	ActionID string `json:"action_id"`
}

type ApproveRejectRequest struct {
	ActionID string `json:"action_id"`
	Approve  bool   `json:"approve"`
}
