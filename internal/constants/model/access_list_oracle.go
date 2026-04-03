package model

import "time"

// AccessListOracle maps ACCESS_LIST rows in Oracle.
// Boolean flags are represented as NUMBER(1) in DB and int in Go.
type AccessListOracle struct {
	ID string ` json:"id,omitempty"`

	Key            string ` json:"key"`
	Enabled        int    ` json:"enabled"`
	AccessListName string ` json:"access_list_name"`
	USSDEnabled    int    ` json:"ussd_enabled"`

	IsDeleted int ` json:"is_deleted"`

	CreateAt time.Time ` json:"create_at"`
	UpdateAt time.Time ` json:"update_at"`
}
