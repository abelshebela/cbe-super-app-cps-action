package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NewsCategory struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	CategoryName   string        `bson:"category_name" json:"category_name"`
	IsDeleted      bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time     `bson:"created_at,omitempty" json:"created_at"`
	DeletedAt      time.Time     `bson:"deleted_at,omitempty" json:"deleted_at"`
	LastModifiedAt time.Time     `bson:"last_modified_at,omitempty" json:"last_modified_at"`
}

type NewsCategoryCPSAction struct {
	ID               bson.ObjectID `bson:"_id" json:"id"`
	CategoryNameList []string      `bson:"category_name_list" json:"category_name_list"`
	CategoryName     string        `bson:"category_name" json:"category_name"`
	IsDeleted        bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt        time.Time     `bson:"created_at,omitempty" json:"created_at"`
	DeletedAt        time.Time     `bson:"deleted_at,omitempty" json:"deleted_at"`
	LastModifiedAt   time.Time     `bson:"last_modified_at,omitempty" json:"last_modified_at"`
}
