package core

import (
	"time"

	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

// BuildFeedbackEntity creates a Feedback model from userID and request with timestamps set.
func BuildFeedbackEntity(userID string, req fbdto.FeedbackRequest) *model.Feedback {
	now := time.Now()
	fb := &model.Feedback{
		UserID:    userID,
		Responses: make(map[string]shared_types.Response),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if len(req.Responses) > 0 {
		fb.Responses = req.Responses
	}
	return fb
}
func BuildSurveyFeedbackEntity(surveyFeedback fbdto.SurveyFeedbackReq) *imodel.SurveyFeedback {
	return &imodel.SurveyFeedback{
		UserID:     surveyFeedback.UserID,
		StarRating: surveyFeedback.StarRating,
		Comment:    surveyFeedback.Comment,
		CreatedAt:  surveyFeedback.CreatedAt,
		Metadata:   surveyFeedback.Metadata,
	}
}
