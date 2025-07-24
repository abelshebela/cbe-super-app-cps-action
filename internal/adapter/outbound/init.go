// Package adapter provides outbound adapters for database and external service interactions
package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cps_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	dep_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	portalCardDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	bpscalls "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/bps_calls"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	serviceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	passwordRuleOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	userOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_services"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type outboundStore struct {
	MongoDalDepartment        *infra_mongo.MongoDal[dep_entities.Department, dep_entities.Department]
	MongoDalCPSAction         *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
	MongoDalBPSUser           *infra_mongo.MongoDal[bps.BPSUser, bps.BPSUser]
	MongoDalServices          *infra_mongo.MongoDal[model.Service, model.Service]
	MongoDalMember            *infra_mongo.MongoDal[member.User, member.User]
	MongoDalAccounts          *infra_mongo.MongoDal[model.LinkedAccount, model.LinkedAccount]
	BpsCalls                  bpscalls.BpsCallsInterface
	MongoDalMiniApp           *infra_mongo.MongoDal[model.MiniApp, model.MiniApp]
	MongoDalCPSUser           *infra_mongo.MongoDal[model.CPSUser, model.CPSUser]
	MongoDalAccountValidation *infra_mongo.MongoDal[model.ValidationRule, model.ValidationRule]
	MongoDalServiceDetails    *infra_mongo.MongoDal[model.ServiceDetails, model.ServiceDetails]
	MongoDalPasswordRule      *infra_mongo.MongoDal[model.PasswordRule, model.PasswordRule]
	MongoDalPortalCard        *infra_mongo.MongoDal[model.Card, model.Card]
}

func NewOutboundPasswordRuleInfra(client *mongo.Client, dbName string, collectionNames []string) passwordRuleOutbound.OutboundPasswordRuleInfra {
	mongoDalPasswordRule := infra_mongo.NewMongoDal[model.PasswordRule, model.PasswordRule](client, dbName, collectionNames[0])
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[1])
	return &outboundStore{MongoDalPasswordRule: mongoDalPasswordRule, MongoDalCPSAction: mongoDalCPSAction}
}

func NewCPSUserPersistence(client *mongo.Client, dbName string, collectionNames []string) userOutbound.OutboundInfra {
	mongoDalDepartment := infra_mongo.NewMongoDal[dep_entities.Department, dep_entities.Department](client, dbName, collectionNames[6])
	mongoDalCPSUser := infra_mongo.NewMongoDal[model.CPSUser, model.CPSUser](client, dbName, collectionNames[0])
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[1])
	mongoDalBPSUser := infra_mongo.NewMongoDal[bps.BPSUser, bps.BPSUser](client, dbName, collectionNames[2])
	mongoDalAccountValidation := infra_mongo.NewMongoDal[model.ValidationRule, model.ValidationRule](client, dbName, collectionNames[3])
	mongoDalServiceDetails := infra_mongo.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, "service")
	mongoDalPortalCard := infra_mongo.NewMongoDal[model.Card, model.Card](client, dbName, collectionNames[5])

	return &outboundStore{
		MongoDalDepartment:        mongoDalDepartment,
		MongoDalCPSUser:           mongoDalCPSUser,
		MongoDalCPSAction:         mongoDalCPSAction,
		MongoDalBPSUser:           mongoDalBPSUser,
		MongoDalAccountValidation: mongoDalAccountValidation,
		MongoDalServiceDetails:    mongoDalServiceDetails,
		MongoDalPortalCard:        mongoDalPortalCard,
	}
}

func NewOutBoundStore(client *mongo.Client, dbName string, collectionNames []string, logger utils.Logger, cfg *config.VaultConfig) outbound.OutboundInfra {

	mongoDalBPSUser := infra_mongo.NewMongoDal[bps.BPSUser, bps.BPSUser](client, dbName, collectionNames[0])
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[1])
	mongoDalCPSUser := infra_mongo.NewMongoDal[model.CPSUser, model.CPSUser](client, dbName, collectionNames[2])
	mongoDalService := infra_mongo.NewMongoDal[model.Service, model.Service](client, dbName, collectionNames[3])
	mongoDalMember := infra_mongo.NewMongoDal[member.User, member.User](client, dbName, collectionNames[4])
	mongoDalAccounts := infra_mongo.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, dbName, collectionNames[5])
	mongoDalMiniApp := infra_mongo.NewMongoDal[model.MiniApp, model.MiniApp](client, dbName, collectionNames[6])
	mongoDalServiceDetail := infra_mongo.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, collectionNames[3])
	MongoDalPortalCard := infra_mongo.NewMongoDal[model.Card, model.Card](client, dbName, collectionNames[7])

	return &outboundStore{
		MongoDalCPSAction:      mongoDalCPSAction,
		MongoDalBPSUser:        mongoDalBPSUser,
		MongoDalServices:       mongoDalService,
		MongoDalMember:         mongoDalMember,
		MongoDalAccounts:       mongoDalAccounts,
		BpsCalls:               bpscalls.NewBpsCalls(cfg, logger),
		MongoDalMiniApp:        mongoDalMiniApp,
		MongoDalCPSUser:        mongoDalCPSUser,
		MongoDalServiceDetails: mongoDalServiceDetail,

		MongoDalPortalCard: MongoDalPortalCard,
	}
}

func NewPortalCardPersistence(client *mongo.Client, dbName string, logger utils.Logger) *outboundStore {
	mongoDalPortalCard := infra_mongo.NewMongoDal[model.Card, model.Card](client, "cbe", "portal_cards")

	return &outboundStore{
		MongoDalPortalCard: mongoDalPortalCard,
	}
}

func NewServiceDetailsPersistence(client *mongo.Client, dbName string, collection []string, logger utils.Logger) *outboundStore {
	mongoDalServiceDetails := infra_mongo.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, collection[0])
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collection[1])

	return &outboundStore{
		MongoDalServiceDetails: mongoDalServiceDetails,
		MongoDalCPSAction:      mongoDalCPSAction,
	}
}

