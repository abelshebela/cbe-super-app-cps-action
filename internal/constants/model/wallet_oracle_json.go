package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// UnmarshalJSON accepts Oracle/CPS-style UPPERCASE keys and string "0"/"1" flags, as well as normal lowercase API JSON.
func (w *WalletOracle) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	w.ID = pickString(m, "id", "ID")
	w.Name = pickString(m, "name", "NAME")
	w.UniqueCode = pickString(m, "unique_code", "UNIQUE_CODE")
	w.Avatar = pickString(m, "avatar", "AVATAR")
	w.Enabled = pickBool(m, "enabled", "ENABLED")
	w.Self = pickBool(m, "self", "SELF", "services_self", "SERVICES_SELF")
	w.Other = pickBool(m, "other", "OTHER", "services_other", "SERVICES_OTHER")
	w.Agent = pickBool(m, "agent", "AGENT", "services_agent", "SERVICES_AGENT")
	w.IsDeleted = pickBool(m, "is_deleted", "IS_DELETED")

	if t, ok := pickTime(m, "created_at", "CREATED_AT"); ok {
		w.CreatedAt = t
	}
	if t, ok := pickTime(m, "last_modified_at", "LAST_MODIFIED_AT"); ok {
		w.LastModifiedAt = t
	}
	if raw, ok := pickRaw(m, "deleted_at", "DELETED_AT"); ok && raw != nil {
		switch v := raw.(type) {
		case string:
			if v == "" || strings.EqualFold(v, "null") {
				w.DeletedAt = nil
			} else if tt, err := time.Parse(time.RFC3339Nano, v); err == nil {
				w.DeletedAt = &tt
			} else if tt, err := time.Parse(time.RFC3339, v); err == nil {
				w.DeletedAt = &tt
			}
		case nil:
			w.DeletedAt = nil
		}
	}
	return nil
}

func pickString(m map[string]interface{}, keys ...string) string {
	v, _ := pickRaw(m, keys...)
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	default:
		return fmt.Sprint(x)
	}
}

func pickRaw(m map[string]interface{}, keys ...string) (interface{}, bool) {
	for _, k := range keys {
		for mk, mv := range m {
			if strings.EqualFold(mk, k) {
				return mv, true
			}
		}
	}
	return nil, false
}

func pickBool(m map[string]interface{}, keys ...string) bool {
	v, ok := pickRaw(m, keys...)
	if !ok || v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case float64:
		return x != 0
	case json.Number:
		i, err := x.Int64()
		return err == nil && i != 0
	case string:
		s := strings.TrimSpace(x)
		return s == "1" || strings.EqualFold(s, "true") || s == "yes"
	default:
		s := strings.TrimSpace(fmt.Sprint(x))
		return s == "1" || strings.EqualFold(s, "true")
	}
}

func pickTime(m map[string]interface{}, keys ...string) (time.Time, bool) {
	v, ok := pickRaw(m, keys...)
	if !ok || v == nil {
		return time.Time{}, false
	}
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprint(v)
	}
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05.000000000Z07:00"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
