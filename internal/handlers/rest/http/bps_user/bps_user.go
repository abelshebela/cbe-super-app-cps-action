package bpsmakerhandler

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	// model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	bps_user_dto "cbe-super-app-cps-action/internal/constants/dto/bps_user"
)

type paginated_resp *types.PaginatedResponse[[]*model.BPSUser]
type bps_user_resp *model.BPSUser

type BPSUserHandler struct {
	Service service.BPSUserService
	logger  utils.Logger
}

func InitBPSUserMakerHandler(service service.BPSUserService, logger utils.Logger) bps_user.BPSUserHandler {
	return BPSUserHandler{
		Service: service,
		logger:  logger,
	}
}

// FetchUserByUserCode retrieves a BPS user by user code
//
//	@Summary		Get BPS user by user code
//	@Description	Retrieves a specific BPS user by their user code
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string												true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=bps_user_resp}	"BPS user retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}				"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}				"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/bps_users/{user_code} [get]
func (h BPSUserHandler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "fetchBpsUserByCode", "handler", "bpsUser")
	defer span.End()
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bps_user.code", userCode))

	user, err := h.Service.FetchUserByUserCode(ctx, userCode)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[FetchUserByUserCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[FetchUserByUserCode] BPS user retrieved successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, user)
}

// GetAllBPSUsers retrieves all BPS users with pagination
//
//	@Summary		Get all BPS users
//	@Description	Retrieves a paginated list of all BPS users with optional filtering
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int	false	"Page number"		default(1)
//	@Param			per_page	query	int	false	"Items per page"	default(10)

// @Param		branch_code			query		string												false	"Branch code filter"
// @Param		branch_name			query		string												false	"Branch name filter"
// @Param		enabled				query		bool												false	"Enabled status filter"
// @Param		role				query		string												false	"Role filter"
// @Param		first_password_set	query		bool												false	"First password set filter"
// @Param		full_name			query		string												false	"Full name filter"
// @Param		user_code			query		string												false	"User code filter"
// @Param		phone_number		query		string												false	"Phone number filter"
// @Param		username			query		string												false	"Username filter"
// @Param		search				query		string												false	"searchable fieldes (full_name,username,user_code,phone_number)"
// @Success	200					{object}	localization.StandardResponse{data=paginated_resp}	"BPS users retrieved successfully"
// @Failure	400					{object}	localization.StandardResponse{data=nil}				"Bad request"
// @Failure	500					{object}	localization.StandardResponse{data=nil}				"Internal server error"
// @Security	BearerAuth
// @Router		/bps_users/ [get]
func (h BPSUserHandler) GetAllBPSUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getAllBpsUsers", "handler", "bpsUser")
	defer span.End()
	filterParams := common_utils.ExtractFilterParams(r)

	users, err := h.Service.GetAllBPSUsers(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAllBPSUsers] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("bps_user.count", len(users.Data)))
	h.logger.Infof("[GetAllBPSUsers] retrieved %d BPS users", len(users.Data))
	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, users)
}

// DisableUser disables a BPS user
//
//	@Summary		Disable BPS user
//	@Description	Disables a BPS user account
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string													true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=map[string]string}	"BPS user disabled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}					"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}					"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/bps_users/disable/{user_code} [post]
func (h BPSUserHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableBpsUser", "handler", "bpsUser")
	defer span.End()
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bps_user.code", userCode))

	err := h.Service.UpdateBpsUser(ctx, userCode, false)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[DisableUser] request sent successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessBpsUserDisableRequestSent, map[string]string{})
}

// EnableUser enables a BPS user
//
//	@Summary		Enable BPS user
//	@Description	Enables a BPS user account
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string													true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=map[string]string}	"BPS user enabled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}					"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}					"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/bps_users/enable/{user_code} [post]
func (h BPSUserHandler) EnableUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableBpsUser", "handler", "bpsUser")
	defer span.End()
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bps_user.code", userCode))

	err := h.Service.UpdateBpsUser(ctx, userCode, true)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[EnableUser] request sent successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessBpsUserEnableRequestSent, map[string]string{})
}

// CreateBPSUser creates a new BPS user
//
//	@Summary		Create BPS user
//	@Description	Creates a new BPS user account
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		bps_user_dto.BPSUserCreateRequest						true	"BPS user create request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}					"BPS user created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/bps_users/ [post]
func (h BPSUserHandler) CreateBPSUser(w http.ResponseWriter, r *http.Request) {
	var req bps_user_dto.BPSUserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, "Invalid request payload")
		return
	}

	// Validate required fields
	if req.UserID == "" {
		localization.SendBadRequestResponse(w, "user id requed")
		return
	}

	concatinated_name := req.FullName.FirstName + " " + req.FullName.MiddleName + " " + req.FullName.LastName
	now := time.Now()
	NewUser := model.BPSUser{
		// ID:       primitive.NewObjectID(),
		UserCode: req.UserID,
		FullName: concatinated_name,
		// JobTitle:    req.JobTitle,
		// UserName:    req.ImpowerID,
		PhoneNumber: req.PhoneNumber,
		// Email:       req.Email,
		CreatedAt:      now,
		LastModifiedAt: now,

		Enabled: false,
	}

	err := h.Service.CreateBPSUser(r.Context(), NewUser)
	if err != nil {
		// span.RecordError(err)
		h.logger.Errorf("[CreateUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSUserCreated, nil)
}
