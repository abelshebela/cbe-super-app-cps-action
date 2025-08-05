package customerhandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/feedback"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/feedback"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type FeedbackHTTPHandler struct {
	feedbackService feedback.FeedbackService
	logger          utils.Logger
}

func NewFeedbackHTTPHandler(feedbackService feedback.FeedbackService, logger utils.Logger) inbound.Feedback {
	return FeedbackHTTPHandler{
		feedbackService: feedbackService,
		logger:          logger,
	}
}

func (f FeedbackHTTPHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var req entity.FeedbackRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.logger.Errorf("failed to bind feedback data: %v", err)
		util.SendErrorResponse(w, "Invalid request data", 400, nil)
		return
	}

	if err := req.Validate(); err != nil {
		util.SendErrorResponse(w, err.Error(), 400, nil)
		return
	}

	// Extract user ID from context or use a default for anonymous feedback
	userID := "anonymous"
	if userIDFromContext := r.Context().Value("user_id"); userIDFromContext != nil {
		if id, ok := userIDFromContext.(string); ok {
			userID = id
		}
	}

	feedback, err := f.feedbackService.CreateFeedback(r.Context(), req, userID)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	util.WriteSuccessResponse(w, feedback, "Feedback created successfully")
}

func (f FeedbackHTTPHandler) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	ctx := r.Context()
	feedbacks, err := f.feedbackService.GetFeedbacks(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	util.WriteSuccessResponse(w, feedbacks, "successfully fetched the feedback")

}

func (f FeedbackHTTPHandler) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	feedback, err := f.feedbackService.GetFeedbackByID(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, feedback, "successfully fetched the feedback")
}
