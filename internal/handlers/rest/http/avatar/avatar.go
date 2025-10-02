package avatar

import (
	avatarInbound "cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type paginatedAvatarResp types.PaginatedResponse[[]*model.Avatar]
type avatarAdapter struct {
	avatarApplication service.AvatarService
	logger            utils.Logger
}

func InitAvatarAdapter(avatarApplication service.AvatarService, logger utils.Logger) avatarInbound.AvatarInbound {
	return &avatarAdapter{
		logger:            logger,
		avatarApplication: avatarApplication,
	}
}

// CreateAvatar godoc
// @Summary Create avatar (maker)
// @Description Upload an avatar image with a label.
// @Tags Avatars
// @Accept mpfd
// @Produce json
// @Param label formData string true "Label" example("Gold")
// @Param avatar formData file true "Avatar image (<=2MB; jpeg/png/gif)"
// @Success 200 {object} localization.StandardResponse{data=nil} "Avatar create request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid input"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar [post]
func (a *avatarAdapter) CreateAvatar(w http.ResponseWriter, r *http.Request) {
	req, err := ReqFileParse(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	// avatarUrl, err := lib.UploadFileToMinio(r.Context(),a)
	if err := a.avatarApplication.CreateAvatar(r.Context(), &model.Avatar{Label: req.Label}, req.Avatar); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarCreated, nil)
}

// DeleteAvatar godoc
// @Summary Delete avatar (maker)
// @Description Submit delete request for an avatar by ID.
// @Tags Avatars
// @Accept json
// @Produce json
// @Param id path string true "Avatar ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Avatar delete request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Avatar not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar/{id} [delete]
func (a *avatarAdapter) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	if err := a.avatarApplication.DeleteAvatar(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAvatarDeleted, nil)
	a.logger.Infof("Successfuly Delete request sent")

}

// Enable godoc
// @Summary Enable avatar (checker)
// @Description Approve enable request for an avatar.
// @Tags Avatars
// @Accept json
// @Produce json
// @Param id path string true "Avatar ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Avatar enabled"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Avatar not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar/enable/{id} [post]
func (a *avatarAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	if err := a.avatarApplication.EnableDisable(r.Context(), id, true); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarEnabled, nil)
}

// Disable godoc
// @Summary Disable avatar (checker)
// @Description Approve disable request for an avatar.
// @Tags Avatars
// @Accept json
// @Produce json
// @Param id path string true "Avatar ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Avatar disabled"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Avatar not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar/disable/{id} [post]
func (a *avatarAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	if err := a.avatarApplication.EnableDisable(r.Context(), id, false); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarDisabled, nil)
}

// FetchAvatar godoc
// @Summary Get avatar by ID
// @Description Retrieve a single avatar by ID.
// @Tags Avatars
// @Accept json
// @Produce json
// @Param id path string true "Avatar ID"
// @Success 200 {object} localization.StandardResponse{data=model.Avatar} "Avatar retrieved"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid ID"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Avatar not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar/{id} [get]
func (a *avatarAdapter) FetchAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	res, err := a.avatarApplication.FetchAvatarById(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarEnabled, res)
}

// FetchAvatars godoc
// @Summary List avatars
// @Description Retrieve avatars with pagination.
// @Tags Avatars
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1) example(1)
// @Param per_page query int false "Items per page" default(10) minimum(1) maximum(100) example(10)
// @Param search query string false "Search term" example("Gold")
// @Success 200 {object} localization.StandardResponse{data=paginatedAvatarResp} "Avatars retrieved"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar [get]
func (a *avatarAdapter) FetchAvatars(w http.ResponseWriter, r *http.Request) {
	filterParam := local_util.ExtractFilterParams(r)

	avatars, err := a.avatarApplication.FetchAllAvatar(r.Context(), *filterParam)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarRetrieved, avatars)
}

// UpdateAvatar godoc
// @Summary Update avatar (maker)
// @Description Update avatar label and/or image. Provide only fields to change.
// @Tags Avatars
// @Accept mpfd
// @Produce json
// @Param id path string true "Avatar ID"
// @Param label formData string false "Label" example("Silver")
// @Param avatar formData file false "Avatar image (<=2MB; jpeg/png/gif)"
// @Success 200 {object} localization.StandardResponse{data=nil} "Avatar update request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid input"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Avatar not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /avatar/{id} [patch]
func (a *avatarAdapter) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var inputData *multipart.FileHeader
	var label string
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	req, err := ReqFileParse(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if req.Avatar != nil {
		inputData = req.Avatar
		label = req.Label
	}
	if err := a.avatarApplication.UpdateAvatar(r.Context(), id, &model.Avatar{Label: label}, inputData, false); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarUpdated, nil)
}
