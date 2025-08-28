package model

import (
	"time"
)

type Department struct {
	ID               string    `bson:"id,omitempty" json:"id,omitempty"`
	DepartmentCode   string    `bson:"department_code" json:"department_code"`
	Department       string    `bson:"department" json:"department"`
	PermissionGroups []string  `bson:"permission_groups" json:"permission_groups"`
	PortalCards      []string  `bson:"portal_cards" json:"portal_cards"`
	Enabled          bool      `bson:"enabled" json:"enabled"`
	IsDeleted        bool      `bson:"is_deleted" json:"is_deleted"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	LastModified     time.Time `bson:"last_modified" json:"last_modified"`
}
