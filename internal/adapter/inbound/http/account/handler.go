package account

import (
	"net/http"

	route "cbe-super-app-member-users/internal/adapter/inbound/http"
	account_application "cbe-super-app-member-users/internal/application/account"
	"cbe-super-app-member-users/internal/application/middleware"
	account_inbound "cbe-super-app-member-users/internal/port/inbound/account"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountAdapter struct {
	Application account_application.ApplicationService
	logger      utils.Logger
}

func InitAccountAdapter(app account_application.ApplicationService, logger utils.Logger) account_inbound.InBound {
	return AccountAdapter{
		Application: app,
		logger:      logger,
	}
}

func InitAccountRoutes(router chi.Router, handler account_inbound.InBound, authMiddleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/add-account",
			Handler: handler.CreateAccount,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
