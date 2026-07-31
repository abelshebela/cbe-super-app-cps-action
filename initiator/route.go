package initiator

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	cps_auth "github.com/abelshebela/cbe-super-app-cps-action/grpc/auth/proto"
	access_list_segmentation "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/access_list_segmentaion"
	active_state "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/active_state"
	accountblock "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/account_block"
	ap_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/account_product"
	apc_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/account_product_category"
	account_sub_type_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/account_sub_type"
	accountvalidation "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/account_validation"
	advert "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/ad"
	amountBasedAuth "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/amount_based_auth"
	avatar "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/avatar"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/bank"
	bps_action "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/bps_action"
	bps_actionrole_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/bps_action_role"
	cps_actionrole_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/cps_action_role"
	customerkyc "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/customer_kyc"
	selfActivationKYC "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/customer_kyc_self"
	device_version "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/device_version"
	ecommerce_merchant "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/ecommerce-merchant"
	event_merchant_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/event_merchant"
	kyc_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/kyc_verifier"
	logistic_merchant_router "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/logistic_merchant"
	newscategory_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/news_category"
	newstag_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/news_tag"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/role_delegation"
	tac_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/term_and_condition"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/services"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/transaction"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/platform/telemetry"

	bankvaultroutes "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/bankvault"
	bpsUser "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/bps_user"
	budgetCategory "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/budget_category"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/bulk_service"
	cpsaction "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/cps_action"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/customer"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/department"
	eventhandler "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/event"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/notification"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/topup"
	vaultcategory "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/vault"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/wallet"

	cps_user_det "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/cps_user"
	customer_group_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/customer_group"
	donation "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/donation"
	donation_category "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/donation_category"
	donation_company "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/donation_company"
	fayda "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/fayda"
	feedback "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/feedback"
	hqRoute "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/hq"
	jobRole "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/job_roles"
	password "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/password_rule"
	superapp_role_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/superapp_role"

	// permission_details "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/permission"
	portalcard "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/portal_card"
	roles "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/roles"

	// servicesMiddleware "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/services"
	sitota "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/sitota"
	unlink "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/unlink"
	customeMiddleware "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"

	ussd_merchant_rout "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/ussd_merchant"
	utility_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/utility"
	survey_sampling_routing "github.com/abelshebela/cbe-super-app-cps-action/internal/glue/routing/survey_sampling"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	sharedMiddleware "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.uber.org/zap"
	// "github.com/go-chi/httprate"
)

