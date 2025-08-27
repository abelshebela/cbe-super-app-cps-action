package bulk_service

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/bulk_service"
	"cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	"fmt"

	// "fmt"

	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bulk_serviceAdapter struct {
	bulkService service.BulkService
	logger      utils.Logger
}

func InitBulkServiceAdapter(bulk_service service.BulkService, logger utils.Logger) bulk_service.BulkServiceHandler {
	return &bulk_serviceAdapter{
		logger:      logger,
		bulkService: bulk_service,
	}
}

func (h *bulk_serviceAdapter) GetAllBulkServices(w http.ResponseWriter, r *http.Request) {
	filter_params := util.ExtractFilterParams(r)
	bulk_services, err := h.bulkService.GetAllBulkServices(r.Context(), filter_params)
	if err != nil {
		h.logger.Errorf("Error while fetch all bulk services: %v\n", err)
		localization.SendErrorByCodeResponse(w, localization.UnableToFetchBulkService.Code)
		return
	}
	localization.SendSuccessResponse(w, localization.BulkServiceFetchSuccessfully, bulk_services)

}

func (h *bulk_serviceAdapter) EnableBulkService(w http.ResponseWriter, r *http.Request) {
	var req dto.BulkServiceDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode payload")
	}

	action_code, err := h.bulkService.EnableBulkService(r.Context(), req.Keys)
	if err != nil {
		h.logger.Errorf("Enable bulk service request failed: %v\n", err)
		localization.SendBadRequestResponse(w, localization.MsgBadRequest)
		return
	}

	localization.SendSuccessResponse(w, localization.BulkServiceEnableRequestSuccess, action_code)
}

func (h *bulk_serviceAdapter) DisableBulkService(w http.ResponseWriter, r *http.Request) {
	var req dto.BulkServiceDTO
	fmt.Println("handelr 1")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode payload")
	}
	fmt.Println("handelr 2")
	action_code, err := h.bulkService.DisableBulkService(r.Context(), req.Keys)
	fmt.Println(err)
	if err != nil {
		h.logger.Errorf("Disable bulk service request failed: %v\n", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	fmt.Println("handelr 3")
	localization.SendSuccessResponse(w, localization.BulkServiceDisableRequestSuccess, action_code)
}
