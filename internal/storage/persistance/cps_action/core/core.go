package cps_action_core

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func SanitizePipeline(exclude []string) bson.D {

	values := bson.A{}
	for i, _ := range exclude {
		values = append(values, fmt.Sprintf("previous_action.%s", exclude[i]))
		values = append(values, fmt.Sprintf("current_action.%s", exclude[i]))
	}
	return bson.D{{
		Key:   "$unset",
		Value: values,
	}}
}
