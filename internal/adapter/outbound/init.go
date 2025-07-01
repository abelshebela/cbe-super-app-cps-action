package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	portalCardDomain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/portal_card"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	bpscalls "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/bps_calls"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
	serviceDomain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
	passwordRuleOutbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound"
	userOutbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound"
	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/bulk_services"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
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
	mongoDalServiceDetails := infra_mongo.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, "CPSServices")
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "CPSActions")

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
		return domain.ServiceDetails{}, err
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
		"$set": map[string]interface{}{
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
		},
	}
	_, err := o.MongoDalServiceDetails.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func stringPointer(s string) *string {
	return &s
}

func (o *outboundStore) CreateCpsAction(ctx context.Context, Action domain.CPSAction) (domain.CPSAction, error) {
	makerUser := model.User{
		UserCode:    Action.Maker.UserID,
		FullName:    Action.Maker.FullName,
		PhoneNumber: Action.Maker.PhoneNumber,
	}
	checkerUser := model.User{
		UserCode:    Action.Checker.UserID,
		FullName:    Action.Checker.FullName,
		PhoneNumber: Action.Checker.PhoneNumber,
	}

	CpsAction := model.CPSAction{
		ActionCode:  Action.ActionCode,
		MakerUser:   makerUser,
		CheckerUser: checkerUser,
		UniqueId: func() string {
			if Action.UniqueId != "" {
				return Action.UniqueId
			}
			return Action.Maker.UserID
		}(),
		CheckerID:          stringPointer(Action.Checker.UserID),
		CheckerName:        stringPointer(Action.Checker.FullName),
		CheckerPhoneNumber: stringPointer(Action.Checker.PhoneNumber),
		Department:         Action.Department,
		RejectionReason:    Action.RejectionReason,
		PreviosAction: func() json.RawMessage {
			b, _ := json.Marshal(Action.PreviosAction)
			return b
		}(),
		CurrentAction: func() json.RawMessage {
			b, _ := json.Marshal(Action.CurrentAction)
			return b
		}(),
		ActionStatus:   model.ActionStatus(Action.ActionStatus),
		ActionType:     model.ActionType(Action.ActionType),
		RequestAction:  model.RequestAction(Action.RequestAction),
		CreatedAt:      Action.CreatedAt,
		LastModifiedAt: Action.LastModifiedAt,
	}

	data, err := o.MongoDalCPSAction.InsertOne(ctx, CpsAction)
	if err != nil {
		return domain.CPSAction{}, err
	}

	result := domain.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Maker:           Action.Maker,
		Checker:         Action.Checker,
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction: func() interface{} {
			var v interface{}
			_ = json.Unmarshal(data.PreviosAction, &v)
			return v
		}(),
		CurrentAction: func() interface{} {
			var v interface{}
			_ = json.Unmarshal(data.CurrentAction, &v)
			return v
		}(),
		ActionStatus:   domain.ActionStatus(data.ActionStatus),
		ActionType:     domain.ActionType(data.ActionType),
		RequestAction:  domain.RequestAction(data.RequestAction),
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
	}
	return result, nil
}
func (o *outboundStore) UpdateCpsAction(ctx context.Context, Action domain.CPSAction) error {
	makerUser := model.User{
		UserCode:    Action.Maker.UserID,
		FullName:    Action.Maker.FullName,
		PhoneNumber: Action.Maker.PhoneNumber,
	}
	checkerUser := model.User{
		UserCode:    Action.Checker.UserID,
		FullName:    Action.Checker.FullName,
		PhoneNumber: Action.Checker.PhoneNumber,
	}
	CpsAction := model.CPSAction{
		ActionCode:         Action.ActionCode,
		MakerUser:          makerUser,
		CheckerUser:        checkerUser,
		UniqueId:           Action.Maker.UserID,
		CheckerID:          stringPointer(Action.Checker.UserID),
		CheckerName:        stringPointer(Action.Checker.FullName),
		CheckerPhoneNumber: stringPointer(Action.Checker.PhoneNumber),
		Department:         Action.Department,
		RejectionReason:    Action.RejectionReason,
		PreviosAction: func() json.RawMessage {
			b, _ := json.Marshal(Action.PreviosAction)
			return b
		}(),
		CurrentAction: func() json.RawMessage {
			b, _ := json.Marshal(Action.CurrentAction)
			return b
		}(),
		ActionStatus:   model.ActionStatus(Action.ActionStatus),
		ActionType:     model.ActionType(Action.ActionType),
		RequestAction:  model.RequestAction(Action.RequestAction),
		CreatedAt:      Action.CreatedAt,
		LastModifiedAt: Action.LastModifiedAt,
	}

	filter := map[string]interface{}{
		"action_code": CpsAction.ActionCode,
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"maker_user":           CpsAction.MakerUser,
			"checker_user":         CpsAction.CheckerUser,
			"unique_id":            CpsAction.UniqueId,
			"checker_id":           CpsAction.CheckerID,
			"checker_name":         CpsAction.CheckerName,
			"checker_phone_number": CpsAction.CheckerPhoneNumber,
			"department":           CpsAction.Department,
			"rejection_reason":     CpsAction.RejectionReason,
			"previos_action":       CpsAction.PreviosAction,
			"current_action":       CpsAction.CurrentAction,
			"action_status":        CpsAction.ActionStatus,
			"action_type":          CpsAction.ActionType,
			"request_action":       CpsAction.RequestAction,
			"created_at":           CpsAction.CreatedAt,
			"last_modified_at":     CpsAction.LastModifiedAt,
		},
	}
	_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return err
}
func (o *outboundStore) FetchCpsActionById(ctx context.Context, Action_Id string) (domain.CPSAction, error) {
	filter := map[string]interface{}{"action_code": Action_Id}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return domain.CPSAction{}, err
	}
	result := domain.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction:   data.PreviosAction,
		CurrentAction: domain.CurrentAction{
			Id: func() []string {
				var currentAction model.CurrentAction
				if err := json.Unmarshal(data.CurrentAction, &currentAction); err != nil {
					return nil
				}
				return currentAction.Id
			}(),
			Action: func() bool {
				var currentAction model.CurrentAction
				if err := json.Unmarshal(data.CurrentAction, &currentAction); err != nil {
					return false
				}
				return currentAction.Action
			}(),
		},
		ActionStatus:   domain.ActionStatus(data.ActionStatus),
		ActionType:     domain.ActionType(data.ActionType),
		RequestAction:  domain.RequestAction(data.RequestAction),
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
	}
	return result, nil
}