func (o *outboundStore) GetAllHqServices(ctx context.Context) ([]domain.ServiceDetails, error) {
	data, err := o.MongoDalServiceDetails.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	var d []domain.ServiceDetails
	for _, v := range data {
		x := domain.ServiceDetails{
			ID:          stringPointer(v.ID.Hex()),
			ServiceCode: "",
			ServiceName: v.ServiceName,
			ServiceType: "",
			Key:         v.Key,
			Cap: domain.Cap{
				KYCLevel:  domain.KYCLevel(v.Cap.KYCLevel),
				SingleCap: v.Cap.SingleCap,
				DailyCap:  v.Cap.DailyCap,
				MinAmount: v.Cap.MinAmount,
			},
			CBEProductCodes: domain.ProductCodes{
				PRD:    v.CBEProductCodes.PRD,
				VATPRD: v.CBEProductCodes.VATPRD,
				SFPRD:  v.CBEProductCodes.SFPRD,
				TRXN:   v.CBEProductCodes.TRXN,
			},
			CBEIFBProductCodes: domain.ProductCodes{
				PRD:    v.CBEIFBProductCodes.PRD,
				VATPRD: v.CBEIFBProductCodes.VATPRD,
				SFPRD:  v.CBEIFBProductCodes.SFPRD,
				TRXN:   v.CBEIFBProductCodes.TRXN,
			},
			AboveAmount:     v.AboveAmount,
			AboveServiceFee: v.AboveServiceFee,
			PaymentType:     v.PaymentType,
			Tiers: func() []domain.Tier {
				var tiers []domain.Tier
				for _, tier := range v.Tiers {
					tiers = append(tiers, domain.Tier{
						ID:        stringPointer(tier.ID.Hex()),
						Min:       tier.Min,
						Max:       tier.Max,
						FeeAmount: tier.FeeAmount,
					})
				}
				return tiers
			}(),
			CBEGLEntry: domain.GLEntry{
				ProductAccount:    v.CBEGLEntry.ProductAccount,
				ProductBranchCode: v.CBEGLEntry.ProductBranchCode,
				ServiceAccount:    v.CBEGLEntry.ServiceAccount,
				ServiceBranchCode: v.CBEGLEntry.ServiceBranchCode,
				VatAccount:        v.CBEGLEntry.VatAccount,
				VatBranchCode:     v.CBEGLEntry.VatBranchCode,
			},
			CBEIFBGLEntry: domain.GLEntry{
				ProductAccount:    v.CBEIFBGLEntry.ProductAccount,
				ProductBranchCode: v.CBEIFBGLEntry.ProductBranchCode,
				ServiceAccount:    v.CBEIFBGLEntry.ServiceAccount,
				ServiceBranchCode: v.CBEIFBGLEntry.ServiceBranchCode,
				VatAccount:        v.CBEIFBGLEntry.VatAccount,
				VatBranchCode:     v.CBEIFBGLEntry.VatBranchCode,
			},
			Enabled:        v.Enabled,
			IsDeleted:      v.IsDeleted,
			CreatedAt:      v.CreatedAt,
			LastModifiedAt: v.LastModifiedAt,
			DeletedAt:      v.DeletedAt,
		}
		d = append(d, x)
	}
	return d, nil
}
func (o *outboundStore) GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, error) {
	skip := int64(offset)
	limitInt64 := int64(limit)

	data, err := o.MongoDalServiceDetails.FindAllWithPagination(ctx, nil, nil, skip, limitInt64)
	if err != nil {
		return nil, err
	}

	var d []domain.ServiceDetails
	for _, v := range data {
		x := domain.ServiceDetails{
			ID:          stringPointer(v.ID.Hex()), // Convert bson.ObjectID to *string
			ServiceCode: "",
			ServiceName: v.ServiceName,
			ServiceType: "",
			Key:         v.Key,
			Cap: domain.Cap{
				KYCLevel:  domain.KYCLevel(v.Cap.KYCLevel),
				SingleCap: v.Cap.SingleCap,
				DailyCap:  v.Cap.DailyCap,
				MinAmount: v.Cap.MinAmount,
			},
			CBEProductCodes: domain.ProductCodes{
				PRD:    v.CBEProductCodes.PRD,
				VATPRD: v.CBEProductCodes.VATPRD,
				SFPRD:  v.CBEProductCodes.SFPRD,
				TRXN:   v.CBEProductCodes.TRXN,
			},
			CBEIFBProductCodes: domain.ProductCodes{
				PRD:    v.CBEIFBProductCodes.PRD,
				VATPRD: v.CBEIFBProductCodes.VATPRD,
				SFPRD:  v.CBEIFBProductCodes.SFPRD,
				TRXN:   v.CBEIFBProductCodes.TRXN,
			},
			AboveAmount:     v.AboveAmount,
			AboveServiceFee: v.AboveServiceFee,
			PaymentType:     v.PaymentType,
			Tiers: func() []domain.Tier {
				var tiers []domain.Tier
				for _, tier := range v.Tiers {
					tiers = append(tiers, domain.Tier{
						ID:        stringPointer(tier.ID.Hex()),
						Min:       tier.Min,
						Max:       tier.Max,
						FeeAmount: tier.FeeAmount,
					})
				}
				return tiers
			}(),
			CBEGLEntry: domain.GLEntry{
				ProductAccount:    v.CBEGLEntry.ProductAccount,
				ProductBranchCode: v.CBEGLEntry.ProductBranchCode,
				ServiceAccount:    v.CBEGLEntry.ServiceAccount,
				ServiceBranchCode: v.CBEGLEntry.ServiceBranchCode,
				VatAccount:        v.CBEGLEntry.VatAccount,
				VatBranchCode:     v.CBEGLEntry.VatBranchCode,
			},
			CBEIFBGLEntry: domain.GLEntry{
				ProductAccount:    v.CBEIFBGLEntry.ProductAccount,
				ProductBranchCode: v.CBEIFBGLEntry.ProductBranchCode,
				ServiceAccount:    v.CBEIFBGLEntry.ServiceAccount,
				ServiceBranchCode: v.CBEIFBGLEntry.ServiceBranchCode,
				VatAccount:        v.CBEIFBGLEntry.VatAccount,
				VatBranchCode:     v.CBEIFBGLEntry.VatBranchCode,
			},
			Enabled:        v.Enabled,
			IsDeleted:      v.IsDeleted,
			CreatedAt:      v.CreatedAt,
			LastModifiedAt: v.LastModifiedAt,
			DeletedAt:      v.DeletedAt,
		}
		d = append(d, x)
	}
	return d, nil
}

func (o *outboundStore) GetHqServiceById(ctx context.Context, id string) (domain.ServiceDetails, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ServiceDetails{}, err
	}
	filter := bson.M{"_id": objID}
	// fmt.Println("chkkkkkkkkkkkkkkkkk", filter)
	data, err := o.MongoDalServiceDetails.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ServiceDetails{}, fmt.Errorf("NOT_FOUND")
		}
		return domain.ServiceDetails{}, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}
	return domain.ServiceDetails{
		ID:          stringPointer(data.ID.Hex()), // Convert bson.ObjectID to *string
		ServiceCode: "",
		ServiceName: data.ServiceName,
		ServiceType: "",
		Key:         data.Key,
		Cap: domain.Cap{
			KYCLevel:  domain.KYCLevel(data.Cap.KYCLevel),
			SingleCap: data.Cap.SingleCap,
			DailyCap:  data.Cap.DailyCap,
			MinAmount: data.Cap.MinAmount,
		},
		CBEProductCodes: domain.ProductCodes{
			PRD:    data.CBEProductCodes.PRD,
			VATPRD: data.CBEProductCodes.VATPRD,
			SFPRD:  data.CBEProductCodes.SFPRD,
			TRXN:   data.CBEProductCodes.TRXN,
		},
		CBEIFBProductCodes: domain.ProductCodes{
			PRD:    data.CBEIFBProductCodes.PRD,
			VATPRD: data.CBEIFBProductCodes.VATPRD,
			SFPRD:  data.CBEIFBProductCodes.SFPRD,
			TRXN:   data.CBEIFBProductCodes.TRXN,
		},
		AboveAmount:     data.AboveAmount,
		AboveServiceFee: data.AboveServiceFee,
		PaymentType:     data.PaymentType,
		Tiers: func() []domain.Tier {
			var tiers []domain.Tier
			for _, tier := range data.Tiers {
				tiers = append(tiers, domain.Tier{
					ID:        stringPointer(tier.ID.Hex()),
					Min:       tier.Min,
					Max:       tier.Max,
					FeeAmount: tier.FeeAmount,
				})
			}
			return tiers
		}(),
		CBEGLEntry: domain.GLEntry{
			ProductAccount:    data.CBEGLEntry.ProductAccount,
			ProductBranchCode: data.CBEGLEntry.ProductBranchCode,
			ServiceAccount:    data.CBEGLEntry.ServiceAccount,
			ServiceBranchCode: data.CBEGLEntry.ServiceBranchCode,
			VatAccount:        data.CBEGLEntry.VatAccount,
			VatBranchCode:     data.CBEGLEntry.VatBranchCode,
		},
		CBEIFBGLEntry: domain.GLEntry{
			ProductAccount:    data.CBEIFBGLEntry.ProductAccount,
			ProductBranchCode: data.CBEIFBGLEntry.ProductBranchCode,
			ServiceAccount:    data.CBEIFBGLEntry.ServiceAccount,
			ServiceBranchCode: data.CBEIFBGLEntry.ServiceBranchCode,
			VatAccount:        data.CBEIFBGLEntry.VatAccount,
			VatBranchCode:     data.CBEIFBGLEntry.VatBranchCode,
		},
		Enabled:        data.Enabled,
		IsDeleted:      data.IsDeleted,
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
		DeletedAt:      data.DeletedAt,
	}, nil
}
func (o *outboundStore) UpdateHqService(ctx context.Context, service domain.ServiceDetails) error {

	objectId, err := bson.ObjectIDFromHex(*service.ID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": objectId,
	}
	update := bson.M{

		"key":         service.Key,
		"serviceName": service.ServiceName,
		"serviceType": service.ServiceType,
		"cap": bson.M{
			"kyc_level":  service.Cap.KYCLevel,
			"single_cap": service.Cap.SingleCap,
			"daily_cap":  service.Cap.DailyCap,
			"min_amount": service.Cap.MinAmount,
		},
		"cbe_product_codes": bson.M{
			"prd":    service.CBEProductCodes.PRD,
			"vatprd": service.CBEProductCodes.VATPRD,
			"sfprd":  service.CBEProductCodes.SFPRD,
			"trxn":   service.CBEProductCodes.TRXN,
		},
		"cbe_ifb_product_codes": bson.M{
			"prd":    service.CBEIFBProductCodes.PRD,
			"vatprd": service.CBEIFBProductCodes.VATPRD,
			"sfprd":  service.CBEIFBProductCodes.SFPRD,
			"trxn":   service.CBEIFBProductCodes.TRXN,
		},
		"above_amount":      service.AboveAmount,
		"above_service_fee": service.AboveServiceFee,
		"payment_type":      service.PaymentType,
		"tiers": func() []bson.M {
			var tiers []bson.M
			for _, tier := range service.Tiers {
				tiers = append(tiers, bson.M{
					"id":         tier.ID,
					"min":        tier.Min,
					"max":        tier.Max,
					"fee_amount": tier.FeeAmount,
				})
			}
			return tiers
		}(),
		"cbe_gl_entry": bson.M{
			"product_account":     service.CBEGLEntry.ProductAccount,
			"product_branch_code": service.CBEGLEntry.ProductBranchCode,
			"service_account":     service.CBEGLEntry.ServiceAccount,
			"service_branch_code": service.CBEGLEntry.ServiceBranchCode,
			"vat_account":         service.CBEGLEntry.VatAccount,
			"vat_branch_code":     service.CBEGLEntry.VatBranchCode,
		},
		"cbe_ifb_gl_entry": bson.M{
			"product_account":     service.CBEIFBGLEntry.ProductAccount,
			"product_branch_code": service.CBEIFBGLEntry.ProductBranchCode,
			"service_account":     service.CBEIFBGLEntry.ServiceAccount,
			"service_branch_code": service.CBEIFBGLEntry.ServiceBranchCode,
			"vat_account":         service.CBEIFBGLEntry.VatAccount,
			"vat_branch_code":     service.CBEIFBGLEntry.VatBranchCode,
		},
		"enabled":          service.Enabled,
		"is_deleted":       service.IsDeleted,
		"last_modified_at": service.LastModifiedAt,
	}

	_, err = o.MongoDalServiceDetails.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)

	}
	return nil
}
func stringPointer(s string) *string {
	return &s
}

