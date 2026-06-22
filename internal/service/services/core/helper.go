package core

import (
	"cbe-super-app-cps-action/internal/constants"
	service_dto "cbe-super-app-cps-action/internal/constants/dto/services"
	"cbe-super-app-cps-action/internal/constants/lib"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"log"
	"strconv"
	"strings"
	"time"
	// shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func formatFloatPointer(f *float64) string {
	if f == nil {
		return ""
	}
	return strconv.FormatFloat(*f, 'f', -1, 64)
}

func sourceAppFromPtr(source *string) constants.SourceApp {
	return constants.SourceApp(service_dto.StringPointer(source, ""))
}

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func MapToServiceModel(req service_dto.CreateServiceRequest, accessList imodel.ServiceKey) imodel.Service {
	var glAccount string
	var glCurrency string
	if req.ProductGlAccount != "" {
		glAccount = req.ProductGlAccount
	}
	if req.ProductGlAccountCurrency != "" {
		glCurrency = req.ProductGlAccountCurrency
	}

	mapped := imodel.Service{
		ServiceName:              accessList.ServiceName,
		ServiceKey:               accessList.ServiceKey,
		ServiceCode:              req.ServiceCode,
		ServiceKeyId:             req.ServiceKeyId,
		IsPlAccount:              *req.IsPlAccount,
		ProductGlAccount:         glAccount,
		ProductGlAccountCurrency: glCurrency,
		Cap: func() []imodel.Cap {
			caps := make([]imodel.Cap, 0, len(req.Cap))

			for _, c := range req.Cap {
				caps = append(caps, imodel.Cap{
					Source:             sourceAppFromPtr(c.Source),
					Currency:           service_dto.StringPointer(c.Currency, ""),
					SingleCap:          formatFloatPointer(c.SingleCap),
					MinimumTransferCap: formatFloatPointer(c.MinimumTransferCap),
				})
			}

			return caps
		}(),
		MinimumFraudAmount: strconv.FormatFloat(req.MinimumFraudAmount, 'f', -1, 64),

		CreatedAt: time.Now(),
	}

	return mapped
}