func InitRoute(ctx context.Context, router *chi.Mux, encryptionMiddleware sharedMiddleware.EncMiddleware, handlerLayer Handler, client cps_auth.CpsAuthServiceClient, redisRepository storage.RedisRepository, logger utils.Logger, cfg *config.VaultConfig) {

	r := chi.NewRouter()

	router.Use(customeMiddleware.CORS(cfg))
	// middleware for encryption and decryption of request and response body
	// router.Use(encryptionMiddleware.SecureTunnelMiddleware())
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			header := sharedMiddleware.ExtractHeader(req)
			if header.EnableEncryption.IsValid() && header.EnableEncryption == sharedMiddleware.Enabled {
				encryptionMiddleware.SecureTunnelMiddleware()(next).ServeHTTP(w, req)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	// CORS must run early so preflight OPTIONS requests are handled before auth/logging
	// Security http rate limitter
	/*
		router.Use(httprate.LimitByIP(100, 1*time.Minute))
	*/
	// Security headers: HSTS, X-Content-Type-Options, X-Frame-Options, CSP, Cache-Control
	router.Use(customeMiddleware.SecurityHeaders)
	// Inject trace and span ids from OpenTelemetry span into context for logger extraction
	router.Use(telemetry.TraceContextMiddleware())
	// Optional debug middleware to detect missing spans. Enable by setting OTEL_DEBUG_TRACE_PRESENCE=true
	if os.Getenv("OTEL_DEBUG_TRACE_PRESENCE") == "true" {
		router.Use(telemetry.SpanPresenceMiddleware(logger))
	}
	// Logger middleware runs after trace context is injected so logs include trace/span ids
	router.Use(customeMiddleware.ChiLogger(logger))
	// Bind request context to response writer so Send*Response helpers include trace_id/request_id
	router.Use(customeMiddleware.BindRequestContext)

	router.Use(customeMiddleware.HandlePanic(logger))
	router.Use(chiMiddleware.Timeout(30 * time.Second))
	router.Use(chiMiddleware.Compress(5, "application/json"))
	// router.Use(sharedMiddleware.SecureTunnelMiddleware)

	authMiddleware := customeMiddleware.InitAuthMiddleware(client, redisRepository, nil, cfg.JwtSecretKey, cfg.Key, cfg.IV, *cfg, logger)

	// Apply CPS Action Route Guard for comprehensive path protection
	// Whitelist CPS Action endpoints that don't need action-based validation
	// whitelist := []string{"cps_action", "cps_actions"}
	// actionRouteGuard := customeMiddleware.CPSActionRouteGuard(whitelist)

	r.Get("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "CPS ACTION IS ACTIVE"}); err != nil {
			logger.Errorf("Failed to write health check response", zap.Error(err))
		}
	})

	cpsaction.Init(r, handlerLayer.CpsActionHandler, authMiddleware)
	bps_action.Init(r, handlerLayer.BpsActionHandler, authMiddleware)
	budgetCategory.Init(r, handlerLayer.BudgetCategoryHandler, authMiddleware)
	avatar.Init(r, handlerLayer.AvatarHandler, authMiddleware)
	unlink.Init(r, handlerLayer.UnlinkHandler, authMiddleware)
	bpsUser.Init(r, handlerLayer.BpsHandler, authMiddleware)
	bank.Init(r, handlerLayer.BankHandler, authMiddleware)
	account_sub_type_routing.Init(r, handlerLayer.AccountSubTypeHandler, authMiddleware)
	apc_routing.Init(r, handlerLayer.AccountProductCategoryHandler, authMiddleware)
	ap_routing.Init(r, handlerLayer.AccountProductHandler, authMiddleware)
	tac_routing.Init(r, handlerLayer.TermAndConditionHandler, authMiddleware)
	eventhandler.Init(r, handlerLayer.EventHandler, authMiddleware)
	active_state.Init(r, handlerLayer.ActiveStateHandler, authMiddleware)
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
	services.Init(r, handlerLayer.ServicesHandler, authMiddleware)
	// permission_details.Init(r, handlerLayer.Permission, authMiddleware)
	cps_user_det.Init(r, handlerLayer.CPSUser, authMiddleware)
	device_version.Init(r, handlerLayer.DeviceVersionHandler, authMiddleware)
	bankvaultroutes.Init(r, handlerLayer.BankVaultHandler, authMiddleware)
	vaultcategory.Init(r, handlerLayer.VaultCategoryHandler, authMiddleware)

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
	// encryption.Init(r, handlerLayer.EncryptionHandler, authMiddleware)
	transaction.Init(r, handlerLayer.TransactionHandler, authMiddleware)
	event_merchant_routing.Init(r, handlerLayer.EventMerchantHandler, authMiddleware)
	access_list_segmentation.Init(r, handlerLayer.AccessLostSegmentationHandler, authMiddleware)
	ecommerce_merchant.Init(r, handlerLayer.EcommerceMerchantHandler, authMiddleware)
	// customer_seg.Init(r, handlerLayer.CustomerSegmentationHandler, authMiddleware)
	customer_group_routing.Init(r, handlerLayer.CustomerGroupHandler, authMiddleware)
	superapp_role_routing.Init(r, handlerLayer.SuperAppRoleHandler, authMiddleware)
	// cps_roles.Init(r, handlerLayer.CPSRolesHandler, authMiddleware)
	logistic_merchant_router.Init(r, handlerLayer.LogisticsMerchantHandler, authMiddleware)
	customerkyc.Init(r, handlerLayer.CustomerKYCHandler, authMiddleware)
	selfActivationKYC.Init(r, handlerLayer.SelfActivationKYCHandler, authMiddleware)
	roles.Init(r, handlerLayer.RoleHandler, authMiddleware)
	role_delegation.Init(r, handlerLayer.RoleDelegationHandler, authMiddleware)
	utility_routing.Init(r, handlerLayer.UtilityHandler, authMiddleware)
	survey_sampling_routing.Init(r, handlerLayer.SurveySamplingHandler, authMiddleware)

	secured := chi.NewRouter()

	// secured.Route("/password_rule", func(r chi.Router) {
	// 	r.Get("/", handlerLayer.PasswordHandler.GetPasswordRule)
	// 	r.Group(func(r chi.Router) {
	// 		r.Use(authMiddleware.AuthenticateToken)
	// 		// r.Use(customeMiddleware.CPSActionRouteGuard([]string{}))
	// 		r.Patch("/{id}", handlerLayer.PasswordHandler.RequestPasswordRuleUpdate)
	// 	})
	// })

	secured.Group(func(r chi.Router) {
		// Auth first
		r.Use(authMiddleware.AuthenticateToken)
		// Global role validation - only allow viewer, maker, checker, auditor roles
		// r.Use(authMiddleware.ValidateRequiredRoles)
		// r.Use(actionRouteGuard)

		// CPS Action Guard
		// // r.Use(customeMiddleware.CPSActionRouteGuard([]string{
		// // 	"/actions",
		// // 	"/actions/{action_code}/approve",
		// // 	"/actions/{action_code}/reject",
		// // }))
		// Routes will be mounted below

	})

	// Swagger routes - only mount for dev/qa/uat environments and require authentication
	if isSwaggerEnabled(cfg.GoEnv) {
		secured.Group(func(r chi.Router) {
			// Require authentication for swagger routes
			// r.Use(authMiddleware.AuthenticateToken)

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
	// data := docs.SwaggerJSONBytes
	return func(w http.ResponseWriter, r *http.Request) {
		// if len(data) == 0 {
		http.Error(w, "Swagger spec not found", http.StatusNotFound)
		return
		// }
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		// _, _ = w.Write(data)
	}
}
