package action_role_repo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ExistsByRoleAndAction returns true if a record exists for the given roleID and actionName
// and the corresponding cps_action_role is enabled.
func (r *CPSActionApproveIndexRepository) ExistsByRoleAndAction(ctx context.Context, roleID string, actionName string) (bool, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"role_id":     roleID,
			"action_name": strings.ToUpper(strings.TrimSpace(actionName)),
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cps_action_roles",
			"localField":   "action_name",
			"foreignField": "action_name",
			"as":           "action_role_info",
		}}},
		{{Key: "$match", Value: bson.M{
			"action_role_info": bson.M{
				"$elemMatch": bson.M{"enabled": true},
			},
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	return cursor.Next(ctx), nil
}

// HasActiveActionRoles returns true if the given role has any allocations
// in cps_action_approver_index where the corresponding cps_action_role is enabled.
func (r *CPSActionApproveIndexRepository) HasActiveActionRoles(ctx context.Context, roleCode string) (bool, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"role_id": roleCode,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cps_action_roles",
			"localField":   "action_name",
			"foreignField": "action_name",
			"as":           "action_role_info",
		}}},
		{{Key: "$match", Value: bson.M{
			"action_role_info": bson.M{
				"$elemMatch": bson.M{"enabled": true},
			},
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	return cursor.Next(ctx), nil
}
