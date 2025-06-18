package ad

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	dal "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/ad"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type ADPersistence struct {
	adDal  dal.MongoDal[entity.Advert, entity.Advert]
	cpsDal dal.MongoDal[entity.CPSAction, entity.CPSAction]
	logger utils.Logger
}

var _ ad.ADRepo = (*ADPersistence)(nil)

func InitAD(client *mongo.Client, database string, collections []string, logger utils.Logger) *ADPersistence {
	adDal := dal.NewMongoDal[entity.Advert, entity.Advert](client, database, collections[0])
	cpsDal := dal.NewMongoDal[entity.CPSAction, entity.CPSAction](client, database, collections[1])
	return &ADPersistence{
		adDal:  adDal,
		cpsDal: cpsDal,
		logger: logger,
	}
}

func (a *ADPersistence) CreateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsAction.MakerUser.PhoneNumber,
		"status":                  entity.ActionPending,
		"department":              cpsAction.Department,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAd, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return nil, err
	}

	cpsAction.Status = entity.ActionPending
	cpsAction.RequestAction = entity.RequestCreateAdvert
	cpsAction.ActionType = entity.ActionCreate
	cpsAction.MakerActionTime = time.Now()

	cps, err := a.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (a *ADPersistence) UpdateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	filter := bson.M{
		"maker_user.phone_number": cpsAction.MakerUser.PhoneNumber,
		"status":                  entity.ActionPending,
		"department":              cpsAction.Department,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAd, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return nil, err
	}

	adFilter := bson.M{
		"id":         cpsAction.ActionData.ID,
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
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cpsAction.Status = entity.ActionPending
	cpsAction.RequestAction = entity.RequestUpdateAdvert
	cpsAction.ActionType = entity.ActionUpdate
	cpsAction.MakerActionTime = time.Now()

	cpsAction.PreviousData = map[string]any{
		"title":       ad.Title,
		"description": ad.Description,
		"advert_for":  ad.AdvertFor,
		"date":        ad.Date,
	}

	cpsAction.CurrentData = map[string]any{
		"title":       cpsAction.ActionData.Title,
		"description": cpsAction.ActionData.Description,
		"advert_for":  cpsAction.ActionData.AdvertFor,
		"date":        cpsAction.ActionData.Date,
	}

	cps, err := a.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (a *ADPersistence) DeleteOneAdvert(ctx context.Context, cpsAction entity.CPSAction) error {
	filter := bson.M{
		"maker_user.phone_number": cpsAction.MakerUser.PhoneNumber,
		"status":                  entity.ActionPending,
		"department":              cpsAction.Department,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAd, err := a.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	} else if existingAd != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		a.logger.Infof("pending cps action present", cpsAction.MakerUser.FullName,
			cpsAction.MakerUser.UserCode, cpsAction.Department)
		return err
	}

	adFilter := bson.M{
		"id":         cpsAction.ActionData.ID,
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
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	}

	cpsAction.PreviousData = map[string]any{
		"title":        ad.Title,
		"description":  ad.Description,
		"banner_image": ad.BannerImage,
		"advert_for":   ad.AdvertFor,
		"date":         ad.Date,
		"is_deleted":   ad.IsDeleted,
		"deleted_at":   ad.DeletedAt,
	}

	cpsAction.CurrentData = map[string]any{
		"is_deleted": true,
	}

	cpsAction.Status = entity.ActionPending
	cpsAction.RequestAction = entity.RequestDeleteAdvert
	cpsAction.ActionType = entity.ActionDelete
	cpsAction.MakerActionTime = time.Now()

	_, err = a.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	}

	return nil
}

func (a *ADPersistence) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}
	projection := bson.M{}

	advert, err := a.adDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("ad not found", err)
			err = fmt.Errorf("ad not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "ad not found",
			})
			return nil, err
		}
		a.logger.Errorf("failed to get ad", err)
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return advert, nil
}

