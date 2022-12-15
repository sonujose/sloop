package template

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))

	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {

		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}

		out[k] = v
	}
	return out
}

func generateValuesForOverrides(keys []string, val string) map[string]interface{} {

	finalElement := map[string]interface{}{
		keys[len(keys)-1]: val,
	}

	for i := len(keys) - 2; i > 0; i-- {
		newElement := map[string]interface{}{
			keys[i]: finalElement,
		}
		finalElement = newElement
	}

	return finalElement
}
