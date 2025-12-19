package core

import (
	"encoding/json"
	"strconv"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func AmountTierMapper(action map[string]interface{}) model.VaultAmountTier {
	amountTier := model.VaultAmountTier{}

	if v, ok := action["vault_category_id"]; ok {
		if id, ok := v.(string); ok {
			amountTier.VaultCategoryID = id
		}
	}
	if v, ok := action["min_amount"]; ok {
		amountTier.MinAmount = getFloat64(v)
	}
	if v, ok := action["max_amount"]; ok {
		amountTier.MaxAmount = getFloat64(v)
	}
	if v, ok := action["interest"]; ok {
		amountTier.Interest = getFloat64(v)
	}
	if v, ok := action["is_active"]; ok {
		if isActive, ok := v.(bool); ok {
			amountTier.IsActive = isActive
		}
	}

	return amountTier
}

func BindVaultAmountTierFromCPSAction(current interface{}) (model.VaultAmountTier, error) {
	var vaultAmountTier model.VaultAmountTier
	if v, ok := current.(model.VaultAmountTier); ok {
		return v, nil
	}
	if s, ok := current.(string); ok {
		if err := json.Unmarshal([]byte(s), &vaultAmountTier); err == nil {
			return vaultAmountTier, nil
		}
	}
	bytes, err := json.Marshal(current)
	if err != nil {
		return vaultAmountTier, err
	}
	return MapCamelCaseToVaultAmountTier(bytes)
}

func MapCamelCaseToVaultAmountTier(jsonBytes []byte) (model.VaultAmountTier, error) {
	var result model.VaultAmountTier

	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return result, err
	}

	result.ID = getString(data, "id")
	result.VaultCategoryID = getString(data, "vault_category_id")
	result.MinAmount = getFloat64(data["min_amount"])
	result.MaxAmount = getFloat64(data["max_amount"])
	result.Interest = getFloat64(data["interest"])
	result.IsActive = getBool(data, "is_active")
	result.IsDeleted = getBool(data, "is_deleted")
	result.CreatedAt = getTime(data, "created_at")
	result.UpdatedAt = getTime(data, "updated_at")
	result.DeletedAt = getTimePtr(data, "deleted_at")

	return result, nil
}

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

func getFloat64(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	if i, ok := v.(int); ok {
		return float64(i)
	}
	if i, ok := v.(int64); ok {
		return float64(i)
	}
	if s, ok := v.(string); ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return 0
}
