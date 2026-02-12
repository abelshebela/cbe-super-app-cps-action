package core

import (
	"cbe-super-app-cps-action/internal/constants"
	service_dto "cbe-super-app-cps-action/internal/constants/dto/services"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"log"
	"strconv"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func formatFloatPointer(f *float64) string {
	if f == nil {
		return ""
	}
	return strconv.FormatFloat(*f, 'f', -1, 64)
}

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	// if incomplet := local_util.IsIncomplete(userData); incomplet {
	// 	log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
	// 	return errors.New(localization.ErrorIncompleteUserInfo.Code)
	// }

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func MapToServiceModel(req service_dto.CreateServiceRequest) model.Service {
	mapped := model.Service{
		ServiceCode:      req.ServiceCode,
		ServiceKey:       req.ServiceKey,
		ServiceName:      req.ServiceName,
		ProductGlAccount: req.ProductGlAccount,
		Cap: model.Cap{
			SingleCap:          formatFloatPointer(req.Cap.SingleCap),
			MinimumTransferCap: formatFloatPointer(req.Cap.MinimumTransferCap),
		},
		// Tiers: func() []model.Tier {
		// 	tiers := make([]model.Tier, 0, len(req.Tiers))
		// 	for _, t := range req.Tiers {
		// 		tiers = append(tiers, model.Tier{
		// 			// FeeType:   model.FeeType(service_dto.StringPointer(t.FeeType, "")),
		// 			FeeType:   shared_constants.FeeType(service_dto.StringPointer(t.FeeType, string(*t.FeeType))),
		// 			FeeAmount: formatFloatPointer(t.FeeAmount),
		// 			Min:       formatFloatPointer(t.Min),
		// 			Max:       formatFloatPointer(t.Max),
		// 		})
		// 	}
		// 	return tiers
		// }(),
		// ServiceList: func() []model.ServiceList {
		// 	lists := make([]model.ServiceList, 0, len(req.ServiceList))
		// 	for _, sl := range req.ServiceList {
		// 		lists = append(lists, model.ServiceList{
		// 			ServiceName:             service_dto.StringPointer(sl.ServiceName, ""),
		// 			ServiceKey:              service_dto.StringPointer(sl.ServiceKey, ""),
		// 			OverideProductGlAccount: service_dto.StringPointer(sl.OverideProductGlAccount, ""),
		// 			IsEnabled:               service_dto.BoolPointer(sl.IsEnabled, true),
		// 			OverideCap: func() model.Cap {
		// 				if sl.OverideCap == nil {
		// 					return model.Cap{}
		// 				}
		// 				return model.Cap{
		// 					SingleCap:          formatFloatPointer(sl.OverideCap.SingleCap),
		// 					MinimumTransferCap: formatFloatPointer(sl.OverideCap.MinimumTransferCap),
		// 				}
		// 			}(),
		// 			// OverideTiers: func() []model.Tier {
		// 			// 	tiers := make([]model.Tier, 0, len(sl.OverideTiers))
		// 			// 	for _, t := range sl.OverideTiers {
		// 			// 		tiers = append(tiers, model.Tier{
		// 			// 			// FeeType:   model.FeeType(service_dto.StringPointer(t.FeeType, "")),
		// 			// 			FeeType:   shared_constants.FeeType(service_dto.StringPointer(t.FeeType, string(*t.FeeType))),
		// 			// 			FeeAmount: formatFloatPointer(t.FeeAmount),
		// 			// 			Min:       formatFloatPointer(t.Min),
		// 			// 			Max:       formatFloatPointer(t.Max),
		// 			// 		})
		// 			// 	}
		// 			// 	return tiers
		// 			// }(),
		// 		})
		// 	}
		// 	return lists
		// }(),
		Enabled:   service_dto.BoolPointer(req.Enabled, true),
		IsDeleted: service_dto.BoolPointer(req.IsDeleted, false),
		CreatedAt: time.Now(),
	}

	return mapped
}

func MapToServiceUpdateModel(req service_dto.UpdateServiceRequest, existing model.Service) model.Service {
	existing.ServiceCode = service_dto.StringPointer(req.ServiceCode, existing.ServiceCode)
	existing.ServiceKey = service_dto.StringPointer(req.ServiceKey, existing.ServiceKey)
	existing.ServiceName = service_dto.StringPointer(req.ServiceName, existing.ServiceName)
	existing.ProductGlAccount = service_dto.StringPointer(req.ProductGlAccount, existing.ProductGlAccount)

	if req.Cap != nil {
		existing.Cap = model.Cap{
			SingleCap:          formatFloatPointer(req.Cap.SingleCap),
			MinimumTransferCap: formatFloatPointer(req.Cap.MinimumTransferCap),
		}
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
