package avatar

import (
	avatarInbound "cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

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

func (a *avatarAdapter) CreateAvatar(w http.ResponseWriter, r *http.Request) {
	req, err := ReqFileParse(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// avatarUrl, err := lib.UploadFileToMinio(r.Context(),a)
	if err := a.avatarApplication.CreateAvatar(r.Context(), &model.Avatar{Label: req.Label}, req.Avatar); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarCreated, nil)
}
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
}
func (a *avatarAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	if err := a.avatarApplication.UpdateAvatar(r.Context(), id, &model.Avatar{Enable: true}); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarEnabled, nil)
}
func (a *avatarAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	if err := a.avatarApplication.UpdateAvatar(r.Context(), id, &model.Avatar{Enable: false}); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarDisabled, nil)
}
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
func (a *avatarAdapter) FetchAvatars(w http.ResponseWriter, r *http.Request) {
	filterParam := local_util.ExtractFilterParams(r)

	avatars, err := a.avatarApplication.FetchAllAvatar(r.Context(), *filterParam)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarRetrieved, avatars)
}
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
	if err := a.avatarApplication.UpdateAvatar(r.Context(), id, &model.Avatar{Label: label}, inputData); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAvatarUpdated, nil)
}
