package donation_oracle

import (
	"database/sql"
	"time"
)

// boolToInt converts a Go bool to 0/1 for Oracle NUMBER(1) columns.
// godror rejects direct bool binds (ORA-00932).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullStringFromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullTimeFromTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}
