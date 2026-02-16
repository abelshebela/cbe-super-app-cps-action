package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRole struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Code           string        `json:"code" bson:"code"`
	Name           string        `json:"name" bson:"name"`
	Enable         bool          `json:"enabled" bson:"enabled"`
	PortalCards    []string      `json:"portal_cards" bson:"portal_cards"`
	MakerActions   []string      `json:"maker_actions,omitempty" bson:"-"`
	CheckerActions []string      `json:"checker_actions,omitempty" bson:"-"`
	AuditorActions []string      `json:"auditor_actions,omitempty" bson:"-"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
}
