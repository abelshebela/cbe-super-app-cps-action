package cps_action

import (
	"cbe-super-app-cps-action/internal/constants/model"

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

	if cps.ID.IsZero() {
		filter["_id"] = cps.ID
	}
	if cps.ActionCode != "" {
		filter["action_code"] = cps.ActionCode
	}
	if cps.UniqueId != "" {
		filter["unique_id"] = cps.UniqueId
	}
	if cps.MakerID != "" {
		filter["maker_id"] = cps.MakerID
	}
	if cps.MakerName != "" {
		filter["maker_name"] = cps.MakerName
	}
	if cps.MakerPhoneNumber != "" {
		filter["maker_phone_number"] = cps.MakerPhoneNumber
	}
	if cps.CheckerID != "" {
		filter["checker_id"] = cps.CheckerID
	}
	if cps.CheckerName != "" {
		filter["checker_name"] = cps.CheckerName
	}
	if cps.CheckerPhoneNumber != "" {
		filter["checker_phone_number"] = cps.CheckerPhoneNumber
	}
	if cps.Department != "" {
		filter["department"] = cps.Department
	}
	if cps.RejectionReason != "" {
		filter["rejection_reason"] = cps.RejectionReason
	}
	if cps.PreviousAction != nil {
		filter["previous_action"] = cps.PreviousAction
	}
	if cps.CurrentAction != nil {
		filter["current_action"] = cps.CurrentAction
	}
	if cps.ActionStatus != "" {
		filter["action_status"] = cps.ActionStatus
	}
	if cps.ActionType != "" {
		filter["action_type"] = cps.ActionType
	}
	// Only filter by is_deleted if it's true, to avoid filtering out non-deleted by default
	if cps.IsDeleted {
		filter["is_deleted"] = cps.IsDeleted
	}
	if cps.RequestAction != "" {
		filter["request_action"] = cps.RequestAction
	}
	// Branch fields
	if len(cps.BranchCodes) > 0 {
		filter["branch_codes"] = cps.BranchCodes
	}
	if len(cps.BranchNames) > 0 {
		filter["branch_names"] = cps.BranchNames
	}
	// City fields
	if len(cps.CityCodes) > 0 {
		filter["city_codes"] = cps.CityCodes
	}
	if len(cps.CityNames) > 0 {
		filter["city_names"] = cps.CityNames
	}
	// Region fields
	if len(cps.RegionCodes) > 0 {
		filter["region_codes"] = cps.RegionCodes
	}
	if len(cps.RegionNames) > 0 {
		filter["region_names"] = cps.RegionNames
	}
	// District fields
	if len(cps.DistrictCodes) > 0 {
		filter["district_codes"] = cps.DistrictCodes
	}
	if len(cps.DistrictNames) > 0 {
		filter["district_names"] = cps.DistrictNames
	}
	// Timestamps
	if !cps.CreatedAt.IsZero() {
		filter["created_at"] = cps.CreatedAt
	}
	if !cps.LastModifiedAt.IsZero() {
		filter["last_modified_at"] = cps.LastModifiedAt
	}
	if !cps.MakerActionTime.IsZero() {
		filter["maker_action_time"] = cps.MakerActionTime
	}
	if cps.CheckerActionTime != nil && !cps.CheckerActionTime.IsZero() {
		filter["checker_action_time"] = cps.CheckerActionTime
	}
	// Additional fields
	if cps.Reason != "" {
		filter["reason"] = cps.Reason
	}
	if cps.Status != "" {
		filter["status"] = cps.Status
	}
	if !cps.LastUpdated.IsZero() {
		filter["last_updated"] = cps.LastUpdated
	}

	return filter
}

func BuildCPSActionUpdateMap(cps model.CPSAction) bson.M {
	update := bson.M{}

	if cps.ActionCode != "" {
		update["action_code"] = cps.ActionCode
	}
	if cps.UniqueId != "" {
		update["unique_id"] = cps.UniqueId
	}
	if cps.MakerID != "" {
		update["maker_id"] = cps.MakerID
	}
	if cps.MakerName != "" {
		update["maker_name"] = cps.MakerName
	}
	if cps.MakerPhoneNumber != "" {
		update["maker_phone_number"] = cps.MakerPhoneNumber
	}
	if cps.CheckerID != "" {
		update["checker_id"] = cps.CheckerID
	}
	if cps.CheckerName != "" {
		update["checker_name"] = cps.CheckerName
	}
	if cps.CheckerPhoneNumber != "" {
		update["checker_phone_number"] = cps.CheckerPhoneNumber
	}
	if cps.Department != "" {
		update["department"] = cps.Department
	}
	if cps.RejectionReason != "" {
		update["rejection_reason"] = cps.RejectionReason
	}
	if cps.PreviousAction != nil {
		update["previous_action"] = cps.PreviousAction
	}
	if cps.CurrentAction != nil {
		update["current_action"] = cps.CurrentAction
	}
	if cps.ActionStatus != "" {
		update["action_status"] = cps.ActionStatus
	}
	if cps.ActionType != "" {
		update["action_type"] = cps.ActionType
	}
	if cps.IsDeleted {
		update["is_deleted"] = cps.IsDeleted
	}
	if cps.RequestAction != "" {
		update["request_action"] = cps.RequestAction
	}
	// Branch fields
	if len(cps.BranchCodes) > 0 {
		update["branch_codes"] = cps.BranchCodes
	}
	if len(cps.BranchNames) > 0 {
		update["branch_names"] = cps.BranchNames
	}
	// City fields
	if len(cps.CityCodes) > 0 {
		update["city_codes"] = cps.CityCodes
	}
	if len(cps.CityNames) > 0 {
		update["city_names"] = cps.CityNames
	}
	// Region fields
	if len(cps.RegionCodes) > 0 {
		update["region_codes"] = cps.RegionCodes
	}
	if len(cps.RegionNames) > 0 {
		update["region_names"] = cps.RegionNames
	}
	// District fields
	if len(cps.DistrictCodes) > 0 {
		update["district_codes"] = cps.DistrictCodes
	}
	if len(cps.DistrictNames) > 0 {
		update["district_names"] = cps.DistrictNames
	}
	// Timestamps
	if !cps.CreatedAt.IsZero() {
		update["created_at"] = cps.CreatedAt
	}
	if !cps.LastModifiedAt.IsZero() {
		update["last_modified_at"] = cps.LastModifiedAt
	}
	if !cps.MakerActionTime.IsZero() {
		update["maker_action_time"] = cps.MakerActionTime
	}
	if cps.CheckerActionTime != nil && !cps.CheckerActionTime.IsZero() {
		update["checker_action_time"] = cps.CheckerActionTime
	}
	// Additional fields
	if cps.Reason != "" {
		update["reason"] = cps.Reason
	}
	if cps.Status != "" {
		update["status"] = cps.Status
	}
	if !cps.LastUpdated.IsZero() {
		update["last_updated"] = cps.LastUpdated
	}

	return update
}
