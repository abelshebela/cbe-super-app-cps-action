package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NewsArticle struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID      bson.ObjectID `bson:"category_id" json:"category_id"`
	Title           string        `bson:"title" json:"title"`
	Content         string        `bson:"content" json:"content"`
	Tags            []string      `bson:"tags" json:"tags"`
	Language        string        `bson:"language" json:"language"`
	Author          string        `bson:"author" json:"author"`
	Thumbnail       string        `bson:"thumbnail" json:"thumbnail"`
	IsPublished     bool          `bson:"is_published" json:"is_published"`
	PublishedAt     *time.Time    `bson:"published_at,omitempty" json:"published_at,omitempty"`
	Views           int           `bson:"views" json:"views"`
	CopyLinkCounter int           `bson:"copy_link_counter" json:"copy_link_counter"`
	Slug            string        `bson:"slug" json:"slug"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
	IsDeleted       bool          `bson:"is_deleted" json:"is_deleted"`
}

type NewsCategoryModel struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Slug        string        `bson:"slug" json:"slug"`
	IsActive    bool          `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	IsDeleted   bool          `bson:"is_deleted" json:"is_deleted"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

type ShortVideo struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID  bson.ObjectID `bson:"category_id" json:"category_id"`
	Title       string        `bson:"title" json:"title"`
	Caption     string        `bson:"caption" json:"caption"`
	Tags        []string      `bson:"tags" json:"tags"`
	Thumbnail   string        `bson:"thumbnail" json:"thumbnail"`
	VideoURL    string        `bson:"video_url" json:"video_url"`
	Duration    float64       `bson:"duration" json:"duration"`
	ViewsCount  int64         `bson:"views_count" json:"views_count"`
	SharesCount int64         `bson:"shares_count" json:"shares_count"`
	IsDeleted   bool          `bson:"is_deleted" json:"is_deleted"`
	IsPublished bool          `bson:"is_published" json:"is_published"`
	Slug        string        `bson:"slug" json:"slug"`
	Author      string        `bson:"author" json:"author"`
	Language    string        `bson:"language" json:"language"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	PublishedAt *time.Time    `bson:"published_at,omitempty" json:"published_at,omitempty"`
}
