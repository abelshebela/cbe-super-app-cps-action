package department_dto

type CreateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}

type UpdateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}
