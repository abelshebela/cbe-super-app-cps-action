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

	portalCardDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	bpscalls "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/bps_calls"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	serviceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	passwordRuleOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	userOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_services"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type outboundStore struct {
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
	mongoDalCPSUser := infra_mongo.NewMongoDal[model.CPSUser, model.CPSUser](client, dbName, collectionNames[0])
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[1])
	mongoDalBPSUser := infra_mongo.NewMongoDal[bps.BPSUser, bps.BPSUser](client, dbName, collectionNames[2])
	mongoDalAccountValidation := infra_mongo.NewMongoDal[model.ValidationRule, model.ValidationRule](client, dbName, collectionNames[3])
	mongoDalServiceDetails := infra_mongo.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, collectionNames[4])
	mongoDalPortalCard := infra_mongo.NewMongoDal[model.Card, model.Card](client, dbName, collectionNames[5])

	return &outboundStore{
		MongoDalCPSUser:           mongoDalCPSUser,
		MongoDalCPSAction:         mongoDalCPSAction,
		MongoDalBPSUser:           mongoDalBPSUser,
		MongoDalAccountValidation: mongoDalAccountValidation,
		MongoDalServiceDetails:    mongoDalServiceDetails,
		MongoDalPortalCard:        mongoDalPortalCard,
	}
}
func NewOutBoundStore(client *mongo.Client, dbName string, collectionNames []string, logger utils.Logger) outbound.OutboundInfra {

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
		BpsCalls:               bpscalls.NewBpsCalls(),
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

func NewServiceDetailsPersistence(client *mongo.Client, dbName string, logger utils.Logger) *outboundStore {
	mongoDalServiceDetails := infra_mongo.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, "cps_services")
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "cps_actions")

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
	filter := map[string]interface{}{"_id": objID}
	data, err := o.MongoDalServiceDetails.FindOne(ctx, filter, nil)
	if err != nil {
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
	filter := map[string]interface{}{
		"_id": service.ID,
	}
	update := map[string]interface{}{

		"key":         service.Key,
		"serviceName": service.ServiceName,
		"serviceType": service.ServiceType,
		"cap": map[string]interface{}{
			"kyc_level":  service.Cap.KYCLevel,
			"single_cap": service.Cap.SingleCap,
			"daily_cap":  service.Cap.DailyCap,
			"min_amount": service.Cap.MinAmount,
		},
		"cbe_product_codes": map[string]interface{}{
			"prd":    service.CBEProductCodes.PRD,
			"vatprd": service.CBEProductCodes.VATPRD,
			"sfprd":  service.CBEProductCodes.SFPRD,
			"trxn":   service.CBEProductCodes.TRXN,
		},
		"cbe_ifb_product_codes": map[string]interface{}{
			"prd":    service.CBEIFBProductCodes.PRD,
			"vatprd": service.CBEIFBProductCodes.VATPRD,
			"sfprd":  service.CBEIFBProductCodes.SFPRD,
			"trxn":   service.CBEIFBProductCodes.TRXN,
		},
		"above_amount":      service.AboveAmount,
		"above_service_fee": service.AboveServiceFee,
		"payment_type":      service.PaymentType,
		"tiers": func() []map[string]interface{} {
			var tiers []map[string]interface{}
			for _, tier := range service.Tiers {
				tiers = append(tiers, map[string]interface{}{
					"id":         tier.ID,
					"min":        tier.Min,
					"max":        tier.Max,
					"fee_amount": tier.FeeAmount,
				})
			}
			return tiers
		}(),
		"cbe_gl_entry": map[string]interface{}{
			"product_account":     service.CBEGLEntry.ProductAccount,
			"product_branch_code": service.CBEGLEntry.ProductBranchCode,
			"service_account":     service.CBEGLEntry.ServiceAccount,
			"service_branch_code": service.CBEGLEntry.ServiceBranchCode,
			"vat_account":         service.CBEGLEntry.VatAccount,
			"vat_branch_code":     service.CBEGLEntry.VatBranchCode,
		},
		"cbe_ifb_gl_entry": map[string]interface{}{
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

	_, err := o.MongoDalServiceDetails.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}
	return nil
}
func stringPointer(s string) *string {
	return &s
}

func domainToModelCPSAction(domainAction domain.CPSAction) model.CPSAction {
	rejectionReason := "rejected by the checker"
	if domainAction.RejectionReason != nil {
		rejectionReason = *domainAction.RejectionReason
	}

	var objID bson.ObjectID
	if domainAction.ID != "" {
		var err error
		objID, err = bson.ObjectIDFromHex(domainAction.ID)
		if err != nil {
			objID = bson.NilObjectID
		}
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
		RejectionReason:    rejectionReason,
		PreviosAction:      domainAction.PreviosAction,
		CurrentAction:      currentAction,
		ActionStatus:       string(domainAction.ActionStatus),
		ActionType:         string(domainAction.ActionType),
		RequestAction:      string(domainAction.RequestAction),
		CreatedAt:          domainAction.CreatedAt,
		LastModifiedAt:     domainAction.LastModifiedAt,
		MakerActionTime:    domainAction.MakerActionTime,
		CheckerActionTime:  domainAction.CheckerActionTime,
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
		PreviosAction:      modelAction.PreviosAction,
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
	filter := map[string]interface{}{
		"action_code": modelAction.ActionCode,
	}
	update := map[string]interface{}{
		"maker_id":             modelAction.MakerID,
		"maker_name":           modelAction.MakerName,
		"maker_phone_number":   modelAction.MakerPhoneNumber,
		"checker_id":           modelAction.CheckerID,
		"checker_name":         modelAction.CheckerName,
		"checker_phone_number": modelAction.CheckerPhoneNumber,
		"unique_id":            modelAction.UniqueId,
		"department":           modelAction.Department,
		"rejection_reason":     modelAction.RejectionReason,
		"previos_action":       modelAction.PreviosAction,
		"current_action":       modelAction.CurrentAction,
		"action_status":        modelAction.ActionStatus,
		"action_type":          modelAction.ActionType,
		"request_action":       modelAction.RequestAction,
		"created_at":           modelAction.CreatedAt,
		"last_modified_at":     modelAction.LastModifiedAt,
	}
	if _, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf(error_codes.GeneralDBUpdateFailed)

	}

	return nil
}
func (o *outboundStore) FetchCpsActionById(ctx context.Context, Action_Id string) (domain.CPSAction, error) {
	filter := map[string]interface{}{"action_code": Action_Id}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return domain.CPSAction{}, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}

	return modelToDomainCPSAction(*data), nil
}

