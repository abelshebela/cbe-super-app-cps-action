package feedback

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (r FeedbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Responses, validation.Required.Error("responses is required")),
	)
}
