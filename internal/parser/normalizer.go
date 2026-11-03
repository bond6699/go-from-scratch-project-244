package parser

import "math"

func Normalize(data map[string]any) map[string]any {
	for key := range data {
		data[key] = NormalizeValue(data[key])
	}

	return data
}

func NormalizeValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return Normalize(v)

	case []any:
		for i := range v {
			v[i] = NormalizeValue(v[i])
		}

		return v

	case int, int32, int64, uint64:
		return normalizeInteger(v)

	case float32, float64:
		return normalizeFloat(v)

	default:
		return value
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
