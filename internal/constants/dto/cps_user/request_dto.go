package cpsuser

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved"`
	Reason  *string `json:"reason,omitempty"`
}

type CreateUserRequest struct {
	UserName           string          `json:"username" bson:"username,omitempty"`
	FullName           string          `json:"full_name" bson:"full_name,omitempty"`
	Department         bson.ObjectID   `json:"department" bson:"department,omitempty"`
	PhoneNumber        string          `json:"phone_number" bson:"phone_number,omitempty"`
	Role               string          `json:"role" bson:"role,omitempty"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category" bson:"permission_category,omitempty"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups" bson:"permission_groups,omitempty"`
}

type UpdateUserRequest struct {
	UserCode           string          `json:"user_code,omitempty" bson:"user_code,omitempty"`
	UserName           string          `json:"username,omitempty" bson:"username,omitempty"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Department         bson.ObjectID   `json:"department,omitempty" bson:"department,omitempty"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Role               string          `json:"role,omitempty" bson:"role,omitempty"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category,omitempty"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups,omitempty" bson:"permission_groups,omitempty"`
}

// Normalize normalizes the request data

type ApproveCPSAction struct {
	Approved bool    `json:"approved"`
	Reason   *string `json:"reason"`
}
