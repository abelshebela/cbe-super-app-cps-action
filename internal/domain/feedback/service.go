package feedback

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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

func (f *FeedbackDomain) GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Feedback], error) {
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
