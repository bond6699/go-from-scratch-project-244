package parser

import (
	"errors"
	"fmt"
	"math"
)

var errUnsupportedKeyType = errors.New("only string keys are supported")

func Normalize(data map[string]any) (map[string]any, error) {
	for key := range data {
		result, err := NormalizeValue(data[key])
		if err != nil {
			return nil, err
		}

		data[key] = result
	}

	return data, nil
}

func NormalizeValue(value any) (any, error) {
	switch v := value.(type) {
	case map[string]any:
		return Normalize(v)

	case map[any]any:
		return nil, fmt.Errorf("value %v not supported: %w", v, errUnsupportedKeyType)

	case []any:
		for i := range v {
			result, err := NormalizeValue(v[i])
			if err != nil {
				return nil, err
			}

			v[i] = result
		}

		return v, nil

	case int, int32, int64, uint64:
		return normalizeInteger(v), nil

	case float32, float64:
		return normalizeFloat(v), nil

	default:
		return value, nil
	}
}

func normalizeInteger(value any) any {
	switch v := value.(type) {
	case int:
		return int64(v)

	case int32:
		return int64(v)

	case int64:
		return v

	case uint64:
		if v <= math.MaxInt64 {
			return int64(v)
		}

		return v

	default:
		return value
	}
}

func normalizeFloat(value any) any {
	switch v := value.(type) {
	case float32:
		return normalizeFloat(float64(v))

	case float64:
		if v == math.Trunc(v) && !math.IsInf(v, 0) && !math.IsNaN(v) &&
			v >= math.MinInt64 && v <= math.MaxInt64 {
			return int64(v)
		}

		return v

	default:
		return value
	}
}
