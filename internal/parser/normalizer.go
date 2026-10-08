package parser

func Normalize(data any) any {
	switch v := data.(type) {
	case map[string]any:
		for key := range v {
			v[key] = Normalize(v[key])
		}
	case []any:
		for i := range v {
			v[i] = Normalize(v[i])
		}
	case nil:
		return nil

	case bool:
		return bool(v)

	case string:
		return string(v)

	case int:
		return float64(v)

	case int64:
		return float64(v)

	case uint64:
		return float64(v)

	case float64:
		return float64(v)

	}

	return data
}
