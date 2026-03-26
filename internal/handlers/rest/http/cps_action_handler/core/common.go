package cpsactioncore

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	mid "cbe-super-app-cps-action/internal/handlers/middleware"
	cpsactionsvc "cbe-super-app-cps-action/internal/service/cps_action"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// CheckerValidationResult holds the output of ValidateCheckerAccess so the
// caller can build the final update payload without re-fetching anything.
type CheckerValidationResult struct {
	IdxDoc  *imodel.CPSActionApproveIndex
	RoleID  string
	Ctx     context.Context
	Request *http.Request
}

// ValidateCheckerAccess performs the common checker/approver role validation
// used by both ApproveCPSAction and RejectCPSAction. It resolves the module
// name, loads the approver index, validates ordering, and checks for duplicate
// or self-approval. Returns a non-nil error string when validation fails (the
// string is already a user-facing message suitable for SendBadRequestResponse).
func ValidateCheckerAccess(
	ctx context.Context,
	r *http.Request,
	action *model.CPSAction,
	userData *types.UserContext,
) (*CheckerValidationResult, string) {

	currentIndex := action.CurrentCheckerIndex

	// Resolve module name from request_action
	actionName := ""
	if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(action.RequestAction)); ok {
		actionName = mod
	}

	repo := mid.GetCPSActionApproveRepo()
	if repo == nil || actionName == "" {
		return nil, localization.ErrorOperationNotAllowed.Message
	}

	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		return nil, localization.ErrorOperationNotAllowed.Message
	}

	roleID := rawRoleID
	upperAction := strings.ToUpper(actionName)
	idxDoc, err := repo.FindByRoleAndAction(ctx, roleID, upperAction, action.Version)
	if err != nil {
		return nil, localization.ErrorOperationNotAllowed.Message
	}
	if idxDoc == nil || idxDoc.CheckerIndex == nil {
		return nil, localization.ErrorOperationNotAllowed.Message
	}

	roleLevel := *idxDoc.CheckerIndex
	expected := int32(*idxDoc.CheckerIndex)
	ctx = context.WithValue(ctx, constants.ContextKey("role_checker_index"), *idxDoc.CheckerIndex)
	ctx = context.WithValue(ctx, constants.ContextKey("role_checker_group"), expected)
	r = r.WithContext(ctx)

	if currentIndex == float64(roleLevel) {
		return nil, localization.MsgCPSActionApprovedByThisRole
	}
	if int64(currentIndex)+1 < int64(roleLevel) {
		return nil, localization.MsgCPSActionWaitPrevious
	}
	if action.MakerID == userData.UserName {
		return nil, localization.ErrorOperationNotAllowed.Message
	}

	for _, cu := range action.CheckerUsers {
		if cu.RoleID == roleID || cu.CheckerID == userData.UserID {
			return nil, localization.ErrorOperationNotAllowed.Message
		}
	}

	return &CheckerValidationResult{
		IdxDoc:  idxDoc,
		RoleID:  roleID,
		Ctx:     ctx,
		Request: r,
	}, ""
}

// CheckActionFinalized returns a user-facing message if the action is already
// in a terminal status (Approved, Rejected, Canceled). Returns "" when the
// action is still actionable.
func CheckActionFinalized(action *model.CPSAction) string {
	switch action.ActionStatus {
	case string(constants.Approved):
		return localization.MsgCPSActionAlreadyApproved
	case string(constants.Rejected):
		return localization.MsgCPSActionAlreadyRejected
	case string(constants.Canceled):
		return localization.MsgCPSActionAlreadyCanceled
	}
	return ""
}

// ResolveRequestActions maps a list of module names (e.g. ["BANK", "TOPUP"])
// to their constituent request action strings using RequestActionGroups.
// Duplicates are removed.
func ResolveRequestActions(modules []string) []string {
	var reqs []string
	seen := map[string]struct{}{}
	for _, mod := range modules {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := cpsactionsvc.RequestActionGroups[upper]; ok {
			for _, ra := range lst {
				key := string(ra)
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				reqs = append(reqs, key)
			}
		}
	}
	return reqs
}

