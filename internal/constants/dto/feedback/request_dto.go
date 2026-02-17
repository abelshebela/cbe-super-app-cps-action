package feedback

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FeedbackRequest struct {
	UserCode string `json:"user_code" bson:"user_code,omitempty"`
	Rating   int    `json:"rating" bson:"rating"`
	Comment  string `json:"comment" bson:"comment"`
}

type FeedbackResponse struct {
	ID         bson.ObjectID `json:"id" bson:"_id,omitempty"`
	User       User          `json:"user" bson:"user,omitempty"`
	StarRating int           `json:"rating" bson:"rating"`
	Comment    string        `json:"comment" bson:"comment"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}

type User struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	UserCode    string `json:"user_code" bson:"user_code,omitempty"`
	FullName    string `json:"full_name" bson:"full_name,omitempty"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}
type SurveyFeedbackReq struct {
	UserID     string                    `json:"user_id"`
	FeedbackID string                    `json:"feedback_id"`
	Responses  map[string]types.Response `json:"responses"`
	CreatedAt  time.Time                 `json:"created_at"`
	Metadata   map[string]interface{}    `json:"metadata,omitempty"`
}
type FeedbackKafkaReq struct {
	UserID     string                 `json:"user_id"`
	FeedbackID string                 `json:"feedback_id"`
	StarRating int                    `json:"star_rating"`
	Comment    string                 `json:"comment"`
	CreatedAt  time.Time              `json:"created_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
