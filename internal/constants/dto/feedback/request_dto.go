package feedback

import (
	"cbe-super-app-cps-action/internal/constants/types"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FeedbackRequest struct {
	Responses map[string]types.Response `json:"responses" bson:"responses"`
}

func (r FeedbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Responses, validation.Required.Error("responses is required")),
	)
}
