package feedback

import (
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"time"

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

func (f *feedbackService) CreateFeedback(ctx context.Context, req feedback.FeedbackRequest, userID string) (*model.Feedback, error) {
	// Create
	now := time.Now()
	feedback := &model.Feedback{
		UserID:    userID,
		Responses: make(map[string]types.Response),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Call the Create method with the Feedback object
	if err := f.repo.Create(ctx, feedback); err != nil {
		f.logger.Errorf("Failed to create feedback: %v", err)
		return nil, err
	}

	f.logger.Infof("Feedback created successfully for user: %s", userID)
	return nil, nil

}

func (f *feedbackService) GetFeedbackByID(ctx context.Context, id string) (*model.Feedback, error) {
	feedback, err := f.repo.FindByID(ctx, id)
	if err != nil {
		f.logger.Errorf("Failed to get feedback by ID: %v", err)
		return nil, err
	}
	return feedback, nil
}

func (f *feedbackService) GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Feedback], error) {
	feedbacks, err := f.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		f.logger.Errorf("Failed to get feedbacks: %v", err)
		return nil, err
	}

	f.logger.Infof("Retrieved %d feedbacks", len(feedbacks.Data))
	return feedbacks, nil
}
