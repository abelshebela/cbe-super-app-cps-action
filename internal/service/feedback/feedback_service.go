package feedback

import (
	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/feedback/core"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	local_util "cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Feedback", "Authorize")
	defer span.End()

	f.logger.Infof("[Authorize] authorizing feedback action: %s", cpsAction.RequestAction)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	f.logger.Infof("[Authorize] feedback action authorized successfully")
	return cpsAction, nil
}

func (f *feedbackService) CreateFeedback(ctx context.Context, req fbdto.FeedbackRequest, userID string) (*model.Feedback, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateFeedback", "Feedback", "CreateFeedback")
	defer span.End()

	// Validate the request
	if err := req.Validate(); err != nil {
		f.logger.Errorf("[CreateFeedback] invalid feedback request: %v", err)
		span.AddEvent("Invalid feedback request", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userID),
		))
		return nil, err
	}

	feedback := core.BuildFeedbackEntity(userID, req)

	// Call the Create method with the Feedback object
	err := f.repo.Create(ctx, feedback)
	if err != nil {
		f.logger.Errorf("[CreateFeedback] failed to create feedback: %v", err)
		span.AddEvent("Failed to create feedback", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userID),
		))
		return nil, err
	}

	f.logger.Infof("[CreateFeedback] feedback created successfully")
	return feedback, nil
}

func (f *feedbackService) GetFeedbackByID(ctx context.Context, id string) (*fbdto.FeedbackResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetFeedbackByID", "Feedback", "GetFeedbackByID")
	defer span.End()

	feedback, err := f.repo.FindByID(ctx, id)
	if err != nil {
		f.logger.Errorf("[GetFeedbackByID] failed to get feedback: %v", err)
		span.AddEvent("Failed to get feedback", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	f.logger.Infof("[GetFeedbackByID] feedback retrieved successfully for id: %s", id)
	return feedback, nil
}

func (f *feedbackService) GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponseForFeedback[[]*fbdto.FeedbackResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetFeedbacks", "Feedback", "GetFeedbacks")
	defer span.End()

	feedbacks, err := f.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		f.logger.Errorf("[GetFeedbacks] failed to fetch feedbacks: %v", err)
		span.AddEvent("Failed to fetch feedbacks", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	f.logger.Infof("[GetFeedbacks] retrieved %d feedbacks", len(feedbacks.Data))
	return feedbacks, nil
}
