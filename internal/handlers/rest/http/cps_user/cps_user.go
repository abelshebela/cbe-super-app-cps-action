package cps_user

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	localization "cbe-super-app-cps-action/internal/constants/localization"

	"go.opentelemetry.io/otel/attribute"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type handler struct {
	svc    service.CPSUserService
	logger utils.Logger
}

func InitCPSUserHandler(svc service.CPSUserService, logger utils.Logger) *handler {
	return &handler{svc: svc, logger: logger}
}

func (h *handler) ExportUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "exportCpsUsers", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	fileType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("file_type")))
	startDateRaw := strings.TrimSpace(r.URL.Query().Get("created_at_from"))
	endDateRaw := strings.TrimSpace(r.URL.Query().Get("created_at_to"))
	userName := strings.TrimSpace(r.URL.Query().Get("user_name"))

	if fileType == "" || startDateRaw == "" || endDateRaw == "" {
		log.Warnf("[ExportUsers] missing required params: file_type=%q, created_at_from=%q, created_at_to=%q", fileType, startDateRaw, endDateRaw)
		localization.SendBadRequestResponse(w, localization.ErrorRequiredFieldMissing.Message)
		return
	}

	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(startDateRaw, endDateRaw)
	if err != nil {
		log.Warnf("[ExportUsers] invalid date format: created_at_from=%s, created_at_to=%s, err=%v", startDateRaw, endDateRaw, err)
		localization.SendErrorResponse(w, localization.ErrorInvalidDateFormat, nil, nil)
		return
	}

	if endDate.Before(startDate) {
		log.Warnf("[ExportUsers] invalid date range: start=%v, end=%v", startDate, endDate)
		localization.SendErrorResponse(w, localization.ErrorInvalidDateFormat, nil, nil)
		return
	}

	fileLink, err := h.svc.ExportUsers(ctx, startDate, endDate, fileType, userName)
	if err != nil {
		log.Errorf("[ExportUsers] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.CpsUserDataExportedSuccess, fileLink)
}

// CreateUserRequest creates a new CPS user request
//
//	@Summary		Create CPS user request
//	@Description	Creates a new CPS user request that requires approval
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		cpsuser.CreateUserRequest				true	"CPS user creation request"
//	@Success		200		{object}	localization.StandardResponse{data=object}	"CPS user creation request submitted successfully"
//	@Success		201		{object}	localization.StandardResponse{data=object}	"CPS user created successfully (maker only)"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/create [post]
func (h *handler) CreateUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createCpsUserRequest", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var req cpsuser.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateUserRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	// Normalize and validate request
	req.Normalize()
	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateUserRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	formattedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)

	req.PhoneNumber = formattedPhone
	span.SetAttributes(attribute.String("cps_user.phone", formattedPhone))
	if err := h.svc.CreateUserRequest(ctx, req); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[CreateUserRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[CreateUserRequest] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessCpsUserCreated, nil)
		return
	}

	log.Infof("[CreateUserRequest] request sent successfully for user_code")
	localization.SendSuccessResponse(w, localization.SuccessCpsUserCreationRequestSubmitted, nil)
}

// UpdateUserRequest updates an existing CPS user request
//
//	@Summary		Update CPS user request
//	@Description	Updates an existing CPS user request that requires approval
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Param			request		body		cpsuser.UpdateUserRequest				true	"CPS user update request"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user update request submitted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/update/{user_code} [patch]
func (h *handler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateCpsUserRequest", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	var req cpsuser.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateUserRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	req.Normalize()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateUserRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("cps_user.code", userCode))

	if err := h.svc.UpdateUserRequest(ctx, userCode, req); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[UpdateUserRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[UpdateUserRequest] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessCpsUserUpdated, nil)
		return
	}

	log.Infof("[UpdateUserRequest] request sent successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessCpsUserUpdateRequestSubmitted, nil)
}

// FetchUserByUserCode retrieves a CPS user by user code
//
//	@Summary		Get CPS user by user code
//	@Description	Retrieves a specific CPS user by their user code
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"CPS user retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/{user_code} [get]
func (h *handler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchCpsUserByCode", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	userData := local_util.ExtractUserFromContext(ctx)
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	if userCode != userData.UserCode {
		h.logger.Errorf("[FetchUserByUserCode] user code does not match maker code: %s", userCode)
		localization.SendErrorByCodeResponse(w, localization.ErrorUserUnauthorized.Code)
		return
	}

	// Build detailed response in the service layer
	span.SetAttributes(attribute.String("cps_user.code", userCode))
	user, err := h.svc.GetCpsUserDetail(ctx, userCode)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FetchUserByUserCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[FetchUserByUserCode] CPS user retrieved successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessCpsUserRetrieved, user)
}

