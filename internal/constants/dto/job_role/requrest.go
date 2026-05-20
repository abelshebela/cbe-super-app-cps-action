package roles

type RequestJobRolesCreate struct {
	JobTitle string `json:"job_title" bson:"job_title"`
	Role     string `json:"role" bson:"role"`
}

type RequestJobRolesUpdate struct {
	JobTitle string `json:"job_title" bson:"job_title"`
	Role     string `json:"role" bson:"role"`
}
