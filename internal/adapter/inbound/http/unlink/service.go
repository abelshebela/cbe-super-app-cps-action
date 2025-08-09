package unlink

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant_util "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*any, error)
	GetAllArchivedUser(ctx context.Context, filterParams *constant_util.Filter) (*local_util.PaginatedResponse[*any], error)
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
		fmt.Println("check")
		local_util.SendErrorResponse(w, req.Validate().Error(), 0, nil)
		return
	}

	user, err := ua.unlinkApp.GetUserByAccount(r.Context(), req.AccountNumbers)
	if err != nil {
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	fmt.Println(user)
	local_util.BaseResponseMaker(user, w, "User successfuly retrived", 200)
}

func (ua *unlinkAdapter) GetArchivedUser(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		local_util.SendErrorResponse(w, "PLEASE_ADD_VALID_PAGE_OR_PERPAGE", 0, nil)
		return
	}

	user, err := ua.unlinkApp.GetAllArchivedUser(r.Context(), filterParams)
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
