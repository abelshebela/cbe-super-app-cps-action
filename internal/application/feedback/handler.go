package feedback

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FeedbackService interface {
	GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Feedback], error)
	GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error)
	CreateFeedback(ctx context.Context, req entity.FeedbackRequest, userID string) (*entity.Feedback, error)
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

func (f FeedbackHandler) GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Feedback], error) {
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

func (f FeedbackHandler) CreateFeedback(ctx context.Context, req entity.FeedbackRequest, userID string) (*entity.Feedback, error) {
	feedback, err := f.domain.CreateFeedback(ctx, req, userID)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}
