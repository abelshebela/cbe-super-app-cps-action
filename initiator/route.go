package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	bpsUser "cbe-super-app-cps-action/internal/glue/routing/bps_user"
	cpsaction "cbe-super-app-cps-action/internal/glue/routing/cps_action"
	feedback "cbe-super-app-cps-action/internal/glue/routing/feedback"
	productcode "cbe-super-app-cps-action/internal/glue/routing/product_code"
	unlink "cbe-super-app-cps-action/internal/glue/routing/unlink"
	customeMiddleware "cbe-super-app-cps-action/internal/handlers/middleware"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func InitRoute(ctx context.Context, router *chi.Mux, handlerLayer Handler, logger utils.Logger) {
	r := chi.NewRouter()
	cfg, _ := config.Load()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(customeMiddleware.ChiLogger(logger))

	router.Use(customeMiddleware.HandlePanic(logger))
	router.Use(customeMiddleware.CORS())
	router.Use(middleware.Timeout(30 * time.Second))

	r.Get("/api/v1/cbesuperapp/cps_action/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "CPS ACTION IS ACTIVE"}); err != nil {
			logger.Errorf("Failed to write health check response", zap.Error(err))
		}
	})
	authMiddleware := customeMiddleware.InitAuthMiddleware(cfg.JwtSecretKey, cfg.Key, cfg.IV, logger)

	cpsaction.Init(r, handlerLayer.CpsActionHandler, authMiddleware)
	unlink.Init(r, handlerLayer.UnlinkHandler, authMiddleware)
	bpsUser.Init(r, handlerLayer.BpsHandler, authMiddleware)
	feedback.Init(r, handlerLayer.FeedbackHandler, authMiddleware)
	productcode.Init(r, handlerLayer.ProductCodeHandler, authMiddleware, nil) //todo: add cps action factory
	router.Mount("/api/v1/cbesuperapp/cps_action", r)
}
