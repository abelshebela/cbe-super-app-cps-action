package dto

import "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"

type GetAccountValidationResponse struct {
	Validation account_validation.ValidationRule `json:"validation"`
}

type UpdateAccountValidationRequest struct {
	ID         string                             `json:"id"`
	Validation account_validation.ValidationRule `json:"validation"`
}

type UpdateAccountValidationResponse struct {
	ActionID string `json:"action_id"`
}

type ApproveRejectRequest struct {
	ActionID string `json:"action_id"`
	Approve  bool   `json:"approve"`
}