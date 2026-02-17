package core

import (
	"time"

	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

// BuildFeedbackEntity creates a Feedback model from userID and request with timestamps set.
func BuildFeedbackEntity(userID string, req fbdto.FeedbackRequest, user *member.User) *imodel.Feedback {
	now := time.Now()
	fb := &imodel.Feedback{
		UserCode: user.UserCode,
		// AccountNumber: user.AccountNumber,
		CustomerName: user.FullName,
		PhoneNumber:  user.PhoneNumber,
		Email:        user.Email,
		SendAt:       time.Now(),
		Rating:       req.Rating,
		Comment:      req.Comment,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return fb
}
func BuildSurveyFeedbackEntity(surveyFeedback fbdto.SurveyFeedbackReq, user *member.User) *imodel.SurveyFeedback {
	return &imodel.SurveyFeedback{
		UserCode:     user.UserCode,
		Email:        user.Email,
		CustomerName: user.FullName,
		PhoneNumber:  user.PhoneNumber,
		SentAt:       time.Now(),
		// AccountNumber: user.AccountNumber,
		Responses: surveyFeedback.Responses,
		CreatedAt: surveyFeedback.CreatedAt,
		Metadata:  surveyFeedback.Metadata,
	}
}
