package types

import (
	"encoding/json"
	"reflect"
)

type PaginatedResponse[T any] struct {
	Data T              `json:"docs"`
	Meta PaginationMeta `json:"meta"`
}

func (p PaginatedResponse[T]) MarshalJSON() ([]byte, error) {
	var data any = p.Data

	v := reflect.ValueOf(p.Data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		// Replace nil slice with empty slice of the same type
		data = reflect.MakeSlice(v.Type(), 0, 0).Interface()
	}

	return json.Marshal(&struct {
		Data any            `json:"docs"`
		Meta PaginationMeta `json:"meta"`
	}{
		Data: data,
		Meta: p.Meta,
	})
}

type PaginatedResponseForFeedback[T any] struct {
	Data T              `json:"docs"`
	Meta PaginationMeta `json:"meta"`
	AverageRatings map[string]float64        `json:"average_ratings,omitempty" bson:"average_ratings,omitempty"`
}