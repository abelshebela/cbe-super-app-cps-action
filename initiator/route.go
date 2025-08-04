package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/glue/routing/users"
	customeMiddleware "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/middleware"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"
	vaultConfig "gitlab.com/bersufekadgetachew/cbe-super-app-shared/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func InitRoute(ctx context.Context, router *chi.Mux, handlerLayer Handler, logger logger.Logger) {
	r := chi.NewRouter()
	cfg := vaultConfig.LoadVault()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(customeMiddleware.ChiLogger(logger))

	router.Use(customeMiddleware.HandlePanic(logger))
	router.Use(customeMiddleware.CORS())
	router.Use(middleware.Timeout(30 * time.Second))

	r.Get("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "It's Working!"}); err != nil {
			logger.Error(ctx, "Failed to write health check response", zap.Error(err))
		}
	})

	users.Init(r, handlerLayer.UserHandler, customeMiddleware.InitAuthMiddleware(cfg.JwtSecretKey, cfg.Key, cfg.IV, logger))

	router.Mount("/api/v1/cbesuperapp/member", r)
}
