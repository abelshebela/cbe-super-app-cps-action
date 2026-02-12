package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PasswordRule struct {
	ID                     bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name                   string        `bson:"name" json:"name"`
	MinLength              int           `bson:"min_length" json:"min_length"`
	MaxLength              int           `bson:"max_length" json:"max_length"`
	Numbers                *bool         `bson:"numbers" json:"numbers"`
	CapitalLetters         *bool         `bson:"capital_letters" json:"capital_letters"`
	AllowSequentialNumbers *bool         `bson:"allow_sequential_numbers" json:"allow_sequential_numbers"`
	SmallLetters           *bool         `bson:"small_letters" json:"small_letters"`
	Characters             *bool         `bson:"characters" json:"characters"`
	IsSpacedAllowed        *bool         `json:"is_spaced_allowed" bson:"is_spaced_allowed"`
	CreatedAt              time.Time     `bson:"created_at" json:"created_at"`
}
