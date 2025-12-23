package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	cps_auth "cbe-super-app-cps-action/grpc/auth/proto"
	access_list_segmentation "cbe-super-app-cps-action/internal/glue/routing/access_list_segmentaion"
	accountblock "cbe-super-app-cps-action/internal/glue/routing/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/glue/routing/account_validation"
	advert "cbe-super-app-cps-action/internal/glue/routing/ad"
	amountBasedAuth "cbe-super-app-cps-action/internal/glue/routing/amount_based_auth"
	avatar "cbe-super-app-cps-action/internal/glue/routing/avatar"
	"cbe-super-app-cps-action/internal/glue/routing/bank"
	bps_actionrole_routing "cbe-super-app-cps-action/internal/glue/routing/bps_action_role"
	cps_actionrole_routing "cbe-super-app-cps-action/internal/glue/routing/cps_action_role"
	device_version "cbe-super-app-cps-action/internal/glue/routing/device_version"
	ecommerce_merchant "cbe-super-app-cps-action/internal/glue/routing/ecommerce-merchant"
	event_merchant_routing "cbe-super-app-cps-action/internal/glue/routing/event_merchant"
	kyc_routing "cbe-super-app-cps-action/internal/glue/routing/kyc_verifier"
	newscategory_routing "cbe-super-app-cps-action/internal/glue/routing/news_category"
	newstag_routing "cbe-super-app-cps-action/internal/glue/routing/news_tag"
	"cbe-super-app-cps-action/internal/glue/routing/transaction"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/platform/telemetry"

	bankvaultroutes "cbe-super-app-cps-action/internal/glue/routing/bankvault"
	bpsUser "cbe-super-app-cps-action/internal/glue/routing/bps_user"
	budgetCategory "cbe-super-app-cps-action/internal/glue/routing/budget_category"
	"cbe-super-app-cps-action/internal/glue/routing/bulk_service"
	cpsaction "cbe-super-app-cps-action/internal/glue/routing/cps_action"
	"cbe-super-app-cps-action/internal/glue/routing/customer"
	"cbe-super-app-cps-action/internal/glue/routing/department"
	eventhandler "cbe-super-app-cps-action/internal/glue/routing/event"
	"cbe-super-app-cps-action/internal/glue/routing/notification"
	"cbe-super-app-cps-action/internal/glue/routing/topup"
	vaultgroupcategory "cbe-super-app-cps-action/internal/glue/routing/vaultgroup_category"
	"cbe-super-app-cps-action/internal/glue/routing/wallet"

	cps_user_det "cbe-super-app-cps-action/internal/glue/routing/cps_user"
	customer_seg "cbe-super-app-cps-action/internal/glue/routing/customer_segmentation"
	donation "cbe-super-app-cps-action/internal/glue/routing/donation"
	donation_category "cbe-super-app-cps-action/internal/glue/routing/donation_category"
	donation_company "cbe-super-app-cps-action/internal/glue/routing/donation_company"
	encryption "cbe-super-app-cps-action/internal/glue/routing/encryption"
	fayda "cbe-super-app-cps-action/internal/glue/routing/fayda"
	feedback "cbe-super-app-cps-action/internal/glue/routing/feedback"
	hqRoute "cbe-super-app-cps-action/internal/glue/routing/hq"
	jobRole "cbe-super-app-cps-action/internal/glue/routing/job_roles"
	password "cbe-super-app-cps-action/internal/glue/routing/password_rule"
	permission_details "cbe-super-app-cps-action/internal/glue/routing/permission"
	portalcard "cbe-super-app-cps-action/internal/glue/routing/portal_card"
	productcode "cbe-super-app-cps-action/internal/glue/routing/product_code"
	roles "cbe-super-app-cps-action/internal/glue/routing/roles"
	service_details "cbe-super-app-cps-action/internal/glue/routing/service_details"
	service "cbe-super-app-cps-action/internal/glue/routing/services"
	sitota "cbe-super-app-cps-action/internal/glue/routing/sitota"
	unlink "cbe-super-app-cps-action/internal/glue/routing/unlink"
	vaultAmountTier "cbe-super-app-cps-action/internal/glue/routing/vault_amount_tier"
	customeMiddleware "cbe-super-app-cps-action/internal/handlers/middleware"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "cbe-super-app-cps-action/docs" // Import generated docs
)

