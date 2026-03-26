package vault

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

func noSpecialChars(value interface{}) error {
	var s string

	switch v := value.(type) {
	case string:
		s = v
	case *string:
		if v == nil {
			return nil
		}
		s = *v
	default:
		return nil
	}

	if s == "" {
		return nil
	}

	// Allow letters, numbers, space, dot, underscore, dash
	re := regexp.MustCompile(`^[a-zA-Z0-9 ._-]+$`)
	if !re.MatchString(s) {
		return errors.New("contains invalid characters (only letters, numbers, spaces, dots, underscores, and dashes are allowed)")
	}
	return nil
}

func validateCoverImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorBankImageMissingOrInvalid
	}

	if !utils.IsValidImage(file) {
		return errors.New(localization.MsgBankImageRequiredOrMissing)
	}

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}

// Validate validates CreateCategoryRequest according to interest type and tier rules.
//
// Rules:
//   - interest_type must be either "flat" or "dynamic" (case insensitive).
//   - FLAT: all tiers must have the same tier_interest value.
//   - DYNAMIC: tier_interest values can differ across tiers.
//   - Tiers numeric rules (for both types):
//   - all min/max must be numeric and > 0
//   - max > min
//   - for tier i>0, Min must equal previous tier's Max
func (r *CreateCategoryRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	r.Name = strings.TrimSpace(r.Name)
	r.InterestType = strings.TrimSpace(r.InterestType)

	if err := validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 100).Error("name must be between 3 and 100 characters"),
			validation.By(noSpecialChars),
		),
		validation.Field(&r.CoverImage,
			validation.When(r.CoverImage != nil,
				validation.By(func(value interface{}) error { return validateCoverImage(value) })),
		),
		validation.Field(&r.InterestType,
			validation.Required.Error("interest_type is required"),
		),
		validation.Field(&r.Tiers,
			validation.Required.Error("tiers are required"),
			validation.Length(1, 0).Error("at least one tier is required"),
		),
	); err != nil {
		return err
	}

	it := strings.ToUpper(r.InterestType)
	if it != "FLAT" && it != "DYNAMIC" {
		return errors.New("interest type must be either 'FLAT' or 'DYNAMIC'")
	}

	// Validate tiers (min/max chaining + interest-type behavior)
	if err := validateCreateTiers(it, r.Tiers); err != nil {
		return err
	}

	return nil
}

func validateCreateTiers(interestType string, tiers []CreateTierDTO) error {
	if len(tiers) == 0 {
		return errors.New("tiers are required")
	}

	var previousMax decimal.Decimal
	var flatTierInterest *decimal.Decimal
	for i, t := range tiers {
		name := strings.TrimSpace(t.Name)
		minStr := strings.TrimSpace(t.Min)
		maxStr := strings.TrimSpace(t.Max)
		tierInterest := strings.TrimSpace(t.TierInterest)

		if err := validation.Validate(name,
			validation.Required.Error(fmt.Sprintf("tier %d: name is required", i+1)),
			validation.Length(1, 100),
			validation.By(noSpecialChars),
		); err != nil {
			return err
		}

		if minStr == "" {
			return fmt.Errorf("tier %d: min is required", i+1)
		}
		if maxStr == "" {
			return fmt.Errorf("tier %d: max is required", i+1)
		}

		min, err := decimal.NewFromString(minStr)
		if err != nil {
			return fmt.Errorf("tier %d: min must be a valid number", i+1)
		}
		max, err := decimal.NewFromString(maxStr)
		if err != nil {
			return fmt.Errorf("tier %d: max must be a valid number", i+1)
		}

		if min.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("tier %d: min must be > 0", i+1)
		}
		if max.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("tier %d: max must be > 0", i+1)
		}
		if !max.GreaterThan(min) {
			return fmt.Errorf("tier %d: max must be greater than min", i+1)
		}

		if tierInterest == "" {
			return fmt.Errorf("tier %d: tier_interest is required", i+1)
		}
		tierInterestVal, err := decimal.NewFromString(strings.TrimSuffix(tierInterest, "%"))
		if err != nil {
			return fmt.Errorf("tier %d: tier_interest must be a valid number", i+1)
		}
		if tierInterestVal.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("tier %d: tier_interest must be > 0", i+1)
		}

		if interestType == "FLAT" {
			if flatTierInterest == nil {
				flatTierInterest = &tierInterestVal
			} else if !tierInterestVal.Equal(*flatTierInterest) {
				return errors.New("all tiers must have the same tier_interest when interest_type is 'FLAT'")
			}
		}

		// Chaining rules
		if i < 1 {
			if min.LessThanOrEqual(decimal.NewFromInt(0)) {
				return errors.New("tier 1: min must be greater than 0")
			}
		} else {
			if !min.Equal(previousMax) {
				return fmt.Errorf("tier %d: min (%s) must be equal to previous tier's max (%s)", i+1, min.String(), previousMax.String())
			}
		}
		previousMax = max

	}

	return nil
}

// Validate validates CreateWithdrawalRequest.
func (r *CreateWithdrawalRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	r.LockedVaultID = strings.TrimSpace(r.LockedVaultID)
	r.WithdrawerName = strings.TrimSpace(r.WithdrawerName)
	r.WithdrawerPhoneNumber = strings.TrimSpace(r.WithdrawerPhoneNumber)

	return validation.ValidateStruct(r,
		validation.Field(&r.LockedVaultID, validation.Required.Error("locked_vault_id is required")),
		validation.Field(&r.WithdrawalAmount, validation.Required.Error("withdrawal_amount is required")),
		validation.Field(&r.WithdrawerName,
			validation.Required.Error("withdrawer_name is required"),
			validation.Length(3, 100),
			validation.By(noSpecialChars),
		),
		validation.Field(&r.WithdrawerPhoneNumber,
			validation.Required.Error("withdrawer_phone_number is required"),
			validation.Length(9, 13).Error("withdrawer_phone_number must be between 9 and 13 characters"),
		),
	)
}

