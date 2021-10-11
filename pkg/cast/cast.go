package cast

func ToMapStringArray(src map[interface{}]interface{}) map[string][]string {
	result := make(map[string][]string, 0)
	for k, v := range src {
		key := k.(string)
		values := v.([]interface{})
		result[key] = make([]string, 0)
		for _, val := range values {
			result[key] = append(result[key], val.(string))
		}
	}
	return result
}

func ToMapString(src map[interface{}]interface{}) map[string]string {
	result := make(map[string]string, 0)
	for key, v := range src {
		result[key.(string)] = v.(string)
	}
	return result
}
