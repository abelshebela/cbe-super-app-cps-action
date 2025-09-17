package avatar

import (
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/avatar"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarHTTPHandler struct {
	Application avatar.AvatarApplicationService
	logger      utils.Logger
}

func InitAvatarHTTPHandler(app avatar.AvatarApplicationService, logger utils.Logger) *AvatarHTTPHandler {
	return &AvatarHTTPHandler{
		Application: app,
		logger:      logger,
	}
}

func (h *AvatarHTTPHandler) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (entities.User, bool) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("Incomplete user context")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusUnauthorized, nil)
		return entities.User{}, false
	}
	maker := entities.User{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
	return maker, true
}

func (h *AvatarHTTPHandler) extractID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, ok := common_util.GetParam(r, "id")
	if !ok || id == "" {
		h.logger.Errorf("Missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return "", false
	}
	return id, true
}

func (h *AvatarHTTPHandler) parseAndValidateAvatarRequest(w http.ResponseWriter, r *http.Request, isCreate bool) (AvatarDTO, bool) {
	req, err := ParseAvatarRequestFromMultipartForm(r, isCreate)
	if err != nil {
		h.logger.Errorf("[AvatarHTTPHandler.parseAndValidateAvatarRequest] Failed to parse avatar request: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return AvatarDTO{}, false
	}

	if err := req.Validate(isCreate); err != nil {
		h.logger.Errorf("[AvatarHTTPHandler.parseAndValidateAvatarRequest] Validation failed: %v", err)
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return AvatarDTO{}, false
	}
	return req, true
}

func (h *AvatarHTTPHandler) sendResponse(w http.ResponseWriter, err error, successMessage string, data interface{}) {
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	common_util.WriteSuccessResponse(w, data, successMessage)
}

func (h *AvatarHTTPHandler) CreateAvatar(w http.ResponseWriter, r *http.Request) {
	req, ok := h.parseAndValidateAvatarRequest(w, r, true)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	domainReq := ToAvatarRequestFromDTO(&req)
	err := h.Application.CreateAvatar(r.Context(), *domainReq, maker)
	h.sendResponse(w, err, "Avatar creation request submitted successfully", nil)
}

func (h *AvatarHTTPHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	req, ok := h.parseAndValidateAvatarRequest(w, r, false)
	if !ok {
		return
	}

	if req.IsEmpty() {
		h.logger.Errorf("[AvatarHTTPHandler.UpdateAvatar] No data provided for update")
		common_util.SendErrorResponse(w, common_util.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	domainReq := ToAvatarRequestFromDTO(&req)
	err := h.Application.UpdateAvatar(r.Context(), id, *domainReq, maker)
	h.sendResponse(w, err, "Avatar update request submitted successfully", nil)
}

func (h *AvatarHTTPHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.DeleteAvatar(r.Context(), id, maker)
	h.sendResponse(w, err, "Avatar delete request submitted successfully", nil)
}

func (h *AvatarHTTPHandler) Enable(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableOrDisableAvatar(r.Context(), id, maker, true)
	h.sendResponse(w, err, "Event enable request submitted successfully", nil)
}

func (h *AvatarHTTPHandler) Disable(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableOrDisableAvatar(r.Context(), id, maker, false)
	h.sendResponse(w, err, "Event disable request submitted successfully", nil)
}

func (h *AvatarHTTPHandler) FetchAvatar(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	data, err := h.Application.GetAvatar(r.Context(), id)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	res := ToAvatarResponseDTO(data)
	h.sendResponse(w, nil, "Avatar successfully retrieved", res)
}

func (h *AvatarHTTPHandler) FetchAvatars(w http.ResponseWriter, r *http.Request) {
	filterParam := common_util.ExtractFilterParams(r)

	list, err := h.Application.GetAllAvatar(r.Context(), filterParam)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	docs := []AvatarResponseDTO{}
	for _, doc := range list.Data {
		res := ToAvatarResponseDTO(doc)
		docs = append(docs, *res)
	}

	res := common_util.PaginatedResponse[*[]AvatarResponseDTO]{
		Data: &docs,
		Meta: list.Meta,
	}
	h.sendResponse(w, nil, "Avatars successfully retrieved", res)
}
