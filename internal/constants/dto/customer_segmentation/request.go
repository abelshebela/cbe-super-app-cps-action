package customersegmentation

type CreateCustomerSegmentationRequest struct {
	CustomerRole       string `json:"customer_role"`
	CustomerSegment    string `json:"customer_segment"`
	CustomerSubSegment string `json:"customer_sub_segment"`
	CustomerGroup      string `json:"customer_group"`
}

type UpdateCustomerSegmentationRequest struct {
	CustomerRole       *string `json:"customer_role,omitempty"`
	CustomerSegment    *string `json:"customer_segment,omitempty"`
	CustomerSubSegment *string `json:"customer_sub_segment,omitempty"`
	CustomerGroup      *string `json:"customer_group,omitempty"`
}
