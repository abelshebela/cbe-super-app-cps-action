package action_role_repo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ExistsByRoleAndAction returns true if a record exists for the given roleID and actionName
func (r *CPSActionApproveIndexRepository) ExistsByRoleAndAction(ctx context.Context, roleID string, actionName string) (bool, error) {
	filter := bson.M{
		"role_id":     roleID,
		"action_name": strings.ToUpper(strings.TrimSpace(actionName)),
		"maker_index": bson.M{"$exists": true, "$ne": nil},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
