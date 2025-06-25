package utils

import (
	"reflect"
)

type ValidationResult struct {
	Success     bool
	Data        interface{}
	Field       interface{}
	FieldErrors map[string][]string
	Formatted   map[string][]string
}

// Validator checks for "required" struct tags.
func Validator(field interface{}, schema interface{}) ValidationResult {
	val := reflect.ValueOf(schema)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	fieldErrors := make(map[string][]string)
	formatted := make(map[string][]string)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "required" {
			zero := reflect.Zero(fieldVal.Type()).Interface()
			if reflect.DeepEqual(fieldVal.Interface(), zero) {
				msg := fieldType.Name + " is required"
				fieldErrors[fieldType.Name] = append(fieldErrors[fieldType.Name], msg)
				formatted[fieldType.Name] = append(formatted[fieldType.Name], msg)
			}
		}
	}

	if len(fieldErrors) > 0 {
		return ValidationResult{
			Success:     false,
			Data:        nil,
			Field:       field,
			FieldErrors: fieldErrors,
			Formatted:   formatted,
		}
	}

	return ValidationResult{
		Success:     true,
		Data:        schema,
		Field:       field,
		FieldErrors: nil,
		Formatted:   nil,
	}
}