func (o *outboundStore) FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]domain.LinkedAccount, error) {
	v, err := o.BpsCalls.FetchLinkedAccount(accountNumber)
	if err != nil {
		return nil, err
	}
	filter := map[string]interface{}{"customer_number": v.CustomerNumber}
	data, err := o.MongoDalAccounts.FindAll(ctx, filter, nil)
	if err != nil {
		return nil, err
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
			"$set": map[string]interface{}{
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
			},
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
		"$set": map[string]interface{}{
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
		},
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
		return domain.CPSAction{}, err
	}
	var data = *d[len(d)-1]
	result := domain.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction:   data.PreviosAction, // Assuming this is directly mapped
		CurrentAction: domain.CurrentAction{
			Id: func() []string {
				var currentAction model.CurrentAction
				if err := json.Unmarshal(data.CurrentAction, &currentAction); err != nil {
					return nil // Handle error or return empty slice
				}
				return currentAction.Id
			}(),
			Action: func() bool {
				var currentAction model.CurrentAction
				if err := json.Unmarshal(data.CurrentAction, &currentAction); err != nil {
					return false // Handle error or return default value
				}
				return currentAction.Action
			}(),
		},
		ActionStatus:   domain.ActionStatus(data.ActionStatus),
		ActionType:     domain.ActionType(data.ActionType),
		RequestAction:  domain.RequestAction(data.RequestAction),
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
	}
	return result, nil
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
			return nil, err
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
		ActionCode: actionCode,
		MakerUser: model.User{
			UserCode:    maker.UserID,
			FullName:    maker.FullName,
			PhoneNumber: maker.PhoneNumber,
		},
		UniqueId:       user.UserCode,
		Department:     department,
		ActionStatus:   model.ActionPending,
		ActionType:     model.ActionCreate,
		RequestAction:  model.RequestUser,
		PreviosAction:  prevActionJSON,
		CurrentAction:  currActionJSON,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	filter := bson.M{
		"unique_id":      user.UserCode,
		"action_status":  model.ActionPending,
		"action_type":    model.ActionCreate,
		"request_action": model.RequestUser,
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
		"$set": bson.M{
			"user_code":           updated.UserCode,
			"user_name":           updated.UserName,
			"full_name":           updated.FullName,
			"phone_number":        updated.PhoneNumber,
			"role":                updated.Role,
			"department":          updated.Department,
			"permission_category": updated.PermissionCategory,
			"permission_group":    updated.PermissionGroup,
			"last_modified":       time.Now(),
		},
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
		ActionCode: actionCode,
		MakerUser: model.User{
			UserCode:    maker.UserID,
			FullName:    maker.FullName,
			PhoneNumber: maker.PhoneNumber,
		},
		UniqueId:       updated.UserCode,
		Department:     departmentCtx,
		ActionStatus:   model.ActionPending,
		ActionType:     model.ActionUpdate,
		RequestAction:  model.RequestUser,
		PreviosAction:  prevActionJSON,
		CurrentAction:  currActionJSON,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
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
		update["action_status"] = model.ActionApproved
	} else {
		update["action_status"] = model.ActionRejected
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
		"action_status": model.ActionPending,
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
			actions = append(actions, domain.CPSAction{
				ID:             ptr.ID.Hex(),
				ActionCode:     ptr.ActionCode,
				Department:     ptr.Department,
				ActionStatus:   domain.ActionStatus(ptr.ActionStatus),
				ActionType:     domain.ActionType(ptr.ActionType),
				RequestAction:  domain.RequestAction(ptr.RequestAction),
				CreatedAt:      ptr.CreatedAt,
				LastModifiedAt: ptr.LastModifiedAt,
			})
		}
	}
	return actions, nil
	// department, _ := ctx.Value("department").(string)
	// filter := bson.M{
	// 	"action_status": model.ActionPending,
	// 	"department":    department,
	// }
	// if actionCode != "" {
	// 	filter["action_code"] = actionCode
	// }
	// fmt.Printf("Fetching pending actions with filter: %+v\n", filter)
	// actionPtrs, err := o.MongoDalCPSAction.FindAll(ctx, filter, bson.M{})
	// if err != nil {
	// 	return nil, err
	// }
	// var actions []domain.CPSAction
	// for _, ptr := range actionPtrs {
	// 	if ptr != nil {
	// 		actions = append(actions, domain.CPSAction{
	// 			ID:             ptr.ID.Hex(),
	// 			ActionCode:     ptr.ActionCode,
	// 			Department:     ptr.Department,
	// 			ActionStatus:   domain.ActionStatus(ptr.ActionStatus),
	// 			ActionType:     domain.ActionType(ptr.ActionType),
	// 			RequestAction:  domain.RequestAction(ptr.RequestAction),
	// 			CreatedAt:      ptr.CreatedAt,
	// 			LastModifiedAt: ptr.LastModifiedAt,
	// 		})
	// 	}
	// }
	// return actions, nil
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
		"$set": map[string]interface{}{
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
		},
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
		"maker_user.phone_number": req.MakerUser.PhoneNumber,
		"status":                  "PENDING",
		"department":              req.Department,
	}

	projection := bson.M{}

	FindAction, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		// o.logger.Errorf("failed to get cps Action", err)
		err = fmt.Errorf("failed to get get cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}

	if FindAction != nil {
		// o.logger.Infof("pending cps action present", req.MakerUser.FullName, req.MakerUser.UserCode, req.Department)
		err = fmt.Errorf("failed to get cps action  %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "pending cps action present",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}

	serviceFilter := bson.M{
		"action_code": req.ActionCode,
	}

	serviceProjection := bson.M{}

	// var serviceFee *service.Service
	serviceFee, err := o.MongoDalServices.FindOne(ctx, serviceFilter, serviceProjection)
	if err != nil {
		// s.logger.Errorf("failed to get service", err)
		err = fmt.Errorf("failed to get service %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}
	if serviceFee == nil {
		// s.logger.Errorf("service not found")
		err = fmt.Errorf("service not found: %w", constant.ErrorDefinition{
			Code:    http.StatusNotFound,
			Message: "service not found",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}

	// Prepare CPSAction for insert (not service.Service)
	cpsAction := model.CPSAction{
		ActionCode: req.ActionCode,
		MakerUser: model.User{
			UserCode:    req.MakerUser.UserCode,
			FullName:    req.MakerUser.FullName,
			PhoneNumber: req.MakerUser.PhoneNumber,
		},
		CheckerUser: model.User{
			UserCode:    req.CheckerUser.UserCode,
			FullName:    req.CheckerUser.FullName,
			PhoneNumber: req.CheckerUser.PhoneNumber,
		},
		UniqueId:      req.MakerUser.UserCode,
		Department:    req.Department,
		ActionStatus:  model.ActionStatus("PENDING"),
		ActionType:    model.ActionType("UPDATE_SERVICE_FEE"),
		RequestAction: model.RequestAction("UPDATE"),
		PreviosAction: func() json.RawMessage {
			b, _ := json.Marshal(map[string]any{
				"tiers": req.PreviousData,
			})
			return b
		}(),
		CurrentAction: func() json.RawMessage {
			b, _ := json.Marshal(map[string]any{
				"is_account_blocked": true,
			})
			return b
		}(),
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	req.MakerActionTime = time.Now()
	cpsActionResult, err := o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		// s.logger.Errorf("failed to create cps action", err)
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
		"action_code": cpsAction.ActionCode,
		"status":      "PENDING",
	}

	projection := bson.M{}
	cpsActionPtr, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil {
		// s.logger.Errorf("failed to find cpsActon", err)
		err = fmt.Errorf("failed to findCpsAction %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return service.UpdateServiceDetailsResponse{}, err
	}

	serviceFilter := bson.M{
		"action_code": cpsActionPtr.ActionCode,
	}
	serviceUpdate := bson.M{
		"set": cpsActionPtr.CurrentAction,
	}
	_, err = o.MongoDalServices.UpdateOne(ctx, serviceFilter, serviceUpdate)
	if err != nil {
		// o.logger.Errorf("failed to update service fee", err)
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
		"action_code": cpsAction.ActionCode,
		"status":      "PENDING",
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    cpsAction.CheckerUser.FullName,
			"phone_number": cpsAction.CheckerUser.PhoneNumber,
			"user_code":    cpsAction.CheckerUser.UserCode,
		},
		"status":                 "REJECTED",
		"rejected_action_reason": cpsAction.RejectedReason,
		"checker_action_time":    time.Now(),
	}

	_, err := o.MongoDalServices.UpdateOne(ctx, filter, update)
	if err != nil {
		// s.logger.Errorf("failed to reject service fee update", err)
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
		"$set": map[string]interface{}{
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
		},
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
		"$set": map[string]interface{}{
			"cap.min_amount": minAmount,
			"lastModifiedAt": time.Now(),
		},
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

	action.Checker.UserID = checkerID
	action.Checker.Timestamp = time.Now()
	action.LastModifiedAt = time.Now()

	if approve {
		action.ActionStatus = domain.ActionApproved

		serviceData, ok := action.CurrentAction.(*serviceDomain.Service)
		if !ok {
			return errors.New("invalid service data in action")
		}
		if err := o.UpdateOneServiceDetail(ctx, action.UniqueId, *serviceData); err != nil {
			return err
		}
	} else {
		action.ActionStatus = domain.ActionRejected
		if rejectionReason != "" {
			action.RejectionReason = &rejectionReason
		} else {
			reason := "Rejected by checker"
			action.RejectionReason = &reason
		}
	}

	return o.UpdateCpsAction(ctx, action)
}
func (o *outboundStore) CreatePasswordRuleUpdateAction(ctx context.Context, rule *action.PasswordRule, maker action.User) (string, error) {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return "", errors.New("department is required in context")
	}

	filter := bson.M{
		"department":           department,
		"maker_user.user_code": maker.UserID,
		"action_status":        model.ActionPending,
		"action_type":          model.ActionUpdate,
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", fmt.Errorf("pending update action already exists for this maker")
	}

	var prevAction json.RawMessage
	prevRulePtr, err := o.MongoDalPasswordRule.FindOne(ctx, bson.M{}, bson.M{})
	if err == nil && prevRulePtr != nil {
		prevAction, _ = json.Marshal(prevRulePtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(rule)
	cpsAction := model.CPSAction{
		ActionCode:     utils.RandomGenerator(24),
		MakerUser:      model.User{UserCode: maker.UserID, FullName: maker.FullName, PhoneNumber: maker.PhoneNumber},
		Department:     department,
		ActionStatus:   model.ActionPending,
		ActionType:     model.ActionUpdate,
		RequestAction:  model.RequestUpdatePasswordRule,
		PreviosAction:  prevAction,
		CurrentAction:  currAction,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return cpsAction.ActionCode, nil
}
func (o *outboundStore) ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error {
	filter := map[string]interface{}{"action_code": actionID}
	action, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	update := map[string]interface{}{
		"checker_user":     model.User{UserCode: checker.UserID, FullName: checker.FullName, PhoneNumber: checker.PhoneNumber},
		"last_modified_at": time.Now(),
	}
	if approve {
		update["action_status"] = model.ActionApproved
		var rule model.PasswordRule
		if err := json.Unmarshal(action.CurrentAction, &rule); err != nil {
			return err
		}
		rule.ID = ""
		_, err := o.MongoDalPasswordRule.InsertOne(ctx, rule)
		if err != nil {
			return err
		}
	} else {
		update["action_status"] = model.ActionRejected
		if rejectionReason != nil {
			update["rejection_reason"] = *rejectionReason
		}
		//
	}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, map[string]interface{}{"$set": update})
	return err
}

func (o *outboundStore) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
	filter := map[string]interface{}{"action_code": actionID}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	result := &action.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Maker:           action.User{UserID: data.MakerUser.UserCode, FullName: data.MakerUser.FullName, PhoneNumber: data.MakerUser.PhoneNumber},
		Checker:         action.User{UserID: data.CheckerUser.UserCode, FullName: data.CheckerUser.FullName, PhoneNumber: data.CheckerUser.PhoneNumber},
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction:   data.PreviosAction,
		CurrentAction:   data.CurrentAction,
		ActionStatus:    action.ActionStatus(data.ActionStatus),
		ActionType:      action.ActionType(data.ActionType),
		RequestAction:   action.RequestAction(data.RequestAction),
		CreatedAt:       data.CreatedAt,
		LastModifiedAt:  data.LastModifiedAt,
	}
	return result, nil
}
func (o *outboundStore) GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error) {
	department, _ := ctx.Value("department").(string)
	filter := bson.M{
		"department":         department,
		"maker.user_id":      maker.UserID,
		"maker.user_code":    maker.UserCode,
		"maker.full_name":    maker.FullName,
		"maker.phone_number": maker.PhoneNumber,
		"action_status":      model.ActionPending,
		"action_type":        model.ActionUpdate,
	}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	result := &action.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Maker:           action.User{UserID: data.MakerUser.UserCode, FullName: data.MakerUser.FullName, PhoneNumber: data.MakerUser.PhoneNumber},
		Checker:         action.User{UserID: data.CheckerUser.UserCode, FullName: data.CheckerUser.FullName, PhoneNumber: data.CheckerUser.PhoneNumber},
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction:   data.PreviosAction,
		CurrentAction:   data.CurrentAction,
		ActionStatus:    action.ActionStatus(data.ActionStatus),
		ActionType:      action.ActionType(data.ActionType),
		RequestAction:   action.RequestAction(data.RequestAction),
		CreatedAt:       data.CreatedAt,
		LastModifiedAt:  data.LastModifiedAt,
	}
	return result, nil
}

func (o *outboundStore) GetCurrentPasswordRule(ctx context.Context) (*action.PasswordRule, error) {
	ruleModel, err := o.MongoDalPasswordRule.FindOne(ctx, bson.M{}, bson.M{})
	if err != nil || ruleModel == nil {
		return nil, err
	}

	rule := &action.PasswordRule{
		ID:             ruleModel.ID,
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
