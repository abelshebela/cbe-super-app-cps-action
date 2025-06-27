package ad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/ad"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/entity"
	inboundAd "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/ad"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	started_at := r.FormValue("started_at")
	started_at_time, err := time.Parse(time.RFC3339, started_at)
	if err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	expired_at := r.FormValue("expired_at")
	expired_at_time, err := time.Parse(time.RFC3339, expired_at)
	if err != nil {
		a.logger.Errorf("failed to parse time", err)
		err = fmt.Errorf("failed to check advert bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = ad.AdvertFor(r.FormValue("advert_for"))
	advertReq.Date.StartedAt = started_at_time
	advertReq.Date.ExpiredAt = expired_at_time

	file, fileHeader, err := r.FormFile("banner_image")
	if err != nil {
		a.logger.Errorf("banner image error: %v", err)
		err = fmt.Errorf("failed to read banner image: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid banner image",
		})
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*ad.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsActionRes,
	}

	res.SendJSON()
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: struct {
			OK      bool   `json:"ok"`
			Message string `json:"message"`
		}{
			OK:      true,
			Message: "Request Sent Successfully",
		},
	}

	res.SendJSON()
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.AdvertResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           adverts,
	}
	res.SendJSON()
}

func (a ADAdapter) GetOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	advert, err := a.adHandler.GetOneAdvert(ctx, id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*ad.AdvertResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           advert,
	}

	res.SendJSON()
}

func (a ADAdapter) UpdateOneAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req ad.UpdateAdvertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("failed to decode advert request", err)
		err = fmt.Errorf("failed to decode advert request error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
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
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	authAction, err := a.adHandler.Authorize(ctx, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           authAction,
	}
	res.SendJSON()
}

func (a ADAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	var cpsReq entity.CPSAction
	action_code := chi.URLParam(r, "action_code")

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode advert request", err)
		err = fmt.Errorf("failed to decode advert request error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           rejectAction,
	}
	res.SendJSON()
}
