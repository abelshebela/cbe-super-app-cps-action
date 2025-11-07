package sitota

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"
	"strings"

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

func (h *handler) GetAllSitotas(w http.ResponseWriter, r *http.Request) {
	sitotas, err := h.svc.GetAllSitotas(r.Context())
	if err != nil {
		h.logger.Errorf("[GetAllSitotas] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAllSitotasRetrieved, sitotas)
}

func (h *handler) GetSitota(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorSitotaRequired.Code)
		return
	}

	sitota, err := h.svc.GetSitotaByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[GetSitota] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessSitotaRetrieved, sitota)
}
