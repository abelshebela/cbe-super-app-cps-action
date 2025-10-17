package cpsuser

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUserDTO struct {
	ID                 bson.ObjectID   `json:"id" example:"507f1f77bcf86cd799439011"`
	UserCode           string          `json:"user_code" example:"USR001"`
	FullName           string          `json:"full_name" example:"John Doe"`
	Role               string          `json:"role" example:"Maker"`
	Department         bson.ObjectID   `json:"department" example:"507f1f77bcf86cd799439011"`
	Gender             string          `json:"gender" example:"Male"`
	PhoneNumber        string          `json:"phone_number" example:"+251911234567"`
	Email              string          `json:"email" example:"john.doe@example.com"`
	UserName           string          `json:"username" example:"john.doe"`
	Realm              string          `json:"realm" example:"cps"`
	PermissionCategory []bson.ObjectID `json:"permission_category" example:"[\"507f1f77bcf86cd799439011\"]"`
	PermissionGroup    []bson.ObjectID `json:"permission_group" example:"[\"507f1f77bcf86cd799439011\"]"`

	Enabled      bool       `json:"enabled" example:"true"`
	DateJoined   *time.Time `json:"date_joined" example:"2024-01-15T10:30:00Z"`
	LastModified *time.Time `json:"last_modified" example:"2024-01-15T10:30:00Z"`

	Country string `json:"country" example:"Ethiopia"`
	Region  string `json:"region" example:"Addis Ababa"`
}

type PermissionResponse struct {
	ID   bson.ObjectID `json:"id" bson:"id"`
	Name string        `json:"permission_name" bson:"permission_name"`
}

type PermissionCategoryResponse struct {
	ID           bson.ObjectID        `json:"id" bson:"id"`
	Access        string               `json:"access" bson:"access"` 
	CategoryName string               `json:"category_name" bson:"category_name"`
	Permissions  []PermissionResponse `json:"permissions" bson:"permissions"`
}

type PermissionGroupResponse struct {
	ID         bson.ObjectID                `json:"id" bson:"id"`
	GroupName  string                       `json:"group_name" bson:"group_name"`
	PermissionCategory []PermissionCategoryResponse `json:"permission_category" bson:"permission_category"`
}

type CpsUserResponse struct {
	ID                 bson.ObjectID   `json:"id" bson:"_id"`
	UserCode           string          `json:"user_code" bson:"user_code"`
	FullName           string          `json:"full_name" bson:"full_name"`
	Role               string          `json:"role" bson:"role"`
	Department         bson.ObjectID   `json:"department" bson:"department"`
	Gender             string          `json:"gender" bson:"gender"`
	PhoneNumber        string          `json:"phone_number" bson:"phone_number"`
	Email              string          `json:"email" bson:"email"`
	UserName           string          `json:"username" bson:"username"`
	Realm              string          `json:"realm" bson:"realm"`

	Enabled            bool            `json:"enabled" bson:"enabled"`
	DateJoined         time.Time       `json:"date_joined" bson:"date_joined"`
	LastModified       time.Time       `json:"last_modified" bson:"last_modified"`
	Country            string          `json:"country" bson:"country"`
	Region             string          `json:"region" bson:"region"`
	DepartmentName   string                    `json:"department_name" bson:"department_name"`
	PortalCards      []string                  `json:"portal_cards" bson:"portal_cards"`
	PermissionGroups []PermissionGroupResponse `json:"permission_groups" bson:"permission_groups"`
}


type Department struct {
	ID   bson.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name string        `json:"name,omitempty" bson:"name,omitempty"`
}

type CPSUserWithDepartment struct {
	ID                 bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode           string          `json:"user_code,omitempty" bson:"user_code"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name"`
	Role               string          `json:"role,omitempty" bson:"role"`
	Department         *Department     `json:"department,omitempty" bson:"department"`
	Gender             string          `json:"gender,omitempty" bson:"gender"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number"`
	Email              string          `json:"email,omitempty" bson:"email"`
	UserName           string          `json:"username,omitempty" bson:"username"`
	Realm              string          `json:"realm,omitempty" bson:"realm"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category"`
	PermissionGroup    []bson.ObjectID `json:"permission_group,omitempty" bson:"permission_group"`
	Enabled            bool            `json:"enabled,omitempty" bson:"enabled"`
	DateJoined         *time.Time      `json:"date_joined,omitempty" bson:"date_joined"`
	LastModified       *time.Time      `json:"last_modified,omitempty" bson:"last_modified"`
	Country            string          `json:"country,omitempty" bson:"country"`
	Region             string          `json:"region,omitempty" bson:"region"`
}
