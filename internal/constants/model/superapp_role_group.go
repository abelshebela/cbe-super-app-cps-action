package model

import "time"

type SuperAppRoleGroup struct {
	SuperappRole      string    `json:"superapp_role"`
	SuperappRoleLabel string    `json:"superapp_role_label"`
	IsEnabled         bool      `json:"is_enabled"`
	Segments          []Segment `json:"segments"`
	CreatedAt         time.Time `json:"created_at"`
	LastModifiedAt    time.Time `json:"last_modified_at"`
}
