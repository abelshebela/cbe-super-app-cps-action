package ussd_merchant

import (
	"cbe-super-app-cps-action/internal/constants"
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	ussd_merchant_interface "cbe-super-app-cps-action/internal/constants/interfaces/ussd_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/ussd_merchant/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"

	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UssdMerchantHandler struct {
	UssdMerchantService service.UssdMerchantService
	Logger              utils.Logger
}

func NewUssdMerchantHandler(ussdMerchantService service.UssdMerchantService, logger utils.Logger) ussd_merchant_interface.UssdMerchantInbound {
	return &UssdMerchantHandler{
		UssdMerchantService: ussdMerchantService,
		Logger:              logger,
	}
}

func (h *UssdMerchantHandler) CreateUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createUssdMerchant", "handler", "ussdMerchant")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var CreateDto ussd_merchant_dto.CreateUssdMerchantRequest

	if err := json.NewDecoder(r.Body).Decode(&CreateDto); err != nil {
		h.Logger.Errorf("[CreateUssdMerchant] failed to read json body error: %v", err)
		localization.ErrorToResponseCode(err.Error(), http.StatusBadRequest, err.Error())
		return
	}

	err := core.ParseInputData(ctx, r, &CreateDto)
	if err != nil {
		h.Logger.Errorf("[CreateUssdMerchant]failed to read json body error: %v", err)
		localization.ErrorToResponseCode(err.Error(), http.StatusBadRequest, err.Error())
		return
	}

	if err := h.UssdMerchantService.CreateUssdMerchant(ctx, CreateDto); err != nil {
		h.Logger.Errorf("[CreateUssdMerchant] failed to create ussd merchant error: %v", err)
		localization.ErrorToResponseCode(err.Error(), http.StatusInternalServerError, err.Error())
		return
	}

	if md.IsMakerOnly {
		h.logger.Infof("[CreateUssdMerchant] request sent successfully: is_maker_only: %v", md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessCpsUserCreated, nil)
		return
	}

	h.Logger.Infof("[CreateUssdMerchant] successfully created ussd merchant")
	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantCreated, "Ussd merchant created successfully")
}

func (h *UssdMerchantHandler) UpdateUssdMerchant(w http.ResponseWriter, r *http.Request) {

}
func (h *UssdMerchantHandler) EnableUssdMerchant(w http.ResponseWriter, r *http.Request) {

}
func (h *UssdMerchantHandler) DisableUssdMerchant(w http.ResponseWriter, r *http.Request) {

}
func (h *UssdMerchantHandler) GetUssdMerchant(w http.ResponseWriter, r *http.Request) {

}
func (h *UssdMerchantHandler) GetAllUssdMerchant(w http.ResponseWriter, r *http.Request) {

}
