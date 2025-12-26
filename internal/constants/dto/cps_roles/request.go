package cpsroles

type CreateCPSRoleRequest struct {
	Name string `json:"name"`
}

type UpdateCPSRoleRequest struct {
	Name *string `json:"name,omitempty"`
}
