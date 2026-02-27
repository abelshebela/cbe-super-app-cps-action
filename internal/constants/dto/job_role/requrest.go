package roles

type RequestRolesCreate struct {
	JobTitle string `json:"job_title" bson:"job_title"`
	Role     string `json:"role" bson:"role"`
}

type RequestRolesUpdate struct {
	JobTitle string `json:"job_title" bson:"job_title"`
	Role     string `json:"role" bson:"role"`
}
