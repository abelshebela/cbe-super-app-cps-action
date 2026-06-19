package bank_core

import (
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"fmt"

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
	fmt.Printf("[Bank_oracle_mapper] Mapping action to BankOracle: %v\n", action)
	bank := imodel.BankOracle{}

	if v, ok := action["bank_name"]; ok {
		if name, ok := v.(string); ok {
			bank.BankName = name
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
		if enabled, ok := v.(bool); ok && enabled {
			bank.IsEnabled = 1
		}
	}
	if v, ok := action["has_alpha_numeric"]; ok {
		if hasAlpha, ok := v.(bool); ok && hasAlpha {
			bank.HasAlphaNumeric = 1
		}
	}
	if v, ok := action["account_length"]; ok {
		if accountLength, ok := v.(float64); ok {
			bank.AccountLength = int(accountLength)
		} else if accountLength, ok := v.(int32); ok {
			bank.AccountLength = int(accountLength)
		} else if accountLength, ok := v.(int64); ok {
			bank.AccountLength = int(accountLength)
		}
	}
	if v, ok := action["is_cbe"]; ok {
		if isCBE, ok := v.(bool); ok && isCBE {
			bank.IS_CBE = 1
		}
	}

	return bank
}

// BankOracleToPayload converts a BankOracle (int-based) to BankOracleResponse (bool-based)
// for storage in the CPS action so the checker sees true/false instead of 0/1.
func BankOracleToPayload(b imodel.BankOracle) bank_dto.BankOracleResponse {
	return bank_dto.BankOracleResponse{
		ID:              b.ID,
		BankName:        b.BankName,
		Logo:            b.Logo,
		BICCode:         b.BICCode,
		IsEnabled:       b.IsEnabled == 1,
		IsCBE:           b.IS_CBE == 1,
		AccountLength:   b.AccountLength,
		HasAlphaNumeric: b.HasAlphaNumeric == 1,
		CreateAt:        b.CreateAt,
		UpdateAt:        b.UpdateAt,
	}
}
