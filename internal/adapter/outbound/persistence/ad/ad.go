package ad

import (
	"context"
	"errors"
	"fmt"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/ad"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ADPersistence struct {
	adDal  dal.MongoDal[entity.Advert, entity.Advert]
	cpsDal dal.MongoDal[model.CPSAction, model.CPSAction]
	logger utils.Logger
}

var _ ad.ADRepo = (*ADPersistence)(nil)

func InitAD(client *mongo.Client, database string, collections []string, logger utils.Logger) ad.ADRepo {
	adDal := dal.NewMongoDal[entity.Advert, entity.Advert](client, database, collections[0])
	cpsDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collections[1])
	return &ADPersistence{
		adDal:  adDal,
		cpsDal: cpsDal,
		logger: logger,
	}
}

func (a *ADPersistence) CheckCPSActionExists(ctx context.Context, cpsAction model.CreateCPSAction) (bool, error) {
	filter := bson.M{
		"maker_phone_number": cpsAction.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsAction.Department,
		"request_action":     cpsAction.RequestAction,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAd, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		return false, fmt.Errorf("FAILED_TO_GET_AD")
	}
	if existingAd != nil {
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return true, fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	return false, nil
}

func (a *ADPersistence) CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {
	filter := bson.M{
		"maker_phone_number": cpsAction.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
	}

	projection := bson.M{
		"action_code": 1,
	}

	existingAction, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		a.logger.Errorf("failed to get cpsAction: %v", err)
		return &model.CPSAction{}, fmt.Errorf("FAILED_TO_GET_CPS_ACTION")
	} else if existingAction != nil {
		a.logger.Infof("pending cps action present for user_code: %s, department: %s", existingAction.MakerID)
		return &model.CPSAction{}, fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	filter = bson.M{
		"maker_phone_number": cpsAction.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsAction.Department,
		"request_action":     cpsAction.RequestAction,
	}

	projection = bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAd, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		return nil, fmt.Errorf("FAILED_TO_GET_AD")

	} else if existingAd != nil {
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return nil, fmt.Errorf("PENDING_CPS_ACTION_PRESENT")

	}

	cpsRes, err := a.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		ActionStatus:     string(model.ActionPending),
		RequestAction:    string(model.RequestCreateAdvert),
		ActionType:       string(model.ActionCreate),
		CurrentAction:    cpsAction.ActionData,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		// LastModifiedAt:   time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	return &cpsRes, nil
}

func (a *ADPersistence) Authorize(ctx context.Context, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error) {
	var advert entity.Advert
	var err error

	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"department":    cpsAction.Department,
		"action_status": entity.ActionPending,
	}

	update := bson.M{
		"checker_id":           cpsAction.CheckerUser.UserCode,
		"checker_phone_number": cpsAction.CheckerUser.PhoneNumber,
		"checker_name":         cpsAction.CheckerUser.FullName,
		"action_status":        model.ActionApproved,
		"checker_action_time":  time.Now(),
		"last_updated_at":      time.Now(),
	}

	cpsRes, err := a.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		// if errors.Is(err,mongo.ErrNoDocuments)
		a.logger.Errorf("failed to update cps action", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	var actionData entity.Advert
	mapData := make(map[string]interface{})

	if cpsRes.ActionType != string(model.ActionDelete) {

		data, err := bson.Marshal(cpsRes.CurrentAction)
		if err != nil {
			a.logger.Errorf("failed to marshal bson: %v", err)
			return nil, fmt.Errorf("INVALID_ACTION_DATA")
		}

		if err := bson.Unmarshal([]byte(data), &actionData); err != nil {
			a.logger.Errorf("failed to unmarshal into Avatar: %v", err)
			return nil, fmt.Errorf("INVALID_ACTION_DATA")
		}
		if err := bson.Unmarshal([]byte(data), &mapData); err != nil {
			a.logger.Errorf("failed to unmarshal into Avatar: %v", err)
			return nil, fmt.Errorf("INVALID_ACTION_DATA")
		}
	} else {
		mapData["_id"] = cpsRes.UniqueId
	}
	fmt.Printf("Type of _id:%T", mapData["_id"])
	objId, ok := mapData["_id"].(bson.ObjectID)
	if !ok {
		a.logger.Errorf("interface not type of objcetID")
		return nil, fmt.Errorf("INVALID_ID_TYPE")
	}

	if cpsRes.ActionType == string(model.ActionCreate) {
		req := entity.Advert{
			ID:            objId,
			Title:         actionData.Title,
			Description:   actionData.Description,
			BannerImage:   actionData.BannerImage,
			AdvertFor:     actionData.AdvertFor,
			Date:          actionData.Date,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		}

		advert, err = a.adDal.InsertOne(ctx, req)
		if err != nil {
			a.logger.Errorf("failed to create cps action", err)
			return nil, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
		}

		cpsRes.CurrentAction = advert

		return &cpsRes, nil

	}

	fmt.Println("check here-----------1-", actionData.ID)

	if cpsRes.ActionType == string(model.ActionUpdate) {
		filter := bson.M{"_id": objId}
		update := bson.M{}

		if actionData.Title != "" {
			update["title"] = actionData.Title
		}

		if actionData.Description != "" {
			update["description"] = actionData.Description
		}

		if !actionData.Date.StartedAt.IsZero() {
			update["date.started_at"] = actionData.Date.StartedAt
		}

		if !actionData.Date.ExpiredAt.IsZero() {
			update["date.expired_at"] = actionData.Date.ExpiredAt
		}

		update["last_updated_at"] = time.Now()

		fmt.Println("check here------------", filter)
		advert, err = a.adDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update advert", err)
			return nil, fmt.Errorf("FAILED_TO_UPDATE_ADVERT")
		}

		cpsRes.CurrentAction = advert

		return &cpsRes, nil
	}

	if cpsRes.ActionType == string(model.ActionDelete) {
		filter := bson.M{
			"_id": objId,
		}

		update := bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}

		advert, err = a.adDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update advert", err)

			return nil, fmt.Errorf("FAILED_TO_UPDATE_ADVERT")

		}
		cpsRes.CurrentAction = advert

		return &cpsRes, nil
	}

	return &cpsRes, nil
}

