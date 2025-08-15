package branch_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func RegisterBranchRoutes(r chi.Router, handler inbound.BranchHandler, authMiddleware middleware.AuthMiddleware) {
	r.Route("/branches", func(rou chi.Router) {
		routes := []sharedhttp.Route{
			{
				Method:  http.MethodPost,
				Path:    "/filter/single",
				Handler: handler.GetBranch,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/filter/multiple",
				Handler: handler.GetAllBranches,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/single",
				Handler: handler.ApproveSingleBranchDisable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/multiple",
				Handler: handler.DisableMultipleBranches,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/single/approve",
				Handler: handler.DisableMultipleBranches,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/multiple/approve",
				Handler: handler.ApproveBulkBranchesDisable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
		}

		sharedhttp.RegisterRoutes(rou, routes)
	})
}
