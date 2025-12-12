package feedback

import (
	"time"
	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FeedbackRequest struct {
	Responses map[string]shared_types.Response `json:"responses" bson:"responses"`
}

type FeedbackResponse struct {
	ID             bson.ObjectID             `json:"id" bson:"_id,omitempty"`
	User           User                      `json:"user" bson:"user,omitempty"`
	Responses      map[string]shared_types.Response `json:"responses" bson:"responses"`
	CreatedAt      time.Time                 `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at" bson:"updated_at"`
}

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	UserCode string `json:"user_code" bson:"user_code,omitempty"`
	FullName string `json:"full_name" bson:"full_name,omitempty"`
}
