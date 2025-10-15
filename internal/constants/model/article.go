package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NewsArticle struct {
	ID              bson.ObjectID `bson:"_id,omitempty"`
	CategoryID      bson.ObjectID `bson:"category_id"`
	Title           string        `bson:"title"`
	Content         string        `bson:"content"`
	Tags            []string      `bson:"tags"`
	Language        string        `bson:"language"`
	Author          string        `bson:"author"`
	Thumbnail       string        `bson:"thumbnail"`
	IsPublished     bool          `bson:"is_published"`
	PublishedAt     *time.Time    `bson:"published_at,omitempty"`
	Views           int           `bson:"views"`
	CopyLinkCounter int           `bson:"copy_link_counter"`
	Slug            string        `bson:"slug"`
	CreatedAt       time.Time     `bson:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at"`
	IsDeleted       bool          `bson:"is_deleted"`
}


type NewsCategoryModel struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Slug        string        `bson:"slug"`
	IsActive    bool          `bson:"is_active"`
	CreatedAt   time.Time     `bson:"created_at"`
	IsDeleted   bool          `bson:"is_deleted"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}