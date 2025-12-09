package sitota

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"
	"strings"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type handler struct {
	svc    service.SitotaService
	logger utils.Logger
}

func InitSitotaHandler(svc service.SitotaService, logger utils.Logger) *handler {
	return &handler{svc: svc, logger: logger}
}

// GetAllSitotas godoc
//
//	@Summary		Get all sitota transactions
//	@Description	Retrieve all sitota transactions from the gRPC service.
//	@Tags			Sitota
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	localization.StandardResponse{data=[]model.SitotaTransaction}	"Sitota transactions retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}							"Server error"
//	@Security		BearerAuth
//	@Router			/sitotas [get]
func (h *handler) GetAllSitotas(w http.ResponseWriter, r *http.Request) {
	params := local_util.ExtractFilterParams(r)
	sitotas, err := h.svc.GetAllSitotas(r.Context(), params)
	if err != nil {
		h.logger.Errorf("[GetAllSitotas] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[GetAllSitotas] retrieved %d sitota transactions", len(sitotas.Data))
	localization.SendSuccessResponse(w, localization.SuccessAllSitotasRetrieved, sitotas)
}

// GetSitota godoc
//
//	@Summary		Get sitota transaction by ID
//	@Description	Retrieve a specific sitota transaction by its transaction ID.
//	@Tags			Sitota
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string														true	"Transaction ID"	example("TXN123456789")
//	@Success		200	{object}	localization.StandardResponse{data=model.SitotaTransaction}	"Sitota transaction retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}						"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}						"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}						"Server error"
//	@Security		BearerAuth
//	@Router			/sitotas/{id} [get]
func (h *handler) GetSitota(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorSitotaRequired.Code)
		return
	}

	sitota, err := h.svc.GetSitotaByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[GetSitota] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[GetSitota] sitota transaction retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessSitotaRetrieved, sitota)
}
