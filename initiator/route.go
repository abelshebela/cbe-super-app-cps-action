package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	accountblock "cbe-super-app-cps-action/internal/glue/routing/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/glue/routing/account_validation"
	advert "cbe-super-app-cps-action/internal/glue/routing/ad"
	"cbe-super-app-cps-action/internal/glue/routing/bank"
	bpsUser "cbe-super-app-cps-action/internal/glue/routing/bps_user"
	budget "cbe-super-app-cps-action/internal/glue/routing/budget"
	"cbe-super-app-cps-action/internal/glue/routing/bulk_service"
	cpsaction "cbe-super-app-cps-action/internal/glue/routing/cps_action"
	"cbe-super-app-cps-action/internal/glue/routing/customer"
	"cbe-super-app-cps-action/internal/glue/routing/department"
	eventhandler "cbe-super-app-cps-action/internal/glue/routing/event"
	fayda "cbe-super-app-cps-action/internal/glue/routing/fayda"
	feedback "cbe-super-app-cps-action/internal/glue/routing/feedback"
	hqRoute "cbe-super-app-cps-action/internal/glue/routing/hq"
	password "cbe-super-app-cps-action/internal/glue/routing/password_rule"
	portalcard "cbe-super-app-cps-action/internal/glue/routing/portal_card"
	unlink "cbe-super-app-cps-action/internal/glue/routing/unlink"
	"cbe-super-app-cps-action/internal/glue/routing/wallet"
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
	budget.Init(r, handlerLayer.BudgetHandler, authMiddleware)
	unlink.Init(r, handlerLayer.UnlinkHandler, authMiddleware)
	bpsUser.Init(r, handlerLayer.BpsHandler, authMiddleware)
	bank.Init(r, handlerLayer.BankHandler, authMiddleware)
	eventhandler.Init(r, handlerLayer.EventHandler, authMiddleware)
	wallet.Init(r, handlerLayer.WalletHandler, authMiddleware)

	customer.Init(r, handlerLayer.customerHandler, authMiddleware)
	bulk_service.Init(r, handlerLayer.bulkServiceHandler, authMiddleware)

	password.Init(r, handlerLayer.PasswordHandler, authMiddleware)

	feedback.Init(r, handlerLayer.FeedbackHandler, authMiddleware)
	advert.Init(r, handlerLayer.AdvertHandler, authMiddleware)
	portalcard.Init(r, handlerLayer.PortalCardHander, authMiddleware)
	accountvalidation.Init(r, handlerLayer.AccountValidation, authMiddleware)
	accountblock.Init(r, handlerLayer.AccountBlockHandler, authMiddleware)
	department.Init(r, &handlerLayer.DepartmentHandler, authMiddleware)
	hqRoute.Init(r, handlerLayer.HqHandler, authMiddleware)
	fayda.Init(r, handlerLayer.FaydaHandler, authMiddleware)

	router.Mount("/api/v1/cbesuperapp/cps_action", r)
}
