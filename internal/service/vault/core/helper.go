package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func CategoryMapper(action map[string]interface{}) imodel.VaultCategory {
	category := imodel.VaultCategory{}

	if v, ok := action["name"]; ok {
		if name, ok := v.(string); ok {
			category.Name = name
		}
	}

	if v, ok := action["cover_image_url"]; ok {
		if img, ok := v.(string); ok {
			category.CoverImageURL = img
		}
	}

	if v, ok := action["interest_type"]; ok {
		if it, ok := v.(string); ok {
			category.InterestType = it
		}
	}

	if v, ok := action["category_interest"]; ok {
		category.CategoryInterest = getAnyAsString(v)
	}

	if v, ok := action["deadlock"]; ok {
		if d, ok := v.(bool); ok {
			category.Deadlock = d
		}
	}

	if v, ok := action["is_active"]; ok {
		if ia, ok := v.(bool); ok {
			category.IsActive = ia
		}
	}

	if v, ok := action["tiers"]; ok {
		switch val := v.(type) {
		case []interface{}:
			category.Tiers = TiersMapper(val)
		case string:
			var raw []interface{}
			if err := json.Unmarshal([]byte(val), &raw); err == nil {
				category.Tiers = TiersMapper(raw)
			}
		default:
			if bytes, err := json.Marshal(val); err == nil {
				var raw []interface{}
				if err := json.Unmarshal(bytes, &raw); err == nil {
					category.Tiers = TiersMapper(raw)
				}
			}
		}
	}

	return category
}

func TiersMapper(tiersRaw []interface{}) []imodel.VaultTiers {
	tiers := make([]imodel.VaultTiers, 0, len(tiersRaw))
	for _, tr := range tiersRaw {
		if tMap, ok := tr.(map[string]interface{}); ok {
			tier := imodel.VaultTiers{}
			if v, ok := tMap["id"]; ok {
				tier.ID = getAnyAsString(v)
			} else if v, ok := tMap["_id"]; ok {
				tier.ID = getAnyAsString(v)
			}
			if v, ok := tMap["name"]; ok {
				if name, ok := v.(string); ok {
					tier.Name = name
				}
			}
			if v, ok := tMap["tier_interest"]; ok {
				tier.TierInterest = getAnyAsString(v)
			}
			if v, ok := tMap["min"]; ok {
				tier.MinAmount = getAnyAsString(v)
			}
			if v, ok := tMap["max"]; ok {
				tier.MaxAmount = getAnyAsString(v)
			}
			tiers = append(tiers, tier)
		}
	}
	return tiers
}

func getAnyAsString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(val, 10)
	case int:
		return strconv.Itoa(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func MapVaultCategoryToResponse(vaultCategory *imodel.VaultCategory) *imodel.VaultCategory {
	return &imodel.VaultCategory{
		ID:               vaultCategory.ID,
		Name:             vaultCategory.Name,
		CoverImageURL:    vaultCategory.CoverImageURL,
		InterestType:     vaultCategory.InterestType,
		CategoryInterest: vaultCategory.CategoryInterest,
		Deadlock:         vaultCategory.Deadlock,
		IsActive:         vaultCategory.IsActive,
		CreatedAt:        vaultCategory.CreatedAt,
		UpdatedAt:        vaultCategory.UpdatedAt,
	}
}

func BindVaultCategoryFromCPSAction(current interface{}) (imodel.VaultCategory, error) {
	var category imodel.VaultCategory
	if v, ok := current.(imodel.VaultCategory); ok {
		return v, nil
	}
	if s, ok := current.(string); ok {
		if err := json.Unmarshal([]byte(s), &category); err == nil {
			return category, nil
		}
	}
	bytes, err := json.Marshal(current)
	if err != nil {
		return category, err
	}
	return MapCamelCaseToVaultCategory(bytes)
}

func MapCamelCaseToVaultCategory(jsonBytes []byte) (imodel.VaultCategory, error) {
	var result imodel.VaultCategory

	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return result, err
	}

	result.ID = getString(data, "id")
	result.Name = getString(data, "name")
	result.CoverImageURL = getString(data, "cover_image_url")
	result.InterestType = getString(data, "interest_type")
	result.CategoryInterest = getAnyAsString(data["category_interest"])
	result.Deadlock = getBool(data, "deadlock")
	result.IsActive = getBool(data, "is_active")
	result.CreatedAt = getTime(data, "created_at")
	result.UpdatedAt = getTime(data, "updated_at")

	return result, nil
}

func BindVaultCategoryUpdateFromCPSAction(current interface{}) (imodel.VaultCategory, error) {
	return BindVaultCategoryFromCPSAction(current)
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

// WithdrawalMapper maps action data (e.g. from CPS CurrentAction) to imodel.Withdrawal.
func WithdrawalMapper(action map[string]interface{}) imodel.Withdrawal {
	w := imodel.Withdrawal{}
	if v, ok := action["id"]; ok {
		w.ID = getAnyAsString(v)
	}
	if v, ok := action["locked_vault_id"]; ok {
		w.LockedVaultID = getAnyAsString(v)
	}
	if v, ok := action["amount"]; ok {
		w.Amount = getAnyAsString(v)
	}
	if v, ok := action["withdrawer_name"]; ok {
		if s, ok := v.(string); ok {
			w.WithdrawerName = s
		}
	}
	if v, ok := action["withdrawer_phone_number"]; ok {
		if s, ok := v.(string); ok {
			w.WithdrawerPhoneNumber = s
		}
	}
	if v, ok := action["status"]; ok {
		if s, ok := v.(string); ok {
			w.Status = s
		}
	}
	if v, ok := action["is_active"]; ok {
		if b, ok := v.(bool); ok {
			w.IsActive = b
		}
	}
	if v, ok := action["created_at"]; ok {
		w.CreatedAt = getTimeFromAny(v)
	}
	if v, ok := action["updated_at"]; ok {
		w.UpdatedAt = getTimeFromAny(v)
	}
	return w
}

func WithdrawalStatusMapper(action map[string]interface{}) string {
	if v, ok := action["withdrawal_status"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getTimeFromAny(v interface{}) time.Time {
	if t, ok := v.(time.Time); ok {
		return t
	}
	if s, ok := v.(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
