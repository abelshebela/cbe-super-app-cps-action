package users

import (
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/glue"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/middleware"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler rest.Users, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		// {
		// 	Method:      http.MethodPost,
		// 	Path:        "/register",
		// 	Handler:     handler.Register,
		// 	Middlewares: []func(next http.Handler) http.Handler{},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/sms/verify_otp",
		// 	Handler: handler.VerifyOtp,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateTempToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/set_pin",
		// 	Handler: handler.SetPin,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateTempToken,
		// 	},
		// },
		{
			Method:      http.MethodGet,
			Path:        "/device_lookup",
			Handler:     handler.DeviceLookup,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/login",
		// 	Handler: handler.Login,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateTempToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/change_pin",
		// 	Handler: handler.ChangePin,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/forget_pin",
		// 	Handler: handler.ForgetPinSendOtp,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateTempToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/forget_pin/verify_otp",
		// 	Handler: handler.VerifyForgetPinOtp,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateTempToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/forget_pin/reset",
		// 	Handler: handler.ResetPin,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateTempToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/set_profile_picture",
		// 	Handler: handler.UpdateProfilePicture,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/set_profile_theme",
		// 	Handler: handler.UpdateProfileTheme,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		{
			Method:      http.MethodPost,
			Path:        "/pre_login",
			Handler:     handler.PreLogin,
			Middlewares: []func(next http.Handler) http.Handler{},
		},
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/check_pin",
		// 	Handler: handler.CheckPin,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		{
			Method:  http.MethodGet,
			Path:    "/healthcheck",
			Handler: handler.Healthcheck,
		},
	}

	glue.RegisterRoutes(router, routes)
}
