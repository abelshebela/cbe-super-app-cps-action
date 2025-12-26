package bpsuser

type BPSUserCreatePayload struct {
	UserCode    string `json:"user_code" bson:"user_code"`
	FullName    string `json:"full_name" bson:"full_name"`
	JobTitle    string `json:"job_title" bson:"job_title"`
	ImpowerID   string `json:"impower_id" bson:"impower_id"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Email       string `json:"email" bson:"email"`
}
