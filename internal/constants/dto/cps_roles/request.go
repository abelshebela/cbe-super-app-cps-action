package cpsroles

type CreateCPSRoleRequest struct {
	Name        string `json:"name"`
	RoleCode    string `json:"role_code,omitempty"`
	Lable       string `json:"label" bson:"lable"`
	Description string `json:"description,omitempty"`
}

type UpdateCPSRoleRequest struct {
	Name        *string `json:"name,omitempty"`
	RoleCode    *string `json:"role_code,omitempty"`
	Lable       *string `json:"label" bson:"lable"`
	Description *string `json:"description,omitempty"`
}

type ToggleServiceAccessRequest struct {
	AccessListKeys []string `json:"access_list_keys"`
}
