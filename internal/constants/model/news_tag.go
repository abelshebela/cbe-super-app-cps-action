package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NewsTag struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	TagName        string        `bson:"tag_name" json:"tag_name"`
	IsDeleted      bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time     `bson:"created_at,omitempty" json:"created_at"`
	DeletedAt      time.Time     `bson:"deleted_at,omitempty" json:"deleted_at"`
	LastModifiedAt time.Time     `bson:"last_modified_at,omitempty" json:"last_modified_at"`
}

type NewsTagCPSAction struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	TagNameList    []string      `bson:"tag_name_list" json:"tag_name_list"`
	TagName        string        `bson:"tag_name" json:"tag_name"`
	IsDeleted      bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time     `bson:"created_at,omitempty" json:"created_at"`
	DeletedAt      time.Time     `bson:"deleted_at,omitempty" json:"deleted_at"`
	LastModifiedAt time.Time     `bson:"last_modified_at,omitempty" json:"last_modified_at"`
}
