package updatedbulkservice

import (
	"net/http"

	application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/updated_bulk_service"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/updated_bulk_service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BulkServiceHandler struct {
	app    application.BulkServiceApplication
	logger utils.Logger
}

func InitBulkServiceHandler(app application.BulkServiceApplication, logger utils.Logger) inbound.BulkServiceHandler {
	return BulkServiceHandler{
		app:    app,
		logger: logger,
	}
}

func (h BulkServiceHandler) GetAllBulkServices(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	bulkServices, err := h.app.GetAllBulkServices(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("Error while getting all bulk services: %v\n", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	doc, _ := common_util.StructToMap(bulkServices)
	common_util.BaseResponseMaker(doc, w, "Bulk Services fetched successfully", 200)
}

// func (h BulkServiceHandler) EnableBulkService(w http.ResponseWriter, r *http.Request)  {}
// func (h BulkServiceHandler) DisableBulkService(w http.ResponseWriter, r *http.Request) {}
