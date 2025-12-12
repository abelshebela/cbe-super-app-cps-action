package core

import (
	"time"

	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// BuildFeedbackEntity creates a Feedback model from userID and request with timestamps set.
func BuildFeedbackEntity(userID string, req fbdto.FeedbackRequest) *model.Feedback {
	now := time.Now()
	fb := &model.Feedback{
		UserID:    userID,
		Responses: make(map[string]types.Response),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if len(req.Responses) > 0 {
		fb.Responses = req.Responses
	}
	return fb
}
