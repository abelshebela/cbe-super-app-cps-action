package encryption

import (
	encryption "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/encryption"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler encryption.EncryptionAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/encryption/encrypt",
			Handler: handler.Encrypt,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken},
		},
	}

	glue.RegisterRoutes(router, routes)

}
