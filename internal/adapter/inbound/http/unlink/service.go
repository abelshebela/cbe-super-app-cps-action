package unlink

import (
	"context"
	"encoding/json"
	"net/http"

	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
}

type unlinkAdapter struct {
	logger    utils.Logger
	unlinkApp UnlinkAccount
}

func InitAdapterUnlinkService(application UnlinkAccount, logger utils.Logger) Inbound.UnlinkHandler {
	return &unlinkAdapter{
		logger:    logger,
		unlinkApp: application,
	}
}

func (ua *unlinkAdapter) GetUserByAccount(w http.ResponseWriter, r *http.Request) {
	var req GetUserByAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_util.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}
	if req.Validate() != nil {
		local_util.SendErrorResponse(w, req.Validate().Error(), 0, nil)
		return
	}
	user, err := ua.unlinkApp.GetUserByAccount(r.Context(), req.AccountNumbers)
	if err != nil {
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(user, w, "User successfuly retrived", 200)
}

func (ua *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter, r *http.Request) {
	var req UnlinkUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_util.SendErrorResponse(w, "UNHANDLE_SERVER_ERROR", 0, nil)
		return
	}
	if req.Validate() != nil {
		local_util.SendErrorResponse(w, req.Validate().Error(), 0, nil)
		return
	}
	err := ua.unlinkApp.UnlinkUserCif(r.Context(), req.UserCode)
	if err != nil {
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(map[string]interface{}{}, w, "Unlink cif request is sent successfuly", 200)
}
