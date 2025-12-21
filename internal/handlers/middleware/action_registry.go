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
	"POST /notifications":               "Notification",
	"PATCH /notifications/{id}":         "Notification",
	"PATCH /notifications/enable/{id}":  "Notification",
	"PATCH /notifications/disable/{id}": "Notification",
	"DELETE /notifications/{id}":        "Notification",

	// MiniAppMerchant
	"POST /mini-app-merchants":               "MiniAppMerchant",
	"PATCH /mini-app-merchants/{id}":         "MiniAppMerchant",
	"PATCH /mini-app-merchants/enable/{id}":  "MiniAppMerchant",
	"PATCH /mini-app-merchants/disable/{id}": "MiniAppMerchant",
	"DELETE /mini-app-merchants/{id}":        "MiniAppMerchant",

	// Advert
	"POST /adverts":               "Advert",
	"PATCH /adverts/{id}":         "Advert",
	"DELETE /advert/{id}":         "Advert",
	"PATCH /adverts/{id}/enable":  "Advert",
	"PATCH /adverts/{id}/disable": "Advert",

	// AccountBlock
	"POST /account_block/branches/enable":   "AccountBlock",
	"POST /account_block/branches/disable":  "AccountBlock",
	"POST /account_block/regions/enable":    "AccountBlock",
	"POST /account_block/regions/disable":   "AccountBlock",
	"POST /account_block/districts/enable":  "AccountBlock",
	"POST /account_block/districts/disable": "AccountBlock",
	"POST /account_block/cities/enable":     "AccountBlock",
	"POST /account_block/cities/disable":    "AccountBlock",

	// AccountValidation
	"PATCH /account_validation/update/{id}": "AccountValidation",

	// AmountBasedAuth
	"PATCH /amount_based_auth/update/{method}/{id}": "AmountBasedAuth",
	"PATCH /amount_based_auth/reject/{id}":          "AmountBasedAuth",

	// Avatar
	"POST /avatar":            "Avatar",
	"DELETE /avatar/{id}":     "Avatar",
	"PATCH /avatar/disable/*": "Avatar",
	"PATCH /avatar/enable/*":  "Avatar",
	"PATCH /avatar/{id}":      "Avatar",

	// Bank
	"POST /banks":               "Bank",
	"PATCH /banks/{id}":         "Bank",
	"DELETE /banks/{id}":        "Bank",
	"PATCH /banks/{id}/enable":  "Bank",
	"PATCH /banks/{id}/disable": "Bank",
	"PATCH /banks/{id}/logo":    "Bank",

	// BankVault
	"POST /vault/products/create":        "BankVault",
	"PATCH /vault/products/update/{id}":  "BankVault",
	"DELETE /vault/products/delete/{id}": "BankVault",
	"PATCH /vault/products/disable/{id}": "BankVault",
	"PATCH /vault/products/enable/{id}":  "BankVault",

	// ActionRole (BPS)
	"POST /bps-action-roles":                 "ActionRole",
	"PATCH /bps-action-roles/{code}":         "ActionRole",
	"PATCH /bps-action-roles/{code}/enable":  "ActionRole",
	"PATCH /bps-action-roles/{code}/disable": "ActionRole",

	// BPSUser
	"POST /bps_users/disable/{user_code}": "BPSUser",
	"POST /bps_users/enable/{user_code}":  "BPSUser",

	// BudgetCategory
	"POST /budget-category":               "BudgetCategory",
	"PATCH /budget-category/{id}":         "BudgetCategory",
	"DELETE /budget-category/{id}":        "BudgetCategory",
	"PATCH /budget-category/enable/{id}":  "BudgetCategory",
	"PATCH /budget-category/disable/{id}": "BudgetCategory",

	// BulkService
	"POST /bulk_services/disable": "BulkService",
	"POST /bulk_services/enable":  "BulkService",

	// CpsActionRole
	"POST /cps-action-roles":                 "CpsActionRole",
	"PATCH /cps-action-roles/{code}":         "CpsActionRole",
	"PATCH /cps-action-roles/{code}/enable":  "CpsActionRole",
	"PATCH /cps-action-roles/{code}/disable": "CpsActionRole",

	// CpsUser
	"POST /cps_users/create":               "CpsUser",
	"PATCH /cps_users/update/{user_code}":  "CpsUser",
	"DELETE /cps_users/delete/{user_code}": "CpsUser",
	"POST /cps_users/disable/{user_code}":  "CpsUser",
	"POST /cps_users/enable/{user_code}":   "CpsUser",

	// Customer
	"PATCH /customers/enable/{id}":            "Customer",
	"PATCH /customers/enable_otp_verify/{id}": "Customer",
	"PATCH /customers/disable/{id}":           "Customer",
	"PATCH /customers/fayda/enable/{id}":      "Customer",

	// Department
	"POST /departments":               "Department",
	"PATCH /departments/{id}":         "Department",
	"PATCH /departments/enable/{id}":  "Department",
	"PATCH /departments/disable/{id}": "Department",

	// DeviceVersion
	"POST /device_versions":               "DeviceVersion",
	"PATCH /device_versions/{id}":         "DeviceVersion",
	"PATCH /device_versions/enable/{id}":  "DeviceVersion",
	"PATCH /device_versions/disable/{id}": "DeviceVersion",

	// Donation
	"POST /donation":               "Donation",
	"PATCH /donation/{id}":         "Donation",
	"PATCH /donation/image/{id}":   "Donation",
	"DELETE /donation/image/{id}":  "Donation",
	"POST /donation/image/{id}":    "Donation",
	"PATCH /donation/enable/{id}":  "Donation",
	"PATCH /donation/disable/{id}": "Donation",

	// DonationCategory
	"POST /donation_category":               "DonationCategory",
	"PATCH /donation_category/{id}":         "DonationCategory",
	"PATCH /donation_category/enable/{id}":  "DonationCategory",
	"PATCH /donation_category/disable/{id}": "DonationCategory",

	// DonationCompany
	"POST /donation_company":               "DonationCompany",
	"PATCH /donation_company/{id}":         "DonationCompany",
	"PATCH /donation_company/enable/{id}":  "DonationCompany",
	"PATCH /donation_company/disable/{id}": "DonationCompany",

	// Event
	"POST /events":               "Event",
	"PATCH /events/{id}":         "Event",
	"PATCH /events/enable/{id}":  "Event",
	"PATCH /events/disable/{id}": "Event",
	"DELETE /events/{id}":        "Event",

	// HQ
	"POST /hq/block_time":      "BlockTime",
	"POST /hq/archive_time":    "Archive",
	"POST /hq/password_expiry": "PasswordExpiry",

	// Fayda
	"POST /fayda_account/disable/{user_code}": "Fayda",
	"POST /fayda_account/enable/{user_code}":  "Fayda",

	// KYCVerifier
	"PATCH /kyc_verifier/update/{id}":  "KYCVerifier",
	"PATCH /kyc_verifier/approve/{id}": "KYCVerifier",

	// NewsCategory
	"POST /news/category/create": "NewsCategory",
	"DELETE /news/category/{id}": "NewsCategory",
	"PATCH /news/category/{id}":  "NewsCategory",

	// NewsTag
	"POST /news/tags/create": "NewsTag",
	"DELETE /news/tags/{id}": "NewsTag",
	"PATCH /news/tags/{id}":  "NewsTag",

	// PasswordRule
	"PATCH /password_rule/{id}": "PasswordRule",
	"POST /password_rule/check": "PasswordRule",

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
	"POST /services":               "Service",
	"PATCH /services/{id}":         "Service",
	"PATCH /services/{id}/enable":  "Service",
	"PATCH /services/{id}/disable": "Service",

	// Topup
	"POST /topups":               "Topup",
	"PATCH /topups/{id}":         "Topup",
	"DELETE /topups/{id}":        "Topup",
	"PATCH /topups/{id}/enable":  "Topup",
	"PATCH /topups/{id}/disable": "Topup",

	// UnlinkDevice
	"PATCH /unlink/user_cif/{user_code}": "UnlinkDevice",

	// VaultGroupCategory
	"POST /vaultgroupcategory/create":        "VaultGroupCategory",
	"PATCH /vaultgroupcategory/update/{id}":  "VaultGroupCategory",
	"DELETE /vaultgroupcategory/delete/{id}": "VaultGroupCategory",
	"PATCH /vaultgroupcategory/enable/{id}":  "VaultGroupCategory",
	"PATCH /vaultgroupcategory/disable/{id}": "VaultGroupCategory",

	// Wallet
	"POST /wallets":               "Wallet",
	"PATCH /wallets/{id}":         "Wallet",
	"DELETE /wallets/{id}":        "Wallet",
	"PATCH /wallets/{id}/enable":  "Wallet",
	"PATCH /wallets/{id}/disable": "Wallet",

	// 	ROLE
	"POST /job_role":               "JOBROLE",
	"PATCH /job_role/{id}":         "JOBROLE",
	"DELETE /job_role/{id}":        "JOBROLE",
	"PATCH /job_role/{id}/enable":  "JOBROLE",
	"PATCH /job_role/{id}/disable": "JOBROLE",
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
