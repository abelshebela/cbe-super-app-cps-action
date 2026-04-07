package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
)

// UnmarshalJSON supports API JSON (min_amount) and CPS/Oracle-style keys (MIN_AMOUNT, MAX_AMOUNT).
func (a *AuthTierOracle) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	a.ID = pickString(m, "id", "ID")
	a.Currency = constants.CurrencyType(pickString(m, "currency", "CURRENCY"))
	a.MinAmount = pickUint64Flex(m, "min_amount", "MIN_AMOUNT")
	a.MaxAmount = pickUint64Flex(m, "max_amount", "MAX_AMOUNT")
	a.Method = constants.Method(pickString(m, "method", "METHOD"))
	a.Enabled = pickInt01(m, "enabled", "ENABLED")
	a.IsDeleted = pickInt01(m, "is_deleted", "IS_DELETED")
	if t, ok := pickTime(m, "created_at", "CREATED_AT"); ok {
		a.CreatedAt = t
	}
	if t, ok := pickTime(m, "last_modified", "LAST_MODIFIED", "last_modified_at", "LAST_MODIFIED_AT"); ok {
		a.LastModified = t
	}
	return nil
}

func pickUint64Flex(m map[string]interface{}, keys ...string) uint64 {
	v, ok := pickRaw(m, keys...)
	if !ok || v == nil {
		return 0
	}
	switch x := v.(type) {
	case float64:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case json.Number:
		n, err := x.Int64()
		if err != nil || n < 0 {
			return 0
		}
		return uint64(n)
	case int:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int64:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return 0
		}
		var n uint64
		_, _ = fmt.Sscanf(s, "%d", &n)
		return n
	default:
		var n uint64
		_, _ = fmt.Sscanf(fmt.Sprint(x), "%d", &n)
		return n
	}
}

func pickInt01(m map[string]interface{}, keys ...string) int {
	v, ok := pickRaw(m, keys...)
	if !ok || v == nil {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return int(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case json.Number:
		i, _ := x.Int64()
		return int(i)
	case int:
		return x
	case string:
		s := strings.TrimSpace(strings.ToLower(x))
		if s == "1" || s == "true" {
			return 1
		}
		return 0
	default:
		return 0
	}
}
