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
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchOwner,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchAccountNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[0-9]+$`)).
				Error("Invalid branch account number"),
		),
	)
}

func (b BranchInformation) ValidateUpdate() error {
	return validation.ValidateStruct(&b,

		validation.Field(&b.BranchCode,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchName,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchAddress,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&b.BranchOwner,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
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

		if err := b.ValidateUpdate(); err != nil {
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
	)
}

func (r UpdateEcommerceMerchant) ValidateUpdate() error {
	return validation.ValidateStruct(&r,

		validation.Field(&r.MerchantName,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&r.MerchantCode,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&r.AccountNumber,
			validation.Match(regexp.MustCompile(`^[0-9]+$`)).
				Error("Invalid account number"),
		),

		validation.Field(&r.SettlementMethod,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),

		validation.Field(&r.Branches,
			validation.By(validateBranchesUpdate),
		),
	)
}
