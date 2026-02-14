package cpsroles

type CreateCPSRoleRequest struct {
	Name        string `json:"name"`
	RoleCode    string `json:"role_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateCPSRoleRequest struct {
	Name        string `json:"name,omitempty"`
	RoleCode    string `json:"role_code,omitempty"`
	Description string `json:"description,omitempty"`
}
