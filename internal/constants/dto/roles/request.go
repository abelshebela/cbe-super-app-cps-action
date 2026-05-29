package roles_dto

// CreateJobRoleRequest represents the payload to create a job role
type CreateJobRoleRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// UpdateJobRoleRequest represents the payload to update a job role
// All fields are optional; at least one must be provided
type UpdateJobRoleRequest struct {
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}
