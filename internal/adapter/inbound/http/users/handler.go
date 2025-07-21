package users

import (
	"net/http"

	route "cbe-super-app-member-users/internal/adapter/inbound/http"
	"cbe-super-app-member-users/internal/application/middleware"
	users_application "cbe-super-app-member-users/internal/application/users"
	users_inbound "cbe-super-app-member-users/internal/port/inbound/users"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UsersAdapter struct {
	Application users_application.ApplicationService
	logger      utils.Logger
	config      *config.VaultConfig
}

func InitUsersAdapter(app users_application.ApplicationService, logger utils.Logger, config *config.VaultConfig) users_inbound.InBound {
	return UsersAdapter{
		Application: app,
		logger:      logger,
		config:      config,
	}
}

func InitUserRoutes(router chi.Router, handler users_inbound.InBound, authMiddleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/member_auth/linked-accounts",
			Handler: handler.FetchLinkedAccounts,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/email/generate-otp",
			Handler: handler.GenerateEmailOTP,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/email/verify-otp",
			Handler: handler.VerifyEmailOTP,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/set-profile-picture",
			Handler: handler.UpdateProfilePicture,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/set-profile-theme",
			Handler: handler.UpdateProfileTheme,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/unlink-device",
			Handler: handler.UnlinkDevice,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/change-pin",
			Handler: handler.ChangePin,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/sms/verify-otp",
			Handler: handler.VerifyOtp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/set-pin",
			Handler: handler.SetPin,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/cbesuperapp/member_auth/register",
			Handler:     handler.Register,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/register/complete",
			Handler: handler.CompleteRegistration,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/login",
			Handler: handler.Login,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/forget-pin",
			Handler: handler.ForgetPinSendOtp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/forget-pin/verify-otp",
			Handler: handler.VerifyForgetPinOtp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/forget-pin/reset",
			Handler: handler.ResetPin,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/cbesuperapp/member_auth/reset-pin-with-token",
			Handler:     handler.ResetPinWithToken,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/v1/cbesuperapp/member_auth/device-lookup",
			Handler:     handler.DeviceLookup,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/v1/cbesuperapp/member_auth/pre-login",
			Handler:     handler.PreLogin,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/member_auth/check-pin",
			Handler: handler.CheckPin,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/member_auth/healthcheck",
			Handler: handler.Healthcheck,
		},
	}

	route.RegisterRoutes(router, routes)
}
