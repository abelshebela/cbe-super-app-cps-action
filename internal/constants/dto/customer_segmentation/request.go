package customersegmentation

type CustomerSubSegments struct {
	Name            string `json:"name"`
	CustomerSegment string `json:"cust_segment"`
	CustomerGroup   string `json:"cust_group"`
}

type CreateCustomerSegmentationRequest struct {
	CustomerRole        string                `json:"customer_role"`
	CustomerSubSegments []CustomerSubSegments `json:"t24_customer_sub_segments"`
}

type UpdateCustomerSegmentationRequest struct {
	OldName             []string              `json:"old_name"`
	CustomerSubSegments []CustomerSubSegments `json:"t24_customer_sub_segments,omitempty"`
}
