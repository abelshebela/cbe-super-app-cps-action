package permission

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
