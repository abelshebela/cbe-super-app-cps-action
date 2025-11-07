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

	if block.ParentID != nil {
		update["parent_id"] = block.ParentID
	}

	update["is_enabled"] = block.IsEnabled

	update["is_deleted"] = block.IsDeleted

	return update
}

// FindAccountBlocksWithParentPopulated executes an aggregation pipeline to find account blocks
// with their parent documents populated. This function can be used for any account block type
// (Region, District, City, Branch) by passing the appropriate filter.
func FindAccountBlocksWithParentPopulated(
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
	skip int64,
	limit int64,
	logger utils.Logger,
) ([]*model.AccountBlock, error) {
	pipeline := mongo.Pipeline{
		// Match documents based on filter
		{{Key: "$match", Value: filter}},
		// Lookup parent document
		{{Key: "$lookup", Value: bson.M{
			"from":         "account_block",
			"localField":   "parent_id",
			"foreignField": "_id",
			"as":           "parent_docs",
		}}},
		// Unwind parent array (preserve null for documents without parent)
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$parent_docs",
			"preserveNullAndEmptyArrays": true,
		}}},
		// Project fields including parent
		{{Key: "$project", Value: bson.M{
			"_id":       1,
			"name":      1,
			"code":      1,
			"address":   1,
			"parent_id": 1,
			"parent": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$ne": []interface{}{"$parent_docs", nil}},
					"then": "$parent_docs",
					"else": nil,
				},
			},
			"slug":       1,
			"type":       1,
			"is_enabled": 1,
			"is_deleted": 1,
			"created_at": 1,
			"updated_at": 1,
		}}},
		// Apply pagination
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	// Execute aggregation
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Errorf("Error aggregating account blocks with parent: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	// Decode results
	var results []*model.AccountBlock
	if err := cursor.All(ctx, &results); err != nil {
		logger.Errorf("Error decoding account blocks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return results, nil
}

// FindAccountBlockByCodeWithParentPopulated executes an aggregation pipeline to find a single account block
// by code with its parent document populated. This function can be used for any account block type
// (Region, District, City, Branch) by passing the appropriate code and optional type filter.
func FindAccountBlockByCodeWithParentPopulated(
	ctx context.Context,
	collection *mongo.Collection,
	code string,
	accountBlockType string, // Optional: "R", "D", "C", "B" or empty string for any type
	logger utils.Logger,
) (*model.AccountBlock, error) {
	// Build filter
	filter := bson.M{"code": code}
	if accountBlockType != "" {
		filter["type"] = accountBlockType
	}

	// Build aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match document by code (and optionally by type)
		{{Key: "$match", Value: filter}},
		// Lookup parent document
		{{Key: "$lookup", Value: bson.M{
			"from":         "account_block",
			"localField":   "parent_id",
			"foreignField": "_id",
			"as":           "parent_docs",
		}}},
		// Unwind parent array (preserve null for documents without parent)
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$parent_docs",
			"preserveNullAndEmptyArrays": true,
		}}},
		// Project fields including parent
		{{Key: "$project", Value: bson.M{
			"_id":       1,
			"name":      1,
			"code":      1,
			"address":   1,
			"parent_id": 1,
			"parent": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$ne": []interface{}{"$parent_docs", nil}},
					"then": "$parent_docs",
					"else": nil,
				},
			},
			"slug":       1,
			"type":       1,
			"is_enabled": 1,
			"is_deleted": 1,
			"created_at": 1,
			"updated_at": 1,
		}}},
		// Limit to 1 result
		{{Key: "$limit", Value: 1}},
	}

	// Execute aggregation
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Errorf("Error aggregating account block by code with parent: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	// Check if document exists
	if !cursor.Next(ctx) {
		return nil, mongo.ErrNoDocuments
	}

	// Decode result
	var result model.AccountBlock
	if err := cursor.Decode(&result); err != nil {
		logger.Errorf("Error decoding account block: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, nil
}

// func BranchMapperForUpdate(branch model.Branch) bson.M {
// 	update := bson.M{}
// 	now := time.Now()
// 	update["updated_at"] = now

// 	if branch.BranchCode != "" {
// 		update["branch_code"] = branch.BranchCode
// 	}
// 	if branch.BranchName != "" {
// 		update["branch_name"] = branch.BranchName
// 	}
// 	if branch.BranchAddress != "" {
// 		update["branch_address"] = branch.BranchAddress
// 	}
// 	if branch.DistrictCode != "" {
// 		update["district_code"] = branch.DistrictCode
// 	}
// 	if branch.DistrictName != "" {
// 		update["district_name"] = branch.DistrictName
// 	}
// 	if branch.RegionName != "" {
// 		update["region_name"] = branch.RegionName
// 	}
// 	if branch.RecordStat != "" {
// 		update["record_stat"] = branch.RecordStat
// 	}
// 	// Booleans are tricky: include them only if they are explicitly meant to be updated
// 	update["enabled"] = branch.Enabled

// 	return update
// }

// func RegionMapperForUpdate(region model.Region) bson.M {
// 	update := bson.M{}
// 	now := time.Now()
// 	update["updated_at"] = now

// 	if region.RegionCode != "" {
// 		update["region_code"] = region.RegionCode
// 	}
// 	if region.RegionName != "" {
// 		update["region_name"] = region.RegionName
// 	}
// 	if region.RegionAddress != "" {
// 		update["region_address"] = region.RegionAddress
// 	}
// 	update["enabled"] = region.Enabled

// 	return update
// }

// func DistrictMapperForUpdate(district model.District) bson.M {
// 	update := bson.M{}
// 	now := time.Now()
// 	update["updated_at"] = now

// 	if district.DistrictCode != "" {
// 		update["district_code"] = district.DistrictCode
// 	}
// 	if district.DistrictName != "" {
// 		update["district_name"] = district.DistrictName
// 	}
// 	if district.DistrictAddress != "" {
// 		update["district_address"] = district.DistrictAddress
// 	}
// 	if district.RegionID != "" {
// 		update["region_id"] = district.RegionID
// 	}
// 	if district.RegionName != "" {
// 		update["region_name"] = district.RegionName
// 	}
// 	update["enabled"] = district.Enabled

// 	return update
// }

// func CityMapperForUpdate(city model.City) bson.M {
// 	update := bson.M{}
// 	now := time.Now()
// 	update["updated_at"] = now

// 	if city.CityCode != "" {
// 		update["city_code"] = city.CityCode
// 	}
// 	if city.CityName != "" {
// 		update["city_name"] = city.CityName
// 	}
// 	if city.CityAddress != "" {
// 		update["city_address"] = city.CityAddress
// 	}
// 	if city.DistrictID != "" {
// 		update["district_id"] = city.DistrictID
// 	}
// 	if city.DistrictName != "" {
// 		update["district_name"] = city.DistrictName
// 	}
// 	if city.RegionID != "" {
// 		update["region_id"] = city.RegionID
// 	}
// 	if city.RegionName != "" {
// 		update["region_name"] = city.RegionName
// 	}
// 	update["enabled"] = city.Enabled

// 	return update
// }
