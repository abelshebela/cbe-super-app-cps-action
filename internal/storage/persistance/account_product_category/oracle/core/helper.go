package account_product_category_oracle_core

import "database/sql"

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func NullStringFromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
