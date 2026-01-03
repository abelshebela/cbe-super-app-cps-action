package feedback

import (
	"net/http"

	feedback "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, feedbackHandler feedback.FeedbackAdapter, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:      http.MethodPost,
			Path:        "/feedback-surveys/create",
			Handler:     feedbackHandler.CreateFeedback,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:  http.MethodGet,
			Path:    "/feedback-surveys",
			Handler: feedbackHandler.GetFeedbacks,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/feedback-surveys/{id}",
			Handler: feedbackHandler.GetFeedbackByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customer-feedbacks",
			Handler: feedbackHandler.GetAllCustomerFeedbacks,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customer-feedbacks/{id}",
			Handler: feedbackHandler.GetCustomerFeedback,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
