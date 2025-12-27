package customersegmentation

import (
	cust_seg "cbe-super-app-cps-action/internal/constants/dto/customer_segmentation"
	seg "cbe-super-app-cps-action/internal/constants/interfaces/customer_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

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
//	@Param			body	body		cust_seg.CreateCustomerSegmentationRequest	true	"Customer Segmentation DTO"
//	@Success		201		{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations [post]
func (c *CustomerSegmentationAdapter) CreateCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	var reqs []cust_seg.CreateCustomerSegmentationRequest

	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		c.logger.Errorf("[CreateCustomerSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}

	for i, req := range reqs {
		if err := req.Validate(); err != nil {
			c.logger.Errorf("[CreateCustomerSegmentation] validation error at index %d: %v", i, err)
			localization.SendErrorByCodeResponse(w, fmt.Sprintf("validation error at index %d: %v", i, err))
			return
		}
	}

	if err := c.svc.CreateBulk(r.Context(), reqs); err != nil {
		c.logger.Errorf("[CreateCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
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
//	@Param			body	body		cust_seg.UpdateCustomerSegmentationRequest	true	"Customer Segmentation DTO"
//	@Success		200		{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations/{id}/update [patch]
func (c *CustomerSegmentationAdapter) UpdateCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	_, span := local_util.TraceLogger(r.Context(), "handler", "UpdateCustomerSegmentation", "handler", "customerSegmentation")
	defer span.End()

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[UpdateCustomerSegmentation] extractID: %v", err)
		return
	}

	var req cust_seg.UpdateCustomerSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorf("[UpdateCustomerSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		c.logger.Errorf("[UpdateCustomerSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.Update(r.Context(), id, req); err != nil {
		c.logger.Errorf("[UpdateCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
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
//	@Success		200			{object}	localization.StandardResponse{data=types.PaginatedResponse[[]*imodel.CustomerSegmentation]}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations [get]
func (c *CustomerSegmentationAdapter) GetAllCustomerSegmentations(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	segs, err := c.svc.FindAllWithPagination(r.Context(), filterParams)
	if err != nil {
		c.logger.Errorf("[GetAllCustomerSegmentations] service error: %v", err)
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
//	@Success		200				{object}	localization.StandardResponse{data=imodel.CustomerSegmentation}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/customer-segmentations/{id} [get]
func (c *CustomerSegmentationAdapter) GetCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[GetCustomerSegmentation] extractID: %v", err)
		return
	}

	seg, err := c.svc.FindById(r.Context(), id)
	if err != nil {
		c.logger.Errorf("[GetCustomerSegmentation] service error: %v", err)
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
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[DeleteCustomerSegmentation] extractID: %v", err)
		return
	}

	if err := c.svc.Delete(r.Context(), id); err != nil {
		c.logger.Errorf("[DeleteCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerSegmentationDeleteddSuccessfully, nil)
}
