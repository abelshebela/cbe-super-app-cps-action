package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

type BranchInformation struct {
	BranchCode    string `json:"branch_code"`
	BranchName    string `json:"branch_name"`
	BranchAddress string `json:"branch_address"`
	BranchOwner   string `json:"branch_owner"`
}

func isBranchEmpty(branch model.BranchInformation) bool {
	return strings.TrimSpace(branch.BranchCode) == "" &&
		strings.TrimSpace(branch.BranchName) == "" &&
		strings.TrimSpace(branch.BranchAddress) == "" &&
		strings.TrimSpace(branch.BranchOwner) == ""
}

func (dto EcommerceMerchant) IsEmpty() bool {
	if strings.TrimSpace(dto.MerchantName) != "" ||
		// strings.TrimSpace(dto.PhoneNumber) != "" ||
		// strings.TrimSpace(dto.Email) != "" ||
		strings.TrimSpace(dto.AccountNumber) != "" {
		return false
	}

	if len(dto.Branches) == 0 {
		return true
	}

	for _, b := range dto.Branches {
		if !isBranchEmpty(b) {
			return false
		}
	}

	return true
}

func (dto EcommerceMerchant) Validate(isCreate bool) error {
	if !isCreate && dto.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	if isCreate {
		rules = []*validation.FieldRules{
			validation.Field(&dto.MerchantName,
				validation.Required.Error("merchant name is required"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
			validation.Field(&dto.MerchantCode,
				validation.By(func(value interface{}) error {
					validation.By(utils.TrimWhiteSpace)
					if err := validation.Required.Error("mercahnt code is required").Validate(value); err != nil {
						return err
					}
					validation.By(utils.NoSpecialChars)
					return nil
				}),
			),
			// validation.Field(&dto.PhoneNumber,
			// 	validation.Required.Error("phone number is required"),
			// 	validation.Match(regexp.MustCompile(`^(?:\+251|251|0)(9|7)\d{8}$`)).Error("invalid phone number format"),
			// 	validation.By(utils.TrimWhiteSpace),
			// ),
			validation.Field(&dto.SettlementMethod,
				validation.Required.Error("Settlement method is required"),
				validation.By(utils.TrimWhiteSpace),
				validation.By(utils.NoSpecialChars),
			),
			// validation.Field(&dto.Email,
			// 	validation.Required.Error("email is required"),
			// 	is.Email.Error("email must be a valid email address"),
			// 	validation.By(utils.TrimWhiteSpace),
			// ),
			validation.Field(&dto.AccountNumber,
				validation.By(func(value interface{}) error {
					validation.By(utils.TrimWhiteSpace)
					validation.By(utils.NoSpecialChars)
					validation.By(utils.NumbersOnly)
					if err := validation.Length(13, 13).Error("account number must be 13 digits").Validate(value); err != nil {
						return err
					}

					return nil
				}),
			),
			validation.Field(
				dto.IsEcommerceMerchant,
				validation.Required.Error(localization.ErrorEcommernceMerchantInvalidIsEcommerceMerchant.Error()),
				validation.By(func(value interface{}) error {
					v, ok := value.(*bool)
					if !ok || v == nil {
						return nil // Required will catch nil
					}
					if !*v {
						return localization.ErrorEcommernceMerchantInvalidIsEcommerceMerchant
					}
					return nil
				}),
			),
		}
	} else {
		if strings.TrimSpace(dto.MerchantName) != "" {
			rules = append(rules, validation.Field(&dto.MerchantName, validation.By(utils.NoSpecialChars)))
		}
		// if strings.TrimSpace(dto.PhoneNumber) != "" {
		// 	rules = append(rules, validation.Field(&dto.PhoneNumber,
		// 		validation.Length(9, 15).Error("phone number must be between 9 and 15 digits"),
		// 		validation.Match(regexp.MustCompile(`^(?:\+251|251|0)9\d{8}$`)).Error("invalid phone number format"),
		// 	))
		// }
		// if strings.TrimSpace(dto.Email) != "" {

		// 	rules = append(rules, validation.Field(&dto.Email,
		// 		is.Email.Error("email must be a valid email address"),
		// 		validation.By(utils.TrimWhiteSpace),
		// 	))
		// }

		if strings.TrimSpace(dto.AccountNumber) != "" {
			rules = append(rules, validation.Field(&dto.AccountNumber, validation.By(utils.NoSpecialChars)))
		}

		if strings.TrimSpace(dto.SettlementMethod) != "" {
			rules = append(rules,
				validation.Field(&dto.SettlementMethod,
					validation.By(utils.NoSpecialChars),
				))
		}
		if dto.IsEcommerceMerchant != nil {
			rules = append(rules,
				validation.Field(
					dto.IsEcommerceMerchant,
					validation.Required.Error(localization.ErrorEcommernceMerchantInvalidIsEcommerceMerchant.Error()),
					validation.By(func(value interface{}) error {
						v, ok := value.(*bool)
						if !ok || v == nil {
							return nil // Required will catch nil
						}
						if !*v {
							return localization.ErrorEcommernceMerchantInvalidIsEcommerceMerchant
						}
						return nil
					}),
				),
			)
		}
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&dto, rules...); err != nil {
			return err
		}
	}

	return nil
}
