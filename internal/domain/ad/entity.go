package ad

import (
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