// Validate validates UpdateWithdrawalStatusRequest.
// func (r *UpdateWithdrawalStatusRequest) Validate() error {
// 	if r == nil {
// 		return errors.New("request is required")
// 	}

// 	// r.WithdrawalID = strings.TrimSpace(r.WithdrawalID)
// 	r.WithdrawalStatus = strings.TrimSpace(r.WithdrawalStatus)

// 	return validation.ValidateStruct(r,
// 		// validation.Field(&r.WithdrawalID, validation.Required.Error("withdrawal_id is required")),
// 		validation.Field(&r.WithdrawalStatus, validation.Required.Error("withdrawal_status is required"),
// 			validation.In("APPROVED", "REJECTED").Error("withdrawal status should be 'APPROVED' or 'REJECTED'"),
// 		),
// 	)
// }

func (r *UpdateCategoryRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	// Trim simple string fields
	if r.Name != nil {
		trimmed := strings.TrimSpace(*r.Name)
		*r.Name = trimmed
	}
	if r.InterestType != nil {
		trimmed := strings.TrimSpace(*r.InterestType)
		*r.InterestType = trimmed
	}

	// At least one updatable field must be provided
	if r.Name == nil && r.InterestType == nil && r.Tiers == nil && r.Deadlock == nil {
		return errors.New("at least one field (name, cover_image, interest_type, interest, tiers, deadlock) must be provided")
	}

	// Field-level validation
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.When(r.Name != nil,
				validation.Length(3, 100),
				validation.By(noSpecialChars),
			),
		),
		validation.Field(&r.CoverImage,
			validation.When(r.CoverImage != nil,
				validation.By(func(value interface{}) error { return validateCoverImage(value) })),
		),
		validation.Field(&r.InterestType, validation.By(func(value interface{}) error {
			ptr, ok := value.(*string)
			if !ok || ptr == nil || *ptr == "" {
				return nil
			}
			it := strings.ToUpper(strings.TrimSpace(*ptr))
			if it != "FLAT" && it != "DYNAMIC" {
				return errors.New("interest_type must be either 'flat' or 'dynamic'")
			}
			return nil
		})),
		// validation.Field(&r.CategoryInterest, validation.By(func(value interface{}) error {
		// 	ptr, ok := value.(*string)
		// 	if !ok || ptr == nil || *ptr == "" {
		// 		return nil
		// 	}
		// 	if _, err := decimal.NewFromString(*ptr); err != nil {
		// 		return errors.New("category_interest must be a valid decimal number")
		// 	}
		// 	return nil
		// })),
	); err != nil {
		return err
	}

	// Determine interest type for cross-field behavior.
	var it string
	if r.InterestType != nil && *r.InterestType != "" {
		it = strings.ToUpper(*r.InterestType)
		if it != "FLAT" && it != "DYNAMIC" {
			return errors.New("interest_type must be either 'flat' or 'dynamic'")
		}
	}

	// Optional tier validation when present
	if r.Tiers != nil {
		if err := validateUpdateTier(r.Tiers, it); err != nil {
			return err
		}
	}

	return nil
}

func validateUpdateTier(tiers []UpdateTierDTO, interestType string) error {
	if tiers == nil {
		return nil
	}

	var flatTierInterest *decimal.Decimal

	// Validate provided numeric fields
	for _, t := range tiers {
		var min, max *decimal.Decimal
		if t.Min != nil && strings.TrimSpace(*t.Min) != "" {
			v, err := decimal.NewFromString(strings.TrimSpace(*t.Min))
			if err != nil {
				return errors.New("tiers.min must be a valid number")
			}
			if v.LessThanOrEqual(decimal.Zero) {
				return errors.New("tiers.min must be > 0")
			}
			min = &v
		}
		if t.Max != nil && strings.TrimSpace(*t.Max) != "" {
			v, err := decimal.NewFromString(strings.TrimSpace(*t.Max))
			if err != nil {
				return errors.New("tiers.max must be a valid number")
			}
			if v.LessThanOrEqual(decimal.Zero) {
				return errors.New("tiers.max must be > 0")
			}
			max = &v
		}

		if min != nil && max != nil && !max.GreaterThan(*min) {
			return errors.New("tiers.max must be greater than tiers.min")
		}

		if t.TierInterest != nil && strings.TrimSpace(*t.TierInterest) != "" {
			raw := strings.TrimSpace(*t.TierInterest)
			clean := strings.TrimSpace(strings.TrimSuffix(raw, "%"))
			if clean == "" {
				return errors.New("tiers.interest must be a valid number")
			}
			v, err := decimal.NewFromString(clean)
			if err != nil {
				return errors.New("tiers.interest must be a valid number")
			}

			switch strings.ToUpper(interestType) {
			case "FLAT":
				if v.LessThanOrEqual(decimal.Zero) {
					return errors.New("tiers.interest must be > 0 when interest_type is 'flat'")
				}
				if flatTierInterest == nil {
					flatTierInterest = &v
				} else if !v.Equal(*flatTierInterest) {
					return errors.New("all provided tiers.interest values must be equal when interest_type is 'flat'")
				}
			case "DYNAMIC":
				// In dynamic mode, when provided it must be > 0
				if v.LessThanOrEqual(decimal.Zero) {
					return errors.New("tiers.interest must be > 0 when interest_type is 'dynamic'")
				}
			default:
				// Generic rule if we somehow don't know the type.
				if v.LessThan(decimal.Zero) {
					return errors.New("tiers.interest must be >= 0")
				}
			}
		}
	}

	return nil
}
