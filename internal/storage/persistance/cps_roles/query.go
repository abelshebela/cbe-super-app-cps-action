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
		SR.NAME,
		SR.ROLE_CODE,
		SR.DESCRIPTION,
		SR.IS_ENABLED,
		SR.IS_DELETED,
		SR.CREATED_AT,
		SR.LAST_MODIFIED_AT,
		SR.DELETED_AT,

		-- ENABLED SERVICES
		(
			SELECT JSON_ARRAYAGG(
				JSON_OBJECT(
					'access_list_id' VALUE RAWTOHEX(AL.ID),
					'name' VALUE RAWTOHEX(AL.NAME),
					'service_code' VALUE AL.SERVICE_KEY
				)
			)
			FROM ACCESS_LISTS AL
			WHERE AL.IS_ENABLED = 1
			  AND AL.IS_DELETED = 0
			  AND NOT EXISTS (
				SELECT 1
				FROM ACCESS_LIST_CUSTOMER_SEG ACS
				WHERE ACS.ACCESS_LIST_KEY = AL.SERVICE_KEY
				  AND ACS.SEGMENTED_ID = SR.ID
				  AND ACS.ENABLED = 1
			  )
		) AS ENABLED_SERVICES,

		-- DISABLED SERVICES
		(
			SELECT JSON_ARRAYAGG(
				JSON_OBJECT(
					'access_list_id' VALUE RAWTOHEX(AL.ID),
					'name' VALUE RAWTOHEX(AL.NAME),
					'key' VALUE AL.SERVICE_KEY
				)
			)
			FROM ACCESS_LISTS AL
			WHERE AL.IS_ENABLED = 1
			  AND AL.IS_DELETED = 0
			  AND EXISTS (
				SELECT 1
				FROM ACCESS_LIST_CUSTOMER_SEG ACS
				WHERE ACS.ACCESS_LIST_KEY = AL.SERVICE_KEY
				  AND ACS.SEGMENTED_ID = SR.ID
				  AND ACS.ENABLED = 1
			  )
		) AS DISABLED_SERVICES

	FROM SUPERAPP_ROLE SR
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
