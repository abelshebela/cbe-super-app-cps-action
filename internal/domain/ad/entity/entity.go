package entity

import (
	"mime/multipart"
	"time"
)

type AdvertFor string

const (
	IFB  AdvertFor = "IFB"
	CB   AdvertFor = "CB"
	Both AdvertFor = "ALL"
)

type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type Advert struct {
	ID            string     `json:"_id,omitempty" bson:"_id"`
	Title         string     `json:"title" bson:"title"`
	Description   string     `json:"description" bson:"description"`
	BannerImage   string     `json:"banner_image" bson:"banner_image"`
	AdvertFor     AdvertFor  `json:"advert_for" bson:"advert_for"`
	Date          AdvertDate `json:"advert_date" bson:"advert_date"`
	Enabled       bool       `json:"enabled" bson:"enabled"`
	IsDeleted     bool       `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time  `json:"deleted_at" bson:"deleted_at"`
	CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at" bson:"last_updated_at"`
}

type AdvertResponse struct {
	ID            string     `json:"id,omitempty" bson:"_id"`
	Title         string     `json:"title" bson:"title"`
	Description   string     `json:"description" bson:"description"`
	BannerImage   string     `json:"banner_image" bson:"banner_image"`
	AdvertFor     AdvertFor  `json:"advert_for" bson:"advert_for"`
	Date          AdvertDate `json:"advert_date" bson:"advert_date"`
	Enabled       bool       `json:"enabled" bson:"enabled"`
	IsDeleted     bool       `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time  `json:"deleted_at" bson:"deleted_at"`
	CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at" bson:"last_updated_at"`
}

func ToAdvertResponse(ad Advert) AdvertResponse {
	return AdvertResponse{
		ID:            ad.ID,
		Title:         ad.Title,
		Description:   ad.Description,
		BannerImage:   ad.BannerImage,
		AdvertFor:     ad.AdvertFor,
		Date:          ad.Date,
		Enabled:       ad.Enabled,
		IsDeleted:     ad.IsDeleted,
		DeletedAt:     ad.DeletedAt,
		CreatedAt:     ad.CreatedAt,
		LastUpdatedAt: ad.LastUpdatedAt,
	}
}

type DeletedAdvert struct {
	ID            string    `json:"_id,omitempty" bson:"_id"`
	IsDeleted     bool      `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time `json:"deleted_at" bson:"deleted_at"`
	LastUpdatedAt time.Time `json:"last_updated_at" bson:"last_updated_at"`
}

type UpdateAdvert struct {
	Title         string      `json:"title,omitempty" bson:"title"`
	Description   string      `json:"description,omitempty" bson:"description"`
	BannerImage   string      `json:"banner_image,omitempty" bson:"banner_image"`
	AdvertFor     AdvertFor   `json:"advert_for,omitempty" bson:"advert_for"`
	Date          *AdvertDate `json:"advert_date,omitempty" bson:"advert_date"`
	Enabled       bool        `json:"enabled" bson:"enabled"`
	IsDeleted     bool        `json:"is_deleted" bson:"is_deleted"`
	LastUpdatedAt time.Time   `json:"last_updated_at" bson:"last_updated_at"`
}

type CreateAdvert struct {
	ID          string                `json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"descritption"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   AdvertFor             `form:"advert_for" json:"advert_for"`
	Date        AdvertDate            `form:"date" json:"date"`
}
