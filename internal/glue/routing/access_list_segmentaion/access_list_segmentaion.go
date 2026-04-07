package access_list_segmentation

import (
	accesslistsegmentation "cbe-super-app-cps-action/internal/constants/interfaces/access_list_segmentation"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler accesslistsegmentation.AccessListSegmentationHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/access_list_segmentation",
			Handler: handler.CreateAccessListSegmentation,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/access_list_segmentation",
			Handler: handler.GetAllAccessListSegmentation,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/access_list_segmentation/{id}",
			Handler: handler.UpdateAccessListSegmentation,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/access_list_segmentation/{id}",
			Handler: handler.GetAccessListSegmentationByID,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/access_list_segmentation/all/{id}",
			Handler: handler.GetAllAccessListSegmentationBySegmentIDorSegmentCode,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/access_list_segmentation/enable/{id}",
			Handler: handler.EnableAccessListSegmentation,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/access_list_segmentation/disable/{id}",
			Handler: handler.DisableAccessListSegmentation,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
