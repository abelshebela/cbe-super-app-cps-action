package avatar

import (
	avatarInbound "cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

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
	// req,err := ReqFileParse(r)
	// if err != nil {
	// 	localization.SendErrorResponse(w,err.Error())
	// 	return
	// }

	// if err := a.avatarApplication.CreateAvatar(r.Context(),&model.Avatar{Avatar: req.Avatar,Label: req.Label}); err != nil {
	// 	localization.SendErrorByCodeResponse(w,err.Error())
	// 	return
	// }

	// localization.SendSuccessResponse(w,localization.SuccessAvatarCreated,nil,nil)
}
func (a *avatarAdapter) DeleteAvatar(w http.ResponseWriter, r *http.Request) {}
func (a *avatarAdapter) Enable(w http.ResponseWriter, r *http.Request)       {}
func (a *avatarAdapter) Disable(w http.ResponseWriter, r *http.Request)      {}
func (a *avatarAdapter) FetchAvatar(w http.ResponseWriter, r *http.Request)  {}
func (a *avatarAdapter) FetchAvatars(w http.ResponseWriter, r *http.Request) {}
func (a *avatarAdapter) UpdateAvatar(w http.ResponseWriter, r *http.Request) {}
