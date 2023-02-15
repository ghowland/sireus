package util

import (
	"errors"
	"math"
)

var (
	// Float64 is our primary data type, custom error for tracking problems
	InvalidTypeFloat64 = errors.New("Value could not be converted to Float64")
)

// Convert any value we can into a float64 in a predictable manner
func ConvertInterfaceToFloat(value interface{}) (float64, error) {
	switch value := value.(type) {
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint:
		return float64(value), nil
	default:
		return math.NaN(), InvalidTypeFloat64
	}
}
