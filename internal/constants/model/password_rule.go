package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PasswordRule struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name           string        `bson:"name" json:"name"`
	MinLength      int           `bson:"min_length" json:"min_length"`
	MaxLength      int           `bson:"max_length" json:"max_length"`
	Numbers        *bool         `bson:"numbers" json:"numbers"`
	CapitalLetters *bool         `bson:"capital_letters" json:"capital_letters"`
	SmallLetters   *bool         `bson:"small_letters" json:"small_letters"`
	Characters     *bool         `bson:"characters" json:"characters"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
}
