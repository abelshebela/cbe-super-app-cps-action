package roles

type RequestJobRolesCreate struct {
	Code     string `json:"code" bson:"code"`
	JobTitle string `json:"job_title" bson:"job_title"`
	Role     string `json:"role" bson:"role"`
}

type RequestJobRolesUpdate struct {
	Code     string `json:"code,omitempty" bson:"code"`
	JobTitle string `json:"job_title" bson:"job_title"`
	Role     string `json:"role" bson:"role"`
}
