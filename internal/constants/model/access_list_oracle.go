package model

import (
	"database/sql"
	"encoding/hex"
	"time"
)

// AccessListOracle maps ACCESS_LIST rows in Oracle.
//
//	CREATE TABLE ACCESS_LIST (
//	    ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
//	    NAME VARCHAR(64) NOT NULL,
//	    SERVICE_KEY VARCHAR2(32) NOT NULL,
//	    IS_ENABLED NUMBER(1) DEFAULT 1,
//	    IS_DELETED NUMBER(1) DEFAULT 0,
//	    CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//	    LAST_MODIFIED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//	    DELETED_AT TIMESTAMP
//	);
type AccessListOracle struct {
	ID             []byte       `json:"-"` // RAW(16)
	IDHex          string       `json:"id,omitempty"`
	Name           string       `json:"name"`
	ServiceKey     string       `json:"service_key"`
	IsEnabled      int          `json:"is_enabled"`
	IsDeleted      int          `json:"is_deleted"`
	CreatedAt      time.Time    `json:"created_at"`
	LastModifiedAt time.Time    `json:"last_modified_at"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty"`
}

// SetIDHex encodes RAW(16) ID for JSON / logging.
func (a *AccessListOracle) SetIDHex() {
	if len(a.ID) == 0 {
		a.IDHex = ""
		return
	}
	a.IDHex = hex.EncodeToString(a.ID)
}
