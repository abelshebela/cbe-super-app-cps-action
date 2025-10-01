package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/bankvault"
	"cbe-super-app-cps-action/internal/constants/model"

	"github.com/shopspring/decimal"
)

// ConvertBankVaultToMongoSafe converts BankVaultProduct to a MongoDB-safe format
// by converting decimal.Decimal fields to float64 for proper serialization
func ConvertBankVaultToMongoSafe(product *model.BankVaultProduct) map[string]interface{} {
	result := map[string]interface{}{
		"id":                product.ID,
		"name":              product.Name,
		"description":       product.Description,
		"currency":          product.Currency,
		"ratebps":           func() float64 { f, _ := product.RateBps.Float64(); return f }(),
		"method":            string(product.Method),
		"frequency":         string(product.Frequency),
		"lockperiod":        product.LockPeriod.Nanoseconds(),
		"minamount":         func() float64 { f, _ := product.MinAmount.Float64(); return f }(),
		"maxamount":         func() float64 { f, _ := product.MaxAmount.Float64(); return f }(),
		"earlyunlockfeebps": func() float64 { f, _ := product.EarlyUnlockFeeBps.Float64(); return f }(),
		"isactive":          product.IsActive,
		"createdat":         product.CreatedAt,
		"updatedat":         product.UpdatedAt,
		"deletedat":         product.DeletedAt,
		"createdby":         product.CreatedBy,
		"updatedby":         product.UpdatedBy,
		"isdeleted":         product.IsDeleted,
	}
	return result
}

func MapBankVaultToResponse(bankVault *model.BankVaultProduct) *bankvault.BankVaultProductResponse {
	return &bankvault.BankVaultProductResponse{
		ID:                bankVault.ID,
		Name:              bankVault.Name,
		Description:       bankVault.Description,
		Currency:          bankVault.Currency,
		RateBps:           bankVault.RateBps,
		Method:            bankVault.Method,
		Frequency:         bankVault.Frequency,
		LockPeriod:        bankVault.LockPeriod,
		MinAmount:         bankVault.MinAmount,
		MaxAmount:         bankVault.MaxAmount,
		EarlyUnlockFeeBps: bankVault.EarlyUnlockFeeBps,
		IsActive:          bankVault.IsActive,
		IsDeleted:         bankVault.IsDeleted,
		CreatedAt:         bankVault.CreatedAt,
		UpdatedAt:         bankVault.UpdatedAt,
		DeletedAt:         bankVault.DeletedAt,
	}
}

func BuildUpdateBankVault(prev *model.BankVaultProduct, req *model.UpdateBankVault) model.BankVaultProduct {
	if req.Description != nil {
		prev.Description = *req.Description
	}
	if req.MinAmount != nil {
		prev.MinAmount = decimal.NewFromFloat(*req.MinAmount)
	}
	if req.MaxAmount != nil {
		prev.MaxAmount = decimal.NewFromFloat(*req.MaxAmount)
	}
	return *prev
}

// bind
func BindBankVaultFromCPSAction(current interface{}) (model.BankVaultProduct, error) {
	var BV model.BankVaultProduct

	if v, ok := current.(model.BankVaultProduct); ok {
		return v, nil
	}

	// If stored as JSON string
	if s, ok := current.(string); ok {
		if err := json.Unmarshal([]byte(s), &BV); err == nil {
			return BV, nil
		}
	}

	// Generic path: marshal then unmarshal
	bytes, err := json.Marshal(current)
	if err != nil {
		return BV, err
	}

	// Use custom mapping function to handle camelCase -> snake_case conversion
	return MapCamelCaseToBankVaultProduct(bytes)
}

