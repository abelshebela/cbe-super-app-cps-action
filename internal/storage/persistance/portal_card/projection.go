package portal_card

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CardMapper(card model.Card) bson.M {
	return bson.M{
		"$set": bson.M{
			"card_name":   card.CardName,
			"sub_cards":   card.SubCards,
			"updated_at":  time.Now(),
		},
	}
}