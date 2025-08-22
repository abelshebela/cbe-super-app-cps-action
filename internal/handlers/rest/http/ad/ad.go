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
	req, ok := core.ParseAndValidateAdvertRequest(w, r, false, a.logger)
	if !ok {
		return
	}

	maker, ok := core.ExtractUserAndMaker(w, r, a.logger)
	if !ok {
		return
	}

	domainReq, _ := core.ToAdvert(req)
	err := a.advertApplication.CreateAdvert(r.Context(), &domainReq, req.BannerImage, maker)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
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
	id, ok := core.ExtractID(w, r, a.logger)
	if !ok {
		return
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
	id, ok := core.ExtractID(w, r, a.logger)
	if !ok {
		return
	}

	req, ok := core.ParseAndValidateAdvertRequest(w, r, true, a.logger)
	if !ok {
		return
	}

	maker, ok := core.ExtractUserAndMaker(w, r, a.logger)
	if !ok {
		return
	}

	domainReq, _ := core.ToAdvert(req)

	if domainReq.Title == "" && domainReq.Description == "" && domainReq.AdvertFor == "" && req.BannerImage == nil && domainReq.Date.StartedAt.IsZero() && domainReq.Date.ExpiredAt.IsZero() {
		a.logger.Errorf("[event.UpdateAdvert] no data provided for update, id: %s", id)
		localization.SendBadRequestResponse(w, localization.ErrorNoDataProvidedForUpdate.Message)
		return
	}

	err := a.advertApplication.UpdateAdvert(r.Context(), id, &domainReq, req.BannerImage, maker)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertUpdateRequestSent, nil)
}
func (a *advertAdapter) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ExtractID(w, r, a.logger)
	if !ok {
		return
	}

	maker, ok := core.ExtractUserAndMaker(w, r, a.logger)
	if !ok {
		return
	}

	err := a.advertApplication.DeleteAdvert(r.Context(), id, maker)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w,localization.SuccessAdvertDeleteRequestSent, nil)
}
func (a *advertAdapter) EnableAdvert(w http.ResponseWriter, r *http.Request)  {
	id, ok := core.ExtractID(w, r, a.logger)
	if !ok {
		return
	}

	maker, ok := core.ExtractUserAndMaker(w, r, a.logger)
	if !ok {
		return
	}

	err := a.advertApplication.EnableDisableAdvert(r.Context(), id, maker, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return	
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertEnableRequestSent, nil)
}
func (a *advertAdapter) DisableAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ExtractID(w, r, a.logger)
	if !ok {
		return
	}

	maker, ok := core.ExtractUserAndMaker(w, r, a.logger)
	if !ok {
		return
	}

	err := a.advertApplication.EnableDisableAdvert(r.Context(), id, maker, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return	
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertDisableRequestSent, nil)
}
