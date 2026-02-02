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
	JobTitle           string          `json:"job_title" example:"Branch Manager"`
	Enabled            bool            `json:"enabled" example:"true"`
	DateJoined         *time.Time      `json:"date_joined" example:"2024-01-15T10:30:00Z"`
	LastModified       *time.Time      `json:"last_modified" example:"2024-01-15T10:30:00Z"`
	MakerAllocations   []string        `json:"maker_allocations" example:"[\"Action1\", \"Action2\"]"`
	CheckerAllocations []string        `json:"checker_allocations" example:"[[\"Action1_Checker1\", \"Action1_Checker2\"], [\"Action2_Checker1\"]]"`
	AuditorAllocations []string        `json:"auditor_allocations" example:"[\"Action1\", \"Action2\"]"`
	PortalCards        []string        `json:"portal_cards" bson:"portal_cards"` // top-level from job_roles
	Country            string          `json:"country" example:"Ethiopia"`
	Region             string          `json:"region" example:"Addis Ababa"`
}

type CPSUserResponse struct {
	ID         string             `json:"id" example:"507f1f77bcf86cd799439011"`
	UserCode   string             `json:"user_code" example:"USR001"`
	FullName   string             `json:"full_name" example:"John Doe"`
	Role       string             `json:"role" example:"Maker"`
	Department DepartmentResponse `json:"department" example:"507f1f77bcf86cd799439011"`
	// Department         bson.ObjectID   `json:"department" example:"507f1f77bcf86cd799439011"`
	Gender             string          `json:"gender" example:"Male"`
	PhoneNumber        string          `json:"phone_number" example:"+251911234567"`
	Email              string          `json:"email" example:"john.doe@example.com"`
	UserName           string          `json:"username" example:"john.doe"`
	Realm              string          `json:"realm" example:"cps"`
	PermissionCategory []bson.ObjectID `json:"permission_category" example:"[\"507f1f77bcf86cd799439011\"]"`
	PermissionGroup    []bson.ObjectID `json:"permission_group" example:"[\"507f1f77bcf86cd799439011\"]"`
	JobTitle           string          `json:"job_title" example:"Branch Manager"`
	Enabled            bool            `json:"enabled" example:"true"`
	DateJoined         *time.Time      `json:"date_joined" example:"2024-01-15T10:30:00Z"`
	LastModified       *time.Time      `json:"last_modified" example:"2024-01-15T10:30:00Z"`
	MakerAllocations   []string        `json:"maker_allocations" example:"[\"Action1\", \"Action2\"]"`
	CheckerAllocations []string        `json:"checker_allocations" example:"[[\"Action1_Checker1\", \"Action1_Checker2\"], [\"Action2_Checker1\"]]"`
	AuditorAllocations []string        `json:"auditor_allocations" example:"[\"Action1\", \"Action2\"]"`
	PortalCards        []string        `json:"portal_cards" bson:"portal_cards"` // top-level from job_roles
	Country            string          `json:"country" example:"Ethiopia"`
	Region             string          `json:"region" example:"Addis Ababa"`
}

type PermissionResponse struct {
	ID   bson.ObjectID `json:"id" bson:"id"`
	Name string        `json:"permission_name" bson:"permission_name"`
}

type PermissionCategoryResponse struct {
	ID           bson.ObjectID        `json:"id" bson:"id"`
	Access       string               `json:"access" bson:"access"`
	CategoryName string               `json:"category_name" bson:"category_name"`
	Permissions  []PermissionResponse `json:"permissions" bson:"permissions"`
}

type PermissionGroupResponse struct {
	ID                 bson.ObjectID                `json:"id" bson:"id"`
	GroupName          string                       `json:"group_name" bson:"group_name"`
	PermissionCategory []PermissionCategoryResponse `json:"permission_category" bson:"permission_category"`
}

type DepartmentResponse struct {
	ID          bson.ObjectID `json:"id" bson:"id"`
	Name        string        `json:"name" bson:"name"`
	PortalCards []string      `json:"portal_cards" bson:"portal_cards"`
}

