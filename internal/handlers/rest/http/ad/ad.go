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

type paginated_advert_response types.PaginatedResponse[[]*ad_dto.AdvertResponse]
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

// CreateAdvert godoc
// @Summary Create a new advert (maker)
// @Description Submit an advert create request. Requires multipart/form-data with optional banner image.
// @Tags Adverts
// @Accept mpfd
// @Produce json
// @Param title formData string true "Title" minLength(3) maxLength(20) example("New Promo")
// @Param description formData string true "Description" minLength(30) maxLength(100) example("Enjoy our new promotion valid this weekend only.")
// @Param advert_for formData string true "Advert audience" Enums(IFB,CB,ALL) example(ALL)
// @Param banner_image formData file false "Banner image (<=2MB; jpeg/png/gif/webp)"
// @Success 200 {object} localization.StandardResponse{data=nil} "Advert create request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid request"
// @Failure 409 {object} localization.StandardResponse{data=nil} "Duplicate or conflict"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /adverts [post]
func (a *advertAdapter) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	req, err := core.ParseAndValidateAdvertRequest(r, false, a.logger)
	if err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
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

// FetchAdverts godoc
// @Summary List adverts
// @Description Fetch adverts with pagination and optional text search.
// @Tags Adverts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1) example(1)
// @Param per_page query int false "Items per page" default(10) minimum(1) maximum(100) example(10)
// @Param search query string false "Search by title/description" example("promo")
// @Success 200 {object} localization.StandardResponse{data=paginated_advert_response} "Adverts fetched successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid pagination params"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /adverts [get]
func (a *advertAdapter) FetchAdverts(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	list, err := a.advertApplication.FetchAdverts(r.Context(), *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	docs := core.ToAdvertResponses(list.Data)
	res := paginated_advert_response{
		Data: docs,
		Meta: list.Meta,
	}
	localization.SendSuccessResponse(w, localization.SuccessAdvertsFetched, res)
}

// FetchAdvertByID godoc
// @Summary Get advert by ID
// @Description Retrieve a single advert by its identifier.
// @Tags Adverts
// @Accept json
// @Produce json
// @Param id path string true "Advert ID"
// @Success 200 {object} localization.StandardResponse{data=ad_dto.AdvertResponse} "Advert fetched successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Advert not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /adverts/{id} [get]
func (a *advertAdapter) FetchAdvertByID(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
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

// UpdateAdvert godoc
// @Summary Update existing advert (maker)
// @Description Submit an advert update request. Provide only fields to change. Multipart/form-data supported for banner_image.
// @Tags Adverts
// @Accept mpfd
// @Produce json
// @Param id path string true "Advert ID"
// @Param title formData string false "Title" minLength(3) maxLength(20) example("Weekend Promo")
// @Param description formData string false "Description" minLength(30) maxLength(100)
// @Param advert_for formData string false "Advert audience" Enums(IFB,CB,ALL)
// @Param banner_image formData file false "Banner image (<=2MB; jpeg/png/gif/webp)"
// @Success 200 {object} localization.StandardResponse{data=nil} "Advert update request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "No data provided for update / invalid payload"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Advert not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /adverts/{id} [patch]
func (a *advertAdapter) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	req, err := core.ParseAndValidateAdvertRequest(r, true, a.logger)
	if err != nil {
		a.logger.Errorf("[event.UpdateAdvert] failed to parse and validate advert request, id: %s, error: %v", id, err.Error())
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	domainReq, _ := core.ToAdvert(*req)

	if domainReq.Title == "" && domainReq.Description == "" && domainReq.AdvertFor == "" && req.BannerImage == nil {
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

// DeleteAdvert godoc
// @Summary Delete advert (maker)
// @Description Submit an advert delete request.
// @Tags Adverts
// @Accept json
// @Produce json
// @Param id path string true "Advert ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Delete request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Advert not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /advert/{id} [delete]
func (a *advertAdapter) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err = a.advertApplication.DeleteAdvert(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAdvertDeleteRequestSent, nil)
}

// EnableAdvert godoc
// @Summary Enable advert (checker)
// @Description Approve enable request for an advert.
// @Tags Adverts
// @Accept json
// @Produce json
// @Param id path string true "Advert ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Enable request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Advert not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /adverts/{id}/enable [patch]
func (a *advertAdapter) EnableAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err = a.advertApplication.EnableDisableAdvert(r.Context(), id, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertEnableRequestSent, nil)
}

// DisableAdvert godoc
// @Summary Disable advert (checker)
// @Description Approve disable request for an advert.
// @Tags Adverts
// @Accept json
// @Produce json
// @Param id path string true "Advert ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Disable request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Advert not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /adverts/{id}/disable [patch]
func (a *advertAdapter) DisableAdvert(w http.ResponseWriter, r *http.Request) {
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err = a.advertApplication.EnableDisableAdvert(r.Context(), id, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAdvertDisableRequestSent, nil)
}
