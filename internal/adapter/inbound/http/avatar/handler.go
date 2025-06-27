package avatar

import (
	"fmt"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	avatarAPP "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	dto "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/avatar"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
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
	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsRequest.CurrentData = req
	cpsRequest.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
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