func domainToModelCPSAction(domainAction domain.CPSAction) model.CPSAction {

	var objID bson.ObjectID
	if domainAction.ID != "" {
		var err error
		objID, err = bson.ObjectIDFromHex(domainAction.ID)
		if err != nil {
			objID = bson.NilObjectID
		}
	} else {
		objID = bson.NewObjectID()
	}

	// Convert CurrentAction to map[string]interface{} if needed
	var currentAction interface{}
	switch v := domainAction.CurrentAction.(type) {
	case map[string]interface{}:
		currentAction = v
	default:
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &currentAction)
	}

	return model.CPSAction{
		ID:                 objID,
		ActionCode:         domainAction.ActionCode,
		UniqueId:           domainAction.UniqueId,
		MakerID:            domainAction.MakerID,
		MakerName:          domainAction.MakerName,
		MakerPhoneNumber:   domainAction.MakerPhoneNumber,
		CheckerID:          domainAction.CheckerID,
		CheckerName:        domainAction.CheckerName,
		CheckerPhoneNumber: domainAction.CheckerPhoneNumber,
		Department:         domainAction.Department,
		RejectionReason: func() string {
			if domainAction.RejectionReason != nil {
				return *domainAction.RejectionReason
			}
			return ""
		}(),
		PreviousAction:    domainAction.PreviousAction,
		CurrentAction:     currentAction,
		ActionStatus:      string(domainAction.ActionStatus),
		ActionType:        string(domainAction.ActionType),
		RequestAction:     string(domainAction.RequestAction),
		CreatedAt:         domainAction.CreatedAt,
		LastModifiedAt:    domainAction.LastModifiedAt,
		MakerActionTime:   domainAction.MakerActionTime,
		CheckerActionTime: domainAction.CheckerActionTime,
	}
}

// Helper to recursively convert bson.D to map[string]interface{}
func bsonDToMap(i interface{}) interface{} {
	switch v := i.(type) {
	case bson.D:
		m := make(map[string]interface{})
		for _, e := range v {
			m[e.Key] = bsonDToMap(e.Value)
		}
		return m
	case []interface{}:
		for i, e := range v {
			v[i] = bsonDToMap(e)
		}
		return v
	default:
		return v
	}
}

func modelToDomainCPSAction(modelAction model.CPSAction) domain.CPSAction {
	var rejectionReason *string
	if modelAction.RejectionReason != "" {
		rejectionReason = &modelAction.RejectionReason
	}

	return domain.CPSAction{
		ID:                 modelAction.ID.Hex(),
		ActionCode:         modelAction.ActionCode,
		UniqueId:           modelAction.UniqueId,
		MakerID:            modelAction.MakerID,
		MakerName:          modelAction.MakerName,
		MakerPhoneNumber:   modelAction.MakerPhoneNumber,
		CheckerID:          modelAction.CheckerID,
		CheckerName:        modelAction.CheckerName,
		CheckerPhoneNumber: modelAction.CheckerPhoneNumber,
		Department:         modelAction.Department,
		RejectionReason:    rejectionReason,
		PreviousAction:     modelAction.PreviousAction,
		CurrentAction:      bsonDToMap(modelAction.CurrentAction), // always map/slice, never bson.D
		ActionStatus:       domain.ActionStatus(modelAction.ActionStatus),
		ActionType:         domain.ActionType(modelAction.ActionType),
		RequestAction:      domain.RequestAction(modelAction.RequestAction),
		CreatedAt:          modelAction.CreatedAt,
		LastModifiedAt:     modelAction.LastModifiedAt,
		MakerActionTime:    modelAction.MakerActionTime,
		CheckerActionTime:  modelAction.CheckerActionTime,
	}
}

func (o *outboundStore) CreateCpsAction(ctx context.Context, Action domain.CPSAction) (domain.CPSAction, error) {
	Action.ActionCode = utils.RandomGenerator(20)

	modelAction := domainToModelCPSAction(Action)
	data, err := o.MongoDalCPSAction.InsertOne(ctx, modelAction)
	if err != nil {
		return domain.CPSAction{}, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}
	return modelToDomainCPSAction(data), nil
}
func (o *outboundStore) UpdateCpsAction(ctx context.Context, Action domain.CPSAction) error {
	modelAction := domainToModelCPSAction(Action)
	filter := bson.M{
		"action_code": modelAction.ActionCode,
	}
	update := bson.M{
		"maker_id":             modelAction.MakerID,
		"maker_name":           modelAction.MakerName,
		"maker_phone_number":   modelAction.MakerPhoneNumber,
		"checker_id":           modelAction.CheckerID,
		"checker_name":         modelAction.CheckerName,
		"checker_phone_number": modelAction.CheckerPhoneNumber,
		"unique_id":            modelAction.UniqueId,
		"department":           modelAction.Department,
		"rejection_reason":     modelAction.RejectionReason,
		"previous_action":      modelAction.PreviousAction,
		"current_action":       modelAction.CurrentAction,
		"action_status":        modelAction.ActionStatus,
		"action_type":          modelAction.ActionType,
		"request_action":       modelAction.RequestAction,
		"created_at":           modelAction.CreatedAt,
		"last_modified_at":     modelAction.LastModifiedAt,
	}
	fmt.Println("Update cps action********", filter)
	if _, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)

	}
	fmt.Println("Update cps action/////////////////", filter)

	return nil
}
func (o *outboundStore) FetchCpsActionById(ctx context.Context, Action_code string) (domain.CPSAction, error) {
	filter := bson.M{"action_code": Action_code}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.CPSAction{}, fmt.Errorf(error_codes.ActionNotFound)
		}
		return domain.CPSAction{}, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}

	return modelToDomainCPSAction(*data), nil

}

func (o *outboundStore) FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]domain.LinkedAccount, error) {
	v, err := o.BpsCalls.FetchLinkedAccount(accountNumber)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf(error_codes.AccountNotFound)
	}
	filter := map[string]interface{}{"customer_number": v.CustomerNumber}
	data, err := o.MongoDalAccounts.FindAll(ctx, filter, nil)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}
	var accounts []domain.LinkedAccount
	for _, v := range data {
		account := domain.LinkedAccount{
			ID:                stringToPointer(v.ID.Hex()),
			UserID:            stringToPointer(v.UserID.Hex()),
			CustomerNumber:    v.CustomerNumber,
			AccountNumber:     v.AccountNumber,
			AccountHolderName: v.AccountHolderName,
			AccountType:       v.AccountType,
			BranchCode:        v.BranchCode,
			LinkedStatus:      v.LinkedStatus,
			LastLinkedStatus:  v.LastLinkedStatus,
			LinkedAt:          v.LinkedAt,
			LinkedBranch:      v.LinkedBranch,
			RegistrationType:  domain.RegistrationType(v.RegistrationType),
			IsAccountActive:   v.IsAccountActive,
			AndOrStatus:       v.AndOrStatus,
			AccountBranchCode: v.AccountBranchCode,
			CurrencyCode:      v.CurrencyCode,
			IsMain:            v.IsMain,
			MakerAndChecker: struct {
				Linkers struct {
					Maker   string
					Checker string
				}
				Unlinkers struct {
					Maker   string
					Checker string
				}
			}{
				Linkers: struct {
					Maker   string
					Checker string
				}{
					Maker:   v.MakerAndChecker.Linkers.Maker,
					Checker: v.MakerAndChecker.Linkers.Checker,
				},
				Unlinkers: struct {
					Maker   string
					Checker string
				}{
					Maker:   v.MakerAndChecker.Unlinkers.Maker,
					Checker: v.MakerAndChecker.Unlinkers.Checker,
				},
			},
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (o *outboundStore) UpdateAccounts(ctx context.Context, linkedAccounts []domain.LinkedAccount) ([]domain.LinkedAccount, error) {
	var updatedAccounts []domain.LinkedAccount
	for _, acc := range linkedAccounts {
		filter := map[string]interface{}{"_id": acc.ID}
		update := map[string]interface{}{

			"user_id":             acc.UserID,
			"customer_number":     acc.CustomerNumber,
			"account_number":      acc.AccountNumber,
			"account_holder_name": acc.AccountHolderName,
			"account_type":        acc.AccountType,
			"branch_code":         acc.BranchCode,
			"linked_status":       acc.LinkedStatus,
			"last_linked_status":  acc.LastLinkedStatus,
			"linked_at":           acc.LinkedAt,
			"linker_branch":       acc.LinkedBranch,
			"registration_type":   acc.RegistrationType,
			"is_account_active":   acc.IsAccountActive,
			"and_or_status":       acc.AndOrStatus,
			"account_branch_code": acc.AccountBranchCode,
			"currency":            acc.CurrencyCode,
			"is_main":             acc.IsMain,
			"maker_and_checker": map[string]interface{}{
				"linkers": map[string]interface{}{
					"maker":   acc.MakerAndChecker.Linkers.Maker,
					"checker": acc.MakerAndChecker.Linkers.Checker,
				},
				"unlinkers": map[string]interface{}{
					"maker":   acc.MakerAndChecker.Unlinkers.Maker,
					"checker": acc.MakerAndChecker.Unlinkers.Checker,
				},
			},
			"created_at": acc.CreatedAt,
			"updated_at": acc.UpdatedAt,
		}

		_, err := o.MongoDalAccounts.UpdateOne(ctx, filter, update)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf(error_codes.AccountNotFound)
			}
			return nil, err
		}
		updatedAccounts = append(updatedAccounts, acc)
	}
	return updatedAccounts, nil
}

