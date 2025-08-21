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
		localization.SendErrorByCodeResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
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
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessFeedbackCreated, nil)
	// data := interface{}{}
	localization.SendSuccessResponse(w, localization.SuccessFeedbackCreated, nil)
}

func (f *feedbackAdapter) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	feedbacks, err := f.feedbackApplication.GetFeedbacks(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())

		return
	}
	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedbacks)
}

func (f *feedbackAdapter) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	feedback, err := f.feedbackApplication.GetFeedbackByID(ctx, id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessFeedbackFetched, feedback)
}
