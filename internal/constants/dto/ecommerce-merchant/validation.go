package ecommercemerchant

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (b BranchInformation) ValidateCreate() error {
	return validation.ValidateStruct(&b,

		validation.Field(&b.BranchCode,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchName,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchAddress,
			validation.When(b.BranchAddress != nil && *b.BranchAddress != "",
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&b.BranchOwner,
			validation.When(b.BranchOwner != nil && *b.BranchOwner != "",
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&b.BranchAccountNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[0-9]+$`)).
				Error("Invalid branch account number"),
		),
	)
}

func (b BranchInformation) ValidateBranchInfoUpdate() error {
	return validation.ValidateStruct(&b,

		validation.Field(&b.BranchCode,
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchName,
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchAddress,
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchOwner,
			validation.When(b.BranchOwner != nil && *b.BranchOwner != "",
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&b.BranchAccountNumber,
			validation.Match(regexp.MustCompile(`^[0-9]+$`)).
				Error("Invalid branch account number"),
		),
	)
}

func validateBranchesCreate(value interface{}) error {
	branches, ok := value.([]BranchInformation)
	if !ok {
		return fmt.Errorf("invalid branches format")
	}

	seen := make(map[string]bool)

	for i, b := range branches {

		if seen[b.BranchCode] {
			return fmt.Errorf("duplicate branch_code: %s", b.BranchCode)
		}
		seen[b.BranchCode] = true

		if err := b.ValidateCreate(); err != nil {
			return fmt.Errorf("branch[%d]: %w", i, err)
		}
	}

	return nil
}

func validateBranchesUpdate(value interface{}) error {
	branches, ok := value.([]BranchInformation)
	if !ok {
		return fmt.Errorf("invalid branches format")
	}

	seen := make(map[string]bool)

	for i, b := range branches {

		if b.BranchCode != "" {
			if seen[b.BranchCode] {
				return fmt.Errorf("duplicate branch_code: %s", b.BranchCode)
			}
			seen[b.BranchCode] = true
		}

		if err := b.ValidateBranchInfoUpdate(); err != nil {
			return fmt.Errorf("branch[%d]: %w", i, err)
		}
	}

	return nil
}

func (r EcommerceMerchant) ValidateCreate() error {
	return validation.ValidateStruct(&r,

		validation.Field(&r.MerchantName,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&r.MerchantCode,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&r.AccountNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[0-9]+$`)).
				Error("Invalid account number"),
		),

		validation.Field(&r.SettlementMethod,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&r.Branches,
			validation.By(validateBranchesCreate),
		),
		validation.Field(&r.IsEcommerceMerchant,
			validation.Required,
			validation.By(func(value interface{}) error {
				isEcommerceMerchant, ok := value.(*bool)
				if !ok || !*isEcommerceMerchant {
					return fmt.Errorf("merchant type should be ecommerce")
				}
				return nil
			}),
		),
	)
}

func (r UpdateEcommerceMerchant) ValidateUpdate() error {
	return validation.ValidateStruct(&r,

		// validation.Field(&r.IsEcommerceMerchant,
		// 	validation.Required,
		// 	validation.By(func(value interface{}) error {
		// 		isEcommerceMerchant, ok := value.(*bool)
		// 		if !ok || isEcommerceMerchant == nil {
		// 			return fmt.Errorf("is_ecommerce_merchant is required")
		// 		}
		// 		if !*isEcommerceMerchant {
		// 			return fmt.Errorf("is_ecommerce_merchant must be true")
		// 		}
		// 		return nil
		// 	}),
		// ),

		validation.Field(&r.MerchantName,
			validation.When(r.MerchantName != nil && *r.MerchantName != "",
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&r.MerchantCode,
			validation.When(r.MerchantCode != nil && *r.MerchantCode != "",
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&r.AccountNumber,
			validation.When(r.AccountNumber != nil && *r.AccountNumber != "",
				validation.Match(regexp.MustCompile(`^[0-9]+$`)).
					Error("Invalid account number"),
			),
		),

		validation.Field(&r.SettlementMethod,
			validation.When(r.SettlementMethod != nil && *r.SettlementMethod != "",
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
		),

		validation.Field(&r.Branches,
			validation.When(r.Branches != nil,
				validation.By(validateBranchesUpdate),
			),
		),
	)
}

func (r EnableOrDisableMerchantsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MerchantIDs,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(
				validation.Required,
				validation.By(utils.TrimWhiteSpace),
			),
		),
	)
}
