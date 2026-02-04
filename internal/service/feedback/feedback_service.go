package feedback

import (
	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/feedback/core"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type feedbackService struct {
	repo       storage.FeedbackRepository
	memberRepo storage.CustomerRepository
	logger     utils.Logger
}

// GetSurveyFeedback implements [service.FeedbackService].
func (f *feedbackService) GetSurveyFeedback(ctx context.Context, id string) (*local_model.SurveyFeedback, error) {
	panic("unimplemented")
}

func NewFeedbackService(repo storage.FeedbackRepository, memberRepo storage.CustomerRepository, logger utils.Logger) service.FeedbackService {
	return &feedbackService{
		repo:       repo,
		memberRepo: memberRepo,
		logger:     logger,
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

func (f *feedbackService) CreateFeedback(ctx context.Context, req fbdto.FeedbackRequest, userCode string) (*local_model.Feedback, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateFeedback", "Feedback", "CreateFeedback")
	defer span.End()

	// Validate the request
	if err := req.Validate(); err != nil {
		f.logger.Errorf("[CreateFeedback] invalid feedback request: %v", err)
		span.AddEvent("Invalid feedback request", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userCode),
		))
		return nil, err
	}

	user, err := f.memberRepo.FindCustomerByID(ctx, userCode)
	if err != nil {
		f.logger.Errorf("[CreateFeedback] failed to find user by code %s: %v", userCode, err)
		span.AddEvent("Failed to find user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userCode),
		))
		return nil, errors.New("user not found")
	}
	feedback := core.BuildFeedbackEntity(userCode, req, user)

	// Call the Create method with the Feedback object
	err = f.repo.Create(ctx, feedback)
	if err != nil {
		f.logger.Errorf("[CreateFeedback] failed to create feedback: %v", err)
		span.AddEvent("Failed to create feedback", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userCode),
		))
		return nil, err
	}

	f.logger.Infof("[CreateFeedback] feedback created successfully")
	return feedback, nil
}

// CreateSurveyFeedback implements [service.FeedbackService].
func (f *feedbackService) CreateSurveyFeedback(ctx context.Context, surveyFeedback fbdto.SurveyFeedbackReq) (*local_model.SurveyFeedback, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateSurveyFeedback", "Feedback", "CreateSurveyFeedback")
	defer span.End()

	// Validate the request
	if err := surveyFeedback.Validate(); err != nil {
		f.logger.Errorf("[CreateSurveyFeedback] invalid survey feedback request: %v", err)
		span.AddEvent("Invalid feedback request", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", surveyFeedback.UserID),
		))
		return nil, err
	}
	user, err := f.memberRepo.FindCustomerByID(ctx, surveyFeedback.UserID)
	if err != nil {
		f.logger.Errorf("[CreateFeedback] failed to find user by code %s: %v", surveyFeedback.UserID, err)
		span.AddEvent("Failed to find user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", surveyFeedback.UserID),
		))
		return nil, errors.New("user not found")
	}

	surveyFeedbackEntity := core.BuildSurveyFeedbackEntity(surveyFeedback, user)

	// Call the Create method with the Feedback object
	err = f.repo.CreateSurveyFeedback(ctx, surveyFeedbackEntity)
	if err != nil {
		f.logger.Errorf("[CreateSurveyFeedback] failed to create survey feedback: %v", err)
		span.AddEvent("Failed to create survey feedback", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", surveyFeedback.UserID),
		))
		return nil, err
	}

	f.logger.Infof("[CreateSurveyFeedback] survey feedback created successfully")
	return surveyFeedbackEntity, nil
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
func (f *feedbackService) GetAllCustomerFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]local_model.CustomerFeedback], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllCustomerFeedbacks", "Feedback", "GetAllCustomerFeedbacks")
	defer span.End()

	feedbacks, err := f.repo.FindAllCustomerFeedbacks(ctx, *filterParams)
	if err != nil {
		f.logger.Errorf("[GetAllCustomerFeedbacks] failed to fetch customer feedbacks: %v", err)
		span.RecordError(err)
		return nil, err
	}

	f.logger.Infof("[GetAllCustomerFeedbacks] retrieved %d customer feedbacks", len(feedbacks.Data))
	return feedbacks, nil
}
func (f *feedbackService) GetAllSurveyFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]local_model.SurveyFeedback], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllSurveyFeedbacks", "Feedback", "GetAllSurveyFeedbacks")
	defer span.End()

	feedbacks, err := f.repo.FindAllSurveyFeedbacks(ctx, *filterParams)
	if err != nil {
		f.logger.Errorf("[GetAllSurveyFeedbacks] failed to fetch survey feedbacks: %v", err)
		span.RecordError(err)
		return nil, err
	}

	f.logger.Infof("[GetAllSurveyFeedbacks] retrieved %d survey feedbacks", len(feedbacks.Data))
	return feedbacks, nil
}

func (f *feedbackService) GetCustomerFeedback(ctx context.Context, id string) (*local_model.CustomerFeedback, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCustomerFeedback", "Feedback", "GetCustomerFeedback")
	defer span.End()

	feedback, err := f.repo.FindCustomerFeedbackByID(ctx, id)
	if err != nil {
		f.logger.Errorf("[GetCustomerFeedback] failed to get customer feedback for id %s: %v", id, err)
		span.RecordError(err)
		return nil, err
	}

	f.logger.Infof("[GetCustomerFeedback] customer feedback retrieved successfully for id: %s", id)
	return feedback, nil
}
