package bulk

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"errors"
	// "strconv"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"

	// "time"

	"cbe-super-app-cps-action/internal/constants/lib"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bulkService struct {
	cpsActionRepo service.CPSActionService
	repo          storage.BulkServiceRepository
	logger        utils.Logger
}

func NewBulkServicwAdapter(repo storage.BulkServiceRepository, CpsActionRepo service.CPSActionService, logger utils.Logger) service.BulkService {
	return &bulkService{
		cpsActionRepo: CpsActionRepo,
		repo:          repo,
		logger:        logger,
	}
}

func (s *bulkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	Now := time.Now()
	cpsAction.MakerActionTime = Now
	cpsAction.LastModifiedAt = Now
	updateData := cpsAction.CurrentAction.([]*model.APPAccessList)

	allAccessLists, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	validAccessList := storedAccessListToMAP(allAccessLists)
	invalidDatas, validKeys := validaterAccessKey(validAccessList, updateData)

	if len(invalidDatas) > 0 {
		s.logger.Errorf("one or more Servive not found")
		return nil, fmt.Errorf("one or more Service not found")
	}
	for key, value := range validAccessList {
		switch strings.ToUpper(cpsAction.RejectionReason) {
		case string(constants.RequestBulkServiceEnable):
			if value == true {
				s.logger.Infof("you enter already enabled service: %v", key)
				return nil, errors.New("you enter already enabled service")
			}
		case string(constants.RequestBulkServiceDisable):
			if value == false {
				s.logger.Infof("you enter already disabled service: %v", key)
				return nil, errors.New("you enter already disabled service")
			}
		}
	}
	keys := GetAllKeysFromMaps(validKeys)

	switch cpsAction.RequestAction {
	case string(constants.RequestBulkServiceEnable):
		return nil, s.repo.Update(ctx, keys, true)
	case string(constants.RequestBulkServiceDisable):
		return nil, s.repo.Update(ctx, keys, false)
	default:
		return nil, fmt.Errorf("Unknown Request Action")
	}

	// return s.repo.Update(ctx, &updateData)

}
func GetAllKeysFromMaps(maps map[string]bool) []string {
	var keys []string
	// keySet := make(map[string]struct{}) // to avoid duplicates

	for key, _ := range maps {
		keys = append(keys, key)
	}
	return keys
}

// this function change the accesslist to key and enable map user to check the validity of  key given and check the current state
func storedAccessListToMAP(AccessLists []*model.APPAccessList) map[string]bool {
	validKeys := make(map[string]bool)
	for _, access := range AccessLists {
		validKeys[access.Key] = access.Enabled
		for _, sub_access := range access.SubAccessList {
			validKeys[sub_access.Key] = sub_access.Enabled
		}
	}

	return validKeys
}

func validaterAccessKey(validAccessMap map[string]bool, accessList []*model.APPAccessList) ([]string, map[string]bool) {
	var invalidKeys []string
	validKeys := make(map[string]bool)

	for _, access := range accessList {
		if enabled, exists := validAccessMap[access.Key]; exists {
			validKeys[access.Key] = enabled
		} else {
			invalidKeys = append(invalidKeys, access.Key)
		}
		for _, subAccess := range access.SubAccessList {
			if enabled, exists := validAccessMap[subAccess.Key]; exists {
				validKeys[subAccess.Key] = enabled
			} else {
				invalidKeys = append(invalidKeys, subAccess.Key)
			}
		}
	}
	return invalidKeys, validKeys
}

func (s *bulkService) GetAllBulkServices(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error) {
	return s.repo.FindAllWithPagination(ctx, filterParams)
}
func (s *bulkService) EnableBulkService(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}
	userPayload := local_util.ExtractUserFromContext(ctx)

	type currAction struct {
		Keys []string
	}
	cpsAction := lib.CpsModelBuilder("", userPayload, nil, currAction{
		Keys: keys,
	}, string(constants.RequestBulkServiceDisable), string(constants.UpdateAction))

	return s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction)
}

func (s *bulkService) DisableBulkService(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}

	userPayload := local_util.ExtractUserFromContext(ctx)

	type currAction struct {
		Keys []string
	}
	cpsAction := lib.CpsModelBuilder("", userPayload, nil, currAction{
		Keys: keys,
	}, string(constants.RequestBulkServiceDisable), string(constants.UpdateAction))

	return s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction)
}

// package bulk_service

// import (
// 	"cbe-super-app-cps-action/internal/constants/dto"
// 	"cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"

// 	"cbe-super-app-cps-action/internal/service"
// 	"encoding/json"
// 	"net/http"

// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	util "cbe-super-app-cps-action/pkgs/utils"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// )

// type bulk_serviceAdapter struct {
// 	bulkService service.BulkService
// 	logger      utils.Logger
// }

// func InitBulkServiceAdapter(bulk_service service.BulkService, logger utils.Logger) bulk_service.BulkServiceHandler {
// 	return &bulk_serviceAdapter{
// 		logger:      logger,
// 		bulkService: bulk_service,
// 	}
// }

// func (h *bulk_serviceAdapter) GetAllBulkServices(w http.ResponseWriter, r *http.Request) {
// 	filter_params := util.ExtractFilterParams(r)
// 	bulk_services, err := h.bulkService.GetAllBulkServices(r.Context(), filter_params)
// 	if err != nil {
// 		h.logger.Errorf("Error while fetch all bulk services: %v\n", err)
// 		localization.SendErrorByCodeResponse(w, localization.MsgUnableToFetchBulkService.Code)
// 		return
// 	}

// 	localization.SendSuccessResponse(w, localization.BulkServiceFetchSuccessfully, bulk_services)

// }

// func (h *bulk_serviceAdapter) EnableBulkService(w http.ResponseWriter, r *http.Request) {

// 	var req dto.BulkServiceDTO

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.logger.Errorf("failed to decode payload")
// 	}

// 	err := h.bulkService.EnableBulkService(r.Context(), req.Keys)
// 	if err != nil {
// 		h.logger.Errorf("Enable bulk service request failed: %v\n", err)
// 		localization.SendBadRequestResponse(w, localization.MsgBadRequest)
// 		return
// 	}

// 	localization.SendSuccessResponse(w, localization.BulkServiceEnableRequestSuccess, nil)
// }

// func (h *bulk_serviceAdapter) DisableBulkService(w http.ResponseWriter, r *http.Request) {
// 	var req dto.BulkServiceDTO
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.logger.Errorf("failed to decode payload")
// 	}

// 	err := h.bulkService.DisableBulkService(r.Context(), req.Keys)
// 	if err != nil {
// 		h.logger.Errorf("Disable bulk service request failed: %v\n", err)
// 		localization.SendBadRequestResponse(w, localization.MsgBadRequest)
// 		return
// 	}

// 	localization.SendSuccessResponse(w, localization.BulkServiceDisableRequestSuccess, nil)
// }
