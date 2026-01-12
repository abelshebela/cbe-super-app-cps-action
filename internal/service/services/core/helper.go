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
	"time"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
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

func MapToServiceModel(req service_dto.CreateServiceRequest) imodel.Service {
	mapped := imodel.Service{
		ServiceCode:      req.ServiceCode,
		ServiceKey:       req.ServiceKey,
		ServiceName:      req.ServiceName,
		ProductGlAccount: req.ProductGlAccount,
		Cap: imodel.Cap{
			SingleCap:          formatFloatPointer(req.Cap.SingleCap),
			MinimumTransferCap: formatFloatPointer(req.Cap.MinimumTransferCap),
		},
		Tiers: func() []imodel.Tier {
			tiers := make([]imodel.Tier, 0, len(req.Tiers))
			for _, t := range req.Tiers {
				tiers = append(tiers, imodel.Tier{
					FeeType:   imodel.FeeType(service_dto.StringPointer(t.FeeType, "")),
					FeeAmount: formatFloatPointer(t.FeeAmount),
					Min:       formatFloatPointer(t.Min),
					Max:       formatFloatPointer(t.Max),
				})
			}
			return tiers
		}(),
		ServiceList: func() []imodel.ServiceLists {
			lists := make([]imodel.ServiceLists, 0, len(req.ServiceList))
			for _, sl := range req.ServiceList {
				lists = append(lists, imodel.ServiceLists{
					ServiceName:             service_dto.StringPointer(sl.ServiceName, ""),
					ServiceKey:              service_dto.StringPointer(sl.ServiceKey, ""),
					OverideProductGlAccount: service_dto.StringPointer(sl.OverideProductGlAccount, ""),
					IsEnabled:               service_dto.BoolPointer(sl.IsEnabled, true),
					OverideCap: func() imodel.Cap {
						if sl.OverideCap == nil {
							return imodel.Cap{}
						}
						return imodel.Cap{
							SingleCap:          formatFloatPointer(sl.OverideCap.SingleCap),
							MinimumTransferCap: formatFloatPointer(sl.OverideCap.MinimumTransferCap),
						}
					}(),
					OverideTiers: func() []imodel.Tier {
						tiers := make([]imodel.Tier, 0, len(sl.OverideTiers))
						for _, t := range sl.OverideTiers {
							tiers = append(tiers, imodel.Tier{
								FeeType:   imodel.FeeType(service_dto.StringPointer(t.FeeType, "")),
								FeeAmount: formatFloatPointer(t.FeeAmount),
								Min:       formatFloatPointer(t.Min),
								Max:       formatFloatPointer(t.Max),
							})
						}
						return tiers
					}(),
				})
			}
			return lists
		}(),
		Enabled:   service_dto.BoolPointer(req.Enabled, true),
		IsDeleted: service_dto.BoolPointer(req.IsDeleted, false),
		CreatedAt: time.Now(),
	}

	return mapped
}

func MapToServiceUpdateModel(req service_dto.UpdateServiceRequest, existing imodel.Service) imodel.Service {
	existing.ServiceCode = service_dto.StringPointer(req.ServiceCode, existing.ServiceCode)
	existing.ServiceKey = service_dto.StringPointer(req.ServiceKey, existing.ServiceKey)
	existing.ServiceName = service_dto.StringPointer(req.ServiceName, existing.ServiceName)
	existing.ProductGlAccount = service_dto.StringPointer(req.ProductGlAccount, existing.ProductGlAccount)

	if req.Cap != nil {
		existing.Cap = imodel.Cap{
			SingleCap:          formatFloatPointer(req.Cap.SingleCap),
			MinimumTransferCap: formatFloatPointer(req.Cap.MinimumTransferCap),
		}
	}

	if req.Tiers != nil {
		existing.Tiers = func() []imodel.Tier {
			tiers := make([]imodel.Tier, 0, len(req.Tiers))
			for _, t := range req.Tiers {
				tiers = append(tiers, imodel.Tier{
					FeeType:   imodel.FeeType(service_dto.StringPointer(t.FeeType, string(existing.Tiers[0].FeeType))), // Use existing or first tier if matching is complex
					FeeAmount: formatFloatPointer(t.FeeAmount),
					Min:       formatFloatPointer(t.Min),
					Max:       formatFloatPointer(t.Max),
				})
			}
			return tiers
		}()
	}

	if req.ServiceList != nil {
		existing.ServiceList = func() []imodel.ServiceLists {
			lists := make([]imodel.ServiceLists, 0, len(req.ServiceList))
			for _, sl := range req.ServiceList {
				lists = append(lists, imodel.ServiceLists{
					ServiceName:             service_dto.StringPointer(sl.ServiceName, ""),
					ServiceKey:              service_dto.StringPointer(sl.ServiceKey, ""),
					OverideProductGlAccount: service_dto.StringPointer(sl.OverideProductGlAccount, ""),
					OverideCap: func() imodel.Cap {
						if sl.OverideCap == nil {
							return imodel.Cap{}
						}
						return imodel.Cap{
							SingleCap:          formatFloatPointer(sl.OverideCap.SingleCap),
							MinimumTransferCap: formatFloatPointer(sl.OverideCap.MinimumTransferCap),
						}
					}(),
					OverideTiers: func() []imodel.Tier {
						tiers := make([]imodel.Tier, 0, len(sl.OverideTiers))
						for _, t := range sl.OverideTiers {
							tiers = append(tiers, imodel.Tier{
								FeeType:   imodel.FeeType(service_dto.StringPointer(t.FeeType, "")),
								FeeAmount: formatFloatPointer(t.FeeAmount),
								Min:       formatFloatPointer(t.Min),
								Max:       formatFloatPointer(t.Max),
							})
						}
						return tiers
					}(),
					// IsEnabled: service_dto.BoolPointer(sl.IsEnabled, true),
					IsEnabled: true,
				})
			}
			return lists
		}()
		existing.Enabled = true
	}

	return existing
}
