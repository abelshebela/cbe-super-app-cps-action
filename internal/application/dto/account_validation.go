package dto

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetAccountValidationResponse struct {
	Validation ValidationRuleDTO `json:"validation"`
}

type ValidationRuleDTO struct {
	ID            string `json:"id"`
	EntityType    string `json:"entity_type"`
	ValidationFor string `json:"validation_for"`
	Identifier    string `json:"identifier"`
	MinLength     uint8  `json:"min_length"`
	MaxLength     uint8  `json:"max_length"`
	Enabled       bool   `json:"enabled"`
	IsDeleted     bool   `json:"is_deleted"`
	ServiceID     string `json:"service_id"`
}

func (v ValidationRuleDTO) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Identifier, validation.Required.Error("identifier is required")),
	)
}

type UpdateAccountValidationRequest struct {
	ValidationRuleDTO
}

func (r UpdateAccountValidationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required.Error("id is required")),
		validation.Field(&r, validation.By(func(value interface{}) error {
			if v, ok := value.(ValidationRuleDTO); ok {
				return v.Validate()
			}
			return nil
		})),
	)
}

// DecisionEnum as a string type with constants
// DecisonEnum represents the allowed values for decision
// Use string type for enum-like behavior in Go

type UpdateAccountValidationResponse struct {
	ActionID string `json:"action_id"`
}

type ApproveRejectRequest struct {
	ActionCode     string            `json:"action_code"`
	Decison        utils.DecisonEnum `json:"decison"`
	RejectedReason string            `json:"rejected_reason"`
}

// Validate checks that Decison is valid and if DENIED, RejectedReason is required
func (r ApproveRejectRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ActionCode, validation.Required.Error("action_code is required")),
		validation.Field(&r.Decison, validation.Required.Error("decison is required"), validation.In(utils.DecisionApproved, utils.DecisionDenied).Error("decison must be either APPROVED or DENIED")),
		validation.Field(&r.RejectedReason, validation.By(func(value interface{}) error {
			if r.Decison == utils.DecisionDenied {
				if str, ok := value.(string); !ok || str == "" {
					return validation.NewError("validation_rejected_reason", "rejected_reason is required when decison is DENIED")
				}
			}
			return nil
		})),
	)
}

// Conversion function from domain to DTO
func ToValidationRuleDTO(rule account_validation.ValidationRule) ValidationRuleDTO {
	return ValidationRuleDTO{
		ID:            rule.ID,
		EntityType:    rule.EntityType,
		ValidationFor: rule.ValidationFor,
		Identifier:    rule.Identifier,
		MinLength:     rule.MinLength,
		MaxLength:     rule.MaxLength,
		Enabled:       rule.Enabled,
		IsDeleted:     rule.IsDeleted,
		ServiceID:     rule.ServiceID,
	}
}

// Conversion function from DTO to domain
func ToDomainValidationRule(dto ValidationRuleDTO) account_validation.ValidationRule {
	return account_validation.ValidationRule{
		ID:            dto.ID,
		EntityType:    dto.EntityType,
		ValidationFor: dto.ValidationFor,
		Identifier:    dto.Identifier,
		MinLength:     dto.MinLength,
		MaxLength:     dto.MaxLength,
		Enabled:       dto.Enabled,
		IsDeleted:     dto.IsDeleted,
		ServiceID:     dto.ServiceID,
	}
}
