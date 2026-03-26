package customersegmentation

import (
	"cbe-super-app-cps-action/internal/constants"
	cust_seg "cbe-super-app-cps-action/internal/constants/dto/customer_segmentation"
	seg "cbe-super-app-cps-action/internal/constants/interfaces/customer_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CustomerSegmentationAdapter struct {
	svc    service.CustomerSegmentationService
	logger utils.Logger
}

func NewCustomerSegmentation(svc service.CustomerSegmentationService, logger utils.Logger) seg.CustomerSegmentation {
	return &CustomerSegmentationAdapter{svc: svc, logger: logger}
}

// Create Customer Segmentation
//
//	@Summary		Create Customer Segmentation
//	@Description	Creates a new customer segmentation request
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		customersegmentation.CreateCustomerSegmentationRequest	true	"Customer Segmentation DTO"
//	@Success		201		{object}	localization.StandardResponse{data=nil}	"Customer segmentation creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		401		{object}	localization.StandardResponse{data=nil}	"Unauthorized"
//	@Failure		422		{object}	localization.StandardResponse{data=nil}	"Validation error"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Router			/customer-segmentations [post]
func (c *CustomerSegmentationAdapter) CreateCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "CreateCustomerSegmentation", "handler", "customerSegmentation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req cust_seg.CreateCustomerSegmentationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[CreateCustomerSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[CreateCustomerSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, fmt.Sprintf("validation error: %v", err))
		return
	}
	req.CapitilizeCustomerSegmentationRequest()

	if err := c.svc.Create(ctx, req); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[CreateCustomerSegmentation] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationCreated, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationCreationSubmittedSuccessfully, nil)
}

// Update Customer Segmentation
//
//	@Summary		Update Customer Segmentation
//	@Description	Updates an existing customer segmentation request
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string										true	"Segmentation ID"
//	@Param			body	body		customersegmentation.UpdateCustomerSegmentationRequest	true	"Customer Segmentation DTO"
//	@Success		200		{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations/{id}/update [patch]
func (c *CustomerSegmentationAdapter) UpdateCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "UpdateCustomerSegmentation", "handler", "customerSegmentation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateCustomerSegmentation] extractID: %v", err)
		return
	}

	var req cust_seg.UpdateCustomerSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[UpdateCustomerSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[UpdateCustomerSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.Update(ctx, id, req); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[UpdateCustomerSegmentation] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationUpdated, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationUpdateSubmittedSuccessfully, nil)
}

// List Customer Segmentations
//
//	@Summary		List Customer Segmentations
//	@Description	Retrieves a paginated list of customer segmentations
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Customer segmentations retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		401			{object}	localization.StandardResponse{data=nil}		"Unauthorized"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Router			/customer-segmentations [get]
func (c *CustomerSegmentationAdapter) GetAllCustomerSegmentations(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), c.logger)
	filterParams := local_util.ExtractFilterParams(r)

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

	segs, err := c.svc.FindAllWithPagination(r.Context(), filterParams)
	if err != nil {
		log.Errorf("[GetAllCustomerSegmentations] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationFetchedSuccessfully, segs)
}

// Get Customer Segmentation by ID
//
//	@Summary		Get Customer Segmentation by ID
//	@Description	Retrieves a customer segmentation by ID
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Segmentation ID"
//	@Success		200				{object}	localization.StandardResponse{data=object}	"Customer segmentation retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		401				{object}	localization.StandardResponse{data=nil}		"Unauthorized"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}		"Customer segmentation not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Router			/customer-segmentations/{id} [get]
func (c *CustomerSegmentationAdapter) GetCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), c.logger)
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[GetCustomerSegmentation] extractID: %v", err)
		return
	}

	seg, err := c.svc.FindById(r.Context(), id)
	if err != nil {
		log.Errorf("[GetCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationFetchedSuccessfully, seg)
}

// Delete Customer Segmentation
//
//	@Summary		Delete Customer Segmentation
//	@Description	Deletes a customer segmentation by ID
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Segmentation ID"
//	@Success		200				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations/{id} [delete]
func (c *CustomerSegmentationAdapter) DeleteCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "DeleteCustomerSegmentation", "handler", "customerSegmentation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[DeleteCustomerSegmentation] extractID: %v", err)
		return
	}

	if err := c.svc.Delete(ctx, id); err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[DeleteCustomerSegmentation] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationDeleteddSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationDeleteddSuccessfully, nil)
}

// Enable Customer Segmentation
//
//	@Summary		Enable Customer Segmentation
//	@Description	Enables a customer segmentation by ID
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Segmentation ID"
//	@Success		200				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations/{id}/enable [patch]
func (c *CustomerSegmentationAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "Enable", "handler", "customerSegmentation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[Enable] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, true); err != nil {
		span.RecordError(err)
		log.Errorf("[Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Enable] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationEnableSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationEnableSuccessfully, nil)
}

// Disable Customer Segmentation
//
//	@Summary		Disable Customer Segmentation
//	@Description	Disables a customer segmentation by ID
//	@Tags			CustomerSegmentation
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Segmentation ID"
//	@Success		200				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations/{id}/disable [patch]
func (c *CustomerSegmentationAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "Disable", "handler", "customerSegmentation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[Disable] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, false); err != nil {
		span.RecordError(err)
		log.Errorf("[Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Disable] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationDisableSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationDisableSuccessfully, nil)
}
