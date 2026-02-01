package bps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
)

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
	"BANK",
	"KYCVERIFIER",
	"BLOCK",
	"ACCOUNT",
	"ADVERT",
	"SERVICE",
	"SERVICESCATALOG",
	"DEVICEVERSION",
	"FAYDA",
	"MINIAPPMERCHANT",
	"MINIAPP",
	"HQ",
	"PASSWORD",
	"PERMISSION",
	"UNLINKDEVICE",
	"WALLET",
	"TOPUP",
	"AMOUNTBASEDAUTH",
	"BULKSERVICEALLUSER",
	"EVENT",
	"NOTIFICATION",
	"PRODUCTCODE",
	"BPSUSER",
	"AVATAR",
	"DONATIONCATEGORY",
	"DONATIONCOMPANY",
	"DONATION",
	"DEPARTMENT",
	"CPSUSER",
	"BANKVAULT",
	"VAULTGROUPCATEGORY",
	"ARTICLE",
	"ARTICLECATEGORY",
	"SHORTVIDEO",
	"CUSTOMER",
	"NEWSTAG",
	"ROLE",
	"ACTIONROLE",
	"NEWSCATEGORY",
	"BUDGETCATEGORY",
	"MINIAPPCATEGORY",
	"CPSACTIONROLE",
	"LOGISTICSMERCHANT",
	"EVENTMERCHANT",
	"USSDMERCHANT",
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
	if mod, ok := lib.FallbackModuleForRA(constants.RequestAction(action)); ok {
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
