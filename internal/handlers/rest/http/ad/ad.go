package ad

import (
	"cbe-super-app-cps-action/internal/handlers/rest/http/ad/core"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	ad_dto "cbe-super-app-cps-action/internal/constants/dto/ad"
	advertInbound "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type advertAdapter struct {
	advertApplication service.AdvertService
	logger            utils.Logger
}

func InitAdvertAdapter(advertApplication service.AdvertService, logger utils.Logger) advertInbound.ADAdapter {
	return &advertAdapter{
		logger:            logger,
		advertApplication: advertApplication,
	}
}

func (a *advertAdapter) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	req, err := core.ParseAndValidateAdvertRequest(w, r, false, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	domainReq, err := core.ToAdvert(*req)
	if err != nil {
		a.logger.Errorf("[event.CreateAdvert] failed to convert to domain advert, error: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequest.Message)
		return
	}

	err = a.advertApplication.CreateAdvert(r.Context(), &domainReq, req.BannerImage)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertCreateRequestSent, nil)

}
func (a *advertAdapter) FetchAdverts(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	list, err := a.advertApplication.FetchAdverts(r.Context(), *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	docs := core.ToAdvertResponses(list.Data)
	res := types.PaginatedResponse[[]*ad_dto.AdvertResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	localization.SendSuccessResponse(w, localization.SuccessAdvertsFetched, res)
}
func (a *advertAdapter) FetchAdvertByID(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(w, r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	data, err := a.advertApplication.FetchAdvertByID(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	res := core.ToAdvertResponse(*data)
	localization.SendSuccessResponse(w, localization.SuccessAdvertFetched, res)
}
func (a *advertAdapter) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(w, r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	req, err := core.ParseAndValidateAdvertRequest(w, r, true, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	domainReq, _ := core.ToAdvert(*req)

	if domainReq.Title == "" && domainReq.Description == "" && domainReq.AdvertFor == "" && req.BannerImage == nil && domainReq.Date.StartedAt.IsZero() && domainReq.Date.ExpiredAt.IsZero() {
		a.logger.Errorf("[event.UpdateAdvert] no data provided for update, id: %s", id)
		localization.SendBadRequestResponse(w, localization.ErrorNoDataProvidedForUpdate.Message)
		return
	}

	err = a.advertApplication.UpdateAdvert(r.Context(), id, &domainReq, req.BannerImage)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertUpdateRequestSent, nil)
}
func (a *advertAdapter) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(w, r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	err = a.advertApplication.DeleteAdvert(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAdvertDeleteRequestSent, nil)
}
func (a *advertAdapter) EnableAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(w, r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	err = a.advertApplication.EnableDisableAdvert(r.Context(), id, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertEnableRequestSent, nil)
}
func (a *advertAdapter) DisableAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(w, r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
	}

	err = a.advertApplication.EnableDisableAdvert(r.Context(), id, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertDisableRequestSent, nil)
}
