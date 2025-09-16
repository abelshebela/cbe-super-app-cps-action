package utils

import (
	"fmt"
	"strconv"
)

// NumericConverter provides methods to convert interface{} values to various numeric types
type NumericConverter struct{}

// NewNumericConverter creates a new NumericConverter instance
func NewNumericConverter() *NumericConverter {
	return &NumericConverter{}
}

// ToUint64 converts interface{} to uint64 with validation
func (nc *NumericConverter) ToUint64(value interface{}, fieldName string) (uint64, error) {
	var result uint64

	switch v := value.(type) {
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid number", fieldName)
		}
		result = parsed
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("%s cannot be negative", fieldName)
		}
		result = uint64(v)
	case int:
		if v < 0 {
			return 0, fmt.Errorf("%s cannot be negative", fieldName)
		}
		result = uint64(v)
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("%s cannot be negative", fieldName)
		}
		result = uint64(v)
	case uint:
		result = uint64(v)
	case uint64:
		result = v
	case uint32:
		result = uint64(v)
	case uint16:
		result = uint64(v)
	case uint8:
		result = uint64(v)
	default:
		return 0, fmt.Errorf("%s must be a valid number, got %T", fieldName, value)
	}

	return result, nil
}

// ToUint converts interface{} to uint with validation
func (nc *NumericConverter) ToUint(value interface{}, fieldName string) (uint, error) {
	uint64Val, err := nc.ToUint64(value, fieldName)
	if err != nil {
		return 0, err
	}
	return uint(uint64Val), nil
}

// ToInt converts interface{} to int with validation
func (nc *NumericConverter) ToInt(value interface{}, fieldName string) (int, error) {
	var result int

	switch v := value.(type) {
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid number", fieldName)
		}
		result = parsed
	case float64:
		result = int(v)
	case int:
		result = v
	case int64:
		result = int(v)
	case uint:
		result = int(v)
	case uint64:
		result = int(v)
	case uint32:
		result = int(v)
	case uint16:
		result = int(v)
	case uint8:
		result = int(v)
	default:
		return 0, fmt.Errorf("%s must be a valid number, got %T", fieldName, value)
	}

	return result, nil
}

// ToInt32 converts interface{} to int32 with validation
func (nc *NumericConverter) ToInt32(value interface{}, fieldName string) (int32, error) {
	var result int32

	switch v := value.(type) {
	case string:
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid number", fieldName)
		}
		result = int32(parsed)
	case float64:
		result = int32(v)
	case int:
		result = int32(v)
	case int64:
		result = int32(v)
	case uint:
		result = int32(v)
	case uint64:
		result = int32(v)
	case uint32:
		result = int32(v)
	case uint16:
		result = int32(v)
	case uint8:
		result = int32(v)
	case int32:
		result = v
	case int16:
		result = int32(v)
	case int8:
		result = int32(v)
	default:
		return 0, fmt.Errorf("%s must be a valid number, got %T", fieldName, value)
	}

	return result, nil
}

// ToInt64 converts interface{} to int64 with validation
func (nc *NumericConverter) ToInt64(value interface{}, fieldName string) (int64, error) {
	var result int64

	switch v := value.(type) {
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid number", fieldName)
		}
		result = parsed
	case float64:
		result = int64(v)
	case int:
		result = int64(v)
	case int64:
		result = v
	case uint:
		result = int64(v)
	case uint64:
		result = int64(v)
	case uint32:
		result = int64(v)
	case uint16:
		result = int64(v)
	case uint8:
		result = int64(v)
	default:
		return 0, fmt.Errorf("%s must be a valid number, got %T", fieldName, value)
	}

	return result, nil
}

