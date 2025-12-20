package customersegmentation

import (
	cust_seg "cbe-super-app-cps-action/internal/constants/dto/customer_segmentation"
	seg "cbe-super-app-cps-action/internal/constants/interfaces/customer_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CustomerSegmentationAdapter struct {
	svc    service.CustomerSegmentationService
	logger utils.Logger
}

func NewCustomerSegmentation(svc service.CustomerSegmentationService, logger utils.Logger) seg.CustomerSegmentation {
	return &CustomerSegmentationAdapter{svc: svc, logger: logger}
}

func (c *CustomerSegmentationAdapter) CreateCustomerSegmentation(w http.ResponseWriter, r *http.Request) {
	var req cust_seg.CreateCustomerSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorf("[CreateCustomerSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		c.logger.Errorf("[CreateCustomerSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := c.svc.Create(r.Context(), req); err != nil {
		c.logger.Errorf("[CreateCustomerSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerSegmentationCreationSubmittedSuccessfully, nil)
}
