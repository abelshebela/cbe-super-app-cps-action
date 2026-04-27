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
		RAWTOHEX(SR.ID) AS ID,
		NVL(SR.NAME, '') AS NAME,
		NVL(SR.LABEL, '') AS LABEL,
		NVL(SR.ROLE_CODE, '') AS ROLE_CODE,
		NVL(SR.DESCRIPTION, '') AS DESCRIPTION,
		NVL(SR.ACCOUNT_TYPE, '') AS ACCOUNT_TYPE,
		SR.IS_ENABLED,
		SR.IS_DELETED,
		SR.CREATED_AT,
		SR.LAST_MODIFIED_AT,
		SR.DELETED_AT

	FROM SUPERAPP_ROLES SR
	WHERE SR.IS_DELETED = 0 AND %s
	ORDER BY SR.CREATED_AT DESC
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
