package miniappmerchant

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	miniApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/mini_app_merchant"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type miniAppMerchantPersistence struct {
	MongoDalMiniApp dal.MongoDal[model.MiniAppMerchant, model.MiniAppMerchant]
	logger          utils.Logger
	client          *mongo.Client
}

func NewMiniAppMerchantPersistence(client *mongo.Client, DB_name string, collection string, logger utils.Logger) miniApp.MiniAppMerchantRepository {
	return &miniAppMerchantPersistence{
		MongoDalMiniApp: dal.NewMongoDal[model.MiniAppMerchant, model.MiniAppMerchant](client, DB_name, collection),
		logger:          logger,
		client:          client,
	}
}

func (p *miniAppMerchantPersistence) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	p.logger.Debugf("Starting MongoDB session for transaction")

	session, err := p.client.StartSession()
	if err != nil {
		p.logger.Errorf("failed to start MongoDB session: %v", err)
		return fmt.Errorf(common_util.UnhandledServerError)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(txCtx context.Context) error {
		p.logger.Debugf("Starting MongoDB transaction")

		if err := session.StartTransaction(); err != nil {
			p.logger.Errorf("failed to start transaction: %v", err)
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		err := fn(txCtx)
		if err != nil {
			p.logger.Errorf("transaction logic failed: %v", err)
			if abortErr := session.AbortTransaction(txCtx); abortErr != nil {
				p.logger.Errorf("failed to abort transaction: %v", abortErr)
			} else {
				p.logger.Debugf("Transaction aborted successfully")
			}
			return err
		}

		if err := session.CommitTransaction(txCtx); err != nil {
			p.logger.Errorf("failed to commit transaction: %v", err)
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		p.logger.Debugf("Transaction committed successfully")
		return nil
	})
}

func (p *miniAppMerchantPersistence) CreateMiniAppMerchant(ctx context.Context, merchant *entities.MiniAppMerchant) (*entities.MiniAppMerchant, error) {

	merchant.CreatedAt = time.Now()
	merchant.LastModifiedAt = time.Now()
	doc, err := mappers.ToMiniAppMerchantModel(merchant)
	if err != nil {
		return nil, err
	}

	result, err := p.MongoDalMiniApp.InsertOne(ctx, *doc)
	if err != nil {
		p.logger.Errorf("failed to create miniapp merchant: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBInsertFailed)
	}
	return mappers.ToMiniAppMerchantDomain(&result), nil
}

func (p *miniAppMerchantPersistence) ListMiniAppMerchant(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*entities.MiniAppMerchant], error) {
	fmt.Println("I have been called ListMiniAppMerchant")

	filter := bson.M{"is_deleted": false}

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}

		filter["$or"] = []bson.M{
			{"merchant_name": searchRegex},
			{"merchant_representative_name": searchRegex},
			{"phone_number": searchRegex},
			{"email": searchRegex},
			{"account_number": searchRegex},
			{"type": searchRegex},
			{"mini_app_id": searchRegex},
		}
	}

	if filterParams.Filters != nil {
		allowedKeys := []string{"enabled", "merchant_type"}
		handlers := map[string]func(interface{}) interface{}{
			"is_deleted": func(value interface{}) interface{} {
				if str, ok := value.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						return parsed
					}
				}
				return value
			},
		}

		enhancedFilter := common_util.BuildMongoFilterWithHandlers(filterParams.Filters, allowedKeys, handlers)
		for key, value := range enhancedFilter {
			filter[key] = value
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	merchantDocs, err := p.MongoDalMiniApp.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		p.logger.Errorf("failed to paginate miniapp merchants: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	var merchants []*entities.MiniAppMerchant
	for _, doc := range merchantDocs {
		d := mappers.ToMiniAppMerchantDomain(doc)
		merchants = append(merchants, d)
	}

	total, err := p.MongoDalMiniApp.TotalCount(ctx, filter)
	if err != nil {
		p.logger.Errorf("failed to count miniapp merchants: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	meta := common_util.BuildPaginationMeta(total, filterParams.Page, limit)
	return &common_util.PaginatedResponse[[]*entities.MiniAppMerchant]{
		Data: merchants,
		Meta: meta,
	}, nil
}

func (p *miniAppMerchantPersistence) DetailMiniAppByID(ctx context.Context, id string) (*entities.MiniAppMerchant, error) {
	objID, err := mappers.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	doc, err := p.MongoDalMiniApp.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		p.logger.Errorf("failed to get miniapp merchant by ID: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}
	return mappers.ToMiniAppMerchantDomain(doc), nil
}

func (p *miniAppMerchantPersistence) UpdateMiniAppMerchant(ctx context.Context, merchant *entities.MiniAppMerchant) (*entities.MiniAppMerchant, error) {
	_, err := p.DetailMiniAppByID(ctx, merchant.ID)
	if err != nil {
		return nil, err
	}

	set := buildUpdateSet(merchant)

	if len(set) == 0 {
		return nil, fmt.Errorf(common_util.NoDataProvidedForUpdate)
	}

	objID, err := mappers.ObjectIDFromHex(merchant.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := p.MongoDalMiniApp.UpdateOne(ctx, filter, set)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		p.logger.Warnf("failed to update merchant fields: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	return mappers.ToMiniAppMerchantDomain(&result), nil
}

func (p *miniAppMerchantPersistence) EnableMiniAppMerchant(ctx context.Context, id string) (*entities.MiniAppMerchant, error) {
	return p.updateMerchantState(ctx, id, true)
}

func (p *miniAppMerchantPersistence) DisableMiniAppMerchant(ctx context.Context, id string) (*entities.MiniAppMerchant, error) {
	return p.updateMerchantState(ctx, id, false)
}

func (p *miniAppMerchantPersistence) DeleteMiniAppMerchant(ctx context.Context, id string) (*entities.MiniAppMerchant, error) {
	objID, err := mappers.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	set := bson.M{"is_deleted": true, "enabled": false}

	result, err := p.MongoDalMiniApp.UpdateOne(ctx, filter, set)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		p.logger.Warnf("failed to soft delete merchant: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	return mappers.ToMiniAppMerchantDomain(&result), nil
}

func (p *miniAppMerchantPersistence) MiniAppMerchantInfoExists(
	ctx context.Context,
	data entities.CheckMiniAppMerchant,
	opts *entities.MiniAppMerchantExistOptions) (bool, error) {
	filter := bson.M{
		"is_deleted": false,
		"$or":        []bson.M{},
	}

	if data.BankAccountNumber != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"bank_account_number": data.BankAccountNumber})
	}
	if data.Email != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"email": data.Email})
	}
	if data.PhoneNumber != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"phone_number": data.PhoneNumber})
	}

	if len(filter["$or"].([]bson.M)) == 0 {
		return false, nil
	}

	if opts != nil && opts.ExcludeID != "" {
		id, err := bson.ObjectIDFromHex(opts.ExcludeID)
		if err != nil {
			p.logger.Errorf("invalid exclude ID: %v", err)
			return false, fmt.Errorf("INVALID_ID")
		}
		filter["_id"] = bson.M{"$ne": id}
	}

	count, err := p.MongoDalMiniApp.TotalCount(ctx, filter)
	if err != nil {
		p.logger.Errorf("failed to check merchant info exists: %v", err)
		return false, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	return count > 0, nil
}

