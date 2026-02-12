package roles_dto

// CreateJobRoleRequest represents the payload to create a job role
type CreateJobRoleRequest struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	PortalCards []string `json:"portal_cards"`
}

// UpdateJobRoleRequest represents the payload to update a job role
// All fields are optional; at least one must be provided
type UpdateJobRoleRequest struct {
	Code        string   `json:"code,omitempty"`
	Name        string   `json:"name,omitempty"`
	PortalCards []string `json:"portal_cards,omitempty"`
}