func (a *ADPersistence) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	ad, err := a.adDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("no ad data found", err)
			err = fmt.Errorf("ad not found %w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "ad data not found",
			})
			return nil, err
		}
		a.logger.Errorf("failed to get ad data", err)
		err = fmt.Errorf("failed to get ad %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	total, err := a.adDal.TotalCount(ctx, bson.M{})
	if err != nil {
		a.logger.Errorf("failed to get ad total counts", err)
		err := fmt.Errorf("failed to get ad total counts %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &entity.AdvertResponse{
		Page:   filterParams.Page,
		Advert: ad,
		Limit:  constant.DefaultPerPage,
		Total:  total,
	}, nil
}

func (a *ADPersistence) Authorize(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	var advert entity.Advert
	var err error

	filter := bson.M{
		"action_code": cpsAction.ActionCode,
		"department":  cpsAction.Department,
		"status":      entity.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    cpsAction.CheckerUser.FullName,
			"phone_number": cpsAction.CheckerUser.PhoneNumber,
			"user_code":    cpsAction.CheckerUser.UserCode,
		},
		"status":              entity.ActionApproved,
		"checker_action_time": time.Now(),
	}
	cpsAction, err = a.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("failed to update cps action", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	if cpsAction.ActionType == entity.ActionCreate {
		req := entity.Advert{
			ID:          bson.NewObjectID().Hex(),
			Title:       cpsAction.ActionData.Title,
			Description: cpsAction.ActionData.Description,
			BannerImage: cpsAction.ActionData.BannerImage,
			AdvertFor:   cpsAction.ActionData.AdvertFor,
			Date:        cpsAction.ActionData.Date,
			CreatedAt:   time.Now(),
		}
		advert, err = a.adDal.InsertOne(ctx, req)
		if err != nil {
			a.logger.Errorf("failed to create cps action", err)
			err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}

		cpsAction.ActionData = advert

		return &cpsAction, nil

	}

	if cpsAction.ActionType == entity.ActionUpdate {
		filter := bson.M{"id": cpsAction.ActionData.ID}
		update := bson.M{}

		if cpsAction.ActionData.Title != "" {
			update["title"] = cpsAction.ActionData.Title
		}

		if cpsAction.ActionData.Description != "" {
			update["description"] = cpsAction.ActionData.Description
		}

		if !cpsAction.ActionData.Date.StartedAt.IsZero() {
			update["date.started_at"] = cpsAction.ActionData.Date.StartedAt
		}

		if !cpsAction.ActionData.Date.ExpiredAt.IsZero() {
			update["date.expired_at"] = cpsAction.ActionData.Date.ExpiredAt
		}

		update["last_updated_at"] = time.Now()

		advert, err = a.adDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update advert", err)
			err = fmt.Errorf("failed to update advert %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}

		cpsAction.ActionData = advert

		return &cpsAction, nil
	}

	if cpsAction.ActionType == entity.ActionDelete {
		filter := bson.M{
			"id": cpsAction.ActionData.ID,
		}

		update := bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}

		advert, err = a.adDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update advert", err)
			err = fmt.Errorf("failed to update advert %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return nil, err
		}
		cpsAction.ActionData = advert

		return &cpsAction, nil
	}

	cpsAction.ActionData = entity.Advert{
		Title: advert.Title,
	}

	return &cpsAction, nil
}

func (a *ADPersistence) Reject(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	filter := bson.M{
		"action_code": cpsAction.ActionCode,
		"department":  cpsAction.Department,
		"status":      entity.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    cpsAction.CheckerUser.FullName,
			"phone_number": cpsAction.CheckerUser.PhoneNumber,
			"user_code":    cpsAction.CheckerUser.UserCode,
		},
		"status":              entity.ActionRejected,
		"rejected_reason":     cpsAction.RejectedReason,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := a.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("failed to update advert status", err)
		err = fmt.Errorf("failed to update advert status %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return &cpsAction, nil
}
