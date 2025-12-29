package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
)

// cpsActionRegistry maps METHOD + " " + RoutePattern to CPS action name
// Only write/approval routes should be included here
var cpsActionRegistry = map[string]string{
	// Notification
	"POST notifications":   "NOTIFICATIONS",
	"PATCH notifications":  "NOTIFICATIONS",
	"DELETE notifications": "NOTIFICATIONS",

	// MiniAppMerchant
	"POST mini-app-merchants":   "MINIAPPMERCHANT",
	"PATCH mini-app-merchants":  "MINIAPPMERCHANT",
	"DELETE mini-app-merchants": "MINIAPPMERCHANT",

	// Advert
	"POST adverts":  "ADVERT",
	"PATCH adverts": "ADVERT",
	"DELETE advert": "ADVERT",

	// AccountBlock
	"POST account_block": "ACCOUNTBLOCK",

	// AccountValidation
	"GET account_validation":   "ACCOUNTVALIDATION",
	"PATCH account_validation": "ACCOUNTVALIDATION",

	// AmountBasedAuth
	"PATCH amount_based_auth": "AMOUNTBASEDAUTH",
	"GET amount_based_auth":   "AMOUNTBASEDAUTH",

	// Avatar
	"POST avatar":   "AVATAR",
	"DELETE avatar": "AVATAR",
	"PATCH avatar":  "AVATAR",
	// Bank
	"POST banks":   "BANK",
	"PATCH banks":  "BANK",
	"DELETE banks": "BANK",

	// BankVault
	"POST vault":   "BANKVAULT",
	"PATCH vault":  "BANKVAULT",
	"DELETE vault": "BANKVAULT",

	// ActionRole (BPS)
	"POST bps-action-roles":  "BPSACTIONEROLE",
	"PATCH bps-action-roles": "BPSACTIONEROLE",

	// BPSUser
	"POST bps_users": "BPSUSER",

	// BudgetCategory
	"POST budget-category":   "BUDGETCATEGORY",
	"PATCH budget-category":  "BUDGETCATEGORY",
	"DELETE budget-category": "BUDGETCATEGORY",

	// BulkService
	"POST bulk_services": "BULKSERVICE",

	// CpsActionRole
	"POST cps-action-roles":  "CPSACTIONROLE",
	"PATCH cps-action-roles": "CPSACTIONROLE",

	// CpsUser
	"POST cps_users":   "CPSUSER",
	"PATCH cps_users":  "CPSUSER",
	"DELETE cps_users": "CPSUSER",

	// Customer
	"PATCH customers": "CUSTOMER",

	// Department
	"POST departments":  "DEPARTMENT",
	"PATCH departments": "DEPARTMENT",

	// DeviceVersion
	"POST device_versions":  "DEVICEVERSION",
	"PATCH device_versions": "DEVICEVERSION",

	// Donation
	"POST donation":   "DONATION",
	"PATCH donation":  "DONATION",
	"DELETE donation": "DONATION",

	// DonationCategory
	"POST donation_category":  "DONATIONCATEGORY",
	"PATCH donation_category": "DONATIONCATEGORY",

	// DonationCompany
	"POST donation_company":  "DONATIONCOMPANY",
	"PATCH donation_company": "DONATIONCOMPANY",

	// Ecommerce merrchant
	"POST ecommerce-merchant":   "ECOMMERCEMERCHANT",
	"PATCH ecommerce-merchant":  "ECOMMERCEMERCHANT",
	"DELETE ecommerce-merchant": "ECOMMERCEMERCHANT",

	// Encryption
	"POST encryption": "ENCRYPTION",

	// Event merchant
	"POST event_merchants":  "EVENTMERCHANT",
	"PATCH event_merchants": "EVENTMERCHANT",
	"DELETE event_merchant": "EVENTMERCHANT",

	// Event
	"POST events":   "Event",
	"PATCH events":  "Event",
	"DELETE events": "Event",

	// HQ
	"POST hq": "HQ",

	// Fayda
	"POST /fayda_account/disable/{user_code}": "FAYDA",
	"POST /fayda_account/enable/{user_code}":  "FAYDA",

	// KYCVerifier
	"PATCH kyc_verifier": "KYCVERIFIER",

	// NewsCategory
	"POST /news/category/create": "NEWSCATEGORY",
	"DELETE /news/category/{id}": "NEWSCATEGORY",
	"PATCH /news/category/{id}":  "NEWSCATEGORY",

	// NewsTag
	"POST /news/tags/create": "NEWSTAG",
	"DELETE /news/tags/{id}": "NEWSTAG",
	"PATCH /news/tags/{id}":  "NEWSTAG",

	// Job roles
	"POST job_roles":  "JOBROLE",
	"PATCH job_roles": "JOBROLE",

	// Access list segmentation
	"POST access_list_segmentation": "ACCESSLISTSEGMENTATION",

	// PasswordRule
	"PATCH password_rule": "PASSWORDRULE",
	"POST password_rule":  "PASSWORDRULE",

	// PermissionGroup
	"POST permissions":  "PermissionGroup",
	"PATCH permissions": "PermissionGroup",

	// ProductCode (nested router)
	"PATCH productcodes": "ProductCode",

	// ServiceDetails -> Service
	"PATCH service":  "Service",
	"DELETE service": "Service",

	// Services module
	"POST services":  "SERVICE",
	"PATCH services": "SERVICE",

	// Topup
	"POST topups":   "TOPUP",
	"PATCH topups":  "TOPUP",
	"DELETE topups": "TOPUP",

	// UnlinkDevice
	"PATCH unlink": "UNLINKDEVICE",

	// VaultGroupCategory
	"POST vaultgroupcategory":   "VAULTCATEGORY",
	"PATCH vaultgroupcategory":  "VAULTCATEGORY",
	"DELETE vaultgroupcategory": "VAULTCATEGORY",

	// Vault amount tier
	"POST vault-amount-tier":   "VAULTAMOUNTTIER",
	"PATCH vault-amount-tier":  "VAULTAMOUNTTIER",
	"DELETE vault-amount-tier": "VAULTAMOUNTTIER",

	// Wallet
	"POST wallets":   "WALLET",
	"PATCH wallets":  "WALLET",
	"DELETE wallets": "WALLET",

	// 	ROLE
	"POST roles":      "ROLE",
	"PATCH roles":     "JOBROLE",
	"POST job_role":   "JOBROLE",
	"PATCH job_role":  "JOBROLE",
	"DELETE job_role": "JOBROLE",

	// Customer Segmentation
	"POST customer-segmentations":   "CUSTOMERSEGMENTATIONS",
	"PATCH customer-segmentations":  "CUSTOMERSEGMENTATIONS",
	"DELETE customer-segmentations": "CUSTOMERSEGMENTATIONS",

	// CPS Roles
	"POST cps-roles":  "CPSROLES",
	"PATCH cps-roles": "CPSROLES",
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
			if strings.HasPrefix(relPath, "/api/v1/cbesuperapp/cps_action/") {
				relPath = strings.TrimPrefix(relPath, "/api/v1/cbesuperapp/cps_action/")
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
			found := false
			if strings.Contains(relPath, "news/category") {
				actionName = "NEWSCATEGORY"
				found = true
			} else if strings.Contains(relPath, "news/tag") {
				actionName = "NEWSTAG"
				found = true
			}

			if !found {
				if rparts := strings.Split(relPattern, "/"); rparts != nil && len(rparts) >= 1 {
					relPattern = rparts[0]
				}
				path := method + " " + relPath
				for k, v := range cpsActionRegistry {
					if strings.EqualFold(path, k) {
						actionName = v
						break
					}
				}
			}

			roleCode, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
			// roleID := utils.FirstHex24(rawRoleID)
			if roleCode == "" {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			action := strings.ToUpper(strings.TrimSpace(actionName))
			cacheKey := roleCode + ":" + action
			if ent, ok := cpsGuardCache.get(cacheKey); ok && ent.allow {
				next.ServeHTTP(w, r)
				return
			}

			allowed, err := cpsApproveRepo.ExistsByRoleAndAction(r.Context(), roleCode, action)
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

func actionKey(method, path string) string {
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return method
	}

	// METHOD + /first-segment
	return method + " /" + parts[0]
}
