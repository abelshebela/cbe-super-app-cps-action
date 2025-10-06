package unlink

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PaginatedArchieveUserResponse types.PaginatedResponse[[]*model.ArchivedUser]
type unlinkAdapter struct {
	logger    utils.Logger
	unlinkApp service.UnlinkService
}

func InitUnlinkAdapter(unlinkApp service.UnlinkService, logger utils.Logger) inbound.UnlinkAdapter {
	return &unlinkAdapter{
		logger:    logger,
		unlinkApp: unlinkApp,
	}
}

// GetArchivedUser godoc
// @Summary Get archived user
// @Description Get archived user
// @Tags Unlink
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=PaginatedArchieveUserResponse} "User retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /unlink/archived_user [get]
func (a *unlinkAdapter) GetArchivedUser(w http.ResponseWriter, r *http.Request) {
	// filter parameter
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	archivedUser, err := a.unlinkApp.GetAllArchivedUser(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, archivedUser)
}

// GetUserByAccount godoc
// @Summary Get user by account
// @Description Get user by account
// @Tags Unlink
// @Accept json
// @Produce json
// @Param account_number path string true "Account number"
// @Success 200 {object} localization.StandardResponse{data=model.User} "User retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "User not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /unlink/user-by-account/{account_number} [get]
func (a *unlinkAdapter) GetUserByAccount(w http.ResponseWriter, r *http.Request) {
	accNumber := chi.URLParam(r, "account_number")
	if accNumber == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	user, err := a.unlinkApp.GetUserByAccount(r.Context(), accNumber)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, user)
}

// UnlinkUserCif godoc
// @Summary Unlink user CIF
// @Description Unlink user CIF
// @Tags Unlink
// @Accept json
// @Produce json
// @Param user_code path string true "User code"
// @Success 200 {object} localization.StandardResponse{data=nil} "Unlink CIF request sent successfully "
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "User not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /unlink/user_cif/{user_code} [patch]
func (a *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter, r *http.Request) {

	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	err := a.unlinkApp.UnlinkUserCif(r.Context(), userCode)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessUnlinkCifRequestSent, nil)
}