func (p *miniAppMerchantPersistence) updateMerchantState(ctx context.Context, id string, enabled bool) (*entities.MiniAppMerchant, error) {
	objID, err := mappers.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	set := bson.M{"enabled": enabled}

	result, err := p.MongoDalMiniApp.UpdateOne(ctx, filter, set)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		p.logger.Warnf("failed to update enable/disable: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	return mappers.ToMiniAppMerchantDomain(&result), nil
}

func (p *miniAppMerchantPersistence) AddMiniApp(ctx context.Context, merchantID string, miniApp entities.MiniApps) error {
	p.logger.Infof("AddMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)

	objID, err := mappers.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Warnf("AddMiniApp: invalid merchant_id=%s, error=%v", merchantID, err)
		return err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	_, err = p.MongoDalMiniApp.PushToArray(ctx, filter, "mini_apps", miniApp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("AddMiniApp: merchant not found, merchant_id=%s", merchantID)
			return fmt.Errorf(common_util.NotFound)
		}
		p.logger.Warnf("AddMiniApp: failed to add mini app, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniApp.ID, err)
		return fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	p.logger.Infof("AddMiniApp: successfully added mini app, merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)
	return nil
}

func (p *miniAppMerchantPersistence) UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error {
	p.logger.Infof("UpdateMiniAppEnabledState: merchant_id=%s, mini_app_id=%s, enabled=%v", merchantID, miniAppID, enabled)

	objID, err := mappers.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Warnf("UpdateMiniAppEnabledState: invalid merchant_id=%s, error=%v", merchantID, err)
		return err
	}

	filter := bson.M{
		"_id":                  objID,
		"mini_apps.id":         miniAppID,
		"mini_apps.is_deleted": false,
		"is_deleted":           false,
	}
	update := bson.M{"mini_apps.$.enabled": enabled}

	_, err = p.MongoDalMiniApp.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("UpdateMiniAppEnabledState: mini app not found, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
			return fmt.Errorf(common_util.NotFound)
		}
		p.logger.Warnf("UpdateMiniAppEnabledState: failed to update enabled state, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniAppID, err)
		return fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	p.logger.Infof("UpdateMiniAppEnabledState: successfully updated enabled=%v, merchant_id=%s, mini_app_id=%s", enabled, merchantID, miniAppID)
	return nil
}

func (p *miniAppMerchantPersistence) SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error {
	p.logger.Infof("SoftDeleteMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)

	objID, err := mappers.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Warnf("SoftDeleteMiniApp: invalid merchant_id=%s, error=%v", merchantID, err)
		return err
	}

	filter := bson.M{
		"_id":                  objID,
		"mini_apps.id":         miniAppID,
		"mini_apps.is_deleted": false,
		"is_deleted":           false,
	}
	update := bson.M{"mini_apps.$.is_deleted": true}

	_, err = p.MongoDalMiniApp.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("SoftDeleteMiniApp: mini app not found, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
			return fmt.Errorf(common_util.NotFound)
		}
		p.logger.Warnf("SoftDeleteMiniApp: failed to soft delete mini app, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniAppID, err)
		return fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	p.logger.Infof("SoftDeleteMiniApp: successfully soft deleted mini app, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
	return nil
}
