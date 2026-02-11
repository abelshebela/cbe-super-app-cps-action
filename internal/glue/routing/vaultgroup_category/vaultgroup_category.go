package vaultgroupcategory

import (
	bankgroupcategory "cbe-super-app-cps-action/internal/constants/interfaces/vaultgroup_category"

	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler bankgroupcategory.VaultGroupCategoryHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/vaultgroupcategory/create",
			Handler: handler.CreateVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vaultgroupcategory",
			Handler: handler.FindAllVaultGroupCategories,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vaultgroupcategory/{id}",
			Handler: handler.GetVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vaultgroupcategory/update/{id}",
			Handler: handler.UpdateVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/vaultgroupcategory/delete/{id}",
			Handler: handler.DeleteVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vaultgroupcategory/enable/{id}",
			Handler: handler.EnableVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vaultgroupcategory/disable/{id}",
			Handler: handler.DisableVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)

}
