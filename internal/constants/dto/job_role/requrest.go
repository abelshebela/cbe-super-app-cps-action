package roles

type RequestRolesCreate struct {
	Code        string `json:"code" bson:"code"`
	JobTitle    string `json:"job_title" bson:"job_title"`
	Description string `json:"description" bson:"description"`
	Role        string `json:"role" bson:"role"`
}

type RequestRolesUpdate struct {
	Code        string `json:"code" bson:"code"`
	JobTitle    string `json:"job_title" bson:"job_title"`
	Description string `json:"description" bson:"description"`
	Role        string `json:"role" bson:"role"`
}
