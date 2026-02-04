package core

import (
	"time"

	customer_dto "cbe-super-app-cps-action/internal/constants/dto/customer"
	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	imodel "cbe-super-app-cps-action/internal/constants/model"
)

// BuildFeedbackEntity creates a Feedback model from userID and request with timestamps set.
func BuildFeedbackEntity(userID string, req fbdto.FeedbackRequest, user *customer_dto.CustomerListResponse) *imodel.Feedback {
	now := time.Now()
	fb := &imodel.Feedback{
		UserID:        userID,
		AccountNumber: user.AccountNumber,
		CustomerName:  user.FullName,
		PhoneNumber:   user.PhoneNumber,
		Email:         user.Email,
		SendAt:        time.Now(),
		Rating:        req.Rating,
		Comment:       req.Comment,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return fb
}
func BuildSurveyFeedbackEntity(surveyFeedback fbdto.SurveyFeedbackReq, user *customer_dto.CustomerListResponse) *imodel.SurveyFeedback {
	return &imodel.SurveyFeedback{
		UserID:        surveyFeedback.UserID,
		Email:         user.Email,
		CustomerName:  user.FullName,
		PhoneNumber:   user.PhoneNumber,
		SentAt:        time.Now(),
		AccountNumber: user.AccountNumber,
		Responses:     surveyFeedback.Responses,
		CreatedAt:     surveyFeedback.CreatedAt,
		Metadata:      surveyFeedback.Metadata,
	}
}
