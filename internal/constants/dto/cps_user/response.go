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
