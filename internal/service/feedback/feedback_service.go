package feedback

import (
	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/feedback/core"
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

func (f *feedbackService) CreateFeedback(ctx context.Context, req fbdto.FeedbackRequest, userID string) (*model.Feedback, error) {

	// Validate the request
	if err := req.Validate(); err != nil {
		f.logger.Errorf("Invalid feedback request: %v", err)
		return nil, err
	}

	feedback := core.BuildFeedbackEntity(userID, req)

	// Call the Create method with the Feedback object
	err := f.repo.Create(ctx, feedback)
	if err != nil {
		f.logger.Errorf("Failed to create feedback: %v", err)
		return nil, err
	}

	f.logger.Infof("Feedback created successfully for user: %s", userID)
	return feedback, nil
}

func (f *feedbackService) GetFeedbackByID(ctx context.Context, id string) (*fbdto.FeedbackResponse, error) {
	feedback, err := f.repo.FindByID(ctx, id)
	if err != nil {
		f.logger.Errorf("Failed to get feedback by ID: %v", err)
		return nil, err
	}
	return feedback, nil
}

func (f *feedbackService) GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponseForFeedback[[]*fbdto.FeedbackResponse], error) {
	feedbacks, err := f.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		f.logger.Errorf("Failed to get feedbacks: %v", err)
		return nil, err
	}

	f.logger.Infof("Retrieved %d feedbacks", len(feedbacks.Data))
	return feedbacks, nil
}
