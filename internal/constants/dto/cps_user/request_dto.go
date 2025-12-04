package cpsuser

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved" example:"true"`
	Reason  *string `json:"reason,omitempty" example:"Approved after review"`
}

type CreateUserRequest struct {
	UserName           string          `json:"username" bson:"username,omitempty" example:"john.doe"`
	FullName           string          `json:"full_name" bson:"full_name,omitempty" example:"John Doe"`
	Department         bson.ObjectID   `json:"department" bson:"department,omitempty" example:"507f1f77bcf86cd799439011"`
	PhoneNumber        string          `json:"phone_number" bson:"phone_number,omitempty" example:"+251911234567"`
	Role               string          `json:"role" bson:"role,omitempty" example:"Maker"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty" example:"Male"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty" example:"john.doe@example.com"`
	PermissionCategory []bson.ObjectID `json:"permission_category" bson:"permission_category,omitempty" example:"[\"507f1f77bcf86cd799439011\"]"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups" bson:"permission_groups,omitempty" example:"[\"507f1f77bcf86cd799439011\"]"`
}

type UpdateUserRequest struct {
	UserCode           string          `json:"user_code,omitempty" bson:"user_code,omitempty" example:"USR001"`
	UserName           string          `json:"username,omitempty" bson:"username,omitempty" example:"john.doe.updated"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name,omitempty" example:"John Doe Updated"`
	Department         bson.ObjectID   `json:"department,omitempty" bson:"department,omitempty" example:"507f1f77bcf86cd799439011"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number,omitempty" example:"+251911234567"`
	Role               string          `json:"role,omitempty" bson:"role,omitempty" example:"Checker"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty" example:"Male"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty" example:"john.doe.updated@example.com"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category,omitempty" example:"[\"507f1f77bcf86cd799439011\"]"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups,omitempty" bson:"permission_groups,omitempty" example:"[\"507f1f77bcf86cd799439011\"]"`
}

type ApproveCPSAction struct {
	Approved bool    `json:"approved" example:"true"`
	Reason   *string `json:"reason" example:"Approved after thorough review"`
}

// CPSUserActionPayload is used for CPS action CurrentAction field
// It includes the user data plus populated permission structures for display
type CPSUserActionPayload struct {
	User                 interface{}                  `json:"user" bson:"user"`
	PermissionCategories []PermissionCategoryResponse `json:"permission_categories,omitempty" bson:"permission_categories,omitempty"`
	PermissionGroups     []PermissionGroupResponse    `json:"permission_groups,omitempty" bson:"permission_groups,omitempty"`
	PortalCards          []string                     `json:"portal_cards,omitempty" bson:"portal_cards,omitempty"`
}