// ValidateFilterParams extracts filter params from the request and validates
// the search and filter query parameters for special characters. Returns the
// filter and an error code string if validation fails ("" on success).
func ValidateFilterParams(r *http.Request) (*types.Filter, string) {
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		return nil, err.Error()
	}
	if err := local_util.NoSpecialChars(filter); err != nil {
		return nil, err.Error()
	}
	return filterParams, ""
}

// EnsureFilterParams guarantees that filterParams and its Filters map are non-nil,
// then injects the request_action $in filter when reqs is non-empty.
func EnsureFilterParams(filterParams *types.Filter, reqs []string) *types.Filter {
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	if len(reqs) > 0 {
		filterParams.Filters["request_action"] = map[string]interface{}{"$in": reqs}
	}
	return filterParams
}

// ModInfo describes a module's request actions and derived action types.
type ModInfo struct {
	RequestActions []string `json:"request_actions"`
	ActionTypes    []string `json:"action_types"`
}

// DeriveActionTypes inspects a list of request action strings and returns
// the unique CRUD action types they imply (CREATE, UPDATE, DELETE, ENABLE, DISABLE).
func DeriveActionTypes(actions []string) []string {
	seen := map[string]struct{}{}
	for _, ra := range actions {
		u := strings.ToUpper(ra)
		switch {
		case strings.Contains(u, "CREATE"):
			seen["CREATE"] = struct{}{}
		case strings.Contains(u, "UPDATE"):
			seen["UPDATE"] = struct{}{}
		case strings.Contains(u, "DELETE"):
			seen["DELETE"] = struct{}{}
		case strings.Contains(u, "ENABLE"):
			seen["ENABLE"] = struct{}{}
		case strings.Contains(u, "DISABLE"):
			seen["DISABLE"] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

// BuildModuleAllocations takes a list of module names, resolves their request
// actions, derives action types, and returns a map keyed by upper-cased module
// name plus a deduplicated list of module names.
func BuildModuleAllocations(mods []string) (map[string]ModInfo, []string) {
	m := map[string]ModInfo{}
	uniq := map[string]struct{}{}
	for _, raw := range mods {
		mod := strings.ToUpper(strings.TrimSpace(raw))
		if mod == "" {
			continue
		}
		uniq[mod] = struct{}{}
		var reqs []string
		if group, ok := cpsactionsvc.RequestActionGroups[mod]; ok {
			reqs = make([]string, 0, len(group))
			for _, ga := range group {
				reqs = append(reqs, string(ga))
			}
		} else {
			reqs = []string{}
		}
		m[mod] = ModInfo{RequestActions: reqs, ActionTypes: DeriveActionTypes(reqs)}
	}
	list := make([]string, 0, len(uniq))
	for k := range uniq {
		list = append(list, k)
	}
	return m, list
}

// GetRoleCode extracts the role_code from the request context and returns it.
// Returns an empty string if not found.
func GetRoleCode(r *http.Request) string {
	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	return rawRoleID
}

// GetApproveRepo is a convenience wrapper around mid.GetCPSActionApproveRepo.
func GetApproveRepo() interface {
	FindByRoleAndAction(ctx context.Context, roleID string, actionName string, version int64) (*imodel.CPSActionApproveIndex, error)
	PopulateUserApproverAllocations(ctx context.Context, roleID string) ([]string, []string, []string, []string, []string, error)
} {
	return mid.GetCPSActionApproveRepo()
}

// ResolveModuleName resolves the module name from a request action string.
// Returns ("", false) if no mapping is found.
func ResolveModuleName(requestAction string) (string, bool) {
	return cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(requestAction))
}

// FormatCheckerIndex is a helper to avoid repeated nil-check + dereference.
func FormatCheckerIndex(idx *float64) string {
	if idx == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *idx)
}
