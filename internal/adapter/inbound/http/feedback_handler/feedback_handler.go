package customerhandler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/feedback"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/feedback"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
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

func (f FeedbackHTTPHandler) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	per_page := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt <= 10 && perPageInt > 0 {
		per_page = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: per_page,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()
	feedbacks, err := f.feedbackService.GetFeedbacks(ctx, filterParams)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           feedbacks,
	}
	res.SendJSON()
}

func (f FeedbackHTTPHandler) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	feedback, err := f.feedbackService.GetFeedbackByID(ctx, id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           feedback,
	}
	res.SendJSON()
}
