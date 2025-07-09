package ad

import (
	"encoding/json"
	"fmt"

	// "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	inboundAd "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
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
	var advertReq ad.CreateAdvertRequest

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	started_at := r.FormValue("started_at")
	fmt.Println("==========================", r.FormValue("description"))
	started_at_time, err := time.Parse(time.RFC3339, started_at)
	if err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	fmt.Println("==========FormValue==============")
	fmt.Println(r.FormValue("title"))

	expired_at := r.FormValue("expired_at")
	expired_at_time, err := time.Parse(time.RFC3339, expired_at)
	if err != nil {
		a.logger.Errorf("failed to parse time", err)
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	advertReq.Title = r.FormValue("title")
	fmt.Printf("title %v", advertReq.Title)
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = ad.AdvertFor(r.FormValue("advert_for"))
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

	var cpsReq ad.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.ActionData = advertReq
	cpsReq.MakerUser = ad.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	ctx := r.Context()
	cpsActionRes, err := a.adHandler.CreateOneAdvert(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsActionRes)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}
func (a ADAdapter) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: status,
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
func (a ADAdapter) DeleteOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq ad.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.ActionData.ID = id
	cpsReq.MakerUser = ad.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	ctx := r.Context()
	if err := a.adHandler.DeleteOneAdvert(ctx, cpsReq); err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	util.BaseResponseMaker(nil, w, def.Message, http.StatusAccepted)
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

	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(adverts)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
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

	var req ad.UpdateAdvertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("failed to decode advert request", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	var cpsReq ad.CreateCPSAction

	cpsReq.MakerUser = ad.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	cpsReq.ActionData = ad.CreateAdvertRequest{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		AdvertFor:   req.AdvertFor,
		Date:        req.Date,
	}

	ctx := r.Context()
	cpsAction, err := a.adHandler.UpdateOneAdvert(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a ADAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq entity.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.CheckerUser = entity.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	fmt.Println("===========================================================")
	fmt.Println(cpsReq.CheckerUser)
	fmt.Println("===========================================================")

	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	authAction, err := a.adHandler.Authorize(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(authAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a ADAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	var cpsReq entity.CPSAction
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

	cpsReq.CheckerUser = entity.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	rejectAction, err := a.adHandler.Reject(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(rejectAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}
