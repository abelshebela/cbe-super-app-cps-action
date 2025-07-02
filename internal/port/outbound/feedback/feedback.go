package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type FeedbackRepository interface {
	GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*entity.FeedbackResponse, error)
	GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error)
}
