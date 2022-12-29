package template

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