func (o *outboundStore) UpdateAccount(ctx context.Context, linkedAccount domain.LinkedAccount) (domain.LinkedAccount, error) {
	objID, err := bson.ObjectIDFromHex(*linkedAccount.ID)
	if err != nil {
		return domain.LinkedAccount{}, fmt.Errorf(error_codes.AccountNotFound)
	}
	filter := map[string]interface{}{"_id": objID}
	update := map[string]interface{}{

		"user_id":             linkedAccount.UserID,
		"customer_number":     linkedAccount.CustomerNumber,
		"account_number":      linkedAccount.AccountNumber,
		"account_holder_name": linkedAccount.AccountHolderName,
		"account_type":        linkedAccount.AccountType,
		"branch_code":         linkedAccount.BranchCode,
		"linked_status":       linkedAccount.LinkedStatus,
		"last_linked_status":  linkedAccount.LastLinkedStatus,
		"linked_at":           linkedAccount.LinkedAt,
		"linker_branch":       linkedAccount.LinkedBranch,
		"registration_type":   linkedAccount.RegistrationType,
		"is_account_active":   linkedAccount.IsAccountActive,
		"and_or_status":       linkedAccount.AndOrStatus,
		"account_branch_code": linkedAccount.AccountBranchCode,
		"currency":            linkedAccount.CurrencyCode,
		"is_main":             linkedAccount.IsMain,
		"maker_and_checker": map[string]interface{}{
			"linkers": map[string]interface{}{
				"maker":   linkedAccount.MakerAndChecker.Linkers.Maker,
				"checker": linkedAccount.MakerAndChecker.Linkers.Checker,
			},
			"unlinkers": map[string]interface{}{
				"maker":   linkedAccount.MakerAndChecker.Unlinkers.Maker,
				"checker": linkedAccount.MakerAndChecker.Unlinkers.Checker,
			},
		},
		"created_at": linkedAccount.CreatedAt,
		"updated_at": linkedAccount.UpdatedAt,
	}

	_, err = o.MongoDalAccounts.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.LinkedAccount{}, fmt.Errorf(error_codes.AccountNotFound)
		}
		return domain.LinkedAccount{}, err
	}
	return linkedAccount, nil
}

func (o *outboundStore) FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (domain.CPSAction, error) {
	filter := map[string]interface{}{"maker_id": makerId}

	d, err := o.MongoDalCPSAction.FindAll(ctx, filter, nil)
	if err != nil {
		return domain.CPSAction{}, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}

	if len(d) == 0 {
		return domain.CPSAction{}, nil

	}
	return modelToDomainCPSAction(*d[len(d)-1]), nil
}
func (o *outboundStore) GetAllPortalCard(ctx context.Context) ([]*portalCardDomain.Card, error) {

	data, err := o.MongoDalPortalCard.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*portalCardDomain.Card, len(data))

	for i, s := range data {
		if s == nil {
			continue
		}
		result[i] = &portalCardDomain.Card{
			ID:       s.ID,
			CardName: s.CardName,
			SubCards: s.SubCards,
		}
	}
	return result, nil

}

func (o *outboundStore) FetchLinkedAccountById(ctx context.Context, ids []string) ([]domain.LinkedAccount, error) {
	var result []domain.LinkedAccount
	for _, i := range ids {

		filter := bson.M{"customer_number": i}
		item, err := o.MongoDalAccounts.FindOne(ctx, filter, nil)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf(error_codes.AccountNotFound)
			}
			return nil, fmt.Errorf(err.Error())
		}
		result = append(result, domain.LinkedAccount{
			ID:                stringToPointer(item.ID.Hex()),
			UserID:            stringToPointer(item.UserID.Hex()),
			CustomerNumber:    item.CustomerNumber,
			AccountNumber:     item.AccountNumber,
			AccountHolderName: item.AccountHolderName,
			AccountType:       item.AccountType,
			BranchCode:        item.BranchCode,
			LinkedStatus:      item.LinkedStatus,
			LastLinkedStatus:  item.LastLinkedStatus,
			LinkedAt:          item.LinkedAt,
			LinkedBranch:      item.LinkedBranch,
			RegistrationType:  domain.RegistrationType(item.RegistrationType),
			IsAccountActive:   item.IsAccountActive,
			AndOrStatus:       item.AndOrStatus,
			AccountBranchCode: item.AccountBranchCode,
			CurrencyCode:      item.CurrencyCode,
			IsMain:            item.IsMain,
			MakerAndChecker: struct {
				Linkers struct {
					Maker   string
					Checker string
				}
				Unlinkers struct {
					Maker   string
					Checker string
				}
			}{
				Linkers: struct {
					Maker   string
					Checker string
				}{
					Maker:   item.MakerAndChecker.Linkers.Maker,
					Checker: item.MakerAndChecker.Linkers.Checker,
				},
				Unlinkers: struct {
					Maker   string
					Checker string
				}{
					Maker:   item.MakerAndChecker.Unlinkers.Maker,
					Checker: item.MakerAndChecker.Unlinkers.Checker,
				},
			},
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result, nil
}

func stringToPointer(s string) *string {
	return &s
}

func (o *outboundStore) CreateUserRequest(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error) {
	projection := bson.M{}
	prevFilter := bson.M{
		"unique_id":  cpsAction.UniqueId,
		"maker_id":   cpsAction.MakerID,
		"maker_name": cpsAction.MakerName,
	}
	previosCPSAction, err := o.MongoDalCPSAction.FindRecentDocument(ctx, prevFilter, projection)
	if err == nil && previosCPSAction != nil {
		cpsAction.PreviousAction = previosCPSAction.CurrentAction
	}

	data, err := o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return nil, fmt.Errorf("database error while creating user request action")
	}

	return &data, nil
}

func (o *outboundStore) UpdateUserRequest(ctx context.Context, cpsAction model.CPSAction, userCode string) (*model.CPSAction, error) {
	cpsUserFilter := bson.M{
		"user_code":  userCode,
		"is_deleted": false,
	}

	projection := bson.M{}

	// Check if the user does exist
	_, err := o.MongoDalCPSUser.FindOne(ctx, cpsUserFilter, projection)
	if err != nil {
		return nil, err
	}

	prevFilter := bson.M{
		"unique_id":  cpsAction.UniqueId,
		"maker_id":   cpsAction.MakerID,
		"maker_name": cpsAction.MakerName,
	}
	previosCPSAction, err := o.MongoDalCPSAction.FindRecentDocument(ctx, prevFilter, projection)
	if err == nil && previosCPSAction != nil {
		cpsAction.PreviousAction = previosCPSAction.CurrentAction
	}

	// Create CPS action
	data, err := o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return nil, fmt.Errorf("database error while updating user request action")
	}

	return &data, nil
}

func (o *outboundStore) DeleteUserRequest(ctx context.Context, userCode string, cpsAction model.CPSAction) (*model.CPSAction, error) {
	filter := bson.M{"user_code": userCode}
	// Check if the user exists
	_, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("USER_NOT_FOUND")
		}
		return nil, fmt.Errorf("database error while finding user")
	}

	// If user exists, create a CPS action for deletion
	cpsAction, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return nil, fmt.Errorf("database error while creating user request action")
	}

	return &cpsAction, nil
}

