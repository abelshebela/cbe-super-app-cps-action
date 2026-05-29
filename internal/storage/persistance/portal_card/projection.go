package portal_card

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CardMapper(card model.Card) bson.M {
	return bson.M{
		"$set": bson.M{
			"card_name":  card.CardName,
			"sub_cards":  card.SubCards,
			"updated_at": time.Now(),
		},
	}
}
