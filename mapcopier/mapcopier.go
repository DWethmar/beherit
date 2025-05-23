package mapcopier

// Copy makes a deep copy of a map[string]any, including nested slices and maps.
func Copy(input map[string]any) map[string]any {
	copy := make(map[string]any, len(input))

	for key, value := range input {
		copy[key] = deepCopyValue(value)
	}

	return copy
}

func deepCopyValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return Copy(v)
	case []any:
		return deepCopySlice(v)
	default:
		return v
	}
}

func deepCopySlice(slice []any) []any {
	result := make([]any, len(slice))
	for i, val := range slice {
		result[i] = deepCopyValue(val)
	}
	return result
}
