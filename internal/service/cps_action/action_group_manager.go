package cpsaction

import "strings"

// ActionGroupManager provides requestAction -> parent module resolution
// without touching the dispatcher. It mirrors the dispatcher's priority
// so that resolution remains consistent across the codebase.

type ActionGroupManager interface {
	// ResolveModuleFor resolves the parent module for a request action string.
	ResolveModuleFor(action string) (string, bool)
	// ResolveModuleForRA resolves the parent module for a RequestAction.
	ResolveModuleForRA(action RequestAction) (string, bool)
	// ListModules returns the known modules (priority first, then remaining).
	ListModules() []string
}

type actionGroupManager struct{}

// DefaultActionGroupManager is a ready-to-use manager instance.
var DefaultActionGroupManager ActionGroupManager = &actionGroupManager{}

// modulePriority mirrors the case order in dispatcher.go to ensure
// deterministic resolution when actions belong to multiple groups.
var modulePriority = []string{
	"Bank",
	"KYCVerifier",
	"Block",
	"Account",
	"Advert",
	"Service",
	"ServicesCatalog",
	"DeviceVersion",
	"Fayda",
	"MiniAppMerchant",
	"MiniApp",
	"HQ",
	"Password",
	"Permission",
	"UnlinkDevice",
	"Wallet",
	"Topup",
	"AmountBasedAuth",
	"BulkService",
	"Event",
	"Notification",
	"ProductCode",
	"BPSUser",
	"Avatar",
	"donationCategory",
	"donationCompany",
	"Donation",
	"Department",
	"CPSUser",
	"BankVault",
	"VaultGroupCategory",
	"article",
	"articleCategory",
	"short_video",
	"customer",
	"news_tag",
	"ActionRole",
	"news_category",
	"BudgetCategory",
	"MiniAppCategory",
	"CpsActionRole",
}

func (m *actionGroupManager) ResolveModuleFor(action string) (string, bool) {
	return ResolveModuleFor(action)
}

func (m *actionGroupManager) ResolveModuleForRA(action RequestAction) (string, bool) {
	return ResolveModuleForRA(action)
}

func (m *actionGroupManager) ListModules() []string {
	return ListModules()
}

// ResolveModuleFor resolves a request action string to a parent module name.
func ResolveModuleFor(action string) (string, bool) {
	return ResolveModuleForRA(RequestAction(action))
}

// ResolveModuleForRA resolves a RequestAction to a parent module name.
func ResolveModuleForRA(action RequestAction) (string, bool) {

	// First, check in priority order to mirror dispatcher behavior
	for _, mod := range modulePriority {
		if IsActionInGroup(action, mod) {
			return mod, true
		}
	}

	// Then, scan any remaining groups not explicitly prioritized
	for mod := range RequestActionGroups {
		// skip already-checked modules
		if contains(modulePriority, mod) {
			continue
		}

		if IsActionInGroup(action, mod) {
			return mod, true
		}
	}
	// Fallback: infer module from request action string patterns
	if mod, ok := fallbackModuleForRA(action); ok {
		return mod, true
	}
	return "", false
}

// ListModules returns known module names in deterministic order
// (priority first, then any additional groups not in the priority list).
func ListModules() []string {
	out := make([]string, 0, len(RequestActionGroups))
	seen := map[string]struct{}{}

	for _, mod := range modulePriority {
		if _, ok := RequestActionGroups[mod]; ok {
			out = append(out, mod)
			seen[mod] = struct{}{}
		}
	}
	for mod := range RequestActionGroups {
		if _, ok := seen[mod]; ok {
			continue
		}
		out = append(out, mod)
	}
	return out
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

// fallbackModuleForRA tries to infer a module name from the action string when
// it isn't mapped in RequestActionGroups. Best-effort, case-insensitive.
func fallbackModuleForRA(action RequestAction) (string, bool) {
	s := strings.ToUpper(string(action))
	switch {
	case strings.Contains(s, "WALLET"):
		return "Wallet", true
	case strings.Contains(s, "TOPUP"):
		return "Topup", true
	case strings.Contains(s, "BANK_VAULT"):
		return "BankVault", true
	case strings.Contains(s, "BANK"):
		return "Bank", true
	case strings.Contains(s, "KYC"):
		return "KYCVerifier", true
	case strings.Contains(s, "FAYDA"):
		return "Fayda", true
	case strings.Contains(s, "MINI_APP_MERCHANT"):
		return "MiniAppMerchant", true
	case strings.Contains(s, "MINI_APP_CATEGORY"):
		return "MiniAppCategory", true
	case strings.Contains(s, "MINI_APP"):
		return "MiniApp", true
	case strings.Contains(s, "DEVICE_VERSION"):
		return "DeviceVersion", true
	case strings.Contains(s, "PERMISSION"):
		return "Permission", true
	case strings.Contains(s, "PASSWORD"):
		return "Password", true
	case strings.Contains(s, "AMOUNT_BASED_AUTH") || strings.Contains(s, "AUTHTIER"):
		return "AmountBasedAuth", true
	case strings.Contains(s, "NOTIFICATION"):
		return "Notification", true
	case strings.Contains(s, "AVATAR"):
		return "Avatar", true
	case strings.Contains(s, "DONATION_CATEGORY"):
		return "donationCategory", true
	case strings.Contains(s, "DONATION_COMPANY"):
		return "donationCompany", true
	case strings.Contains(s, "DONATION"):
		return "Donation", true
	case strings.Contains(s, "DEPARTMENT"):
		return "Department", true
	case strings.Contains(s, "CPS_USER"):
		return "CPSUser", true
	case strings.Contains(s, "VAULT_GROUP_CATEGORY"):
		return "VaultGroupCategory", true
	case strings.Contains(s, "ARTICLE_CATEGORY"):
		return "articleCategory", true
	case strings.Contains(s, "ARTICLE"):
		return "article", true
	case strings.Contains(s, "SHORT_VIDEO"):
		return "short_video", true
	case strings.Contains(s, "CUSTOMER"):
		return "customer", true
	case strings.Contains(s, "NEWS_TAG"):
		return "news_tag", true
	case strings.Contains(s, "NEWS_CATEGORY"):
		return "news_category", true
	case strings.Contains(s, "BUDGET_CATEGORY"):
		return "BudgetCategory", true
	case strings.Contains(s, "SERVICE_FEE") || strings.Contains(s, "DAILY_LIMIT") || strings.Contains(s, "MINIMUM") || strings.Contains(s, "TOTAL") || strings.Contains(s, "ACCESS_CONFIG"):
		return "Service", true
	case strings.Contains(s, "SERVICE"):
		return "ServicesCatalog", true
	case strings.Contains(s, "PRODUCT_CODE"):
		return "ProductCode", true
	case strings.Contains(s, "EVENT"):
		return "Event", true
	case strings.Contains(s, "BULK_SERVICE"):
		return "BulkService", true
	case strings.Contains(s, "UNLINK"):
		return "UnlinkDevice", true
	case strings.Contains(s, "ACTION_ROLE"):
		return "ActionRole", true
	case strings.Contains(s, "CPS_ACTION_ROLE"):
		return "CpsActionRole", true
	}
	return "", false
}
