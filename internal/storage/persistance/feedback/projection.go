package feedback

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// FeedbackMapper maps a Feedback model to a bson.M for updates
func FeedbackMapper(feedback model.Feedback) bson.M {
	return bson.M{
		"$set": bson.M{
			"user_id":    feedback.UserID,
			"responses":  feedback.Responses,
			"updated_at": time.Now(),
		},
	}
}
