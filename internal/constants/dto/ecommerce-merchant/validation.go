package miniappmerchant

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"regexp"
	"strings"

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

	if isCreate {
		if regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(strings.TrimSpace(dto.MerchantName)) == false {
			return errors.New("Invalid merchant name is required")
		}
		if regexp.MustCompile(`^[a-zA-Z0-9_\s]+$`).MatchString(strings.TrimSpace(dto.MerchantCode)) == false {
			return errors.New("Invalid merchant code is required")
		}
		if regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(strings.TrimSpace(dto.SettlementMethod)) == false {
			return errors.New("Invalid settlement method is required")
		}
		if regexp.MustCompile(`^[0-9]+$`).MatchString(strings.TrimSpace(dto.AccountNumber)) == false {
			return errors.New("Invalid account number is required")
		}
		if dto.IsEcommerceMerchant == nil || !*dto.IsEcommerceMerchant {
			return errors.New(localization.ErrorEcommernceMerchantInvalidIsEcommerceMerchant.Message)
		}

	} else {
		// if regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(strings.TrimSpace(dto.MerchantName)) == false {
		// 	return errors.New("Invalid merchant name is required")
		// }
		if regexp.MustCompile(`^[a-zA-Z0-9_\s]+$`).MatchString(strings.TrimSpace(dto.MerchantCode)) == false {
			return errors.New("Invalid merchant code is required")
		}
		if regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(strings.TrimSpace(dto.SettlementMethod)) == false {
			return errors.New("Invalid settlement method is required")
		}
		if regexp.MustCompile(`^[0-9]+$`).MatchString(strings.TrimSpace(dto.AccountNumber)) == false {
			return errors.New("Invalid account number is required")
		}
		if dto.IsEcommerceMerchant == nil || !*dto.IsEcommerceMerchant {
			return errors.New(localization.ErrorEcommernceMerchantInvalidIsEcommerceMerchant.Message)
		}

	}

	return nil
}
