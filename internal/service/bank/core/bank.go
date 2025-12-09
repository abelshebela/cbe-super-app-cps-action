package bank_core

import (
	"cbe-super-app-cps-action/internal/constants/model"
)

// Bank_mapper maps a map[string]interface{} to a model.Bank.
// Only sets fields if they exist in the action map.
func Bank_mapper(action map[string]interface{}) model.Bank {
	bank := model.Bank{}

	if v, ok := action["name"]; ok {
		if name, ok := v.(string); ok {
			bank.Name = name
		}
	}
	if v, ok := action["logo"]; ok {
		if logo, ok := v.(string); ok {
			bank.Logo = logo
		}
	}
	if v, ok := action["code"]; ok {
		if code, ok := v.(string); ok {
			bank.Code = code
		}
	}
	if v, ok := action["bic"]; ok {
		if bic, ok := v.(string); ok {
			bank.BIC = bic
		}
	}
	if v, ok := action["account_length"]; ok {
		if accountLength, ok := v.(float64); ok {
			bank.AccountLength = int(accountLength)
		}
	}
	if v, ok := action["enabled"]; ok {
		if enabled, ok := v.(bool); ok {
			bank.Enabled = enabled
		}
	}
	if v, ok := action["is_deleted"]; ok {
		if isDeleted, ok := v.(bool); ok {
			bank.IsDeleted = isDeleted
		}
	}

	return bank
}
