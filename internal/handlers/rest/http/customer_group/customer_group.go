package customergroup

import (
	"context"
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	customer_group_dto "cbe-super-app-cps-action/internal/constants/dto/customer_group"
	cg_iface "cbe-super-app-cps-action/internal/constants/interfaces/customer_group"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CustomerGroupAdapter struct {
	svc    service.CustomerGroupService
	logger utils.Logger
}

func NewCustomerGroupHandler(svc service.CustomerGroupService, logger utils.Logger) cg_iface.CustomerGroup {
	return &CustomerGroupAdapter{svc: svc, logger: logger}
}

// Create Customer Group
//
//	@Summary		Create Customer Group
//	@Description	Creates a new customer group segment request
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		customer_group_dto.CreateSegmentRequest	true	"Customer Group DTO"
//	@Success		201		{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups [post]
func (h *CustomerGroupAdapter) CreateCustomerGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "CreateCustomerGroup", "handler", "customerGroup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var req customer_group_dto.CreateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[CreateCustomerGroup] decode err: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}
	if err := req.Validate(); err != nil {
		log.Errorf("[CreateCustomerGroup] validation err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	req.Normalize()

	if err := h.svc.Create(ctx, req); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateCustomerGroup] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.CustomerGroupCreated, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerGroupCreationSubmittedSuccessfully, nil)
}

// Update Customer Group
//
//	@Summary		Update Customer Group
//	@Description	Updates an existing customer group segment request
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Segment ID"
//	@Param			body	body		customer_group_dto.UpdateSegmentRequest	true	"Customer Group DTO"
//	@Success		200		{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups/{id}/update [patch]
func (h *CustomerGroupAdapter) UpdateCustomerGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "UpdateCustomerGroup", "handler", "customerGroup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[UpdateCustomerGroup] extractID: %v", err)
		return
	}

	var req customer_group_dto.UpdateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[UpdateCustomerGroup] decode err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		log.Errorf("[UpdateCustomerGroup] validation err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	req.Normalize()

	if err := h.svc.Update(ctx, id, req); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateCustomerGroup] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.CustomerGroupUpdated, nil)
		return
	}
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.CustomerGroupUpdateSubmittedSuccessfully, nil)
}

// List Customer Groups
//
//	@Summary		List Customer Groups
//	@Description	Retrieves a paginated list of customer group segments
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200		{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups [get]
func (h *CustomerGroupAdapter) GetAllCustomerGroups(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.svc.FindAllWithPagination(r.Context(), filterParams)
	if err != nil {
		log.Errorf("[GetAllCustomerGroups] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerGroupFetchedSuccessfully, result)
}

// Get Customer Group by ID
//
//	@Summary		Get Customer Group by ID
//	@Description	Retrieves a customer group segment by ID
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Segment ID"
//	@Success		200	{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups/{id} [get]
func (h *CustomerGroupAdapter) GetCustomerGroup(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[GetCustomerGroup] extractID: %v", err)
		return
	}

	result, err := h.svc.FindByID(r.Context(), id)
	if err != nil {
		log.Errorf("[GetCustomerGroup] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerGroupFetchedSuccessfully, result)
}

// Delete Customer Group
//
//	@Summary		Delete Customer Group
//	@Description	Deletes a customer group segment by ID
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Segment ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups/{id} [delete]
func (h *CustomerGroupAdapter) DeleteCustomerGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "DeleteCustomerGroup", "handler", "customerGroup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[DeleteCustomerGroup] extractID: %v", err)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteCustomerGroup] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.CustomerGroupDeletedSuccessfully, nil)
		return
	}
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.CustomerGroupDeleteRequestSubmittedSuccessfully, nil)
}

// Enable Customer Group
//
//	@Summary		Enable Customer Group
//	@Description	Enables a customer group segment by ID
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Segment ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups/{id}/enable [patch]
func (h *CustomerGroupAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "Enable", "handler", "customerGroup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[Enable] extractID: %v", err)
		return
	}

	if err := h.svc.EnableOrDisable(ctx, id, true); err != nil {
		span.RecordError(err)
		log.Errorf("[Enable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.CustomerGroupEnabledSuccessfully, nil)
		return
	}
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.CustomerGroupEnableSubmittedSuccessfully, nil)
}

// Disable Customer Group
//
//	@Summary		Disable Customer Group
//	@Description	Disables a customer group segment by ID
//	@Tags			CustomerGroup
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Segment ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-groups/{id}/disable [patch]
func (h *CustomerGroupAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "Disable", "handler", "customerGroup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[Disable] extractID: %v", err)
		return
	}

	if err := h.svc.EnableOrDisable(ctx, id, false); err != nil {
		span.RecordError(err)
		log.Errorf("[Disable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.CustomerGroupDisabledSuccessfully, nil)
		return
	}
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.CustomerGroupDisableSubmittedSuccessfully, nil)
}
