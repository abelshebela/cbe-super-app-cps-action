package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	avatarAPP "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/avatar"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/avatar"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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

	cpsRes, err := a.avatarHandler.CreateAvatar(r.Context(), cpsRequest)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, cpsRes, "Avatar Create Request Created successfully")

}

func (a *AvatarHTTPHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

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

	cpsAction, err := a.avatarHandler.DeleteAvatar(r.Context(), id, cpsReq)
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

	approvedAction, err := a.avatarHandler.Authorize(r.Context(), cpsReq)
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

	rejectAction, err := a.avatarHandler.Reject(r.Context(), cpsReq)
	if err != nil {
		a.logger.Errorf("failed to decode avatar request", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, rejectAction, "Avatar Rejected sucessfully")
}

func (a *AvatarHTTPHandler) Disable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
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

	cpsAction, err := a.avatarHandler.EnableOrDisableAvatar(r.Context(), id, model.RequestDisableAvatar, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, cpsAction, "Avatar Disabled Request Create sucessfully")
}

func (a *AvatarHTTPHandler) Enable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

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

	cpsAction, err := a.avatarHandler.EnableOrDisableAvatar(r.Context(), id, model.RequestEnableAvatar, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, cpsAction, "Avatar Enable Request Create sucessfully")

}

func (a *AvatarHTTPHandler) GetAllAvatar(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	data, err := a.avatarHandler.GetAllAvatar(r.Context(), *filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	util.WriteSuccessResponse(w, data, "Avatar list feached sucessfully")

}

func (a *AvatarHTTPHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	avatar, err := a.avatarHandler.GetAvatar(r.Context(), id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, avatar, "Avatar list feached sucessfully")

}

func (a *AvatarHTTPHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	fmt.Println(id, ok, "id and ok")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

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

	cpsAction, err := a.avatarHandler.UpdateAvatar(r.Context(), id, updateRequest)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, cpsAction, "Avatar Update Request created sucessfully")

}

func (a *AvatarHTTPHandler) extractUserFromContext(r *http.Request) (model.User, string, error) {
	userContext := ctx_util.ExtractUserContext(r)

	if userContext.IsIncomplete() {
		return model.User{}, "", fmt.Errorf(common_util.IncompleteUserInfo)
	}

	return model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}, userContext.Department, nil
}
