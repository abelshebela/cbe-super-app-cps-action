package action_role_repo

import (
	"context"
	"strings"

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
