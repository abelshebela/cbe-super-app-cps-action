package entities

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/entities/type_definition"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AdvertFor string

const (
	IFB  AdvertFor = "IFB"
	CB   AdvertFor = "CB"
	Both AdvertFor = "ALL"
)

type Advert struct {
	ID            bson.ObjectID `bson:"_id" json:"id"`
	Title         string
	Description   string
	BannerImage   string
	AdvertFor     AdvertFor
	Date          type_definition.AdvertDate
	Enabled       bool
	IsDeleted     bool
	DeletedAt     time.Time
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}
