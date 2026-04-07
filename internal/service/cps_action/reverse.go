package cpsaction

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"

	// "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ReverseCPSAction performs an auditor-driven reversal of an APPROVED CPS action.
func (ca *cpsActionService) ReverseCPSAction(ctx context.Context, actionCode string) error {
	// 1) Fetch original action
	orig, err := ca.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		return err
	}
	if orig == nil {
		return localization.ErrorResourceNotFound
	}
	if strings.ToUpper(strings.TrimSpace(orig.ActionStatus)) != string(constants.Approved) {
		return localization.ErrorOperationNotAllowed
	}

	// 2) Determine reversed request action and payload
	req := strings.ToUpper(strings.TrimSpace(orig.RequestAction))
	reversedReq := req

	// Toggle ENABLE/DISABLE pair first
	if strings.Contains(req, "ENABLE") && !strings.Contains(req, "DISABLE") {
		reversedReq = strings.Replace(req, "ENABLE", "DISABLE", 1)
	} else if strings.Contains(req, "DISABLE") {
		reversedReq = strings.Replace(req, "DISABLE", "ENABLE", 1)
	}

	var payload any
	switch ActionType(strings.ToUpper(strings.TrimSpace(orig.ActionType))) {
	case ActionCreate:
		// CREATE -> DELETE
		if strings.Contains(reversedReq, "CREATE") {
			reversedReq = strings.Replace(reversedReq, "CREATE", "DELETE", 1)
		}
		payload = orig.CurrentAction
	case ActionUpdate, ActionEnable, ActionDisable:
		// Revert to previous state
		payload = orig.PreviousAction
	case ActionDelete:
		// DELETE -> CREATE (restore)
		if strings.Contains(reversedReq, "DELETE") {
			reversedReq = strings.Replace(reversedReq, "DELETE", "CREATE", 1)
		}
		payload = orig.PreviousAction
	default:
		payload = orig.PreviousAction
	}

	roleCode := ctx.Value(constants.ContextKey("role_code")).(string)
	// 3) Dispatch the reverse operation via module Authorize
	rev := &model.CPSAction{
		RequestAction: reversedReq,
		UniqueId:      orig.UniqueId,
		CurrentAction: payload,
		RoleCode:      roleCode,
	}

	// Ensure a valid UniqueId for modules that require it (e.g., UPDATE/DELETE/ENABLE/DISABLE)
	if id := strings.TrimSpace(rev.UniqueId); id == "" || local_util.FirstHex24(id) == "" {
		tryExtract := func(v any) string {
			b, _ := json.Marshal(v)
			return local_util.FirstHex24(string(b))
		}
		if cand := tryExtract(payload); cand != "" {
			rev.UniqueId = cand
		} else if cand := tryExtract(orig.CurrentAction); cand != "" {
			rev.UniqueId = cand
		} else if cand := tryExtract(orig.PreviousAction); cand != "" {
			rev.UniqueId = cand
		}
	}

	if _, err := ca.dispatcher.Authorize(ctx, rev); err != nil {
		return err
	}

	// 4) Record auditor metadata on the original action
	auditor := local_util.ExtractUserFromContext(ctx)
	roleID, _ := ctx.Value(constants.ContextKey("role_id")).(string)

	update := bson.M{
		"reversed_by_role_id": roleID,
		"reversed_by_id":      auditor.UserID,
		"reversed_by_name":    auditor.FullName,
		"reversed_at":         time.Now(),
		"action_status":       string(constants.Reversed),
		"action_type":         string(constants.UPDATE),
	}

	if err := ca.repo.UpdateCustome(ctx, bson.M{"action_code": actionCode}, update); err != nil {
		return err
	}
	return nil
}
