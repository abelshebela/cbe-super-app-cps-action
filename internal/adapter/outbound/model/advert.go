package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)


type AdvertFor string

type AdvertDate struct {
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type Advert struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id"`
	Title         string        `json:"title" bson:"title"`
	Description   string        `json:"description" bson:"description"`
	BannerImage   string        `json:"banner_image" bson:"banner_image"`
	AdvertFor     AdvertFor     `json:"advert_for" bson:"advert_for"`
	Date          AdvertDate    `json:"advert_date" bson:"advert_date"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time     `json:"deleted_at" bson:"deleted_at"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	LastUpdatedAt time.Time     `json:"last_updated_at" bson:"last_updated_at"`
}
