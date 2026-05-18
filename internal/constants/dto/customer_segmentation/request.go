package customersegmentation

import imodel "cbe-super-app-cps-action/internal/constants/model"

type CreateCustomerSegmentationRequest struct {
	CustomerRole string                `json:"customer_role"`
	Customer     []imodel.CustomerEntry `json:"customer"`
}

type UpdateCustomerSegmentationRequest struct {
	Customer                  []imodel.CustomerEntry `json:"customer,omitempty"`
	RemovedCustomerSegmentIds []string               `json:"removed_customer_segment_ids,omitempty"`
}
