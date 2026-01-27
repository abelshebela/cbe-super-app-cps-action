package feedback

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (r FeedbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Responses, validation.Required.Error("responses is required")),
	)
}

func (r SurveyFeedbackReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required.Error("user_id is required")),
		validation.Field(&r.FeedbackID, validation.Required.Error("feedback_id is required")),
		validation.Field(&r.StarRating, validation.Required.Error("star_rating is required")),
	)
}