func InitRoute(ctx context.Context, router *chi.Mux, handlerLayer Handler, client cps_auth.CpsAuthServiceClient, redisRepository storage.RedisRepository, logger utils.Logger, cfg *config.VaultConfig) {

	r := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	// Inject trace and span ids from OpenTelemetry span into context for logger extraction
	router.Use(telemetry.TraceContextMiddleware())
	// Optional debug middleware to detect missing spans. Enable by setting OTEL_DEBUG_TRACE_PRESENCE=true
	if os.Getenv("OTEL_DEBUG_TRACE_PRESENCE") == "true" {
		router.Use(telemetry.SpanPresenceMiddleware(logger))
	}
	// Logger middleware runs after trace context is injected so logs include trace/span ids
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
	authMiddleware := customeMiddleware.InitAuthMiddleware(client, redisRepository, cfg.JwtSecretKey, cfg.Key, cfg.IV, *cfg, logger)

	cpsaction.Init(r, handlerLayer.CpsActionHandler, authMiddleware)
	budgetCategory.Init(r, handlerLayer.BudgetCategoryHandler, authMiddleware)
	avatar.Init(r, handlerLayer.AvatarHandler, authMiddleware)
	unlink.Init(r, handlerLayer.UnlinkHandler, authMiddleware)
	bpsUser.Init(r, handlerLayer.BpsHandler, authMiddleware)
	bank.Init(r, handlerLayer.BankHandler, authMiddleware)
	eventhandler.Init(r, handlerLayer.EventHandler, authMiddleware)
	wallet.Init(r, handlerLayer.WalletHandler, authMiddleware)
	topup.Init(r, handlerLayer.TopupHandler, authMiddleware)
	jobRole.Init(r, handlerLayer.jobRoleHandler, authMiddleware)
	customer.Init(r, handlerLayer.customerHandler, authMiddleware)
	bulk_service.Init(r, handlerLayer.bulkServiceHandler, authMiddleware)

	password.Init(r, handlerLayer.PasswordHandler, authMiddleware)

	feedback.Init(r, handlerLayer.FeedbackHandler, authMiddleware)
	productcode.Init(r, &handlerLayer.ProductCodeHandler, authMiddleware)
	advert.Init(r, handlerLayer.AdvertHandler, authMiddleware)
	portalcard.Init(r, handlerLayer.PortalCardHander, authMiddleware)
	accountvalidation.Init(r, handlerLayer.AccountValidation, authMiddleware)
	accountblock.Init(r, handlerLayer.AccountBlockHandler, authMiddleware)
	department.Init(r, &handlerLayer.DepartmentHandler, authMiddleware)
	hqRoute.Init(r, handlerLayer.HqHandler, authMiddleware)
	fayda.Init(r, handlerLayer.FaydaHandler, authMiddleware)
	service_details.Init(r, handlerLayer.ServiceDetailsHandler, authMiddleware)
	service.Init(r, handlerLayer.ServicesHandler, authMiddleware)
	permission_details.Init(r, handlerLayer.Permission, authMiddleware)
	cps_user_det.Init(r, handlerLayer.CPSUser, authMiddleware)
	device_version.Init(r, handlerLayer.DeviceVersionHandler, authMiddleware)
	bankvaultroutes.Init(r, handlerLayer.BankVaultHandler, authMiddleware)
	vaultgroupcategory.Init(r, handlerLayer.VaultGroupCategoryHandler, authMiddleware)

	donation.Init(r, handlerLayer.DonationHandler, authMiddleware)
	donation_category.Init(r, handlerLayer.DonationCategoryHandler, authMiddleware)
	donation_company.Init(r, handlerLayer.DonationCompanyHandler, authMiddleware)

	amountBasedAuth.Init(r, handlerLayer.AmountBasedAuthHandler, authMiddleware)
	notification.Init(r, handlerLayer.NotificationHandler, authMiddleware)
	kyc_routing.Init(r, handlerLayer.KYCVerifierHandler, authMiddleware)

	newscategory_routing.Init(r, handlerLayer.NewsCategoryHandler, authMiddleware)
	newstag_routing.Init(r, handlerLayer.NewsTagHandler, authMiddleware)
	bps_actionrole_routing.Init(r, handlerLayer.BPSActionRoleHandler, authMiddleware)
	cps_actionrole_routing.Init(r, handlerLayer.CPSActionRoleHandler, authMiddleware)
	sitota.Init(r, handlerLayer.SitotaHandler, authMiddleware)
	encryption.Init(r, handlerLayer.EncryptionHandler, authMiddleware)
	transaction.Init(r, handlerLayer.TransactionHandler, authMiddleware)
	vaultAmountTier.Init(r, handlerLayer.AmountTierHandler, authMiddleware)
	event_merchant_routing.Init(r, handlerLayer.EventMerchantHandler, authMiddleware)
	access_list_segmentation.Init(r, handlerLayer.AccessLostSegmentationHandler, authMiddleware)
	ecommerce_merchant.Init(r, handlerLayer.EcommerceMerchantHandler, authMiddleware)
	customer_seg.Init(r, handlerLayer.CustomerSegmentationHandler, authMiddleware)

	roles.Init(r, handlerLayer.RoleHandler, authMiddleware)
	// router.Mount("/api/v1/cbesuperapp/cps_action", r)

	// // Mini App Proxy Routes
	// miniAppProxyHandler := miniapp.CreateMiniAppProxyHandler(logger, cfg)
	// if miniAppProxyHandler == nil {
	// 	logger.Fatalf("Failed to create mini-app proxy handler")
	// }

	// r.Route("/mini-apps", func(r chi.Router) {
	// 	r.Use(authMiddleware.AuthenticateToken)
	// 	r.Handle("/*", miniAppProxyHandler)
	// })

	// Mini App Category Proxy Routes
	// miniAppCategoryProxyHandler := miniapp.CreateMiniAppCategoryProxyHandler(logger, cfg)
	// if miniAppProxyHandler == nil {
	// 	logger.Fatalf("Failed to create mini-app proxy handler")
	// }

	// r.Route("/mini-apps/categories", func(r chi.Router) {
	// 	r.Use(authMiddleware.AuthenticateToken)
	// 	r.Handle("/*", miniAppCategoryProxyHandler)
	// })

	// Public routes for password_rule

	// Wrap all CPS routes in a secured router that authenticates first, then applies the central guard
	// secured := chi.NewRouter()

	// secured.Use("/password_rule", r.HandleFunc("/",handlerLayer.PasswordHandler.GetPasswordRule))

	// secured.Use(authMiddleware.AuthenticateToken)
	// // Central CPS Action Guard (Option B): authorize by role_id + action_name with cache.
	// // Whitelist CPSAction endpoints under /actions (approve/reject/list...), allow them all.

	// secured.Use(customeMiddleware.CPSActionRouteGuard([]string{"/password_rule", "/actions", "/actions/{action_code}/approve", "/actions/{action_code}/reject"}))

	// // router.Mount("/api/v1/cbesuperapp/cps_action", public)
	// router.Mount("/api/v1/cbesuperapp/cps_action", secured)
	// // router.Use(customeMiddleware.ChiCORS())

	secured := chi.NewRouter()

	// --------------------
	// Public (no auth)
	// --------------------
	secured.Route("/password_rule", func(r chi.Router) {
		r.Get("/", handlerLayer.PasswordHandler.GetPasswordRule)
	})

	// --------------------
	// Protected routes
	// --------------------
	secured.Group(func(r chi.Router) {
		// Auth first
		r.Use(authMiddleware.AuthenticateToken)

		// CPS Action Guard
		r.Use(customeMiddleware.CPSActionRouteGuard([]string{
			"/actions",
			"/actions/{action_code}/approve",
			"/actions/{action_code}/reject",
		}))
		r.Mount("/", r)

		// CPS Actions
		// r.Route("/actions", func(r chi.Router) {
		// 	r.Get("/", handlerLayer.CpsActionHandler.GetCPSActionsByDepartment)
		// 	r.Post("/{action_code}/approve", handlerLayer.CpsActionHandler.ApproveCPSAction)
		// 	r.Post("/{action_code}/reject", handlerLayer.CpsActionHandler.RejectCPSAction)
		// })
	})
	// Mount
	router.Mount("/api/v1/cbesuperapp/cps_action", secured)

	// Swagger documentation routes
	router.Get("/api/v1/cbesuperapp/cps_action/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/cbesuperapp/cps_action/swagger/doc.json"),
	))
	router.Get("/api/v1/cbesuperapp/cps_action/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/cbesuperapp/cps_action/swagger/index.html", http.StatusMovedPermanently)
	})
}
