package bps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	bps_action "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// objectIDOrStringToHex converts ObjectID or string to hex string
func objectIDOrStringToHex(v interface{}) string {
	switch val := v.(type) {
	case bson.ObjectID:
		return val.Hex()
	case string:
		return val
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// objectIDToObjectID converts ObjectID or string to bson.ObjectID
func objectIDToObjectID(v interface{}) bson.ObjectID {
	switch val := v.(type) {
	case bson.ObjectID:
		return val
	case string:
		if objID, err := bson.ObjectIDFromHex(val); err == nil {
			return objID
		}
		return bson.NilObjectID
	case nil:
		return bson.NilObjectID
	default:
		return bson.NilObjectID
	}
}

// mapBpsActionToEntity maps the BPSActionDocument to local BPSAction entity
func mapBpsActionToEntity(dbAction types.BPSActionDocument) bps_action.BPSAction {
	// decode helper for interface{} JSON/BSON stored fields
	decode := func(v interface{}) map[string]interface{} {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
		if m, ok := v.(bson.M); ok {
			return map[string]interface{}(m)
		}
		if d, ok := v.(bson.D); ok {
			result := make(map[string]interface{})
			for _, elem := range d {
				result[elem.Key] = elem.Value
			}
			return result
		}
		switch val := v.(type) {
		case []byte:
			var out map[string]interface{}
			// try JSON first
			if err := json.Unmarshal(val, &out); err == nil {
				return out
			}
			// then BSON
			if err := bson.Unmarshal(val, &out); err == nil {
				return out
			}
			return nil
		case string:
			var out map[string]interface{}
			if err := json.Unmarshal([]byte(val), &out); err == nil {
				return out
			}
			return nil
		case bson.Raw:
			var out map[string]interface{}
			if err := bson.Unmarshal(val, &out); err == nil {
				return out
			}
			return nil
		case bson.Binary:
			var out map[string]interface{}
			if err := json.Unmarshal(val.Data, &out); err == nil {
				return out
			}
			if err := bson.Unmarshal(val.Data, &out); err == nil {
				return out
			}
			return nil
		default:
			return nil
		}
	}

	var currentAction map[string]interface{}
	if dbAction.CurrentAction != nil {
		// Handle bson.Binary case specifically for CurrentAction
		if currentActionBytes, ok := dbAction.CurrentAction.(bson.Binary); ok {
			err := json.Unmarshal(currentActionBytes.Data, &currentAction)
			if err != nil {
				fmt.Println("JSON unmarshal error for CurrentAction:", err)
				currentAction = nil
			}
		} else {
			currentAction = decode(dbAction.CurrentAction)
		}
	}

	var previousAction map[string]interface{}
	if dbAction.PreviousAction != nil {
		previousAction = decode(dbAction.PreviousAction)
	}

	// Convert ObjectID and handle flexible ID types
	actionID := dbAction.ID
	userID := objectIDOrStringToHex(dbAction.UserInformation.UserID)
	businessIDObjectID := objectIDToObjectID(dbAction.BusinessInformation.BusinessID)

	// Use time.Time values directly
	var verifiedAt time.Time
	if dbAction.VerifiedAt != nil {
		verifiedAt = *dbAction.VerifiedAt
	}

	createdAt := dbAction.CreatedAt
	lastModifiedAt := dbAction.LastModifiedAt

	return bps_action.BPSAction{
		ID:                actionID,
		ActionCode:        dbAction.ActionCode,
		IsAuditorApproved: dbAction.IsAuditorApproved,
		UserInformation: struct {
			UserID         string   `json:"user_id,omitempty" bson:"user_id,omitempty"`
			UserCode       string   `json:"user_code" bson:"user_code"`
			FullName       string   `json:"full_name" bson:"full_name"`
			AccountNumbers []string `json:"account_numbers" bson:"account_numbers"`
			PhoneNumbers   string   `json:"phone_numbers" bson:"phone_numbers"`
			BranchCode     string   `json:"branch_code" bson:"branch_code"`
		}{
			UserID:         userID,
			UserCode:       dbAction.UserInformation.UserCode,
			FullName:       dbAction.UserInformation.FullName,
			AccountNumbers: dbAction.UserInformation.AccountNumbers,
			PhoneNumbers:   dbAction.UserInformation.PhoneNumbers,
			BranchCode:     dbAction.UserInformation.BranchCode,
		},
		BusinessInformation: struct {
			BusinessID   bson.ObjectID `json:"business_id" bson:"business_id"`
			TILLNumber   string        `json:"till_number" bson:"till_number"`
			BusinessName string        `json:"business_name" bson:"business_name"`
		}{
			BusinessID:   businessIDObjectID,
			TILLNumber:   dbAction.BusinessInformation.TILLNumber,
			BusinessName: dbAction.BusinessInformation.BusinessName,
		},
		CheckersNeeded:     dbAction.CheckersNeeded,
		CheckersApproved:   dbAction.CheckersApproved,
		CheckerID:          dbAction.CheckerID,
		MakerID:            dbAction.MakerID,
		MakerName:          dbAction.MakerName,
		MakerReason:        dbAction.MakerReason,
		MakerPhoneNumber:   dbAction.MakerPhoneNumber,
		CheckerName:        dbAction.CheckerName,
		CheckerPhoneNumber: dbAction.CheckerPhoneNumber,
		ActionReason: struct {
			ActionType constants.ActionType `json:"action_type" bson:"action_type"` //  reject, enable, disable
			ActionNote string               `json:"action_note" bson:"action_note"`
			Identifier string               `json:"identifier" bson:"identifier"`
		}{
			ActionType: constants.ActionType(dbAction.ActionReason.ActionType),
			ActionNote: dbAction.ActionReason.ActionNote,
			Identifier: dbAction.ActionReason.Identifier,
		},
		MakerMID:        dbAction.MakerMID,
		CheckerMID:      dbAction.CheckerMID,
		AuditorMID:      dbAction.AuditorMID,
		CheckerNameList: dbAction.CheckerNameList,
		AuditorNameList: dbAction.AuditorNameList,
		Auditors: struct {
			AuditorName        string   `json:"auditor_name" bson:"auditor_name"`
			AuditorPhoneNumber string   `json:"auditor_phone_number" bson:"auditor_phone_number"`
			AuditorsRequired   int      `json:"auditors_required" bson:"auditors_required"`
			AuditorID          []string `json:"auditor_id" bson:"auditor_id"`
			Audited            bool     `json:"audited" bson:"audited"`
			AuditorApproval    bool     `json:"auditor_approval" bson:"auditor_approval"`
			Reason             string   `json:"reason" bson:"reason"`
		}{
			AuditorName:        dbAction.Auditors.AuditorName,
			AuditorPhoneNumber: dbAction.Auditors.AuditorPhoneNumber,
			AuditorsRequired:   dbAction.Auditors.AuditorsRequired,
			AuditorID:          dbAction.Auditors.AuditorID,
			Audited:            dbAction.Auditors.Audited,
			AuditorApproval:    dbAction.Auditors.AuditorApproval,
			Reason:             dbAction.Auditors.Reason,
		},
		CheckerTime:        dbAction.CheckerTime,
		AuditorTime:        dbAction.AuditorTime,
		RequestAction:      constants.RequestAction(dbAction.RequestAction),
		EntityIdentifyer:   dbAction.EntityIdentifyer,
		HomeBranch:         dbAction.HomeBranch,
		AccountBranchCode:  dbAction.AccountBranchCode,
		DistrictCode:       dbAction.DistrictCode,
		BranchCode:         dbAction.BranchCode,
		LinkedDistrictCode: dbAction.LinkedDistrictCode,
		AccountNumber:      dbAction.AccountNumber,
		AccountHolderName:  dbAction.AccountHolderName,
		ServiceName:        dbAction.ServiceName,
		CurrentAction:      currentAction,
		PreviousAction:     previousAction,
		VerifiedAt:         &verifiedAt,
		Status:             dbAction.Status,
		CreatedAt:          createdAt,
		LastModifiedAt:     lastModifiedAt,
	}
}
