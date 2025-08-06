package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FeedbackRequest struct {
	Responses map[string]Response `json:"responses" bson:"responses"`
}

func (f FeedbackRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Responses, validation.Required.Error("responses are required")),
	)
}
 