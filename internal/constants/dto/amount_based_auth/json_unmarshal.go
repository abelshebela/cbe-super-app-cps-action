package amount_based_auth

import (
	"encoding/json"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
)

// UnmarshalJSON accepts min_amount / MIN_AMOUNT and max_amount / MAX_AMOUNT (CPS/Oracle style).
func (r *UpdateAmountBasedAuthRequest) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	r.MinAmount = pickUint64Keys(m, "min_amount", "MIN_AMOUNT", "MinAmount")
	r.MaxAmount = pickUint64Keys(m, "max_amount", "MAX_AMOUNT", "MaxAmount")
	return nil
}

// UnmarshalJSON accepts method/METHOD, min/max in several casings and numeric forms.
func (t *TierInput) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	t.Method = CanonicalMethodFromPath(pickStringKeys(m, "method", "METHOD"))
	t.MinAmount = int64(pickUint64Keys(m, "min_amount", "MIN_AMOUNT", "MinAmount"))
	t.MaxAmount = int64(pickUint64Keys(m, "max_amount", "MAX_AMOUNT", "MaxAmount"))
	return nil
}

func pickStringKeys(m map[string]interface{}, keys ...string) string {
	v, ok := pickAnyKey(m, keys...)
	if !ok || v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x)
	default:
		return strings.TrimSpace(fmt.Sprint(x))
	}
}

func pickUint64Keys(m map[string]interface{}, keys ...string) uint64 {
	v, ok := pickAnyKey(m, keys...)
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

func pickAnyKey(m map[string]interface{}, keys ...string) (interface{}, bool) {
	for _, want := range keys {
		for k, v := range m {
			if strings.EqualFold(k, want) {
				return v, true
			}
		}
	}
	return nil, false
}

// UnmarshalJSON supports currency/CURRENCY, methods/METHODS, and tiers with flexible tier fields.
func (r *AddCurrencyRequest) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	r.Currency = constants.CurrencyType(pickStringKeys(m, "currency", "CURRENCY"))
	if rawMethods, ok := pickAnyKey(m, "methods", "METHODS"); ok && rawMethods != nil {
		mb, err := json.Marshal(rawMethods)
		if err != nil {
			return err
		}
		var methods []constants.Method
		if err := json.Unmarshal(mb, &methods); err != nil {
			return err
		}
		for i := range methods {
			methods[i] = CanonicalMethodFromPath(string(methods[i]))
		}
		r.Methods = methods
	}
	rawTiers, ok := pickAnyKey(m, "tiers", "TIERS")
	if !ok || rawTiers == nil {
		return nil
	}
	tb, err := json.Marshal(rawTiers)
	if err != nil {
		return err
	}
	var tiers []TierInput
	if err := json.Unmarshal(tb, &tiers); err != nil {
		return err
	}
	r.Tiers = tiers
	return nil
}

// UnmarshalJSON supports methods/METHODS and tiers/TIERS.
func (r *ResetConfigRequest) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	rawMethods, ok := pickAnyKey(m, "methods", "METHODS")
	if ok && rawMethods != nil {
		mb, err := json.Marshal(rawMethods)
		if err != nil {
			return err
		}
		var methods []constants.Method
		if err := json.Unmarshal(mb, &methods); err != nil {
			return err
		}
		for i := range methods {
			methods[i] = CanonicalMethodFromPath(string(methods[i]))
		}
		r.Methods = methods
	}
	rawTiers, ok := pickAnyKey(m, "tiers", "TIERS")
	if ok && rawTiers != nil {
		tb, err := json.Marshal(rawTiers)
		if err != nil {
			return err
		}
		var tiers []TierInput
		if err := json.Unmarshal(tb, &tiers); err != nil {
			return err
		}
		r.Tiers = tiers
	}
	return nil
}