func (o *outboundStore) ApproveUserAction(ctx context.Context, actionCode string, cpsAction model.CPSAction) (*model.CPSAction, error) {
	// Find the CPS action
	filter := bson.M{"action_code": actionCode}
	projection := bson.M{}

	cps_action, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("CPS action not found")
	} else if err != nil {
		return nil, fmt.Errorf("database error while finding CPS action")
	}

	if cps_action.ActionStatus == string(model.ActionApproved) {
		return nil, fmt.Errorf("ACTION_HAS_ALREADY_APPROVED")
	} else if cps_action.ActionStatus == string(model.ActionRejected) {
		return nil, fmt.Errorf("ACTION_HAS_ALREADY_REJECTED")
	}

	now := time.Now()

	// Update it
	cps_action.CheckerID = cpsAction.CheckerID
	cps_action.CheckerName = cpsAction.CheckerName
	cps_action.CheckerPhoneNumber = cpsAction.CheckerPhoneNumber
	cps_action.Department = cpsAction.Department
	cps_action.ActionStatus = cpsAction.ActionStatus
	cps_action.RejectionReason = cpsAction.RejectionReason
	cps_action.LastModifiedAt = time.Now()
	cps_action.CheckerActionTime = &now

	// Conver the struct to bson.M
	updateBytes, err := bson.Marshal(cps_action)
	if err != nil {
		return nil, err
	}

	var updateData bson.M
	if err := bson.Unmarshal(updateBytes, &updateData); err != nil {
		return nil, err
	}

	// Apply the update
	updatedCPSAction, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return nil, err
	}

	// If cps action is rejected, return early
	if updatedCPSAction.ActionStatus == string(model.ActionRejected) {
		return nil, nil
	}

	// Update the users data
	data, err := json.Marshal(updatedCPSAction.CurrentAction)
	if err != nil {
		return nil, err
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, err
	}

	userCode, _ := dataMap["user_code"].(string)
	userName, _ := dataMap["username"].(string)
	fullName, _ := dataMap["full_name"].(string)
	phoneNumber, _ := dataMap["phone_number"].(string)
	role, _ := dataMap["role"].(string)
	gender, _ := dataMap["gender"].(string)
	email, _ := dataMap["email"].(string)
	realm, _ := dataMap["realm"].(string)
	country, _ := dataMap["country"].(string)
	region, _ := dataMap["region"].(string)
	deptStr, ok := dataMap["department"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing department")
	}
	department, err := bson.ObjectIDFromHex(deptStr)
	if err != nil {
		return nil, err
	}

	// Convert string slices to []bson.ObjectID
	var permissionCategory []bson.ObjectID
	if raw, ok := dataMap["permission_category"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				str, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("permission_category element is not a string")
				}
				id, err := bson.ObjectIDFromHex(str)
				if err != nil {
					return nil, err
				}
				permissionCategory = append(permissionCategory, id)
			}
		}
	}

	var permissionGroups []bson.ObjectID
	if raw, ok := dataMap["permissiongroups"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				str, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("permissiongroups element is not a string")
				}
				id, err := bson.ObjectIDFromHex(str)
				if err != nil {
					return nil, err
				}
				permissionGroups = append(permissionGroups, id)
			}
		}
	}

	// Creating a new user
	now = time.Now()

	newUser := model.CPSUser{}
	newUser.ID = bson.NewObjectID()
	newUser.UserCode = userCode
	newUser.UserName = userName
	newUser.FullName = fullName
	newUser.Gender = gender
	newUser.Email = email
	newUser.Realm = realm
	newUser.PhoneNumber = phoneNumber
	newUser.Role = role
	newUser.Department = department
	newUser.Country = country
	newUser.Region = region
	newUser.PermissionCategory = permissionCategory
	newUser.PermissionGroup = permissionGroups
	newUser.PasswordDisable = false
	newUser.SyncDisabled = false
	newUser.LoginAttemptCount = 0
	newUser.LastLoginAttempt = time.Now()
	newUser.NextLoginAttempt = time.Now()
	newUser.LastOnlineDate = time.Now()
	newUser.LastLogin = time.Now()
	newUser.LoginPassword = ""
	newUser.AccountAuthorizationCode = ""
	newUser.UnlockAccountRequested = false
	newUser.PasswordChangedAt = nil
	newUser.OTPStatus = ""
	newUser.OTPLastTriedAt = nil
	newUser.OTPVerifyCount = 0
	newUser.Enabled = true
	newUser.IsDeleted = false
	newUser.DateJoined = &now

	filterUser := bson.M{"user_code": userCode}

	// If the action type is "CREATE", create a new user, else if it is "UPDATE" update the exsting user
	if updatedCPSAction.ActionType == string(model.ActionCreate) {
		_, err := o.MongoDalCPSUser.InsertOne(ctx, newUser)
		if err != nil {
			return nil, err
		}

	} else if updatedCPSAction.ActionType == string(model.ActionUpdate) {
		// Convert data (*model.CPSUser) to bson.M for update
		updateBytes, err := bson.Marshal(data)
		if err != nil {
			return nil, err
		}
		var updateDoc bson.M
		if err := bson.Unmarshal(updateBytes, &updateDoc); err != nil {
			return nil, err
		}

		// Update the user document
		_, err = o.MongoDalCPSUser.UpdateOne(ctx, filterUser, bson.M{"$set": updateDoc})
		if err != nil {
			return nil, err
		}
	}

	return &updatedCPSAction, nil
}

func (o *outboundStore) AuthorizeCreate(ctx context.Context, action *cps_entity.CPSAction) (*cps_entity.CPSAction, error) {
	// Unmarshal the CurrentAction to a model.CPSUser
	data, err := common_util.JsonUnmarshal[model.CPSUser](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	// Check for existing user by phone, email, username, or user_code
	filter := bson.M{
		"$or": []bson.M{
			{"phone_number": data.PhoneNumber},
			{"email": data.Email},
			{"username": data.UserName},
			{"user_code": data.UserCode},
		},
	}
	existing, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	// 	if existing != nil {
	// 	if data.UserCode == existing.UserCode ||
	// 		data.PhoneNumber == existing.PhoneNumber ||
	// 		data.Email == existing.Email ||
	// 		data.UserName == existing.UserName {
	// 		return nil, fmt.Errorf("CPS_USER_NOT_FOUND")
	// 	}
	// }

	if existing != nil {
		if data.UserCode == existing.UserCode {
			return nil, fmt.Errorf("USER_CODE_ALREADY_EXIST")
		} else if data.PhoneNumber == existing.PhoneNumber {
			return nil, fmt.Errorf("PHONE_NUMBER_EXISTS")
		} else if data.Email == existing.Email {
			return nil, fmt.Errorf("EMAIL_ALREADY_EXISTS")
		} else if data.UserName == existing.UserName {
			return nil, fmt.Errorf("USERNAME_ALREADY_EXISTS")
		}
	}

	// Add filds to the User collection
	now := time.Now()
	data.Country = ""
	data.Region = ""
	data.Enabled = true
	data.Realm = "bank"
	data.PasswordDisable = false
	data.SyncDisabled = false
	data.LoginAttemptCount = 0
	data.NextLoginAttempt = time.Now()
	data.LastLoginAttempt = time.Now()
	data.LastLogin = time.Now()
	data.LoginPassword = ""
	data.AccountAuthorizationCode = ""
	data.UnlockAccountRequested = false
	data.PasswordChangedAt = nil
	data.OTPStatus = ""
	data.OTPLastTriedAt = nil
	data.OTPVerifyCount = 0
	data.Enabled = true
	data.IsDeleted = false
	data.DateJoined = &now
	data.LastModified = &now

	// Insert the new user
	_, err = o.MongoDalCPSUser.InsertOne(ctx, *data)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (o *outboundStore) AuthorizeUpdate(ctx context.Context, action *cps_entity.CPSAction) (*cps_entity.CPSAction, error) {
	// Unmarshal the CurrentAction to a model.CPSUser
	data, err := common_util.JsonUnmarshal[model.CPSUser](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	// Find the user by user_code
	filter := bson.M{"user_code": data.UserCode}
	existing, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("USER_NOT_FOUND")
	}

	now := time.Now()
	data.LastModified = &now

	// Update the user document
	updateBytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	var updateDoc bson.M
	if err := bson.Unmarshal(updateBytes, &updateDoc); err != nil {
		return nil, err
	}

	_, err = o.MongoDalCPSUser.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (o *outboundStore) AuthorizeDelete(ctx context.Context, action *cps_entity.CPSAction) (*cps_entity.CPSAction, error) {
	// Unmarshal the CurrentAction to a model.CPSUser
	data, err := common_util.JsonUnmarshal[model.CPSUser](action.CurrentAction)
	if err != nil {
		return nil, err
	}
	// Find the user by user_code
	filter := bson.M{"user_code": data.UserCode}
	existing, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("USER_NOT_FOUND")
	}

	// Soft delete the user by setting is_deleted to true
	update := bson.M{"is_deleted": true}
	_, err = o.MongoDalCPSUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (o *outboundStore) GetDepartmentByID(ctx context.Context, id string) (*dep_entities.Department, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	department, err := o.MongoDalDepartment.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}
	if department == nil {
		return nil, fmt.Errorf("DEPARTMENT_NOT_FOUND")
	}

	return department, nil

}

func (o *outboundStore) GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error) {
	filter := bson.M{
		"action_status":  model.ActionPending,
		"request_action": bson.M{"$in": []string{string(model.RequestUser), string(model.RequestUpdateUser)}},
	}
	projection := bson.M{}

	data, err := o.MongoDalCPSAction.FindAll(ctx, filter, projection)
	if err != nil {
		return nil, err
	}

	// Convert []*model.CPSAction to []model.CPSAction
	result := make([]model.CPSAction, 0, len(data))
	for _, d := range data {
		if d != nil {
			result = append(result, *d)

		}
	}

	return result, nil
}

func (o *outboundStore) FetchUserByUserCode(ctx context.Context, userCode string) (*userDTO.CPSUserDTO, error) {
	filter := bson.M{"user_code": userCode}
	modelUser, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil && modelUser == nil {
		return nil, fmt.Errorf("CPS_USER_NOT_FOUND")
	}

	safeUser := userDTO.NewCPSUserDTO(*modelUser)
	return &safeUser, nil
}

func (o *outboundStore) GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*userDTO.CPSUserDTO], error) {
	filter := bson.M{
		// "is_deleted": false,
	}
	projection := bson.M{}

	// Add search functionality
	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"user_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"full_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"phone_number": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"email": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"role": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"department": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	// Add filter functionality
	if filterParams.Filters != "" {
		filter["role"] = filterParams.Filters
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	users, err := o.MongoDalCPSUser.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	// Covert to DTOs
	dtoUsers := make([]*userDTO.CPSUserDTO, 0, len(users))
	for _, user := range users {
		dto := userDTO.NewCPSUserDTO(*user)
		dtoUsers = append(dtoUsers, &dto)
	}

	totalDocs, err := o.MongoDalCPSUser.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := common_util.BuildPaginationMeta(totalDocs, page, limit)

	return &common_util.PaginatedResponse[[]*userDTO.CPSUserDTO]{
		Data: dtoUsers,
		Meta: meta,
	}, nil
}

func (o *outboundStore) GetAllServiceDetails(ctx context.Context) ([]*serviceDomain.Service, error) {
	data, err := o.MongoDalServiceDetails.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*serviceDomain.Service, len(data))
	for i, s := range data {
		if s == nil {
			continue
		}
		result[i] = &serviceDomain.Service{
			ID:                 s.ID.Hex(),
			ServiceCode:        s.ServiceCode,
			ServiceName:        s.ServiceName,
			ServiceType:        s.ServiceType,
			Key:                s.Key,
			Cap:                convertToServiceCap(s.Cap),
			CBEProductCodes:    convertToServiceProductCodes(s.CBEProductCodes),
			CBEIFBProductCodes: convertToServiceProductCodes(s.CBEIFBProductCodes),
			AboveAmount:        s.AboveAmount,
			AboveServiceFee:    s.AboveServiceFee,
			PaymentType:        s.PaymentType,
			Tiers:              convertTiers(s.Tiers),
			CBEGLEntry:         convertToServiceGLEntry(s.CBEGLEntry),
			CBEIFBGLEntry:      convertToServiceGLEntry(s.CBEIFBGLEntry),
			Enabled:            s.Enabled,
			IsDeleted:          s.IsDeleted,
			CreatedAt:          s.CreatedAt,
			LastModifiedAt:     s.LastModifiedAt,
			DeletedAt:          s.DeletedAt,
		}
	}
	return result, nil
}

