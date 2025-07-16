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

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	inboundAd "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
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

	file, fileHeader, err := util.ParseMultipartFormFile(r, "banner_image", 10<<20)
	if err != nil {
		a.logger.Errorf("error parsing file: %v", err)
		util.SendErrorResponse(w, util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	startedAt := r.FormValue("started_at")
	startedAtTime, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	expiredAt := r.FormValue("expired_at")
	expiredAtTime, err := time.Parse(time.RFC3339, expiredAt)
	if err != nil {
		a.logger.Errorf("failed to parse time", err)
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = dto.AdvertFor(r.FormValue("advert_for"))
	advertReq.Date.StartedAt = startedAtTime
	advertReq.Date.ExpiredAt = expiredAtTime

	advertReq.BannerImage = fileHeader

	var cpsReq model.CreateCPSAction

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq.ActionData = advertReq

	cpsReq.MakerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department

	ctx := r.Context()
	cpsActionResponse, err := a.adHandler.CreateOneAdvert(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	util.WriteSuccessResponse(w, cpsActionResponse, "AD Create Request Created successfully")
}

func (a ADAdapter) DeleteOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.CreateCPSAction

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}
	cpsReq.MakerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department
	ctx := r.Context()

	response, err := a.adHandler.DeleteOneAdvert(ctx, id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	util.WriteSuccessResponse(w, response, "AD Delete Request Created successfully")

}
func (a ADAdapter) GetAllAdvert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}
	perPage := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt <= 10 && perPageInt > 0 {
		perPage = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()

	adverts, err := a.adHandler.GetAllAdvert(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, adverts, "AD featch successfully")

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

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	var cpsReq model.CreateCPSAction

	cpsReq.MakerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department
	req.ID = id
	cpsReq.ActionData = req

	ctx := r.Context()
	cpsActionResponse, err := a.adHandler.UpdateOneAdvert(ctx, id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, cpsActionResponse, "AD Update Request Created successfully")

}

func (a ADAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")

	var cpsReq model.AuthorizeCPSAction

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq.CheckerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}

	cpsReq.Department = userContext.Department
	cpsReq.ActionCode = actionCode

	ctx := r.Context()
	AuthorizeAction, err := a.adHandler.Authorize(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, AuthorizeAction, "AD Approved")

}

func (a ADAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	var cpsReq model.RejectCPSAction
	actionCode := chi.URLParam(r, "action_code")

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode advert request", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq.CheckerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department
	cpsReq.ActionCode = actionCode

	ctx := r.Context()
	rejectedAction, err := a.adHandler.Reject(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, rejectedAction, "AD REJECTED")

}