func (a *ADPersistence) UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {

	_, err := a.CheckCPSActionExists(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("failed to check existing cps action", err)
		return nil, err
	}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid id providedd", err)
		return nil, fmt.Errorf("INVALID_ID")
	}

	adFilter := bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}

	adProjection := bson.M{
		"title":        1,
		"description":  1,
		"banner_image": 1,
		"advert_for":   1,
		"date":         1,
	}

	ad, err := a.adDal.FindOne(ctx, adFilter, adProjection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("advert not found", err)
			return nil, fmt.Errorf("ADVERT_NOT_FOUND")
		}
		a.logger.Errorf("failed to get ad", err)

		return nil, fmt.Errorf("FAILED_TO_GET_AD")
	}

	cps, err := a.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		CurrentAction:    cpsAction.ActionData,
		RequestAction:    string(model.RequestUpdateAdvert),
		PreviousAction: map[string]any{
			"title":        ad.Title,
			"description":  ad.Description,
			"banner_image": ad.BannerImage,
			"advert_for":   ad.AdvertFor,
			"date":         ad.Date,
		},
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	return &cps, nil
}

func (a *ADPersistence) DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {
	filter := bson.M{
		"maker_phone_number": cpsAction.MakerUser.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         cpsAction.Department,
		"request_action":     cpsAction.RequestAction,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAd, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		return nil, fmt.Errorf("FAILED_TO_GET_AD")
	} else if existingAd != nil {
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return nil, fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid id provided", err)
		return nil, fmt.Errorf("INVALID_ID")
	}

	adFilter := bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}

	adProjection := bson.M{
		"title":        1,
		"descripttion": 1,
		"banner_image": 1,
		"advert_for":   1,
		"date":         1,
	}

	ad, err := a.adDal.FindOne(ctx, adFilter, adProjection)
	if err != nil {
		a.logger.Errorf("failed to get ad", err)
		return nil, fmt.Errorf("FAILED_TO_GET_AD")
	}

	cpsRes, err := a.cpsDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		ActionStatus:     string(model.ActionPending),
		RequestAction:    string(model.RequestDeleteAdvert),
		ActionType:       string(model.ActionDelete),
		CurrentAction:    cpsAction.ActionData,
		UniqueId:         id,
		PreviousAction: map[string]any{
			"title":        ad.Title,
			"description":  ad.Description,
			"banner_image": ad.BannerImage,
			"advert_for":   ad.AdvertFor,
			"date":         ad.Date,
			"is_deleted":   true,
			"deleted_at":   ad.DeletedAt,
		},
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	return &cpsRes, nil
}

func (a *ADPersistence) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid id provided", err)
		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}
	projection := bson.M{}

	advert, err := a.adDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("ad not found", err)
			return nil, fmt.Errorf("AD_NOT_FOUND")
		}
		a.logger.Errorf("failed to get ad", err)
		return nil, fmt.Errorf("AD_NOT_FOUND")
	}
	return advert, nil
}

func (a *ADPersistence) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Advert], error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["action_status"] = filterParams.Filters
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	ad, err := a.adDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("no ad data found", err)
			return nil, fmt.Errorf("NO_DATA_FOUND")
		}
		a.logger.Errorf("failed to get ad data", err)
		return nil, fmt.Errorf("AD_NOT_FOUND")

	}

	total, err := a.adDal.TotalCount(ctx, bson.M{})
	if err != nil {
		a.logger.Errorf("failed to get ad total counts", err)
		return nil, fmt.Errorf("AD_NOT_FOUND")

	}

	meta := common_util.BuildPaginationMeta(total, page, limit)

	return &common_util.PaginatedResponse[[]*entity.Advert]{
		Data: ad,
		Meta: meta,
	}, nil
}

func (a *ADPersistence) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {

	checkerUser := contexts.ExtractContext(ctx)

	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_id":           checkerUser.UserCode,
		"checker_name":         checkerUser.FullName,
		"checker_phone_number": checkerUser.PhoneNumber,
		"action_status":        model.ActionRejected,
		"rejection_reason":     req.RejectedReason,
		"checker_action_time":  time.Now(),
		"last_updated_at":      time.Now(),
	}

	cpsAction, err := a.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("cps action not found", err)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to update advert status", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_ADVERT")

	}
	return &cpsAction, nil
}
