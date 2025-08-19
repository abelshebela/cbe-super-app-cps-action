package feedback

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type feedbackService struct {
	repo   storage.FeedbackRepository
	logger utils.Logger
}

func NewFeedbackService(repo storage.FeedbackRepository, logger utils.Logger) service.FeedbackService {
	return &feedbackService{
		repo:   repo,
		logger: logger,
	}
}

func (f *feedbackService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	f.logger.Infof("Feedback service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
