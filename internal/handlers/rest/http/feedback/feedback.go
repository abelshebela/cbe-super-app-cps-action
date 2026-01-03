package feedback

import (
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type feedbackAdapter struct {
	logger              utils.Logger
	feedbackApplication service.FeedbackService
}

func InitFeedbackAdapter(feedbackApplication service.FeedbackService, logger utils.Logger) *feedbackAdapter {
	return &feedbackAdapter{
		logger:              logger,
		feedbackApplication: feedbackApplication,
	}
}

func (f *feedbackAdapter) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createFeedback", "handler", "feedback")
	defer span.End()
	var req feedback.FeedbackRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		f.logger.Errorf("failed to bind feedback data: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		f.logger.Errorf("validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	// Extract user ID from context or use a default for anonymous feedback
	userID := "anonymous"
	if userIDFromContext := r.Context().Value("user_id"); userIDFromContext != nil {
		if id, ok := userIDFromContext.(string); ok {
			userID = id
		}
	}

	span.SetAttributes(
		attribute.String("feedback.user_id", userID),
	)
	_, err := f.feedbackApplication.CreateFeedback(ctx, req, userID)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[CreateFeedback] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	f.logger.Infof("[CreateFeedback] feedback created successfully by user: %s", userID)
	localization.SendSuccessResponse(w, localization.SuccessFeedbackCreated, nil)
}

// GetFeedbacks godoc
//
//	@Summary		Get list of feedbacks
//	@Description	Fetch feedbacks with pagination, filtering, and search. Filterable fields: created_at, user_id, responses. Searchable field: responses.
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int							false	"Page number (default 1)"
//	@Param			per_page	query		int							false	"Items per page (default 10, max 100)"
//	@Param			sort		query		string						false	"Sort field"
//	@Param			order		query		string						false	"Sort order (asc/desc)"
//	@Param			created_at	query		string						false	"Filter by created_at"
//	@Param			user_id		query		string						false	"Filter by user_id"
//	@Param			responses	query		string						false	"Filter by responses"
//	@Param			search		query		string						false	"Search term (searches responses)"
//	@Success		200			{object}	localization.ResponseCode	"Feedbacks fetched successfully"
//	@Failure		400			{object}	localization.ResponseCode	"Invalid query params"
//	@Failure		500			{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/feedback [get]
func (f *feedbackAdapter) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getFeedbacks", "handler", "feedback")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)

	// Enhanced pagination validation
	if filterParams.Page < 1 {
		filterParams.Page = 1
	}
	if filterParams.PerPage < 1 {
		filterParams.PerPage = 10
	}
	if filterParams.PerPage > 100 {
		filterParams.PerPage = 100
	}

	feedbacks, err := f.feedbackApplication.GetFeedbacks(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[GetFeedbacks] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("feedback.count", len(feedbacks.Data)))
	f.logger.Infof("[GetFeedbacks] retrieved %d feedbacks", len(feedbacks.Data))
	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedbacks)
}

// GetFeedbackByID godoc
//
//	@Summary		Get feedback by ID
//	@Description	Fetch a single feedback by its ID
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Feedback ID (24-char hex)"
//	@Success		200	{object}	localization.ResponseCode	"Feedback fetched successfully"
//	@Failure		400	{object}	localization.ResponseCode	"Invalid feedback ID"
//	@Failure		404	{object}	localization.ResponseCode	"Feedback not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/feedback/{id} [get]
func (f *feedbackAdapter) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getFeedbackById", "handler", "feedback")
	defer span.End()
	id := chi.URLParam(r, "id")

	// Enhanced ID validation
	if id == "" {
		f.logger.Errorf("empty feedback ID provided")
		localization.SendErrorByCodeResponse(w, localization.ErrorFeedbackIDRequired.Code)
		return
	}

	if len(id) != 24 {
		f.logger.Errorf("invalid feedback ID format: %s", id)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidIDFormat.Code)
		return
	}

	span.SetAttributes(attribute.String("feedback.id", id))
	feedback, err := f.feedbackApplication.GetFeedbackByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[GetFeedbackByID] service error for id %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	f.logger.Infof("[GetFeedbackByID] feedback retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedback)
}

// GetAllCustomerFeedbacks godoc
//
//	@Summary		Get all customer feedbacks
//	@Description	Fetch all customer feedbacks with pagination and search
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int							false	"Page number"
//	@Param			per_page	query		int							false	"Items per page"
//	@Param			search		query		string						false	"Search term"
//	@Success		200			{object}	localization.ResponseCode	"Customer feedbacks fetched successfully"
//	@Failure		500			{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customer-feedbacks [get]
func (f *feedbackAdapter) GetAllCustomerFeedbacks(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllCustomerFeedbacks", "handler", "feedback")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)

	feedbacks, err := f.feedbackApplication.GetAllCustomerFeedbacks(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[GetAllCustomerFeedbacks] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("feedback.count", len(feedbacks.Data)))
	f.logger.Infof("[GetAllCustomerFeedbacks] retrieved %d customer feedbacks", len(feedbacks.Data))
	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedbacks)
}

// GetCustomerFeedback godoc
//
//	@Summary		Get customer feedback by ID
//	@Description	Fetch a single customer feedback by its ID
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Feedback ID"
//	@Success		200	{object}	localization.ResponseCode	"Customer feedback fetched successfully"
//	@Failure		404	{object}	localization.ResponseCode	"Customer feedback not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customer-feedbacks/{id} [get]
func (f *feedbackAdapter) GetCustomerFeedback(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCustomerFeedback", "handler", "feedback")
	defer span.End()
	id := chi.URLParam(r, "id")

	if id == "" {
		f.logger.Errorf("empty feedback ID provided")
		localization.SendErrorByCodeResponse(w, localization.ErrorFeedbackIDRequired.Code)
		return
	}

	span.SetAttributes(attribute.String("feedback.id", id))
	feedback, err := f.feedbackApplication.GetCustomerFeedback(ctx, id)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[GetCustomerFeedback] service error for id %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	f.logger.Infof("[GetCustomerFeedback] customer feedback retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedback)
}
