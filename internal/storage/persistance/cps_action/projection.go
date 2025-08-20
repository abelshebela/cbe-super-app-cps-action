package cps_action

import (
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

	// Helper to add slice fields
	addSlice := func(key string, value []string) {
		if len(value) > 0 {
			filter[key] = value
		}
	}

	// Helper to add time fields
	addTime := func(key string, value time.Time) {
		if !value.IsZero() {
			filter[key] = value
		}
	}

	// Helper to add pointer time fields
	addPtrTime := func(key string, value *time.Time) {
		if value != nil && !value.IsZero() {
			filter[key] = value
		}
	}

	if !cps.ID.IsZero() {
		filter["_id"] = cps.ID
	}
	addString("action_code", cps.ActionCode)
	addString("unique_id", cps.UniqueId)
	addString("maker_id", cps.MakerID)
	addString("maker_name", cps.MakerName)
	addString("maker_phone_number", cps.MakerPhoneNumber)
	addString("checker_id", cps.CheckerID)
	addString("checker_name", cps.CheckerName)
	addString("checker_phone_number", cps.CheckerPhoneNumber)
	addString("department", cps.Department)
	addString("rejection_reason", cps.RejectionReason)
	if cps.PreviousAction != nil {
		filter["previous_action"] = cps.PreviousAction
	}
	if cps.CurrentAction != nil {
		filter["current_action"] = cps.CurrentAction
	}
	addString("action_status", cps.ActionStatus)
	addString("action_type", cps.ActionType)
	if cps.IsDeleted {
		filter["is_deleted"] = cps.IsDeleted
	}
	addString("request_action", cps.RequestAction)
	addSlice("branch_codes", cps.BranchCodes)
	addSlice("branch_names", cps.BranchNames)
	addSlice("city_codes", cps.CityCodes)
	addSlice("city_names", cps.CityNames)
	addSlice("region_codes", cps.RegionCodes)
	addSlice("region_names", cps.RegionNames)
	addSlice("district_codes", cps.DistrictCodes)
	addSlice("district_names", cps.DistrictNames)
	addTime("created_at", cps.CreatedAt)
	addTime("last_modified_at", cps.LastModifiedAt)
	addTime("maker_action_time", cps.MakerActionTime)
	addPtrTime("checker_action_time", cps.CheckerActionTime)
	addString("reason", cps.Reason)
	addString("status", cps.Status)
	addTime("last_updated", cps.LastUpdated)

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

	// Helper to add slice fields
	addSlice := func(key string, value []string) {
		if len(value) > 0 {
			update[key] = value
		}
	}

	// Helper to add time fields
	addTime := func(key string, value time.Time) {
		if !value.IsZero() {
			update[key] = value
		}
	}

	// Helper to add pointer time fields
	addPtrTime := func(key string, value *time.Time) {
		if value != nil && !value.IsZero() {
			update[key] = value
		}
	}

	addString("action_code", cps.ActionCode)
	addString("unique_id", cps.UniqueId)
	addString("maker_id", cps.MakerID)
	addString("maker_name", cps.MakerName)
	addString("maker_phone_number", cps.MakerPhoneNumber)
	addString("checker_id", cps.CheckerID)
	addString("checker_name", cps.CheckerName)
	addString("checker_phone_number", cps.CheckerPhoneNumber)
	addString("department", cps.Department)
	addString("rejection_reason", cps.RejectionReason)
	if cps.PreviousAction != nil {
		update["previous_action"] = cps.PreviousAction
	}
	if cps.CurrentAction != nil {
		update["current_action"] = cps.CurrentAction
	}
	addString("action_status", cps.ActionStatus)
	addString("action_type", cps.ActionType)
	if cps.IsDeleted {
		update["is_deleted"] = cps.IsDeleted
	}
	addString("request_action", cps.RequestAction)
	addSlice("branch_codes", cps.BranchCodes)
	addSlice("branch_names", cps.BranchNames)
	addSlice("city_codes", cps.CityCodes)
	addSlice("city_names", cps.CityNames)
	addSlice("region_codes", cps.RegionCodes)
	addSlice("region_names", cps.RegionNames)
	addSlice("district_codes", cps.DistrictCodes)
	addSlice("district_names", cps.DistrictNames)
	addTime("created_at", cps.CreatedAt)
	addTime("last_modified_at", cps.LastModifiedAt)
	addTime("maker_action_time", cps.MakerActionTime)
	addPtrTime("checker_action_time", cps.CheckerActionTime)
	addString("reason", cps.Reason)
	addString("status", cps.Status)
	addTime("last_updated", cps.LastUpdated)

	return update
}
