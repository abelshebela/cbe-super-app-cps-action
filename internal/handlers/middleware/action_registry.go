package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
)

var cpsActionRegistry = map[string]string{
	// Notification
	"POST notifications":               "NOTIFICATIONS",
	"PATCH notifications":              "NOTIFICATIONS",
	"DELETE notifications":             "NOTIFICATIONS",
	"GET notifications":                "NOTIFICATIONS",
	"GET notifications/{id}":           "NOTIFICATIONS",
	"PATCH notifications/enable/{id}":  "NOTIFICATIONS",
	"PATCH notifications/disable/{id}": "NOTIFICATIONS",

	// MiniAppMerchant
	"POST mini-app-merchants":   "MINIAPPMERCHANT",
	"PATCH mini-app-merchants":  "MINIAPPMERCHANT",
	"DELETE mini-app-merchants": "MINIAPPMERCHANT",

	// Advert
	"POST adverts":  "ADVERT",
	"PATCH adverts": "ADVERT",
	"DELETE advert": "ADVERT",

	// AccountBlock - Branch Operations
	"GET account_block/branches":             "ACCOUNTBLOCK",
	"GET account_block/branches/{branch_id}": "ACCOUNTBLOCK",
	"POST account_block/branches/enable":     "SINGLEBRANCHENABLEACCOUNTBLOCK",
	"POST account_block/branches/disable":    "SINGLEBRANCHDISABLEACCOUNTBLOCK",

	// AccountBlock - Region Operations
	"GET account_block/regions":             "ACCOUNTBLOCK",
	"GET account_block/regions/{region_id}": "ACCOUNTBLOCK",
	"POST account_block/regions/enable":     "MULTIBRANCHENABLEACCOUNTBLOCK",
	"POST account_block/regions/disable":    "MULTIBRANCHDISABLEACCOUNTBLOCK",

	// AccountBlock - District Operations
	"GET account_block/districts":               "ACCOUNTBLOCK",
	"GET account_block/districts/{district_id}": "ACCOUNTBLOCK",
	"POST account_block/districts/enable":       "MULTIBRANCHENABLEACCOUNTBLOCK",
	"POST account_block/districts/disable":      "MULTIBRANCHDISABLEACCOUNTBLOCK",

	// AccountBlock - Details
	"GET account_block/details/{id}": "ACCOUNTBLOCK",

	// AccountValidation
	"GET account_validation":   "ACCOUNTVALIDATION",
	"PATCH account_validation": "ACCOUNTVALIDATION",

	// AmountBasedAuth
	"PATCH amount_based_auth": "AMOUNTBASEDAUTH",
	"GET amount_based_auth":   "AMOUNTBASEDAUTH",

	// Banks
	"GET banks":                "BANK",
	"GET banks/{id}":           "BANK",
	"POST banks":               "BANK",
	"PATCH banks":              "BANK",
	"DELETE banks":             "BANK",
	"PATCH banks/{id}/enable":  "BANK",
	"PATCH banks/{id}/disable": "BANK",
	"PATCH banks/{id}/logo":    "BANK",

	// Wallet
	"GET wallets":                "WALLET",
	"GET wallets/{id}":           "WALLET",
	"POST wallets":               "WALLET",
	"PATCH wallets":              "WALLET",
	"DELETE wallets":             "WALLET",
	"PATCH wallets/{id}/enable":  "WALLET",
	"PATCH wallets/{id}/disable": "WALLET",

	// Roles (GET operations)
	"GET roles":     "ROLE",
	"GET job_roles": "JOBROLE",

	// Services (GET operations)
	"GET services": "SERVICE",

	// Topup (GET operations)
	"GET topups":      "TOPUP",
	"GET topups/{id}": "TOPUP",

	// Customers
	"GET customers":      "CUSTOMER",
	"GET customers/{id}": "CUSTOMER",
	"PATCH customers":    "CUSTOMER",

	// Departments
	"GET departments":      "DEPARTMENT",
	"GET departments/{id}": "DEPARTMENT",
	"POST departments":     "DEPARTMENT",
	"PATCH departments":    "DEPARTMENT",

	// Events
	"GET events":      "EVENT",
	"GET events/{id}": "EVENT",
	"POST events":     "EVENT",
	"PATCH events":    "EVENT",
	"DELETE events":   "EVENT",

	// CPS Users
	"GET cps_users":                       "CPSUSER",
	"GET cps_users/{user_code}":           "CPSUSER",
	"GET cps_users/code/{code}":           "CPSUSER",
	"POST cps_users/create":               "CPSUSER",
	"PATCH cps_users/update/{user_code}":  "CPSUSER",
	"DELETE cps_users/delete/{user_code}": "CPSUSER",
	"POST cps_users/disable/{user_code}":  "CPSUSER",
	"POST cps_users/enable/{user_code}":   "CPSUSER",

	// BPS Users (GET operations)
	"GET bps_users":      "BPSUSER",
	"GET bps_users/{id}": "BPSUSER",

	// Avatar
	"POST avatar":   "AVATAR",
	"DELETE avatar": "AVATAR",
	"PATCH avatar":  "AVATAR",

	// BankVault
	"POST vault":   "BANKVAULT",
	"PATCH vault":  "BANKVAULT",
	"DELETE vault": "BANKVAULT",

	// ActionRole (BPS)

	"POST bps-action-roles":  "BPSACTIONEROLE",
	"PATCH bps-action-roles": "BPSACTIONEROLE",

	// BPSUser
	"POST bps_users":       "BPSUSER",
	"PATCH bps_users/{id}": "BPSUSER",

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

	// Logistics merchant
	"POST logistics_merchants":  "LOGISTICSMERCHANT",
	"PATCH logistics_merchants": "LOGISTICSMERCHANT",
	"DELETE logistics_merchant": "LOGISTICSMERCHANT",

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
	"POST access_list_segmentation":  "ACCESSLISTSEGMENTATION",
	"PATCH access_list_segmentation": "ACCESSLISTSEGMENTATION",

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
	"POST topups":                     "TOPUP",
	"PATCH topups":                    "TOPUP",
	"DELETE topups":                   "TOPUP",
	"PATCH topups/{topup_id}/disable": "TOPUP",
	"PATCH topups/{topup_id}/enable":  "TOPUP",

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
			// Security: Fail-closed - deny access if repo is not initialized
			if cpsApproveRepo == nil {
				if guardLogger != nil {
					guardLogger.Errorf("[ActionRegistry][RouteGuard] security: CPSActionApproveRepo is not initialized - denying access")
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			rc := chi.RouteContext(r.Context())
			if rc == nil {
				if guardLogger != nil {
					guardLogger.Errorf("[ActionRegistry][RouteGuard] security: route context is nil - denying access")
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			pattern := routeFullPattern(rc)
			if pattern == "" {
				if guardLogger != nil {
					guardLogger.Errorf("[ActionRegistry][RouteGuard] security: route pattern is empty - denying access")
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
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
			relPath := strings.TrimPrefix(r.URL.Path, "/api/v1/cbesuperapp/cps_action/")

			// Allowlist (e.g., CPSAction endpoints)
			for _, p := range whitelist {
				if strings.HasPrefix(relPath, p) {
					next.ServeHTTP(w, r)
					return
				}
			}

			method := strings.ToUpper(r.Method)

			// Use the GetActionNameFromPath function to resolve action name from registry
			actionName := GetActionNameFromPath(method, r.URL.Path)
			if actionName == "" {
				if guardLogger != nil {
					guardLogger.Errorf("[ActionRegistry][RouteGuard] no action name found for %s %s - denying access", method, r.URL.Path)
				}
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			roleCode, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
			// roleID := utils.FirstHex24(rawRoleID)
			if roleCode == "" {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}

			// Check if the user's job_title role is enabled
			if roleRepo != nil {
				roleCacheKey := "role_enabled:" + roleCode
				if ent, ok := roleEnabledCache.get(roleCacheKey); ok {
					if !ent.allow {
						localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
						return
					}
				} else {
					role, err := roleRepo.FindByRole(r.Context(), roleCode)
					if err != nil || role == nil {
						if guardLogger != nil {
							guardLogger.Errorf("[ActionRegistry][RouteGuard] role lookup err code: %s: %v", roleCode, err)
						}
						localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
						return
					}
					roleEnabledCache.set(roleCacheKey, allowEntry{allow: role.Enabled, exp: nowPlus(roleEnabledCache.ttl)})
					if !role.Enabled {
						localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
						return
					}
				}
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
					guardLogger.Errorf("[ActionRegistry][RouteGuard] guard lookup err: %v", err)
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

// GetActionNameFromPath returns the action name from the registry based on HTTP method and path
func GetActionNameFromPath(method, path string) string {
	// Normalize the path - remove base API prefix
	relPath := strings.TrimPrefix(path, "/api/v1/cbesuperapp/cps_action")
	if relPath == "" {
		relPath = "/"
	}

	// Handle special cases for news endpoints
	if strings.Contains(relPath, "news/category") {
		return "NEWSCATEGORY"
	}
	if strings.Contains(relPath, "news/tag") {
		return "NEWSTAG"
	}

	// Handle special cases for fayda endpoints
	if strings.Contains(relPath, "fayda_account") {
		return "FAYDA"
	}

	// Try exact match first (case-insensitive)
	keyPattern := strings.ToUpper(method) + " " + relPath
	for key, actionName := range cpsActionRegistry {
		if strings.EqualFold(keyPattern, key) {
			return actionName
		}
	}

	// Try pattern matching (handle dynamic segments like {id})
	for key, actionName := range cpsActionRegistry {
		if strictAvatarMatch(keyPattern, key) {
			return actionName
		}
	}

	// Try resource-based matching as fallback
	resource := extractResource(relPath)
	normalizedResource := normalize(resource)

	for key, actionName := range cpsActionRegistry {
		keyParts := strings.SplitN(key, " ", 2)
		if len(keyParts) == 2 && strings.EqualFold(keyParts[0], method) {
			keyResource := extractResource(keyParts[1])
			if normalize(keyResource) == normalizedResource {
				return actionName
			}
		}
	}

	return ""
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
