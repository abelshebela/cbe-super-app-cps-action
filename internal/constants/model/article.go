package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NewsArticle struct {
	ID               bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID       bson.ObjectID   `bson:"category_id" json:"category_id"`
	Title            string          `bson:"title" json:"title"`
	Content          string          `bson:"content" json:"content"`
	Tags             []bson.ObjectID `bson:"tags" json:"tags"`
	Author           string          `bson:"author" json:"author"`
	Thumbnail        string          `bson:"thumbnail" json:"thumbnail"`
	ThumbnailAltText string          `bson:"thumbnail_alt_text" json:"thumbnail_alt_text"`
	IsPublished      bool            `bson:"is_published" json:"is_published"`
	PublishedAt      *time.Time      `bson:"published_at,omitempty" json:"published_at,omitempty"`
	Views            int             `bson:"views" json:"views"`
	CopyLinkCounter  int             `bson:"copy_link_counter" json:"copy_link_counter"`
	Slug             string          `bson:"slug" json:"slug"`
	CreatedAt        time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `bson:"updated_at" json:"updated_at"`
	IsDeleted        bool            `bson:"is_deleted" json:"is_deleted"`
	DeletedAt        *time.Time      `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
type NewsArticleDetail struct {
	ID               bson.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
	Category         NewsCategoryModel `bson:"category" json:"category"`
	Title            string            `bson:"title" json:"title"`
	Content          string            `bson:"content" json:"content"`
	Tags             []NewsTags        `bson:"tags" json:"tags"`
	Language         string            `bson:"language" json:"language"`
	Author           string            `bson:"author" json:"author"`
	Thumbnail        string            `bson:"thumbnail" json:"thumbnail"`
	ThumbnailAltText string            `bson:"thumbnail_alt_text" json:"thumbnail_alt_text"`
	IsPublished      bool              `bson:"is_published" json:"is_published"`
	PublishedAt      *time.Time        `bson:"published_at,omitempty" json:"published_at,omitempty"`
	Views            int               `bson:"views" json:"views"`
	CopyLinkCounter  int               `bson:"copy_link_counter" json:"copy_link_counter"`
	Slug             string            `bson:"slug" json:"slug"`
	CreatedAt        time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `bson:"updated_at" json:"updated_at"`
	IsDeleted        bool              `bson:"is_deleted" json:"is_deleted"`
	DeletedAt        *time.Time        `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type NewsCategoryModel struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Color     string        `bson:"color" json:"color"`
	IsActive  bool          `bson:"is_active" json:"is_active"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	IsDeleted bool          `bson:"is_deleted" json:"is_deleted"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type ShortVideo struct {
	ID               bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	CategoryID       bson.ObjectID   `bson:"category_id" json:"category_id"`
	Title            string          `bson:"title" json:"title"`
	Caption          string          `bson:"caption" json:"caption"`
	Tags             []bson.ObjectID `bson:"tags" json:"tags"`
	Thumbnail        string          `bson:"thumbnail" json:"thumbnail"`
	ThumbnailAltText string          `bson:"thumbnail_alt_text" json:"thumbnail_alt_text"`
	VideoURL         string          `bson:"video_url" json:"video_url"`
	Duration         float64         `bson:"duration" json:"duration"`
	ViewsCount       int64           `bson:"views_count" json:"views_count"`
	SharesCount      int64           `bson:"shares_count" json:"shares_count"`
	IsDeleted        bool            `bson:"is_deleted" json:"is_deleted"`
	IsPublished      bool            `bson:"is_published" json:"is_published"`
	Slug             string          `bson:"slug" json:"slug"`
	Author           string          `bson:"author" json:"author"`
	CreatedAt        time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time      `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	PublishedAt      *time.Time      `bson:"published_at,omitempty" json:"published_at,omitempty"`
}

type ShortVideoDetail struct {
	ID               bson.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Category         NewsCategoryModel `bson:"category" json:"category"`
	Title            string            `bson:"title" json:"title"`
	Caption          string            `bson:"caption" json:"caption"`
	Tags             []NewsTags        `bson:"tags" json:"tags"`
	Thumbnail        string            `bson:"thumbnail" json:"thumbnail"`
	ThumbnailAltText string            `bson:"thumbnail_alt_text" json:"thumbnail_alt_text"`
	VideoURL         string            `bson:"video_url" json:"video_url"`
	Duration         float64           `bson:"duration" json:"duration"`
	ViewsCount       int64             `bson:"views_count" json:"views_count"`
	SharesCount      int64             `bson:"shares_count" json:"shares_count"`
	IsDeleted        bool              `bson:"is_deleted" json:"is_deleted"`
	IsPublished      bool              `bson:"is_published" json:"is_published"`
	Slug             string            `bson:"slug" json:"slug"`
	Author           string            `bson:"author" json:"author"`
	Language         string            `bson:"language" json:"language"`
	CreatedAt        time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time        `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	PublishedAt      *time.Time        `bson:"published_at,omitempty" json:"published_at,omitempty"`
}

type NewsTags struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Color     string        `bson:"color" json:"color"`
	IsEnabled bool          `bson:"is_enabled" json:"is_enabled"`
	IsDeleted bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (m *NewsArticleDetail) ToNewsArticle() *NewsArticle {
	return &NewsArticle{
		ID:               m.ID,
		CategoryID:       m.Category.ID,
		Title:            m.Title,
		Content:          m.Content,
		Tags:             extractTagIDs(m.Tags),
		Thumbnail:        m.Thumbnail,
		ThumbnailAltText: m.ThumbnailAltText,
		IsDeleted:        m.IsDeleted,
		IsPublished:      m.IsPublished,
		Slug:             m.Slug,
		Author:           m.Author,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		PublishedAt:      m.PublishedAt,
		Views:            m.Views,
		CopyLinkCounter:  m.CopyLinkCounter,
		DeletedAt:        m.DeletedAt,
	}
}

func (m *ShortVideoDetail) ToShortVideo() *ShortVideo {
	return &ShortVideo{
		ID:               m.ID,
		CategoryID:       m.Category.ID,
		Title:            m.Title,
		Caption:          m.Caption,
		Tags:             extractTagIDs(m.Tags),
		Thumbnail:        m.Thumbnail,
		ThumbnailAltText: m.ThumbnailAltText,
		VideoURL:         m.VideoURL,
		Duration:         m.Duration,
		IsDeleted:        m.IsDeleted,
		IsPublished:      m.IsPublished,
		Slug:             m.Slug,
		Author:           m.Author,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		PublishedAt:      m.PublishedAt,
		ViewsCount:       m.ViewsCount,
		SharesCount:      m.SharesCount,
		DeletedAt:        m.DeletedAt,
	}
}

func extractTagIDs(tags []NewsTags) []bson.ObjectID {
	var tagIDs []bson.ObjectID
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
	}
	return tagIDs
}
