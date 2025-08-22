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
	var req feedback.FeedbackRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.logger.Errorf("failed to bind feedback data: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
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

	_, err := f.feedbackApplication.CreateFeedback(r.Context(), req, userID)
	if err != nil {
		f.logger.Errorf("failed to create feedback: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	localization.SendSuccessResponse(w, localization.SuccessFeedbackCreated, nil)
}

func (f *feedbackAdapter) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
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

	feedbacks, err := f.feedbackApplication.GetFeedbacks(r.Context(), filterParams)
	if err != nil {
		f.logger.Errorf("failed to get feedbacks: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedbacks)
}

func (f *feedbackAdapter) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Enhanced ID validation
	if id == "" {
		f.logger.Errorf("empty feedback ID provided")
		localization.SendErrorByCodeResponse(w, "feedback ID is required")
		return
	}

	if len(id) != 24 {
		f.logger.Errorf("invalid feedback ID format: %s", id)
		localization.SendErrorByCodeResponse(w, "invalid feedback ID format")
		return
	}

	ctx := r.Context()
	feedback, err := f.feedbackApplication.GetFeedbackByID(ctx, id)
	if err != nil {
		f.logger.Errorf("failed to get feedback by ID %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedback)
}
