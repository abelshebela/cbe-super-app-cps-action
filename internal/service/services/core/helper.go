package core

import (
	"cbe-super-app-cps-action/internal/constants"
	service_dto "cbe-super-app-cps-action/internal/constants/dto/services"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
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

func ValidateCreate(req service_dto.CreateServiceRequest, service storage.ServicesRepository) error {
	if req.ServiceCode == "" || req.ServiceName == "" {
		return localization.ErrorRequiredFieldMissing
	}

	filterCode := types.Filter{
		Filters: map[string]interface{}{
			"service_code": req.ServiceCode,
		},
	}
	serviceDocCode, err := service.FindAllWithPagination(context.Background(), filterCode)
	if err != nil {
		return err
	}
	if len(serviceDocCode.Data) > 0 {
		return localization.ErrorServiceExists
	}

	filterName := types.Filter{
		Filters: map[string]interface{}{
			"service_name": req.ServiceName,
		},
	}
	serviceDocName, err := service.FindAllWithPagination(context.Background(), filterName)
	if err != nil {
		return err
	}
	if len(serviceDocName.Data) > 0 {
		return localization.ErrorServiceExists
	}

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
					FeeType:   imodel.FeeType(t.FeeType),
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
					ServiceName:             sl.ServiceName,
					ServiceKey:              sl.ServiceKey,
					OverideProductGlAccount: sl.OverideProductGlAccount,
					IsEnabled:               true,
					OverideCap: imodel.Cap{
						SingleCap:          formatFloatPointer(sl.OverideCap.SingleCap),
						MinimumTransferCap: formatFloatPointer(sl.OverideCap.MinimumTransferCap),
					},
					OverideTiers: func() []imodel.Tier {
						tiers := make([]imodel.Tier, 0, len(sl.OverideTiers))
						for _, t := range sl.OverideTiers {
							tiers = append(tiers, imodel.Tier{
								FeeType:   imodel.FeeType(t.FeeType),
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
		Enabled:   true,
		IsDeleted: false,
		CreatedAt: time.Now(),
	}

	return mapped
}

func MapToServiceUpdateModel(req service_dto.UpdateServiceRequest) imodel.Service {
	return imodel.Service{
		ServiceCode:      req.ServiceCode,
		ServiceKey:       req.ServiceKey,
		ServiceName:      req.ServiceName,
		ProductGlAccount: req.ProductGlAccount,
		Cap: imodel.Cap{
			SingleCap:          formatFloatPointer(req.Cap.SingleCap),
			MinimumTransferCap: formatFloatPointer(req.Cap.MinimumTransferCap),
		},
		Tiers: func() []imodel.Tier {
			if len(req.Tiers) == 0 {
				return nil
			}
			tiers := make([]imodel.Tier, 0, len(req.Tiers))
			for _, t := range req.Tiers {
				tiers = append(tiers, imodel.Tier{
					FeeType:   imodel.FeeType(t.FeeType),
					FeeAmount: formatFloatPointer(t.FeeAmount),
					Min:       formatFloatPointer(t.Min),
					Max:       formatFloatPointer(t.Max),
				})
			}
			return tiers
		}(),
		ServiceList: func() []imodel.ServiceLists {
			if len(req.ServiceList) == 0 {
				return nil
			}
			lists := make([]imodel.ServiceLists, 0, len(req.ServiceList))
			for _, sl := range req.ServiceList {
				lists = append(lists, imodel.ServiceLists{
					ServiceName:             sl.ServiceName,
					ServiceKey:              sl.ServiceKey,
					OverideProductGlAccount: sl.OverideProductGlAccount,
					OverideCap: imodel.Cap{
						SingleCap:          formatFloatPointer(sl.OverideCap.SingleCap),
						MinimumTransferCap: formatFloatPointer(sl.OverideCap.MinimumTransferCap),
					},
					OverideTiers: func() []imodel.Tier {
						if len(sl.OverideTiers) == 0 {
							return nil
						}
						tiers := make([]imodel.Tier, 0, len(sl.OverideTiers))
						for _, t := range sl.OverideTiers {
							tiers = append(tiers, imodel.Tier{
								FeeType:   imodel.FeeType(t.FeeType),
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
	}
}
