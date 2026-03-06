package action_role_repo

import (
	"context"
	"errors"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// FindByRoleAndAction returns a merged approver index document for the given roleID and actionName.
// Multiple documents may exist (maker, checker, auditor as separate rows), so this fetches all
// matching documents and merges checker_index and auditor_index into a single result.
func (r *CPSActionApproveIndexRepository) FindByRoleAndAction(ctx context.Context, roleID string, actionName string, version int64) (*imodel.CPSActionApproveIndex, error) {
	filter := bson.M{
		"role_id":     roleID,
		"action_name": strings.ToUpper(strings.TrimSpace(actionName)),
		"version":     version,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []imodel.CPSActionApproveIndex
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, nil
	}

	// Merge all documents into a single result, picking the first non-nil value for each index.
	merged := docs[0]
	for _, doc := range docs[1:] {
		if merged.MakerIndex == nil && doc.MakerIndex != nil {
			merged.MakerIndex = doc.MakerIndex
		}
		if merged.CheckerIndex == nil && doc.CheckerIndex != nil {
			merged.CheckerIndex = doc.CheckerIndex
		}
		if merged.AuditorIndex == nil && doc.AuditorIndex != nil {
			merged.AuditorIndex = doc.AuditorIndex
		}
	}
	return &merged, nil
}

// FindVersionsByActionName returns distinct versions for a given action_name from cps_action_approver_index.
func (r *CPSActionApproveIndexRepository) FindVersionsByActionName(ctx context.Context, actionName string) ([]int64, error) {
	r.logger.Infof("FindVersionsByActionName: actionName=%s", actionName)
	filter := bson.M{
		"action_name": strings.ToUpper(strings.TrimSpace(actionName)),
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("FindVersionsByActionName: Find failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	seen := make(map[int64]bool)
	var versions []int64
	for cursor.Next(ctx) {
		var doc struct {
			Version int64 `bson:"version"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		if !seen[doc.Version] {
			seen[doc.Version] = true
			versions = append(versions, doc.Version)
		}
	}
	r.logger.Infof("FindVersionsByActionName: found %d versions for actionName=%s", len(versions), actionName)
	return versions, nil
}

// UpdateRoleInIndices updates all cps_action_approver_index documents matching
// the given actionName and version, swapping oldRoleCode to newRoleCode.
func (r *CPSActionApproveIndexRepository) UpdateRoleInIndices(ctx context.Context, actionName string, version int64, oldRoleCode, newRoleCode string) (int64, error) {
	r.logger.Infof("UpdateRoleInIndices: actionName=%s version=%d old=%s new=%s", actionName, version, oldRoleCode, newRoleCode)
	filter := bson.M{
		"action_name": strings.ToUpper(strings.TrimSpace(actionName)),
		"version":     version,
		"role_id":     oldRoleCode,
	}
	update := bson.M{
		"$set": bson.M{
			"role_id":    newRoleCode,
			"updated_at": bson.M{"$currentDate": true},
		},
	}
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("UpdateRoleInIndices: UpdateMany failed: %v", err)
		return 0, errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("UpdateRoleInIndices: updated %d documents", result.ModifiedCount)
	return result.ModifiedCount, nil
}

// FindAllocationsWithVersions returns a map of action_name -> versions
// where the given role has the specified index type (e.g. "checker_index" or "auditor_index").
// This is used for version-aware fetch filtering of pending/unclaimed actions.
func (r *CPSActionApproveIndexRepository) FindAllocationsWithVersions(ctx context.Context, roleID string, indexField string) (map[string][]int64, error) {
	r.logger.Infof("FindAllocationsWithVersions: roleID=%s indexField=%s", roleID, indexField)
	filter := bson.M{
		"role_id":  roleID,
		indexField: bson.M{"$ne": nil},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Errorf("FindAllocationsWithVersions: Find failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var docs []imodel.CPSActionApproveIndex
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Errorf("FindAllocationsWithVersions: Cursor.All failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	result := make(map[string][]int64)
	for _, doc := range docs {
		seen := false
		for _, v := range result[doc.ActionName] {
			if v == doc.Version {
				seen = true
				break
			}
		}
		if !seen {
			result[doc.ActionName] = append(result[doc.ActionName], doc.Version)
		}
	}
	r.logger.Infof("FindAllocationsWithVersions: found %d action_name groups for roleID=%s", len(result), roleID)
	return result, nil
}
