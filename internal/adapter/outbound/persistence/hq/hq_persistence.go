package hq

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	models "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HQPersistence struct {
	hqDal   dal.MongoDal[hq.HQ, hq.HQ]
	timeout time.Duration
	logger  utils.Logger
}

func NewHQPersistence(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *HQPersistence {
	hqDal := dal.NewMongoDal[hq.HQ, hq.HQ](client, dbName, "hq")
	return &HQPersistence{
		hqDal:   hqDal,
		timeout: timeout,
		logger:  logger,
	}
}

func (p *HQPersistence) GetHQByID(ctx context.Context, id string) (hq.HQ, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if id == "" {
		p.logger.Errorf("invalid HQ ID: empty")
		return hq.HQ{}, mongo.ErrNoDocuments
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("invalid HQ ObjectID: %v", err)
		return hq.HQ{}, err
	}

	filter := bson.M{"_id": objID}

	result, err := p.hqDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			p.logger.Errorf("HQ not found: id=%s", id)
			return hq.HQ{}, err
		}
		p.logger.Errorf("failed to fetch HQ: %v", err)
		return hq.HQ{}, err
	}
	if result == nil {
		p.logger.Errorf("HQ not found: id=%s", id)
		return hq.HQ{}, mongo.ErrNoDocuments
	}
	return *result, nil
}

func modelToDomainHQ(m models.HQ) hq.HQ {
	return hq.HQ{
		ID:           m.ID,
		Name:         m.Name,
		BlockTime:    m.BlockTime,
		ArchiveTime:  m.ArchiveTime,
		CreatedAt:    m.CreatedAt,
		LastModified: m.LastModifiedAt,
	}
}

func normalizePhone(search string) bson.M {
	re := regexp.MustCompile(`^\+?251[79]\d{8}$`)

	trimmedSearch := strings.TrimSpace(search)

	if re.MatchString(trimmedSearch) {
		return bson.M{
			"phone_number": bson.M{
				"$regex":   trimmedSearch,
				"$options": "i",
			},
		}
	}

	return bson.M{
		"phone_number.number": bson.M{
			"$regex":   "^$",
			"$options": "i",
		},
	}
}
func (h *HQPersistence) GetAllHQ(ctx context.Context, filterParams *constant.Filter) (*hq.HQRespose, error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{}

	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"address": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"email": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"enabled": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"created_at": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"last_modified": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"latest_android_version": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"latest_ios_version": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"linked_accounts": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"phone_number": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"archive_time": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"block_time": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	if filterParams.Filters != "" {
		filter["account_status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	hqData, err := h.hqDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.logger.Errorf("no HQ data found", err)

			return nil, fmt.Errorf("NO_HQ_DATA_FOUND")
		}
		h.logger.Errorf("failed to get hq data", err)
		return nil, fmt.Errorf("FAILED_TO_GET_HQ")
	}

	total, err := h.hqDal.TotalCount(ctx, bson.M{})
	if err != nil {
		h.logger.Errorf("failed to get total counts", err)

		return nil, fmt.Errorf("FAILED_TO_GET_HQ_COUNT")
	}

	return &hq.HQRespose{
		Page:  1,
		HQ:    hqData,
		Limit: constant.DefaultPerPage,
		Total: total,
	}, nil
}
func (p *HQPersistence) UpdateHQ(ctx context.Context, id string, update hq.HQ) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("invalid HQ ObjectID: %v", err)
		return err
	}
	// fmt.Println("update", update)
	updateDoc := bson.M{
		"block_time":       update.BlockTime,
		"archive_time":     update.ArchiveTime,
		"last_modified_at": time.Now(),
	}
	_, err = p.hqDal.UpdateOne(ctx, bson.M{"_id": objID}, updateDoc)
	if err != nil {
		p.logger.Errorf("failed to update HQ: %v", err)
		return err
	}
	return nil
}
