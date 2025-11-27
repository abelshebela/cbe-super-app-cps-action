package account_block

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AccountBlockMapperForUpdate(block model.AccountBlock) bson.M {
	update := bson.M{}
	now := time.Now()
	update["updated_at"] = now

	if block.Name != "" {
		update["name"] = block.Name
	}
	if block.Code != "" {
		update["code"] = block.Code
	}
	if block.Address != "" {
		update["address"] = block.Address
	}
	if block.Slug != "" {
		update["slug"] = block.Slug
	}
	if block.Type != "" {
		update["type"] = block.Type
	}
	if block.CityID != "" {
		update["city_id"] = block.CityID
	}
	if block.DistrictID != "" {
		update["district_id"] = block.DistrictID
	}
	if block.RegionID != "" {
		update["region_id"] = block.RegionID
	}
	if block.ParentID != nil {
		update["parent_id"] = block.ParentID
	}

	update["is_enabled"] = block.IsEnabled

	update["is_deleted"] = block.IsDeleted

	return update
}

func FindAccountBlocksWithParentPopulatedRecursive(
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
	skip int64,
	limit int64,
	logger utils.Logger,
) ([]*model.AccountBlock, error) {
	pipeline := mongo.Pipeline{
		// 1. Match documents based on the initial filter
		{{Key: "$match", Value: filter}},

		// 2. Perform recursive lookup for ancestors (parents)
		{{Key: "$graphLookup", Value: bson.M{
			"from":                    "account_block",
			"startWith":               "$parent_id",
			"connectFromField":        "parent_id",
			"connectToField":          "_id",
			"as":                      "ancestors",
			"depthField":              "depth",
			"restrictSearchWithMatch": bson.M{"is_deleted": false},
		}}},

		// 3. Apply pagination before processing
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},

		// 4. Add a field to help with sorting ancestors by depth
		{{Key: "$addFields", Value: bson.M{
			"sortedAncestors": bson.M{
				"$sortArray": bson.M{
					"input":  "$ancestors",
					"sortBy": bson.M{"depth": -1}, // Sort by depth descending (deepest first)
				},
			},
		}}},

		// 5. Project final structure
		{{Key: "$project", Value: bson.M{
			"_id":             1,
			"name":            1,
			"code":            1,
			"address":         1,
			"parent_id":       1,
			"slug":            1,
			"type":            1,
			"is_enabled":      1,
			"city_id":         1,
			"district_id":     1,
			"region_id":       1,
			"is_deleted":      1,
			"created_at":      1,
			"updated_at":      1,
			"sortedAncestors": 1,
		}}},
	}

	// Execute aggregation
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Errorf("Error aggregating account blocks with recursive parent: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	// Decode results into intermediate structure
	type ResultWithAncestors struct {
		ID              bson.ObjectID          `bson:"_id,omitempty"`
		Name            string                 `bson:"name"`
		Code            string                 `bson:"code"`
		Address         string                 `bson:"address"`
		ParentID        *bson.ObjectID         `bson:"parent_id,omitempty"`
		Slug            string                 `bson:"slug"`
		Type            model.AccountBlockType `bson:"type"`
		IsEnabled       bool                   `bson:"is_enabled"`
		CityID          string                 `bson:"city_id,omitempty"`
		DistrictID      string                 `bson:"district_id,omitempty"`
		RegionID        string                 `bson:"region_id,omitempty"`
		IsDeleted       bool                   `bson:"is_deleted,omitempty"`
		CreatedAt       time.Time              `bson:"created_at"`
		UpdatedAt       time.Time              `bson:"updated_at"`
		SortedAncestors []model.AccountBlock   `bson:"sortedAncestors"`
	}

	var intermediateResults []ResultWithAncestors
	if err := cursor.All(ctx, &intermediateResults); err != nil {
		logger.Errorf("Error decoding account blocks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Build the nested parent structure
	results := make([]*model.AccountBlock, 0, len(intermediateResults))
	for _, item := range intermediateResults {
		accountBlock := &model.AccountBlock{
			ID:         item.ID,
			Name:       item.Name,
			Code:       item.Code,
			Address:    item.Address,
			ParentID:   item.ParentID,
			Slug:       item.Slug,
			Type:       item.Type,
			IsEnabled:  item.IsEnabled,
			CityID:     item.CityID,
			DistrictID: item.DistrictID,
			RegionID:   item.RegionID,
			IsDeleted:  item.IsDeleted,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		}

		// Build nested parent structure
		if len(item.SortedAncestors) > 0 {
			accountBlock.Parent = buildParentHierarchy(item.SortedAncestors)
		}

		results = append(results, accountBlock)
	}

	return results, nil
}

// buildParentHierarchy constructs the nested parent structure from a sorted list of ancestors
func buildParentHierarchy(ancestors []model.AccountBlock) *model.AccountBlock {
	if len(ancestors) == 0 {
		return nil
	}

	// Create a map for quick lookup by ID
	ancestorMap := make(map[bson.ObjectID]*model.AccountBlock)
	for i := range ancestors {
		ancestor := &model.AccountBlock{
			ID:         ancestors[i].ID,
			Name:       ancestors[i].Name,
			Code:       ancestors[i].Code,
			Address:    ancestors[i].Address,
			ParentID:   ancestors[i].ParentID,
			Slug:       ancestors[i].Slug,
			Type:       ancestors[i].Type,
			IsEnabled:  ancestors[i].IsEnabled,
			CityID:     ancestors[i].CityID,
			DistrictID: ancestors[i].DistrictID,
			RegionID:   ancestors[i].RegionID,
			IsDeleted:  ancestors[i].IsDeleted,
			CreatedAt:  ancestors[i].CreatedAt,
			UpdatedAt:  ancestors[i].UpdatedAt,
		}
		ancestorMap[ancestor.ID] = ancestor
	}

	// Link each ancestor to its parent
	for _, ancestor := range ancestorMap {
		if ancestor.ParentID != nil {
			if parent, exists := ancestorMap[*ancestor.ParentID]; exists {
				ancestor.Parent = parent
			}
		}
	}

	// Find the immediate parent (the one at depth 0)
	for i := range ancestors {
		if ancestors[i].ParentID != nil {
			if immediateParent, exists := ancestorMap[ancestors[i].ID]; exists {
				// Check if this is depth 0 (immediate parent)
				// Since we sorted by depth descending, the last one should be depth 0
				if i == len(ancestors)-1 {
					return immediateParent
				}
			}
		}
	}

	// Fallback: return the first ancestor
	if len(ancestorMap) > 0 {
		for _, v := range ancestorMap {
			return v
		}
	}

	return nil
}