// MapCamelCaseToBankVaultProduct converts MongoDB camelCase JSON to BankVaultProduct model
func MapCamelCaseToBankVaultProduct(jsonBytes []byte) (model.BankVaultProduct, error) {
	var BV model.BankVaultProduct

	// First, unmarshal into a generic map to handle field name conversion
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return BV, err
	}

	// Map camelCase fields to snake_case fields manually
	BV.ID = getString(data, "id")
	BV.Name = getString(data, "name")
	BV.Description = getString(data, "description")
	BV.Currency = getString(data, "currency")
	BV.Method = constants.AccrualMethod(getString(data, "method"))
	BV.Frequency = constants.AccrualFrequency(getString(data, "frequency"))
	BV.IsActive = getBool(data, "isactive")
	BV.CreatedBy = getString(data, "createdby")
	BV.UpdatedBy = getString(data, "updatedby")
	BV.IsDeleted = getBool(data, "isdeleted")

	// Handle decimal fields
	BV.RateBps = getDecimal(data, "ratebps")
	BV.MinAmount = getDecimal(data, "minamount")
	BV.MaxAmount = getDecimal(data, "maxamount")
	BV.EarlyUnlockFeeBps = getDecimal(data, "earlyunlockfeebps")

	// Handle time.Duration field
	BV.LockPeriod = getDuration(data, "lockperiod")

	// Handle time fields
	BV.CreatedAt = getTime(data, "createdat")
	BV.UpdatedAt = getTime(data, "updatedat")
	BV.DeletedAt = getTimePtr(data, "deletedat")

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

func BindBankVaultUpdateFromCPSAction(current interface{}) (model.UpdateBankVault, error) {
	var BV model.UpdateBankVault
	if v, ok := current.(model.UpdateBankVault); ok {
		return v, nil
	}

	// If stored as JSON string
	if s, ok := current.(string); ok {
		if err := json.Unmarshal([]byte(s), &BV); err == nil {
			return BV, nil
		}
	}

	// Generic path: marshal then unmarshal
	bytes, err := json.Marshal(current)
	if err != nil {
		return BV, err
	}

	// Use custom mapping function to handle camelCase -> UpdateBankVault conversion
	return MapCamelCaseToUpdateBankVault(bytes)
}

// MapCamelCaseToUpdateBankVault converts MongoDB camelCase JSON to UpdateBankVault model
func MapCamelCaseToUpdateBankVault(jsonBytes []byte) (model.UpdateBankVault, error) {
	var result model.UpdateBankVault

	// First, unmarshal into a generic map to handle field name conversion
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return result, err
	}

	// Map camelCase fields to UpdateBankVault fields
	if desc, ok := data["description"].(string); ok && desc != "" {
		result.Description = &desc
	}

	// Handle min_amount
	if minAmount, ok := data["minamount"]; ok {
		if floatVal, ok := minAmount.(float64); ok && floatVal > 0 {
			result.MinAmount = &floatVal
		}
	}

	// Handle max_amount
	if maxAmount, ok := data["maxamount"]; ok {
		if floatVal, ok := maxAmount.(float64); ok && floatVal > 0 {
			result.MaxAmount = &floatVal
		}
	}

	return result, nil
}

func BankUpdateVault(req *model.UpdateBankVault) model.BankVaultProduct {
	result := model.BankVaultProduct{}

	if req.Description != nil {
		result.Description = *req.Description
	}
	if req.MinAmount != nil {
		result.MinAmount = decimal.NewFromFloat(*req.MinAmount)
	}
	if req.MaxAmount != nil {
		result.MaxAmount = decimal.NewFromFloat(*req.MaxAmount)
	}

	return result
}

func MapBankVaultToCreate(bankVault bankvault.CreateBankVaultProductRequest) (*model.BankVaultProduct, error) {
	// Convert string days → int
	days, err := strconv.Atoi(bankVault.LockPeriodDays)
	if err != nil {
		return nil, fmt.Errorf("invalid lock period days: %w", err)
	}

	// Convert days → time.Duration
	lockPeriod := time.Duration(days) * 24 * time.Hour

	return &model.BankVaultProduct{
		Name:              bankVault.Name,
		Description:       bankVault.Description,
		Currency:          bankVault.Currency,
		RateBps:           bankVault.RateBps,
		Method:            bankvault.ToDomainMethod(bankVault.Method),
		Frequency:         bankvault.ToDomainFrequency(bankVault.Frequency),
		LockPeriod:        lockPeriod,
		MinAmount:         bankVault.MinAmount,
		MaxAmount:         bankVault.MaxAmount,
		EarlyUnlockFeeBps: bankVault.EarlyUnlockFeeBps,
		IsActive:          false,
		IsDeleted:         false,
	}, nil
}