func (o *outboundStore) GetOneServiceDetail(ctx context.Context, id string) (serviceDomain.Service, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return serviceDomain.Service{}, err
	}

	filter := bson.M{"_id": objID}

	data, err := o.MongoDalServiceDetails.FindOne(ctx, filter, nil)

	if err != nil {
		return serviceDomain.Service{}, err
	}

	return serviceDomain.Service{
		ID:                 data.ID.Hex(),
		ServiceCode:        data.ServiceCode,
		ServiceName:        data.ServiceName,
		ServiceType:        data.ServiceType,
		Key:                data.Key,
		Cap:                convertToServiceCap(data.Cap),
		CBEProductCodes:    convertToServiceProductCodes(data.CBEProductCodes),
		CBEIFBProductCodes: convertToServiceProductCodes(data.CBEIFBProductCodes),
		AboveAmount:        data.AboveAmount,
		AboveServiceFee:    data.AboveServiceFee,
		PaymentType:        data.PaymentType,
		Tiers:              convertTiers(data.Tiers),
		CBEGLEntry:         convertToServiceGLEntry(data.CBEGLEntry),
		CBEIFBGLEntry:      convertToServiceGLEntry(data.CBEIFBGLEntry),
		Enabled:            data.Enabled,
		IsDeleted:          data.IsDeleted,
		CreatedAt:          data.CreatedAt,
		LastModifiedAt:     data.LastModifiedAt,
		DeletedAt:          data.DeletedAt,
	}, nil
}

func (o *outboundStore) UpdateOneServiceDetail(ctx context.Context, id string, update serviceDomain.Service) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := map[string]interface{}{"_id": objID}
	updateDoc := map[string]interface{}{

		"serviceCode":        update.ServiceCode,
		"serviceName":        update.ServiceName,
		"serviceType":        update.ServiceType,
		"key":                update.Key,
		"cap":                convertToModelCap(update.Cap),
		"cbeProductCodes":    convertToModelProductCodes(update.CBEProductCodes),
		"cbeIfbProductCodes": convertToModelProductCodes(update.CBEIFBProductCodes),
		"aboveAmount":        update.AboveAmount,
		"aboveServiceFee":    update.AboveServiceFee,
		"paymentType":        update.PaymentType,
		"tiers":              convertToModelTiers(update.Tiers),
		"cbeGLEntry":         convertToModelGLEntry(update.CBEGLEntry),
		"cbeIfbGLEntry":      convertToModelGLEntry(update.CBEIFBGLEntry),
		"enabled":            update.Enabled,
		"isDeleted":          update.IsDeleted,
		"lastModifiedAt":     update.LastModifiedAt,
		"deletedAt":          update.DeletedAt,
	}

	_, err = o.MongoDalServiceDetails.UpdateOne(ctx, filter, updateDoc)
	return err
}

func convertToServiceCap(c model.Cap) serviceDomain.Cap {
	return serviceDomain.Cap{
		KYCLevel:  serviceDomain.KYCLevel(c.KYCLevel),
		SingleCap: c.SingleCap,
		DailyCap:  c.DailyCap,
		MinAmount: c.MinAmount,
	}
}

func convertToModelCap(c serviceDomain.Cap) model.Cap {
	return model.Cap{
		KYCLevel:  model.KYCLevel(c.KYCLevel),
		SingleCap: c.SingleCap,
		DailyCap:  c.DailyCap,
		MinAmount: c.MinAmount,
	}
}

func convertToServiceProductCodes(p model.ProductCodes) serviceDomain.ProductCodes {
	return serviceDomain.ProductCodes{
		PRD:    p.PRD,
		VATPRD: p.VATPRD,
		SFPRD:  p.SFPRD,
		TRXN:   p.TRXN,
	}
}

func convertToModelProductCodes(p serviceDomain.ProductCodes) model.ProductCodes {
	return model.ProductCodes{
		PRD:    p.PRD,
		VATPRD: p.VATPRD,
		SFPRD:  p.SFPRD,
		TRXN:   p.TRXN,
	}
}

func convertTiers(tiers []model.Tier) []serviceDomain.Tier {
	result := make([]serviceDomain.Tier, len(tiers))
	for i, t := range tiers {
		result[i] = serviceDomain.Tier{
			ID:        t.ID.Hex(),
			Min:       t.Min,
			Max:       t.Max,
			FeeAmount: t.FeeAmount,
		}
	}
	return result
}

func convertToModelTiers(tiers []serviceDomain.Tier) []model.Tier {
	result := make([]model.Tier, len(tiers))
	for i, t := range tiers {
		objID, _ := bson.ObjectIDFromHex(t.ID)
		result[i] = model.Tier{
			ID:        objID,
			Min:       t.Min,
			Max:       t.Max,
			FeeAmount: t.FeeAmount,
		}
	}
	return result
}

func convertToServiceGLEntry(g model.GLEntry) serviceDomain.GLEntry {
	return serviceDomain.GLEntry{
		ProductAccount:    g.ProductAccount,
		ProductBranchCode: g.ProductBranchCode,
		ServiceAccount:    g.ServiceAccount,
		ServiceBranchCode: g.ServiceBranchCode,
		VatAccount:        g.VatAccount,
		VatBranchCode:     g.VatBranchCode,
	}
}

func convertToModelGLEntry(g serviceDomain.GLEntry) model.GLEntry {
	return model.GLEntry{
		ProductAccount:    g.ProductAccount,
		ProductBranchCode: g.ProductBranchCode,
		ServiceAccount:    g.ServiceAccount,
		ServiceBranchCode: g.ServiceBranchCode,
		VatAccount:        g.VatAccount,
		VatBranchCode:     g.VatBranchCode,
	}
}

