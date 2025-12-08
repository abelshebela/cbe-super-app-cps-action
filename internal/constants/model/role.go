package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID          bson.ObjectID `json:"_id" bson:"_id"`
	JobTitle    string        `json:"job_title" bson:"job_title"`
	Role        string        `json:"role" bson:"role"`
	BranchGrade string        `json:"branch_grade" bson:"branch_grade"`
	Department  string        `json:"department" bson:"department"`
	Position    string        `json:"position" bson:"position"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}
