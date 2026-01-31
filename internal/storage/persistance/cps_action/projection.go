package cps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var Projection = bson.M{
	"action_code":           1,
	"action_name":           1,
	"maker_id":              1,
	"maker_name":            1,
	"maker_phone_number":    1,
	"checker_users":         1,
	"checker_count":         1,
	"current_checker_index": 1,
	"action_description":    1,
	"action_type":           1,
	"action_status":         1,
	"auditor_status":        1,
	"auditor_users":         1,
	"auditor_count":         1,
	"current_auditor_index": 1,
	"request_action":        1,
	"action_created_at":     1,
	"maker_action_time":     1,
	"checker_action_time":   1,
	"previous_action":       1,
	"current_action":        1,
	"department":            1,
	"rejection_reason":      1,
	"unique_id":             1,
	"created_at":            1,
	"last_modified_at":      1,
	"reversed_by_role_id":   1,
	"reversed_by_id":        1,
	"reversed_by_name":      1,
	"reversed_at":           1,
}

func BuildCPSActionFilter(cps model.CPSAction) bson.M {
	filter := bson.M{}

	// Helper to add string fields
	addString := func(key, value string) {
		if value != "" {
			filter[key] = value
		}
	}

	if !cps.ID.IsZero() {
		filter["_id"] = cps.ID
	}
	addString("action_code", cps.ActionCode)
	// addString("department", cps.Department)
	addString("request_action", cps.RequestAction)
	addString("action_status", string(constants.Pending))

	return filter
}

func BuildCPSActionUpdateMap(cps model.CPSAction) bson.M {
	update := bson.M{}

	// Helper to add string fields
	addString := func(key, value string) {
		if value != "" {
			update[key] = value
		}
	}

	// Helper to add time fields
	addTime := func(key string, value time.Time) {
		if !value.IsZero() {
			update[key] = value
		}
	}

	addString("maker_id", cps.MakerID)
	addString("maker_name", cps.MakerName)
	addString("maker_phone_number", cps.MakerPhoneNumber)
	addString("action_code", cps.ActionCode)
	addString("action_status", cps.ActionStatus)
	addString("auditor_status", string(cps.AuditorStatus))
	addString("rejection_reason", cps.RejectionReason)
	// Multi-checker fields
	if cps.CheckerUsers != nil && len(cps.CheckerUsers) > 0 {
		update["checker_users"] = cps.CheckerUsers
	}
	if cps.AuditorUsers != nil && len(cps.AuditorUsers) > 0 {
		update["auditor_users"] = cps.AuditorUsers
	}
	if cps.CheckerCount > 0 {
		update["checker_count"] = cps.CheckerCount
	}
	if cps.CurrentCheckerIndex > 0 {
		update["current_checker_index"] = cps.CurrentCheckerIndex
	}
	if cps.AuditorCount > 0 {
		update["auditor_count"] = cps.AuditorCount
	}
	if cps.CurrentAuditorIndex > 0 {
		update["current_auditor_index"] = cps.CurrentAuditorIndex
	}
	addString("auditor_status", string(cps.AuditorStatus))
	addString("role_code", cps.RoleCode)
	addTime("last_modified_at", time.Now())
	addTime("checker_action_time", time.Now())

	return update
}

func BuildCPSActionFilterAuditor(cps imodel.CPSAction) bson.M {
	filter := bson.M{}

	// Helper to add string fields
	addString := func(key, value string) {
		if value != "" {
			filter[key] = value
		}
	}

	if !cps.ID.IsZero() {
		filter["_id"] = cps.ID
	}
	addString("action_code", cps.ActionCode)
	// addString("department", cps.Department)
	addString("request_action", cps.RequestAction)
	addString("action_status", string(constants.Pending))

	return filter
}

func BuildCPSActionUpdateMapAuditor(cps imodel.CPSAction) bson.M {
	update := bson.M{}

	// Helper to add string fields
	addString := func(key, value string) {
		if value != "" {
			update[key] = value
		}
	}

	// Helper to add time fields
	addTime := func(key string, value time.Time) {
		if !value.IsZero() {
			update[key] = value
		}
	}

	addString("maker_id", cps.MakerID)
	addString("maker_name", cps.MakerName)
	addString("maker_phone_number", cps.MakerPhoneNumber)
	addString("action_code", cps.ActionCode)
	addString("action_status", cps.ActionStatus)
	addString("rejection_reason", cps.RejectionReason)
	// Multi-checker fields
	if cps.CheckerUsers != nil && len(cps.CheckerUsers) > 0 {
		update["checker_users"] = cps.CheckerUsers
	}
	if cps.CheckerCount > 0 {
		update["checker_count"] = cps.CheckerCount
	}
	if cps.CurrentCheckerIndex > 0 {
		update["current_checker_index"] = cps.CurrentCheckerIndex
	}
	addString("role_code", cps.RoleCode)
	addTime("last_modified_at", time.Now())
	addTime("checker_action_time", time.Now())

	return update
}
