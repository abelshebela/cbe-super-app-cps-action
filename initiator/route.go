package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	accountblock "cbe-super-app-cps-action/internal/glue/routing/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/glue/routing/account_validation"
	advert "cbe-super-app-cps-action/internal/glue/routing/ad"
	amountBasedAuth "cbe-super-app-cps-action/internal/glue/routing/amount_based_auth"
	avatar "cbe-super-app-cps-action/internal/glue/routing/avatar"
	"cbe-super-app-cps-action/internal/glue/routing/bank"

	// bankvaultroutes "cbe-super-app-cps-action/internal/glue/routing/bankvault"
	// vaultgroupcategory "cbe-super-app-cps-action/internal/glue/routing/vaultgroup_category"
	bpsUser "cbe-super-app-cps-action/internal/glue/routing/bps_user"
	budget "cbe-super-app-cps-action/internal/glue/routing/budget"
	"cbe-super-app-cps-action/internal/glue/routing/bulk_service"
	cpsaction "cbe-super-app-cps-action/internal/glue/routing/cps_action"
	"cbe-super-app-cps-action/internal/glue/routing/customer"
	"cbe-super-app-cps-action/internal/glue/routing/department"
	eventhandler "cbe-super-app-cps-action/internal/glue/routing/event"
	miniapp "cbe-super-app-cps-action/internal/glue/routing/mini_app"
	miniappmerchant "cbe-super-app-cps-action/internal/glue/routing/mini_app_merchant"
	"cbe-super-app-cps-action/internal/glue/routing/notification"
	"cbe-super-app-cps-action/internal/glue/routing/wallet"

	cps_user_det "cbe-super-app-cps-action/internal/glue/routing/cps_user"
	fayda "cbe-super-app-cps-action/internal/glue/routing/fayda"
	feedback "cbe-super-app-cps-action/internal/glue/routing/feedback"
	hqRoute "cbe-super-app-cps-action/internal/glue/routing/hq"
	password "cbe-super-app-cps-action/internal/glue/routing/password_rule"
	permission_details "cbe-super-app-cps-action/internal/glue/routing/permission"
	portalcard "cbe-super-app-cps-action/internal/glue/routing/portal_card"
	service_details "cbe-super-app-cps-action/internal/glue/routing/service_details"

	donation "cbe-super-app-cps-action/internal/glue/routing/donation"
	donation_category "cbe-super-app-cps-action/internal/glue/routing/donation_category"
	donation_company "cbe-super-app-cps-action/internal/glue/routing/donation_company"
	productcode "cbe-super-app-cps-action/internal/glue/routing/product_code"
	unlink "cbe-super-app-cps-action/internal/glue/routing/unlink"
	customeMiddleware "cbe-super-app-cps-action/internal/handlers/middleware"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "cbe-super-app-cps-action/docs" // Import generated docs
)

func InitRoute(ctx context.Context, router *chi.Mux, handlerLayer Handler, logger utils.Logger, cfg *config.VaultConfig) {

	r := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(customeMiddleware.ChiLogger(logger))

	router.Use(customeMiddleware.HandlePanic(logger))
	router.Use(customeMiddleware.CORS())
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.Compress(5, "application/json"))

	r.Get("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "CPS ACTION IS ACTIVE"}); err != nil {
			logger.Errorf("Failed to write health check response", zap.Error(err))
		}
	})
	authMiddleware := customeMiddleware.InitAuthMiddleware(cfg.JwtSecretKey, cfg.Key, cfg.IV, logger)

	cpsaction.Init(r, handlerLayer.CpsActionHandler, authMiddleware)
	budget.Init(r, handlerLayer.BudgetHandler, authMiddleware)
	avatar.Init(r, handlerLayer.AvatarHandler, authMiddleware)
	unlink.Init(r, handlerLayer.UnlinkHandler, authMiddleware)
	bpsUser.Init(r, handlerLayer.BpsHandler, authMiddleware)
	bank.Init(r, handlerLayer.BankHandler, authMiddleware)
	eventhandler.Init(r, handlerLayer.EventHandler, authMiddleware)
	wallet.Init(r, handlerLayer.WalletHandler, authMiddleware)

	customer.Init(r, handlerLayer.customerHandler, authMiddleware)
	bulk_service.Init(r, handlerLayer.bulkServiceHandler, authMiddleware)

	password.Init(r, handlerLayer.PasswordHandler, authMiddleware)

	feedback.Init(r, handlerLayer.FeedbackHandler, authMiddleware)
	productcode.Init(r, &handlerLayer.ProductCodeHandler, authMiddleware)
	advert.Init(r, handlerLayer.AdvertHandler, authMiddleware)
	portalcard.Init(r, handlerLayer.PortalCardHander, authMiddleware)
	accountvalidation.Init(r, handlerLayer.AccountValidation, authMiddleware)
	miniappmerchant.Init(r, handlerLayer.MiniAppMerchantHandler, authMiddleware)
	accountblock.Init(r, handlerLayer.AccountBlockHandler, authMiddleware)
	department.Init(r, &handlerLayer.DepartmentHandler, authMiddleware)
	hqRoute.Init(r, handlerLayer.HqHandler, authMiddleware)
	miniapp.Init(r, handlerLayer.MiniAPPHandler, authMiddleware)
	fayda.Init(r, handlerLayer.FaydaHandler, authMiddleware)
	service_details.Init(r, handlerLayer.ServiceDetailsHandler, authMiddleware)
	permission_details.Init(r, handlerLayer.Permission, authMiddleware)
	cps_user_det.Init(r, handlerLayer.CPSUser, authMiddleware)
	// bankvaultroutes.Init(r, handlerLayer.BankVaultHandler, authMiddleware)
	// vaultgroupcategory.Init(r, handlerLayer.VaultGroupCategoryHandler, authMiddleware)

	donation.Init(r, handlerLayer.DonationHandler, authMiddleware)
	donation_category.Init(r, handlerLayer.DonationCategoryHandler, authMiddleware)
	donation_company.Init(r, handlerLayer.DonationCompanyHandler, authMiddleware)

	amountBasedAuth.Init(r, handlerLayer.AmountBasedAuthHandler, authMiddleware)
	notification.Init(r, handlerLayer.NotificationHandler, authMiddleware)

	router.Mount("/api/v1/cbesuperapp/cps_action", r)

	// Swagger documentation routes
	router.Get("/api/v1/cbesuperapp/cps_action/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/cbesuperapp/cps_action/swagger/doc.json"),
	))
	router.Get("/api/v1/cbesuperapp/cps_action/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/cbesuperapp/cps_action/swagger/index.html", http.StatusMovedPermanently)
	})
}
