package cpsuser

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved" example:"true"`
	Reason  *string `json:"reason,omitempty" example:"Approved after review"`
}

type CreateUserRequest struct {
	UserName    string `json:"username" bson:"username,omitempty" example:"john.doe"`
	FullName    string `json:"full_name" bson:"full_name,omitempty" example:"John Doe"`
	PhoneNumber string `json:"phone_number" bson:"phone_number,omitempty" example:"+251911234567"`
	JobTitle    string `json:"job_title" bson:"job_title"`
	Gender      string `json:"gender,omitempty" bson:"gender,omitempty" example:"Male"`
	Email       string `json:"email,omitempty" bson:"email,omitempty" example:"john.doe@example.com"`
}

type UpdateUserRequest struct {
	UserName    string `json:"username,omitempty" bson:"username,omitempty" example:"john.doe.updated"`
	FullName    string `json:"full_name,omitempty" bson:"full_name,omitempty" example:"John Doe Updated"`
	PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number,omitempty" example:"+251911234567"`
	Gender      string `json:"gender,omitempty" bson:"gender,omitempty" example:"Male"`
	Email       string `json:"email,omitempty" bson:"email,omitempty" example:"john.doe.updated@example.com"`
	JobTitle    string `json:"job_title,omitempty" bson:"job_title,omitempty"`
}

type ApproveCPSAction struct {
	Approved bool    `json:"approved" example:"true"`
	Reason   *string `json:"reason" example:"Approved after thorough review"`
}

type CPSUserActionPayload struct {
	User                 interface{}                  `json:"user" bson:"user"`
	PermissionCategories []PermissionCategoryResponse `json:"permission_categories,omitempty" bson:"permission_categories,omitempty"`
	PermissionGroups     []PermissionGroupResponse    `json:"permission_groups,omitempty" bson:"permission_groups,omitempty"`
	PortalCards          []string                     `json:"portal_cards,omitempty" bson:"portal_cards,omitempty"`
}
