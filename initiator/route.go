package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	cps_auth "cbe-super-app-cps-action/grpc/auth/proto"
	access_list_segmentation "cbe-super-app-cps-action/internal/glue/routing/access_list_segmentaion"
	accountblock "cbe-super-app-cps-action/internal/glue/routing/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/glue/routing/account_validation"
	advert "cbe-super-app-cps-action/internal/glue/routing/ad"
	amountBasedAuth "cbe-super-app-cps-action/internal/glue/routing/amount_based_auth"
	avatar "cbe-super-app-cps-action/internal/glue/routing/avatar"
	"cbe-super-app-cps-action/internal/glue/routing/bank"
	bps_action "cbe-super-app-cps-action/internal/glue/routing/bps_action"
	bps_actionrole_routing "cbe-super-app-cps-action/internal/glue/routing/bps_action_role"
	cps_actionrole_routing "cbe-super-app-cps-action/internal/glue/routing/cps_action_role"
	customerkyc "cbe-super-app-cps-action/internal/glue/routing/customer_kyc"
	device_version "cbe-super-app-cps-action/internal/glue/routing/device_version"
	ecommerce_merchant "cbe-super-app-cps-action/internal/glue/routing/ecommerce-merchant"
	event_merchant_routing "cbe-super-app-cps-action/internal/glue/routing/event_merchant"
	kyc_routing "cbe-super-app-cps-action/internal/glue/routing/kyc_verifier"
	logistic_merchant_router "cbe-super-app-cps-action/internal/glue/routing/logistic_merchant"
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

	cps_roles "cbe-super-app-cps-action/internal/glue/routing/cps_roles"
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
	roles "cbe-super-app-cps-action/internal/glue/routing/roles"
	service "cbe-super-app-cps-action/internal/glue/routing/services"
	sitota "cbe-super-app-cps-action/internal/glue/routing/sitota"
	unlink "cbe-super-app-cps-action/internal/glue/routing/unlink"
	vaultAmountTier "cbe-super-app-cps-action/internal/glue/routing/vault_amount_tier"
	customeMiddleware "cbe-super-app-cps-action/internal/handlers/middleware"

	ussd_merchant_rout "cbe-super-app-cps-action/internal/glue/routing/ussd_merchant"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"cbe-super-app-cps-action/docs" 
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
	bps_action.Init(r, handlerLayer.BpsActionHandler, authMiddleware)
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
	ussd_merchant_rout.Init(r, handlerLayer.UssdMerchantHandler, authMiddleware)
	password.Init(r, handlerLayer.PasswordHandler, authMiddleware)

	feedback.Init(r, handlerLayer.FeedbackHandler, authMiddleware)
	advert.Init(r, handlerLayer.AdvertHandler, authMiddleware)
	portalcard.Init(r, handlerLayer.PortalCardHander, authMiddleware)
	accountvalidation.Init(r, handlerLayer.AccountValidation, authMiddleware)
	accountblock.Init(r, handlerLayer.AccountBlockHandler, authMiddleware)
	department.Init(r, &handlerLayer.DepartmentHandler, authMiddleware)
	hqRoute.Init(r, handlerLayer.HqHandler, authMiddleware)
	fayda.Init(r, handlerLayer.FaydaHandler, authMiddleware)
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
	cps_roles.Init(r, handlerLayer.CPSRolesHandler, authMiddleware)
	logistic_merchant_router.Init(r, handlerLayer.LogisticsMerchantHandler, authMiddleware)
	customerkyc.Init(r, handlerLayer.CustomerKYCHandler, authMiddleware)

	roles.Init(r, handlerLayer.RoleHandler, authMiddleware)

	secured := chi.NewRouter()

	secured.Route("/password_rule", func(r chi.Router) {
		r.Get("/", handlerLayer.PasswordHandler.GetPasswordRule)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.AuthenticateToken)
			// r.Use(customeMiddleware.CPSActionRouteGuard([]string{}))
			r.Patch("/{id}", handlerLayer.PasswordHandler.RequestPasswordRuleUpdate)
		})
	})

	secured.Group(func(r chi.Router) {
		// Auth first
		r.Use(authMiddleware.AuthenticateToken)

		// CPS Action Guard
		// // r.Use(customeMiddleware.CPSActionRouteGuard([]string{
		// // 	"/actions",
		// // 	"/actions/{action_code}/approve",
		// // 	"/actions/{action_code}/reject",
		// // }))
		// r.Mount("/", r)

	})

	// Swagger routes - only mount for dev/qa/uat environments and require authentication
	if isSwaggerEnabled(cfg.GoEnv) {
		secured.Group(func(r chi.Router) {
			// Require authentication for swagger routes
			r.Use(authMiddleware.AuthenticateToken)

			// Build after: swag init -g cmd/main.go -o docs && go run scripts/merge_swagger_examples.go
			r.Get("/swagger/doc.json", serveSwaggerDocEmbedded())
			r.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL("/api/v1/cbesuperapp/cps_action/swagger/doc.json"),
			))
			r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/api/v1/cbesuperapp/cps_action/swagger/index.html", http.StatusMovedPermanently)
			})
		})
	}

	secured.Mount("/", r)
	// Mount
	router.Mount("/api/v1/cbesuperapp/cps_action", secured)
}

// Swagger is disabled for: "staging", "production"
func isSwaggerEnabled(goEnv string) bool {
	env := strings.ToLower(strings.TrimSpace(goEnv))
	return env == "dev" || env == "qa" || env == "uat"
}

// serveSwaggerDocEmbedded serves the embedded docs.SwaggerJSONBytes (run merge script before build to include examples).
func serveSwaggerDocEmbedded() http.HandlerFunc {
	data := docs.SwaggerJSONBytes
	return func(w http.ResponseWriter, r *http.Request) {
		if len(data) == 0 {
			http.Error(w, "Swagger spec not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		_, _ = w.Write(data)
	}
}
