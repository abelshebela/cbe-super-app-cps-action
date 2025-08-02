package initiator

import (
	"cbe-super-app-budget/internal/glue/routing"
	customeMiddleware "cbe-super-app-budget/internal/handlers/middleware"

	"cbe-super-app-budget/platform/logger"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func InitRoute(ctx context.Context, router *chi.Mux, handlerLayer HandlerLayer, logger logger.Logger) {
	r := chi.NewRouter()

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

	routing.InitSpending(r, handlerLayer.spending, logger)

	router.Mount("/api/v1/cbesuperapp/member", r)
}
