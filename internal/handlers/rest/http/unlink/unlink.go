package unlink

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
//
//	@Summary		Get archived user
//	@Description	Get archived users with pagination, filtering, and search. Filterable fields: enabled, kyc_level, is_blocked, is_verified. Searchable fields: full_name, username, user_code, phone_number.
//	@Tags			Unlink
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"
//	@Param			per_page	query		int																	false	"Items per page"
//	@Param			enabled		query		bool																false	"Filter by enabled status"
//	@Param			kyc_level	query		string																false	"Filter by KYC level"
//	@Param			is_blocked	query		bool																false	"Filter by blocked status"
//	@Param			is_verified	query		bool																false	"Filter by verified status"
//	@Param			search		query		string																false	"Search term (searches full_name, username, user_code, phone_number)"
//	@Success		200			{object}	localization.StandardResponse{data=PaginatedArchieveUserResponse}	"User retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}								"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"Internal server error"
//	@Security		BearerAuth
//	@Router			/unlink/archived_user [get]
func (a *unlinkAdapter) GetArchivedUser(w http.ResponseWriter, r *http.Request) {
	// filter parameter
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "unlink", "unlinkAdapter", "GetArchivedUser")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		span.AddEvent("Invalid pagination params", trace.WithAttributes(attribute.Int("page", filterParams.Page), attribute.Int("per_page", filterParams.PerPage)))
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	archivedUser, err := a.unlinkApp.GetAllArchivedUser(ctx, filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		a.logger.Errorf("[GetArchivedUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Archived users retrieved", trace.WithAttributes(attribute.Int("count", len(archivedUser.Data))))
	a.logger.Infof("[GetArchivedUser] retrieved %d archived users", len(archivedUser.Data))
	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, archivedUser)
}

// GetUserByAccount godoc
//
//	@Summary		Get user by account
//	@Description	Get user by account
//	@Tags			Unlink
//	@Accept			json
//	@Produce		json
//	@Param			account_number	path		string											true	"Account number"
//	@Success		200				{object}	localization.StandardResponse{data=model.User}	"User retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}			"Bad request"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}			"User not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}			"Internal server error"
//	@Security		BearerAuth
//	@Router			/unlink/user-by-account/{account_number} [get]
func (a *unlinkAdapter) GetUserByAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "unlink", "unlinkAdapter", "GetUserByAccount")
	defer span.End()
	accNumber := chi.URLParam(r, "account_number")
	if accNumber == "" {
		span.AddEvent("Missing account_number param")
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	user, err := a.unlinkApp.GetUserByAccount(ctx, accNumber)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("account_number", accNumber)))
		a.logger.Errorf("[GetUserByAccount] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("User retrieved", trace.WithAttributes(attribute.String("account_number", accNumber)))
	a.logger.Infof("[GetUserByAccount] user retrieved successfully")
	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, user)
}

// UnlinkUserCif godoc
//
//	@Summary		Unlink user CIF
//	@Description	Unlink user CIF
//	@Tags			Unlink
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"Unlink CIF request sent successfully "
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/unlink/user_cif/{user_code} [patch]
func (a *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter, r *http.Request) {

	ctx, span := local_util.TraceLogger(r.Context(), "handler", "unlink", "unlinkAdapter", "UnlinkUserCif")
	defer span.End()
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		span.AddEvent("Missing user_code param")
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	err := a.unlinkApp.UnlinkUserCif(ctx, userCode)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("user_code", userCode)))
		a.logger.Errorf("[UnlinkUserCif] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("UnlinkUserCif request sent", trace.WithAttributes(attribute.String("user_code", userCode)))
	a.logger.Infof("[UnlinkUserCif] request sent successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessUnlinkCifRequestSent, nil)
}
