package action_role_repo

import (
	"context"
	"strings"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// FindByRoleAndAction returns the approver index document for the given roleID and actionName
// Requires maker_index to exist and be non-null
func (r *BPSActionApproveIndexRepository) FindByRoleAndAction(ctx context.Context, roleID string, actionName string) (*imodel.BPSActionApproveIndex, error) {
	filter := bson.M{
		"role_id":     roleID,
		"action_name": strings.ToUpper(strings.TrimSpace(actionName)),
		// "maker_index": bson.M{"$exists": true, "$ne": nil},
	}
	var res imodel.BPSActionApproveIndex
	err := r.collection.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		r.logger.Errorf("[FindByRoleAndAction] failed to find approver index: %v", err)
		return nil, err
	}
	return &res, nil
}
