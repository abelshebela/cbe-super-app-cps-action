package bank_core

import (
	"encoding/json"
	"fmt"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
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
	if v, ok := action["bic_code"]; ok {
		if bicCode, ok := v.(string); ok {
			bank.BICCode = bicCode
		}
	}
	if v, ok := action["is_enabled"]; ok {
		if enabled, ok := v.(bool); ok {
			bank.Enabled = enabled
		}
	}
	if v, ok := action["is_deleted"]; ok {
		if isDeleted, ok := v.(bool); ok {
			bank.IsDeleted = isDeleted
		}
	}
	if v, ok := action["type"]; ok {
		if bankType, ok := v.(string); ok {
			bank.Type = bankType
		}
	}

	return bank
}
func Bank_oracle_mapper(action map[string]interface{}) imodel.BankOracle {
	bank := imodel.BankOracle{}

	// CurrentAction is JSON from BankOracle (tags: bank_name, bic_code, …) or legacy/alternate keys.
	bank.BankName = stringField(action, "bank_name", "bankname", "name")
	bank.Logo = stringField(action, "logo")
	bank.BICCode = stringField(action, "bic_code", "biccode", "bic")
	bank.IsEnabled = intField(action, "is_enabled", "isenabled", "enabled")
	bank.HasAlphaNumeric = intField(action, "has_alpha_numeric", "hasalphanumeric")
	bank.AccountLength = intField(action, "account_length", "accountlength")
	bank.ID = stringField(action, "id")
	bank.CreateAt = stringField(action, "create_at", "created_at")
	bank.UpdateAt = stringField(action, "update_at", "last_modified_at")

	if v, ok := action["is_cbe"]; ok {
		switch x := v.(type) {
		case bool:
			if x {
				bank.IS_CBE = 1
			} else {
				bank.IS_CBE = 0
			}
		case float64:
			if int(x) != 0 {
				bank.IS_CBE = 1
			} else {
				bank.IS_CBE = 0
			}
		case int:
			if x != 0 {
				bank.IS_CBE = 1
			} else {
				bank.IS_CBE = 0
			}
		}
	}

	return bank
}

func stringField(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch x := v.(type) {
			case string:
				return x
			case fmt.Stringer:
				return x.String()
			default:
				return fmt.Sprint(x)
			}
		}
	}
	return ""
}

func intField(m map[string]interface{}, keys ...string) int {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch x := v.(type) {
			case float64:
				return int(x)
			case int:
				return x
			case int64:
				return int(x)
			case json.Number:
				i, err := x.Int64()
				if err == nil {
					return int(i)
				}
			case bool:
				if x {
					return 1
				}
				return 0
			}
		}
	}
	return 0
}
