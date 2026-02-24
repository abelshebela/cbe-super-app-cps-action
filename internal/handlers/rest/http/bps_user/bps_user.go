package bpsmakerhandler

import (
	constants "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	model "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_utils "cbe-super-app-cps-action/pkgs/utils"

	// model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	bps_user_dto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
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
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bps_user.code", userCode))

	user, err := h.Service.FetchUserByUserCode(ctx, userCode)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FetchUserByUserCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[FetchUserByUserCode] BPS user retrieved successfully for user_code: %s", userCode)
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
// @Success	200					{object}	localization.StandardResponse{data=object}	"BPS users retrieved successfully"
// @Failure	400					{object}	localization.StandardResponse{data=nil}				"Bad request"
// @Failure	500					{object}	localization.StandardResponse{data=nil}				"Internal server error"
// @Security	BearerAuth
// @Router		/bps_users/ [get]
func (h BPSUserHandler) GetAllBPSUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getAllBpsUsers", "handler", "bpsUser")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	filterParams := common_utils.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := common_utils.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := common_utils.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	users, err := h.Service.GetAllBPSUsers(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAllBPSUsers] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("bps_user.count", len(users.Data)))
	log.Infof("[GetAllBPSUsers] retrieved %d BPS users", len(users.Data))
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
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bps_user.code", userCode))

	err := h.Service.UpdateStatusBpsUser(ctx, userCode, false)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DisableUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		log.Infof("[DisableUser] request sent successfully for user_code: %s", userCode)
		localization.SendSuccessResponse(w, localization.SuccessBpsUserDisableRequestSentSP, map[string]string{})
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBpsUserDisableRequestSent, map[string]string{})

	}
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
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bps_user.code", userCode))

	err := h.Service.UpdateStatusBpsUser(ctx, userCode, true)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[EnableUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		log.Infof("[EnableUser] request sent successfully for user_code: %s", userCode)
		localization.SendSuccessResponse(w, localization.SuccessBpsUserEnableRequestSentSP, map[string]string{})
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBpsUserEnableRequestSent, map[string]string{})
	}
}

// CreateBPSUser creates a new BPS user
//
//	@Summary		Create BPS user
//	@Description	Creates a new BPS user account
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		bpsuser.BPSUserCreateRequest	true	"BPS user create request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}					"BPS user creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/bps_users/ [post]
func (h BPSUserHandler) CreateBPSUser(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	var req bps_user_dto.BPSUserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	// Validate required fields
	if req.UserID == "" {
		localization.SendBadRequestResponse(w, "user id requed")
		return
	}
	userCodeGenerated := local_utils.RandomGenerator(8)

	NewUser := bps_model.BPSUser{
		Username:         req.UserID,
		FullName:         req.FullName,
		JobTitle:         req.JobTitle,
		UserCode:         userCodeGenerated,
		PhoneNumber:      req.PhoneNumber,
		BranchCode:       req.BranchCode,
		Email:            req.Email,
		FirstPasswordSet: true,
		IsFirstTimeLogin: true,
		Enabled:          true,
	}

	err := h.Service.CreateBPSUser(ctx, NewUser)
	if err != nil {
		// span.RecordError(err)
		log.Errorf("[CreateUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessBPSUserCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBPSUserCreated, nil)

	}
}

// UpdateBPSUser updates an existing BPS user
//
//	@Summary		Update BPS user
//	@Description	Updates an existing BPS user account
//	@Tags			BPS Users
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string											true	"BPS User ID"
//	@Param			request		body		bpsuser.BPSUserUpdateRequest	true	"BPS user update request"
//	@Success		200			{object}	localization.StandardResponse{data=nil}			"BPS user update request submitted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}			"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}			"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}			"Internal server error"
//	@Security		BearerAuth
//	@Router			/bps_users/{id} [patch]
func (h BPSUserHandler) UpdateBPSUser(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, "user code required")
		return
	}

	var req bps_user_dto.BPSUserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		localization.SendBadRequestResponse(w, "unable to convert object id")
		return
	}

	updatedUser := bps_model.BPSUser{
		ID:             objID,
		LastModifiedAt: time.Now(),
	}

	// Populate fields if not empty
	if req.FullName != nil {
		updatedUser.FullName = *req.FullName
	}
	if req.PhoneNumber != nil {
		updatedUser.PhoneNumber = *req.PhoneNumber
	}
	if req.JobTitle != nil {
		updatedUser.JobTitle = *req.JobTitle
	}
	if req.UserID != nil {
		updatedUser.Username = *req.UserID
	}
	if req.Email != nil {
		updatedUser.Email = *req.Email
	}
	// if req.Role != "" {
	// 	updatedUser.Role = req.Role
	// }
	if len(req.BranchCode) > 0 {
		updatedUser.BranchCode = req.BranchCode
	}
	// if req.HomeBranch != "" {
	// 	updatedUser.HomeBranch = req.HomeBranch
	// }
	// if req.Realm != "" {
	// 	updatedUser.Realm = req.Realm
	// }
	// Always update enabled (bool, so default is false if not set)
	// updatedUser.Enabled = req.Enabled

	err = h.Service.UpdateBPSUser(ctx, id, updatedUser)
	if err != nil {
		log.Errorf("[UpdateBPSUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessBPSUserUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBPSUserUpdated, nil)
	}
}
