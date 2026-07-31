package newstag_routing

import (
	newstag_adaptor "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/news_tag"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler newstag_adaptor.NewsTagAdaptor, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/news/tags",
			Handler: handler.FetchNewsTags,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/news/tags/create",
			Handler: handler.CreateNewsTags,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/news/tags/{id}",
			Handler: handler.GetNewsTagByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/news/tags/{id}",
			Handler: handler.DeleteNewsTag,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/news/tags/{id}",
			Handler: handler.UpdateNewsTag,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
