package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Card struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CardName string        `bson:"card_name" json:"card_name"`
	SubCards []string      `bson:"sub_cards" json:"sub_cards"`
}
