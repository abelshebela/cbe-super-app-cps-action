package cpsuser

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved" example:"true"`
	Reason  *string `json:"reason,omitempty" example:"Approved after review"`
}

type CreateUserRequest struct {
	UserName           string          `json:"username" example:"john.doe"`
	FullName           string          `json:"full_name" example:"John Doe"`
	Department         bson.ObjectID   `json:"department" example:"507f1f77bcf86cd799439011"`
	PhoneNumber        string          `json:"phone_number" example:"+251911234567"`
	Role               string          `json:"role" example:"Maker"`
	Gender             string          `json:"gender,omitempty" example:"Male"`
	Email              string          `json:"email,omitempty" example:"john.doe@example.com"`
	PermissionCategory []bson.ObjectID `json:"permission_category" example:"[\"507f1f77bcf86cd799439011\"]"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups" example:"[\"507f1f77bcf86cd799439011\"]"`
}

type UpdateUserRequest struct {
	UserCode           string          `json:"user_code,omitempty" example:"USR001"`
	UserName           string          `json:"username,omitempty" example:"john.doe.updated"`
	FullName           string          `json:"full_name,omitempty" example:"John Doe Updated"`
	Department         bson.ObjectID   `json:"department,omitempty" example:"507f1f77bcf86cd799439011"`
	PhoneNumber        string          `json:"phone_number,omitempty" example:"+251911234567"`
	Role               string          `json:"role,omitempty" example:"Checker"`
	Gender             string          `json:"gender,omitempty" example:"Male"`
	Email              string          `json:"email,omitempty" example:"john.doe.updated@example.com"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" example:"[\"507f1f77bcf86cd799439011\"]"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups,omitempty" example:"[\"507f1f77bcf86cd799439011\"]"`
}


type ApproveCPSAction struct {
	Approved bool    `json:"approved" example:"true"`
	Reason   *string `json:"reason" example:"Approved after thorough review"`
}