// ToFloat64 converts interface{} to float64 with validation
func (nc *NumericConverter) ToFloat64(value interface{}, fieldName string) (float64, error) {
	var result float64

	switch v := value.(type) {
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid number", fieldName)
		}
		result = parsed
	case float64:
		result = v
	case float32:
		result = float64(v)
	case int:
		result = float64(v)
	case int64:
		result = float64(v)
	case uint:
		result = float64(v)
	case uint64:
		result = float64(v)
	case uint32:
		result = float64(v)
	case uint16:
		result = float64(v)
	case uint8:
		result = float64(v)
	default:
		return 0, fmt.Errorf("%s must be a valid number, got %T", fieldName, value)
	}

	return result, nil
}

// ValidateUint64Range validates uint64 value within a range
func (nc *NumericConverter) ValidateUint64Range(value interface{}, fieldName string, min, max uint64) (uint64, error) {
	result, err := nc.ToUint64(value, fieldName)
	if err != nil {
		return 0, err
	}

	if result < min {
		return 0, fmt.Errorf("%s must be at least %d", fieldName, min)
	}

	if result > max {
		return 0, fmt.Errorf("%s cannot exceed %d", fieldName, max)
	}

	return result, nil
}

// ValidateUintRange validates uint value within a range
func (nc *NumericConverter) ValidateUintRange(value interface{}, fieldName string, min, max uint) (uint, error) {
	result, err := nc.ToUint(value, fieldName)
	if err != nil {
		return 0, err
	}

	if result < min {
		return 0, fmt.Errorf("%s must be at least %d", fieldName, min)
	}

	if result > max {
		return 0, fmt.Errorf("%s cannot exceed %d", fieldName, max)
	}

	return result, nil
}

// ValidateIntRange validates int value within a range
func (nc *NumericConverter) ValidateIntRange(value interface{}, fieldName string, min, max int) (int, error) {
	result, err := nc.ToInt(value, fieldName)
	if err != nil {
		return 0, err
	}

	if result < min {
		return 0, fmt.Errorf("%s must be at least %d", fieldName, min)
	}

	if result > max {
		return 0, fmt.Errorf("%s cannot exceed %d", fieldName, max)
	}

	return result, nil
}

// ValidateFloat64Range validates float64 value within a range
func (nc *NumericConverter) ValidateFloat64Range(value interface{}, fieldName string, min, max float64) (float64, error) {
	result, err := nc.ToFloat64(value, fieldName)
	if err != nil {
		return 0, err
	}

	if result < min {
		return 0, fmt.Errorf("%s must be at least %f", fieldName, min)
	}

	if result > max {
		return 0, fmt.Errorf("%s cannot exceed %f", fieldName, max)
	}

	return result, nil
}

// IsPositiveNumber checks if the value is a positive number
func (nc *NumericConverter) IsPositiveNumber(value interface{}, fieldName string) (uint64, error) {
	return nc.ValidateUint64Range(value, fieldName, 1, ^uint64(0))
}

// IsNonNegativeNumber checks if the value is a non-negative number
func (nc *NumericConverter) IsNonNegativeNumber(value interface{}, fieldName string) (uint64, error) {
	return nc.ValidateUint64Range(value, fieldName, 0, ^uint64(0))
}

// IsPositiveInt checks if the value is a positive integer
func (nc *NumericConverter) IsPositiveInt(value interface{}, fieldName string) (int, error) {
	return nc.ValidateIntRange(value, fieldName, 1, ^int(0))
}

// IsNonNegativeInt checks if the value is a non-negative integer
func (nc *NumericConverter) IsNonNegativeInt(value interface{}, fieldName string) (int, error) {
	return nc.ValidateIntRange(value, fieldName, 0, ^int(0))
}

// IsPositiveFloat checks if the value is a positive float
func (nc *NumericConverter) IsPositiveFloat(value interface{}, fieldName string) (float64, error) {
	return nc.ValidateFloat64Range(value, fieldName, 0.000001, 1e308)
}

// IsNonNegativeFloat checks if the value is a non-negative float
func (nc *NumericConverter) IsNonNegativeFloat(value interface{}, fieldName string) (float64, error) {
	return nc.ValidateFloat64Range(value, fieldName, 0, 1e308)
}
