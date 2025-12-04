package vaultgroupcategory

import (
	role "cbe-super-app-cps-action/internal/constants"
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
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vaultgroupcategory",
			Handler: handler.FindAllVaultGroupCategories,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vaultgroupcategory/{id}",
			Handler: handler.GetVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vaultgroupcategory/update/{id}",
			Handler: handler.UpdateVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/vaultgroupcategory/delete/{id}",
			Handler: handler.DeleteVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vaultgroupcategory/enable/{id}",
			Handler: handler.EnableVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vaultgroupcategory/disable/{id}",
			Handler: handler.DisableVaultGroupCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}
	glue.RegisterRoutes(router, routes)

}
