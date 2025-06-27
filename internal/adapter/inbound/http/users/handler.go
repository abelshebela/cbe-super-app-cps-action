package users

import (
	"net/http"

	route "cbe-super-app-member-users/internal/adapter/inbound/http"
	"cbe-super-app-member-users/internal/application/middleware"
	users_application "cbe-super-app-member-users/internal/application/users"
	users_inbound "cbe-super-app-member-users/internal/port/inbound/users"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UsersAdapter struct {
	Application users_application.ApplicationService
	logger      utils.Logger
}

func InitUsersAdapter(app users_application.ApplicationService, logger utils.Logger) users_inbound.InBound {
	return UsersAdapter{
		Application: app,
		logger:      logger,
	}
}

func InitUserRoutes(router chi.Router, handler users_inbound.InBound, authMiddleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/user/linked-accounts",
			Handler: handler.FetchLinkedAccounts,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/email/generate-otp",
			Handler: handler.GenerateEmailOTP,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/email/verify-otp",
			Handler: handler.VerifyEmailOTP,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/set-profile-picture",
			Handler: handler.UpdateProfilePicture,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/unlink-device",
			Handler: handler.UnlinkDevice,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/change-pin",
			Handler: handler.ChangePin,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/sms/verify-otp",
			Handler: handler.VerifyOtp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/user/set-pin",
			Handler: handler.SetPin,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/cbesuperapp/user/register",
			Handler:     handler.Register,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/cbesuperapp/user/login",
			Handler:     handler.Login,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/cbesuperapp/user/forget-pin/send-otp",
			Handler:     handler.ForgetPinSendOtp,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/v1/cbesuperapp/user/device-lookup",
			Handler:     handler.DeviceLookup,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
	}

	route.RegisterRoutes(router, routes)
}
