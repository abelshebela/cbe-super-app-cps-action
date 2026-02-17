package model

import "go.mongodb.org/mongo-driver/v2/bson"

type ServiceList struct {
	ID          bson.ObjectID `json:"_id" bson:"_id"`
	ServiceName string        `json:"service_name" bson:"service_name"`
	ServiceKey  string        `json:"service_key" bson:"service_key"`
}