func (o *outboundStore) FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]domain.LinkedAccount, error) {
	v, err := o.BpsCalls.FetchLinkedAccount(accountNumber)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		updatedAccounts = append(updatedAccounts, acc)
	}
	return updatedAccounts, nil
}

func (o *outboundStore) UpdateAccount(ctx context.Context, linkedAccount domain.LinkedAccount) (domain.LinkedAccount, error) {
	filter := map[string]interface{}{"_id": linkedAccount.ID}
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

	_, err := o.MongoDalAccounts.UpdateOne(ctx, filter, update)
	if err != nil {
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

func (o *outboundStore) FetchLinkedAccountById(ctx context.Context, id []string) ([]domain.LinkedAccount, error) {

	var result []domain.LinkedAccount
	for _, i := range id {
		objID, err := bson.ObjectIDFromHex(i)
		if err != nil {
			return nil, err
		}
		filter := map[string]interface{}{"_id": objID}
		item, err := o.MongoDalAccounts.FindOne(ctx, filter, nil)
		if err != nil {
			return nil, fmt.Errorf(error_codes.GeneralDBQueryFailed)
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
func (o *outboundStore) CreateUserRequest(ctx context.Context, user domain.CPSUser, maker domain.User) error {
	department, _ := ctx.Value("department").(string)
	actionCode := utils.RandomGenerator(24)

	prevActionJSON := json.RawMessage("null")
	currActionJSON, _ := json.Marshal(user)

	cpsAction := model.CPSAction{
		ActionCode:       actionCode,
		MakerID:          maker.UserID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		UniqueId:         user.UserCode,
		Department:       department,
		ActionStatus:     "PENDING",
		ActionType:       "CREATE",
		RequestAction:    "USER",
		PreviosAction:    prevActionJSON,
		CurrentAction:    currActionJSON,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	filter := bson.M{
		"unique_id":      user.UserCode,
		"action_status":  "PENDING",
		"action_type":    "CREATE",
		"request_action": "USER",
	}
	projection := bson.M{"_id": 1}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err == nil && existing != nil {
		return errors.New("a pending user creation action already exists for this user code")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || (err != nil && strings.Contains(err.Error(), "E11000")) {
			return errors.New("duplicate action code or user creation request")
		}
		return err
	}
	return nil
}
func (o *outboundStore) UpdateUserRequest(ctx context.Context, updated domain.CPSUser, maker domain.User) error {

	if strings.TrimSpace(updated.UserCode) == "" {
		return errors.New("user_code is required")
	}
	if strings.TrimSpace(updated.UserCode) == "" {
		return errors.New("user_code is required")
	}

	filter := bson.M{"user_code": updated.UserCode}
	modelUser, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil {
		return errors.New("user not found")
	}
	// filter := bson.M{"user_code": updated.UserCode}
	// modelUser, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	// if err != nil {
	// 	return errors.New("user not found")
	// }

	// var department string
	// if len(modelUser.Department) > 0 {
	// 	department = modelUser.Department.Hex()
	// }
	var permissionCategory []string
	for _, oid := range modelUser.PermissionCategory {
		permissionCategory = append(permissionCategory, oid.Hex())
	}
	// var permissionGroup []string
	// for _, oid := range modelUser.PermissionGroup {
	// 	permissionGroup = append(permissionGroup, oid.Hex())
	// }
	var department string
	if len(modelUser.Department) > 0 {
		department = modelUser.Department.Hex()
	}
	// var permissionCategory []string
	// for _, oid := range modelUser.PermissionCategory {
	// 	permissionCategory = append(permissionCategory, oid.Hex())
	// }
	var permissionGroup []string
	for _, oid := range modelUser.PermissionGroup {
		permissionGroup = append(permissionGroup, oid.Hex())
	}

	prevUser := domain.CPSUser{
		UserCode:           modelUser.UserCode,
		UserName:           modelUser.UserName,
		FullName:           modelUser.FullName,
		PhoneNumber:        modelUser.PhoneNumber,
		Role:               modelUser.Role,
		Department:         department,
		PermissionCategory: permissionCategory,
		PermissionGroup:    permissionGroup,
	}
	// prevUser := domain.CPSUser{
	// 	UserCode:           modelUser.UserCode,
	// 	UserName:           modelUser.UserName,
	// 	FullName:           modelUser.FullName,
	// 	PhoneNumber:        modelUser.PhoneNumber,
	// 	Role:               modelUser.Role,
	// 	Department:         department,
	// 	PermissionCategory: permissionCategory,
	// 	PermissionGroup:    permissionGroup,
	// }

	prevActionJSON, _ := json.Marshal(prevUser)
	currActionJSON, _ := json.Marshal(updated)
	// prevActionJSON, _ := json.Marshal(prevUser)
	// currActionJSON, _ := json.Marshal(updated)

	updateDoc := bson.M{

		"user_code":           updated.UserCode,
		"user_name":           updated.UserName,
		"full_name":           updated.FullName,
		"phone_number":        updated.PhoneNumber,
		"role":                updated.Role,
		"department":          updated.Department,
		"permission_category": updated.PermissionCategory,
		"permission_group":    updated.PermissionGroup,
		"last_modified":       time.Now(),
	}
	_, err = o.MongoDalCPSUser.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
		return err
	}
	// updateDoc := bson.M{
	// 	"$set": bson.M{
	// 		"user_code":           updated.UserCode,
	// 		"user_name":           updated.UserName,
	// 		"full_name":           updated.FullName,
	// 		"phone_number":        updated.PhoneNumber,
	// 		"role":                updated.Role,
	// 		"department":          updated.Department,
	// 		"permission_category": updated.PermissionCategory,
	// 		"permission_group":    updated.PermissionGroup,
	// 		"last_modified":       time.Now(),
	// 	},
	// }
	// _, err = o.MongoDalCPSUser.UpdateOne(ctx, filter, updateDoc)
	// if err != nil {
	// 	fmt.Printf("Failed to update user: %v\n", err)
	// 	return err
	// }

	actionCode := utils.RandomGenerator(24)
	departmentCtx, _ := ctx.Value("department").(string)
	cpsAction := model.CPSAction{
		ActionCode:       actionCode,
		MakerID:          maker.UserID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		UniqueId:         updated.UserCode,
		Department:       departmentCtx,
		ActionStatus:     "PENDING",
		ActionType:       "UPDATE",
		RequestAction:    "USER",
		PreviosAction:    prevActionJSON,
		CurrentAction:    currActionJSON,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		fmt.Printf("Failed to log update action: %v\n", err)
		return err
	}
	// actionCode := utils.RandomGenerator(24)
	// departmentCtx, _ := ctx.Value("department").(string)
	// cpsAction := model.CPSAction{
	// 	ActionCode: actionCode,
	// 	MakerUser: model.User{
	// 		UserCode:    maker.UserID,
	// 		FullName:    maker.FullName,
	// 		PhoneNumber: maker.PhoneNumber,
	// 	},
	// 	UniqueId:      updated.UserCode,
	// 	Department:     departmentCtx,
	// 	ActionStatus:   model.ActionPending,
	// 	ActionType:     model.ActionUpdate,
	// 	RequestAction:  model.RequestUser,
	// 	PreviosAction:  prevActionJSON,
	// 	CurrentAction:  currActionJSON,
	// 	CreatedAt:      time.Now(),
	// 	LastModifiedAt: time.Now(),
	// }
	// _, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	// if err != nil {
	// 	fmt.Printf("Failed to log update action: %v\n", err)
	// 	return err
	// }

	return nil
	// return nil
}
func (o *outboundStore) ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return errors.New("actionID is required")
	}
	filter := bson.M{"action_code": actionID}
	update := bson.M{
		"last_modified_at": time.Now(),
	}
	if approve {
		update["action_status"] = "APPROVED"
	} else {
		update["action_status"] = "REJECTED"
	}
	if reason != nil {
		update["rejection_reason"] = *reason
	}
	_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return err
}
func (o *outboundStore) GetPendingUserActions(ctx context.Context, actionCode string) ([]domain.CPSAction, error) {
	department, _ := ctx.Value("department").(string)
	filter := bson.M{
		"action_status": "PENDING",
		"department":    department,
	}
	if actionCode != "" {
		filter["action_code"] = actionCode
	}
	fmt.Printf("Fetching pending actions with filter: %+v\n", filter)
	actionPtrs, err := o.MongoDalCPSAction.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	var actions []domain.CPSAction
	for _, ptr := range actionPtrs {
		if ptr != nil {
			actions = append(actions, modelToDomainCPSAction(*ptr))
		}
	}
	return actions, nil
}

func (o *outboundStore) FetchUserByUserCode(ctx context.Context, userCode string) (*domain.CPSUser, error) {
	filter := bson.M{"user_code": userCode}
	modelUser, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	// filter := bson.M{"user_code": userCode}
	// modelUser, err := o.MongoDalCPSUser.FindOne(ctx, filter, nil)
	// if err != nil {
	// 	return nil, err
	// }

	var department string
	if len(modelUser.Department) > 0 {
		department = modelUser.Department.Hex()
	}

	var permissionCategory []string
	for _, oid := range modelUser.PermissionCategory {
		permissionCategory = append(permissionCategory, oid.Hex())
	}
	var permissionGroup []string
	for _, oid := range modelUser.PermissionGroup {
		permissionGroup = append(permissionGroup, oid.Hex())
	}
	// var permissionCategory []string
	// for _, oid := range modelUser.PermissionCategory {
	// 	permissionCategory = append(permissionCategory, oid.Hex())
	// }
	// var permissionGroup []string
	// for _, oid := range modelUser.PermissionGroup {
	// 	permissionGroup = append(permissionGroup, oid.Hex())
	// }

	user := &domain.CPSUser{
		UserCode:           modelUser.UserCode,
		UserName:           modelUser.UserName,
		FullName:           modelUser.FullName,
		PhoneNumber:        modelUser.PhoneNumber,
		Role:               modelUser.Role,
		Department:         department,
		PermissionCategory: permissionCategory,
		PermissionGroup:    permissionGroup,
	}
	return user, nil
	// user := &domain.CPSUser{
	// 	UserCode:           modelUser.UserCode,
	// 	UserName:           modelUser.UserName,
	// 	FullName:           modelUser.FullName,
	// 	PhoneNumber:        modelUser.PhoneNumber,
	// 	Role:               modelUser.Role,
	// 	Department:         department,
	// 	PermissionCategory: permissionCategory,
	// 	PermissionGroup:    permissionGroup,
	// }
	// return user, nil
}
func ptrTime(t time.Time) *time.Time {
	return &t
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

	filter := map[string]interface{}{"_id": objID}
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
		PreviosAction:      req.PreviosAction,
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

func (o *outboundStore) ApproveServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": "PENDING",
	}
	projection := bson.M{}
	cpsActionPtr, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil {
		err = fmt.Errorf("failed to find cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	serviceFilter := bson.M{
		"action_code": cpsActionPtr.ActionCode,
	}
	serviceUpdate := bson.M{
		"$set": cpsActionPtr.CurrentAction,
	}
	_, err = o.MongoDalServices.UpdateOne(ctx, serviceFilter, serviceUpdate)
	if err != nil {
		err = fmt.Errorf("failed to update service fee %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{
		ActionID: cpsActionPtr.ActionCode,
	}, nil
}

func (o *outboundStore) RejectServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": "PENDING",
	}
	update := bson.M{

		"checker_id":           cpsAction.CheckerID,
		"checker_name":         cpsAction.CheckerName,
		"checker_phone_number": cpsAction.CheckerPhoneNumber,
		"action_status":        "REJECTED",
		"rejection_reason":     cpsAction.RejectionReason,
		"checker_action_time":  time.Now(),
	}
	_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		err = fmt.Errorf("failed to reject service fee update %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{
		ActionID: cpsAction.ActionCode,
	}, nil
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

	action.CheckerID = checkerID
	action.CheckerActionTime = time.Now()
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
		PreviosAction:      data.PreviosAction,
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
		PreviosAction:      data.PreviosAction,
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
		PasswordID:     ruleModel.PasswordId,
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

		"password_id":     rule.PasswordID,
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
		if err.Error() == "mongo: no documents in result" {
			return []action.ActionResponse{}, nil
		}
		return nil, fmt.Errorf("failed to fetch pending action: %w", err)
	}
	if doc == nil {
		return []action.ActionResponse{}, nil
	}
	response := action.ActionResponse{
		ID:       doc.ID.Hex(),
		ActionId: doc.ActionCode,
	}
	return []action.ActionResponse{response}, nil
}
