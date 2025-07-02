package customerhandler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/feedback"

	"github.com/go-chi/chi/v5"
)

func InitFeedbackRoutes(router chi.Router, feedbackHandler inbound.Feedback) {
	router.Route("/api/v1/cbesuperapp/cps_action/feedback", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: feedbackHandler.GetFeedbacks,
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: feedbackHandler.GetFeedbackByID,
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
