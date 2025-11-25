package core

import (
	"encoding/json"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/bankvault"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/pkgs/utils"

	"github.com/shopspring/decimal"
)

func ConvertBankVaultToMongoSafe(product *model.BankVaultProduct) map[string]interface{} {
	result := map[string]interface{}{
		"id":                             product.ID,
		"name":                           product.Name,
		"currency":                       product.Currency,
		"rate_bps":                       func() float64 { f, _ := product.RateBps.Float64(); return f }(),
		"method":                         string(product.Method),
		"frequency":                      string(product.Frequency),
		"lock_period":                    product.LockPeriod.Nanoseconds(),
		"min_amount":                     func() float64 { f, _ := product.MinAmount.Float64(); return f }(),
		"max_amount":                     func() float64 { f, _ := product.MaxAmount.Float64(); return f }(),
		"apply_interest_on_early_unlock": product.ApplyInterestOnEarlyUnlock,
		"is_active":                      product.IsActive,
		"created_at":                     product.CreatedAt,
		"updated_at":                     product.UpdatedAt,
		"deleted_at":                     product.DeletedAt,
		"created_by":                     product.CreatedBy,
		"updated_by":                     product.UpdatedBy,
		"is_deleted":                     product.IsDeleted,
	}
	return result
}

func MapBankVaultToResponse(bankVault *model.BankVaultProduct) *bankvault.BankVaultProductResponse {
	return &bankvault.BankVaultProductResponse{
		ID:                         bankVault.ID,
		Name:                       bankVault.Name,
		Currency:                   bankVault.Currency,
		Interest:                   bankVault.RateBps.Div(decimal.NewFromInt(100)),
		Method:                     bankVault.Method,
		Frequency:                  bankVault.Frequency,
		LockPeriod:                 fmt.Sprintf("%d months", utils.DurationToMonths(bankVault.LockPeriod)),
		MinAmount:                  bankVault.MinAmount,
		MaxAmount:                  bankVault.MaxAmount,
		ApplyInterestOnEarlyUnlock: bankVault.ApplyInterestOnEarlyUnlock,
		IsActive:                   bankVault.IsActive,
		IsDeleted:                  bankVault.IsDeleted,
		CreatedAt:                  bankVault.CreatedAt,
		UpdatedAt:                  bankVault.UpdatedAt,
		DeletedAt:                  bankVault.DeletedAt,
	}
}

func BuildUpdateBankVault(prev *model.BankVaultProduct, req *model.UpdateBankVault) model.BankVaultProduct {
	if req.MinAmount != nil {
		prev.MinAmount = decimal.NewFromFloat(*req.MinAmount)
	}
	if req.MaxAmount != nil {
		prev.MaxAmount = decimal.NewFromFloat(*req.MaxAmount)
	}
	return *prev
}

func BindBankVaultFromCPSAction(current interface{}) (*model.BankVaultProduct, error) {
	var actionMap map[string]interface{}

	marshaled, err := json.Marshal(current)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %v", err)
	}

	BV, err := MapBankVaultProduct(actionMap)
	if err != nil {
		return nil, err
	}

	return &BV, nil
}

func MapBankVaultProduct(data map[string]interface{}) (model.BankVaultProduct, error) {
	var BV model.BankVaultProduct

	BV.Name = getString(data, "name")
	BV.Currency = getString(data, "currency")
	BV.Method = constants.AccrualMethod(getString(data, "method"))
	BV.Frequency = constants.AccrualFrequency(getString(data, "frequency"))

	BV.RateBps = getDecimal(data, "rate_bps")
	BV.MinAmount = getDecimal(data, "min_amount")
	BV.MaxAmount = getDecimal(data, "max_amount")

	BV.LockPeriod = getDuration(data, "lock_period")

	BV.ApplyInterestOnEarlyUnlock = getBool(data, "apply_interest_on_early_unlock")
	BV.IsActive = getBool(data, "is_active")
	BV.IsDeleted = getBool(data, "is_deleted")

	BV.CreatedAt = getTime(data, "created_at")
	BV.UpdatedAt = getTime(data, "updated_at")
	BV.DeletedAt = getTimePtr(data, "deleted_at")

	BV.CreatedBy = getString(data, "created_by")
	BV.UpdatedBy = getString(data, "updated_by")

	return BV, nil
}

// Helper functions for type conversion
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

func getDecimal(data map[string]interface{}, key string) decimal.Decimal {
	val, exists := data[key]
	if !exists {
		return decimal.Zero
	}

	// Handle empty decimal object {} (old format)
	if valMap, ok := val.(map[string]interface{}); ok {
		if len(valMap) == 0 {
			return decimal.Zero
		}
		// If it's a decimal object with value, try to extract it
		if strVal, ok := valMap["value"].(string); ok {
			if d, err := decimal.NewFromString(strVal); err == nil {
				return d
			}
		}
	}

	// Try to parse as number directly (new format - float64)
	if floatVal, ok := val.(float64); ok {
		return decimal.NewFromFloat(floatVal)
	}
	if intVal, ok := val.(int64); ok {
		return decimal.NewFromInt(intVal)
	}
	if intVal, ok := val.(int); ok {
		return decimal.NewFromInt(int64(intVal))
	}

	// Try to parse as string
	if strVal, ok := val.(string); ok {
		if d, err := decimal.NewFromString(strVal); err == nil {
			return d
		}
	}

	return decimal.Zero
}

func getDuration(data map[string]interface{}, key string) time.Duration {
	if val, ok := data[key].(int64); ok {
		// Convert nanoseconds to duration
		return time.Duration(val)
	}
	if val, ok := data[key].(float64); ok {
		return time.Duration(int64(val))
	}
	return 0
}

func getTime(data map[string]interface{}, key string) time.Time {
	if val, ok := data[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t
		}
	}
	return time.Time{}
}

func getTimePtr(data map[string]interface{}, key string) *time.Time {
	if val, ok := data[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return &t
		}
	}
	return nil
}

func BindBankVaultUpdateFromCPSAction(current interface{}) (*model.BankVaultProduct, error) {
	var actionMap map[string]interface{}

	marshaled, err := json.Marshal(current)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %v", err)
	}

	BV, err := MapToUpdateBankVault(actionMap)
	if err != nil {
		return nil, err
	}

	return &BV, nil
}

func MapToUpdateBankVault(data map[string]interface{}) (model.BankVaultProduct, error) {
	var BV model.BankVaultProduct

	BV.MinAmount = getDecimal(data, "min_amount")
	BV.MaxAmount = getDecimal(data, "max_amount")

	return BV, nil
}
