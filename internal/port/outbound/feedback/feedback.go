package outbound

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/feedback/entity"
	constant "cbe-super-app-cps-action/utils"
)

type FeedbackRepository interface {
	GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*entity.FeedbackResponse, error)
	GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error)
}