func (o *outboundStore) InitiateServiceFeeUpdate(ctx context.Context, req service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"maker_id":      req.MakerID,
		"action_status": "PENDING",
		"department":    req.Department,
	}
	projection := bson.M{}
	FindAction, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		err = fmt.Errorf("failed to get cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	if FindAction != nil {
		err = fmt.Errorf("pending cps action present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "pending cps action present",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}

	serviceFilter := bson.M{
		"action_code": req.ActionCode,
	}
	serviceProjection := bson.M{}
	serviceFee, err := o.MongoDalServices.FindOne(ctx, serviceFilter, serviceProjection)
	if err != nil {
		err = fmt.Errorf("failed to get service %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	if serviceFee == nil {
		err = fmt.Errorf("service not found: %w", constant.ErrorDefinition{
			Code:    http.StatusNotFound,
			Message: "service not found",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}

	cpsAction := model.CPSAction{
		ActionCode:         req.ActionCode,
		MakerID:            req.MakerID,
		MakerName:          req.MakerName,
		MakerPhoneNumber:   req.MakerPhoneNumber,
		CheckerID:          req.CheckerID,
		CheckerName:        req.CheckerName,
		CheckerPhoneNumber: req.CheckerPhoneNumber,
		UniqueId:           req.UniqueId, // feature/entity ID
		Department:         req.Department,
		ActionStatus:       "PENDING",
		ActionType:         "UPDATE_SERVICE_FEE",
		RequestAction:      "UPDATE",
		PreviousAction:     req.PreviousAction,
		CurrentAction:      req.CurrentAction,
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
	}
	cpsActionResult, err := o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{
		ActionID: cpsActionResult.ActionCode,
	}, nil
}

func (o *outboundStore) ApproveServiceFeeUpdate(ctx context.Context, action_code string) error {
	filter := bson.M{
		"action_code":   action_code,
		"action_status": "PENDING",
	}

	update := bson.M{
		"action_status": "APPROVED",
	}

	projection := bson.M{}
	cpsAction, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil {
		err = fmt.Errorf("failed to find cps action")
		return err
	}

	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("failed to find cps action")
		return err
	}

	switch cpsAction.RequestAction {
	case "CREATE_SERVICE_FEE":
		// Convert CurrentAction to model.Service
		serviceDataMap := cpsAction.CurrentAction

		svc := make(map[string]interface{})
		b, err := json.Marshal(serviceDataMap)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, &svc); err != nil {
			return fmt.Errorf("failed to unmarshal service data")
		}
		serviceData := o.serviceMapper(svc)

		fmt.Println("Service Data", serviceData)

		_, err = o.MongoDalServices.InsertOne(ctx, serviceData)

		if err != nil {
			return fmt.Errorf("failed to create service")
		}
	case "UPDATE_SERVICE_FEE":
		serviceData, ok := cpsAction.CurrentAction.(map[string]interface{})
		if !ok {
			return errors.New("invalid service data for update")
		}
		serviceCode, ok := serviceData["service_code"].(string)
		if !ok {
			return errors.New("invalid service_code in current action")
		}
		serviceFilter := bson.M{"service_code": serviceCode}
		serviceUpdate := bson.M{"$set": serviceData}
		_, err := o.MongoDalServices.UpdateOne(ctx, serviceFilter, serviceUpdate)
		if err != nil {
			return fmt.Errorf("failed to update service fee: %w", err)
		}
	case "DELETE_SERVICE_FEE":
		serviceData, ok := cpsAction.CurrentAction.(map[string]interface{})
		if !ok {
			return errors.New("invalid service data for delete")
		}
		serviceCode, ok := serviceData["service_code"].(string)
		if !ok {
			return errors.New("invalid service_code in current action for delete")
		}
		serviceFilter := bson.M{"service_code": serviceCode}
		serviceUpdate := bson.M{"is_deleted": true, "deleted_at": time.Now(), "last_modified_at": time.Now()}
		_, err := o.MongoDalServices.UpdateOne(ctx, serviceFilter, serviceUpdate)
		if err != nil {
			return fmt.Errorf("failed to delete service: %w", err)
		}
	default:
		return errors.New("unsupported request action")
	}
	return nil
}

func (o *outboundStore) serviceMapper(data map[string]interface{}) model.Service {
	// Helper for safe string from map
	safeStringFromMap := func(m map[string]interface{}, key string) string {
		if v, ok := m[key].(string); ok {
			return v
		}
		return ""
	}
	// Helper for safe uint64 from map
	safeUint64FromMap := func(m map[string]interface{}, key string) uint64 {
		if v, ok := m[key].(float64); ok {
			return uint64(v)
		}
		if v, ok := m[key].(uint64); ok {
			return v
		}
		return 0
	}
	// Helper for Cap
	parseCap := func(val interface{}) model.Cap {
		capMap, ok := val.(map[string]interface{})
		if !ok {
			return model.Cap{}
		}
		return model.Cap{
			KYCLevel:  model.KYCLevel(safeStringFromMap(capMap, "kyc_level")),
			SingleCap: safeUint64FromMap(capMap, "single_cap"),
			DailyCap:  safeUint64FromMap(capMap, "daily_cap"),
			MinAmount: safeUint64FromMap(capMap, "min_amount"),
			MaxAmount: safeUint64FromMap(capMap, "max_amount"),
		}
	}
	// Helper for ProductCodes
	parseProductCodes := func(val interface{}) model.ProductCodes {
		pcMap, ok := val.(map[string]interface{})
		if !ok {
			return model.ProductCodes{}
		}
		return model.ProductCodes{
			PRD:    safeStringFromMap(pcMap, "prd"),
			VATPRD: safeStringFromMap(pcMap, "vatprd"),
			SFPRD:  safeStringFromMap(pcMap, "sfprd"),
			TRXN:   safeStringFromMap(pcMap, "trxn"),
		}
	}
	// Helper for GLEntry
	parseGLEntry := func(val interface{}) model.GLEntry {
		glMap, ok := val.(map[string]interface{})
		if !ok {
			return model.GLEntry{}
		}
		return model.GLEntry{
			ProductAccount:    safeStringFromMap(glMap, "product_account"),
			ProductBranchCode: safeStringFromMap(glMap, "product_branch_code"),
			ServiceAccount:    safeStringFromMap(glMap, "service_account"),
			ServiceBranchCode: safeStringFromMap(glMap, "service_branch_code"),
			VatAccount:        safeStringFromMap(glMap, "vat_account"),
			VatBranchCode:     safeStringFromMap(glMap, "vat_branch_code"),
		}
	}
	// Helper for Tiers
	parseTiers := func(val interface{}) []model.Tier {
		arr, ok := val.([]interface{})
		if !ok {
			return nil
		}
		var tiers []model.Tier
		for _, t := range arr {
			tierMap, ok := t.(map[string]interface{})
			if !ok {
				continue
			}
			tiers = append(tiers, model.Tier{
				// ID:        safeObjectIDFromMap(tierMap, "id"), // omit ID
				Min:       safeUint64FromMap(tierMap, "min"),
				Max:       safeUint64FromMap(tierMap, "max"),
				FeeAmount: safeUint64FromMap(tierMap, "fee_amount"),
			})
		}
		return tiers
	}
	// Helper for safe string
	safeString := func(key string) string {
		if v, ok := data[key].(string); ok {
			return v
		}
		return ""
	}
	// Helper for safe uint64
	safeUint64 := func(key string) uint64 {
		if v, ok := data[key].(float64); ok {
			return uint64(v)
		}
		if v, ok := data[key].(uint64); ok {
			return v
		}
		return 0
	}
	// Helper for safe bool
	safeBool := func(key string) bool {
		if v, ok := data[key].(bool); ok {
			return v
		}
		return false
	}
	// Helper for safe time.Time
	safeTime := func(key string) time.Time {
		if v, ok := data[key].(string); ok {
			t, _ := time.Parse(time.RFC3339, v)
			return t
		}
		return time.Time{}
	}

	return model.Service{
		ServiceCode:        safeString("service_code"),
		ServiceName:        safeString("service_name"),
		ServiceType:        safeString("service_type"),
		Key:                safeString("key"),
		Cap:                parseCap(data["cap"]),
		CBEProductCodes:    parseProductCodes(data["cbe_product_codes"]),
		CBEIFBProductCodes: parseProductCodes(data["cbe_ifb_product_codes"]),
		AboveAmount:        safeUint64("above_amount"),
		AboveServiceFee:    safeUint64("above_service_fee"),
		PaymentType:        safeString("payment_type"),
		Tiers:              parseTiers(data["tiers"]),
		CBEGLEntry:         parseGLEntry(data["cbe_gl_entry"]),
		CBEIFBGLEntry:      parseGLEntry(data["cbe_ifb_gl_entry"]),
		Enabled:            safeBool("enabled"),
		IsDeleted:          safeBool("is_deleted"),
		CreatedAt:          safeTime("created_at"),
		LastModifiedAt:     safeTime("last_modified_at"),
		DeletedAt:          safeTime("deleted_at"),
	}
}

func (o *outboundStore) jsonParser(data any) (map[string]interface{}, error) {

	dataResult := make(map[string]interface{})
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteData, &dataResult); err != nil {
		return nil, err
	}

	return dataResult, nil
}

func (o *outboundStore) RejectServiceFeeUpdate(ctx context.Context, action_code string, rejection_reason string) error {
	filter := bson.M{
		"action_code":   action_code,
		"action_status": "PENDING",
	}
	cpsAction := contexts.ExtractContext(ctx)
	update := bson.M{

		"checker_id":           cpsAction.UserID,
		"checker_name":         cpsAction.FullName,
		"checker_phone_number": cpsAction.PhoneNumber,
		"action_status":        "REJECTED",
		"rejection_reason":     rejection_reason,
		"checker_action_time":  time.Now(),
	}
	_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("failed to reject service fee update %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	}
	return nil
}

func (o *outboundStore) UpdateOneServiceDetailRequest(ctx context.Context, id string, update serviceDomain.Service) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := map[string]interface{}{"_id": objID}
	updateDoc := map[string]interface{}{

		"serviceCode":        update.ServiceCode,
		"serviceName":        update.ServiceName,
		"serviceType":        update.ServiceType,
		"key":                update.Key,
		"cap":                convertToModelCap(update.Cap),
		"cbeProductCodes":    convertToModelProductCodes(update.CBEProductCodes),
		"cbeIfbProductCodes": convertToModelProductCodes(update.CBEIFBProductCodes),
		"aboveAmount":        update.AboveAmount,
		"aboveServiceFee":    update.AboveServiceFee,
		"paymentType":        update.PaymentType,
		"tiers":              convertToModelTiers(update.Tiers),
		"cbeGLEntry":         convertToModelGLEntry(update.CBEGLEntry),
		"cbeIfbGLEntry":      convertToModelGLEntry(update.CBEIFBGLEntry),
		"enabled":            update.Enabled,
		"isDeleted":          update.IsDeleted,
		"lastModifiedAt":     update.LastModifiedAt,
		"deletedAt":          update.DeletedAt,
	}

	_, err = o.MongoDalServiceDetails.UpdateOne(ctx, filter, updateDoc)
	return err
}
func (o *outboundStore) UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := map[string]interface{}{"_id": objID}
	update := map[string]interface{}{

		"cap.min_amount": minAmount,
		"lastModifiedAt": time.Now(),
	}
	_, err = o.MongoDalServiceDetails.UpdateOne(ctx, filter, update)
	return err
}

