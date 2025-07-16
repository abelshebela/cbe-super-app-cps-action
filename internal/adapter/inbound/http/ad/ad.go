package ad

import (
	"encoding/json"

	// "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"

	inboundAd "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type ADAdapter struct {
	adHandler ad.ADHandlers
	logger    utils.Logger
}

func InitADAdapter(adHandler ad.ADHandlers, logger utils.Logger) inboundAd.ADAdapter {
	return ADAdapter{
		adHandler: adHandler,
		logger:    logger,
	}
}

func (a ADAdapter) CreateOneAdvert(w http.ResponseWriter, r *http.Request) {
	var advertReq dto.CreateAdvertRequest

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	started_at := r.FormValue("started_at")
	started_at_time, err := time.Parse(time.RFC3339, started_at)
	if err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	expired_at := r.FormValue("expired_at")
	expired_at_time, err := time.Parse(time.RFC3339, expired_at)
	if err != nil {
		a.logger.Errorf("failed to parse time", err)
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = dto.AdvertFor(r.FormValue("advert_for"))
	advertReq.Date.StartedAt = started_at_time
	advertReq.Date.ExpiredAt = expired_at_time

	file, fileHeader, err := r.FormFile("banner_image")
	if err != nil {
		a.logger.Errorf("banner image error: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	advertReq.BannerImage = fileHeader

	var cpsReq model.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.ActionData = advertReq

	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	ctx := r.Context()
	cpsActionResponse, err := a.adHandler.CreateOneAdvert(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsActionResponse, "AD Create Request Created successfully")
}

func (a ADAdapter) DeleteOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	ctx := r.Context()

	response, err := a.adHandler.DeleteOneAdvert(ctx, id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	common_util.WriteSuccessResponse(w, response, "AD Delete Request Created successfully")

}
func (a ADAdapter) GetAllAdvert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}
	per_page := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt <= 10 && perPageInt > 0 {
		per_page = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: per_page,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()

	adverts, err := a.adHandler.GetAllAdvert(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, adverts, "AD featch successfully")

}

func (a ADAdapter) GetOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	advert, err := a.adHandler.GetOneAdvert(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(advert)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a ADAdapter) UpdateOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateAdvertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("failed to decode advert request", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	var cpsReq model.CreateCPSAction

	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	req.ID = id
	cpsReq.ActionData = req

	ctx := r.Context()
	cpsActionResponse, err := a.adHandler.UpdateOneAdvert(ctx, id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsActionResponse, "AD Update Request Created successfully")

}

func (a ADAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.AuthorizeCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.CheckerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}

	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	AuthorizeAction, err := a.adHandler.Authorize(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, AuthorizeAction, "AD Approved")

}

func (a ADAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	var cpsReq model.RejectCPSAction
	action_code := chi.URLParam(r, "action_code")

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode advert request", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.CheckerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	rejectedAction, err := a.adHandler.Reject(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, rejectedAction, "AD REJECTED")

}
