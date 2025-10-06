package permission

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreatePermissionGroupRequest struct {
	GroupName               string   `json:"group_name"`
	Role                    string   `json:"role"`
	PermissionCategoryLists []string `json:"permission_category_list"`
}

type UpdatePermissionGroupRequest struct {
	OldGroupName            string   `json:"old_group_name"`
	NewGroupName            string   `json:"new_group_name"`
	Role                    string   `json:"role"`
	PermissionCategoryLists []string `json:"permission_category_list"`
}

type PaginatedPermissionGroupResponse struct {
	Data []PermissionCategoryResponse `json:"docs"`
	Meta types.PaginationMeta         `json:"meta"`
}

type PermissionCategoryResponse struct {
	ID           bson.ObjectID
	CategoryName string
	Access       string
	Permissions  interface{}
	Enabled      bool
	IsDeleted    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
