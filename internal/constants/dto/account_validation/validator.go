package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v ValidationRuleDTO) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Identifier, validation.Required.Error("identifier is required")),
	)
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

// Validate checks that Decison is valid and if DENIED, RejectedReason is required
func (r ApproveRejectRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RejectedReason, validation.By(func(value interface{}) error {
			if r.Decison == DecisionDenied {
				if str, ok := value.(string); !ok || str == "" {
					return validation.NewError("validation_rejected_reason", "rejected_reason is required when decison is DENIED")
				}
			}
			return nil
		})),
	)
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RejectedReason, validation.By(func(value interface{}) error {
			if !r.Decison {
				if str, ok := value.(string); !ok || str == "" {
					return validation.NewError("validation_rejected_reason", "rejected_reason is required when decison is DENIED")
				}
			}
			return nil
		})),
	)
}

// Conversion function from domain to DTO
func ToValidationRuleDTO(rule ValidationRule) ValidationRuleDTO {
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

// ToDomainValidationRule Conversion function from DTO to domain
func ToDomainValidationRule(dto ValidationRuleDTO) ValidationRule {
	return ValidationRule{
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

func ToModel(dto ValidationRuleDTO) *model.ValidationRule {
	return &model.ValidationRule{
		EntityType:     dto.EntityType,
		ValidationFor:  dto.ValidationFor,
		Identifier:     dto.Identifier,
		MinLength:      int(dto.MinLength),
		MaxLength:      int(dto.MaxLength),
		Enabled:        dto.Enabled,
		IsDeleted:      dto.IsDeleted,
		ServiceID:      dto.ServiceID,
		LastModifiedAt: time.Now(),
	}
}
