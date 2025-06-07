package ad

import "time"

type AdvertFor string

const (
	IFB  AdvertFor = "IFB"
	CB   AdvertFor = "CB"
	Both AdvertFor = "ALL"
)

type AdvertDate struct {
	StartedAt time.Time
	ExpiredAt time.Time
}

type Advert struct {
	ID            string
	Title         string
	Description   string
	BannerImage   string
	AdvertFor     AdvertFor
	Date          AdvertDate
	Enabled       bool
	IsDeleted     bool
	DeletedAt     time.Time
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}