type CpsUserResponse struct {
	ID                 bson.ObjectID             `json:"id" bson:"_id"` // use primitive.ObjectID instead of bson.ObjectID
	UserCode           string                    `json:"user_code" bson:"user_code"`
	FullName           string                    `json:"full_name" bson:"full_name"`
	Role               string                    `json:"role,omitempty" bson:"role"`   // optional, can keep empty
	Department         *DepartmentResponse       `json:"department" bson:"department"` // populated via $lookup
	JobTitle           string                    `json:"job_title" bson:"job_title"`
	Gender             string                    `json:"gender" bson:"gender"`
	PhoneNumber        string                    `json:"phone_number" bson:"phone_number"`
	Email              string                    `json:"email" bson:"email"`
	UserName           string                    `json:"username,omitempty" bson:"username"` // optional
	Realm              string                    `json:"realm" bson:"realm"`
	Enabled            bool                      `json:"enabled" bson:"enabled"`
	DateJoined         time.Time                 `json:"date_joined" bson:"date_joined"`
	LastModified       time.Time                 `json:"last_modified" bson:"last_modified"`
	Country            string                    `json:"country,omitempty" bson:"country"`                     // optional
	Region             string                    `json:"region,omitempty" bson:"region"`                       // optional
	PortalCards        []string                  `json:"portal_cards" bson:"portal_cards"`                     // top-level from job_roles
	PermissionGroups   []PermissionGroupResponse `json:"permission_groups,omitempty" bson:"permission_groups"` // optional
	PermissionCategory any                       `json:"permission_category,omitempty" bson:"permission_category"`
	MakerAllocations   []string                  `json:"maker_allocations" example:"[\"Action1\", \"Action2\"]"`
	CheckerAllocations []string                  `json:"checker_allocations" example:"[[\"Action1_Checker1\", \"Action1_Checker2\"], [\"Action2_Checker1\"]]"`
	AuditorAllocations []string                  `json:"auditor_allocations" example:"[\"Action1\", \"Action2\"]"`
}

type RoleResponse struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

type CpsUserPopulatedResponse struct {
	ID                 bson.ObjectID             `json:"id" bson:"_id"` // use primitive.ObjectID instead of bson.ObjectID
	UserCode           string                    `json:"user_code" bson:"user_code"`
	FullName           string                    `json:"full_name" bson:"full_name"`
	Role               RoleResponse              `json:"role,omitempty" bson:"role"` // optional, can keep empty
	RoleCode           string                    `json:"role_code" bson:"role_code"`
	Department         *DepartmentResponse       `json:"department" bson:"department"` // populated via $lookup
	JobTitle           string                    `json:"job_title" bson:"job_title"`
	Gender             string                    `json:"gender" bson:"gender"`
	PhoneNumber        string                    `json:"phone_number" bson:"phone_number"`
	Email              string                    `json:"email" bson:"email"`
	UserName           string                    `json:"username,omitempty" bson:"username"` // optional
	Realm              string                    `json:"realm" bson:"realm"`
	Enabled            bool                      `json:"enabled" bson:"enabled"`
	DateJoined         time.Time                 `json:"date_joined" bson:"date_joined"`
	LastModified       time.Time                 `json:"last_modified" bson:"last_modified"`
	Country            string                    `json:"country,omitempty" bson:"country"`                     // optional
	Region             string                    `json:"region,omitempty" bson:"region"`                       // optional
	PortalCards        []string                  `json:"portal_cards" bson:"portal_cards"`                     // top-level from job_roles
	PermissionGroups   []PermissionGroupResponse `json:"permission_groups,omitempty" bson:"permission_groups"` // optional
	PermissionCategory any                       `json:"permission_category,omitempty" bson:"permission_category"`
	MakerAllocations   []string                  `json:"maker_allocations" example:"[\"Action1\", \"Action2\"]"`
	CheckerAllocations []string                  `json:"checker_allocations" example:"[[\"Action1_Checker1\", \"Action1_Checker2\"], [\"Action2_Checker1\"]]"`
	AuditorAllocations []string                  `json:"auditor_allocations" example:"[\"Action1\", \"Action2\"]"`
}

type Department struct {
	ID   bson.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name string        `json:"name,omitempty" bson:"name,omitempty"`
}

type CPSUserWithDepartment struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserCode     string        `json:"user_code,omitempty" bson:"user_code"`
	FullName     string        `json:"full_name,omitempty" bson:"full_name"`
	Role         string        `json:"role,omitempty" bson:"role"`
	Gender       string        `json:"gender,omitempty" bson:"gender"`
	PhoneNumber  string        `json:"phone_number,omitempty" bson:"phone_number"`
	Email        string        `json:"email,omitempty" bson:"email"`
	UserName     string        `json:"username,omitempty" bson:"username"`
	JobTitle     string        `json:"job_title,omitempty" bson:"job_title"`
	Enabled      bool          `json:"enabled" bson:"enabled"`
	DateJoined   *time.Time    `json:"date_joined,omitempty" bson:"date_joined"`
	LastModified *time.Time    `json:"last_modified,omitempty" bson:"last_modified"`
	Country      string        `json:"country,omitempty" bson:"country"`
	Region       string        `json:"region,omitempty" bson:"region"`
}
