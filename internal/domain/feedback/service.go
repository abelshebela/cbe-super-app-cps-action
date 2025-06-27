package feedback

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/feedback/entity"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FeedbackDomain struct {
	feedbackService FeedbackRepository
	logger          utils.Logger
}

func InitFeedbackDomain(feedbackRepo FeedbackRepository, logger utils.Logger) *FeedbackDomain {
	return &FeedbackDomain{
		feedbackService: feedbackRepo,
		logger:          logger,
	}
}

func (f *FeedbackDomain) GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*entity.FeedbackResponse, error) {
	feedbacks, err := f.feedbackService.GetFeedbacks(ctx, filterParams)
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (f *FeedbackDomain) GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error) {
	feedback, err := f.feedbackService.GetFeedbackByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return feedback, nil
}