func (h *handler) FetchUserByUserName(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchCpsUserByCode", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	userCode := strings.TrimSpace(chi.URLParam(r, "user_name"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}
	span.SetAttributes(attribute.String("cps_user.code", userCode))
	user, err := h.svc.GetCpsUserDetailByUserName(ctx, userCode)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FetchUserByUserName] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[FetchUserByUserName] CPS user retrieved successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessCpsUserRetrieved, user)
}

func (h *handler) FetchUserByCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchCpsUserByCode", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	userCode := strings.TrimSpace(chi.URLParam(r, "code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserIDRequired.Code)
		return
	}

	// Build detailed response in the service layer
	span.SetAttributes(attribute.String("cps_user.code", userCode))
	user, err := h.svc.GetCpsUserDetailByCode(ctx, userCode)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FetchUserByCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[FetchUserByCode] CPS user retrieved successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessCpsUserRetrieved, user)
}

// GetAllCPSUsers retrieves all CPS users with pagination
//
//	@Summary		Get all CPS users
//	@Description	Retrieves a paginated list of all CPS users with optional filtering
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int										false	"Page number"		default(1)
//	@Param			per_page	query		int										false	"Items per page"	default(10)
//	@Param			search		query		string									false	"Search terms(full_name,email,username,user_code,phone_number)"
//	@Param			enabled		query		bool									false	"true or false"
//	@Param			department	query		string									false	"CPS user department"
//	@Param			role		query		string									false	"CPS user department"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"CPS users retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users [get]
func (h *handler) GetAllCPSUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllCpsUsers", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	filterParasm, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if filterParasm.Filters["search"] != nil {

		search := strings.TrimSpace(search)
		if phoneNumber, ok := filterParasm.Filters["search"].(string); ok {
			phoneNumber = local_util.FormatPhoneNumber(phoneNumber)
			filterParasm.Filters["search"] = phoneNumber
		}
		filterParasm.Filters["search"] = search
	}

	users, err := h.svc.GetAllCPSUsers(ctx, filterParasm)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAllCPSUsers] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("cps_user.count", len(users.Data)))
	log.Infof("[GetAllCPSUsers] retrieved %d CPS users", len(users.Data))
	localization.SendSuccessResponse(w, localization.SuccessCpsUsersRetrieved, users)
}

// DeleteUserRequest deletes a CPS user request
//
//	@Summary		Delete CPS user request
//	@Description	Deletes a CPS user request that requires approval
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user deleted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/delete/{user_code} [delete]
func (h *handler) DeleteUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteCpsUserRequest", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	span.SetAttributes(attribute.String("cps_user.code", userCode))

	if err := h.svc.DeleteUserRequest(ctx, userCode); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[DeleteUserRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[DeleteUserRequest] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessCPSUserDeleted, nil)
		return
	}

	log.Infof("[DeleteUserRequest] request sent successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SuccessCpsUserDeleted, nil)
}

// DisableUser disables a CPS user
//
//	@Summary		Disable CPS user
//	@Description	Disables a CPS user account
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user disabled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/disable/{user_code} [post]
func (h *handler) DisableUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableCpsUser", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	userData := local_util.ExtractUserFromContext(r.Context())

	if userCode == userData.UserCode {
		log.Errorf("[EnableUser] cannot enable own account")
		localization.SendErrorByCodeResponse(w, localization.ErrorCannotEnableOwnAccount.Code)
		return
	}
	span.SetAttributes(attribute.String("cps_user.code", userCode))

	if err := h.svc.DisableUser(ctx, userCode); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[DisableUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		log.Infof("[DisableUser] CPS user with user_code: %s is successfully disabled and is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessCPSUserDisable, nil)
		return
	}

	log.Infof("[DisableUser] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
	localization.SendSuccessResponse(w, localization.SuccessCpsUserDisabled, nil)
}

// EnableUser enables a CPS user
//
//	@Summary		Enable CPS user
//	@Description	Enables a CPS user account
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user enabled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/enable/{user_code} [post]
func (h *handler) EnableUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableCpsUser", "handler", "cpsUser")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	userData := local_util.ExtractUserFromContext(r.Context())
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		log.Errorf("[EnableUser] user code is required")
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	if userCode == userData.UserCode {
		log.Errorf("[EnableUser] cannot enable own account")
		localization.SendErrorByCodeResponse(w, localization.ErrorCannotDisableOwnAccount.Code)
		return
	}

	span.SetAttributes(attribute.String("cps_user.code", userCode))

	if err := h.svc.EnableUser(ctx, userCode); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[EnableUser] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[EnableUser] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessCPSUserEnabled, nil)
		return
	}

	log.Infof("[EnableUser] request sent successfully for user_code: %s", userCode)
	localization.SendSuccessResponse(w, localization.SucessCpsUserEnabled, nil)
}
