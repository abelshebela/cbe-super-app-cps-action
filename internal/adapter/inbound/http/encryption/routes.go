package encryption

import (
	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	encryption "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/encryption"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitEncryptionRoutes(router chi.Router, encryptionHandler encryption.Encryption, authMiddleware middleware.AuthMiddleware) {
	router.Route("/encryption", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/encrypt",
				Handler: encryptionHandler.Encrypt,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})
}
