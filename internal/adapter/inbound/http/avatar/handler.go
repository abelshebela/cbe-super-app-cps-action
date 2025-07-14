package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	avatarAPP "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/avatar"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/avatar"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	util_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type AvatarHTTPHandler struct {
	avatarHandler avatarAPP.AvatarApplicationService
	logger        utils.Logger
}

func InitAvatarHTTPHandler(avatarHandler avatarAPP.AvatarApplicationService, logger utils.Logger) avatar.AvatarInbound {
	return &AvatarHTTPHandler{
		avatarHandler: avatarHandler,
		logger:        logger,
	}
}

func (a *AvatarHTTPHandler) CreateAvatar(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAvatar

	if err := r.ParseMultipartForm(2 << 20); err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	req.Label = r.FormValue("label")

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		a.logger.Errorf("avatar error: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)

		return
	}
	defer file.Close()

	req.Avatar = fileHeader
	var cpsRequest model.CreateCPSAction
	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	cpsRequest.CurrentData = req
	cpsRequest.MakerUser = userData
	cpsRequest.Department = department
	cpsRequest.ActionData = req

	ctx := r.Context()
	cpsRes, err := a.avatarHandler.CreateAvatar(ctx, cpsRequest)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsRes)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)

}

func (a *AvatarHTTPHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.CreateCPSAction

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		a.logger.Errorf("failed to extract department", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	cpsReq.ActionData = dto.Avatar{
		ID: id,
	}
	cpsReq.MakerUser = userData
	cpsReq.Department = department

	ctx := r.Context()
	cpsAction, err := a.avatarHandler.DeleteAvatar(ctx, id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)

}

func (a *AvatarHTTPHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.AuthorizeCPSAction

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	cpsReq.CheckerUser = userData
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	authAction, err := a.avatarHandler.Authorize(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(authAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) Reject(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode avatar request", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	cpsReq.CheckerUser = userData
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	rejectAction, err := a.avatarHandler.Reject(ctx, cpsReq)
	if err != nil {
		a.logger.Errorf("failed to decode avatar request", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(rejectAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	var cpsReq model.CreateCPSAction

	cpsReq.MakerUser = userData
	cpsReq.Department = department
	cpsReq.ActionData = dto.Avatar{
		ID: id,
	}

	ctx := r.Context()
	cpsAction, err := a.avatarHandler.EnableOrDisableAvatar(ctx, id, model.RequestDisableAvatar, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	var cpsReq model.CreateCPSAction

	cpsReq.MakerUser = userData
	cpsReq.Department = department
	cpsReq.ActionData = dto.Avatar{
		ID: id,
	}

	ctx := r.Context()
	cpsAction, err := a.avatarHandler.EnableOrDisableAvatar(ctx, id, model.RequestEnableAvatar, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) GetAllAvatar(w http.ResponseWriter, r *http.Request) {
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

	filterParams := constant.Filter{
		Page:    page,
		PerPage: per_page,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()

	avatars, err := a.avatarHandler.GetAllAvatar(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(avatars)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	avatar, err := a.avatarHandler.GetAvatar(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(avatar)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateAvatar

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		a.logger.Errorf("avatar error: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)

		return
	}
	defer file.Close()

	req.Avatar = fileHeader
	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	var updateRequest model.CreateCPSAction

	updateRequest.MakerUser = userData
	updateRequest.Department = department
	updateRequest.ActionData = req

	ctx := r.Context()
	cpsAction, err := a.avatarHandler.UpdateAvatar(ctx, id, updateRequest)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (a *AvatarHTTPHandler) extractUserFromContext(r *http.Request) (model.User, string, error) {
	userCode, ok := r.Context().Value(constant.ContextKey("user_code")).(string)
	if !ok {
		a.logger.Errorf("failed to get user code from context")
		return model.User{}, "", fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	fullName, ok := r.Context().Value(constant.ContextKey("full_name")).(string)
	if !ok {
		a.logger.Errorf("failed to get full name from context")
		return model.User{}, "", fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	phoneNumber, ok := r.Context().Value(constant.ContextKey("phone_number")).(string)
	if !ok {
		a.logger.Errorf("failed to get phone number from context")
		return model.User{}, "", fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	department, ok := r.Context().Value(constant.ContextKey("department")).(string)
	if !ok {
		a.logger.Errorf("failed to get department from context")
		return model.User{}, "", fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		})
	}

	return model.User{
		UserCode:    userCode,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}, department, nil
}
