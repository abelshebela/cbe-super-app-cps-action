package cpsroles

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func FetchCPSRolesQuery(where string) string {
	if where == "" {
		where = "1=1"
	}

	return fmt.Sprintf(`
	SELECT
		RAWTOHEX(ID) AS ID,
		NAME,
		ROLE_CODE,
		DESCRIPTION,
		IS_ENABLED,
		IS_DELETED,
		CREATED_AT,
		LAST_MODIFIED_AT,
		DELETED_AT,
	FROM SUPERAPP_ROLE
	WHERE IS_DELETED = 0 AND %s
	ORDER BY CREATED_AT DESC
	OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
	`, where)
}

func SetTime(dst *time.Time, src sql.NullTime) {
	if src.Valid {
		*dst = src.Time
	}
}

func UnmarshalJSON(src sql.NullString, dst any) {
	if !src.Valid || src.String == "" || src.String == "null" {
		return
	}

	_ = json.Unmarshal([]byte(src.String), dst)
}

func PtrBool(b bool) *bool {
	return &b
}
