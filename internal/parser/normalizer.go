package parser

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

	case int:
		return float64(v)

	case int64:
		return float64(v)

	case uint64:
		return float64(v)

	default:
		return value
	}

	return value
}