func MapToServiceUpdateModel(req service_dto.UpdateServiceRequest, existing service_dto.ServiceResponse) service_dto.ServiceResponse {
	existing.ServiceCode = service_dto.StringPointer(req.ServiceCode, existing.ServiceCode)
	existing.ServiceKeyId = service_dto.StringPointer(req.ServiceKeyId, existing.ServiceKeyId)
	existing.IsPlAccount = service_dto.BoolPointer(req.IsPlAccount, existing.IsPlAccount)
	existing.ProductGlAccountCurrency = service_dto.StringPointer(req.ProductGlAccountCurrency, existing.ProductGlAccountCurrency)
	// existing.ServiceKey = service_dto.StringPointer(req.ServiceKey, existing.ServiceKey)
	// existing.ServiceName = service_dto.StringPointer(req.ServiceName, existing.ServiceName)
	existing.ProductGlAccount = service_dto.StringPointer(req.ProductGlAccount, existing.ProductGlAccount)

	if req.Cap != nil {
		existing.Cap = func() []imodel.Cap {
			caps := make([]imodel.Cap, 0, len(req.Cap))

			for _, c := range req.Cap {
				caps = append(caps, imodel.Cap{
					Source:             constants.SourceApp(*c.Source),
					Currency:           service_dto.StringPointer(c.Currency, ""),
					SingleCap:          formatFloatPointer(c.SingleCap),
					MinimumTransferCap: formatFloatPointer(c.MinimumTransferCap),
				})
			}

			return caps
		}()
	}
	if req.MinimumFraudAmount != nil {
		existing.MinimumFraudAmount = strconv.FormatFloat(*req.MinimumFraudAmount, 'f', -1, 64)
	}

	// if req.Tiers != nil {
	// 	existing.Tiers = func() []model.Tier {
	// 		tiers := make([]model.Tier, 0, len(req.Tiers))
	// 		for _, t := range req.Tiers {
	// 			tiers = append(tiers, model.Tier{
	// 				// FeeType:   model.FeeType(service_dto.StringPointer(t.FeeType, string(existing.Tiers[0].FeeType))), // Use existing or first tier if matching is complex
	// 				FeeType:   shared_constants.FeeType(service_dto.StringPointer(t.FeeType, string(*t.FeeType))),
	// 				FeeAmount: formatFloatPointer(t.FeeAmount),
	// 				Min:       formatFloatPointer(t.Min),
	// 				Max:       formatFloatPointer(t.Max),
	// 			})
	// 		}
	// 		return tiers
	// 	}()
	// }

	// if req.ServiceList != nil {
	// 	existing.ServiceList = func() []model.ServiceList {
	// 		lists := make([]model.ServiceList, 0, len(req.ServiceList))
	// 		for _, sl := range req.ServiceList {
	// 			lists = append(lists, model.ServiceList{
	// 				ServiceName:             service_dto.StringPointer(sl.ServiceName, ""),
	// 				ServiceKey:              service_dto.StringPointer(sl.ServiceKey, ""),
	// 				OverideProductGlAccount: service_dto.StringPointer(sl.OverideProductGlAccount, ""),
	// 				OverideCap: func() model.Cap {
	// 					if sl.OverideCap == nil {
	// 						return model.Cap{}
	// 					}
	// 					return model.Cap{
	// 						SingleCap:          formatFloatPointer(sl.OverideCap.SingleCap),
	// 						MinimumTransferCap: formatFloatPointer(sl.OverideCap.MinimumTransferCap),
	// 					}
	// 				}(),
	// 				// OverideTiers: func() []model.Tier {
	// 				// 	tiers := make([]model.Tier, 0, len(sl.OverideTiers))
	// 				// 	for _, t := range sl.OverideTiers {
	// 				// 		tiers = append(tiers, model.Tier{
	// 				// 			// FeeType:   model.FeeType(service_dto.StringPointer(t.FeeType, "")),
	// 				// 			FeeType:   shared_constants.FeeType(service_dto.StringPointer(t.FeeType, string(*t.FeeType))),
	// 				// 			FeeAmount: formatFloatPointer(t.FeeAmount),
	// 				// 			Min:       formatFloatPointer(t.Min),
	// 				// 			Max:       formatFloatPointer(t.Max),
	// 				// 		})
	// 				// 	}
	// 				// 	return tiers
	// 				// }(),
	// 				// IsEnabled: service_dto.BoolPointer(sl.IsEnabled, true),
	// 				IsEnabled: true,
	// 			})
	// 		}
	// 		return lists
	// 	}()
	// 	existing.Enabled = true
	// }

	return existing
}

func MapServiceListDtoToModel(req *service_dto.CreateServiceList) imodel.ServiceKey {
	return imodel.ServiceKey{
		ServiceName:       req.ServiceName,
		ServiceKey:        req.ServiceKey,
		AccountType:       req.AccountType,
		IsSuperAppEnabled: *req.IsSuperAppEnabled,
		IsUSSDEnabled:     *req.IsUSSDEnabled,
	}
}

func MapServiceListDtoUpdateToModel(existing imodel.ServiceKey, req *service_dto.UpdateServiceList) imodel.ServiceKey {
	var name, key, aType string

	if strings.TrimSpace(req.ServiceName) == "" {
		name = existing.ServiceName
	}
	if strings.TrimSpace(req.ServiceKey) == "" {
		key = existing.ServiceKey
	}
	if strings.TrimSpace(req.AccountType) == "" {
		aType = existing.AccountType
	}

	return imodel.ServiceKey{
		ServiceName:       name,
		ServiceKey:        key,
		AccountType:       aType,
		IsSuperAppEnabled: *req.IsSuperAppEnabled,
		IsUSSDEnabled:     *req.IsUSSDEnabled,
	}
}

// BuildServiceDeleteSnapshot preserves the full service record for CPS audit before hard delete.
func BuildServiceDeleteSnapshot(prev service_dto.ServiceResponse) service_dto.ServiceResponse {
	snap := prev
	now := time.Now()
	snap.IsDeleted = true
	snap.DeletedAt = &now
	snap.LastModifiedAt = now
	return snap
}

// BuildServiceKeyDeleteSnapshot preserves the full access list record for CPS audit before hard delete.
func BuildServiceKeyDeleteSnapshot(prev imodel.ServiceKey) imodel.ServiceKey {
	snap := prev
	now := time.Now()
	snap.DeletedAt = &now
	snap.LastModifiedAt = now
	return snap
}
