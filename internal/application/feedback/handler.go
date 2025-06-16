package feedback

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/feedback/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/feedback"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type FeedbackService interface {
    GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*entity.FeedbackResponse, error)
	GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error)
}

type FeedbackHandler struct {
	domain *feedback.FeedbackDomain
	logger utils.Logger
}

func InitFeedbackHandler(feedbackDomain *feedback.FeedbackDomain, logger utils.Logger) FeedbackService {
	return FeedbackHandler{
		domain: feedbackDomain,
		logger: logger,
	}
}

func (f FeedbackHandler) GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*entity.FeedbackResponse, error) {
	feedbacks, err := f.domain.GetFeedbacks(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (f FeedbackHandler) GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error) {
	feedback, err := f.domain.GetFeedbackByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}
