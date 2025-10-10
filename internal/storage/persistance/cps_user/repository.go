package cps_user

import (
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"


	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSUserStorage struct {
	dal        dal.MongoDal[model.CPSUser, model.CPSUser]
	cpsAction  dal.MongoDal[model.CPSAction, model.CPSAction]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewCPSUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.CpsUserRepository {
	return &CPSUserStorage{
		dal:        dal.NewMongoDal[model.CPSUser, model.CPSUser](client, dbName, collection),
		cpsAction:  dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "cps_actions"),
		client:     client,
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

// Implement actual repository methods for CPS action authorization
func (r *CPSUserStorage) Create(ctx context.Context, cpsUser *model.CPSUser) error {
	_, err := r.dal.InsertOne(ctx, *cpsUser)
	if err != nil {
		r.logger.Errorf("failed to create CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) Update(ctx context.Context, userCode string, cpsUser *model.CPSUser) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := CPSUserUpdateMapper(cpsUser)

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to update CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) Delete(ctx context.Context, userCode string) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"is_deleted": true}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to delete CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) EnableOrDisable(ctx context.Context, userCode string, enable bool) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified": time.Now()}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to enable/disable CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByID(ctx context.Context, id string) (*model.CPSUser, error) {
	var filter bson.M
	if objID, ok := local_util.StringToObjectID(id); ok {
		filter = bson.M{"_id": objID, "is_deleted": false}
	} else {
		filter = bson.M{"user_code": id, "is_deleted": false}
	}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	return result, nil
}

func (r *CPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CPSUser], error) {
	filter := bson.M{}
	searchKeys := bson.M{}

	allowedKeys := []string{"enabled", "department", "role"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := r.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*model.CPSUser]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *CPSUserStorage) GetPopulatedByID(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
    const (
        departmentColl         = "department"
        portalCardsColl        = "cards"
        permissionGroupsColl   = "permission_groups"
        permissionCategoryColl = "permission_category"
        permissionColl         = "permission"
    )

    // Universal ID converter function
    convertIDs := func(fieldName string) bson.M {
        return bson.M{
            "$map": bson.M{
                "input": fieldName,
                "as":    "id",
                "in": bson.M{
                    "$cond": bson.M{
                        "if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$$id"}, "string"}},
                        "then": bson.M{"$toObjectId": "$$id"},
                        "else": "$$id",
                    },
                },
            },
        }
    }

    // Lookup pipeline for permission categories
    categoryLookup := bson.D{{Key: "$lookup", Value: bson.M{
        "from": permissionCategoryColl,
        "let":  bson.M{"categoryIds": "$permission_category"},
        "pipeline": mongo.Pipeline{
            bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
                "$in": []interface{}{"$_id", convertIDs("$$categoryIds")},
            }}}},
            bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
            bson.D{{Key: "$project", Value: bson.M{
                "category_name": 1,
                "access":        1,
                "permissions":   1,
            }}},
            bson.D{{Key: "$lookup", Value: bson.M{
                "from": permissionColl,
                "let":  bson.M{"permissionIds": "$permissions"},
                "pipeline": mongo.Pipeline{
                    bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
                        "$in": []interface{}{"$_id", convertIDs("$$permissionIds")},
                    }}}},
                    bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
                    bson.D{{Key: "$project", Value: bson.M{"permission_name": 1}}},
                },
                "as": "permissions_docs",
            }}},
        },
        "as": "permission_category_docs",
    }}}

    pipeline := mongo.Pipeline{
        bson.D{{Key: "$match", Value: bson.M{"user_code": userCode, "is_deleted": false}}},
        
        // Lookup department
        bson.D{{Key: "$lookup", Value: bson.M{
            "from":         departmentColl,
            "localField":   "department",
            "foreignField": "_id",
            "as":           "department_doc",
        }}},
        bson.D{{Key: "$unwind", Value: bson.M{"path": "$department_doc", "preserveNullAndEmptyArrays": true}}},
        
        // Lookup portal cards
        bson.D{{Key: "$lookup", Value: bson.M{
            "from": portalCardsColl,
            "let":  bson.M{"portalCards": "$department_doc.portal_cards"},
            "pipeline": mongo.Pipeline{
                bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
                    "$in": []interface{}{"$_id", convertIDs("$$portalCards")},
                }}}},
                bson.D{{Key: "$project", Value: bson.M{"_id": 1, "card_name": 1}}},
            },
            "as": "portal_cards_docs",
        }}},
        
        // Lookup permission groups
        bson.D{{Key: "$lookup", Value: bson.M{
            "from":         permissionGroupsColl,
            "localField":   "permission_group",
            "foreignField": "_id",
            "as":           "permission_groups_raw",
            "pipeline": mongo.Pipeline{
                bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
                bson.D{{Key: "$project", Value: bson.M{
                    "group_name":         1,
                    "permission_category": 1,
                }}},
                categoryLookup,
            },
        }}},
        
        // Project final response
        bson.D{{Key: "$project", Value: bson.M{
            "_id":            1,
            "user_code":      1,
            "full_name":      1,
            "role":           1,
            "department":     1,
            "gender":         1,
            "phone_number":   1,
            "email":          1,
            "username":       1,
            "realm":          1,
            "enabled":        1,
            "date_joined":    1,
            "last_modified":  1,
            "country":        1,
            "region":         1,
            "department_name": "$department_doc.department",
            "portal_cards": bson.M{
                "$map": bson.M{
                    "input": "$portal_cards_docs",
                    "as":    "card",
                    "in":    "$$card.card_name",
                },
            },
            "permission_groups": bson.M{
                "$map": bson.M{
                    "input": "$permission_groups_raw",
                    "as":    "pg",
                    "in": bson.M{
                        "id":         "$$pg._id",
                        "group_name": "$$pg.group_name",
                        "permission_category": bson.M{
                            "$map": bson.M{
                                "input": "$$pg.permission_category_docs",
                                "as":    "cat",
                                "in": bson.M{
                                    "id":            "$$cat._id",
                                    "category_name": "$$cat.category_name",
                                    "access":        "$$cat.access",
                                    "permissions": bson.M{
                                        "$map": bson.M{
                                            "input": "$$cat.permissions_docs",
                                            "as":    "perm",
                                            "in": bson.M{
                                                "id":               "$$perm._id",
                                                "permission_name":  "$$perm.permission_name",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
        }}},
    }

    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        r.logger.Errorf("failed to aggregate cps user by id: %v", err)
        return nil, errors.New(localization.ErrorUnexpectedError.Message)
    }
    
    defer cursor.Close(ctx)

    if !cursor.Next(ctx) {
        return nil, errors.New(localization.ErrorFileNotFound.Code)
    }

    var resp cpsuser.CpsUserResponse
    if err := cursor.Decode(&resp); err != nil {
        r.logger.Errorf("failed to decode cps user response: %v", err)
        return nil, errors.New(localization.ErrorUnexpectedError.Message)
    }

    return &resp, nil
}