func (o *outboundStore) ApproveServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error {
	action, err := o.FetchCpsActionById(ctx, actionID)
	if err != nil {
		return err
	}

	if action.ActionStatus != domain.ActionPending {
		return errors.New("action is not in pending status")
	}
	now := time.Now()
	action.CheckerID = checkerID
	action.CheckerActionTime = &now
	action.LastModifiedAt = time.Now()

	if approve {
		action.ActionStatus = domain.ActionApproved
		var serviceData serviceDomain.Service
		if b, ok := action.CurrentAction.([]byte); ok {
			if err := json.Unmarshal(b, &serviceData); err != nil {
				return errors.New("invalid service data in action")
			}
		} else {
			return errors.New("invalid service data in action")
		}
		if err := o.UpdateOneServiceDetail(ctx, action.UniqueId, serviceData); err != nil {
			return err
		}
	} else {
		action.ActionStatus = domain.ActionRejected
		if rejectionReason != "" {
			action.RejectionReason = &rejectionReason
		} else {
			defaultReason := "Rejected by checker"
			action.RejectionReason = &defaultReason
		}
	}

	return o.UpdateCpsAction(ctx, action)
}

func (o *outboundStore) GetAllPasswordRules(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*action.PasswordRule], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{}

	// Add search functionality
	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"password_id": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	// Add filter functionality (if needed, e.g., by status)
	if filterParams.Filters != "" {
		filter["name"] = filterParams.Filters
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	data, err := o.MongoDalPasswordRule.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}
	totalDocs, err := o.MongoDalPasswordRule.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	var dataList []*action.PasswordRule
	for _, v := range data {
		dataList = append(dataList, &action.PasswordRule{
			ID:             v.ID.Hex(),
			MinLength:      v.MinLength,
			MaxLength:      v.MaxLength,
			CapitalLetters: v.CapitalLetters,
			SmallLetters:   v.SmallLetters,
			Numbers:        v.Numbers,
			Name:           v.Name,
			Characters:     v.Characters,
			CreatedAt:      v.CreatedAt,
		})
	}

	meta := common_util.BuildPaginationMeta(totalDocs, page, limit)

	return &common_util.PaginatedResponse[[]*action.PasswordRule]{
		Data: dataList,
		Meta: meta,
	}, nil
}

func (o *outboundStore) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
	filter := map[string]interface{}{"action_code": actionID}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	var rejectionReason *string
	if data.RejectionReason != "" {
		rejectionReason = &data.RejectionReason
	}
	result := &action.CPSAction{

		ID:                 data.ID.Hex(),
		ActionCode:         data.ActionCode,
		MakerID:            data.MakerID,
		MakerName:          data.MakerName,
		MakerPhoneNumber:   data.MakerPhoneNumber,
		CheckerID:          data.CheckerID,
		CheckerName:        data.CheckerName,
		CheckerPhoneNumber: data.CheckerPhoneNumber,
		Department:         data.Department,
		RejectionReason:    rejectionReason,
		PreviousAction:     data.PreviousAction,
		CurrentAction:      data.CurrentAction,
		ActionStatus:       action.ActionStatus(data.ActionStatus),
		ActionType:         action.ActionType(data.ActionType),
		RequestAction:      action.RequestAction(data.RequestAction),
		CreatedAt:          data.CreatedAt,
		LastModifiedAt:     data.LastModifiedAt,
	}
	return result, nil
}
func (o *outboundStore) GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error) {
	department, _ := ctx.Value("department").(string)
	filter := bson.M{
		"department":         department,
		"maker_id":           maker.UserID,
		"maker_name":         maker.FullName,
		"maker_phone_number": maker.PhoneNumber,
		"action_status":      "PENDING",
		"action_type":        "UPDATE",
	}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	var rejectionReason *string
	if data.RejectionReason != "" {
		rejectionReason = &data.RejectionReason
	}
	result := &action.CPSAction{

		ID:                 data.ID.Hex(),
		ActionCode:         data.ActionCode,
		MakerID:            data.MakerID,
		MakerName:          data.MakerName,
		MakerPhoneNumber:   data.MakerPhoneNumber,
		CheckerID:          data.CheckerID,
		CheckerName:        data.CheckerName,
		CheckerPhoneNumber: data.CheckerPhoneNumber,
		Department:         data.Department,
		RejectionReason:    rejectionReason,
		PreviousAction:     data.PreviousAction,
		CurrentAction:      data.CurrentAction,
		ActionStatus:       action.ActionStatus(data.ActionStatus),
		ActionType:         action.ActionType(data.ActionType),
		RequestAction:      action.RequestAction(data.RequestAction),
		CreatedAt:          data.CreatedAt,
		LastModifiedAt:     data.LastModifiedAt,
	}
	return result, nil
}

func (o *outboundStore) GetCurrentPasswordRule(ctx context.Context) (*action.PasswordRule, error) {
	ruleModel, err := o.MongoDalPasswordRule.FindOne(ctx, bson.M{}, bson.M{})

	if err != nil || ruleModel == nil {

		return nil, err
	}

	rule := &action.PasswordRule{
		ID:             ruleModel.ID.Hex(),
		Name:           ruleModel.Name,
		MinLength:      ruleModel.MinLength,
		MaxLength:      ruleModel.MaxLength,
		Numbers:        ruleModel.Numbers,
		CapitalLetters: ruleModel.CapitalLetters,
		SmallLetters:   ruleModel.SmallLetters,
		Characters:     ruleModel.Characters,
		CreatedAt:      ruleModel.CreatedAt,
	}
	return rule, nil
}

// Add this method to outboundStore:
func (o *outboundStore) UpdatePasswordRule(ctx context.Context, rule action.PasswordRule) error {
	if strings.TrimSpace(rule.ID) == "" {
		return errors.New("password rule ID is required")
	}
	objID, err := bson.ObjectIDFromHex(rule.ID)
	if err != nil {
		return errors.New("invalid password rule ID format")
	}
	filter := map[string]interface{}{"_id": objID}
	update := map[string]interface{}{

		"name":            rule.Name,
		"min_length":      rule.MinLength,
		"max_length":      rule.MaxLength,
		"numbers":         rule.Numbers,
		"capital_letters": rule.CapitalLetters,
		"small_letters":   rule.SmallLetters,
		"characters":      rule.Characters,
		"created_at":      rule.CreatedAt,
	}

	_, err = o.MongoDalPasswordRule.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update password rule: %w", err)
	}
	return nil
}

func (o *outboundStore) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
	if strings.TrimSpace(uniqueID) == "" {
		return nil, fmt.Errorf("invalid unique ID")
	}
	filter := map[string]interface{}{
		"unique_id":     uniqueID,
		"action_status": "PENDING",
	}
	projection := map[string]interface{}{
		"action_code": 1,
		"_id":         1,
	}

	doc, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, err
	}

	response := action.ActionResponse{
		ID:       doc.ID.Hex(),
		ActionId: doc.ActionCode,
	}
	return []action.ActionResponse{response}, nil
}
