package miniappmerchant

import (
	"context"
	"errors"
	"fmt"
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
}

func NewMiniAppMerchantPersistence(client *mongo.Client, DB_name string, collection string, logger utils.Logger) miniApp.MiniAppMerchantRepository {
	return &miniAppMerchantPersistence{
		MongoDalMiniApp: dal.NewMongoDal[model.MiniAppMerchant, model.MiniAppMerchant](client, DB_name, collection),
		logger:          logger,
	}
}

func (p *miniAppMerchantPersistence) CreateMiniAppMerchant(ctx context.Context, merchant *entities.MiniAppMerchant) (*entities.MiniAppMerchant, error) {

	fmt.Println("THis is the persistence")
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

func (p *miniAppMerchantPersistence) ListMiniAppMerchant(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.MiniAppMerchant], error) {
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

func (p *miniAppMerchantPersistence) MiniAppMerchantInfoExists(ctx context.Context, data entities.CheckMiniAppMerchant) (bool, error) {
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

func buildUpdateSet(m *entities.MiniAppMerchant) bson.M {
	set := bson.M{}

	if m.MerchantType != "" {
		set["merchant_type"] = m.MerchantType
	}
	if m.MerchantName != "" {
		set["merchant_name"] = m.MerchantName
	}
	// Update KYC fields if present
	if m.KYC.Status != "" {
		set["kyc.status"] = m.KYC.Status
	}
	if m.KYC.Representative.Name != "" {
		set["kyc.representative.name"] = m.KYC.Representative.Name
	}
	if m.KYC.Representative.Email != "" {
		set["kyc.representative.email"] = m.KYC.Representative.Email
	}
	if m.KYC.Representative.Phone != "" {
		set["kyc.representative.phone"] = m.KYC.Representative.Phone
	}
	if m.PhoneNumber != "" {
		set["phone_number"] = m.PhoneNumber
	}
	if m.Email != "" {
		set["email"] = m.Email
	}
	if m.BankAccountNumber != "" {
		set["bank_account_number"] = m.BankAccountNumber
	}
	if len(m.Branches) > 0 {
		set["branches"] = m.Branches
	}
	if len(m.MiniAppIDs) > 0 {
		set["mini_apps"] = m.MiniAppIDs
	}
	set["enabled"] = m.Enabled // include enabled even if false, to allow toggling

	return set
}
