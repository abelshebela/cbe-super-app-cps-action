package cpsuser

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUserDTO struct {
	ID                 bson.ObjectID   `json:"id"`
	UserCode           string          `json:"user_code"`
	FullName           string          `json:"full_name"`
	Role               string          `json:"role"`
	Department         bson.ObjectID   `json:"department"`
	Gender             string          `json:"gender"`
	PhoneNumber        string          `json:"phone_number"`
	Email              string          `json:"email"`
	UserName           string          `json:"username"`
	Realm              string          `json:"realm"`
	PermissionCategory []bson.ObjectID `json:"permission_category"`
	PermissionGroup    []bson.ObjectID `json:"permission_group"`

	Enabled      bool       `json:"enabled"`
	DateJoined   *time.Time `json:"date_joined"`
	LastModified *time.Time `json:"last_modified"`

	Country string `json:"country"`
	Region  string `json:"region"`
}
