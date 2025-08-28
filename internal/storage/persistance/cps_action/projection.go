package cps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var Projection = bson.M{
	"action_code":        1,
	"action_name":        1,
	"action_description": 1,
	"action_type":        1,
	"action_status":      1,
	"action_created_at":  1,
	"action_updated_at":  1,
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
	addString("department", cps.Department)
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

	addString("action_code", cps.ActionCode)
	addString("action_status", cps.ActionStatus)
	addString("checker_id", cps.CheckerID)
	addString("checker_name", cps.CheckerName)
	addString("rejection_reason", cps.RejectionReason)

	addString("checker_phone_number", cps.CheckerPhoneNumber)
	addTime("last_modified_at", time.Now())
	addTime("checker_action_time", time.Now())

	return update
}
