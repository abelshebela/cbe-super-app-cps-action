package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
)

// cpsActionRegistry maps METHOD + " " + RoutePattern to CPS action name
// Only write/approval routes should be included here
var cpsActionRegistry = map[string]string{
	// Notification
	"POST /notifications":               "NOTIFICATIONS",
	"PATCH /notifications/{id}":         "NOTIFICATIONS",
	"PATCH /notifications/enable/{id}":  "NOTIFICATIONS",
	"PATCH /notifications/disable/{id}": "NOTIFICATIONS",
	"DELETE /notifications/{id}":        "NOTIFICATIONS",

	// MiniAppMerchant
	"POST /mini-app-merchants":               "MINIAPPMERCHANT",
	"PATCH /mini-app-merchants/{id}":         "MINIAPPMERCHANT",
	"PATCH /mini-app-merchants/enable/{id}":  "MINIAPPMERCHANT",
	"PATCH /mini-app-merchants/disable/{id}": "MINIAPPMERCHANT",
	"DELETE /mini-app-merchants/{id}":        "MINIAPPMERCHANT",

	// Advert
	"POST /adverts":               "ADVERT",
	"PATCH /adverts/{id}":         "ADVERT",
	"DELETE /advert/{id}":         "ADVERT",
	"PATCH /adverts/{id}/enable":  "ADVERT",
	"PATCH /adverts/{id}/disable": "ADVERT",

	// AccountBlock
	"POST /account_block/branches/enable":   "ACCOUNTBLOCK",
	"POST /account_block/branches/disable":  "ACCOUNTBLOCK",
	"POST /account_block/regions/enable":    "ACCOUNTBLOCK",
	"POST /account_block/regions/disable":   "ACCOUNTBLOCK",
	"POST /account_block/districts/enable":  "ACCOUNTBLOCK",
	"POST /account_block/districts/disable": "ACCOUNTBLOCK",
	"POST /account_block/cities/enable":     "ACCOUNTBLOCK",
	"POST /account_block/cities/disable":    "ACCOUNTBLOCK",

	// AccountValidation
	"GET /account_validation":               "ACCOUNTVALIDATION",
	"GET /account_validation/{id}":          "ACCOUNTVALIDATION",
	"PATCH /account_validation/update/{id}": "AccountValidation",

	// AmountBasedAuth
	"PATCH /amount_based_auth/update/{method}/{id}": "AMOUNTBASEDAUTH",
	"GET /amount_based_auth":                        "AMOUNTBASEDAUTH",
	"PATCH /amount_based_auth/reject/{id}":          "AMOUNTBASEDAUTH",

	// Avatar
	"POST /avatar":            "AVATAR",
	"DELETE /avatar/{id}":     "AVATAR",
	"PATCH /avatar/disable/*": "AVATAR",
	"PATCH /avatar/enable/*":  "AVATAR",
	"PATCH /avatar/{id}":      "AVATAR",

	// Bank
	"POST /banks":               "BANK",
	"PATCH /banks/{id}":         "BANK",
	"DELETE /banks/{id}":        "BANK",
	"PATCH /banks/{id}/enable":  "BANK",
	"PATCH /banks/{id}/disable": "BANK",
	"PATCH /banks/{id}/logo":    "BANK",

	// BankVault
	"POST /vault/products/create":        "BANKVAULT",
	"PATCH /vault/products/update/{id}":  "BANKVAULT",
	"DELETE /vault/products/delete/{id}": "BANKVAULT",
	"PATCH /vault/products/disable/{id}": "BANKVAULT",
	"PATCH /vault/products/enable/{id}":  "BANKVAULT",

	// ActionRole (BPS)
	"POST /bps-action-roles":                 "BPSACTIONEROLE",
	"PATCH /bps-action-roles/{code}":         "BPSACTIONEROLE",
	"PATCH /bps-action-roles/{code}/enable":  "BPSACTIONEROLE",
	"PATCH /bps-action-roles/{code}/disable": "BPSACTIONEROLE",

	// BPSUser
	"POST /bps_users/disable/{user_code}": "BPSUSER",
	"POST /bps_users/enable/{user_code}":  "BPSUSER",

	// BudgetCategory
	"POST /budget-category":               "BUDGETCATEGORY",
	"PATCH /budget-category/{id}":         "BUDGETCATEGORY",
	"DELETE /budget-category/{id}":        "BUDGETCATEGORY",
	"PATCH /budget-category/enable/{id}":  "BUDGETCATEGORY",
	"PATCH /budget-category/disable/{id}": "BUDGETCATEGORY",

	// BulkService
	"POST /bulk_services/disable": "BULKSERVICE",
	"POST /bulk_services/enable":  "BULKSERVICE",

	// CpsActionRole
	"POST /cps-action-roles":                 "CPSACTIONROLE",
	"PATCH /cps-action-roles/{code}":         "CPSACTIONROLE",
	"PATCH /cps-action-roles/{code}/enable":  "CPSACTIONROLE",
	"PATCH /cps-action-roles/{code}/disable": "CPSACTIONROLE",

	// CpsUser
	"POST /cps_users/create":               "CPSUSER",
	"PATCH /cps_users/update/{user_code}":  "CPSUSER",
	"DELETE /cps_users/delete/{user_code}": "CPSUSER",
	"POST /cps_users/disable/{user_code}":  "CPSUSER",
	"POST /cps_users/enable/{user_code}":   "CPSUSER",

	// Customer
	"PATCH /customers/enable/{id}":            "CUSTOMER",
	"PATCH /customers/enable_otp_verify/{id}": "CUSTOMER",
	"PATCH /customers/disable/{id}":           "CUSTOMER",
	"PATCH /customers/fayda/enable/{id}":      "CUSTOMER",

	// Department
	"POST /departments":               "DEPARTMENT",
	"PATCH /departments/{id}":         "DEPARTMENT",
	"PATCH /departments/enable/{id}":  "DEPARTMENT",
	"PATCH /departments/disable/{id}": "DEPARTMENT",

	// DeviceVersion
	"POST /device_versions":               "DEVICEVERSION",
	"PATCH /device_versions/{id}":         "DEVICEVERSION",
	"PATCH /device_versions/enable/{id}":  "DEVICEVERSION",
	"PATCH /device_versions/disable/{id}": "DEVICEVERSION",

	// Donation
	"POST /donation":               "DONATION",
	"PATCH /donation/{id}":         "DONATION",
	"PATCH /donation/image/{id}":   "DONATION",
	"DELETE /donation/image/{id}":  "DONATION",
	"POST /donation/image/{id}":    "DONATION",
	"PATCH /donation/enable/{id}":  "DONATION",
	"PATCH /donation/disable/{id}": "DONATION",

	// DonationCategory
	"POST /donation_category":               "DONATIONCATEGORY",
	"PATCH /donation_category/{id}":         "DONATIONCATEGORY",
	"PATCH /donation_category/enable/{id}":  "DONATIONCATEGORY",
	"PATCH /donation_category/disable/{id}": "DONATIONCATEGORY",

	// DonationCompany
	"POST /donation_company":               "DONATIONCOMPANY",
	"PATCH /donation_company/{id}":         "DONATIONCOMPANY",
	"PATCH /donation_company/enable/{id}":  "DONATIONCOMPANY",
	"PATCH /donation_company/disable/{id}": "DONATIONCOMPANY",

	// Ecommerce merrchant
	"POST /ecommerce-merchant":               "ECOMMERCEMERCHANT",
	"PATCH /ecommerce-merchant/enable/{id}":  "ECOMMERCEMERCHANT",
	"PATCH /ecommerce-merchant/disable/{id}": "ECOMMERCEMERCHANT",
	"PATCH /ecommerce-merchant/{id}":         "ECOMMERCEMERCHANT",
	"DELETE /ecommerce-merchant/{id}":        "ECOMMERCEMERCHANT",

	// Encryption
	"POST /encryption/encrypt": "ENCRYPTION",

	// Event merchant
	"POST /event_merchants":               "EVENTMERCHANT",
	"PATCH /event_merchants/{id}":         "EVENTMERCHANT",
	"PATCH /event_merchants/enable/{id}":  "EVENTMERCHANT",
	"PATCH /event_merchants/disable/{id}": "EVENTMERCHANT",
	"DELETE /event_merchant/{id}":         "EVENTMERCHANT",

	// Event
	"POST /events":               "Event",
	"PATCH /events/{id}":         "Event",
	"PATCH /events/enable/{id}":  "Event",
	"PATCH /events/disable/{id}": "Event",
	"DELETE /events/{id}":        "Event",

	// HQ
	"POST /hq/block_time":      "HQ",
	"POST /hq/archive_time":    "HQ",
	"POST /hq/password_expiry": "HQ",

	// Fayda
	"POST /fayda_account/disable/{user_code}": "FAYDA",
	"POST /fayda_account/enable/{user_code}":  "FAYDA",

	// KYCVerifier
	"PATCH /kyc_verifier/update/{id}":  "KYCVERIFIER",
	"PATCH /kyc_verifier/approve/{id}": "KYCVERIFIER",

	// NewsCategory
	"POST /news/category/create": "NEWSCATEGORY",
	"DELETE /news/category/{id}": "NEWSCATEGORY",
	"PATCH /news/category/{id}":  "NEWSCATEGORY",

	// NewsTag
	"POST /news/tags/create": "NEWSTAG",
	"DELETE /news/tags/{id}": "NEWSTAG",
	"PATCH /news/tags/{id}":  "NEWSTAG",

	// Job roles
	"POST /job_roles":       "JOBROLES",
	"PATCH /job_roles/{id}": "JOBROLES",

	// Access list segmentation
	"POST /access_list_segmentation":              "ACCESSLISTEGMENTATION",
	"POST /access_list_segmentation/{id}":         "ACCESSLISTEGMENTATION",
	"POST /access_list_segmentation/enable/{id}":  "ACCESSLISTEGMENTATION",
	"POST /access_list_segmentation/disable/{id}": "ACCESSLISTEGMENTATION",

	// PasswordRule
	"PATCH /password_rule/{id}": "PASSWORDRULES",
	"POST /password_rule/check": "PASSWORDRULES",

	// PermissionGroup
	"POST /permissions":       "PermissionGroup",
	"PATCH /permissions/{id}": "PermissionGroup",

	// ProductCode (nested router)
	"PATCH /productcodes/{id}": "ProductCode",

	// ServiceDetails -> Service
	"PATCH /service/service_fee/update/{id}":         "Service",
	"PATCH /service/single_transfer_max/update/{id}": "Service",
	"PATCH /service/total_transfer_max/update":       "Service",
	"PATCH /service/minimum_transfer/update/{id}":    "Service",
	"DELETE /service/service_fee/delete/{id}":        "Service",

	// Services module
	"POST /services":               "SERVICE",
	"PATCH /services/{id}":         "SERVICE",
	"PATCH /services/{id}/enable":  "SERVICE",
	"PATCH /services/{id}/disable": "SERVICE",

	// Topup
	"POST /topups":               "TOPUP",
	"PATCH /topups/{id}":         "TOPUP",
	"DELETE /topups/{id}":        "TOPUP",
	"PATCH /topups/{id}/enable":  "TOPUP",
	"PATCH /topups/{id}/disable": "TOPUP",

	// UnlinkDevice
	"PATCH /unlink/user_cif/{user_code}": "UNLINKDEVICE",

	// VaultGroupCategory
	"POST /vaultgroupcategory/create":        "VAULTCATEGORY",
	"PATCH /vaultgroupcategory/update/{id}":  "VAULTCATEGORY",
	"DELETE /vaultgroupcategory/delete/{id}": "VAULTCATEGORY",
	"PATCH /vaultgroupcategory/enable/{id}":  "VAULTCATEGORY",
	"PATCH /vaultgroupcategory/disable/{id}": "VAULTCATEGORY",

	// Vault amount tier
	"POST /vault-amount-tier/create":        "VAULTAMOUNTTIER",
	"PATCH /vault-amount-tier/{id}/update":  "VAULTAMOUNTTIER",
	"DELETE /vault-amount-tier/{id}/delete": "VAULTAMOUNTTIER",
	"PATCH /vault-amount-tier/{id}/disable": "VAULTAMOUNTTIER",
	"PATCH /vault-amount-tier/{id}/enable":  "VAULTAMOUNTTIER",

	// Wallet
	"POST /wallets":               "WALLET",
	"PATCH /wallets/{id}":         "WALLET",
	"DELETE /wallets/{id}":        "WALLET",
	"PATCH /wallets/{id}/enable":  "WALLET",
	"PATCH /wallets/{id}/disable": "WALLET",

	// 	ROLE
	"POST /job_role":               "JOBROLE",
	"PATCH /job_role/{id}":         "JOBROLE",
	"DELETE /job_role/{id}":        "JOBROLE",
	"PATCH /job_role/{id}/enable":  "JOBROLE",
	"PATCH /job_role/{id}/disable": "JOBROLE",

	// Customer Segmentation
	"POST /customer-segmentations":              "CUSTOMERSEGMENTATIONS",
	"PATCH /customer-segmentations/{id}/update": "CUSTOMERSEGMENTATIONS",
	"GET /customer-segmentations":               "CUSTOMERSEGMENTATIONS",
	"GET /customer-segmentations/{id}":          "CUSTOMERSEGMENTATIONS",
	"DELETE /customer-segmentations/{id}":       "CUSTOMERSEGMENTATIONS",
}

