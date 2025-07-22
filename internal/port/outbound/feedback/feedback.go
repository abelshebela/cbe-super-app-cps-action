package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type FeedbackRepository interface {
	GetFeedbacks(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Feedback], error)
	GetFeedbackByID(ctx context.Context, id string) (*entity.Feedback, error)
}
