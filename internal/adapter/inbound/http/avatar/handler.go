package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	avatarAPP "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/avatar"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/avatar"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

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

	file, fileHeader, err := util.ParseMultipartFormFile(r, "avatar", 2<<20)
	if err != nil {
		a.logger.Errorf("error parsing file: %v", err)
		util.SendErrorResponse(w, util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	req.Label = r.FormValue("label")
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

	util.WriteSuccessResponse(w, cpsRes, "Avatar Create Request Created successfully")

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

	util.WriteSuccessResponse(w, cpsAction, "Avatar Delete Request Created successfully")

}

func (a *AvatarHTTPHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	action_code, ok := common_util.GetParam(r, "action_code")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

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
	approvedAction, err := a.avatarHandler.Authorize(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, approvedAction, "AD approved Sucessfully")
}

func (a *AvatarHTTPHandler) Reject(w http.ResponseWriter, r *http.Request) {
	action_code, ok := common_util.GetParam(r, "action_code")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

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

	util.WriteSuccessResponse(w, rejectAction, "Avatar Rejected sucessfully")
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

	util.WriteSuccessResponse(w, cpsAction, "Avatar Disabled Request Create sucessfully")
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

	util.WriteSuccessResponse(w, cpsAction, "Avatar Enable Request Create sucessfully")

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

	util.WriteSuccessResponse(w, avatars, "Avatar list feached sucessfully")

}

func (a *AvatarHTTPHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	avatar, err := a.avatarHandler.GetAvatar(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, avatar, "Avatar list feached sucessfully")

}

func (a *AvatarHTTPHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateAvatar

	file, fileHeader, err := util.ParseMultipartFormFile(r, "avatar", 10<<20)
	if err != nil {
		a.logger.Errorf("error parsing file: %v", err)
		util.SendErrorResponse(w, util.MissingOrInvalidImage, 0, nil)
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

	util.WriteSuccessResponse(w, cpsAction, "Avatar Update Request created sucessfully")

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
