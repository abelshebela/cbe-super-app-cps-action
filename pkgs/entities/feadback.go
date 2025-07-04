package entities

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities/type_definition"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Feedback struct {
	ID              bson.ObjectID                    `bson:"_id,omitempty" json:"id,omitempty"`
	Rating          int                              `bson:"rating" json:"rating" validate:"oneof=1 2 3 4 5"`
	Name            string                           `bson:"name" json:"name,omitempty"`
	Message         string                           `bson:"message" json:"message,omitempty"`
	Source          string                           `bson:"source" json:"source,omitempty" validate:"oneof=app website"`
	SurveyResponses []type_definition.SurveyResponse `bson:"survey_responses" json:"survey_responses,omitempty"`
	Enabled         bool                             `bson:"enabled" json:"enabled"`
	IsDeleted       bool                             `bson:"is_deleted" json:"is_deleted"`
	CreatedAt       time.Time                        `bson:"created_at" json:"created_at"`
	LastModified    time.Time                        `bson:"last_modified" json:"last_modified"`
}