func extractResource(path string) string {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func normalize(s string) string {
	s = strings.ToUpper(s)
	replacer := strings.NewReplacer(
		"-", "",
		"_", "",
	)
	return replacer.Replace(s)
}

func resolveActionName(relPath string, registry map[string]string) string {
	resource := extractResource(relPath)
	normalizedResource := normalize(resource)

	for _, action := range registry {
		if normalize(action) == normalizedResource {
			return action
		}
	}
	return ""
}

func CPSActionRouteGuard(whitelist []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ensure repo is initialized
			if cpsApproveRepo == nil {
				next.ServeHTTP(w, r) // fail-open if not configured
				return
			}

			rc := chi.RouteContext(r.Context())
			if rc == nil {
				next.ServeHTTP(w, r)
				return
			}

			pattern := routeFullPattern(rc)
			if pattern == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Normalize pattern-relative route (keeps placeholders like {id})
			relPattern := pattern
			if strings.HasPrefix(relPattern, "/api/v1/cbesuperapp/cps_action") {
				relPattern = strings.TrimPrefix(relPattern, "/api/v1/cbesuperapp/cps_action")
				if relPattern == "" {
					relPattern = "/"
				}
			}

			// Normalize actual path route (concrete values like /banks/567...)
			relPath := r.URL.Path
			if strings.HasPrefix(relPath, "/api/v1/cbesuperapp/cps_action") {
				relPath = strings.TrimPrefix(relPath, "/api/v1/cbesuperapp/cps_action")
				if relPath == "" {
					relPath = "/"
				}
			}

			// Allowlist (e.g., CPSAction endpoints)
			for _, p := range whitelist {
				if strings.HasPrefix(relPath, p) {
					next.ServeHTTP(w, r)
					return
				}
			}

			method := strings.ToUpper(r.Method)

			if method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			//actionName := resolveActionName(relPath, cpsActionRegistry)
			// keyPattern := method + " " + relPattern
			// actionName, ok := cpsActionRegistry[keyPattern]
			actionName := ""

			for _, v := range cpsActionRegistry {
				path := strings.ReplaceAll(relPath, "_", "")

				if strings.Contains(path, strings.ToLower(v)) {
					actionName = v
					break
				}
			}

			rawRoleID, _ := r.Context().Value(constants.ContextKey("role_id")).(string)
			roleID := utils.FirstHex24(rawRoleID)
			if roleID == "" {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			action := strings.ToUpper(strings.TrimSpace(actionName))
			cacheKey := roleID + ":" + action
			if ent, ok := cpsGuardCache.get(cacheKey); ok && ent.allow {
				next.ServeHTTP(w, r)
				return
			}

			allowed, err := cpsApproveRepo.ExistsByRoleAndAction(r.Context(), roleID, action)
			if err != nil {
				if guardLogger != nil {
					guardLogger.Errorf("central guard lookup failed: %v", err)
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			if allowed {
				cpsGuardCache.set(cacheKey, allowEntry{allow: true, exp: nowPlus(cpsGuardCache.ttl)})
			}
			if !allowed {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ResolveActionKey(request string) string {
	request = strings.TrimSpace(request)

	for key := range cpsActionRegistry {
		if strictAvatarMatch(key, request) {
			return key
		}
	}
	return ""
}

func strictAvatarMatch(template, request string) bool {
	tpl := strings.SplitN(template, " ", 2)
	req := strings.SplitN(request, " ", 2)

	if len(tpl) != 2 || len(req) != 2 {
		return false
	}

	// 1. METHOD must match
	if tpl[0] != req[0] {
		return false
	}

	tplParts := strings.Split(strings.Trim(tpl[1], "/"), "/")
	reqParts := strings.Split(strings.Trim(req[1], "/"), "/")

	// 2. root resource must match (avatar)
	if tplParts[0] != reqParts[0] {
		return false
	}

	// 3. template length must match request length
	if len(tplParts) != len(reqParts) {
		return false
	}

	// 4. strict ordered matching
	for i := range tplParts {
		tplSeg := tplParts[i]
		reqSeg := reqParts[i]

		// dynamic segment
		if strings.HasPrefix(tplSeg, "{") &&
			strings.HasSuffix(tplSeg, "}") {
			continue
		}

		if tplSeg != reqSeg {
			return false
		}
	}

	return true
}

func nowPlus(dur time.Duration) time.Time {
	return time.Now().Add(dur)
}

// routeFullPattern builds the full route pattern from chi's context,
// joining nested Route() segments (e.g., "/productcodes" + "/{id}")
func routeFullPattern(rc *chi.Context) string {
	if rc == nil {
		return ""
	}
	parts := rc.RoutePatterns
	if len(parts) == 0 {
		return ""
	}
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, "/") && b.Len() > 0 {
			b.WriteString("/")
		}
		b.WriteString(p)
	}
	s := strings.ReplaceAll(b.String(), "//", "/")
	return s
}

func deriveModuleFromPattern(pattern string) string {
	p := strings.TrimPrefix(pattern, "/")
	if p == "" {
		return ""
	}
	if i := strings.IndexByte(p, '/'); i >= 0 {
		p = p[:i]
	}
	// normalize underscores to hyphens then TitleCase hyphenated words
	p = strings.ReplaceAll(p, "_", "-")
	parts := strings.Split(p, "-")
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		r := []rune(part)
		// Uppercase first rune, lowercase rest
		b.WriteString(strings.ToUpper(string(r[0:1])))
		if len(r) > 1 {
			b.WriteString(strings.ToLower(string(r[1:])))
		}
	}
	res := b.String()
	return res
}
