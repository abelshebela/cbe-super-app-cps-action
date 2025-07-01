package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cbe-super-app-cps-action/internal/adapter/outbound/model"
	avatarAPP "cbe-super-app-cps-action/internal/application/avatar"
	"cbe-super-app-cps-action/internal/application/middleware"
	dto "cbe-super-app-cps-action/internal/domain/avatar"
	"cbe-super-app-cps-action/internal/port/inbound/avatar"
	constant "cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	req.Label = r.FormValue("label")

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		a.logger.Errorf("avatar error: %v", err)
		err = fmt.Errorf("failed to read avatar: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid avatar",
		})
		middleware.ErrorHandler(w, err)
		return
	}
	defer file.Close()

	req.Avatar = fileHeader
	var cpsRequest model.CreateCPSAction
	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	cpsRequest.CurrentData = req
	cpsRequest.MakerUser = userData
	cpsRequest.Department = department
	cpsRequest.ActionData = req

	ctx := r.Context()
	cpsRes, err := a.avatarHandler.CreateAvatar(ctx, cpsRequest)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[constant.SuccesResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: constant.SuccesResponse{
			Ok:         true,
			StatusCode: http.StatusOK,
			Data:       cpsRes,
		},
	}

	res.SendJSON()
}

func (a *AvatarHTTPHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.CreateCPSAction

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
}

func (a *AvatarHTTPHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.AuthorizeCPSAction

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	cpsReq.CheckerUser = userData
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	authAction, err := a.avatarHandler.Authorize(ctx, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           authAction,
	}
	res.SendJSON()
}

func (a *AvatarHTTPHandler) Reject(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode avatar request", err)
		err = fmt.Errorf("failed to decode avatar request error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	cpsReq.CheckerUser = userData
	cpsReq.Department = department
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	rejectAction, err := a.avatarHandler.Reject(ctx, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           rejectAction,
	}
	res.SendJSON()
}

func (a *AvatarHTTPHandler) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
}

func (a *AvatarHTTPHandler) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[dto.AvatarResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           avatars,
	}
	res.SendJSON()
}

func (a *AvatarHTTPHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	avatar, err := a.avatarHandler.GetAvatar(ctx, id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*dto.Avatar]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           avatar,
	}

	res.SendJSON()
}

func (a *AvatarHTTPHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateAvatar

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		a.logger.Errorf("avatar error: %v", err)
		err = fmt.Errorf("failed to read avatar: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid avatar",
		})
		middleware.ErrorHandler(w, err)
		return
	}
	defer file.Close()

	req.Avatar = fileHeader
	userData, department, err := a.extractUserFromContext(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	var updateRequest model.CreateCPSAction

	updateRequest.MakerUser = userData
	updateRequest.Department = department
	updateRequest.ActionData = req

	ctx := r.Context()
	cpsAction, err := a.avatarHandler.UpdateAvatar(ctx, id, updateRequest)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
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
