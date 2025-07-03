package entity

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type SurveyResponse struct {
	Question string `bson:"question" json:"question" validate:"required"`
	Answer   string `bson:"answer" json:"answer" validate:"required"`
}

type Feedback struct {
	ID              bson.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	Rating          int              `bson:"rating" json:"rating" validate:"oneof=1 2 3 4 5"`
	Name            string           `bson:"name" json:"name,omitempty"`
	Message         string           `bson:"message" json:"message,omitempty"`
	Source          string           `bson:"source" json:"source,omitempty" validate:"oneof=app website"`
	SurveyResponses []SurveyResponse `bson:"survey_responses" json:"survey_responses,omitempty"`
	Enabled         bool             `bson:"enabled" json:"enabled"`
	IsDeleted       bool             `bson:"is_deleted" json:"is_deleted"`
	CreatedAt       time.Time        `bson:"created_at" json:"created_at"`
	LastModified    time.Time        `bson:"last_modified" json:"last_modified"`
}

type FeedbackResponse struct {
	Page      int         `json:"page"`
	Feedbacks []*Feedback `json:"feedbacks"`
	Limit     int         `json:"limit"`
	Total     int64       `json:"total"`
}
