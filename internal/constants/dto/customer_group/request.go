package customer_group

type CreateSegmentRequest struct {
	CustomerGroup           string `json:"customer_group"`
	CustomerGroupLabel      string `json:"customer_group_label"`
	CustomerSegment         string `json:"customer_segment"`
	CustomerSegmentLabel    string `json:"customer_segment_label"`
	CustomerSubsegment      string `json:"customer_subsegment"`
	CustomerSubsegmentLabel string `json:"customer_subsegment_label"`
	SuperappRole            string `json:"superapp_role"`
	SuperappRoleLabel       string `json:"superapp_role_label"`
}

type UpdateSegmentRequest struct {
	CustomerGroup           string `json:"customer_group"`
	CustomerGroupLabel      string `json:"customer_group_label"`
	CustomerSegment         string `json:"customer_segment"`
	CustomerSegmentLabel    string `json:"customer_segment_label"`
	CustomerSubsegment      string `json:"customer_subsegment"`
	CustomerSubsegmentLabel string `json:"customer_subsegment_label"`
	SuperappRole            string `json:"superapp_role"`
	SuperappRoleLabel       string `json:"superapp_role_label"`
}
